// src/services/api.js
import axios from 'axios';

const API_URL = 'http://localhost:8080/api'; // Update with your server URL

const api = axios.create({
    baseURL: API_URL,
});

api.interceptors.request.use((config) => {
    const user = JSON.parse(localStorage.getItem('user'));
    if (user && user.token) {
        config.headers.Authorization = `Bearer ${user.token}`;
    }
    return config;
});

export const register = (username, password) => api.post('/register', { username, password });
export const login = (username, password) => api.post('/login', { username, password });
export const createMatch = () => api.post('/create-match');
export const joinMatch = (matchId) => api.post(`/join-match/${matchId}`);

export default api;