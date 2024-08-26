# Makefile

# Variables
API_SERVER_CMD = cmd/apiserver/main.go
WS_SERVER_CMD = cmd/wsserver/main.go

# Default target
all: run

# Run both servers
run:
	@echo "Starting API server..."
	@go run $(API_SERVER_CMD) &
	@echo "Starting WebSocket server..."
	@go run $(WS_SERVER_CMD)

.PHONY: all run