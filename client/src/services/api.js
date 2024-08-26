import axios from 'axios';
import { API_BASE_URL } from '../config';

const api = axios.create({
    baseURL: API_BASE_URL,
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