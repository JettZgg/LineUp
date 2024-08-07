next steps:
-  Fix "Join match" and "get match"
-  Add comprehensive input validation for all API endpoints
-  Implement detailed logging throughout the application
-  Set up unit tests for core functionalities (game logic, auth, API handlers)
-  Create a basic client-side interface for testing the game
-  Set up a CI/CD pipeline for automated testing and deployment
-  Optimize database queries and add indexing if needed
-  Implement rate limiting and additional security measures
-  Add documentation for API endpoints and WebSocket events
-  Consider adding monitoring and analytics for the application

CREATE TABLE users (
    uid BIGINT PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE matches (
    id BIGINT PRIMARY KEY,
    player1_id BIGINT REFERENCES users(uid),
    player2_id BIGINT REFERENCES users(uid),
    status VARCHAR(50),
    start_time TIMESTAMP WITH TIME ZONE,
    end_time TIMESTAMP WITH TIME ZONE,
    winner BIGINT,
    board_width INTEGER,
    board_height INTEGER,
    win_length INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE moves (
    id SERIAL PRIMARY KEY,
    match_id BIGINT NOT NULL,
    player_id BIGINT NOT NULL,
    x INTEGER NOT NULL,
    y INTEGER NOT NULL,
    move_number INTEGER NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (match_id) REFERENCES matches(id) ON DELETE CASCADE,
    FOREIGN KEY (player_id) REFERENCES users(uid)
);

DROP TABLE moves;
DROP TABLE matches;
DROP TABLE users;