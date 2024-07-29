// src/services/websocket.js
import { useEffect, useRef } from 'react';

const WS_URL = 'ws://localhost:8080/ws'; // Update with your WebSocket server URL

export const useWebSocket = (matchId) => {
    const socket = useRef(null);

    useEffect(() => {
        // Create WebSocket connection
        socket.current = new WebSocket(`${WS_URL}/${matchId}`);

        socket.current.onopen = () => {
            console.log('WebSocket Connected');
        };

        socket.current.onmessage = (event) => {
            const data = JSON.parse(event.data);
            // Handle incoming messages
            console.log('Received:', data);
            // You'll want to update your game state here
        };

        socket.current.onclose = () => {
            console.log('WebSocket Disconnected');
        };

        return () => {
            socket.current.close();
        };
    }, [matchId]);

    const sendMessage = (message) => {
        if (socket.current.readyState === WebSocket.OPEN) {
            socket.current.send(JSON.stringify(message));
        }
    };

    return { sendMessage };
};