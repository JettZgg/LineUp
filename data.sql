-- Drop tables if they exist
DROP TABLE IF EXISTS matches, users;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    uid BIGINT PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create matches table
CREATE TABLE IF NOT EXISTS matches (
    id BIGINT PRIMARY KEY,
    player1_id BIGINT REFERENCES users(uid),
    player2_id BIGINT REFERENCES users(uid),
    winner BIGINT,
    first_move_player_id BIGINT,
    moves TEXT,
    date TIMESTAMP WITH TIME ZONE
);