package artifacts

// Type identifies a production artifact type.
type Type string

const (
	TypeCopyDraft    Type = "copy_draft"
	TypeImagePrompt  Type = "image_prompt"
	TypeImage        Type = "image"
	TypeAudioScript  Type = "audio_script"
	TypeVideoShot    Type = "video_shot"
	TypeReviewReport Type = "review_report"
)

// Artifact is a typed production output with labels and provenance.
type Artifact struct {
	ID         string
	Type       Type
	Content    string
	Labels     map[string]string
	Provenance string
}
