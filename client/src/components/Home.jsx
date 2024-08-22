import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { TextField, Button, Typography, Box } from '@mui/material';
import { createMatch, joinMatch } from '../services/api';

const Home = () => {
    const [boardWidth, setBoardWidth] = useState(10);
    const [boardHeight, setBoardHeight] = useState(10);
    const [winLength, setWinLength] = useState(5);
    const [matchId, setMatchId] = useState('');
    const navigate = useNavigate();

    const handleCreateMatch = async (e) => {
        e.preventDefault();
        try {
            console.log('Sending create match request with:', { boardWidth, boardHeight, winLength });
            const response = await createMatch({ boardWidth, boardHeight, winLength });
            console.log('Create match response:', response);
            if (response && response.data) {
                const matchId = String(response.data.id || response.data.match?.id);
                if (matchId) {
                    navigate(`/match/${matchId}/waiting`);
                } else {
                    console.error('Invalid response from server: No match ID found', response);
                }
            } else {
                console.error('Invalid response from server:', response);
            }
        } catch (error) {
            console.error('Failed to create match:', error);
        }
    };

    const handleJoinMatch = async (e) => {
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
                alert(`Failed to join match: ${error.response.data.error || 'Unknown error'}`);
            } else {
                alert('Failed to join match. Please check your connection and try again.');
            }
        }
    };

    return (
        <Box sx={{ width: '100%', maxWidth: 400 }}>
            <Typography variant="h4" gutterBottom>Create Match</Typography>
            <Box component="form" onSubmit={handleCreateMatch} sx={{ mb: 4 }}>
                <TextField
                    fullWidth
                    margin="normal"
                    label="Board Width"
                    type="number"
                    value={boardWidth}
                    onChange={(e) => setBoardWidth(parseInt(e.target.value))}
                    inputProps={{ min: 3, max: 99 }}
                />
                <TextField
                    fullWidth
                    margin="normal"
                    label="Board Height"
                    type="number"
                    value={boardHeight}
                    onChange={(e) => setBoardHeight(parseInt(e.target.value))}
                    inputProps={{ min: 3, max: 99 }}
                />
                <TextField
                    fullWidth
                    margin="normal"
                    label="Win Length"
                    type="number"
                    value={winLength}
                    onChange={(e) => setWinLength(parseInt(e.target.value))}
                    inputProps={{ min: 3, max: 19 }}
                />
                <Button fullWidth variant="contained" type="submit" sx={{ mt: 2 }}>
                    Create Match
                </Button>
            </Box>

            <Typography variant="h4" gutterBottom>Join Match</Typography>
            <Box component="form" onSubmit={handleJoinMatch}>
                <TextField
                    fullWidth
                    margin="normal"
                    label="Match ID"
                    value={matchId}
                    onChange={(e) => setMatchId(e.target.value)}
                />
                <Button fullWidth variant="contained" type="submit" sx={{ mt: 2 }}>
                    Join Match
                </Button>
            </Box>
        </Box>
    );
};

export default Home;