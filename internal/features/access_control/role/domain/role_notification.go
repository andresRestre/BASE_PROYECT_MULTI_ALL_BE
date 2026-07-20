package domain

// RoleNotificationRule defines a rule where members of TargetRoleID receive notifications
// when users of CreatorRoleID perform Action on MenuID.
type RoleNotificationRule struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	TargetRoleID  uint   `gorm:"not null;index" json:"target_role_id"`
	MenuID        uint   `gorm:"not null;index" json:"menu_id"`
	CreatorRoleID uint   `gorm:"not null;index" json:"creator_role_id"`
	Action        string `gorm:"type:varchar(50);not null" json:"action"` // CREATE, EDIT, DELETE
	IsEnabled     bool   `gorm:"default:true" json:"is_enabled"`
}

func (RoleNotificationRule) TableName() string {
	return "administrative.role_notification_rules"
}

// NotificationRuleRequest represents payload for notification matrix rules
type NotificationRuleRequest struct {
	MenuID        uint   `json:"menu_id" binding:"required"`
	CreatorRoleID uint   `json:"creator_role_id" binding:"required"`
	Action        string `json:"action" binding:"required"` // CREATE, EDIT, DELETE
	IsEnabled     bool   `json:"is_enabled"`
}
