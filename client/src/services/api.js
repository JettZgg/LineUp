import axios from 'axios';
import { API_BASE_URL } from '../config';

const api = axios.create({
    baseURL: API_BASE_URL,
});

export const createMatch = (token) => api.post('/create-match', {}, {
    headers: { Authorization: `Bearer ${token}` }
});
export const joinMatch = (matchId, token) => api.post(`/join-match/${matchId}`, {}, {
    headers: { Authorization: `Bearer ${token}` }
});

export default api;