package domain

type NotificationRepository interface {
	Create(notification *Notification) error
	FindAllByUserAndCompany(userID uint, companyID uint) ([]Notification, error)
	FindByID(id uint) (*Notification, error)
	Update(notification *Notification) error
	MarkAllAsRead(userID uint, companyID uint) error
	FindAdminsAndSuperadminsByCompany(companyID uint) ([]uint, error)
	GetCreatorInfo(userID uint) (string, string, string, error) // firstName, lastName, roleCode, error
}

type NotificationService interface {
	CreateNotification(companyID uint, userID uint, title, message, notifType, route string) (*Notification, error)
	GetNotifications(userID uint, companyID uint) ([]Notification, error)
	MarkAsRead(id uint, userID uint) (*Notification, error)
	MarkAllRead(userID uint, companyID uint) error
	TriggerArticleCreatedNotification(companyID uint, creatorID uint, articleName string) error
}
