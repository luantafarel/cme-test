package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocql/gocql"

	"chat-system/types"
)

var session *gocql.Session

func init() {
	cassandraHost := os.Getenv("CASSANDRA_HOST")
	if cassandraHost == "" {
		cassandraHost = "localhost:9042"
	}

	// Split the host and port
	hostPort := strings.Split(cassandraHost, ":")
	if len(hostPort) != 2 {
		log.Fatal("CASSANDRA_HOST environment variable is not in the correct format (host:port)")
	}
	host := hostPort[0]
	port := hostPort[1]
	portInt, errStrconv := strconv.Atoi(port)
	if errStrconv != nil {
		log.Fatalf("Invalid port: %v", errStrconv)
	}
	cluster := gocql.NewCluster(host)
	cluster.Port = portInt
	cluster.Keyspace = "chat"
	cluster.Consistency = gocql.Quorum

	var err error
	session, err = cluster.CreateSession()

	session.Query(``)
	if err != nil {
		log.Fatalf("Unable to connect to Cassandra: %v", err)
	}
}

func CreateUser(username, password string) error {
	return session.Query(`INSERT INTO users (id, username, password, created_at) VALUES (?, ?, ?, ?)`, gocql.TimeUUID(), username, password, time.Now()).Exec()
}

func ValidateUser(username, password string) (bool, gocql.UUID, error) {
	var storedPassword string
	var userId gocql.UUID
	err := session.Query(`SELECT password, id FROM users WHERE username = ?`, username).Scan(&storedPassword, &userId)
	if err != nil {
		if err == gocql.ErrNotFound {
			return false, userId, gocql.Error{Code: 404, Message: "Not found"}
		}
		return false, userId, err
	}
	return password == storedPassword, userId, nil // Password is not hashed for brevity
}

func CreateSession(userID gocql.UUID) (string, error) {
	uuid := gocql.TimeUUID().String()
	err := session.Query(`INSERT INTO sessions (id, created_at, "token", user_id) VALUES (?, ?, ?, ?)`, gocql.TimeUUID(), time.Now(), uuid, userID).Exec()
	if err != nil {
		return "", err
	}
	return uuid, nil
}

func SaveMessage(sender, recipient gocql.UUID, content string) error {
	fmt.Println(sender, recipient, content)
	return session.Query(`INSERT INTO messages (id, sender, recipient, content, timestamp) VALUES (?, ?, ?, ?, ?)`, gocql.MustRandomUUID(), sender, recipient, content, time.Now()).Exec()
}

func GetMessageHistory() ([]map[string]interface{}, error) {
	var messages []map[string]interface{}
	iter := session.Query(`SELECT sender, recipient, content, timestamp FROM messages ALLOW FILTERING`).Iter()
	for {
		row := make(map[string]interface{})
		if !iter.MapScan(row) {
			break
		}
		messages = append(messages, row)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return messages, nil
}

func ValidateToken(providedToken string) (gocql.UUID, bool, error) {
	var dbToken string
	var userID gocql.UUID
	token := strings.TrimPrefix(providedToken, "Bearer ")
	err := session.Query(`SELECT user_id, "token" FROM sessions WHERE "token" = ? ALLOW FILTERING`, token).Scan(&userID, &dbToken)
	if err != nil {
		if err == gocql.ErrNotFound {
			return gocql.UUID{}, false, gocql.Error{Code: 404, Message: "Token not found"}
		}
		return gocql.UUID{}, false, err
	}
	return userID, true, nil
}

func GetUserIDByUsername(username string) (gocql.UUID, error) {
	var userID gocql.UUID
	err := session.Query(`SELECT id FROM users WHERE username = ?  ALLOW FILTERING`, username).Scan(&userID)
	if err != nil {
		log.Printf("Get by username error: %v", err)
		return gocql.UUID{}, err
	}
	log.Printf("Get by username Token: %s", userID)
	return userID, nil
}

func GetMessagesBetweenUsers(user1, user2 gocql.UUID) ([]types.Message, error) {
	var messages []types.Message

	query := `SELECT id, recipient, sender, content, timestamp FROM messages WHERE (sender = ? AND recipient = ?) OR (sender = ? AND recipient = ?) ORDER BY timestamp`
	iter := session.Query(query, user1, user2, user2, user1).Iter()
	for {
		var msg types.Message
		if !iter.Scan(&msg.ID, &msg.Recipient, &msg.Sender, &msg.Content, &msg.Timestamp) {
			break
		}
		messages = append(messages, msg)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}
	return messages, nil
}

func GetUsernameByID(userID gocql.UUID) (string, error) {
	var username string
	err := session.Query(`SELECT username FROM users WHERE id = ?`, userID).Scan(&username)
	return username, err
}

func GetMessagesForUser(userID gocql.UUID) ([]types.Message, error) {
	// this is to prevent messages with the same recipient and sender from being duplicated
	messagesMap := make(map[gocql.UUID]types.Message)

	iterSender := session.Query(`SELECT id, recipient, sender, content, timestamp FROM messages WHERE sender = ? ALLOW FILTERING`, userID).Iter()
	for {
		var msg types.Message
		if !iterSender.Scan(&msg.ID, &msg.Recipient, &msg.Sender, &msg.Content, &msg.Timestamp) {
			break
		}
		messagesMap[msg.ID] = msg
	}

	if err := iterSender.Close(); err != nil {
		fmt.Println("Error fetching messages:", err)
		return nil, err
	}

	iterRecipient := session.Query(`SELECT id, recipient, sender, content, timestamp FROM messages WHERE recipient = ? ALLOW FILTERING`, userID).Iter()
	for {
		var msg types.Message
		if !iterRecipient.Scan(&msg.ID, &msg.Recipient, &msg.Sender, &msg.Content, &msg.Timestamp) {
			break
		}
		messagesMap[msg.ID] = msg
	}

	if err := iterRecipient.Close(); err != nil {
		fmt.Println("Error fetching messages:", err)
		return nil, err
	}

	var messages []types.Message
	for _, msg := range messagesMap {
		messages = append(messages, msg)
	}

	return messages, nil
}
