package config

import (
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	AWS        AWSConfig
	API        APIConfig
	JWT        JWTConfig
	Cognito    CognitoConfig
	VideoProvider string
	Environment string
}

// AWSConfig holds AWS-related configuration
type AWSConfig struct {
	Region              string
	DynamoDBEndpoint    string
	S3Endpoint          string
	S3BucketName        string
	AccessKeyID         string
	SecretAccessKey     string
}

// APIConfig holds API-related configuration
type APIConfig struct {
	Port int
	Host string
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret string
	TTL    int // in seconds
}

// CognitoConfig holds Cognito-related configuration
type CognitoConfig struct {
	Region       string
	UserPoolID   string
	ClientID     string
	ClientSecret string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{
		AWS: AWSConfig{
			Region:           getEnv("AWS_REGION", "ap-northeast-1"),
			DynamoDBEndpoint: getEnv("DYNAMODB_ENDPOINT", ""),
			S3Endpoint:       getEnv("S3_ENDPOINT", ""),
			S3BucketName:     getEnv("S3_BUCKET_NAME", "inspection-evidence"),
			AccessKeyID:      getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey:  getEnv("AWS_SECRET_ACCESS_KEY", ""),
		},
		API: APIConfig{
			Port: getEnvInt("API_PORT", 8080),
			Host: getEnv("API_HOST", "0.0.0.0"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "change-me-in-production"),
			TTL:    getEnvInt("JWT_TTL", 3600),
		},
		Cognito: CognitoConfig{
			Region:       getEnv("COGNITO_REGION", "ap-northeast-1"),
			UserPoolID:   getEnv("COGNITO_USER_POOL_ID", ""),
			ClientID:     getEnv("COGNITO_CLIENT_ID", ""),
			ClientSecret: getEnv("COGNITO_CLIENT_SECRET", ""),
		},
		VideoProvider: getEnv("VIDEO_PROVIDER", "mock"),
		Environment:   getEnv("ENVIRONMENT", "development"),
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}
