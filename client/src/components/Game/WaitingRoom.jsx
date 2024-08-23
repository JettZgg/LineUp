import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Box, Typography, Button, Grid } from '@mui/material';
import { useWebSocket } from '../../services/websocket';
import { useAuth } from '../../contexts/AuthContext';

const WaitingRoom = () => {
    const { matchId } = useParams();
    const { user } = useAuth();
    const navigate = useNavigate();
    const [players, setPlayers] = useState([]);
    const [isReady, setIsReady] = useState(false);
    const [gameConfig, setGameConfig] = useState(null);
    const { sendMessage, lastMessage } = useWebSocket(matchId);

    useEffect(() => {
        if (lastMessage) {
            const data = JSON.parse(lastMessage.data);
            console.log("Received WebSocket message:", data);
            if (data.type === 'gameInfo') {
                setPlayers(data.players);
                setGameConfig(data.config);
            } else if (data.type === 'playerJoined' || data.type === 'playerLeft') {
                setPlayers(data.players);
            } else if (data.type === 'gameStart') {
                navigate(`/match/${matchId}`);
            }
        }
    }, [lastMessage, matchId, navigate]);

    useEffect(() => {
        sendMessage({ type: 'joinMatch', matchId, token: user.token });
    }, [matchId, sendMessage, user.token]);

    const handleReady = () => {
        setIsReady(true);
        sendMessage({ type: 'playerReady', matchId, token: user.token });
    };

    const handleStart = () => {
        sendMessage({ type: 'startGame', matchId, token: user.token });
    };

    return (
        <Box sx={{ mt: 4 }}>
            <Typography variant="h4" gutterBottom>Waiting Room - Match {matchId}</Typography>
            <Grid container spacing={2}>
                <Grid item xs={12}>
                    <Typography variant="h6">Players:</Typography>
                    <Typography>Player 1: {players[0]?.username || user.username}</Typography>
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