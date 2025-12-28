package domain

import "time"

// Package domain contains core domain model types used across the
// backend. These types represent persisted entities (Inspections,
// Issues, Photos, etc.) and value objects (statuses, roles) used by the
// service and repository layers. Keep these types small and serializable
// so they can be easily converted to/from DynamoDB items or JSON API
// payloads.

// Role represents user roles
type Role string

const (
	RoleViewer   Role = "viewer"
	RoleCapturer Role = "capturer"
	RoleEditor   Role = "editor"
	RoleAdmin    Role = "admin"
	RoleGuestViewer Role = "guestViewer"
)

// User represents an authenticated user. This model stores minimal
// identity information used for authorization and audit trails.
type User struct {
	UserID    string
	Username  string
	Email     string
	Role      Role
	CreatedAt time.Time
}

// InspectionStatus represents inspection statuses
type InspectionStatus string

const (
	StatusDraft   InspectionStatus = "draft"
	StatusRunning InspectionStatus = "running"
	StatusClosed  InspectionStatus = "closed"
)

// Inspection represents an inspection project and its lifecycle
// information. `RealtimeRef` holds provider-specific references (e.g.,
// IVS stage id) used to coordinate live capture sessions.
type Inspection struct {
	InspectionID string
	Title        string
	Place        string
	Status       InspectionStatus
	RealtimeRef  string // IVS Stage ID or other
	CreatedAt    time.Time
	StartedAt    *time.Time
	EndedAt      *time.Time
	CreatedBy    string
}

// ChecklistItem represents a single checklist entry within an
// `Inspection`. `Result` is a small controlled vocabulary (e.g., "ok",
// "ng", "pending").
type ChecklistItem struct {
	InspectionID string
	ItemID       string
	Label        string
	Result       string // "ok", "ng", "pending"
	Comment      string
	UpdatedAt    time.Time
	UpdatedBy    string
}

// IssueStatus represents issue statuses
type IssueStatus string

const (
	IssueOpen     IssueStatus = "open"
	IssueResolved IssueStatus = "resolved"
	IssueClosed   IssueStatus = "closed"
)

// Issue represents an issue discovered during an inspection. `Assignee`
// is optional. Use `Status` to track workflow transitions (open,
// resolved, closed).
type Issue struct {
	InspectionID string
	IssueID      string
	Title        string
	Detail       string
	Status       IssueStatus
	Assignee     *string
	CreatedAt    time.Time
	CreatedBy    string
	UpdatedAt    time.Time
	UpdatedBy    string
}

// CaptureRequestState represents capture request states
type CaptureRequestState string

const (
	CaptureRequested CaptureRequestState = "requested"
	CaptureCapturing CaptureRequestState = "capturing"
	CaptureUploaded  CaptureRequestState = "uploaded"
	CaptureFailed    CaptureRequestState = "failed"
	CaptureTimeout   CaptureRequestState = "timeout"
)

// CaptureRequest represents a request to capture a photo from a
// field device. The lifecycle is tracked in `State`. `LinkedType` and
// `LinkedID` allow associating the photo to a checklist item or an
// issue.
type CaptureRequest struct {
	InspectionID      string
	CaptureRequestID  string
	RequestedAt       time.Time
	RequestedBy       string
	State             CaptureRequestState
	CapturedAt        *time.Time
	UploadedAt        *time.Time
	PhotoID           *string
	LinkedType        *string // "checklistItem", "issue"
	LinkedID          *string
	FailReason        *string
}

// EvidencePhoto stores metadata about an uploaded photo and the S3
// object keys for original and annotated versions. `S3AnnotationJSONKey`
// can hold annotation metadata (shapes, labels) in JSON form.
type EvidencePhoto struct {
	InspectionID       string
	PhotoID            string
	CaptureRequestID   *string
	S3OriginalKey      string
	S3AnnotatedKey     *string
	S3AnnotationJSONKey *string
	CapturedAt         time.Time
	UploadedAt         time.Time
	AnnotatedAt        *time.Time
	LinkedType         *string // "checklistItem", "issue"
	LinkedID           *string
	Note               string
}

// AnnotationTemplate contains reusable annotation definitions that can
// be applied to captured photos. `DefinitionJSON` is a free-form
// structure describing shapes, labels and other metadata used by the
// frontend annotator.
type AnnotationTemplate struct {
	TemplateID     string
	Version        int
	Name           string
	Description    string
	IsActive       bool
	DefinitionJSON map[string]interface{}
	CreatedAt      time.Time
}

// ViewerLink represents a temporary guest access token allowing a
// third-party to view an inspection. `MaxUses` and `ExpiresAt` control
// access lifetime.
type ViewerLink struct {
	ViewerAccessToken string
	InspectionID      string
	ExpiresAt         time.Time
	IsRevoked         bool
	MaxUses           *int
	UsedCount         int
	CreatedAt         time.Time
	CreatedBy         string
}

// AuditLog records actions performed by users for compliance and
// troubleshooting. `Details` can store arbitrary structured metadata
// about the action.
type AuditLog struct {
	LogID        string
	UserID       string
	Action       string // "create_inspection", "capture_request", "upload_photo", etc.
	InspectionID *string
	Details      map[string]interface{}
	Timestamp    time.Time
}
