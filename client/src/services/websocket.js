// src/services/websocket.js
import { useState, useEffect, useCallback } from 'react';

const WS_URL = 'ws://localhost:8080/ws'; // Update with your WebSocket server URL

export const useWebSocket = (matchId) => {
    const [socket, setSocket] = useState(null);
    const [lastMessage, setLastMessage] = useState(null);

    useEffect(() => {
        const ws = new WebSocket(`${WS_URL}/${matchId}`);
        setSocket(ws);

        ws.onmessage = (event) => {
            setLastMessage(event);
        };

        return () => {
            ws.close();
        };
    }, [matchId]);

    const sendMessage = useCallback((message) => {
        if (socket && socket.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify(message));
        }
    }, [socket]);

    return { sendMessage, lastMessage };
};