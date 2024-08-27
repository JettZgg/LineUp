import React from 'react';
import { Box, Typography, IconButton } from '@mui/material';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';

const MatchInfo = ({ matchId, players, onCopyMatchId }) => {
    return (
        <Box>
            <Box sx={{ display: 'flex', alignItems: 'center', marginBottom: '2rem' }}>
                <Typography variant="body1" sx={{ marginRight: '0.5rem' }}>
                    Match ID: {matchId}
                </Typography>
                <IconButton onClick={onCopyMatchId} sx={{ '&:focus': { outline: 'none' } }}>
                    <ContentCopyIcon />
                </IconButton>
            </Box>
            {players.map((player, index) => (
                <Typography key={player.id} variant="body1" sx={{ fontWeight: 600, marginBottom: '0.5rem' }}>
                    Player{index + 1}: {player.username}
                </Typography>
            ))}
        </Box>
    );
};

export default MatchInfo;
