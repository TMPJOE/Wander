package service

import (
	"context"
	"sync"

	"wander/backend/internal/models"
	"wander/backend/internal/repository"
)

// ActiveConnection represents a client's message stream.
type ActiveConnection struct {
	UserID int
	Send   chan models.Message
}

type MessageService struct {
	repo repository.MessageRepository
	mu   sync.RWMutex
	// activeStreams keeps track of users connected to live chat.
	activeStreams map[int][]chan models.Message
}

func NewMessageService(repo repository.MessageRepository) *MessageService {
	return &MessageService{
		repo:          repo,
		activeStreams: make(map[int][]chan models.Message),
	}
}

func (s *MessageService) Create(ctx context.Context, senderID int, receiverID int, content string, bookingID *int) (*models.Message, error) {
	m, err := s.repo.Create(ctx, senderID, receiverID, content, bookingID)
	if err != nil {
		return nil, err
	}

	// Broadcast message to any active channel streams.
	s.mu.RLock()
	streams, ok := s.activeStreams[receiverID]
	s.mu.RUnlock()

	if ok {
		for _, ch := range streams {
			select {
			case ch <- *m:
			default:
			}
		}
	}

	return m, nil
}

func (s *MessageService) ListConversations(ctx context.Context, userID int) ([]models.Conversation, error) {
	return s.repo.ListConversations(ctx, userID)
}

func (s *MessageService) ListMessages(ctx context.Context, userID int, otherID int) ([]models.Message, error) {
	return s.repo.ListMessages(ctx, userID, otherID)
}

// RegisterStream adds a message channel for live updates.
func (s *MessageService) RegisterStream(userID int, ch chan models.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.activeStreams[userID] = append(s.activeStreams[userID], ch)
}

// UnregisterStream removes a message channel.
func (s *MessageService) UnregisterStream(userID int, ch chan models.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	streams := s.activeStreams[userID]
	for i, stream := range streams {
		if stream == ch {
			s.activeStreams[userID] = append(streams[:i], streams[i+1:]...)
			break
		}
	}
	if len(s.activeStreams[userID]) == 0 {
		delete(s.activeStreams, userID)
	}
}
