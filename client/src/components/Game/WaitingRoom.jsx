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
    const [gameConfig, setGameConfig] = useState(null);
    const { sendMessage, lastMessage } = useWebSocket(matchId);

    useEffect(() => {
        if (lastMessage) {
            const data = JSON.parse(lastMessage.data);
            if (data.type === 'playerJoined' || data.type === 'playerLeft') {
                setPlayers(data.players);
            } else if (data.type === 'gameConfig') {
                setGameConfig(data.config);
            } else if (data.type === 'gameStart') {
                // Navigate to the game board or update state to start the game
            }
        }
    }, [lastMessage]);

    useEffect(() => {
        // Request game configuration and player information when component mounts
        sendMessage({ type: 'getGameInfo', matchId });
    }, [matchId, sendMessage]);

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
                <Grid item xs={12}>
                    <Typography variant="h6">Players:</Typography>
                    <Typography>Player 1: {players[0]?.username || 'Waiting...'}</Typography>
                    <Typography>Player 2: {players[1]?.username || 'Waiting...'}</Typography>
                </Grid>
                {gameConfig && (
                    <Grid item xs={12}>
                        <Typography variant="h6">Game Configuration:</Typography>
                        <Typography>Board Width: {gameConfig.boardWidth}</Typography>
                        <Typography>Board Height: {gameConfig.boardHeight}</Typography>
                        <Typography>Length to Win: {gameConfig.winLength}</Typography>
                    </Grid>
                )}
                <Grid item xs={12}>
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
            </Grid>
        </Box>
    );
};

export default WaitingRoom;