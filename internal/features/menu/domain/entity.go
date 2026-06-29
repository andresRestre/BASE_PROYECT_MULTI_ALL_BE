package domain

import (
	"time"
)

type Menu struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Label     string     `gorm:"type:varchar(100);not null" json:"label"`
	LabelEN   string     `gorm:"type:varchar(100);not null" json:"label_en"`
	Route     string     `gorm:"type:varchar(255);not null" json:"route"`
	Icon      string     `gorm:"type:varchar(100)" json:"icon"`
	IsActive  bool       `gorm:"default:true" json:"is_active"`
	SortOrder int        `gorm:"default:0" json:"sort_order"`
	ParentID  *uint      `gorm:"default:null" json:"parent_id"`
	CreateBy  *uint      `json:"create_by"`
	CreateAt  time.Time  `gorm:"autoCreateTime" json:"create_at"`
	UpdateBy  *uint      `json:"update_by"`
	UpdateAt  time.Time  `gorm:"autoUpdateTime" json:"update_at"`
}

func (Menu) TableName() string {
	return "administrative.menus"
}

// AllowedMenuResponse includes options/permissions code list and hierarchical submenus.
type AllowedMenuResponse struct {
	ID          uint                  `json:"id"`
	Label       string                `json:"label"`
	LabelEN     string                `json:"label_en"`
	Route       string                `json:"route"`
	Icon        string                `json:"icon"`
	SortOrder   int                   `json:"sort_order"`
	ParentID    *uint                 `json:"parent_id"`
	Permissions []string              `json:"permissions"` // e.g. ["VIEW", "EDIT"]
	Submenus    []AllowedMenuResponse `json:"submenus"`
}

type CreateMenuRequest struct {
	Label     string `json:"label" binding:"required"`
	LabelEN   string `json:"label_en" binding:"required"`
	Route     string `json:"route" binding:"required"`
	Icon      string `json:"icon"`
	SortOrder int    `json:"sort_order"`
	ParentID  *uint  `json:"parent_id"`
}

type UpdateMenuRequest struct {
	Label     *string `json:"label"`
	LabelEN   *string `json:"label_en"`
	Route     *string `json:"route"`
	Icon      *string `json:"icon"`
	IsActive  *bool   `json:"is_active"`
	SortOrder *int    `json:"sort_order"`
	ParentID  *uint   `json:"parent_id"`
}
