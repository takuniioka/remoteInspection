package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/inspection-tool/backend/internal/domain"
	"github.com/inspection-tool/backend/internal/infrastructure/storage"
)

// InspectionService handles inspection-related business logic
type InspectionService struct {
	inspectionRepo domain.InspectionRepository
	auditLogRepo   domain.AuditLogRepository
}

// NewInspectionService creates a new inspection service
func NewInspectionService(
	inspectionRepo domain.InspectionRepository,
	auditLogRepo domain.AuditLogRepository,
) *InspectionService {
	return &InspectionService{
		inspectionRepo: inspectionRepo,
		auditLogRepo:   auditLogRepo,
	}
}

// CreateInspection creates a new inspection
func (s *InspectionService) CreateInspection(ctx context.Context, user *domain.User, title, place string) (*domain.Inspection, error) {
	inspection := &domain.Inspection{
		InspectionID: uuid.New().String(),
		Title:        title,
		Place:        place,
		Status:       domain.StatusDraft,
		CreatedAt:    time.Now(),
		CreatedBy:    user.UserID,
	}

	err := s.inspectionRepo.Create(ctx, inspection)
	if err != nil {
		return nil, fmt.Errorf("failed to create inspection: %w", err)
	}

	// Log action
	s.auditLogRepo.Create(ctx, &domain.AuditLog{
		LogID:        uuid.New().String(),
		UserID:       user.UserID,
		Action:       "create_inspection",
		InspectionID: &inspection.InspectionID,
		Timestamp:    time.Now(),
	})

	return inspection, nil
}

// GetInspection retrieves an inspection by ID
func (s *InspectionService) GetInspection(ctx context.Context, inspectionID string) (*domain.Inspection, error) {
	return s.inspectionRepo.GetByID(ctx, inspectionID)
}

// ListInspections retrieves all inspections
func (s *InspectionService) ListInspections(ctx context.Context) ([]*domain.Inspection, error) {
	return s.inspectionRepo.ListAll(ctx)
}

// StartInspection starts an inspection (draft -> running)
func (s *InspectionService) StartInspection(ctx context.Context, user *domain.User, inspectionID string) (*domain.Inspection, error) {
	inspection, err := s.inspectionRepo.GetByID(ctx, inspectionID)
	if err != nil {
		return nil, err
	}

	inspection.Status = domain.StatusRunning
	now := time.Now()
	inspection.StartedAt = &now

	err = s.inspectionRepo.Update(ctx, inspection)
	if err != nil {
		return nil, fmt.Errorf("failed to start inspection: %w", err)
	}

	// Log action
	s.auditLogRepo.Create(ctx, &domain.AuditLog{
		LogID:        uuid.New().String(),
		UserID:       user.UserID,
		Action:       "start_inspection",
		InspectionID: &inspectionID,
		Timestamp:    time.Now(),
	})

	return inspection, nil
}

// EndInspection ends an inspection (running -> closed)
func (s *InspectionService) EndInspection(ctx context.Context, user *domain.User, inspectionID string) (*domain.Inspection, error) {
	inspection, err := s.inspectionRepo.GetByID(ctx, inspectionID)
	if err != nil {
		return nil, err
	}

	inspection.Status = domain.StatusClosed
	now := time.Now()
	inspection.EndedAt = &now

	err = s.inspectionRepo.Update(ctx, inspection)
	if err != nil {
		return nil, fmt.Errorf("failed to end inspection: %w", err)
	}

	// Log action
	s.auditLogRepo.Create(ctx, &domain.AuditLog{
		LogID:        uuid.New().String(),
		UserID:       user.UserID,
		Action:       "end_inspection",
		InspectionID: &inspectionID,
		Timestamp:    time.Now(),
	})

	return inspection, nil
}

// PhotoService handles photo-related business logic
type PhotoService struct {
	photoRepo     domain.EvidencePhotoRepository
	captureRepo   domain.CaptureRequestRepository
	s3Storage     *storage.S3Storage
	auditLogRepo  domain.AuditLogRepository
}

// NewPhotoService creates a new photo service
func NewPhotoService(
	photoRepo domain.EvidencePhotoRepository,
	captureRepo domain.CaptureRequestRepository,
	s3Storage *storage.S3Storage,
	auditLogRepo domain.AuditLogRepository,
) *PhotoService {
	return &PhotoService{
		photoRepo:    photoRepo,
		captureRepo:  captureRepo,
		s3Storage:    s3Storage,
		auditLogRepo: auditLogRepo,
	}
}

// GenerateUploadPresignedURL generates a presigned URL for photo upload
func (s *PhotoService) GenerateUploadPresignedURL(ctx context.Context, user *domain.User, inspectionID, captureRequestID, fileName string) (string, error) {
	// Verify capture request exists
	_, err := s.captureRepo.GetByID(ctx, inspectionID, captureRequestID)
	if err != nil {
		return "", err
	}

	// Generate S3 key
	key := fmt.Sprintf("inspections/%s/photos/%s/%s", inspectionID, uuid.New().String(), fileName)

	// Generate presigned URL
	url, err := s.s3Storage.GeneratePresignedUploadURL(ctx, key)
	if err != nil {
		return "", err
	}

	return url, nil
}

// CompletePhotoUpload marks a photo upload as complete
func (s *PhotoService) CompletePhotoUpload(ctx context.Context, user *domain.User, inspectionID, captureRequestID, s3Key string) (*domain.EvidencePhoto, error) {
	// Update capture request
	captureReq, err := s.captureRepo.GetByID(ctx, inspectionID, captureRequestID)
	if err != nil {
		return nil, err
	}

	photoID := uuid.New().String()
	now := time.Now()
	captureReq.UploadedAt = &now
	captureReq.State = domain.CaptureUploaded
	captureReq.PhotoID = &photoID

	err = s.captureRepo.Update(ctx, captureReq)
	if err != nil {
		return nil, err
	}

	// Create evidence photo
	photo := &domain.EvidencePhoto{
		InspectionID:     inspectionID,
		PhotoID:          photoID,
		CaptureRequestID: &captureRequestID,
		S3OriginalKey:    s3Key,
		CapturedAt:       now,
		UploadedAt:       now,
		LinkedType:       captureReq.LinkedType,
		LinkedID:         captureReq.LinkedID,
	}

	err = s.photoRepo.Create(ctx, photo)
	if err != nil {
		return nil, err
	}

	// Log action
	s.auditLogRepo.Create(ctx, &domain.AuditLog{
		LogID:        uuid.New().String(),
		UserID:       user.UserID,
		Action:       "upload_photo",
		InspectionID: &inspectionID,
		Details: map[string]interface{}{
			"photoId":          photoID,
			"captureRequestId": captureRequestID,
		},
		Timestamp: time.Now(),
	})

	return photo, nil
}

// CaptureRequestService handles capture request business logic
type CaptureRequestService struct {
	captureRepo  domain.CaptureRequestRepository
	auditLogRepo domain.AuditLogRepository
}

// NewCaptureRequestService creates a new capture request service
func NewCaptureRequestService(
	captureRepo domain.CaptureRequestRepository,
	auditLogRepo domain.AuditLogRepository,
) *CaptureRequestService {
	return &CaptureRequestService{
		captureRepo:  captureRepo,
		auditLogRepo: auditLogRepo,
	}
}

// CreateCaptureRequest creates a new capture request
func (s *CaptureRequestService) CreateCaptureRequest(ctx context.Context, user *domain.User, inspectionID string, linkedType, linkedID *string) (*domain.CaptureRequest, error) {
	now := time.Now()
	captureReq := &domain.CaptureRequest{
		InspectionID:     inspectionID,
		CaptureRequestID: uuid.New().String(),
		RequestedAt:      now,
		RequestedBy:      user.UserID,
		State:            domain.CaptureRequested,
		LinkedType:       linkedType,
		LinkedID:         linkedID,
	}

	err := s.captureRepo.Create(ctx, captureReq)
	if err != nil {
		return nil, err
	}

	// Log action
	s.auditLogRepo.Create(ctx, &domain.AuditLog{
		LogID:        uuid.New().String(),
		UserID:       user.UserID,
		Action:       "create_capture_request",
		InspectionID: &inspectionID,
		Details: map[string]interface{}{
			"captureRequestId": captureReq.CaptureRequestID,
		},
		Timestamp: time.Now(),
	})

	return captureReq, nil
}

// GetCaptureRequest retrieves a capture request
func (s *CaptureRequestService) GetCaptureRequest(ctx context.Context, inspectionID, captureRequestID string) (*domain.CaptureRequest, error) {
	return s.captureRepo.GetByID(ctx, inspectionID, captureRequestID)
}

// TemplateService handles annotation template business logic
type TemplateService struct {
	templateRepo domain.AnnotationTemplateRepository
}

// NewTemplateService creates a new template service
func NewTemplateService(templateRepo domain.AnnotationTemplateRepository) *TemplateService {
	return &TemplateService{
		templateRepo: templateRepo,
	}
}

// CreateTemplate creates a new annotation template
func (s *TemplateService) CreateTemplate(ctx context.Context, name, description string, definition map[string]interface{}) (*domain.AnnotationTemplate, error) {
	template := &domain.AnnotationTemplate{
		TemplateID:     uuid.New().String(),
		Version:        1,
		Name:           name,
		Description:    description,
		IsActive:       true,
		DefinitionJSON: definition,
		CreatedAt:      time.Now(),
	}

	err := s.templateRepo.Create(ctx, template)
	if err != nil {
		return nil, err
	}

	return template, nil
}

// GetActiveTemplates retrieves all active templates
func (s *TemplateService) GetActiveTemplates(ctx context.Context) ([]*domain.AnnotationTemplate, error) {
	return s.templateRepo.GetActiveTemplates(ctx)
}

// ViewerLinkService handles viewer link business logic
type ViewerLinkService struct {
	viewerLinkRepo domain.ViewerLinkRepository
	auditLogRepo   domain.AuditLogRepository
}

// NewViewerLinkService creates a new viewer link service
func NewViewerLinkService(
	viewerLinkRepo domain.ViewerLinkRepository,
	auditLogRepo domain.AuditLogRepository,
) *ViewerLinkService {
	return &ViewerLinkService{
		viewerLinkRepo: viewerLinkRepo,
		auditLogRepo:   auditLogRepo,
	}
}

// CreateViewerLink creates a new viewer access link
func (s *ViewerLinkService) CreateViewerLink(ctx context.Context, user *domain.User, inspectionID string, expiresInDays int, maxUses *int) (*domain.ViewerLink, error) {
	token := uuid.New().String()
	expiresAt := time.Now().AddDate(0, 0, expiresInDays)

	link := &domain.ViewerLink{
		ViewerAccessToken: token,
		InspectionID:      inspectionID,
		ExpiresAt:         expiresAt,
		IsRevoked:         false,
		MaxUses:           maxUses,
		UsedCount:         0,
		CreatedAt:         time.Now(),
		CreatedBy:         user.UserID,
	}

	err := s.viewerLinkRepo.Create(ctx, link)
	if err != nil {
		return nil, err
	}

	return link, nil
}

// ValidateViewerLink validates a viewer access link
func (s *ViewerLinkService) ValidateViewerLink(ctx context.Context, token string) (*domain.ViewerLink, error) {
	link, err := s.viewerLinkRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Check if revoked
	if link.IsRevoked {
		return nil, fmt.Errorf("link has been revoked")
	}

	// Check if expired
	if time.Now().After(link.ExpiresAt) {
		return nil, fmt.Errorf("link has expired")
	}

	// Check max uses
	if link.MaxUses != nil && link.UsedCount >= *link.MaxUses {
		return nil, fmt.Errorf("link usage limit exceeded")
	}

	return link, nil
}
