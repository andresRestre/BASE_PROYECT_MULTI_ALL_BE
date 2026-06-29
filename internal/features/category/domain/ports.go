package domain

type CategoryRepository interface {
	Create(category *Category) error
	FindByID(id uint) (*Category, error)
	FindAllByCompany(companyID uint) ([]Category, error)
	Update(category *Category) error
	Delete(id uint) error
}

type CategoryService interface {
	CreateCategory(req *CreateCategoryRequest, companyID uint, createdBy *uint) (*Category, error)
	GetCategoryByID(id uint) (*Category, error)
	GetCategoriesByCompany(companyID uint) ([]Category, error)
	UpdateCategory(id uint, req *UpdateCategoryRequest, updatedBy *uint) (*Category, error)
	DeleteCategory(id uint) error
}
