import React, { useState } from 'react';
import { Typography, Box, Button } from '@mui/material';
import { styled } from '@mui/material/styles';
import { useAuth } from '../contexts/AuthContext';
import { useNavigate } from 'react-router-dom';
import { createMatch } from '../services/api';
import { useTheme } from '@mui/material/styles';
import JoinMatchModal from './Game/JoinMatchModal'; // Added this import

const StyledBox = styled(Box)(({ theme }) => ({
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    minHeight: '100vh',
    width: '100vw',
    backgroundColor: '#BF9D9D',
}));

const StyledButton = styled(Button)(({ theme }) => ({
    backgroundColor: '#DCC2C2',
    color: '#1E1E1E',
    border: '1px solid #1E1E1E',
    borderRadius: '8px',
    padding: '3px 20px',
    margin: '10px 0',
    width: '200px',
    fontSize: '1.4rem',
    '&:hover': {
        backgroundColor: '#C2B0B0',
    },
    '&:focus': {
        outline: 'none',
        boxShadow: 'none',
    },
    '&:active': {
        outline: 'none',
        boxShadow: 'none',
    },
}));

const LogoutButton = styled(StyledButton)(({ theme }) => ({
    color: '#B32D2D',
    position: 'absolute',
    bottom: '20px',
    left: '20px',
}));

const Home = () => {
    const { user, logout } = useAuth();
    const navigate = useNavigate();
    const theme = useTheme();
    const [joinModalOpen, setJoinModalOpen] = useState(false);

    const handlePlay = async () => {
        try {
            const response = await createMatch();
            console.log('Create match response:', response);
            if (response && response.data && response.data.match && response.data.match.id) {
                const matchId = response.data.match.id;
                console.log('Navigating to waiting room with id:', matchId);
                navigate(`/match/${matchId}/waiting`);
            } else {
                console.error('Invalid response from server:', response);
            }
        } catch (error) {
            console.error('Failed to create match:', error);
        }
    };

    const handleJoin = () => {
        setJoinModalOpen(true);
    };

    const handleLogout = () => {
        logout();
        navigate('/login');
    };

    return (
        <StyledBox>
            <Typography variant="h4" gutterBottom align="center" sx={{ fontSize: '6rem', color: '#1E1E1E', marginBottom: '2rem' }}>
                LineUp
            </Typography>
            <Typography variant="body1" sx={{ marginBottom: '1rem' }}>
                Username: {user.username}
            </Typography>
            <Typography variant="body1" sx={{ marginBottom: '2rem' }}>
                UID: {user.userID}
            </Typography>
            <StyledButton onClick={handlePlay}>Play</StyledButton>
            <StyledButton onClick={handleJoin}>Join</StyledButton>
            <StyledButton>History</StyledButton>
            <StyledButton>Settings</StyledButton>
            <LogoutButton onClick={handleLogout}>Logout →</LogoutButton>
            <JoinMatchModal open={joinModalOpen} onClose={() => setJoinModalOpen(false)} />
        </StyledBox>
    );
};

export default Home;