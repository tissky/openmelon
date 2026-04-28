package workflow

// Stage identifies a content production stage.
type Stage string

const (
	StageIntentPlanning       Stage = "intent_planning"
	StageAngleSelection       Stage = "angle_selection"
	StageCopywriting          Stage = "copywriting"
	StageVisualConcretization Stage = "visual_concretization"
	StageGeneration           Stage = "generation"
	StageReview               Stage = "review"
	StagePackaging            Stage = "packaging"
)

// Workflow is a staged content production process.
type Workflow struct {
	ID           string
	Vertical     string
	Stages       []Stage
	CurrentStage Stage
}
