package api

import (
	"VK/internal/api/handlers"
	"VK/internal/repository"
	"VK/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, jwtSecret string) *gin.Engine {
	router := gin.Default()

	userRepo := repository.NewUserRepository(db)
	adRepo := repository.NewAdvertisementRepository(db)

	authService := services.NewAuthService(userRepo, jwtSecret)
	adService := services.NewAdvertisementService(adRepo)

	authHandler := handlers.NewAuthHandler(authService)
	adHandler := handlers.NewAdvertisementHandler(adService)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	adGroup := router.Group("/ads")
	{
		adGroup.GET("", adHandler.GetAds)
		adGroup.POST("", JWTMiddleware(jwtSecret), adHandler.CreateAd)
	}

	return router
}

func JWTMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization token required"})
			c.Abort()
			return
		}

		userID, err := services.NewAuthService(nil, jwtSecret).ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}