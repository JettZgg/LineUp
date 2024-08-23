// src/components/Game/GameBoard.jsx
import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useWebSocket } from '../../services/websocket';
import { useAuth } from '../../contexts/AuthContext';
import { Box, Typography, Button } from '@mui/material';

const GameBoard = () => {
    const { matchId } = useParams();
    const { user } = useAuth();
    const [gameState, setGameState] = useState(null);
    const { sendMessage, lastMessage } = useWebSocket(matchId);

    useEffect(() => {
        if (lastMessage) {
            const data = JSON.parse(lastMessage.data);
            if (data.type === 'gameState') {
                setGameState(data);
            }
        }
    }, [lastMessage]);

    useEffect(() => {
        sendMessage({ type: 'joinGame', matchId, token: user.token });
    }, [matchId, sendMessage, user.token]);

    const handleCellClick = (x, y) => {
        sendMessage({ type: 'move', matchId, x, y, token: user.token });
    };

    if (!gameState) {
        return <Typography>Loading...</Typography>;
    }

    return (
        <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
            <Typography variant="h4" gutterBottom>Match {matchId}</Typography>
            <Typography variant="h6" gutterBottom>
                {gameState.currentPlayer === user.username ? "Your turn" : `${gameState.currentPlayer}'s turn`}
            </Typography>
            <Box sx={{ display: 'inline-grid', gridTemplateColumns: `repeat(${gameState.boardWidth}, 30px)`, gap: 1 }}>
                {gameState.board.map((row, y) => (
                    row.map((cell, x) => (
                        <Button
                            key={`${x}-${y}`}
                            variant="outlined"
                            sx={{ minWidth: 30, height: 30, padding: 0 }}
                            onClick={() => handleCellClick(x, y)}
                            disabled={cell !== '' || gameState.currentPlayer !== user.username}
                        >
                            {cell}
                        </Button>
                    ))
                ))}
            </Box>
            {gameState.winner && (
                <Typography variant="h5" sx={{ mt: 2 }}>
                    {gameState.winner === user.username ? "You won!" : `${gameState.winner} won!`}
                </Typography>
            )}
        </Box>
    );
};

export default GameBoard;