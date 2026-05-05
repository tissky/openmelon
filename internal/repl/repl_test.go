package repl

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/eight-acres-lab/openmelon/internal/llm"
	"github.com/eight-acres-lab/openmelon/internal/projectx"
	"github.com/eight-acres-lab/openmelon/internal/runtime"
	"github.com/eight-acres-lab/openmelon/internal/tools"
)

// scriptedLLM returns a sequence of pre-recorded chat responses, one
// per Run() call. Used to drive the REPL through deterministic turns.
type scriptedLLM struct{ responses []llm.ChatResponse }

func (s *scriptedLLM) Chat(_ context.Context, _ llm.ChatRequest) (*llm.ChatResponse, error) {
	if len(s.responses) == 0 {
		return &llm.ChatResponse{
			Message:      llm.Message{Role: llm.RoleAssistant, Content: "(out of script)"},
			FinishReason: llm.FinishStop,
		}, nil
	}
	r := s.responses[0]
	s.responses = s.responses[1:]
	return &r, nil
}

func newProjectAt(t *testing.T) (string, *projectx.Project) {
	t.Helper()
	wd := t.TempDir()
	p, err := projectx.Init(wd, "test-proj", "Test")
	if err != nil {
		t.Fatalf("project init: %v", err)
	}
	return wd, p
}

func TestRunExitsOnSlashExit(t *testing.T) {
	wd, p := newProjectAt(t)
	reg := tools.NewRegistry()
	rt := &runtime.Runtime{LLM: &scriptedLLM{}, Registry: reg}

	in := strings.NewReader("/exit\n")
	var out, errOut bytes.Buffer
	err := Run(context.Background(), Options{
		Workdir: wd, Project: p, Runtime: rt,
		SystemPrompt: "be terse", SessionIntent: "test",
		In: in, Out: &out, Err: &errOut,
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !strings.Contains(out.String(), "bye") {
		t.Errorf("expected goodbye, got: %q", out.String())
	}
}

func TestRunExitsOnEOF(t *testing.T) {
	wd, p := newProjectAt(t)
	reg := tools.NewRegistry()
	rt := &runtime.Runtime{LLM: &scriptedLLM{}, Registry: reg}

	in := strings.NewReader("") // immediate EOF
	var out, errOut bytes.Buffer
	if err := Run(context.Background(), Options{
		Workdir: wd, Project: p, Runtime: rt,
		SystemPrompt: "x", SessionIntent: "test",
		In: in, Out: &out, Err: &errOut,
	}); err != nil {
		t.Fatalf("Run: %v", err)
	}
}

func TestRunSendsUserInputThroughRuntimeAndStreamsTextOut(t *testing.T) {
	wd, p := newProjectAt(t)
	reg := tools.NewRegistry()
	rt := &runtime.Runtime{
		LLM: &scriptedLLM{responses: []llm.ChatResponse{{
			Message:      llm.Message{Role: llm.RoleAssistant, Content: "hello back"},
			FinishReason: llm.FinishStop,
		}}},
		Registry: reg,
	}
	in := strings.NewReader("hi\n/exit\n")
	var out, errOut bytes.Buffer
	if err := Run(context.Background(), Options{
		Workdir: wd, Project: p, Runtime: rt,
		SystemPrompt: "x", SessionIntent: "test",
		In: in, Out: &out, Err: &errOut,
	}); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !strings.Contains(out.String(), "hello back") {
		t.Errorf("model reply not rendered: %q", out.String())
	}
}

func TestRunPersistsHistoryAcrossTurns(t *testing.T) {
	wd, p := newProjectAt(t)
	reg := tools.NewRegistry()

	// Track what each Chat call sees.
	var seen []int
	wrapping := &recordingLLM{
		inner: &scriptedLLM{responses: []llm.ChatResponse{
			{Message: llm.Message{Role: llm.RoleAssistant, Content: "a"}, FinishReason: llm.FinishStop},
			{Message: llm.Message{Role: llm.RoleAssistant, Content: "b"}, FinishReason: llm.FinishStop},
		}},
		recorder: &seen,
	}
	rt := &runtime.Runtime{LLM: wrapping, Registry: reg}

	in := strings.NewReader("first\nsecond\n/exit\n")
	var out, errOut bytes.Buffer
	if err := Run(context.Background(), Options{
		Workdir: wd, Project: p, Runtime: rt,
		SystemPrompt: "be terse", SessionIntent: "test",
		In: in, Out: &out, Err: &errOut,
	}); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(seen) != 2 {
		t.Fatalf("expected 2 chat calls, got %d", len(seen))
	}
	// Turn 1: system + user. Turn 2: system + user + assistant + user.
	if seen[0] != 2 {
		t.Errorf("first turn message count: %d (want 2)", seen[0])
	}
	if seen[1] != 4 {
		t.Errorf("second turn message count: %d (want 4)", seen[1])
	}
}

func TestSlashClearResetsHistory(t *testing.T) {
	wd, p := newProjectAt(t)
	reg := tools.NewRegistry()

	var seen []int
	wrapping := &recordingLLM{
		inner: &scriptedLLM{responses: []llm.ChatResponse{
			{Message: llm.Message{Role: llm.RoleAssistant, Content: "a"}, FinishReason: llm.FinishStop},
			{Message: llm.Message{Role: llm.RoleAssistant, Content: "b"}, FinishReason: llm.FinishStop},
		}},
		recorder: &seen,
	}
	rt := &runtime.Runtime{LLM: wrapping, Registry: reg}

	in := strings.NewReader("first\n/clear\nsecond\n/exit\n")
	var out, errOut bytes.Buffer
	if err := Run(context.Background(), Options{
		Workdir: wd, Project: p, Runtime: rt,
		SystemPrompt: "x", SessionIntent: "test",
		In: in, Out: &out, Err: &errOut,
	}); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(seen) != 2 {
		t.Fatalf("expected 2 chat calls, got %d", len(seen))
	}
	// After /clear, the second turn should NOT include prior history —
	// it sends the system prompt + user only, just like the first turn.
	if seen[1] != 2 {
		t.Errorf("after /clear, second turn should send 2 messages, got %d", seen[1])
	}
	if !strings.Contains(out.String(), "history cleared") {
		t.Errorf("expected /clear feedback, got: %q", out.String())
	}
}

func TestSlashHelpPrintsCommandList(t *testing.T) {
	wd, p := newProjectAt(t)
	reg := tools.NewRegistry()
	rt := &runtime.Runtime{LLM: &scriptedLLM{}, Registry: reg}

	in := strings.NewReader("/help\n/exit\n")
	var out, errOut bytes.Buffer
	if err := Run(context.Background(), Options{
		Workdir: wd, Project: p, Runtime: rt,
		SystemPrompt: "x", SessionIntent: "test",
		In: in, Out: &out, Err: &errOut,
	}); err != nil {
		t.Fatalf("Run: %v", err)
	}
	body := out.String()
	for _, want := range []string{"/clear", "/history", "/save", "/session", "/exit"} {
		if !strings.Contains(body, want) {
			t.Errorf("/help missing %q", want)
		}
	}
}

// recordingLLM wraps a ToolCaller and records the message-list length
// passed into each Chat call.
type recordingLLM struct {
	inner    llm.ToolCaller
	recorder *[]int
}

func (r *recordingLLM) Chat(ctx context.Context, req llm.ChatRequest) (*llm.ChatResponse, error) {
	*r.recorder = append(*r.recorder, len(req.Messages))
	return r.inner.Chat(ctx, req)
}

// satisfy json import in renderer (used by /save smoke test below)
var _ = json.RawMessage{}
