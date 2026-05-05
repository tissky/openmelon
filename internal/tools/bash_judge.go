package tools

// bash_judge.go — wires a small classifier prompt against an
// llm.ToolCaller (or any client supporting Chat) so cmd_repl can plug
// the main agent LLM straight in as the bash safety judge.
//
// Returning BashAsk on any error keeps the system fail-safe: a network
// blip won't auto-approve a destructive command.

import (
	"context"
	"strings"

	"github.com/eight-acres-lab/openmelon/internal/llm"
)

// JudgeBashWithLLM is a ready-made implementation of Env.JudgeBash that
// asks the given client to classify the command. cmd_repl wires this
// in with the main LLM as the judge — fast enough for normal use and
// avoids a second model dep.
//
// The classifier prompt is kept tight: 1-token output, no rationale,
// no markdown. We ignore everything past the first non-whitespace word.
func JudgeBashWithLLM(client llm.ToolCaller) func(context.Context, string, string) BashJudgement {
	return func(ctx context.Context, command, description string) BashJudgement {
		// Use the simpler Complete API if the client supports it, but
		// ToolCaller.Chat works too — we just send no tools.
		req := llm.ChatRequest{
			Messages: []llm.Message{
				{Role: llm.RoleSystem, Content: bashJudgeSystemPrompt},
				{Role: llm.RoleUser, Content: "Command: " + command + "\nDescription: " + description},
			},
			Temperature: 0.0,
			MaxTokens:   8,
		}
		resp, err := client.Chat(ctx, req)
		if err != nil {
			return BashAsk
		}
		switch firstWord(resp.Message.Content) {
		case "AUTO":
			return BashAuto
		case "BLOCK":
			return BashBlock
		default:
			return BashAsk
		}
	}
}

func firstWord(s string) string {
	for _, f := range strings.Fields(strings.ToUpper(strings.TrimSpace(s))) {
		// Strip surrounding punctuation the model might add.
		f = strings.Trim(f, "`*_.,;:!?\"'")
		if f != "" {
			return f
		}
	}
	return ""
}

const bashJudgeSystemPrompt = `You are a safety classifier for an AI agent's bash tool. Given a shell command, output EXACTLY one of:

AUTO   Read-only inspection commands. Examples: ls, cat, file, head, tail, du, df, wc, grep, find (without -delete/-exec), stat, identify, exiftool, magick identify, pdfinfo, ffprobe, jq (without -i), xmllint, sha256sum, md5sum, base64, date, uname, hostname, uptime, ps, lsof, which, pwd, env, command -v, open (macOS), xdg-open. Anything that reads but doesn't write or call out.
ASK    Writes inside the project workdir, normal-but-side-effecting commands during creative work. Examples: mkdir, touch, mv, cp, npm install, git commit, ImageMagick convert, ffmpeg encode, sed -i, python script.py, curl GET (read-only), pip install.
BLOCK  Destructive or exfiltrating. Examples: rm -rf, dd, mkfs, sudo, chmod 777 /, anything piping to /etc, modifying ~/.ssh, curl/wget POST/PUT to a non-localhost URL, scp/rsync to a remote, eval/exec of remote content, nc -l (network listener), iptables, anything that reads secrets and sends them out.

Respond with ONE WORD on a single line. No prose, no markdown, no explanation. If unsure, output ASK.
`
