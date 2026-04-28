package roles

// Role identifies a production role inside OpenMelon.
type Role string

const (
	RoleContentDirector     Role = "content_director"
	RoleImagePromptDirector Role = "image_prompt_director"
	RoleCopywriter          Role = "copywriter"
	RoleReviewer            Role = "reviewer"
	RoleMemoryManager       Role = "project_memory_manager"
)
