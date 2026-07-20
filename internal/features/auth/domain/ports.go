package domain

import userDomain "multicliente-backend/internal/features/access_control/user/domain"

// AuthService defines the primary port for authentication operations.
type AuthService interface {
	Login(req *LoginRequest) (*LoginResponse, error)
	GetProfile(userID uint) (*userDomain.User, error)
	ChangePassword(userID uint, req *ChangePasswordRequest) error
}
