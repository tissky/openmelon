package imagegen

import "fmt"

// New returns a Generator for the requested provider.
//
// provider must be one of: "openai", "openrouter".
// apiKey, baseURL, defaultModel are all optional — empty values fall
// back to per-provider env vars (OPENAI_API_KEY+OPENAI_BASE_URL or
// OPENROUTER_API_KEY+OPENROUTER_BASE_URL) and built-in defaults.
//
// defaultModel is REQUIRED — no vendor model defaults are baked in.
func New(provider, apiKey, baseURL, defaultModel string) (Generator, error) {
	switch provider {
	case "openai", "":
		return NewOpenAI(apiKey, baseURL, defaultModel)
	case "openrouter":
		return NewOpenRouter(apiKey, baseURL, defaultModel)
	default:
		return nil, fmt.Errorf("imagegen: unknown provider %q (supported: openai, openrouter)", provider)
	}
}
