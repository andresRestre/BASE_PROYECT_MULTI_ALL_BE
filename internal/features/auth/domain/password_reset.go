package domain

import (
	"time"
)

// PasswordResetToken represents a record in database stored for user password recovery.
type PasswordResetToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Email     string    `gorm:"type:varchar(255);not null;index" json:"email"`
	Token     string    `gorm:"type:varchar(255);not null;index" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Used      bool      `gorm:"default:false" json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

func (PasswordResetToken) TableName() string {
	return "administrative.password_reset_tokens"
}

// ForgotPasswordRequest represents the request payload for initiating password reset.
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents the request payload for completing password reset.
type ResetPasswordRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Token           string `json:"token" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6"`
}
