// src/components/Game/CreateMatch.jsx
import React, { useState } from 'react';
import { createMatch } from '../../services/api';
import { useNavigate } from 'react-router-dom';
import { TextField, Button, Typography, Box } from '@mui/material';

const CreateMatch = () => {
    const [boardWidth, setBoardWidth] = useState(10);
    const [boardHeight, setBoardHeight] = useState(10);
    const [winLength, setWinLength] = useState(5);
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            console.log('Sending create match request with:', { boardWidth, boardHeight, winLength });
            const response = await createMatch({ boardWidth, boardHeight, winLength });
            console.log('Create match response:', response);

            if (response && response.data) {
                console.log('Response data:', response.data);
                const matchId = response.data.id || (response.data.match && response.data.match.id);
                console.log('Match ID:', matchId);
                if (matchId) {
                    console.log('Navigating to waiting room with id:', matchId);
                    navigate(`/match/${matchId}/waiting`);
                } else {
                    console.error('No match ID found in response:', response.data);
                }
            } else {
                console.error('Invalid response from server:', response);
            }
        } catch (error) {
            console.error('Failed to create match:', error);
            if (error.response) {
                console.error('Error response:', error.response.data);
            }
        }
    };

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ mt: 3 }}>
            <TextField
                margin="normal"
                required
                fullWidth
                id="boardWidth"
                label="Board Width"
                type="number"
                value={boardWidth}
                onChange={(e) => setBoardWidth(parseInt(e.target.value))}
                inputProps={{ min: "3", max: "99" }}
            />
            <TextField
                margin="normal"
                required
                fullWidth
                id="boardHeight"
                label="Board Height"
                type="number"
                value={boardHeight}
                onChange={(e) => setBoardHeight(parseInt(e.target.value))}
                inputProps={{ min: "3", max: "99" }}
            />
            <TextField
                margin="normal"
                required
                fullWidth
                id="winLength"
                label="Win Length"
                type="number"
                value={winLength}
                onChange={(e) => setWinLength(parseInt(e.target.value))}
                inputProps={{ min: "3", max: "19" }}
            />
            <Button
                type="submit"
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
            >
                Create Match
            </Button>
        </Box>
    );
};

export default CreateMatch;