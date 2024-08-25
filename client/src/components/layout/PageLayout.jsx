import React from 'react';
import { Box, Typography } from '@mui/material';
import { styled } from '@mui/material/styles';

const StyledBox = styled(Box)(({ theme }) => ({
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    minHeight: '100vh',
    width: '100vw',
    backgroundColor: '#BF9D9D',
}));

const PageLayout = ({ children }) => {
    return (
        <StyledBox>
            <Typography variant="h4" gutterBottom align="center" sx={{ fontFamily: 'Explora, cursive', fontSize: '6rem', color: '#1E1E1E', marginBottom: '2rem' }}>
                LineUp
            </Typography>
            {children}
        </StyledBox>
    );
};

export default PageLayout;
