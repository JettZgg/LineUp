// src/services/api.js
import axios from 'axios';

const API_URL = 'http://localhost:8080/api'; // Update with your server URL

const api = axios.create({
    baseURL: API_URL,
});

api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

export const register = (username, password) => api.post('/register', { username, password });
export const login = (username, password) => api.post('/login', { username, password });
export const createMatch = (config) => api.post('/create-match', config);
export const joinMatch = (matchId) => api.post(`/join-match/${matchId}`);

export default api;