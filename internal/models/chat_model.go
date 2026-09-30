package models

import "time"

type ChatRole string

const (
	ChatRoleUser      ChatRole = "user"
	ChatRoleAssistant ChatRole = "assistant"
	ChatRoleSystem    ChatRole = "system"
)

type ChatConversation struct {
	ID          string
	UserID      string
	ItineraryID *string
	Topic       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ChatMessage struct {
	ID             string
	ConversationID string
	Role           ChatRole
	Content        string
	CreatedAt      time.Time
}
