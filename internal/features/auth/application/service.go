package application

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	authDomain "multicliente-backend/internal/features/auth/domain"
	userDomain "multicliente-backend/internal/features/user/domain"
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

	// Generate JWT token
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"role":     roleCode,
		"role_id":  roleIDVal,
		"iat":      now.Unix(),
		"exp":      now.Add(time.Hour * time.Duration(s.jwtExpHours)).Unix(),
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
