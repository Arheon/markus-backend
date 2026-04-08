package outbox

import (
	"time"

	"github.com/google/uuid"
)

type MessageOutput func(*Message)

type Message struct {
	ID       uuid.UUID
	Type     string
	Payload  []byte
	Error    *string
	Metadata map[string]string `gorm:"serializer:json"`

	CreatedAt   time.Time
	ProcessedAt *time.Time

	Attempts  int
	LastError *string
}

func (m *Message) TableName() string {
	return "outbox"
}

func WithID(id uuid.UUID) MessageOutput {
	return func(m *Message) {
		m.ID = id
	}
}

func WithType(t string) MessageOutput {
	return func(m *Message) {
		m.Type = t
	}
}

func WithMetadata(metadata map[string]string) MessageOutput {
	return func(m *Message) {
		m.Metadata = metadata
	}
}

func WithCreatedAt(createdAt time.Time) MessageOutput {
	return func(m *Message) {
		m.CreatedAt = createdAt
	}
}

func WithProcessedAt(processedAt *time.Time) MessageOutput {
	return func(m *Message) {
		m.ProcessedAt = processedAt
	}
}

func WithAttempts(attempts int) MessageOutput {
	return func(m *Message) {
		m.Attempts = attempts
	}
}

func NewMessage(payload []byte, opts ...MessageOutput) *Message {
	now := time.Now().UTC()
	m := &Message{
		ID:          uuid.New(),
		CreatedAt:   now,
		ProcessedAt: &now,
		Type:        "",
		Payload:     payload,
		Attempts:    0,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}
