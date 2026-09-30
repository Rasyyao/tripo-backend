package models

import "time"

type ChatRole string

const (
	ChatRoleUser      ChatRole = "user"
	ChatRoleAssistant ChatRole = "assistant"
	ChatRoleSystem    ChatRole = "system"
)

type ChatConversation struct {
	ID          string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      string     `gorm:"type:uuid;not null;index:idx_chat_conversations_user,priority:1"`
	User        User       `gorm:"constraint:OnDelete:CASCADE"`
	ItineraryID *string    `gorm:"type:uuid"`
	Itinerary   *Itinerary `gorm:"constraint:OnDelete:SET NULL"`
	Topic       string     `gorm:"size:150;not null"`
	CreatedAt   time.Time  `gorm:"not null"`
	UpdatedAt   time.Time  `gorm:"not null;index:idx_chat_conversations_user,priority:2,sort:desc"`
}

type ChatMessage struct {
	ID             string           `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ConversationID string           `gorm:"type:uuid;not null;index:idx_chat_messages_conversation,priority:1"`
	Conversation   ChatConversation `gorm:"constraint:OnDelete:CASCADE"`
	Role           ChatRole         `gorm:"size:20;not null"`
	Content        string           `gorm:"type:text;not null"`
	CreatedAt      time.Time        `gorm:"not null;index:idx_chat_messages_conversation,priority:2"`
}
