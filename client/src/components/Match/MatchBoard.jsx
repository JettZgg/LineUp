import React from 'react';
import { Box } from '@mui/material';
import { styled } from '@mui/material/styles';

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

const MatchBoard = ({ onCellClick }) => {
    return (
        <GameBoard>
            {Array.from({ length: 225 }).map((_, index) => (
                <GameCell
                    key={index}
                    onClick={() => onCellClick(index % 15, Math.floor(index / 15))}
                />
            ))}
        </GameBoard>
    );
};

export default MatchBoard;
