package application

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	userDomain "multicliente-backend/internal/features/access_control/user/domain"
	authDomain "multicliente-backend/internal/features/auth/domain"
	"multicliente-backend/internal/platform/email"
)

type authService struct {
	userRepo     userDomain.UserRepository
	db           *gorm.DB
	emailService email.EmailService
	jwtSecret    string
	jwtExpHours  int
}

// NewAuthService creates a new AuthService.
func NewAuthService(
	userRepo userDomain.UserRepository,
	db *gorm.DB,
	emailService email.EmailService,
	jwtSecret string,
	jwtExpHours string,
) authDomain.AuthService {
	hours, err := strconv.Atoi(jwtExpHours)
	if err != nil || hours <= 0 {
		hours = 24
	}
	return &authService{
		userRepo:     userRepo,
		db:           db,
		emailService: emailService,
		jwtSecret:    jwtSecret,
		jwtExpHours:  hours,
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

func generate6DigitToken() string {
	nBig, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return fmt.Sprintf("%06d", time.Now().Nanosecond()%900000+100000)
	}
	return fmt.Sprintf("%06d", nBig.Int64()+100000)
}

func (s *authService) ForgotPassword(req *authDomain.ForgotPasswordRequest) error {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		// Return nil to avoid email enumeration
		return nil
	}

	if !user.IsActive {
		return errors.New("cuenta de usuario inactiva")
	}

	// Generate 6-digit recovery token
	tokenStr := generate6DigitToken()

	// Invalidate previous unused tokens for this email
	s.db.Model(&authDomain.PasswordResetToken{}).
		Where("email = ? AND used = ?", user.Email, false).
		Update("used", true)

	// Save new reset token record (expires in 15 minutes)
	resetRecord := authDomain.PasswordResetToken{
		UserID:    user.ID,
		Email:     user.Email,
		Token:     tokenStr,
		ExpiresAt: time.Now().Add(15 * time.Minute),
		Used:      false,
	}

	if err := s.db.Create(&resetRecord).Error; err != nil {
		return errors.New("error al generar código de recuperación")
	}

	// Send email / log to console
	return s.emailService.SendPasswordResetEmail(user.Email, tokenStr)
}

func (s *authService) ResetPassword(req *authDomain.ResetPasswordRequest) error {
	if req.NewPassword != req.ConfirmPassword {
		return errors.New("las contraseñas no coinciden")
	}

	var resetRecord authDomain.PasswordResetToken
	err := s.db.Where("email = ? AND token = ? AND used = ?", req.Email, req.Token, false).
		First(&resetRecord).Error
	if err != nil {
		return errors.New("el código de recuperación es inválido o ya fue utilizado")
	}

	if time.Now().After(resetRecord.ExpiresAt) {
		return errors.New("el código de recuperación ha expirado")
	}

	// Find user
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return errors.New("usuario no encontrado")
	}

	// Hash new password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error al procesar la nueva contraseña")
	}

	user.Password = string(hashed)
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("error al actualizar la contraseña")
	}

	// Mark token as used
	resetRecord.Used = true
	s.db.Save(&resetRecord)

	return nil
}
