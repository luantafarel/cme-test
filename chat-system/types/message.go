package types

import (
	"time"

	"github.com/gocql/gocql"
)

type MessageResponse struct {
	SenderUsername    string    `json:"sender_username"`
	RecipientUsername string    `json:"recipient_username"`
	Content           string    `json:"content"`
	Timestamp         time.Time `json:"timestamp"`
}

type Message struct {
	ID        gocql.UUID `json:"id"`
	Recipient gocql.UUID `json:"recipient"`
	Sender    gocql.UUID `json:"sender"`
	Content   string     `json:"content"`
	Timestamp time.Time  `json:"timestamp"`
}

type SendMessageBody struct {
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
}
