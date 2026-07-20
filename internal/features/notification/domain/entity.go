package domain

import (
	"time"
)

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CompanyID uint      `gorm:"not null;index" json:"company_id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"` // Target user who receives the notification
	Title     string    `gorm:"type:varchar(255);not null" json:"title"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	Type      string    `gorm:"type:varchar(50);not null" json:"type"` // e.g. "article_created"
	IsRead    bool      `gorm:"default:false" json:"is_read"`
	Route     string    `gorm:"type:varchar(255)" json:"route"` // Optional route for frontend navigation
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (Notification) TableName() string {
	return "app.notifications"
}
