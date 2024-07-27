import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { sendMessage } from '../WebSocket';

function Register() {
    const [email, setEmail] = useState('');
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const navigate = useNavigate();

    const handleSubmit = (e) => {
        e.preventDefault();
        sendMessage({
            action: 'register',
            email,
            username,
            password
        });
    };

    // Listen for WebSocket messages
    React.useEffect(() => {
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

        const socket = new WebSocket('ws://localhost:8080'); // Replace with your server address
        socket.addEventListener('message', handleMessage);

        return () => {
            socket.removeEventListener('message', handleMessage);
            socket.close();
        };
    }, [navigate]);

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
