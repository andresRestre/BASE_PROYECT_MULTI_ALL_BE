package domain

import (
	"time"
)

type Company struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Name         string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	PhotoURL     string     `gorm:"type:varchar(255);default:''" json:"photo_url"`
	CreateBy     *uint      `json:"create_by"`
	CreateAt     time.Time  `gorm:"autoCreateTime" json:"create_at"`
	UpdateBy     *uint      `json:"update_by"`
	UpdateAt     time.Time  `gorm:"autoUpdateTime" json:"update_at"`
	CreateByName string     `gorm:"-" json:"create_by_name"`
	UpdateByName string     `gorm:"-" json:"update_by_name"`
}

func (Company) TableName() string {
	return "administrative.companies"
}

// DTOs
type CreateCompanyRequest struct {
	Name     string `json:"name" binding:"required"`
	PhotoURL string `json:"photo_url"`
}

type UpdateCompanyRequest struct {
	Name     *string `json:"name"`
	IsActive *bool   `json:"is_active"`
	PhotoURL *string `json:"photo_url"`
}
