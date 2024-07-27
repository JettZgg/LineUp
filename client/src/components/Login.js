import React, { useState } from 'react';
import { useHistory } from 'react-router-dom';
import { sendMessage } from '../WebSocket';

function Login() {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const history = useHistory();

    const handleSubmit = (e) => {
        e.preventDefault();
        sendMessage({
            action: 'login',
            email,
            password
        });
    };

    // Listen for WebSocket messages
    React.useEffect(() => {
        const handleMessage = (event) => {
            const data = JSON.parse(event.data);
            if (data.action === 'login_result') {
                if (data.success) {
                    localStorage.setItem('sessionId', data.session_id);
                    localStorage.setItem('username', data.username);
                    history.push('/lobby');
                } else {
                    setError('Login failed. Please check your credentials and try again.');
                }
            }
        };

        const socket = new WebSocket('ws://localhost:8080'); // Replace with your server address
        socket.addEventListener('message', handleMessage);

        return () => {
            socket.removeEventListener('message', handleMessage);
            socket.close();
        };
    }, [history]);

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
        </div>
    );
}

export default Login;