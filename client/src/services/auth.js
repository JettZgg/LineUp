// src/services/auth.js
import api from './api';

export const login = async (username, password) => {
    const response = await api.post('/login', { username, password });
    return response.data;
};

export const register = async (username, password) => {
    const response = await api.post('/register', { username, password });
    return response.data;
};

export const logout = () => {
    // Clear token from localStorage
    localStorage.removeItem('token');
};

export const getCurrentUser = () => {
    const token = localStorage.getItem('token');
    if (token) {
        // You might want to decode the JWT token here to get user info
        return { token };
    }
    return null;
};