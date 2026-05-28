// Package router wires HTTP routes to handlers (the API map of your app).
package router

import (
	"github.com/deckforge/backend/internal/config"
	"github.com/deckforge/backend/internal/handler"
	"github.com/deckforge/backend/internal/middleware"
	"github.com/deckforge/backend/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Setup builds the Gin engine with all routes and middleware.
func Setup(
	cfg *config.Config,
	authHandler *handler.AuthHandler,
	presentationHandler *handler.PresentationHandler,
	authService *service.AuthService,
) *gin.Engine {
	r := gin.Default()

	// CORS allows the React dev server (different port) to call the API
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.CORSOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check for Docker / load balancers
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "deckforge-api"})
	})

	api := r.Group("/api/v1")
	{
		// Public auth routes
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)

		// Protected routes (require valid JWT)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			protected.POST("/upload", presentationHandler.Upload)
			protected.POST("/generate", presentationHandler.Generate)
			protected.GET("/presentations", presentationHandler.ListPresentations)
			protected.GET("/presentation/:id", presentationHandler.GetPresentation)
			protected.GET("/presentation/:id/export/pptx", presentationHandler.ExportPPTX)
			protected.GET("/presentation/:id/export/pdf", presentationHandler.ExportPDF)
		}
	}

	return r
}
