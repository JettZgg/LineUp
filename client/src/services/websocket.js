// src/services/websocket.js
import { useState, useEffect, useCallback, useRef } from 'react';

const WS_URL = 'ws://localhost:8080/ws'; // Update with your WebSocket server URL
const MAX_RETRIES = 5;
const INITIAL_RETRY_DELAY = 1000; // 1 second
const INITIAL_CONNECTION_DELAY = 500; // 500ms delay before first connection attempt

export const useWebSocket = (matchId, user) => {
    const [socket, setSocket] = useState(null);
    const [lastMessage, setLastMessage] = useState(null);
    const [isConnecting, setIsConnecting] = useState(true);
    const [isConnected, setIsConnected] = useState(false);
    const retryCount = useRef(0);
    const retryDelay = useRef(INITIAL_RETRY_DELAY);

    const connectWebSocket = useCallback(() => {
        if (!matchId || matchId === 'undefined' || !user || !user.token) {
            console.error('Invalid matchId or user');
            setIsConnecting(false);
            return;
        }

        const ws = new WebSocket(`${WS_URL}/${matchId}`);

        ws.onopen = () => {
            console.log('WebSocket connected');
            setIsConnected(true);
            setIsConnecting(false);
            retryCount.current = 0;
            retryDelay.current = INITIAL_RETRY_DELAY;
            ws.send(JSON.stringify({ type: 'joinMatch', matchId, token: user.token }));
        };

        ws.onmessage = (event) => {
            setLastMessage(event);
        };

        ws.onclose = (event) => {
            console.log('WebSocket disconnected', event);
            setIsConnected(false);
            if (retryCount.current < MAX_RETRIES) {
                const delay = retryDelay.current;
                console.log(`Retrying connection in ${delay}ms... Attempt ${retryCount.current + 1}`);
                setTimeout(() => {
                    retryCount.current += 1;
                    retryDelay.current *= 2; // Exponential backoff
                    setIsConnecting(true);
                    connectWebSocket();
                }, delay);
            } else {
                console.error('Max retries reached. WebSocket connection failed.');
                setIsConnecting(false);
            }
        };

        setSocket(ws);

        return () => {
            ws.close();
        };
    }, [matchId, user]);

    useEffect(() => {
        const timer = setTimeout(() => {
            connectWebSocket();
        }, INITIAL_CONNECTION_DELAY);

        return () => {
            clearTimeout(timer);
        };
    }, [connectWebSocket]);

    const sendMessage = useCallback((message) => {
        if (socket && socket.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify(message));
        }
    }, [socket]);

    return { sendMessage, lastMessage, isConnected, isConnecting };
};