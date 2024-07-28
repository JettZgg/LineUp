import React, { useState, useEffect } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { sendMessage, initWebSocket, setHandleMessage } from '../WebSocket';

function Login() {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const navigate = useNavigate();

    useEffect(() => {
        initWebSocket();

        const handleMessage = (event) => {
            const data = JSON.parse(event.data);
            if (data.action === 'login_result') {
                if (data.success) {
                    localStorage.setItem('sessionId', data.session_id);
                    localStorage.setItem('username', data.username);
                    navigate('/lobby');
                } else {
                    setError('Login failed. Please check your credentials and try again.');
                }
            }
        };

        setHandleMessage(handleMessage);
    }, [navigate]);

    const handleSubmit = (e) => {
        e.preventDefault();
        sendMessage({
            action: 'login',
            email,
            password
        });
    };

    return (
        <div>
            <h2>Login</h2>
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
                    <label>Password:</label>
                    <input
                        type="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                    />
                </div>
                <button type="submit">Login</button>
            </form>
            <p>Don't have an account? <Link to="/register">Register</Link></p>
        </div>
    );
}

export default Login;