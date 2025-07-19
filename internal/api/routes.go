package api

import (
	"VK/internal/api/handlers"
	"VK/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authService handlers.AuthService,
	adService handlers.AdvertisementService,
	jwtSecret string,
) *gin.Engine {
	router := gin.Default()

	authHandler := handlers.NewAuthHandler(authService)
	adHandler := handlers.NewAdvertisementHandler(adService)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	apiGroup := router.Group("/ads")
	{
		apiGroup.GET("", Middleware(jwtSecret), adHandler.GetAds)
		apiGroup.POST("", JWTMiddleware(jwtSecret), adHandler.CreateAd)
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

func Middleware(jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        
        if token == "" {
            c.Next()
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