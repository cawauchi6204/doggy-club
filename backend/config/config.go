package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Firebase FirebaseConfig
	External ExternalConfig
	Features FeatureConfig
}

type ServerConfig struct {
	Port        string
	Environment string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret              string
	ExpireHours         time.Duration
	RefreshExpireHours  time.Duration
}

type FirebaseConfig struct {
	CredentialsPath string
	ProjectID       string
}

type ExternalConfig struct {
	CloudflareR2AccessKey  string
	CloudflareR2SecretKey  string
	CloudflareR2Bucket     string
	CloudflareR2Endpoint   string
	FirebaseProjectID      string
	FirebasePrivateKey     string
	FirebaseClientEmail    string
	StripeSecretKey        string
	StripeWebhookSecret    string
	GoogleMapsAPIKey       string
}

type FeatureConfig struct {
	EnableEncounterDetection bool
	EnablePushNotifications  bool
	EnablePremiumFeatures    bool
}

func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	jwtExpire, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))
	jwtRefreshExpire, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRE_HOURS", "168"))

	return &Config{
		Server: ServerConfig{
			Port:        getEnv("PORT", "9090"),
			Environment: getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "doggyclub"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "doggyclub_db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		JWT: JWTConfig{
			Secret:             getEnv("JWT_SECRET", ""),
			ExpireHours:        time.Duration(jwtExpire) * time.Hour,
			RefreshExpireHours: time.Duration(jwtRefreshExpire) * time.Hour,
		},
		Firebase: FirebaseConfig{
			CredentialsPath: getEnv("FIREBASE_CREDENTIALS_PATH", ""),
			ProjectID:       getEnv("FIREBASE_PROJECT_ID", ""),
		},
		External: ExternalConfig{
			CloudflareR2AccessKey:  getEnv("CLOUDFLARE_R2_ACCESS_KEY", ""),
			CloudflareR2SecretKey:  getEnv("CLOUDFLARE_R2_SECRET_KEY", ""),
			CloudflareR2Bucket:     getEnv("CLOUDFLARE_R2_BUCKET", ""),
			CloudflareR2Endpoint:   getEnv("CLOUDFLARE_R2_ENDPOINT", ""),
			FirebaseProjectID:      getEnv("FIREBASE_PROJECT_ID", ""),
			FirebasePrivateKey:     getEnv("FIREBASE_PRIVATE_KEY", ""),
			FirebaseClientEmail:    getEnv("FIREBASE_CLIENT_EMAIL", ""),
			StripeSecretKey:        getEnv("STRIPE_SECRET_KEY", ""),
			StripeWebhookSecret:    getEnv("STRIPE_WEBHOOK_SECRET", ""),
			GoogleMapsAPIKey:       getEnv("GOOGLE_MAPS_API_KEY", ""),
		},
		Features: FeatureConfig{
			EnableEncounterDetection: getEnvBool("ENABLE_ENCOUNTER_DETECTION", true),
			EnablePushNotifications:  getEnvBool("ENABLE_PUSH_NOTIFICATIONS", true),
			EnablePremiumFeatures:    getEnvBool("ENABLE_PREMIUM_FEATURES", true),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		boolValue, err := strconv.ParseBool(value)
		if err == nil {
			return boolValue
		}
	}
	return defaultValue
}