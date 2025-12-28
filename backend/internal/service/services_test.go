package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/inspection-tool/backend/internal/domain"
)

// This test file contains unit tests for service layer logic along with
// simple in-memory mock implementations of repository interfaces. The
// mocks are intentionally lightweight and synchronous to make tests
// deterministic and fast.

// MockInspectionRepository for testing
type MockInspectionRepository struct {
	inspections map[string]*domain.Inspection
	shouldFail  bool
}

func (m *MockInspectionRepository) Create(ctx context.Context, inspection *domain.Inspection) error {
	if m.shouldFail {
		return errors.New("mock error")
	}
	if m.inspections == nil {
		m.inspections = make(map[string]*domain.Inspection)
	}
	m.inspections[inspection.InspectionID] = inspection
	return nil
}

func (m *MockInspectionRepository) GetByID(ctx context.Context, id string) (*domain.Inspection, error) {
	if m.shouldFail {
		return nil, domain.AppError{Code: domain.ErrInspectionNotFound}
	}
	if insp, ok := m.inspections[id]; ok {
		return insp, nil
	}
	return nil, domain.AppError{Code: domain.ErrInspectionNotFound}
}

func (m *MockInspectionRepository) List(ctx context.Context, userID string, limit int) ([]*domain.Inspection, error) {
	if m.shouldFail {
		return nil, errors.New("mock error")
	}
	var result []*domain.Inspection
	for _, insp := range m.inspections {
		if insp.CreatedByUserID == userID {
			result = append(result, insp)
		}
	}
	return result, nil
}

func (m *MockInspectionRepository) Update(ctx context.Context, inspection *domain.Inspection) error {
	if m.shouldFail {
		return errors.New("mock error")
	}
	if _, ok := m.inspections[inspection.InspectionID]; ok {
		m.inspections[inspection.InspectionID] = inspection
		return nil
	}
	return domain.AppError{Code: domain.ErrInspectionNotFound}
}

// MockChecklistItemRepository for testing
type MockChecklistItemRepository struct {
	items      map[string]map[string]*domain.ChecklistItem
	shouldFail bool
}

func (m *MockChecklistItemRepository) Create(ctx context.Context, item *domain.ChecklistItem) error {
	if m.shouldFail {
		return errors.New("mock error")
	}
	if m.items == nil {
		m.items = make(map[string]map[string]*domain.ChecklistItem)
	}
	if m.items[item.InspectionID] == nil {
		m.items[item.InspectionID] = make(map[string]*domain.ChecklistItem)
	}
	m.items[item.InspectionID][item.ItemID] = item
	return nil
}

func (m *MockChecklistItemRepository) GetByID(ctx context.Context, inspectionID, itemID string) (*domain.ChecklistItem, error) {
	if m.shouldFail {
		return nil, errors.New("mock error")
	}
	if items, ok := m.items[inspectionID]; ok {
		if item, ok := items[itemID]; ok {
			return item, nil
		}
	}
	return nil, domain.AppError{Code: domain.ErrChecklistItemNotFound}
}

func (m *MockChecklistItemRepository) ListByInspection(ctx context.Context, inspectionID string) ([]*domain.ChecklistItem, error) {
	if m.shouldFail {
		return nil, errors.New("mock error")
	}
	var result []*domain.ChecklistItem
	if items, ok := m.items[inspectionID]; ok {
		for _, item := range items {
			result = append(result, item)
		}
	}
	return result, nil
}

func (m *MockChecklistItemRepository) Update(ctx context.Context, item *domain.ChecklistItem) error {
	if m.shouldFail {
		return errors.New("mock error")
	}
	if items, ok := m.items[item.InspectionID]; ok {
		if _, ok := items[item.ItemID]; ok {
			items[item.ItemID] = item
			return nil
		}
	}
	return domain.AppError{Code: domain.ErrChecklistItemNotFound}
}

// MockIssueRepository for testing
type MockIssueRepository struct {
	issues     map[string]map[string]*domain.Issue
	shouldFail bool
}

func (m *MockIssueRepository) Create(ctx context.Context, issue *domain.Issue) error {
	if m.shouldFail {
		return errors.New("mock error")
	}
	if m.issues == nil {
		m.issues = make(map[string]map[string]*domain.Issue)
	}
	if m.issues[issue.InspectionID] == nil {
		m.issues[issue.InspectionID] = make(map[string]*domain.Issue)
	}
	m.issues[issue.InspectionID][issue.IssueID] = issue
	return nil
}

func (m *MockIssueRepository) GetByID(ctx context.Context, inspectionID, issueID string) (*domain.Issue, error) {
	if m.shouldFail {
		return nil, errors.New("mock error")
	}
	if issues, ok := m.issues[inspectionID]; ok {
		if issue, ok := issues[issueID]; ok {
			return issue, nil
		}
	}
	return nil, domain.AppError{Code: domain.ErrIssueNotFound}
}

func (m *MockIssueRepository) ListByInspection(ctx context.Context, inspectionID string) ([]*domain.Issue, error) {
	if m.shouldFail {
		return nil, errors.New("mock error")
	}
	var result []*domain.Issue
	if issues, ok := m.issues[inspectionID]; ok {
		for _, issue := range issues {
			result = append(result, issue)
		}
	}
	return result, nil
}

func (m *MockIssueRepository) Update(ctx context.Context, issue *domain.Issue) error {
	if m.shouldFail {
		return errors.New("mock error")
	}
	if issues, ok := m.issues[issue.InspectionID]; ok {
		if _, ok := issues[issue.IssueID]; ok {
			issues[issue.IssueID] = issue
			return nil
		}
	}
	return domain.AppError{Code: domain.ErrIssueNotFound}
}

// MockCaptureRequestRepository for testing
type MockCaptureRequestRepository struct {
	requests   map[string]*domain.CaptureRequest
	shouldFail bool
}

func (m *MockCaptureRequestRepository) Create(ctx context.Context, req *domain.CaptureRequest) error {
	if m.shouldFail {
		return errors.New("mock error")
	}
	if m.requests == nil {
		m.requests = make(map[string]*domain.CaptureRequest)
	}
	m.requests[req.CaptureRequestID] = req
	return nil
}

func (m *MockCaptureRequestRepository) GetByID(ctx context.Context, inspectionID, captureRequestID string) (*domain.CaptureRequest, error) {
	if m.shouldFail {
		return nil, errors.New("mock error")
	}
	if req, ok := m.requests[captureRequestID]; ok && req.InspectionID == inspectionID {
		return req, nil
	}
	return nil, domain.AppError{Code: domain.ErrCaptureRequestNotFound}
}

func (m *MockCaptureRequestRepository) ListByInspection(ctx context.Context, inspectionID string) ([]*domain.CaptureRequest, error) {
	if m.shouldFail {
		return nil, errors.New("mock error")
	}
	var result []*domain.CaptureRequest
	for _, req := range m.requests {
		if req.InspectionID == inspectionID {
			result = append(result, req)
		}
	}
	return result, nil
}

func (m *MockCaptureRequestRepository) Update(ctx context.Context, req *domain.CaptureRequest) error {
	if m.shouldFail {
		return errors.New("mock error")
	}
	if _, ok := m.requests[req.CaptureRequestID]; ok {
		m.requests[req.CaptureRequestID] = req
		return nil
	}
	return domain.AppError{Code: domain.ErrCaptureRequestNotFound}
}

// Test cases for InspectionService
func TestInspectionService_CreateInspection(t *testing.T) {
	repo := &MockInspectionRepository{}
	auditRepo := &MockAuditLogRepository{}
	service := NewInspectionService(repo, auditRepo)

	ctx := context.Background()
	user := &domain.User{
		UserID:   "user1",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     domain.RoleEditor,
	}

	inspection, err := service.CreateInspection(ctx, "Test Inspection", "Tokyo", user)
	if err != nil {
		t.Fatalf("CreateInspection failed: %v", err)
	}

	if inspection.InspectionID == "" {
		t.Error("InspectionID should not be empty")
	}
	if inspection.Status != domain.StatusDraft {
		t.Errorf("Status should be draft, got %s", inspection.Status)
	}
	if inspection.CreatedByUserID != "user1" {
		t.Errorf("CreatedByUserID should be user1, got %s", inspection.CreatedByUserID)
	}
}

func TestInspectionService_StartInspection(t *testing.T) {
	repo := &MockInspectionRepository{}
	auditRepo := &MockAuditLogRepository{}
	service := NewInspectionService(repo, auditRepo)

	ctx := context.Background()

	// Create initial inspection
	inspection := &domain.Inspection{
		InspectionID:   uuid.New().String(),
		Title:          "Test",
		Place:          "Tokyo",
		Status:         domain.StatusDraft,
		CreatedByUserID: "user1",
		CreatedAt:      time.Now(),
	}
	repo.Create(ctx, inspection)

	user := &domain.User{UserID: "user1", Role: domain.RoleEditor}

	// Start inspection
	updated, err := service.StartInspection(ctx, inspection.InspectionID, user)
	if err != nil {
		t.Fatalf("StartInspection failed: %v", err)
	}

	if updated.Status != domain.StatusRunning {
		t.Errorf("Status should be running, got %s", updated.Status)
	}
	if updated.StartedAt == nil {
		t.Error("StartedAt should be set")
	}
}

func TestInspectionService_ListInspections(t *testing.T) {
	repo := &MockInspectionRepository{}
	auditRepo := &MockAuditLogRepository{}
	service := NewInspectionService(repo, auditRepo)

	ctx := context.Background()
	user := &domain.User{UserID: "user1", Role: domain.RoleViewer}

	// Add some inspections
	for i := 0; i < 3; i++ {
		inspection := &domain.Inspection{
			InspectionID:    uuid.New().String(),
			Title:           "Test " + string(rune(i)),
			Place:           "Tokyo",
			Status:          domain.StatusDraft,
			CreatedByUserID: "user1",
			CreatedAt:       time.Now(),
		}
		repo.Create(ctx, inspection)
	}

	inspections, err := service.ListInspections(ctx, user, 10)
	if err != nil {
		t.Fatalf("ListInspections failed: %v", err)
	}

	if len(inspections) != 3 {
		t.Errorf("Expected 3 inspections, got %d", len(inspections))
	}
}

// MockAuditLogRepository for testing
type MockAuditLogRepository struct {
	logs       []*domain.AuditLog
	shouldFail bool
}

func (m *MockAuditLogRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	if m.shouldFail {
		return errors.New("mock error")
	}
	m.logs = append(m.logs, log)
	return nil
}

func (m *MockAuditLogRepository) List(ctx context.Context, inspectionID string, limit int) ([]*domain.AuditLog, error) {
	if m.shouldFail {
		return nil, errors.New("mock error")
	}
	var result []*domain.AuditLog
	for _, log := range m.logs {
		if log.InspectionID == inspectionID {
			result = append(result, log)
		}
	}
	return result, nil
}

// Test PhotoService
func TestPhotoService_GenerateUploadPresignedURL(t *testing.T) {
	s3Storage := NewMockS3Storage()
	photoRepo := &MockEvidencePhotoRepository{}
	service := NewPhotoService(s3Storage, photoRepo)

	ctx := context.Background()
	inspectionID := uuid.New().String()

	url, err := service.GenerateUploadPresignedURL(ctx, inspectionID, "test.jpg")
	if err != nil {
		t.Fatalf("GenerateUploadPresignedURL failed: %v", err)
	}

	if url == "" {
		t.Error("Presigned URL should not be empty")
	}
}

// MockEvidencePhotoRepository for testing
type MockEvidencePhotoRepository struct {
	photos     map[string]map[string]*domain.EvidencePhoto
	shouldFail bool
}

func (m *MockEvidencePhotoRepository) Create(ctx context.Context, photo *domain.EvidencePhoto) error {
	if m.shouldFail {
		return errors.New("mock error")
	}
	if m.photos == nil {
		m.photos = make(map[string]map[string]*domain.EvidencePhoto)
	}
	if m.photos[photo.InspectionID] == nil {
		m.photos[photo.InspectionID] = make(map[string]*domain.EvidencePhoto)
	}
	m.photos[photo.InspectionID][photo.PhotoID] = photo
	return nil
}

func (m *MockEvidencePhotoRepository) GetByID(ctx context.Context, inspectionID, photoID string) (*domain.EvidencePhoto, error) {
	if m.shouldFail {
		return nil, errors.New("mock error")
	}
	if photos, ok := m.photos[inspectionID]; ok {
		if photo, ok := photos[photoID]; ok {
			return photo, nil
		}
	}
	return nil, domain.AppError{Code: domain.ErrPhotoNotFound}
}

func (m *MockEvidencePhotoRepository) ListByInspection(ctx context.Context, inspectionID string) ([]*domain.EvidencePhoto, error) {
	if m.shouldFail {
		return nil, errors.New("mock error")
	}
	var result []*domain.EvidencePhoto
	if photos, ok := m.photos[inspectionID]; ok {
		for _, photo := range photos {
			result = append(result, photo)
		}
	}
	return result, nil
}

func (m *MockEvidencePhotoRepository) Update(ctx context.Context, photo *domain.EvidencePhoto) error {
	if m.shouldFail {
		return errors.New("mock error")
	}
	if photos, ok := m.photos[photo.InspectionID]; ok {
		if _, ok := photos[photo.PhotoID]; ok {
			photos[photo.PhotoID] = photo
			return nil
		}
	}
	return domain.AppError{Code: domain.ErrPhotoNotFound}
}

// MockS3Storage for testing
type MockS3Storage struct {
	urls       map[string]string
	shouldFail bool
}

func NewMockS3Storage() *MockS3Storage {
	return &MockS3Storage{
		urls: make(map[string]string),
	}
}

func (m *MockS3Storage) GeneratePresignedUploadURL(ctx context.Context, bucketName, objectKey string, expiresIn int) (string, error) {
	if m.shouldFail {
		return "", errors.New("mock error")
	}
	url := "https://mock-s3.amazonaws.com/" + bucketName + "/" + objectKey
	m.urls[objectKey] = url
	return url, nil
}

func (m *MockS3Storage) GeneratePresignedDownloadURL(ctx context.Context, bucketName, objectKey string, expiresIn int) (string, error) {
	if m.shouldFail {
		return "", errors.New("mock error")
	}
	return "https://mock-s3.amazonaws.com/" + bucketName + "/" + objectKey, nil
}

func (m *MockS3Storage) ObjectExists(ctx context.Context, bucketName, objectKey string) (bool, error) {
	if m.shouldFail {
		return false, errors.New("mock error")
	}
	_, exists := m.urls[objectKey]
	return exists, nil
}

func (m *MockS3Storage) DeleteObject(ctx context.Context, bucketName, objectKey string) error {
	if m.shouldFail {
		return errors.New("mock error")
	}
	delete(m.urls, objectKey)
	return nil
}

func (m *MockS3Storage) GetObject(ctx context.Context, bucketName, objectKey string) ([]byte, error) {
	if m.shouldFail {
		return nil, errors.New("mock error")
	}
	return []byte("mock file content"), nil
}

// Test CaptureRequestService
func TestCaptureRequestService_CreateCaptureRequest(t *testing.T) {
	repo := &MockCaptureRequestRepository{}
	auditRepo := &MockAuditLogRepository{}
	service := NewCaptureRequestService(repo, auditRepo)

	ctx := context.Background()
	user := &domain.User{UserID: "user1", Role: domain.RoleCapturer}

	request, err := service.CreateCaptureRequest(ctx, "inspection1", "device1", user)
	if err != nil {
		t.Fatalf("CreateCaptureRequest failed: %v", err)
	}

	if request.CaptureRequestID == "" {
		t.Error("CaptureRequestID should not be empty")
	}
	if request.State != domain.CaptureStateRequested {
		t.Errorf("State should be requested, got %s", request.State)
	}
	if request.RequestedAt == nil {
		t.Error("RequestedAt should be set")
	}
}

// Test ViewerLinkService
func TestViewerLinkService_CreateViewerLink(t *testing.T) {
	repo := &MockViewerLinkRepository{}
	auditRepo := &MockAuditLogRepository{}
	service := NewViewerLinkService(repo, auditRepo)

	ctx := context.Background()
	user := &domain.User{UserID: "user1", Role: domain.RoleAdmin}

	link, err := service.CreateViewerLink(ctx, "inspection1", user)
	if err != nil {
		t.Fatalf("CreateViewerLink failed: %v", err)
	}

	if link.ViewerAccessToken == "" {
		t.Error("ViewerAccessToken should not be empty")
	}
	if link.InspectionID != "inspection1" {
		t.Errorf("InspectionID should be inspection1, got %s", link.InspectionID)
	}
}

func TestViewerLinkService_ValidateViewerLink(t *testing.T) {
	repo := &MockViewerLinkRepository{}
	auditRepo := &MockAuditLogRepository{}
	service := NewViewerLinkService(repo, auditRepo)

	ctx := context.Background()

	// Create a valid link
	link := &domain.ViewerLink{
		ViewerAccessToken: "test-token",
		InspectionID:      "inspection1",
		CreatedAt:         time.Now(),
		ExpiresAt:         time.Now().Add(time.Hour),
		IsRevoked:         false,
		MaxUses:           10,
		UsageCount:        0,
	}
	repo.links[link.ViewerAccessToken] = link

	// Validate should succeed
	valid, err := service.ValidateViewerLink(ctx, "test-token")
	if err != nil {
		t.Fatalf("ValidateViewerLink failed: %v", err)
	}

	if !valid {
		t.Error("Link should be valid")
	}

	// Update usage
	link.UsageCount++
	repo.links["test-token"] = link

	// Validate should still work
	valid, err = service.ValidateViewerLink(ctx, "test-token")
	if err != nil {
		t.Fatalf("ValidateViewerLink with usage failed: %v", err)
	}

	if !valid {
		t.Error("Link should still be valid after usage increment")
	}
}

// MockViewerLinkRepository for testing
type MockViewerLinkRepository struct {
	links      map[string]*domain.ViewerLink
	shouldFail bool
}

func (m *MockViewerLinkRepository) Create(ctx context.Context, link *domain.ViewerLink) error {
	if m.shouldFail {
		return errors.New("mock error")
	}
	if m.links == nil {
		m.links = make(map[string]*domain.ViewerLink)
	}
	m.links[link.ViewerAccessToken] = link
	return nil
}

func (m *MockViewerLinkRepository) GetByToken(ctx context.Context, token string) (*domain.ViewerLink, error) {
	if m.shouldFail {
		return nil, errors.New("mock error")
	}
	if link, ok := m.links[token]; ok {
		return link, nil
	}
	return nil, domain.AppError{Code: domain.ErrViewerLinkNotFound}
}

func (m *MockViewerLinkRepository) ListByInspection(ctx context.Context, inspectionID string) ([]*domain.ViewerLink, error) {
	if m.shouldFail {
		return nil, errors.New("mock error")
	}
	var result []*domain.ViewerLink
	for _, link := range m.links {
		if link.InspectionID == inspectionID {
			result = append(result, link)
		}
	}
	return result, nil
}

func (m *MockViewerLinkRepository) Update(ctx context.Context, link *domain.ViewerLink) error {
	if m.shouldFail {
		return errors.New("mock error")
	}
	if _, ok := m.links[link.ViewerAccessToken]; ok {
		m.links[link.ViewerAccessToken] = link
		return nil
	}
	return domain.AppError{Code: domain.ErrViewerLinkNotFound}
}
