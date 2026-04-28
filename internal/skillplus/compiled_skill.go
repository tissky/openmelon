package skillplus

// CompiledSkill is the OpenMelon-facing result of compiling a Skill-Plus package.
type CompiledSkill struct {
	PackageID      string
	PackageVersion string
	Target         string
	ModelProfile   string
	RuntimeVars    map[string]string
	Prompt         string
	Evaluation     []string
}
