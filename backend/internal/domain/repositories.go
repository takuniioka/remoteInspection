package domain

import "context"

// Repository interfaces define storage operations for domain entities.
// Implementations (e.g., DynamoDB) should satisfy these interfaces so
// the service layer can remain storage-agnostic.

// InspectionRepository defines inspection CRUD operations
type InspectionRepository interface {
	Create(ctx context.Context, inspection *Inspection) error
	GetByID(ctx context.Context, inspectionID string) (*Inspection, error)
	ListAll(ctx context.Context) ([]*Inspection, error)
	Update(ctx context.Context, inspection *Inspection) error
}

// ChecklistItemRepository defines checklist item CRUD operations
type ChecklistItemRepository interface {
	Create(ctx context.Context, item *ChecklistItem) error
	GetByID(ctx context.Context, inspectionID, itemID string) (*ChecklistItem, error)
	ListByInspection(ctx context.Context, inspectionID string) ([]*ChecklistItem, error)
	Update(ctx context.Context, item *ChecklistItem) error
}

// IssueRepository defines issue CRUD operations
type IssueRepository interface {
	Create(ctx context.Context, issue *Issue) error
	GetByID(ctx context.Context, inspectionID, issueID string) (*Issue, error)
	ListByInspection(ctx context.Context, inspectionID string) ([]*Issue, error)
	Update(ctx context.Context, issue *Issue) error
}

// CaptureRequestRepository defines capture request CRUD operations
type CaptureRequestRepository interface {
	Create(ctx context.Context, req *CaptureRequest) error
	GetByID(ctx context.Context, inspectionID, captureRequestID string) (*CaptureRequest, error)
	ListByInspection(ctx context.Context, inspectionID string) ([]*CaptureRequest, error)
	Update(ctx context.Context, req *CaptureRequest) error
}

// EvidencePhotoRepository defines evidence photo CRUD operations
type EvidencePhotoRepository interface {
	Create(ctx context.Context, photo *EvidencePhoto) error
	GetByID(ctx context.Context, inspectionID, photoID string) (*EvidencePhoto, error)
	ListByInspection(ctx context.Context, inspectionID string) ([]*EvidencePhoto, error)
	Update(ctx context.Context, photo *EvidencePhoto) error
}

// AnnotationTemplateRepository defines annotation template CRUD operations
type AnnotationTemplateRepository interface {
	Create(ctx context.Context, template *AnnotationTemplate) error
	GetByID(ctx context.Context, templateID string) (*AnnotationTemplate, error)
	GetActiveTemplates(ctx context.Context) ([]*AnnotationTemplate, error)
	Update(ctx context.Context, template *AnnotationTemplate) error
}

// ViewerLinkRepository defines viewer link CRUD operations
type ViewerLinkRepository interface {
	Create(ctx context.Context, link *ViewerLink) error
	GetByToken(ctx context.Context, token string) (*ViewerLink, error)
	Update(ctx context.Context, link *ViewerLink) error
}

// AuditLogRepository defines audit log CRUD operations
type AuditLogRepository interface {
	Create(ctx context.Context, log *AuditLog) error
	ListByInspection(ctx context.Context, inspectionID string) ([]*AuditLog, error)
	ListByUser(ctx context.Context, userID string) ([]*AuditLog, error)
}
