package domain

import (
	"time"
)

type Category struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	CompanyID uint       `gorm:"not null" json:"company_id"`
	Name      string     `gorm:"type:varchar(100);not null" json:"name"`
	IsActive  bool       `gorm:"default:true" json:"is_active"`
	CreateBy  *uint      `json:"create_by"`
	CreateAt  time.Time  `gorm:"autoCreateTime" json:"create_at"`
	UpdateBy  *uint      `json:"update_by"`
	UpdateAt  time.Time  `gorm:"autoUpdateTime" json:"update_at"`
}

func (Category) TableName() string {
	return "administrative.categories"
}

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateCategoryRequest struct {
	Name     *string `json:"name"`
	IsActive *bool   `json:"is_active"`
}
