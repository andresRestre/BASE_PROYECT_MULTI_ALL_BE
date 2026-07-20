package domain

import "time"

type LandingText struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Key      string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"key"`
	Section  string    `gorm:"type:varchar(50);not null" json:"section"`
	TextES   string    `gorm:"type:text;not null" json:"text_es"`
	TextEN   string    `gorm:"type:text;not null" json:"text_en"`
	TextFR   string    `gorm:"type:text;not null" json:"text_fr"`
	UpdateAt time.Time `gorm:"autoUpdateTime" json:"update_at"`
}

func (LandingText) TableName() string {
	return "cms.landing_texts"
}

type UpdateTextRequest struct {
	TextES string `json:"text_es" binding:"required"`
	TextEN string `json:"text_en" binding:"required"`
	TextFR string `json:"text_fr" binding:"required"`
}

type LandingNews struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TitleES     string    `gorm:"type:varchar(255);not null" json:"title_es"`
	TitleEN     string    `gorm:"type:varchar(255);not null" json:"title_en"`
	TitleFR     string    `gorm:"type:varchar(255);not null" json:"title_fr"`
	ContentES   string    `gorm:"type:text;not null" json:"content_es"`
	ContentEN   string    `gorm:"type:text;not null" json:"content_en"`
	ContentFR   string    `gorm:"type:text;not null" json:"content_fr"`
	ImageURL    string    `gorm:"type:varchar(500)" json:"image_url"`
	IsPublished bool      `gorm:"default:true" json:"is_published"`
	CreateAt    time.Time `gorm:"autoCreateTime" json:"create_at"`
	UpdateAt    time.Time `gorm:"autoUpdateTime" json:"update_at"`
}

func (LandingNews) TableName() string {
	return "cms.landing_news"
}

type CreateNewsRequest struct {
	TitleES     string `json:"title_es" binding:"required"`
	TitleEN     string `json:"title_en" binding:"required"`
	TitleFR     string `json:"title_fr" binding:"required"`
	ContentES   string `json:"content_es" binding:"required"`
	ContentEN   string `json:"content_en" binding:"required"`
	ContentFR   string `json:"content_fr" binding:"required"`
	ImageURL    string `json:"image_url"`
	IsPublished bool   `json:"is_published"`
}

type UpdateNewsRequest struct {
	TitleES     *string `json:"title_es"`
	TitleEN     *string `json:"title_en"`
	TitleFR     *string `json:"title_fr"`
	ContentES   *string `json:"content_es"`
	ContentEN   *string `json:"content_en"`
	ContentFR   *string `json:"content_fr"`
	ImageURL    *string `json:"image_url"`
	IsPublished *bool   `json:"is_published"`
}

type LandingBanner struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"type:varchar(255);not null" json:"title"`
	Subtitle  string    `gorm:"type:varchar(255)" json:"subtitle"`
	ImageURL  string    `gorm:"type:varchar(500);not null" json:"image_url"`
	LinkURL   string    `gorm:"type:varchar(500)" json:"link_url"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreateAt  time.Time `gorm:"autoCreateTime" json:"create_at"`
}

func (LandingBanner) TableName() string {
	return "cms.landing_banners"
}

type CreateBannerRequest struct {
	Title     string `json:"title" binding:"required"`
	Subtitle  string `json:"subtitle"`
	ImageURL  string `json:"image_url" binding:"required"`
	LinkURL   string `json:"link_url"`
	SortOrder int    `json:"sort_order"`
	IsActive  bool   `json:"is_active"`
}

type UpdateBannerRequest struct {
	Title     *string `json:"title"`
	Subtitle  *string `json:"subtitle"`
	ImageURL  *string `json:"image_url"`
	LinkURL   *string `json:"link_url"`
	SortOrder *int    `json:"sort_order"`
	IsActive  *bool   `json:"is_active"`
}
