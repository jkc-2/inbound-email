package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               int
	WebhookURL         string
	WebhookConcurrency int
	S3BucketName       string
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	MaxFileSize        int64
	SMTPSecure         bool
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		Port:               getInt("PORT", 25),
		WebhookURL:         getString("WEBHOOK_URL", "https://enkhprqr4n2t.x.pipedream.net/"),
		WebhookConcurrency: getInt("WEBHOOK_CONCURRENCY", 5),
		S3BucketName:       getString("S3_BUCKET_NAME", ""),
		AWSRegion:          getString("AWS_REGION", ""),
		AWSAccessKeyID:     getString("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: getString("AWS_SECRET_ACCESS_KEY", ""),
		MaxFileSize:        getInt64("MAX_FILE_SIZE", 5*1024*1024),
		SMTPSecure:         getBool("SMTP_SECURE", false),
	}
}

func getString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getInt64(key string, defaultValue int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
