package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"wander/backend/internal/models"
)

type MessageRepository interface {
	Create(ctx context.Context, senderID int, receiverID int, content string, bookingID *int) (*models.Message, error)
	ListConversations(ctx context.Context, userID int) ([]models.Conversation, error)
	ListMessages(ctx context.Context, userID int, otherID int) ([]models.Message, error)
}

type PgMessageRepository struct {
	pool *pgxpool.Pool
}

func NewPgMessageRepository(pool *pgxpool.Pool) MessageRepository {
	return &PgMessageRepository{pool: pool}
}

func (r *PgMessageRepository) Create(ctx context.Context, senderID int, receiverID int, content string, bookingID *int) (*models.Message, error) {
	query := `
		INSERT INTO messages (sender_id, receiver_id, content, booking_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, sender_id, receiver_id, booking_id, content, read_at, created_at
	`
	m := &models.Message{}
	err := r.pool.QueryRow(ctx, query, senderID, receiverID, content, bookingID).Scan(
		&m.ID, &m.SenderID, &m.ReceiverID, &m.BookingID, &m.Content, &m.ReadAt, &m.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}
	return m, nil
}

func (r *PgMessageRepository) ListConversations(ctx context.Context, userID int) ([]models.Conversation, error) {
	query := `
		WITH last_messages AS (
			SELECT DISTINCT ON (LEAST(sender_id, receiver_id), GREATEST(sender_id, receiver_id))
				id, sender_id, receiver_id, content, created_at
			FROM messages
			WHERE sender_id = $1 OR receiver_id = $1
			ORDER BY LEAST(sender_id, receiver_id), GREATEST(sender_id, receiver_id), created_at DESC
		)
		SELECT 
			u.id as user_id,
			u.first_name || ' ' || u.last_name as user_name,
			u.avatar_url as user_avatar,
			lm.content as last_message,
			lm.created_at as last_at,
			(SELECT COUNT(*)::INT FROM messages m2 WHERE m2.sender_id = u.id AND m2.receiver_id = $1 AND m2.read_at IS NULL) as unread_count
		FROM last_messages lm
		JOIN users u ON u.id = CASE WHEN lm.sender_id = $1 THEN lm.receiver_id ELSE lm.sender_id END
		ORDER BY lm.created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list conversations: %w", err)
	}
	defer rows.Close()

	var conversations []models.Conversation
	for rows.Next() {
		var c models.Conversation
		err := rows.Scan(&c.UserID, &c.UserName, &c.UserAvatar, &c.LastMessage, &c.LastAt, &c.UnreadCount)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, c)
	}
	return conversations, nil
}

func (r *PgMessageRepository) ListMessages(ctx context.Context, userID int, otherID int) ([]models.Message, error) {
	// Mark incoming messages as read first.
	_, _ = r.pool.Exec(ctx, "UPDATE messages SET read_at = NOW() WHERE sender_id = $1 AND receiver_id = $2 AND read_at IS NULL", otherID, userID)

	query := `
		SELECT id, sender_id, receiver_id, booking_id, content, read_at, created_at
		FROM messages
		WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)
		ORDER BY created_at ASC
	`
	rows, err := r.pool.Query(ctx, query, userID, otherID)
	if err != nil {
		return nil, fmt.Errorf("list messages: %w", err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var m models.Message
		err := rows.Scan(&m.ID, &m.SenderID, &m.ReceiverID, &m.BookingID, &m.Content, &m.ReadAt, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}
