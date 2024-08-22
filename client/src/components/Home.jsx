import React from 'react';
import { Typography, Box } from '@mui/material';
import CreateMatch from './Game/CreateMatch';
import JoinMatch from './Game/JoinMatch';

const Home = () => {
    return (
        <Box sx={{ width: '100%', maxWidth: 600, mx: 'auto' }}>
            <Typography variant="h4" gutterBottom>Create Match</Typography>
            <CreateMatch />

            <Typography variant="h4" gutterBottom sx={{ mt: 3 }}>Join Match</Typography>
            <JoinMatch />
        </Box>
    );
};

export default Home;