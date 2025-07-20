package services

import (
	"github.com/keenetic29/vk-internship/internal/domain"
	"testing"
)

type MockAdRepository struct {
	ads []*domain.Advertisement
}

func (m *MockAdRepository) Create(ad *domain.Advertisement) error {
	m.ads = append(m.ads, ad)
	return nil
}

func (m *MockAdRepository) GetAll(page, limit int, sortBy, order string, minPrice, maxPrice float64) ([]domain.Advertisement, error) {
	var result []domain.Advertisement
	for _, ad := range m.ads {
		if (minPrice == 0 || ad.Price >= minPrice) && (maxPrice == 0 || ad.Price <= maxPrice) {
			result = append(result, *ad)
		}
	}
	return result, nil
}

func TestAdvertisementService_CreateAd(t *testing.T) {
	repo := &MockAdRepository{}
	service := NewAdvertisementService(repo)

	// Успешное создание
	ad, err := service.CreateAd(1, "Title", "Description", "http://example.com/image.jpg", 100)
	if err != nil {
		t.Fatalf("CreateAd failed: %v", err)
	}

	if ad.Title != "Title" {
		t.Error("Ad title mismatch")
	}

	// Невалидные данные
	testCases := []struct {
		title       string
		description string
		price       float64
	}{
		{"", "Desc", 100},          
		{"Title", "", 100},         // пустое описание
		{"Title", "Desc", -100},    // отрицательная цена
		{"T", "Desc", 100},         // короткий заголовок
		{makeString(101), "Desc", 100}, // длинный заголовок
	}

	for _, tc := range testCases {
		_, err := service.CreateAd(1, tc.title, tc.description, "http://valid.url", tc.price)
		if err == nil {
			t.Errorf("Expected error for title=%q, desc=%q, price=%f", tc.title, tc.description, tc.price)
		}
	}
}

// Вспомогательная функция для генерации длинных строк
func makeString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = 'a'
	}
	return string(b)
}