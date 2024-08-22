import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { Box, Typography, Button, Grid } from '@mui/material';
import { useWebSocket } from '../../services/websocket';
import { useAuth } from '../../contexts/AuthContext';

const WaitingRoom = () => {
    const { matchId } = useParams();
    const { user } = useAuth();
    const [players, setPlayers] = useState([]);
    const [isReady, setIsReady] = useState(false);
    const { sendMessage, lastMessage } = useWebSocket(String(matchId));

    useEffect(() => {
        if (lastMessage) {
            const data = JSON.parse(lastMessage.data);
            if (data.type === 'playerJoined' || data.type === 'playerLeft') {
                setPlayers(data.players);
            } else if (data.type === 'gameStart') {
                // Navigate to the game board or update state to start the game
            }
        }
    }, [lastMessage]);

    const handleReady = () => {
        setIsReady(true);
        sendMessage({ type: 'playerReady', matchId, playerId: user.id });
    };

    const handleStart = () => {
        sendMessage({ type: 'startGame', matchId });
    };

    return (
        <Box sx={{ mt: 4 }}>
            <Typography variant="h4" gutterBottom>Waiting Room - Match {matchId}</Typography>
            <Grid container spacing={2}>
                <Grid item xs={4}>
                    <Typography variant="h6">Players:</Typography>
                    {players.map((player, index) => (
                        <Typography key={index}>{player.username} {player.ready ? '(Ready)' : ''}</Typography>
                    ))}
                </Grid>
                <Grid item xs={4}>
                    <Box sx={{ width: '100%', height: '300px', bgcolor: 'grey.300' }}>
                        {/* Placeholder for the board preview */}
                        <Typography variant="h6" sx={{ textAlign: 'center', pt: 2 }}>Board Preview</Typography>
                    </Box>
                    {players.length === 2 && players.every(p => p.ready) && (
                        <Button 
                            variant="contained" 
                            color="primary" 
                            fullWidth 
                            sx={{ mt: 2 }}
                            onClick={handleStart}
                        >
                            Start Game
                        </Button>
                    )}
                </Grid>
                <Grid item xs={4}>
                    <Typography variant="h6">
                        {players.length < 2 ? 'Waiting for opponent...' : 'Opponent joined!'}
                    </Typography>
                    {!isReady && (
                        <Button 
                            variant="contained" 
                            color="secondary" 
                            fullWidth 
                            sx={{ mt: 2 }}
                            onClick={handleReady}
                        >
                            Ready
                        </Button>
                    )}
                </Grid>
            </Grid>
        </Box>
    );
};

export default WaitingRoom;