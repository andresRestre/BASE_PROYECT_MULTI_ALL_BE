package domain

import (
	"time"

	categoryDomain "multicliente-backend/internal/features/category/domain"
)

type Article struct {
	ID         uint                    `gorm:"primaryKey" json:"id"`
	CompanyID  uint                    `gorm:"not null" json:"company_id"`
	CategoryID uint                    `gorm:"not null" json:"category_id"`
	Category   *categoryDomain.Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Name       string                  `gorm:"type:varchar(150);not null" json:"name"`
	CreateBy   *uint                   `json:"create_by"`
	CreateAt   time.Time               `gorm:"autoCreateTime" json:"create_at"`
	UpdateBy   *uint                   `json:"update_by"`
	UpdateAt   time.Time               `gorm:"autoUpdateTime" json:"update_at"`
}

func (Article) TableName() string {
	return "administrative.articles"
}

type CreateArticleRequest struct {
	CategoryID uint   `json:"category_id" binding:"required"`
	Name       string `json:"name" binding:"required"`
}

type UpdateArticleRequest struct {
	CategoryID *uint   `json:"category_id"`
	Name       *string `json:"name"`
}
