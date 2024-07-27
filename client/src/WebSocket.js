let socket;
const messageQueue = [];

export const initWebSocket = () => {
    if (socket && socket.readyState === WebSocket.OPEN) {
        return; // WebSocket is already open
    }

    socket = new WebSocket('ws://localhost:8080');

    socket.onopen = () => {
        console.log('WebSocket connection established');
        while (messageQueue.length > 0) {
            socket.send(messageQueue.shift());
        }
    };

    socket.onclose = (event) => {
        console.log('WebSocket connection closed', event);
        // Reconnect after 1 second
        setTimeout(initWebSocket, 1000);
    };

    socket.onerror = (error) => {
        console.log('WebSocket error', error);
    };

    socket.onmessage = (event) => {
        console.log('Message from server ', event.data);
        if (typeof handleMessage === 'function') {
            handleMessage(event);
        }
    };
};

export const sendMessage = (message) => {
    const msg = JSON.stringify(message);
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(msg);
    } else {
        console.log('WebSocket connection is not open. Queueing message.');
        messageQueue.push(msg);
    }
};

let handleMessage;

export const setHandleMessage = (callback) => {
    handleMessage = callback;
};
