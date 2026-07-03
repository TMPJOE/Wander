package models

import "time"

// Message represents a message between two users.
type Message struct {
	ID         int        `json:"id"`
	SenderID   int        `json:"sender_id"`
	ReceiverID int        `json:"receiver_id"`
	BookingID  *int       `json:"booking_id,omitempty"`
	Content    string     `json:"content"`
	ReadAt     *time.Time `json:"read_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// MessageCreateRequest is used to send a message.
type MessageCreateRequest struct {
	Content   string `json:"content"`
	BookingID *int   `json:"booking_id,omitempty"`
}

// Conversation represents a conversation summary with another user.
type Conversation struct {
	UserID      int       `json:"user_id"`
	UserName    string    `json:"user_name"`
	UserAvatar  string    `json:"user_avatar"`
	LastMessage string    `json:"last_message"`
	LastAt      time.Time `json:"last_at"`
	UnreadCount int       `json:"unread_count"`
}
