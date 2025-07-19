package services

import (
	"VK/internal/domain"
	pass "VK/pkg/password"
	"VK/pkg/jwt"
	"errors"
	"time"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetByUsername(username string) (*domain.User, error)
	Exists(username string) (bool, error)
}

type authService struct {
	userRepo UserRepository
	jwtSecret string
}

func NewAuthService(userRepo UserRepository, jwtSecret string) *authService {
	return &authService{
		userRepo: userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Register(username, password string) (*domain.User, error) {
	exists, err := s.userRepo.Exists(username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	if len(username) < 3 || len(username) > 20 {
		return nil, errors.New("username must be between 3 and 20 characters")
	}

	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	hashedPassword, err := pass.HashPassword(password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &domain.User{
		Username: username,
		Password: hashedPassword, 
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(username, password string) (string, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := pass.CheckPassword(password, user.Password); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := jwt.GenerateToken(user.ID, s.jwtSecret, 1*time.Hour)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) ValidateToken(token string) (uint, error) {
	claims, err := jwt.ParseToken(token, s.jwtSecret)
	if err != nil {
		return 0, err
	}

	return claims.UserID, nil
}