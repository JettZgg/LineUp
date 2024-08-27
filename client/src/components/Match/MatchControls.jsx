import React from 'react';
import { Button } from '@mui/material';
import { styled } from '@mui/material/styles';

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

const MatchControls = ({ onStart, matchStarted }) => {
    return (
        <>
            {!matchStarted && (
                <StyledButton onClick={onStart}>Start</StyledButton>
            )}
        </>
    );
};

export default MatchControls;
