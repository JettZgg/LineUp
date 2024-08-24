import React, { useState } from 'react';
import { Dialog, DialogTitle, DialogContent, DialogActions, TextField, Button } from '@mui/material';
import { styled } from '@mui/material/styles';
import { useNavigate } from 'react-router-dom';
import { joinMatch } from '../../services/api';

const StyledDialog = styled(Dialog)(({ theme }) => ({
    '& .MuiDialog-paper': {
        backgroundColor: '#DCC2C2',
        borderRadius: '10px',
    },
}));

const StyledTextField = styled(TextField)(({ theme }) => ({
    '& .MuiOutlinedInput-root': {
        backgroundColor: '#DCC2C2',
        borderRadius: '20px',
        '& fieldset': {
            borderColor: '#1E1E1E',
        },
        '&:hover fieldset': {
            borderColor: '#1E1E1E',
        },
        '&.Mui-focused fieldset': {
            borderColor: '#1E1E1E',
        },
    },
    '& .MuiInputBase-input': {
        color: '#1E1E1E',
        fontFamily: theme.typography.fontFamily,
    },
    '& .MuiInputLabel-root': {
        color: '#65558F',
    },
    '& .MuiInputLabel-root.Mui-focused': {
        color: '#65558F',
    },
}));

const StyledButton = styled(Button)(({ theme }) => ({
    color: '#65558F',
    '&:hover': {
        backgroundColor: 'rgba(101, 85, 143, 0.04)',
    },
}));

const JoinMatchModal = ({ open, onClose }) => {
    const [matchId, setMatchId] = useState('');
    const navigate = useNavigate();

    const handleJoin = async () => {
        try {
            const response = await joinMatch(matchId);
            if (response && response.data) {
                navigate(`/match/${matchId}/waiting`);
            }
        } catch (error) {
            console.error('Failed to join match:', error);
        }
        onClose();
    };

    const handleKeyDown = (event) => {
        if (event.key === 'Enter') {
            handleJoin();
        } else if (event.key === 'Escape') {
            onClose();
        }
    };

    return (
        <StyledDialog open={open} onClose={onClose}>
            <DialogTitle>Join a match</DialogTitle>
            <DialogContent>
                <StyledTextField
                    autoFocus
                    margin="dense"
                    id="matchId"
                    label="Match ID"
                    type="text"
                    fullWidth
                    variant="outlined"
                    value={matchId}
                    onChange={(e) => setMatchId(e.target.value)}
                    onKeyDown={handleKeyDown}
                />
            </DialogContent>
            <DialogActions>
                <StyledButton onClick={onClose}>
                    Cancel
                </StyledButton>
                <StyledButton onClick={handleJoin}>
                    OK
                </StyledButton>
            </DialogActions>
        </StyledDialog>
    );
};

export default JoinMatchModal;