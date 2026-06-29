package infrastructure

import (
	"gorm.io/gorm"
	"multicliente-backend/internal/features/article/domain"
)

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) domain.ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) Create(art *domain.Article) error {
	return r.db.Create(art).Error
}

func (r *articleRepository) FindByID(id uint) (*domain.Article, error) {
	var art domain.Article
	if err := r.db.Preload("Category").First(&art, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &art, nil
}

func (r *articleRepository) FindAllByCompany(companyID uint) ([]domain.Article, error) {
	var articles []domain.Article
	if err := r.db.Preload("Category").Where("company_id = ?", companyID).Order("id ASC").Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

func (r *articleRepository) Update(art *domain.Article) error {
	return r.db.Save(art).Error
}

func (r *articleRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Article{}, "id = ?", id).Error
}
