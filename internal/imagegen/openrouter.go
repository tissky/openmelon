package imagegen

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// OpenRouter image generation client.
//
// Unlike the official OpenAI image API which has a dedicated
// /v1/images/generations endpoint, OpenRouter exposes image-capable
// models through /v1/chat/completions with modalities=["image","text"].
// The model returns base64-encoded PNGs in choices[0].message.images[].
//
// Image-capable models on OpenRouter today (2026-05):
//   - openai/gpt-5-image, openai/gpt-5-image-mini, openai/gpt-5.4-image-2
//   - google/gemini-2.5-flash-image (cheapest)
//   - google/gemini-3-pro-image-preview, google/gemini-3.1-flash-image-preview
//
// We use only net/http + encoding/json — same zero-dep posture as the
// rest of internal/imagegen.

const openRouterDefaultBaseURL = "https://openrouter.ai/api"

// OpenRouterGenerator implements Generator against OpenRouter's
// chat-completions image generation surface.
type OpenRouterGenerator struct {
	apiKey       string
	baseURL      string
	defaultModel string
	httpClient   *http.Client
}

// NewOpenRouter builds an OpenRouterGenerator.
//
// Env-var fallbacks (only used when the matching argument is ""):
//   - apiKey ← OPENROUTER_API_KEY
//   - baseURL ← OPENROUTER_BASE_URL
//
// defaultModel: required — pass an explicit model id (e.g.
// "google/gemini-2.5-flash-image" or "openai/gpt-5-image"). We do not
// bake in a default because the cheapest-vs-best model menu evolves.
func NewOpenRouter(apiKey, baseURL, defaultModel string) (*OpenRouterGenerator, error) {
	if apiKey == "" {
		apiKey = os.Getenv("OPENROUTER_API_KEY")
	}
	if apiKey == "" {
		return nil, ErrNoAPIKey
	}
	if baseURL == "" {
		baseURL = os.Getenv("OPENROUTER_BASE_URL")
	}
	if baseURL == "" {
		baseURL = openRouterDefaultBaseURL
	}
	if defaultModel == "" {
		return nil, ErrModelRequired
	}
	return &OpenRouterGenerator{
		apiKey:       apiKey,
		baseURL:      baseURL,
		defaultModel: defaultModel,
		httpClient:   &http.Client{Timeout: 5 * time.Minute},
	}, nil
}

func (g *OpenRouterGenerator) Provider() string { return "openrouter" }
func (g *OpenRouterGenerator) Model() string    { return g.defaultModel }
func (g *OpenRouterGenerator) BaseURL() string  { return g.baseURL }

// orChatRequest is a chat-completions request with the modalities field
// that OpenRouter requires for image generation.
type orChatRequest struct {
	Model      string          `json:"model"`
	Messages   []orChatMessage `json:"messages"`
	Modalities []string        `json:"modalities,omitempty"`
}

type orChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// orChatResponse is the chat-completions response. The image lives in
// choices[0].message.images[0].image_url.url as a data:image/png;base64
// URL.
type orChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
			Images  []struct {
				Type     string `json:"type"`
				ImageURL struct {
					URL string `json:"url"`
				} `json:"image_url"`
			} `json:"images"`
		} `json:"message"`
	} `json:"choices"`
}

// Generate sends a chat-completions request asking the model to
// generate an image, parses the data URL, and returns raw image bytes.
func (g *OpenRouterGenerator) Generate(ctx context.Context, opts GenerateOptions) (*Result, error) {
	if opts.Prompt == "" {
		return nil, fmt.Errorf("imagegen[openrouter]: Prompt is required")
	}
	model := opts.Model
	if model == "" {
		model = g.defaultModel
	}
	n := opts.N
	if n == 0 {
		n = 1
	}
	if n != 1 {
		return nil, fmt.Errorf("imagegen[openrouter]: only N=1 is supported today (got %d)", n)
	}
	// Note: opts.Size is ignored — OpenRouter's chat-completions image
	// surface doesn't accept a size hint; the model picks. If the caller
	// needs a specific size, embed it in the prompt itself.

	body, err := json.Marshal(orChatRequest{
		Model:      model,
		Messages:   []orChatMessage{{Role: "user", Content: opts.Prompt}},
		Modalities: []string{"image", "text"},
	})
	if err != nil {
		return nil, fmt.Errorf("imagegen[openrouter]: marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.baseURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("imagegen[openrouter]: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+g.apiKey)
	req.Header.Set("content-type", "application/json")
	// OpenRouter routing telemetry — optional but recommended.
	req.Header.Set("HTTP-Referer", "https://github.com/eight-acres-lab/openmelon")
	req.Header.Set("X-Title", "openmelon")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("imagegen[openrouter]: HTTP: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("imagegen[openrouter]: read response: %w", err)
	}
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("imagegen[openrouter]: HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var parsed orChatResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, fmt.Errorf("imagegen[openrouter]: parse response: %w (body: %s)", err, string(respBody))
	}
	if len(parsed.Choices) == 0 {
		return nil, fmt.Errorf("imagegen[openrouter]: no choices in response")
	}
	images := parsed.Choices[0].Message.Images
	if len(images) == 0 {
		// Sometimes the model refuses to generate; the text content has
		// the explanation. Surface it.
		txt := strings.TrimSpace(parsed.Choices[0].Message.Content)
		if txt != "" {
			return nil, fmt.Errorf("imagegen[openrouter]: no image in response (model said: %s)", txt)
		}
		return nil, fmt.Errorf("imagegen[openrouter]: no image in response (body: %s)", string(respBody))
	}
	url := images[0].ImageURL.URL
	imgBytes, contentType, err := decodeDataURL(url)
	if err != nil {
		return nil, fmt.Errorf("imagegen[openrouter]: %w", err)
	}

	return &Result{
		Data:        imgBytes,
		ContentType: contentType,
		Provider:    "openrouter",
		Model:       model,
		Prompt:      opts.Prompt,
		SizeBytes:   len(imgBytes),
	}, nil
}

// decodeDataURL parses a "data:<mime>;base64,<payload>" URL into raw
// bytes + the MIME content type. Returns a clear error for malformed
// URLs.
func decodeDataURL(dataURL string) ([]byte, string, error) {
	if !strings.HasPrefix(dataURL, "data:") {
		return nil, "", fmt.Errorf("not a data URL")
	}
	rest := strings.TrimPrefix(dataURL, "data:")
	commaIdx := strings.Index(rest, ",")
	if commaIdx < 0 {
		return nil, "", fmt.Errorf("data URL missing comma separator")
	}
	header := rest[:commaIdx]
	payload := rest[commaIdx+1:]

	// header is "<mime>;base64" or "<mime>" or ";base64"
	contentType := "application/octet-stream"
	isBase64 := false
	for _, part := range strings.Split(header, ";") {
		part = strings.TrimSpace(part)
		if part == "base64" {
			isBase64 = true
		} else if strings.Contains(part, "/") {
			contentType = part
		}
	}
	if !isBase64 {
		return nil, "", fmt.Errorf("data URL is not base64-encoded (got header: %q)", header)
	}
	decoded, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, "", fmt.Errorf("decode base64: %w", err)
	}
	return decoded, contentType, nil
}
