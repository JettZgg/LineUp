// src/App.jsx
import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import Container from '@mui/material/Container';
import Box from '@mui/material/Box';
import Header from './components/common/Header';
import Login from './components/Auth/Login';
import Register from './components/Auth/Register';
import Home from './components/Home';
import GameBoard from './components/Game/GameBoard';
import ProtectedRoute from './components/common/ProtectedRoute';
import WaitingRoom from './components/Game/WaitingRoom';

const theme = createTheme();

const App = () => {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <AuthProvider>
        <Router>
          <Header />
          <Container maxWidth="md" sx={{ display: 'flex', flexDirection: 'column', justifyContent: 'flex-start', minHeight: '100vh', pt: 4 }}>
            <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
              <Routes>
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />
                <Route path="/" element={<ProtectedRoute><Home /></ProtectedRoute>} />
                <Route path="/match/:matchId/waiting" element={<ProtectedRoute><WaitingRoom /></ProtectedRoute>} />
                <Route path="/match/:matchId" element={<ProtectedRoute><GameBoard /></ProtectedRoute>} />
              </Routes>
            </Box>
          </Container>
        </Router>
      </AuthProvider>
    </ThemeProvider>
  );
};

export default App;