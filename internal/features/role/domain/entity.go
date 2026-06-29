package domain

import (
	"time"
)

type Role struct {
	ID           uint         `gorm:"primaryKey" json:"id"`
	Name         string       `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Code         string       `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	IsActive     bool         `gorm:"default:true" json:"is_active"`
	CreateBy     *uint        `json:"create_by"`
	CreateAt     time.Time    `gorm:"autoCreateTime" json:"create_at"`
	UpdateBy     *uint        `json:"update_by"`
	UpdateAt     time.Time    `gorm:"autoUpdateTime" json:"update_at"`
	CreateByName string       `gorm:"-" json:"create_by_name"`
	UpdateByName string       `gorm:"-" json:"update_by_name"`
	Permissions  []Permission `json:"permissions,omitempty" gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
}

func (Role) TableName() string {
	return "administrative.roles"
}

type Option struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar(50);not null" json:"name"`
	Code string `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
}

func (Option) TableName() string {
	return "administrative.options"
}

type Permission struct {
	RoleID   uint `gorm:"primaryKey" json:"role_id"`
	MenuID   uint `gorm:"primaryKey" json:"menu_id"`
	OptionID uint `gorm:"primaryKey" json:"option_id"`
}

func (Permission) TableName() string {
	return "administrative.permissions"
}

// DTOs
type PermissionRequest struct {
	MenuID   uint `json:"menu_id" binding:"required"`
	OptionID uint `json:"option_id" binding:"required"`
}

type CreateRoleRequest struct {
	Name        string              `json:"name" binding:"required"`
	Code        string              `json:"code" binding:"required"`
	Permissions []PermissionRequest `json:"permissions"`
}

type UpdateRoleRequest struct {
	Name        *string             `json:"name"`
	Code        *string             `json:"code"`
	IsActive    *bool               `json:"is_active"`
	Permissions []PermissionRequest `json:"permissions"`
}
