package repository

import (
	"github.com/keenetic29/vk-internship/internal/domain"
	"gorm.io/gorm"
)

type advertisementRepository struct {
	db *gorm.DB
}

func NewAdvertisementRepository(db *gorm.DB) *advertisementRepository {
	return &advertisementRepository{db: db}
}

func (r *advertisementRepository) Create(ad *domain.Advertisement) error {
	return r.db.Create(ad).Error
}

func (r *advertisementRepository) GetAll(page, limit int, sortBy, order string, minPrice, maxPrice float64) ([]domain.Advertisement, error) {
	var ads []domain.Advertisement

	query := r.db.Model(&domain.Advertisement{}).Preload("User")

	if minPrice > 0 {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice > 0 {
		query = query.Where("price <= ?", maxPrice)
	}

	if sortBy != "" {
		query = query.Order(sortBy + " " + order)
	} else {
		query = query.Order("created_at DESC")
	}

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&ads).Error

	return ads, err
}