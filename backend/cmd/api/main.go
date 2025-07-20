package main

import (
	"fmt"
	"log"

	"github.com/doggyclub/backend/config"
	"github.com/doggyclub/backend/pkg/db"
	"github.com/doggyclub/backend/pkg/handlers"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database
	database, err := db.InitPostgres(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize Redis
	redisClient, err := db.InitRedis(cfg.Redis)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	// Create Echo instance
	e := echo.New()

	// Cache middleware disabled for simplified schema
	// cacheService := services.NewCacheService(redisClient, *cfg)
	// cacheMiddleware := middleware.NewCacheMiddleware(cacheService)

	// Middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())
	// Cache middleware disabled for simplified schema
	// e.Use(cacheMiddleware.ConditionalCache())
	// e.Use(cacheMiddleware.ETagMiddleware())
	// e.Use(cacheMiddleware.RateLimitMiddleware(100, 15*time.Minute)) // 100 requests per 15 minutes

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status":      "healthy",
			"environment": cfg.Server.Environment,
		})
	})

	// Register handlers
	authHandler := handlers.NewAuthHandler(database, redisClient, *cfg)
	authHandler.RegisterRoutes(e)

	userHandler := handlers.NewUserHandler(database, redisClient, *cfg)
	userHandler.RegisterRoutes(e)

	dogHandler := handlers.NewDogHandler(database, redisClient, *cfg)
	dogHandler.RegisterRoutes(e)

	encounterHandler := handlers.NewEncounterHandler(database, redisClient, *cfg)
	encounterHandler.RegisterRoutes(e)

	postHandler := handlers.NewPostHandler(database, redisClient, *cfg)
	postHandler.RegisterRoutes(e)

	giftHandler := handlers.NewGiftHandler(database, redisClient, *cfg)
	giftHandler.RegisterRoutes(e)

	notificationHandler := handlers.NewNotificationHandler(database, redisClient, *cfg)
	notificationHandler.RegisterRoutes(e)

	subscriptionHandler := handlers.NewSubscriptionHandler(database, redisClient, *cfg)
	subscriptionHandler.RegisterRoutes(e)

	moderationHandler := handlers.NewModerationHandler(database, redisClient, *cfg)
	moderationHandler.RegisterRoutes(e)

	// Start server
	log.Printf("Starting server on port %s in %s mode", cfg.Server.Port, cfg.Server.Environment)
	log.Printf("Database: %s@%s:%s/%s", cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	log.Printf("Redis: %s:%s", cfg.Redis.Host, cfg.Redis.Port)

	// Close database connections when server stops
	defer func() {
		if sqlDB, err := database.DB(); err == nil {
			sqlDB.Close()
		}
		redisClient.Close()
	}()

	// Start server
	if err := e.Start(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
