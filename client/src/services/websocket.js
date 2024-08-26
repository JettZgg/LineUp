import { useState, useEffect, useCallback, useRef } from 'react';
import { WS_BASE_URL } from '../config';

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

        const ws = new WebSocket(`${WS_BASE_URL}/${matchId}`);

        ws.onopen = () => {
            console.log('WebSocket connected');
            setIsConnected(true);
            setIsConnecting(false);
            retryCount.current = 0;
            retryDelay.current = INITIAL_RETRY_DELAY;
            const joinMessage = JSON.stringify({ type: 'joinMatch', matchId, token: user.token });
            console.log('Sending join message:', joinMessage);
            ws.send(joinMessage);
        };

        ws.onmessage = (event) => {
            console.log('Raw WebSocket message:', event.data);
            if (event.data === undefined || event.data === 'undefined') {
                console.error('Received undefined WebSocket message');
                return;
            }
            try {
                const data = JSON.parse(event.data);
                console.log('Parsed WebSocket message:', data);
                setLastMessage(data);
            } catch (error) {
                console.error('Error parsing WebSocket message:', error);
            }
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
            console.log('Sending WebSocket message:', message);
            socket.send(JSON.stringify(message));
        } else {
            console.error('WebSocket is not open. Cannot send message:', message);
        }
    }, [socket]);

    return { sendMessage, lastMessage, isConnected, isConnecting };
};