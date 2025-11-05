package permissions

type Permission string

const (
	// Organisation-level
	OrgEdit          Permission = "org:edit"
	OrgDelete        Permission = "org:delete"
	OrgViewMembers   Permission = "org:view_members"
	OrgInviteMembers Permission = "org:invite_member"
	OrgRemoveMembers Permission = "org:remove_member"
	OrgAssignRole    Permission = "org:assign_role"

	// Project-level
	ProjectCreate Permission = "project:create"
	ProjectEdit   Permission = "project:edit"
	ProjectView   Permission = "project:view"
	ProjectDelete Permission = "project:delete"

	// Kanban
	KanbanCreate Permission = "kanban:create"
	KanbanEdit   Permission = "kanban:edit"
	KanbanView   Permission = "kanban:view"
	KanbanDelete Permission = "kanban:delete"

	// Whiteboard
	WhiteboardCreate Permission = "whiteboard:create"
	WhiteboardEdit   Permission = "whiteboard:edit"
	WhiteboardView   Permission = "whiteboard:view"
	WhiteboardDelete Permission = "whiteboard:delete"
)

var (
	// Default permissions per role
	OwnerPermissions = []Permission{
		OrgEdit, OrgDelete, OrgViewMembers, OrgInviteMembers, OrgRemoveMembers, OrgAssignRole,
		ProjectCreate, ProjectEdit, ProjectView, ProjectDelete,
		KanbanCreate, KanbanEdit, KanbanView, KanbanDelete,
		WhiteboardCreate, WhiteboardEdit, WhiteboardView, WhiteboardDelete,
	}

	AdminPermissions = []Permission{
		OrgEdit, OrgViewMembers, OrgInviteMembers, OrgRemoveMembers, OrgAssignRole,
		ProjectCreate, ProjectEdit, ProjectView, ProjectDelete,
		KanbanCreate, KanbanEdit, KanbanView, KanbanDelete,
		WhiteboardCreate, WhiteboardEdit, WhiteboardView, WhiteboardDelete,
	}

	MemberPermissions = []Permission{
		ProjectView, ProjectCreate, ProjectEdit,
		KanbanView, KanbanCreate, KanbanEdit,
		WhiteboardView, WhiteboardCreate, WhiteboardEdit,
	}

	ViewerPermissions = []Permission{
		ProjectView, KanbanView, WhiteboardView,
	}
)
