import api from './api';

export const login = async (username, password) => {
    const response = await api.post('/login', { username, password });
    return response;
};

export const register = async (username, password) => {
    const response = await api.post('/register', { username, password });
    return response.data;
};

export const logout = () => {
    localStorage.removeItem('user');
};

export const getCurrentUser = () => {
    const user = localStorage.getItem('user');
    return user ? JSON.parse(user) : null;
};

export const setAuthToken = (token) => {
    if (token) {
        api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    } else {
        delete api.defaults.headers.common['Authorization'];
    }
};