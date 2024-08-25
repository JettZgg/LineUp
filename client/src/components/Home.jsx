import React, { useState } from 'react';
import { Typography, Box, Button } from '@mui/material';
import { styled } from '@mui/material/styles';
import { useAuth } from './common/AuthContext';
import { useNavigate } from 'react-router-dom';
import { createMatch } from '../services/api';
import JoinMatchModal from './Game/JoinMatchModal';
import PageLayout from './layout/PageLayout';

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
    const [joinModalOpen, setJoinModalOpen] = useState(false);

    const handlePlay = async () => {
        try {
            const response = await createMatch();
            console.log('Create match response:', response);
            if (response && response.data && response.data.match && response.data.match.id) {
                const matchId = response.data.match.id;
                const initialPlayers = [response.data.player];
                console.log('Navigating to waiting room with id:', matchId);
                navigate(`/match/${matchId}/waiting`, { state: { initialPlayers } });
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
        <PageLayout>
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
        </PageLayout>
    );
};

export default Home;