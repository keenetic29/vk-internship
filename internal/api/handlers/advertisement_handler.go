package handlers

import (
	"github.com/keenetic29/vk-internship/internal/domain"
	"github.com/keenetic29/vk-internship/pkg/logger"
	"errors"
	"strings"
	"time"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	MaxImageSize       = 10 * 1024 * 1024 // 10MB
	ImageCheckTimeout  = 2 * time.Second
	AllowedImageTypes  = "image/jpeg,image/png,image/webp"
)

type HTTPClient interface {
    Head(url string) (*http.Response, error)
}

type AdvertisementService interface {
	CreateAd(userID uint, title, description, imageURL string, price float64) (*domain.Advertisement, error)
	GetAds(page, limit int, sortBy, order string, minPrice, maxPrice float64) ([]domain.Advertisement, error)
}

type AdvertisementHandler struct {
	adService AdvertisementService
	httpClient HTTPClient
}

func NewAdvertisementHandler(adService AdvertisementService) *AdvertisementHandler {
	return &AdvertisementHandler{
        adService: adService,
        httpClient: &http.Client{Timeout: ImageCheckTimeout},
    }
}

// устанавливает HTTPClient (для тестов)
func (h *AdvertisementHandler) SetHTTPClient(client HTTPClient) {
	h.httpClient = client
}

type CreateAdRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	ImageURL    string  `json:"image_url" binding:"required,url"`
	Price       float64 `json:"price" binding:"required"`
}

func (h *AdvertisementHandler) validateImageURL(imageURL string) error {
	logger.Log.Debug("Validating image URL", "url", imageURL)

	resp, err := h.httpClient.Head(imageURL)
	if err != nil {
		logger.Log.Warn("Image URL validation failed", 
			"error", err,
			"url", imageURL,
		)
		return errors.New("invalid image URL or unable to verify")
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(AllowedImageTypes, contentType) {
		logger.Log.Warn("Unsupported image format", 
			"content_type", contentType,
			"allowed_types", AllowedImageTypes,
		)
		return errors.New("only JPEG, PNG and WEBP images are allowed")
	}

	if size := resp.ContentLength; size > MaxImageSize {
		logger.Log.Warn("Image size exceeds limit",
			"size", size,
			"max_allowed", MaxImageSize,
		)
		return errors.New("image size exceeds maximum limit")
	}

	logger.Log.Debug("Image URL validation successful")

	return nil
}

func (h *AdvertisementHandler) CreateAd(c *gin.Context) {
	logger.Log.Info("CreateAd request started",
		"path", c.Request.URL.Path,
		"method", c.Request.Method,
	)
	
	userID, exists := c.Get("userID")
	if !exists {
		logger.Log.Warn("Unauthorized create ad attempt",
			"client_ip", c.ClientIP(),
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateAdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("Invalid request body",
			"error", err,
			"user_id", userID,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Debug("CreateAd request data",
		"title_length", len(req.Title),
		"description_length", len(req.Description),
		"price", req.Price,
	)

	if err := h.validateImageURL(req.ImageURL); err != nil {
		logger.Log.Warn("Image validation failed",
			"error", err,
			"user_id", userID,
			"image_url", req.ImageURL,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ad, err := h.adService.CreateAd(userID.(uint), req.Title, req.Description, req.ImageURL, req.Price)
	if err != nil {
		logger.Log.Error("Failed to create advertisement",
			"error", err,
			"user_id", userID,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Info("Advertisement created successfully",
		"ad_id", ad.ID,
		"user_id", userID,
		"title", ad.Title,
	)

	ad.IsOwner = true
	c.JSON(http.StatusCreated, ad)
}

func (h *AdvertisementHandler) GetAds(c *gin.Context) {
	logger.Log.Info("GetAds request started",
		"path", c.Request.URL.Path,
		"query", c.Request.URL.RawQuery,
	)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")
	minPrice, _ := strconv.ParseFloat(c.Query("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("max_price"), 64)

	logger.Log.Debug("GetAds query parameters",
		"page", page,
		"limit", limit,
		"sort_by", sortBy,
		"order", order,
		"min_price", minPrice,
		"max_price", maxPrice,
	)

	ads, err := h.adService.GetAds(page, limit, sortBy, order, minPrice, maxPrice)
	if err != nil {
		logger.Log.Error("Failed to get advertisements",
			"error", err,
			"query", c.Request.URL.RawQuery,
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

    var currentUserID uint
    if userID, exists := c.Get("userID"); exists {
        if uid, ok := userID.(uint); ok {
			logger.Log.Debug("User authenticated",
				"user_id", currentUserID,
			)
            currentUserID = uid
        }
    }

    type responseAd struct {
        ID          uint      `json:"id"`
        Title       string    `json:"title"`
        Description string    `json:"description"`
        ImageURL    string    `json:"image_url"`
        Price       float64   `json:"price"`
        AuthorLogin string    `json:"author_login"`
        CreatedAt   time.Time `json:"created_at"`
        IsOwner     *bool     `json:"is_owner,omitempty"` 
    }

    var response []responseAd
    for _, ad := range ads {
        item := responseAd{
            ID:          ad.ID,
            Title:       ad.Title,
            Description: ad.Description,
            ImageURL:    ad.ImageURL,
            Price:       ad.Price,
            AuthorLogin: ad.User.Username,
            CreatedAt:   ad.CreatedAt,
        }

        if currentUserID != 0 {
            isOwner := ad.UserID == currentUserID
            item.IsOwner = &isOwner 
        }

        response = append(response, item)
    }

	logger.Log.Info("GetAds request completed",
		"ads_count", len(response),
		"page", page,
	)
	
    c.JSON(http.StatusOK, response)
}