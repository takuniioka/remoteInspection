package video

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/inspection-tool/backend/internal/domain"
)

// MockVideoProvider implements `domain.VideoProvider` for local
// development and testing. It simulates stage creation, token issuance
// and participant lists without communicating with a real streaming
// provider. Use this in development to avoid external dependencies.
type MockVideoProvider struct {
	stages       map[string]*MockStage
	tokens       map[string]*domain.VideoToken
	participants map[string][]*domain.Participant
}

// MockStage represents a lightweight in-memory stage used by the mock provider.
type MockStage struct {
	StageID     string
	StageName   string
	CreatedAt   time.Time
	Participants []*domain.Participant
}

// NewMockVideoProvider constructs a new mock provider instance.
func NewMockVideoProvider() domain.VideoProvider {
	return &MockVideoProvider{
		stages:       make(map[string]*MockStage),
		tokens:       make(map[string]*domain.VideoToken),
		participants: make(map[string][]*domain.Participant),
	}
}

// IssuePubToken issues a publisher token for the given stage. The mock
// generates random identifiers and records the token in-memory.
func (m *MockVideoProvider) IssuePubToken(ctx context.Context, stageID string) (*domain.VideoToken, error) {
	if _, ok := m.stages[stageID]; !ok {
		return nil, fmt.Errorf("stage not found: %s", stageID)
	}

	token := &domain.VideoToken{
		Token:         uuid.New().String(),
		ExpiresAt:     time.Now().Add(1 * time.Hour).Unix(),
		ParticipantID: uuid.New().String(),
	}

	m.tokens[token.Token] = token
	return token, nil
}

// IssueSubToken issues a subscriber token; mock behavior mirrors publisher token issuance.
func (m *MockVideoProvider) IssueSubToken(ctx context.Context, stageID string) (*domain.VideoToken, error) {
	if _, ok := m.stages[stageID]; !ok {
		return nil, fmt.Errorf("stage not found: %s", stageID)
	}

	token := &domain.VideoToken{
		Token:         uuid.New().String(),
		ExpiresAt:     time.Now().Add(1 * time.Hour).Unix(),
		ParticipantID: uuid.New().String(),
	}

	m.tokens[token.Token] = token
	return token, nil
}

// CreateStage creates a mock stage and returns its generated ID.
func (m *MockVideoProvider) CreateStage(ctx context.Context, stageName string) (string, error) {
	stageID := uuid.New().String()
	m.stages[stageID] = &MockStage{
		StageID:      stageID,
		StageName:    stageName,
		CreatedAt:    time.Now(),
		Participants: []*domain.Participant{},
	}
	m.participants[stageID] = []*domain.Participant{}
	return stageID, nil
}

// DeleteStage removes a mock stage from memory.
func (m *MockVideoProvider) DeleteStage(ctx context.Context, stageID string) error {
	if _, ok := m.stages[stageID]; !ok {
		return fmt.Errorf("stage not found: %s", stageID)
	}
	delete(m.stages, stageID)
	delete(m.participants, stageID)
	return nil
}

// ListParticipants returns the list of participants recorded for a stage.
func (m *MockVideoProvider) ListParticipants(ctx context.Context, stageID string) ([]*domain.Participant, error) {
	if _, ok := m.stages[stageID]; !ok {
		return nil, fmt.Errorf("stage not found: %s", stageID)
	}
	return m.participants[stageID], nil
}

// DisconnectParticipant simulates disconnecting a participant; mock just
// returns success.
func (m *MockVideoProvider) DisconnectParticipant(ctx context.Context, stageID, participantID string) error {
	if _, ok := m.stages[stageID]; !ok {
		return fmt.Errorf("stage not found: %s", stageID)
	}
	// In mock, just return success
	return nil
}
