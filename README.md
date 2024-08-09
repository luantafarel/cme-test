# Chat System

This project is a simple chat system built with Go, Docker, and various services like Cassandra and Redis. The project is designed to be scalable and easy to deploy using Docker Compose.

## Table of Contents
- [Chat System](#chat-system)
  - [Table of Contents](#table-of-contents)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Project Structure](#project-structure)
    - [Endpoints](#endpoints)
- [Files Overview](#files-overview)
  - [main.go](#maingo)
  - [auth.go](#authgo)
  - [tokenValidator.go](#tokenvalidatorgo)
  - [message.go](#messagego)
  - [contextKeys.go](#contextkeysgo)
  - [grafana.go](#grafanago)
- [Detailed documentation of functions](#detailed-documentation-of-functions)
  - [Database Package](#database-package)
    - [Initialization](#initialization)
    - [Functions](#functions)
      - [`CreateUser(username, password string) error`](#createuserusername-password-string-error)
      - [`ValidateUser(username, password string) (bool, gocql.UUID, error)`](#validateuserusername-password-string-bool-gocqluuid-error)
      - [`CreateSession(userID gocql.UUID) (string, error)`](#createsessionuserid-gocqluuid-string-error)
      - [`SaveMessage(sender, recipient gocql.UUID, content string) error`](#savemessagesender-recipient-gocqluuid-content-string-error)
      - [`GetMessageHistory() ([]map[string]interface{}, error)`](#getmessagehistory-mapstringinterface-error)
      - [`ValidateToken(providedToken string) (gocql.UUID, bool, error)`](#validatetokenprovidedtoken-string-gocqluuid-bool-error)
      - [`GetUserIDByUsername(username string) (gocql.UUID, error)`](#getuseridbyusernameusername-string-gocqluuid-error)
      - [`GetMessagesBetweenUsers(user1, user2 gocql.UUID) ([]types.Message, error)`](#getmessagesbetweenusersuser1-user2-gocqluuid-typesmessage-error)
      - [`GetUsernameByID(userID gocql.UUID) (string, error)`](#getusernamebyiduserid-gocqluuid-string-error)
      - [`GetMessagesForUser(userID gocql.UUID) ([]types.Message, error)`](#getmessagesforuseruserid-gocqluuid-typesmessage-error)
  - [Handlers Package](#handlers-package)
    - [Auth](#auth)
      - [Functions](#functions-1)
        - [`RegisterHandler(w http.ResponseWriter, r *http.Request)`](#registerhandlerw-httpresponsewriter-r-httprequest)
        - [`LoginHandler(w http.ResponseWriter, r *http.Request)`](#loginhandlerw-httpresponsewriter-r-httprequest)
    - [Message](#message)
      - [Types](#types)
        - [Variables](#variables)
          - [`requestsTotal`](#requeststotal)
      - [Functions](#functions-2)
        - [`SendMessageHandler(w http.ResponseWriter, r *http.Request)`](#sendmessagehandlerw-httpresponsewriter-r-httprequest)
        - [`GetMessageHistoryHandler(w http.ResponseWriter, r *http.Request)`](#getmessagehistoryhandlerw-httpresponsewriter-r-httprequest)
  - [Middlewares Package](#middlewares-package)
    - [Functions](#functions-3)
      - [`TokenAuthMiddleware(next http.Handler) http.Handler`](#tokenauthmiddlewarenext-httphandler-httphandler)
  - [Types Package](#types-package)
    - [Variables](#variables-1)
      - [`UserIDKey`](#useridkey)
      - [`Message`](#message-1)
      - [`MessageResponse`](#messageresponse)
      - [`SendMessageBody`](#sendmessagebody)
  - [Main Package](#main-package)
    - [Functions](#functions-4)
      - [`main()`](#main)
- [Notes about the project](#notes-about-the-project)
- [Simple Payment Processor](#simple-payment-processor)
  - [Table of Contents](#table-of-contents-1)
  - [Overview](#overview)
  - [Program Structure](#program-structure)
    - [Data Types](#data-types)
      - [Account Structure](#account-structure)
  - [Functions](#functions-5)
    - [generateUniqueId](#generateuniqueid)
  - [How to Compile and Run](#how-to-compile-and-run)

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Docker
- Docker Compose

### Project Structure

Usage
To start the application, simply run:

```sh
docker compose up
```
This command will build the Go application and start all the services defined in the docker-compose.yml file.

I could not find a proper way to execute the cassandra init file, because of that i have been using like this.
```sh
docker exec -it <containerId> cqlsh
```
Like:
```sh
docker exec -it ed5126712630c739ff76a2e93db6dae61051fa7bb5866a0c432a214626342e3c cqlsh
```
Then copy pasting the content of ./init-casandra.cql inside it 

### Endpoints
- GET /messages: Retrieve all messages.
- POST /messages: Send a new message.
- GET /auth: Authentication endpoint.
- GET /validate-token: Validate user token.

```
Access localhost:8080/docs to review swagger docs
````

# Files Overview

## main.go
This is the main entry point of the application. It sets up the necessary routes and starts the HTTP server.

## auth.go
Handles user authentication.

## tokenValidator.go
Validates tokens for ensuring secure communication.

## message.go
Manages message-related functionalities including storing and retrieving messages from the database.

## contextKeys.go
Defines context keys for storing and retrieving values in request contexts.

## grafana.go
Sets up monitoring and logging with Grafana.

# Detailed documentation of functions
## Database Package

This package provides functions to interact with a Cassandra database for a chat system. The main functionalities include user management, session management, and message handling.

### Initialization

The `init()` function initializes a connection to the Cassandra database using the `CASSANDRA_HOST` environment variable, which should be in the format `host:port`. If this variable is not set, it defaults to `localhost:9042`. The function establishes a session with the Cassandra cluster, setting the keyspace to `chat` and consistency to `Quorum`.

### Functions

#### `CreateUser(username, password string) error`

- **Description**: Inserts a new user into the `users` table with a unique ID, username, password, and the creation timestamp.
- **Parameters**:
  - `username`: The username of the new user.
  - `password`: The password for the new user.
  - id and timestamp are automatically generated inside the function.
- **Returns**: An error if the user creation fails.

#### `ValidateUser(username, password string) (bool, gocql.UUID, error)`

- **Description**: Validates a user's credentials by checking the provided password against the stored password.
- **Parameters**:
  - `username`: The username of the user.
  - `password`: The password provided for validation.
- **Returns**: 
  - `bool`: A boolean indicating whether the validation was successful.
  - `gocql.UUID`: The user's ID on Cassandra go package.
  - `error`: An error if the operation fails.

#### `CreateSession(userID gocql.UUID) (string, error)`

- **Description**: Creates a new session for a user, generating a unique token and storing it in the `sessions` table.
- **Parameters**:
  - `userID`: The ID of the user for whom the session is being created.
- **Returns**:
  - `string`: The session token.
  - `error`: An error if the session creation fails.

#### `SaveMessage(sender, recipient gocql.UUID, content string) error`

- **Description**: Saves a chat message between two users into the `messages` table.
- **Parameters**:
  - `sender`: The ID of the sender.
  - `recipient`: The ID of the recipient.
  - `content`: The content of the message.
- **Returns**: An error if the message saving fails.

#### `GetMessageHistory() ([]map[string]interface{}, error)`

- **Description**: Retrieves the entire message history from the `messages` table.
- **Returns**: 
  - `[]map[string]interface{}`: A slice of maps, each representing a message with its sender, recipient, content, and timestamp.
  - `error`: An error if the retrieval fails.

#### `ValidateToken(providedToken string) (gocql.UUID, bool, error)`

- **Description**: Validates a session token, checking if it exists in the `sessions` table, saving the user id in the context of the execution to be further gathered on the application.
- **Parameters**:
  - `providedToken`: The token provided for validation.
- **Returns**:
  - `gocql.UUID`: The user ID associated with the token.
  - `bool`: A boolean indicating whether the token is valid.
  - `error`: An error if the operation fails.

#### `GetUserIDByUsername(username string) (gocql.UUID, error)`

- **Description**: Retrieves a user's ID based on their username.
- **Parameters**:
  - `username`: The username of the user.
- **Returns**: 
  - `gocql.UUID`: The user's ID.
  - `error`: An error if the retrieval fails.

#### `GetMessagesBetweenUsers(user1, user2 gocql.UUID) ([]types.Message, error)`

- **Description**: Retrieves all messages exchanged between two users, ordered by timestamp.
- **Parameters**:
  - `user1`: The ID of the first user.
  - `user2`: The ID of the second user.
- **Returns**:
  - `[]types.Message`: A slice of `types.Message` structs representing the messages.
  - `error`: An error if the retrieval fails.

#### `GetUsernameByID(userID gocql.UUID) (string, error)`

- **Description**: Retrieves a username based on the user's ID.
- **Parameters**:
  - `userID`: The ID of the user.
- **Returns**: 
  - `string`: The username.
  - `error`: An error if the retrieval fails.

#### `GetMessagesForUser(userID gocql.UUID) ([]types.Message, error)`

- **Description**: Retrieves all messages for a user, considering both sent and received messages. Duplicates are avoided by storing messages in a map.
- **Parameters**:
  - `userID`: The ID of the user for whom the messages are being retrieved.
- **Returns**:
  - `[]types.Message`: A slice of `types.Message` structs representing the messages.
  - `error`: An error if the retrieval fails.

## Handlers Package

### Auth 

This package provides HTTP handlers for user registration and login in a chat system. These handlers interact with the database to create users, validate credentials, and manage sessions.

#### Functions

##### `RegisterHandler(w http.ResponseWriter, r *http.Request)`

- **Description**: Handles user registration by decoding a `User` object from the request body and creating a new user in the database.
- **Parameters**:
  - `w`: The HTTP response writer.
  - `r`: The HTTP request.
- **Behavior**:
  - Decodes the `User` object from the request body.
  - Creates a new user in the database.
  - Returns a `201 Created` status if successful, or a `500 Internal Server Error` if an error occurs.

##### `LoginHandler(w http.ResponseWriter, r *http.Request)`

- **Description**: Handles user login by validating credentials and creating a session if the credentials are correct.
- **Parameters**:
  - `w`: The HTTP response writer.
  - `r`: The HTTP request.
- **Behavior**:
  - Decodes the `User` object from the request body.
  - Validates the user's credentials.
  - If the credentials are valid, creates a session and returns a `200 OK` status with a session token.
  - Returns a `400 Bad Request` status if the request payload is invalid, or a `401 Unauthorized` status if the credentials are invalid.

### Message

This package provides HTTP handlers for sending messages and retrieving message history in a chat system. It also includes Prometheus metrics for monitoring HTTP requests.

#### Types

##### Variables

###### `requestsTotal`

- **Description**: A Prometheus counter vector that tracks the total number of HTTP requests made to the server, labeled by the request path.
- **Type**: `prometheus.CounterVec`

#### Functions

##### `SendMessageHandler(w http.ResponseWriter, r *http.Request)`

- **Description**: Handles sending a message from the logged-in user to another user.
- **Parameters**:
  - `w`: The HTTP response writer.
  - `r`: The HTTP request.
- **Behavior**:
  - Decodes the `Message` object from the request body.
  - Retrieves the sender's `userID` from the request context.
  - Fetches the recipient's `userID` based on the provided username.
  - Saves the message to the database.
  - Returns a `201 Created` status if successful, or appropriate error messages and status codes if an error occurs.

##### `GetMessageHistoryHandler(w http.ResponseWriter, r *http.Request)`

- **Description**: Handles retrieving the message history for the logged-in user.
- **Parameters**:
  - `w`: The HTTP response writer.
  - `r`: The HTTP request.
- **Behavior**:
  - Increments the Prometheus counter for the request path.
  - Retrieves the user's `userID` from the request context.
  - Fetches all messages for the user from the database.
  - Maps user IDs to usernames for formatting purposes.
  - Groups and formats the messages by the other users involved in the conversations.
  - Sends the grouped messages as a JSON response.
  - Returns a `200 OK` status if successful, or appropriate error messages and status codes if an error occurs.

## Middlewares Package

This package provides middleware functions for the chat system, specifically for handling token-based authentication.

### Functions

#### `TokenAuthMiddleware(next http.Handler) http.Handler`

- **Description**: Middleware that handles token-based authentication for incoming HTTP requests. It validates the token from the `Authorization` header and sets the `userID` in the request context if the token is valid.
- **Parameters**:
  - `next`: The next HTTP handler to be called if the token is valid.
- **Behavior**:
  - Extracts the token from the `Authorization` header of the request.
  - Logs the received token.
  - Validates the token using the `ValidateToken` function from the `database` package.
  - If the token is missing or invalid, responds with a `401 Unauthorized` status.
  - If the token is valid, logs the `userID` and sets it in the request context.
  - Passes the request to the next handler in the chain with the modified context.

## Types Package

This file defines keys used for storing and retrieving values from the context within the application and the structures related to messages within the chat system.

### Variables

#### `UserIDKey`

- **Description**: A key used to store and retrieve the user ID from the request context.
- **Type**: `contextKey` (a custom type to avoid key collisions in the context)

#### `Message`

- **Description**: Represents a chat message in the system.
- **Fields**:
  - `ID`: The unique identifier of the message.
  - `Sender`: The ID of the user who sent the message.
  - `Recipient`: The ID of the user who received the message.
  - `Content`: The content of the message.
  - `Timestamp`: The time when the message was sent.

#### `MessageResponse`

- **Description**: Represents the response format of a chat message, including usernames instead of IDs.
- **Fields**:
  - `SenderUsername`: The username of the user who sent the message.
  - `RecipientUsername`: The username of the user who received the message.
  - `Content`: The content of the message.
  - `Timestamp`: The time when the message was sent.


#### `SendMessageBody`

- **Description**: Represents the response format of a chat message, including usernames instead of IDs.
- **Fields**:
  - `Recipient`: Name of the user that will receive the message
  - `Content`: The content of the message.

## Main Package

This is the entry point of the chat system application. It sets up the HTTP server and routes.

### Functions

#### `main()`

- **Description**: Initializes the HTTP server, sets up the routes, and starts the server.
- **Behavior**:
  - Configures the HTTP server with routes for registration, login, sending messages, and retrieving message history.
  - Applies necessary middleware, such as token-based authentication, to protect routes that require authorization.
  - Listens and serves on the specified port, handling incoming requests.




# Notes about the project
The basic project was done using in memory storage, then i upgraded it to use cassandra, I haven'nt had any access to cassandra before. Being honest it looks a lot with the querying of other SQL languages, but the things that are different are really stressing, like the lack of OR operator and the need to create index to query or having to use ALLOW FILTERING on every request.

The main issue with the docker was to find a way to proper initialize a keyspace, I gave up and did it manually.

I did the code is as simple as i could, you register, generate a session when logging in, then onwards i use the token and body to send messages setting the sender and receiver id.

The project has limitations
- you can send message to yourself (as in whatsapp and telegram), but since i have discovered that theres no OR on the Cassandra i had to do two different queries, to request messages where the user is the sender and where he is the receiver. That generated a bug where if the sender and receiver where the same it would generate duplicated messages.

- known bugs
  - you can have many sessions and valid logins and they will live forever (i feel sorry for that)

  - you can have multiple users with the same username

I have only added test for auth, in fact i had other tests but since i have changed and added the new stuff like cassandra and auth token they were broken and i haven't fixed them
Run 
```sh
go test ./...
```

The code adds a start for prometheus and grafana, but not much. You can access it through localhost:3000 user and password is admin. Its too primitive does'nt work properly.
Redis is implemented but not used.

#Payment Processor 
# Simple Payment Processor

This is a simple C++ program that simulates a basic payment processor. The program allows you to create accounts, process transactions between accounts, and track account balances. This implementation is designed to be easy to understand, making use of basic C++ constructs like structures, arrays, and simple functions.

## Table of Contents

- [Overview](#overview)
- [Program Structure](#program-structure)
  - [Data Types](#data-types)
    - [Account Structure](#account-structure)
    - [Transaction Structure](#transaction-structure)
  - [Functions](#functions)
    - [generateUniqueId](#generateuniqueid)
    - [findAccountIndex](#findaccountindex)
    - [createAccount](#createaccount)
    - [processTransaction](#processtransaction)
- [How to Compile and Run](#how-to-compile-and-run)
- [Example Output](#example-output)

## Overview

This program demonstrates the basic functionality of a payment processor:
- **Account Management**: Create new accounts with unique IDs, store the account owner's name, and track the account balance.
- **Transaction Processing**: Transfer funds from one account to another, update balances accordingly, and record the transaction.

## Program Structure

### Data Types

#### Account Structure

```cpp
struct Account {
    std::string accountId;
    std::string ownerName;
    double balance;
};
```

- accountId: A unique identifier for the account, represented as a string.
- ownerName: The name of the account owner, represented as a string.
- balance: The current balance of the account, represented as a double.

```cpp
struct Transaction {
    std::string transactionId;
    std::string fromAccountId;
    std::string toAccountId;
    double amount;
};
```

- transactionId: A unique identifier for the transaction, represented as a string.
- fromAccountId: The ID of the account from which funds are being transferred, represented as a string.
- toAccountId: The ID of the account to which funds are being transferred, represented as a string.
- amount: The amount of money being transferred, represented as a double.


## Functions

### generateUniqueId

```cpp
std::string generateUniqueId();
```
- Purpose: Generates a simple, random unique identifier that can be used as an account ID or transaction ID.
- Returns: A string representing the unique ID.
- Usage: This function is called whenever a new account or transaction is created to assign it a unique identifier.

```cpp
int findAccountIndex(const std::string& accountId);
```

- Purpose: Searches the accounts array to find the index of an account by its ID.
- Parameters:
  - accountId: The ID of the account to find.
- Returns: The index of the account in the accounts array, or -1 if the account is not found.
- Usage: This function is used internally to locate accounts by their ID before performing operations like transactions.

```cpp
std::string createAccount(const std::string& ownerName, double initialBalance);
```
- Purpose: Creates a new account, assigns it a unique ID, and initializes its balance.
- Parameters:
  - ownerName: The name of the account owner.
  - initialBalance: The starting balance of the account.
- Returns: The unique ID of the newly created account.
- Usage: Call this function to create a new account and add it to the system. It prints the details of the newly created account.

```cpp
bool processTransaction(const std::string& fromAccountId, const std::string& toAccountId, double amount);
```
- Purpose: Transfers money from one account to another and records the transaction.
- Parameters:
  - fromAccountId: The ID of the account from which the funds are to be transferred.
  - toAccountId: The ID of the account to which the funds are to be transferred.
  - amount: The amount of money to transfer.
- Returns: A boolean value indicating whether the transaction was successful (true) or not (false).
- Usage: This function is used to perform a transaction between two accounts. It checks if both accounts exist and if the source account has sufficient funds, then processes the transaction and prints the details.


## How to Compile and Run
Compile the Code:

Open a terminal and navigate to the directory where the file is saved.
Use a C++ compiler like g++ to compile the code:

```sh
g++ -o payment_processor payment_processor.cpp
```

Run 
```sh
./payment_processor
```