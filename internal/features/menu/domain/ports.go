package domain

type MenuRepository interface {
	Create(menu *Menu) error
	FindByID(id uint) (*Menu, error)
	FindAll() ([]Menu, error)
	Update(menu *Menu) error
	Delete(id uint) error
	GetAllowedMenusForRole(roleID uint) ([]AllowedMenuResponse, error)
}

type MenuService interface {
	CreateMenu(req *CreateMenuRequest, createdBy *uint) (*Menu, error)
	GetMenuByID(id uint) (*Menu, error)
	GetAllMenus() ([]Menu, error)
	UpdateMenu(id uint, req *UpdateMenuRequest, updatedBy *uint) (*Menu, error)
	DeleteMenu(id uint) error
	GetAllowedMenus(roleID uint) ([]AllowedMenuResponse, error)
}
