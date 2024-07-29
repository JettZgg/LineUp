// src/components/common/Header.jsx
import React from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';

const Header = () => {
    const { user, logout } = useAuth();

    return (
        <header>
            <nav>
                <ul>
                    <li><Link to="/">Home</Link></li>
                    {user ? (
                        <>
                            <li><Link to="/create-match">Create Match</Link></li>
                            <li><Link to="/join-match">Join Match</Link></li>
                            <li><button onClick={logout}>Logout</button></li>
                        </>
                    ) : (
                        <>
                            <li><Link to="/login">Login</Link></li>
                            <li><Link to="/register">Register</Link></li>
                        </>
                    )}
                </ul>
            </nav>
        </header>
    );
};

export default Header;