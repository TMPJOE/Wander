package repository

import (
	"context"

	"wander/backend/internal/models"
)

type MessageRepository interface {
	Create(ctx context.Context, senderID int, receiverID int, content string, bookingID *int) (*models.Message, error)
	ListConversations(ctx context.Context, userID int) ([]models.Conversation, error)
	ListMessages(ctx context.Context, userID int, otherID int) ([]models.Message, error)
}
