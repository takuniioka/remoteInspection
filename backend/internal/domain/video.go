package domain

import "context"

// VideoToken represents authentication/authorization information returned
// by a video streaming provider that a client (publisher/subscriber)
// uses to connect to a live stage.
type VideoToken struct {
	Token         string
	ExpiresAt     int64 // Unix timestamp
	ParticipantID string
}

// Participant represents a participant in a video stage. `State` and
// `PublishState` are provider-agnostic strings indicating connection and
// publishing status respectively.
type Participant struct {
	ParticipantID string
	State         string // "connected", "disconnected"
	PublishState  string // "live", "idle"
}

// VideoProvider defines an abstraction over third-party or mocked video
// streaming providers. Implementations must handle provider-specific
// authentication and lifecycle actions (stage creation, token issuance,
// participant management) while exposing a consistent API to services.
type VideoProvider interface {
	// IssuePubToken issues a publisher token for the stage
	IssuePubToken(ctx context.Context, stageID string) (*VideoToken, error)

	// IssueSubToken issues a subscriber token for the stage
	IssueSubToken(ctx context.Context, stageID string) (*VideoToken, error)

	// CreateStage creates a new video stage and returns its provider-specific ID
	CreateStage(ctx context.Context, stageName string) (string, error)

	// DeleteStage deletes a video stage
	DeleteStage(ctx context.Context, stageID string) error

	// ListParticipants lists all participants in a stage
	ListParticipants(ctx context.Context, stageID string) ([]*Participant, error)

	// DisconnectParticipant disconnects a participant
	DisconnectParticipant(ctx context.Context, stageID, participantID string) error
}
