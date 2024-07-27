let socket = null;

export const connectWebSocket = (url) => {
    socket = new WebSocket(url);

    socket.onopen = () => {
        console.log('WebSocket connection established');
    };

    socket.onclose = () => {
        console.log('WebSocket connection closed');
    };

    return socket;
};

export const sendMessage = (message) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(message));
    } else {
        console.error('WebSocket is not connected');
    }
};

export const closeWebSocket = () => {
    if (socket) {
        socket.close();
    }
};