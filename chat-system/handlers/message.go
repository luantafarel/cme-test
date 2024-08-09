package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/prometheus/client_golang/prometheus"

	"chat-system/database"
	"chat-system/types"
)

var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"path"},
	)
)

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	var msg types.SendMessageBody
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userIDValue := r.Context().Value(types.UserIDKey)
	if userIDValue == nil {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDValue.(gocql.UUID)
	if !ok {
		http.Error(w, "User ID is not of type gocql.UUID", http.StatusInternalServerError)
		return
	}

	recipientUuid, err := database.GetUserIDByUsername(msg.Recipient)
	if err != nil {
		fmt.Println("Invalid UUID string:", err)
	}

	if err := database.SaveMessage(userID, recipientUuid, msg.Content); err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetMessageHistoryHandler(w http.ResponseWriter, r *http.Request) {
	requestsTotal.WithLabelValues(r.URL.Path).Inc()
	// Extract userID from context
	userIDValue := r.Context().Value(types.UserIDKey)
	if userIDValue == nil {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDValue.(gocql.UUID)
	if !ok {
		http.Error(w, "User ID is not of type gocql.UUID", http.StatusInternalServerError)
		return
	}

	// Fetch messages for the user
	messages, err := database.GetMessagesForUser(userID)
	if err != nil {
		http.Error(w, "Error fetching messages", http.StatusInternalServerError)
		return
	}

	// Create a map to store userID -> username mappings
	userMap := make(map[gocql.UUID]string)
	userMap[userID], _ = database.GetUsernameByID(userID)

	// Fetch usernames for all involved users
	for _, msg := range messages {
		if _, exists := userMap[msg.Sender]; !exists {
			userMap[msg.Sender], _ = database.GetUsernameByID(msg.Sender)
		}
		if _, exists := userMap[msg.Recipient]; !exists {
			userMap[msg.Recipient], _ = database.GetUsernameByID(msg.Recipient)
		}
	}

	// Group and format messages for the response
	groupedMessages := make(map[string][]types.MessageResponse)
	for _, msg := range messages {
		var otherUserID gocql.UUID
		if msg.Sender == userID && msg.Recipient == userID {
			// Both sender and recipient are the same user (current user)
			otherUserID = userID
		} else if msg.Recipient == userID {
			// The current user is the recipient, so the other user is the sender
			otherUserID = msg.Sender
		} else {
			// The current user is the sender, so the other user is the recipient
			otherUserID = msg.Recipient
		}

		otherUsername := userMap[otherUserID]
		messageResponse := types.MessageResponse{
			SenderUsername:    userMap[msg.Sender],
			RecipientUsername: userMap[msg.Recipient],
			Content:           msg.Content,
			Timestamp:         msg.Timestamp,
		}

		// Prevent adding the message twice when sender and recipient are the same
		if msg.Sender != msg.Recipient {
			groupedMessages[otherUsername] = append(groupedMessages[otherUsername], messageResponse)
		} else if msg.Sender == userID {
			// Special handling when sender and recipient are the same user
			selfUsername := userMap[userID]
			groupedMessages[selfUsername] = append(groupedMessages[selfUsername], messageResponse)
		}
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groupedMessages)
}
