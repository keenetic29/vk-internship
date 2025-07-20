package services

import (
	"github.com/keenetic29/vk-internship/internal/domain"
	"errors"
	"testing"
)

type MockUserRepository struct {
	users map[string]*domain.User
}

func (m *MockUserRepository) Create(user *domain.User) error {
	if _, exists := m.users[user.Username]; exists {
		return errors.New("user already exists")
	}
	m.users[user.Username] = user
	return nil
}

func (m *MockUserRepository) GetByUsername(username string) (*domain.User, error) {
	if user, exists := m.users[username]; exists {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) Exists(username string) (bool, error) {
	_, exists := m.users[username]
	return exists, nil
}

func TestAuthService_Register(t *testing.T) {
	repo := &MockUserRepository{users: make(map[string]*domain.User)}
	service := NewAuthService(repo, "test-secret")

	// Успешная регистрация
	user, err := service.Register("testuser", "password123")
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Проверяем, что пароль хэширован
	if user.Password == "password123" {
		t.Error("Password was not hashed")
	}

	// Дублирование пользователя
	_, err = service.Register("testuser", "newpass")
	if err == nil {
		t.Error("Duplicate username should fail")
	}
}

func TestAuthService_Login(t *testing.T) {
	repo := &MockUserRepository{users: make(map[string]*domain.User)}
	service := NewAuthService(repo, "test-secret")

	// Предварительно регистрируем пользователя
	_, _ = service.Register("testuser", "password123")

	// Успешный логин
	token, err := service.Login("testuser", "password123")
	if err != nil || token == "" {
		t.Error("Valid login should succeed")
	}

	// Неверный пароль
	_, err = service.Login("testuser", "wrongpass")
	if err == nil {
		t.Error("Invalid password should fail")
	}

	// Несуществующий пользователь
	_, err = service.Login("nobody", "pass")
	if err == nil {
		t.Error("Non-existent user should fail")
	}
}