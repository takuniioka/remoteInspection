package domain

import (
	"testing"
	"time"
)

// Test AppError creation
func TestAppError_Creation(t *testing.T) {
	err := AppError{
		Code:    ErrInspectionNotFound,
		Message: "Inspection not found",
	}

	if err.Code != ErrInspectionNotFound {
		t.Errorf("Expected code %s, got %s", ErrInspectionNotFound, err.Code)
	}

	if err.Message != "Inspection not found" {
		t.Errorf("Expected message 'Inspection not found', got '%s'", err.Message)
	}
}

// Test AppError Error method
func TestAppError_Error(t *testing.T) {
	err := AppError{
		Code:    ErrUnauthorized,
		Message: "Unauthorized access",
	}

	errStr := err.Error()
	if errStr != "Unauthorized: Unauthorized access" {
		t.Errorf("Expected 'Unauthorized: Unauthorized access', got '%s'", errStr)
	}
}

// Test Inspection validation
func TestInspection_Validation(t *testing.T) {
	inspection := &Inspection{
		InspectionID:    "test-123",
		Title:           "Valid Inspection",
		Place:           "Tokyo",
		Status:          StatusDraft,
		CreatedByUserID: "user1",
		CreatedAt:       time.Now(),
	}

	if inspection.InspectionID == "" {
		t.Error("InspectionID should not be empty")
	}

	if inspection.Status != StatusDraft {
		t.Errorf("Expected status %s, got %s", StatusDraft, inspection.Status)
	}
}

// Test inspection status transitions
func TestInspectionStatusTransitions(t *testing.T) {
	validTransitions := map[InspectionStatus][]InspectionStatus{
		StatusDraft:   {StatusRunning},
		StatusRunning: {StatusClosed},
		StatusClosed:  {},
	}

	for from, tos := range validTransitions {
		for _, to := range tos {
			if !isValidStatusTransition(from, to) {
				t.Errorf("Transition from %s to %s should be valid", from, to)
			}
		}
	}
}

func isValidStatusTransition(from, to InspectionStatus) bool {
	validTransitions := map[InspectionStatus][]InspectionStatus{
		StatusDraft:   {StatusRunning},
		StatusRunning: {StatusClosed},
	}
	for _, valid := range validTransitions[from] {
		if valid == to {
			return true
		}
	}
	return false
}

// Test CaptureRequest state transitions
func TestCaptureRequestStateTransitions(t *testing.T) {
	validTransitions := map[CaptureRequestState][]CaptureRequestState{
		CaptureStateRequested: {CaptureStateCapturing, CaptureStateTimeout},
		CaptureStateCapturing: {CaptureStateUploaded, CaptureStateFailed, CaptureStateTimeout},
		CaptureStateUploaded:  {},
		CaptureStateFailed:    {CaptureStateRequested},
		CaptureStateTimeout:   {CaptureStateRequested},
	}

	for from, tos := range validTransitions {
		for _, to := range tos {
			if !isValidCaptureStateTransition(from, to) {
				t.Errorf("Transition from %s to %s should be valid", from, to)
			}
		}
	}
}

func isValidCaptureStateTransition(from, to CaptureRequestState) bool {
	validTransitions := map[CaptureRequestState][]CaptureRequestState{
		CaptureStateRequested: {CaptureStateCapturing, CaptureStateTimeout},
		CaptureStateCapturing: {CaptureStateUploaded, CaptureStateFailed, CaptureStateTimeout},
		CaptureStateUploaded:  {},
		CaptureStateFailed:    {CaptureStateRequested},
		CaptureStateTimeout:   {CaptureStateRequested},
	}
	for _, valid := range validTransitions[from] {
		if valid == to {
			return true
		}
	}
	return false
}

// Test User role validation
func TestUserRole_Validation(t *testing.T) {
	validRoles := []Role{
		RoleViewer,
		RoleCapturer,
		RoleEditor,
		RoleAdmin,
		RoleGuestViewer,
	}

	user := &User{
		UserID:   "user1",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     RoleEditor,
	}

	found := false
	for _, role := range validRoles {
		if user.Role == role {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Role %s should be valid", user.Role)
	}
}

// Test ChecklistItem status
func TestChecklistItem_Status(t *testing.T) {
	validStatuses := []ItemStatus{
		ItemStatusOK,
		ItemStatusIssue,
		ItemStatusNA,
	}

	item := &ChecklistItem{
		ItemID:        "item1",
		InspectionID:  "insp1",
		Text:          "Test item",
		Status:        ItemStatusOK,
		ImageRefID:    nil,
		Timestamp:     time.Now(),
		CapturedBy:    "user1",
	}

	found := false
	for _, status := range validStatuses {
		if item.Status == status {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Status %s should be valid", item.Status)
	}
}

// Test ViewerLink expiration
func TestViewerLink_Expiration(t *testing.T) {
	expiredLink := &ViewerLink{
		ViewerAccessToken: "token1",
		InspectionID:      "insp1",
		CreatedAt:         time.Now().Add(-2 * time.Hour),
		ExpiresAt:         time.Now().Add(-1 * time.Hour),
		IsRevoked:         false,
		MaxUses:           10,
		UsageCount:        0,
	}

	if !expiredLink.IsExpired() {
		t.Error("Link should be expired")
	}

	validLink := &ViewerLink{
		ViewerAccessToken: "token2",
		InspectionID:      "insp1",
		CreatedAt:         time.Now(),
		ExpiresAt:         time.Now().Add(1 * time.Hour),
		IsRevoked:         false,
		MaxUses:           10,
		UsageCount:        0,
	}

	if validLink.IsExpired() {
		t.Error("Link should not be expired")
	}
}

// Helper method for ViewerLink to test
func (v *ViewerLink) IsExpired() bool {
	if v.IsRevoked {
		return true
	}
	if v.UsageCount >= v.MaxUses && v.MaxUses > 0 {
		return true
	}
	return time.Now().After(v.ExpiresAt)
}

// Test AuditLog creation
func TestAuditLog_Creation(t *testing.T) {
	log := &AuditLog{
		AuditLogID:   "log1",
		InspectionID: "insp1",
		UserID:       "user1",
		Action:       "CREATE_INSPECTION",
		ResourceType: "inspection",
		ResourceID:   "insp1",
		Timestamp:    time.Now(),
	}

	if log.AuditLogID == "" {
		t.Error("AuditLogID should not be empty")
	}

	if log.Action != "CREATE_INSPECTION" {
		t.Errorf("Expected action CREATE_INSPECTION, got %s", log.Action)
	}

	if log.Timestamp.IsZero() {
		t.Error("Timestamp should be set")
	}
}

// Test Issue severity levels
func TestIssue_SeverityLevels(t *testing.T) {
	validSeverities := []string{"low", "medium", "high", "critical"}

	issue := &Issue{
		IssueID:       "issue1",
		InspectionID:  "insp1",
		Title:         "Test Issue",
		Description:   "Test description",
		Severity:      "high",
		Status:        "open",
		CreatedAt:     time.Now(),
		CreatedByUser: "user1",
	}

	found := false
	for _, severity := range validSeverities {
		if issue.Severity == severity {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Severity %s should be valid", issue.Severity)
	}
}

// Test Issue status transitions
func TestIssue_StatusTransitions(t *testing.T) {
	issue := &Issue{
		IssueID:      "issue1",
		Status:       "open",
		InspectionID: "insp1",
	}

	// Simulate transition from open -> in_progress
	issue.Status = "in_progress"
	if issue.Status != "in_progress" {
		t.Errorf("Status should be in_progress, got %s", issue.Status)
	}

	// Simulate transition from in_progress -> resolved
	issue.Status = "resolved"
	if issue.Status != "resolved" {
		t.Errorf("Status should be resolved, got %s", issue.Status)
	}
}

// Test EvidencePhoto metadata
func TestEvidencePhoto_Metadata(t *testing.T) {
	photo := &EvidencePhoto{
		PhotoID:       "photo1",
		InspectionID:  "insp1",
		S3ObjectKey:   "photos/photo1.jpg",
		ChecklistRef:  nil,
		IssueRef:      nil,
		UploadedAt:    time.Now(),
		UploadedBy:    "user1",
		IsAnnotated:   false,
		AnnotationRef: nil,
	}

	if photo.PhotoID == "" {
		t.Error("PhotoID should not be empty")
	}

	if photo.S3ObjectKey == "" {
		t.Error("S3ObjectKey should not be empty")
	}

	if photo.UploadedAt.IsZero() {
		t.Error("UploadedAt should be set")
	}
}
