import React, { useState } from 'react';
import { login as loginApi } from '../../services/auth';
import { useAuth } from '../common/AuthContext';
import { useNavigate } from 'react-router-dom';
import { TextField, Button, Typography, Box } from '@mui/material';
import { styled } from '@mui/material/styles';
import PageLayout from '../layout/PageLayout';

const StyledForm = styled(Box)({
    backgroundColor: '#DCC2C2',
    padding: '2rem',
    borderRadius: '8px',
    width: '300px',
    border: '1px solid #1E1E1E',
    boxShadow: '0px 4px 10px rgba(0, 0, 0, 0.1)',
});

const StyledButton = styled(Button)({
    backgroundColor: '#65558F',
    color: '#F5F5F5',
    border: '1px solid #1E1E1E',
    '&:hover': {
        backgroundColor: '#5048C5',
    },
});

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
        fontFamily: theme.typography.fontFamily,
        fontWeight: theme.typography.fontWeightRegular,
    },
    '& .MuiInputLabel-root': {
        color: '#65558F',
    },
    '& .MuiInputLabel-root.Mui-focused': {
        color: '#65558F',
    },
    '& input:-webkit-autofill': {
        WebkitBoxShadow: '0 0 0 1000px #DCC2C2 inset',
    },
    marginBottom: '10px',
}));

const Login = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const { login } = useAuth();
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const response = await loginApi(username, password);
            if (response && response.data) {
                login({
                    token: response.data.token,
                    userID: response.data.userID,
                    username: response.data.username
                });
                navigate('/');
            } else {
                console.error('Invalid response from server:', response);
            }
        } catch (error) {
            console.error('Login failed:', error);
        }
    };

    return (
        <PageLayout>
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
                    autoComplete="current-password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                />
                <StyledButton
                    type="submit"
                    fullWidth
                    variant="contained"
                    sx={{ mt: '1rem', mb: 2 }}
                >
                    Sign In
                </StyledButton>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', color: '#1E1E1E' }}>
                    <Typography variant="body2" sx={{ textDecoration: 'underline' }}>Forgot password?</Typography>
                    <Typography
                        variant="body2"
                        component="a"
                        href="/register"
                        sx={{ textDecoration: 'underline', color: '#1E1E1E' }}
                    >
                        Register
                    </Typography>
                </Box>
            </StyledForm>
        </PageLayout>
    );
};

export default Login;