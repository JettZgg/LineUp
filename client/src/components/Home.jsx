import React from 'react';
import { Typography, Box, Button } from '@mui/material';
import { styled } from '@mui/material/styles';
import { useAuth } from '../contexts/AuthContext';
import { useNavigate } from 'react-router-dom';

const StyledBox = styled(Box)({
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    minHeight: '100vh',
    width: '100vw',
    backgroundColor: '#BF9D9D',
});

const StyledButton = styled(Button)({
    backgroundColor: '#DCC2C2',
    color: '#1E1E1E',
    border: '1px solid #1E1E1E',
    borderRadius: '8px',
    padding: '5px 20px',
    margin: '10px 0',
    width: '200px',
    fontFamily: 'Lora, serif',
    fontSize: '1.2rem', // Increase this value to make the text bigger
    '&:hover': {
        backgroundColor: '#C2B0B0',
    },
});

const LogoutButton = styled(StyledButton)({
    color: '#B32D2D',
    position: 'absolute',
    bottom: '20px',
    left: '20px',
});

const Home = () => {
    const { user, logout } = useAuth();
    const navigate = useNavigate();

    const handlePlay = () => {
        navigate('/match/new/waiting');
    };

    const handleLogout = () => {
        logout();
        navigate('/login');
    };

    return (
        <StyledBox>
            <Typography variant="h4" gutterBottom align="center" sx={{ fontFamily: 'Explora, cursive', fontSize: '6rem', color: '#1E1E1E', marginBottom: '2rem' }}>
                LineUp
            </Typography>
            <Typography variant="body1" sx={{ fontFamily: 'Lora, serif', fontWeight: 600, marginBottom: '1rem' }}>
                Username: {user.username}
            </Typography>
            <Typography variant="body1" sx={{ fontFamily: 'Lora, serif', fontWeight: 600, marginBottom: '2rem' }}>
                UID: {user.userID}
            </Typography>
            <StyledButton onClick={handlePlay}>Play</StyledButton>
            <StyledButton>Join</StyledButton>
            <StyledButton>History</StyledButton>
            <StyledButton>Settings</StyledButton>
            <LogoutButton onClick={handleLogout}>Logout →</LogoutButton>
        </StyledBox>
    );
};

export default Home;