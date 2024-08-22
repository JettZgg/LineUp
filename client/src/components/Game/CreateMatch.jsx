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
            const response = await createMatch({ boardWidth, boardHeight, winLength });
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

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ mt: 3 }}>
            <Typography variant="h5" gutterBottom>
                Create New Match
            </Typography>
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