package routers

import (
	"net/http"
	"time"

	"github.com/Cprime50/RegulerClub/handlers"
	"github.com/Cprime50/RegulerClub/util"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

func SetupRoutes(router *gin.Engine) {
	// CORS
	//router = gin.New()
	router.Use(cors.Middleware(cors.Config{
		Origins:         "https://*, http://*, *", //change in production to front-end url
		Methods:         "GET, POST, PUT, DELETE, OPTIONS",
		RequestHeaders:  "Origin, Accept, Authorization, Content-Type, X-CSRF-Token",
		ExposedHeaders:  "Link",
		Credentials:     false,
		MaxAge:          50 * time.Second,
		ValidateHeaders: false,
	}))

	// starting route
	router.GET("/", func(c *gin.Context) {
		time.Sleep(2 * time.Second)
		c.String(http.StatusOK, "Welcome Reguler Club Server")
	})

	// Create a group for base routes
	baseRoutes := router.Group("/api")
	{
		baseRoutes.POST("/register", handlers.Register)
		baseRoutes.POST("/login", handlers.Login)
		baseRoutes.PUT("/send-verify-email", handlers.ResendCode)
		baseRoutes.POST("/verify-email", handlers.VerifyEmail)
	}

	// Basic Authenticated routes
	authRoutes := router.Group("/api")
	authRoutes.Use(util.JWTAuth())
	{

	}

	// Admin routes
	adminRoutes := router.Group("/api/admin")
	adminRoutes.Use(util.JWTAuthAdmin())
	{

	}
}
