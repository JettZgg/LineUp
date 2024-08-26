import axios from 'axios';
import { API_BASE_URL } from '../config';

const api = axios.create({
    baseURL: API_BASE_URL,
});

export const createMatch = () => api.post('/create-match');
export const joinMatch = (matchId) => api.post(`/join-match/${matchId}`);

export default api;