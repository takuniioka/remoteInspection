package domain

import "time"

// Role represents user roles
type Role string

const (
	RoleViewer   Role = "viewer"
	RoleCapturer Role = "capturer"
	RoleEditor   Role = "editor"
	RoleAdmin    Role = "admin"
	RoleGuestViewer Role = "guestViewer"
)

// User represents a logged-in user
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

// Inspection represents an inspection project
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

// ChecklistItem represents a checklist item
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

// Issue represents an inspection issue
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

// CaptureRequest represents a photo capture request
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

// EvidencePhoto represents a captured evidence photo
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

// AnnotationTemplate represents an annotation template
type AnnotationTemplate struct {
	TemplateID     string
	Version        int
	Name           string
	Description    string
	IsActive       bool
	DefinitionJSON map[string]interface{}
	CreatedAt      time.Time
}

// ViewerLink represents a guest viewer access link
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

// AuditLog represents an audit log entry
type AuditLog struct {
	LogID        string
	UserID       string
	Action       string // "create_inspection", "capture_request", "upload_photo", etc.
	InspectionID *string
	Details      map[string]interface{}
	Timestamp    time.Time
}
