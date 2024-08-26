import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Box, Typography, Button, IconButton } from '@mui/material';
import { styled } from '@mui/material/styles';
import { useWebSocket } from '../../services/websocket';
import { useAuth } from '../common/AuthContext';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';

const StyledBox = styled(Box)(({ theme }) => ({
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    minHeight: '100vh',
    width: '100vw',
    backgroundColor: '#BF9D9D',
    padding: '20px',
}));

const StyledButton = styled(Button)(({ theme }) => ({
    backgroundColor: '#DCC2C2',
    color: '#1E1E1E',
    border: '1px solid #1E1E1E',
    borderRadius: '8px',
    padding: '3px 20px',
    margin: '10px 0',
    width: '200px',
    fontFamily: theme.typography.fontFamily,
    fontSize: '1.4rem',
    '&:hover': {
        backgroundColor: '#C2B0B0',
    },
}));

const GameBoard = styled(Box)(({ theme }) => ({
    display: 'grid',
    gridTemplateColumns: 'repeat(15, 1fr)',
    gridTemplateRows: 'repeat(15, 1fr)',
    gap: '1px',
    width: '80vmin',
    height: '80vmin',
    backgroundColor: '#DCC2C2',
    border: '2px solid #1E1E1E',
}));

const GameCell = styled(Box)(({ theme }) => ({
    width: '100%',
    height: '100%',
    backgroundColor: '#BF9D9D',
    '&:hover': {
        backgroundColor: '#A88A8A',
    },
}));

const GameRoom = () => {
    const { matchId } = useParams();
    const { user } = useAuth();
    const navigate = useNavigate();
    const [players, setPlayers] = useState([]);
    const [gameStarted, setGameStarted] = useState(false);
    const { sendMessage, lastMessage, isConnected } = useWebSocket(matchId, user);

    useEffect(() => {
        if (lastMessage) {
            const { type, players: updatedPlayers, gameStarted } = lastMessage;
            if (['gameInfo', 'playerJoined', 'playerLeft'].includes(type)) {
                if (updatedPlayers) setPlayers(updatedPlayers);
            } else if (type === 'gameStart') {
                setGameStarted(true);
            }
        }
    }, [lastMessage]);

    const handleStart = () => {
        sendMessage({ type: 'startGame', matchId, token: user.token });
    };

    const handleCopyMatchId = () => {
        navigator.clipboard.writeText(matchId);
    };

    const handleCellClick = (x, y) => {
        if (gameStarted) {
            sendMessage({ type: 'move', matchId, x, y, token: user.token });
        }
    };

    return (
        <StyledBox>
            <Typography variant="h4" gutterBottom align="center" sx={{ fontFamily: 'Explora, cursive', fontSize: '6rem', color: '#1E1E1E', marginBottom: '1rem' }}>
                LineUp
            </Typography>
            <Box sx={{ display: 'flex', alignItems: 'center', marginBottom: '2rem' }}>
                <Typography variant="body1" sx={{ marginRight: '0.5rem' }}>
                    Match ID: {matchId}
                </Typography>
                <IconButton onClick={handleCopyMatchId} sx={{ '&:focus': { outline: 'none' } }}>
                    <ContentCopyIcon />
                </IconButton>
            </Box>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', width: '100%', marginBottom: '2rem' }}>
                <Box>
                    {players.map((player, index) => (
                        <Typography key={player.id} variant="body1" sx={{ fontWeight: 600, marginBottom: '0.5rem' }}>
                            Player{index + 1}: {player.username}
                        </Typography>
                    ))}
                    {players.length === 2 && !gameStarted && (
                        <StyledButton onClick={handleStart}>Start</StyledButton>
                    )}
                </Box>
                <GameBoard>
                    {Array.from({ length: 225 }).map((_, index) => (
                        <GameCell
                            key={index}
                            onClick={() => handleCellClick(index % 15, Math.floor(index / 15))}
                        />
                    ))}
                </GameBoard>
            </Box>
        </StyledBox>
    );
};

export default GameRoom;
