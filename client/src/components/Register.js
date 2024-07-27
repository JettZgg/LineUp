import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { sendMessage, initWebSocket, setHandleMessage } from '../WebSocket';

function Register() {
    const [email, setEmail] = useState('');
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const navigate = useNavigate();

    useEffect(() => {
        initWebSocket();

        const handleMessage = (event) => {
            const data = JSON.parse(event.data);
            if (data.action === 'register_result') {
                if (data.success) {
                    navigate('/login');
                } else {
                    setError('Registration failed. Please try again.');
                }
            }
        };

        setHandleMessage(handleMessage);
    }, [navigate]);

    const handleSubmit = (e) => {
        e.preventDefault();
        sendMessage({
            action: 'register',
            email,
            username,
            password
        });
    };

    return (
        <div>
            <h2>Register</h2>
            {error && <p style={{ color: 'red' }}>{error}</p>}
            <form onSubmit={handleSubmit}>
                <div>
                    <label>Email:</label>
                    <input
                        type="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        required
                    />
                </div>
                <div>
                    <label>Username:</label>
                    <input
                        type="text"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                        required
                    />
                </div>
                <div>
                    <label>Password:</label>
                    <input
                        type="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                    />
                </div>
                <button type="submit">Register</button>
            </form>
        </div>
    );
}

export default Register;
