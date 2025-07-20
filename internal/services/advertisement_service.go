package services

import (
	"github.com/keenetic29/vk-internship/internal/domain"
	"errors"
)

type AdvertisementRepository interface {
	Create(ad *domain.Advertisement) error
	GetAll(page, limit int, sortBy, order string, minPrice, maxPrice float64) ([]domain.Advertisement, error)
}

type advertisementService struct {
	adRepo AdvertisementRepository
}

func NewAdvertisementService(adRepo AdvertisementRepository) *advertisementService {
	return &advertisementService{adRepo: adRepo}
}

func (s *advertisementService) CreateAd(userID uint, title, description, imageURL string, price float64) (*domain.Advertisement, error) {
	if len(title) < 5 || len(title) > 100 {
		return nil, errors.New("title must be between 5 and 100 characters")
	}

	if len(description) < 10 || len(description) > 1000 {
		return nil, errors.New("description must be between 10 and 1000 characters")
	}

	if price <= 0 {
		return nil, errors.New("price must be positive")
	}

	ad := &domain.Advertisement{
		Title:       title,
		Description: description,
		ImageURL:    imageURL,
		Price:       price,
		UserID:      userID,
	}

	if err := s.adRepo.Create(ad); err != nil {
		return nil, err
	}

	return ad, nil
}

func (s *advertisementService) GetAds(page, limit int, sortBy, order string, minPrice, maxPrice float64) ([]domain.Advertisement, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 10
	}

	if sortBy != "" && sortBy != "price" && sortBy != "created_at" {
		sortBy = "created_at"
	}

	if order != "" && order != "asc" && order != "desc" {
		order = "desc"
	}

	return s.adRepo.GetAll(page, limit, sortBy, order, minPrice, maxPrice)
}