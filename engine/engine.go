// Copyright 2026 Point Eight AI Pte. Ltd.
// Licensed under the Apache License, Version 2.0

// Package engine provides the core Skill-Plus execution engine.
//
// The engine loads Skills from a registry, dispatches them based on
// content characteristics, runs them in sandboxed environments, and
// aggregates their outputs into B-face JSON.
package engine

import "context"

// BFace represents the Agent-facing metadata produced by running Skills.
type BFace struct {
	VisualDescription string          `json:"visual_description,omitempty"`
	Entities          []Entity        `json:"entities,omitempty"`
	Topics            []string        `json:"topics,omitempty"`
	Sentiment         *Sentiment      `json:"sentiment,omitempty"`
	RAGAnchors        []string        `json:"rag_anchors,omitempty"`
	AgentPrompts      []string        `json:"agent_prompts,omitempty"`
	Lang              string          `json:"lang,omitempty"`
	Safety            *SafetyResult   `json:"safety,omitempty"`
	SkillsApplied     []SkillRunInfo  `json:"skills_applied"`
}

// Entity represents a detected entity in content.
type Entity struct {
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	Confidence float64 `json:"confidence"`
}

// Sentiment represents sentiment analysis results.
type Sentiment struct {
	Valence float64 `json:"valence"`
	Arousal float64 `json:"arousal"`
}

// SafetyResult represents safety check results.
type SafetyResult struct {
	Flags []string `json:"flags"`
}

// SkillRunInfo records metadata about a single skill execution.
type SkillRunInfo struct {
	Skill      string `json:"skill"`
	Version    string `json:"version"`
	DurationMs int64  `json:"duration_ms"`
	Status     string `json:"status"` // "ok", "timeout", "error", "skipped"
}

// PostInput represents the A-face content input for pipeline processing.
type PostInput struct {
	PostID    int64    `json:"post_id"`
	Text      string   `json:"text"`
	MediaURLs []string `json:"media_urls"`
	Lang      string   `json:"lang"`
}

// Engine is the core Skill-Plus execution engine.
type Engine struct {
	// TODO: registry, dispatcher, sandbox manager
}

// Process runs all matching Skills against the given post input
// and returns the aggregated B-face output.
func (e *Engine) Process(ctx context.Context, input *PostInput) (*BFace, error) {
	// TODO: Implement dispatch → concurrent execution → aggregation
	return &BFace{
		Lang:          input.Lang,
		SkillsApplied: []SkillRunInfo{},
		Safety:        &SafetyResult{Flags: []string{}},
	}, nil
}
