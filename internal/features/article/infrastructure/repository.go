package infrastructure

import (
	"gorm.io/gorm"
	"multicliente-backend/internal/features/article/domain"
	"multicliente-backend/internal/platform/database"
)

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) domain.ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) populateAudits(articles []*domain.Article) {
	var userIDs []uint
	for _, a := range articles {
		if a.CreateBy != nil {
			userIDs = append(userIDs, *a.CreateBy)
		}
		if a.UpdateBy != nil {
			userIDs = append(userIDs, *a.UpdateBy)
		}
	}
	namesMap, err := database.GetUserNamesMap(r.db, userIDs)
	if err != nil {
		return
	}
	for _, a := range articles {
		if a.CreateBy != nil {
			a.CreateByName = namesMap[*a.CreateBy]
		}
		if a.UpdateBy != nil {
			a.UpdateByName = namesMap[*a.UpdateBy]
		}
	}
}

func (r *articleRepository) Create(art *domain.Article) error {
	return r.db.Create(art).Error
}

func (r *articleRepository) FindByID(id uint) (*domain.Article, error) {
	var art domain.Article
	if err := r.db.Preload("Category").First(&art, "id = ?", id).Error; err != nil {
		return nil, err
	}
	r.populateAudits([]*domain.Article{&art})
	return &art, nil
}

func (r *articleRepository) FindAllByCompany(companyID uint) ([]domain.Article, error) {
	var articles []domain.Article
	if err := r.db.Preload("Category").Where("company_id = ?", companyID).Order("id ASC").Find(&articles).Error; err != nil {
		return nil, err
	}
	articlesPtrs := make([]*domain.Article, len(articles))
	for i := range articles {
		articlesPtrs[i] = &articles[i]
	}
	r.populateAudits(articlesPtrs)
	return articles, nil
}

func (r *articleRepository) Update(art *domain.Article) error {
	return r.db.Save(art).Error
}

func (r *articleRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Article{}, "id = ?", id).Error
}
