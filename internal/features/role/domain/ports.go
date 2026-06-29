package domain

type RoleRepository interface {
	Create(role *Role) error
	FindByID(id uint) (*Role, error)
	FindAll() ([]Role, error)
	Update(role *Role) error
	Delete(id uint) error
	ReplacePermissions(roleID uint, permissions []Permission) error
	FindAllOptions() ([]Option, error)
}

type RoleService interface {
	CreateRole(req *CreateRoleRequest, createdBy *uint) (*Role, error)
	GetRoleByID(id uint) (*Role, error)
	GetAllRoles() ([]Role, error)
	UpdateRole(id uint, req *UpdateRoleRequest, updatedBy *uint) (*Role, error)
	DeleteRole(id uint) error
	GetAllOptions() ([]Option, error)
}
