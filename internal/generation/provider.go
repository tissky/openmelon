package generation

import "context"

// Provider executes a generation request and returns artifact content and a trace.
type Provider interface {
	Generate(ctx context.Context, req *Request) (content string, trace *Trace, err error)
}

// Trace records how an artifact was produced (included in provenance).
type Trace struct {
	ProviderType string  `json:"provider_type"`
	Model        string  `json:"model,omitempty"`
	Command      string  `json:"command,omitempty"`
	DurationSec  float64 `json:"duration_sec"`
}

// ProviderError is a typed error from a Provider.Generate call.
type ProviderError struct {
	// Code is a machine-readable error category: "timeout", "non_zero_exit", "empty_output".
	Code    string
	Message string
	Wrapped error
}

func (e *ProviderError) Error() string { return e.Message }
func (e *ProviderError) Unwrap() error { return e.Wrapped }
