package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/keenetic29/vk-internship/internal/api/handlers"
	"github.com/keenetic29/vk-internship/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(username, password string) (*domain.User, error) {
	args := m.Called(username, password)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockAuthService) Login(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ValidateToken(token string) (uint, error) {
	args := m.Called(token)
	return args.Get(0).(uint), args.Error(1)
}



func TestAuthHandler_Register(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  interface{}
		mockSetup    func(*MockAuthService)
		expectedCode int
	}{
		{
			name: "Successful registration",
			requestBody: map[string]string{
				"username": "testuser",
				"password": "testpass",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("Register", "testuser", "testpass").Return(&domain.User{
					ID:       1,
					Username: "testuser",
				}, nil)
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "Invalid request body",
			requestBody: map[string]string{
				"username": "",
			},
			mockSetup:    func(m *MockAuthService) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Username already exists",
			requestBody: map[string]string{
				"username": "existinguser",
				"password": "testpass",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("Register", "existinguser", "testpass").Return(
					(*domain.User)(nil),
					assert.AnError,
				)
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.mockSetup(mockService)

			// Используем handlers.NewAuthHandler вместо NewAuthHandler
			handler := handlers.NewAuthHandler(mockService)
			router := setupTestRouter()
			router.POST("/register", handler.Register)

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  interface{}
		mockSetup    func(*MockAuthService)
		expectedCode int
	}{
		{
			name: "Successful login",
			requestBody: map[string]string{
				"username": "testuser",
				"password": "testpass",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("Login", "testuser", "testpass").Return("testtoken12345", nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Invalid credentials",
			requestBody: map[string]string{
				"username": "testuser",
				"password": "wrongpass",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("Login", "testuser", "wrongpass").Return("", assert.AnError)
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "Invalid request body",
			requestBody: map[string]string{
				"username": "",
			},
			mockSetup:    func(m *MockAuthService) {},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.mockSetup(mockService)

			// Используем handlers.NewAuthHandler вместо NewAuthHandler
			handler := handlers.NewAuthHandler(mockService)
			router := setupTestRouter()
			router.POST("/login", handler.Login)

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}