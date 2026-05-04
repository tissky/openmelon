package generation

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestShellProvider_success(t *testing.T) {
	p := &ShellProvider{Command: "echo hello", Model: "test-model"}
	req := &Request{Prompt: "some prompt", ArtifactType: "image_prompt"}

	content, trace, err := p.Generate(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content != "hello" {
		t.Errorf("content = %q, want %q", content, "hello")
	}
	if trace == nil {
		t.Fatal("trace is nil")
	}
	if trace.ProviderType != "shell" {
		t.Errorf("trace.ProviderType = %q, want %q", trace.ProviderType, "shell")
	}
	if trace.DurationSec < 0 {
		t.Errorf("trace.DurationSec = %f, want >= 0", trace.DurationSec)
	}
}

func TestShellProvider_nonZeroExit(t *testing.T) {
	// "false" command always exits with code 1
	p := &ShellProvider{Command: "false"}
	req := &Request{Prompt: "prompt"}

	_, _, err := p.Generate(context.Background(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var provErr *ProviderError
	if !errors.As(err, &provErr) {
		t.Fatalf("expected *ProviderError, got %T: %v", err, err)
	}
	if provErr.Code != "non_zero_exit" {
		t.Errorf("Code = %q, want %q", provErr.Code, "non_zero_exit")
	}
}

func TestShellProvider_contextTimeout(t *testing.T) {
	// "sleep 10" will be killed by the context timeout
	p := &ShellProvider{Command: "sleep 10"}
	req := &Request{Prompt: "prompt"}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, _, err := p.Generate(ctx, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var provErr *ProviderError
	if !errors.As(err, &provErr) {
		t.Fatalf("expected *ProviderError, got %T: %v", err, err)
	}
	if provErr.Code != "timeout" {
		t.Errorf("Code = %q, want %q", provErr.Code, "timeout")
	}
	if !strings.Contains(provErr.Message, "timed out") {
		t.Errorf("expected 'timed out' in message, got: %q", provErr.Message)
	}
}
