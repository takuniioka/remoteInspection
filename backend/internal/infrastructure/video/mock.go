package video

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/inspection-tool/backend/internal/domain"
)

// MockVideoProvider implements domain.VideoProvider for testing/local development
type MockVideoProvider struct {
	stages       map[string]*MockStage
	tokens       map[string]*domain.VideoToken
	participants map[string][]*domain.Participant
}

// MockStage represents a mock video stage
type MockStage struct {
	StageID     string
	StageName   string
	CreatedAt   time.Time
	Participants []*domain.Participant
}

// NewMockVideoProvider creates a new mock video provider
func NewMockVideoProvider() domain.VideoProvider {
	return &MockVideoProvider{
		stages:       make(map[string]*MockStage),
		tokens:       make(map[string]*domain.VideoToken),
		participants: make(map[string][]*domain.Participant),
	}
}

// IssuePubToken issues a publisher token
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

// IssueSubToken issues a subscriber token
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

// CreateStage creates a new mock video stage
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

// DeleteStage deletes a mock video stage
func (m *MockVideoProvider) DeleteStage(ctx context.Context, stageID string) error {
	if _, ok := m.stages[stageID]; !ok {
		return fmt.Errorf("stage not found: %s", stageID)
	}
	delete(m.stages, stageID)
	delete(m.participants, stageID)
	return nil
}

// ListParticipants lists all participants in a stage
func (m *MockVideoProvider) ListParticipants(ctx context.Context, stageID string) ([]*domain.Participant, error) {
	if _, ok := m.stages[stageID]; !ok {
		return nil, fmt.Errorf("stage not found: %s", stageID)
	}
	return m.participants[stageID], nil
}

// DisconnectParticipant disconnects a participant
func (m *MockVideoProvider) DisconnectParticipant(ctx context.Context, stageID, participantID string) error {
	if _, ok := m.stages[stageID]; !ok {
		return fmt.Errorf("stage not found: %s", stageID)
	}
	// In mock, just return success
	return nil
}
