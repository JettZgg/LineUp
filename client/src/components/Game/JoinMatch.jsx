// src/components/Game/JoinMatch.jsx
import React, { useState } from 'react';
import { joinMatch } from '../../services/api';
import { useNavigate } from 'react-router-dom';
import { TextField, Button, Box } from '@mui/material';

const JoinMatch = () => {
    const [matchId, setMatchId] = useState('');
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            console.log('Attempting to join match:', matchId);
            const response = await joinMatch(matchId);
            console.log('Join match response:', response);
            if (response && response.data) {
                navigate(`/match/${matchId}/waiting`);
            } else {
                console.error('Invalid response from server:', response);
            }
        } catch (error) {
            console.error('Failed to join match:', error);
            if (error.response) {
                console.error('Error response:', error.response.data);
            }
        }
    };

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ display: 'flex', flexDirection: 'column', alignItems: 'stretch', mt: 2 }}>
            <TextField
                value={matchId}
                onChange={(e) => setMatchId(e.target.value)}
                placeholder="Match ID"
                required
                sx={{ mb: 1 }}
                fullWidth
            />
            <Button type="submit" variant="contained" color="primary">
                Join Match
            </Button>
        </Box>
    );
};

export default JoinMatch;