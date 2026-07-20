package domain

type ArticleRepository interface {
	Create(article *Article) error
	FindByID(id uint) (*Article, error)
	FindAllByCompany(companyID uint) ([]Article, error)
	Update(article *Article) error
	Delete(id uint) error
}

type ArticleService interface {
	CreateArticle(req *CreateArticleRequest, companyID uint, createdBy *uint) (*Article, error)
	GetArticleByID(id uint) (*Article, error)
	GetArticlesByCompany(companyID uint) ([]Article, error)
	UpdateArticle(id uint, req *UpdateArticleRequest, updatedBy *uint) (*Article, error)
	DeleteArticle(id uint) error
}
