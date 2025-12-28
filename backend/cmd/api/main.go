package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	appConfig "github.com/inspection-tool/backend/internal/config"
	"github.com/inspection-tool/backend/internal/infrastructure/storage"
	"github.com/inspection-tool/backend/internal/infrastructure/video"
	"github.com/inspection-tool/backend/internal/middleware"
	"github.com/inspection-tool/backend/internal/repository"
	"github.com/inspection-tool/backend/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	_ = godotenv.Load()

	// Setup logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := appConfig.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize AWS SDK
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.AWS.Region),
	)
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	// Initialize DynamoDB client
	var dynamoDBClient *dynamodb.Client
	if cfg.AWS.DynamoDBEndpoint != "" {
		// Local development
		awsCfg.Credentials.Retrieve(context.Background())
		dynamoDBClient = dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = &cfg.AWS.DynamoDBEndpoint
		})
	} else {
		// Production AWS
		dynamoDBClient = dynamodb.NewFromConfig(awsCfg)
	}

	// Initialize S3 client
	var s3Client *s3.Client
	if cfg.AWS.S3Endpoint != "" {
		// Local development (MinIO)
		s3Client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
			o.BaseEndpoint = &cfg.AWS.S3Endpoint
			o.UsePathStyle = true
		})
	} else {
		// Production AWS
		s3Client = s3.NewFromConfig(awsCfg)
	}

	// Initialize repositories
	dbClient := repository.NewDynamoDBClient(dynamoDBClient)
	inspectionRepo := repository.NewInspectionRepository(dbClient)
	checklistRepo := repository.NewChecklistItemRepository(dbClient)

	// Initialize storage
	s3Storage := storage.NewS3Storage(s3Client, cfg.AWS.S3BucketName)

	// Initialize services
	inspectionService := service.NewInspectionService(inspectionRepo, nil)
	photoService := service.NewPhotoService(nil, nil, s3Storage, nil)

	// Initialize video provider
	var videoProvider interface{}
	if cfg.VideoProvider == "mock" {
		videoProvider = video.NewMockVideoProvider()
	} else {
		// TODO: Implement IVS Real-Time provider
		videoProvider = video.NewMockVideoProvider()
	}

	_ = inspectionService
	_ = photoService
	_ = videoProvider

	// Create HTTP router
	r := chi.NewRouter()

	// Add CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Add logging middleware
	r.Use(middleware.LoggingMiddleware(logger))

	// Public routes
	r.Post("/api/public/viewer-sessions", func(w http.ResponseWriter, r *http.Request) {
		middleware.WriteSuccessResponse(w, http.StatusOK, map[string]string{
			"guestSessionToken": "mock-token",
			"expiresIn":         "3600",
		})
	})

	// Protected routes (require auth)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg.JWT.Secret))

		// Health check
		r.Get("/api/me", func(w http.ResponseWriter, r *http.Request) {
			middleware.WriteSuccessResponse(w, http.StatusOK, map[string]string{
				"userId":   "mock-user",
				"username": "mock",
				"email":    "mock@example.com",
				"role":     "admin",
			})
		})

		// Inspection endpoints
		r.Post("/api/inspections", func(w http.ResponseWriter, r *http.Request) {
			middleware.WriteSuccessResponse(w, http.StatusCreated, map[string]string{
				"inspectionId": "mock-id",
				"title":        "Test Inspection",
				"status":       "draft",
			})
		})

		r.Get("/api/inspections/{inspectionId}", func(w http.ResponseWriter, r *http.Request) {
			middleware.WriteSuccessResponse(w, http.StatusOK, map[string]string{
				"inspectionId": "mock-id",
				"title":        "Test Inspection",
				"status":       "draft",
			})
		})
	})

	// Start server
	logger.Info("Starting API server", slog.String("addr", cfg.API.Host+":"+string(rune(cfg.API.Port))))
	if err := http.ListenAndServe(":"+string(rune(cfg.API.Port)), r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
