package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/keenetic29/vk-internship/internal/api/handlers"
	"github.com/keenetic29/vk-internship/internal/domain"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock HTTPClient для имитации HTTP запросов
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Head(url string) (*http.Response, error) {
	args := m.Called(url)
	return args.Get(0).(*http.Response), args.Error(1)
}

type MockAdvertisementService struct {
	mock.Mock
}

func (m *MockAdvertisementService) CreateAd(userID uint, title, description, imageURL string, price float64) (*domain.Advertisement, error) {
	args := m.Called(userID, title, description, imageURL, price)
	return args.Get(0).(*domain.Advertisement), args.Error(1)
}

func (m *MockAdvertisementService) GetAds(page, limit int, sortBy, order string, minPrice, maxPrice float64) ([]domain.Advertisement, error) {
	args := m.Called(page, limit, sortBy, order, minPrice, maxPrice)
	return args.Get(0).([]domain.Advertisement), args.Error(1)
}

// Вспомогательная функция для создания валидного HTTP ответа для изображения
func createValidImageResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type":   []string{"image/jpeg"},
			"Content-Length": []string{"1048576"}, // 1MB
		},
		ContentLength: 1048576,
		Body:          io.NopCloser(bytes.NewReader([]byte{})),
	}
}

// Вспомогательная функция для создания невалидного HTTP ответа
func createInvalidImageResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type":   []string{"text/html"},
			"Content-Length": []string{"10485760"}, // 10MB
		},
		ContentLength: 10485760,
		Body:          io.NopCloser(bytes.NewReader([]byte{})),
	}
}

func TestAdvertisementHandler_CreateAd(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  interface{}
		setupContext func(*gin.Context)
		mockSetup    func(*MockAdvertisementService, *MockHTTPClient)
		expectedCode int
	}{
		{
			name: "Successful ad creation",
			requestBody: map[string]interface{}{
				"title":       "Test Ad",
				"description": "Test Description",
				"image_url":   "http://valid.com/image.jpg",
				"price":       100.50,
			},
			setupContext: func(c *gin.Context) {
				c.Set("userID", uint(1))
			},
			mockSetup: func(as *MockAdvertisementService, hc *MockHTTPClient) {
				hc.On("Head", "http://valid.com/image.jpg").Return(createValidImageResponse(), nil)
				as.On("CreateAd", uint(1), "Test Ad", "Test Description", "http://valid.com/image.jpg", 100.50).
					Return(&domain.Advertisement{
						ID:          1,
						Title:       "Test Ad",
						Description: "Test Description",
						ImageURL:    "http://valid.com/image.jpg",
						Price:       100.50,
						UserID:      1,
					}, nil)
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "Unauthorized request",
			requestBody: map[string]interface{}{
				"title":       "Test Ad",
				"description": "Test Description",
				"image_url":   "http://example.com/image.jpg",
				"price":       100.50,
			},
			setupContext: func(c *gin.Context) {},
			mockSetup:    func(as *MockAdvertisementService, hc *MockHTTPClient) {},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "Invalid request body",
			requestBody: map[string]interface{}{
				"title":       "",
				"description": "Test Description",
				"image_url":   "http://example.com/image.jpg",
				"price":       100.50,
			},
			setupContext: func(c *gin.Context) {
				c.Set("userID", uint(1))
			},
			mockSetup:    func(as *MockAdvertisementService, hc *MockHTTPClient) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Image validation failed - invalid content type",
			requestBody: map[string]interface{}{
				"title":       "Test Ad",
				"description": "Test Description",
				"image_url":   "http://invalid.com/image.jpg",
				"price":       100.50,
			},
			setupContext: func(c *gin.Context) {
				c.Set("userID", uint(1))
			},
			mockSetup: func(as *MockAdvertisementService, hc *MockHTTPClient) {
				hc.On("Head", "http://invalid.com/image.jpg").Return(createInvalidImageResponse(), nil)
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Image validation failed - network error",
			requestBody: map[string]interface{}{
				"title":       "Test Ad",
				"description": "Test Description",
				"image_url":   "http://error.com/image.jpg",
				"price":       100.50,
			},
			setupContext: func(c *gin.Context) {
				c.Set("userID", uint(1))
			},
			mockSetup: func(as *MockAdvertisementService, hc *MockHTTPClient) {
				hc.On("Head", "http://error.com/image.jpg").Return((*http.Response)(nil), errors.New("connection error"))
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Service error",
			requestBody: map[string]interface{}{
				"title":       "Test Ad",
				"description": "Test Description",
				"image_url":   "http://valid.com/image.jpg",
				"price":       100.50,
			},
			setupContext: func(c *gin.Context) {
				c.Set("userID", uint(1))
			},
			mockSetup: func(as *MockAdvertisementService, hc *MockHTTPClient) {
				hc.On("Head", "http://valid.com/image.jpg").Return(createValidImageResponse(), nil)
				as.On("CreateAd", uint(1), "Test Ad", "Test Description", "http://valid.com/image.jpg", 100.50).
					Return((*domain.Advertisement)(nil), errors.New("service error"))
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAdvertisementService)
			mockHTTPClient := new(MockHTTPClient)
			tt.mockSetup(mockService, mockHTTPClient)

			handler := handlers.NewAdvertisementHandler(mockService)
			handler.SetHTTPClient(mockHTTPClient)

			router := setupTestRouter()
			router.POST("/ads", func(c *gin.Context) {
				tt.setupContext(c)
				handler.CreateAd(c)
			})

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/ads", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			mockService.AssertExpectations(t)
			mockHTTPClient.AssertExpectations(t)
		})
	}
}

func TestAdvertisementHandler_GetAds(t *testing.T) {
	now := time.Now()
	testAds := []domain.Advertisement{
		{
			ID:          1,
			Title:       "Ad 1",
			Description: "Description 1",
			ImageURL:    "http://example.com/image1.jpg",
			Price:       100.50,
			UserID:      1,
			User: domain.User{
				ID:       1,
				Username: "user1",
			},
			CreatedAt: now,
		},
		{
			ID:          2,
			Title:       "Ad 2",
			Description: "Description 2",
			ImageURL:    "http://example.com/image2.jpg",
			Price:       200.75,
			UserID:      2,
			User: domain.User{
				ID:       2,
				Username: "user2",
			},
			CreatedAt: now.Add(-time.Hour),
		},
	}

	tests := []struct {
		name         string
		queryParams  string
		setupContext func(*gin.Context)
		mockSetup    func(*MockAdvertisementService)
		expectedCode int
		expectedBody string
	}{
		{
			name:        "Successful get ads without auth",
			queryParams: "",
			setupContext: func(c *gin.Context) {},
			mockSetup: func(m *MockAdvertisementService) {
				m.On("GetAds", 1, 10, "created_at", "desc", 0.0, 0.0).Return(testAds, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `[{"id":1,"title":"Ad 1","description":"Description 1","image_url":"http://example.com/image1.jpg","price":100.5,"author_login":"user1","created_at":"`,
		},
		{
			name:        "Successful get ads with auth",
			queryParams: "",
			setupContext: func(c *gin.Context) {
				c.Set("userID", uint(1))
			},
			mockSetup: func(m *MockAdvertisementService) {
				m.On("GetAds", 1, 10, "created_at", "desc", 0.0, 0.0).Return(testAds, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `"is_owner":true`,
		},
		{
			name:        "With query parameters",
			queryParams: "?page=2&limit=5&sort_by=price&order=asc&min_price=100&max_price=300",
			setupContext: func(c *gin.Context) {
				c.Set("userID", uint(1))
			},
			mockSetup: func(m *MockAdvertisementService) {
				m.On("GetAds", 2, 5, "price", "asc", 100.0, 300.0).Return(testAds, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `"is_owner":true`,
		},
		{
			name:        "Service error",
			queryParams: "",
			setupContext: func(c *gin.Context) {
				c.Set("userID", uint(1))
			},
			mockSetup: func(m *MockAdvertisementService) {
				m.On("GetAds", 1, 10, "created_at", "desc", 0.0, 0.0).
					Return([]domain.Advertisement{}, errors.New("service error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":"service error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAdvertisementService)
			tt.mockSetup(mockService)

			handler := handlers.NewAdvertisementHandler(mockService)

			router := setupTestRouter()
			router.GET("/ads", func(c *gin.Context) {
				tt.setupContext(c)
				handler.GetAds(c)
			})

			req, _ := http.NewRequest("GET", "/ads"+tt.queryParams, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			
			if tt.expectedBody != "" {
				assert.True(t, strings.Contains(w.Body.String(), tt.expectedBody),
					"Response body should contain: %s, but got: %s", 
					tt.expectedBody, w.Body.String())
			}
			
			mockService.AssertExpectations(t)
		})
	}
}