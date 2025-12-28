package domain

import "context"

// VideoToken represents a video streaming token
type VideoToken struct {
	Token         string
	ExpiresAt     int64 // Unix timestamp
	ParticipantID string
}

// Participant represents a video conference participant
type Participant struct {
	ParticipantID string
	State         string // "connected", "disconnected"
	PublishState  string // "live", "idle"
}

// VideoProvider defines the interface for video streaming providers
type VideoProvider interface {
	// IssuePubToken issues a publisher token for the stage
	IssuePubToken(ctx context.Context, stageID string) (*VideoToken, error)

	// IssueSubToken issues a subscriber token for the stage
	IssueSubToken(ctx context.Context, stageID string) (*VideoToken, error)

	// CreateStage creates a new video stage
	CreateStage(ctx context.Context, stageName string) (string, error)

	// DeleteStage deletes a video stage
	DeleteStage(ctx context.Context, stageID string) error

	// ListParticipants lists all participants in a stage
	ListParticipants(ctx context.Context, stageID string) ([]*Participant, error)

	// DisconnectParticipant disconnects a participant
	DisconnectParticipant(ctx context.Context, stageID, participantID string) error
}
