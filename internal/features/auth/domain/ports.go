package domain

import userDomain "multicliente-backend/internal/features/user/domain"

// AuthService defines the primary port for authentication operations.
type AuthService interface {
	Login(req *LoginRequest) (*LoginResponse, error)
	GetProfile(userID uint) (*userDomain.User, error)
}
