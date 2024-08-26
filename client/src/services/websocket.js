import { useState, useEffect, useCallback, useRef } from 'react';
import { WS_BASE_URL } from '../config';

export const useWebSocket = (matchId, user) => {
    const [socket, setSocket] = useState(null);
    const [lastMessage, setLastMessage] = useState(null);
    const [isConnecting, setIsConnecting] = useState(true);
    const [isConnected, setIsConnected] = useState(false);
    const messageQueue = useRef([]);
    const sentMessages = useRef(new Set());

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
            // Send queued messages
            while (messageQueue.current.length > 0) {
                const message = messageQueue.current.shift();
                sendMessageInternal(ws, message);
            }
        };

        ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                setLastMessage(data);
            } catch (error) {
                console.error('Error parsing WebSocket message:', error);
            }
        };

        ws.onclose = () => {
            console.log('WebSocket disconnected');
            setIsConnected(false);
            setIsConnecting(false);
        };

        setSocket(ws);

        return () => {
            ws.close();
        };
    }, [matchId, user]);

    useEffect(() => {
        connectWebSocket();
        return () => {
            if (socket) {
                socket.close();
            }
        };
    }, [connectWebSocket]);

    const sendMessageInternal = (ws, message) => {
        const messageString = JSON.stringify(message);
        if (!sentMessages.current.has(messageString)) {
            ws.send(messageString);
            sentMessages.current.add(messageString);
        }
    };

    const sendMessage = useCallback((message) => {
        if (socket && socket.readyState === WebSocket.OPEN) {
            sendMessageInternal(socket, message);
        } else {
            messageQueue.current.push(message);
        }
    }, [socket]);

    return { sendMessage, lastMessage, isConnected, isConnecting };
};