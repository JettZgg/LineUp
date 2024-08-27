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
            <Typography variant="body1" sx={{ fontWeight: 600, marginBottom: '0.5rem' }}>
                Player1: {players[0] ? players[0].username : 'Waiting'}
            </Typography>
            <Typography variant="body1" sx={{ fontWeight: 600, marginBottom: '0.5rem' }}>
                Player2: {players[1] ? players[1].username : 'Waiting'}
            </Typography>
        </Box>
    );
};

export default MatchInfo;
