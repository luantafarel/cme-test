CREATE KEYSPACE IF NOT EXISTS chat WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'};
USE chat;

CREATE TABLE IF NOT EXISTS auth (
    id UUID PRIMARY KEY,
    username TEXT,
    password TEXT,
    created_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY,
    recipient UUID,
    sender UUID,
    content TEXT,
    "timestamp" TIMESTAMP,
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username TEXT,
    password TEXT,
    created_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY,
    user_id UUID,
    "token" TEXT,
    created_at TIMESTAMP
);

CREATE INDEX ON users (username);
CREATE INDEX ON messages (recipient);