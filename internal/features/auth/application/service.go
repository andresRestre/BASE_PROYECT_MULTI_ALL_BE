package application

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	authDomain "multicliente-backend/internal/features/auth/domain"
	userDomain "multicliente-backend/internal/features/access_control/user/domain"
)

type authService struct {
	userRepo    userDomain.UserRepository
	jwtSecret   string
	jwtExpHours int
}

// NewAuthService creates a new AuthService.
// It depends on the user repository (cross-feature dependency via domain port).
func NewAuthService(userRepo userDomain.UserRepository, jwtSecret string, jwtExpHours string) authDomain.AuthService {
	hours, err := strconv.Atoi(jwtExpHours)
	if err != nil || hours <= 0 {
		hours = 24
	}
	return &authService{
		userRepo:    userRepo,
		jwtSecret:   jwtSecret,
		jwtExpHours: hours,
	}
}

func (s *authService) Login(req *authDomain.LoginRequest) (*authDomain.LoginResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Extract role details
	roleCode := ""
	var roleIDVal *uint
	if user.Role != nil {
		roleCode = user.Role.Code
		roleIDVal = &user.Role.ID
	}

	// Determine token expiration duration
	expirationDuration := time.Hour * time.Duration(s.jwtExpHours)
	if user.Role != nil {
		roleDuration := time.Duration(user.Role.SessionDays)*24*time.Hour +
			time.Duration(user.Role.SessionHours)*time.Hour +
			time.Duration(user.Role.SessionMinutes)*time.Minute
		if roleDuration > 0 {
			expirationDuration = roleDuration
		}
	}

	// Generate JWT token
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"role":     roleCode,
		"role_id":  roleIDVal,
		"iat":      now.Unix(),
		"exp":      now.Add(expirationDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &authDomain.LoginResponse{
		Token: tokenString,
		User: authDomain.UserInfo{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			RoleID:    roleIDVal,
			RoleCode:  roleCode,
		},
		SessionDurationSeconds: int(expirationDuration.Seconds()),
	}, nil
}

func (s *authService) GetProfile(userID uint) (*userDomain.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	// Hide password hash
	user.Password = ""
	return user, nil
}

func (s *authService) ChangePassword(userID uint, req *authDomain.ChangePasswordRequest) error {
	if req.CurrentPassword != req.ConfirmCurrentPassword {
		return errors.New("la contraseña actual y la confirmación no coinciden")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("usuario no encontrado")
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return errors.New("la contraseña actual es incorrecta")
	}

	// Hash new password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("falló al procesar la nueva contraseña")
	}

	user.Password = string(hashed)
	return s.userRepo.Update(user)
}
