// src/components/Game/GameBoard.jsx
import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useWebSocket } from '../../services/websocket';
import { Box, Typography, Button } from '@mui/material';

const GameBoard = () => {
    const { matchId } = useParams();
    const [gameState, setGameState] = useState(null);
    const { sendMessage, lastMessage } = useWebSocket(matchId);

    useEffect(() => {
        if (lastMessage) {
            const data = JSON.parse(lastMessage.data);
            setGameState(data);
        }
    }, [lastMessage]);

    const handleCellClick = (x, y) => {
        sendMessage({ type: 'move', matchId, x, y });
    };

    if (!gameState) {
        return <Typography>Loading...</Typography>;
    }

    return (
        <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
            <Typography variant="h4" gutterBottom>Match {matchId}</Typography>
            <Typography variant="h6" gutterBottom>
                {gameState.currentPlayer === gameState.playerTurn ? "Your turn" : "Opponent's turn"}
            </Typography>
            <Box sx={{ display: 'inline-grid', gridTemplateColumns: `repeat(${gameState.width}, 30px)`, gap: 1 }}>
                {Array.from({ length: gameState.height }, (_, y) => (
                    Array.from({ length: gameState.width }, (_, x) => (
                        <Button
                            key={`${x}-${y}`}
                            variant="outlined"
                            sx={{ minWidth: 30, height: 30, padding: 0 }}
                            onClick={() => handleCellClick(x, y)}
                            disabled={gameState.board[y][x] !== 0 || gameState.currentPlayer !== gameState.playerTurn}
                        >
                            {gameState.board[y][x] === 1 ? 'X' : gameState.board[y][x] === 2 ? 'O' : ''}
                        </Button>
                    ))
                ))}
            </Box>
            {gameState.winner && (
                <Typography variant="h5" sx={{ mt: 2 }}>
                    {gameState.winner === gameState.currentPlayer ? "You won!" : "You lost!"}
                </Typography>
            )}
        </Box>
    );
};

export default GameBoard;