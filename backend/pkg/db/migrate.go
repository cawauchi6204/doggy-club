package db

import (
	"errors"
	"log"

	"github.com/doggyclub/backend/pkg/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	// Enable UUID extension for PostgreSQL
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return err
	}

	// Migrate all models for simplified schema
	err := db.AutoMigrate(
		// Core models
		&models.User{},
		&models.RefreshToken{},
		&models.Dog{},
		&models.Encounter{},
		&models.DeviceLocation{},
		&models.Gift{},
		&models.Post{},
		&models.Like{},
		&models.Comment{},
		&models.Hashtag{},
		&models.Follower{},
		&models.SubscriptionPlan{},
		&models.UserSubscription{},
		&models.DeviceToken{},
		&models.Notification{},
	)
	
	if err != nil {
		return err
	}

	// Create indexes for better performance
	createIndexes(db)

	log.Println("Database migration completed successfully")
	return nil
}

func createIndexes(db *gorm.DB) {
	// Modern GORM uses Migrator to create indexes
	migrator := db.Migrator()
	
	// User indexes
	migrator.CreateIndex(&models.User{}, "email")
	migrator.CreateIndex(&models.User{}, "username")
	
	// Dog indexes
	migrator.CreateIndex(&models.Dog{}, "user_id")
	migrator.CreateIndex(&models.Dog{}, "created_at")
	
	// Encounter indexes
	migrator.CreateIndex(&models.Encounter{}, "dog1_id")
	migrator.CreateIndex(&models.Encounter{}, "dog2_id")
	migrator.CreateIndex(&models.Encounter{}, "timestamp")
	
	// Device location indexes
	migrator.CreateIndex(&models.DeviceLocation{}, "dog_id")
	migrator.CreateIndex(&models.DeviceLocation{}, "updated_at")
	
	// Gift indexes
	migrator.CreateIndex(&models.Gift{}, "sender_dog_id")
	migrator.CreateIndex(&models.Gift{}, "receiver_dog_id")
	migrator.CreateIndex(&models.Gift{}, "sent_at")
	
	// Post indexes
	migrator.CreateIndex(&models.Post{}, "dog_id")
	migrator.CreateIndex(&models.Post{}, "created_at")
	
	// Like indexes
	migrator.CreateIndex(&models.Like{}, "post_id")
	migrator.CreateIndex(&models.Like{}, "dog_id")
	
	// Comment indexes
	migrator.CreateIndex(&models.Comment{}, "post_id")
	migrator.CreateIndex(&models.Comment{}, "dog_id")
	
	// Hashtag indexes
	migrator.CreateIndex(&models.Hashtag{}, "tag")
	
	// Follower indexes
	migrator.CreateIndex(&models.Follower{}, "follower_dog_id")
	migrator.CreateIndex(&models.Follower{}, "followed_dog_id")
	
	// Subscription indexes
	migrator.CreateIndex(&models.UserSubscription{}, "user_id")
	migrator.CreateIndex(&models.UserSubscription{}, "plan_id")
	migrator.CreateIndex(&models.UserSubscription{}, "status")
	
	// Device token indexes
	migrator.CreateIndex(&models.DeviceToken{}, "user_id")
	migrator.CreateIndex(&models.DeviceToken{}, "device_type")
	
	// Notification indexes
	migrator.CreateIndex(&models.Notification{}, "user_id")
	migrator.CreateIndex(&models.Notification{}, "type")
	migrator.CreateIndex(&models.Notification{}, "sent_at")
}

func SeedInitialData(db *gorm.DB) error {
	// Seed subscription plans
	plans := []models.SubscriptionPlan{
		{
			Name:           "Premium",
			Price:          999,  // $9.99
			DurationMonths: 1,
			Features:       `{"unlimited_dogs": true, "premium_support": true}`,
		},
	}

	for _, plan := range plans {
		var existingPlan models.SubscriptionPlan
		if err := db.Where("name = ?", plan.Name).First(&existingPlan).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&plan).Error; err != nil {
					return err
				}
			}
		}
	}

	log.Println("Initial data seeded successfully")
	return nil
}