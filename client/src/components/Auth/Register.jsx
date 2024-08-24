// src/components/Auth/Register.jsx
import React, { useState } from 'react';
import { register } from '../../services/api';
import { useAuth } from '../../contexts/AuthContext';
import { useNavigate } from 'react-router-dom';
import { TextField, Button, Typography, Box } from '@mui/material';
import { styled } from '@mui/material/styles';
import { useTheme } from '@mui/material/styles';

const StyledBox = styled(Box)(({ theme }) => ({
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    minHeight: '100vh',
    width: '100vw',
    backgroundColor: '#BF9D9D',
}));

const StyledForm = styled(Box)(({ theme }) => ({
    backgroundColor: '#DCC2C2',
    padding: '2rem',
    borderRadius: '8px',
    width: '300px',
    border: '1px solid #1E1E1E',
    boxShadow: '0px 4px 10px rgba(0, 0, 0, 0.1)',
}));

const StyledButton = styled(Button)(({ theme }) => ({
    backgroundColor: '#65558F',
    color: '#F5F5F5',
    border: '1px solid #1E1E1E',
    '&:hover': {
        backgroundColor: '#5048C5',
    },
}));

const StyledTextField = styled(TextField)(({ theme }) => ({
    '& .MuiOutlinedInput-root': {
        '& fieldset': {
            borderColor: '#1E1E1E',
        },
        '&:hover fieldset': {
            borderColor: '#1E1E1E',
        },
        '&.Mui-focused fieldset': {
            borderColor: '#1E1E1E',
        },
    },
    '& .MuiInputBase-input': {
        color: '#1E1E1E',
    },
    '& .MuiInputLabel-root': {
        color: '#65558F',
    },
    '& .MuiInputLabel-root.Mui-focused': {
        color: '#65558F',
    },
    '& input:-webkit-autofill': {
        WebkitBoxShadow: '0 0 0 1000px #DCC2C2 inset',
        WebkitTextFillColor: '#1E1E1E',
    },
    marginBottom: '10px',
}));

const Register = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const { login } = useAuth();
    const navigate = useNavigate();
    const theme = useTheme();

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (password !== confirmPassword) {
            console.error('Passwords do not match');
            return;
        }
        try {
            const response = await register(username, password);
            login(response.data);
            navigate('/');
        } catch (error) {
            console.error('Registration failed:', error);
        }
    };

    return (
        <StyledBox>
            <Typography variant="h4" gutterBottom align="center" sx={{ fontSize: '6rem', color: '#1E1E1E', marginBottom: '2rem' }}>
                LineUp
            </Typography>
            <StyledForm component="form" onSubmit={handleSubmit}>
                <StyledTextField
                    fullWidth
                    id="username"
                    label="Username"
                    name="username"
                    autoComplete="username"
                    autoFocus
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                />
                <StyledTextField
                    fullWidth
                    name="password"
                    label="Password"
                    type="password"
                    id="password"
                    autoComplete="new-password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                />
                <StyledTextField
                    fullWidth
                    name="confirmPassword"
                    label="Confirm Password"
                    type="password"
                    id="confirmPassword"
                    autoComplete="new-password"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                />
                <StyledButton
                    type="submit"
                    fullWidth
                    variant="contained"
                    sx={{ mt: '1rem', mb: 2 }}
                >
                    Register
                </StyledButton>
            </StyledForm>
        </StyledBox>
    );
};

export default Register;