package handlers

import (
	"net/http"
	"strconv"
	"VK/internal/services"
	"github.com/gin-gonic/gin"
)

type AdvertisementHandler struct {
	adService services.AdvertisementService
}

func NewAdvertisementHandler(adService services.AdvertisementService) *AdvertisementHandler {
	return &AdvertisementHandler{adService: adService}
}

type CreateAdRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	ImageURL    string  `json:"image_url" binding:"required,url"`
	Price       float64 `json:"price" binding:"required"`
}

func (h *AdvertisementHandler) CreateAd(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateAdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ad, err := h.adService.CreateAd(userID.(uint), req.Title, req.Description, req.ImageURL, req.Price)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ad.IsOwner = true
	c.JSON(http.StatusCreated, ad)
}

func (h *AdvertisementHandler) GetAds(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")
	minPrice, _ := strconv.ParseFloat(c.Query("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("max_price"), 64)

	ads, err := h.adService.GetAds(page, limit, sortBy, order, minPrice, maxPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if exists {
		for i := range ads {
			ads[i].IsOwner = ads[i].UserID == userID.(uint)
		}
	}

	c.JSON(http.StatusOK, ads)
}