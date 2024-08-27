import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { Box, Typography } from '@mui/material';
import { styled } from '@mui/material/styles';
import { useWebSocket } from '../services/websocket';
import { useAuth } from '../components/common/AuthContext';
import MatchBoard from '../components/Match/MatchBoard';
import MatchInfo from '../components/Match/MatchInfo';
import MatchControls from '../components/Match/MatchControls';

const StyledBox = styled(Box)(({ theme }) => ({
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    minHeight: '100vh',
    width: '100vw',
    backgroundColor: '#BF9D9D',
    padding: '20px',
}));

const MatchRoom = () => {
    const { matchId } = useParams();
    const { user } = useAuth();
    const [players, setPlayers] = useState([
        { id: user.userID, username: user.username },
        { id: null, username: 'Waiting' }
    ]);
    const [matchStarted, setMatchStarted] = useState(false);
    const { sendMessage, lastMessage, isConnected } = useWebSocket(matchId, user);

    useEffect(() => {
        if (lastMessage) {
            const { type, players: updatedPlayers, matchStarted } = lastMessage;
            if (['matchInfo', 'playerJoined', 'playerLeft'].includes(type)) {
                if (updatedPlayers) {
                    setPlayers(prevPlayers => {
                        const newPlayers = [...prevPlayers];
                        updatedPlayers.forEach((player, index) => {
                            newPlayers[index] = player;
                        });
                        return newPlayers;
                    });
                }
            } else if (type === 'matchStart') {
                setMatchStarted(true);
            }
        }
    }, [lastMessage]);

    useEffect(() => {
        if (isConnected) {
            sendMessage({ type: 'getMatchInfo', matchId });
        }
    }, [isConnected, matchId, sendMessage]);

    const handleStart = () => {
        sendMessage({ type: 'startMatch', matchId, token: user.token });
    };

    const handleCopyMatchId = () => {
        navigator.clipboard.writeText(matchId);
    };

    const handleCellClick = (x, y) => {
        if (matchStarted) {
            sendMessage({ type: 'move', matchId, x, y, token: user.token });
        }
    };

    return (
        <StyledBox>
            <Typography variant="h4" gutterBottom align="center" sx={{ fontFamily: 'Explora, cursive', fontSize: '6rem', color: '#1E1E1E', marginBottom: '1rem' }}>
                LineUp
            </Typography>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', width: '100%', marginBottom: '2rem' }}>
                <Box>
                    <MatchInfo matchId={matchId} players={players} onCopyMatchId={handleCopyMatchId} />
                    <MatchControls onStart={handleStart} matchStarted={matchStarted} />
                </Box>
                <MatchBoard onCellClick={handleCellClick} />
            </Box>
        </StyledBox>
    );
};

export default MatchRoom;
