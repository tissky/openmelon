// Copyright 2026 Point Eight AI Pte. Ltd.
// Licensed under the Apache License, Version 2.0

// Package dispatcher implements Skill dispatch based on skill.yaml dispatch_hints.
//
// The dispatcher evaluates each registered Skill's invoke_when / do_not_invoke_when
// conditions against the content being processed and selects matching Skills.
package dispatcher

import "github.com/pointeight/skillplus-engine/engine"

// Dispatcher selects which Skills to run for a given content input.
type Dispatcher struct {
	// TODO: skill registry reference
}

// Match returns the list of Skill slugs that should be run for the given input.
func (d *Dispatcher) Match(input *engine.PostInput) []string {
	// TODO: Evaluate dispatch_hints for each registered skill
	return nil
}
