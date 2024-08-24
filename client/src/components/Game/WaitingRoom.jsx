import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Box, Typography, Button, TextField, IconButton } from '@mui/material';
import { styled } from '@mui/material/styles';
import { useWebSocket } from '../../services/websocket';
import { useAuth } from '../../contexts/AuthContext';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import EmojiEventsIcon from '@mui/icons-material/EmojiEvents';
import CheckBoxOutlineBlankIcon from '@mui/icons-material/CheckBoxOutlineBlank';
import CheckBoxIcon from '@mui/icons-material/CheckBox';

const StyledBox = styled(Box)(({ theme }) => ({
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    minHeight: '100vh',
    width: '100vw',
    backgroundColor: '#BF9D9D',
}));

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
    '&:focus': {
        outline: 'none',
        boxShadow: 'none',
    },
    '&:active': {
        outline: 'none',
        boxShadow: 'none',
    },
}));

const StyledTextField = styled(TextField)(({ theme }) => ({
    '& .MuiOutlinedInput-root': {
        backgroundColor: '#F5F5F5',
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
        padding: '10px 14px',
    },
    width: '200px',
}));

const WaitingRoom = () => {
    const { matchId } = useParams();
    const { user } = useAuth();
    const navigate = useNavigate();
    const [players, setPlayers] = useState([{ username: user.username, ready: false }, { username: '', ready: false }]);
    const [isReady, setIsReady] = useState(false);
    const [gameConfig, setGameConfig] = useState({ boardWidth: 10, boardHeight: 10, winLength: 5 });
    const { sendMessage, lastMessage } = useWebSocket(matchId);

    useEffect(() => {
        if (lastMessage) {
            const data = JSON.parse(lastMessage.data);
            if (data.type === 'gameInfo' || data.type === 'playerJoined' || data.type === 'playerLeft' || data.type === 'playerReady') {
                setPlayers(data.players);
                if (data.type === 'gameInfo') {
                    setGameConfig(data.config);
                }
            } else if (data.type === 'gameStart') {
                navigate(`/match/${matchId}`);
            }
        }
    }, [lastMessage, matchId, navigate]);

    useEffect(() => {
        sendMessage({ type: 'joinMatch', matchId, token: user.token });
    }, [matchId, sendMessage, user.token]);

    useEffect(() => {
        const handleBeforeUnload = () => {
            sendMessage({ type: 'leaveMatch', matchId, token: user.token });
        };
        window.addEventListener('beforeunload', handleBeforeUnload);
        return () => {
            window.removeEventListener('beforeunload', handleBeforeUnload);
            handleBeforeUnload();
        };
    }, [matchId, sendMessage, user.token]);

    const handleReady = () => {
        const newIsReady = !isReady;
        setIsReady(newIsReady);
        const updatedPlayers = players.map(player =>
            player.username === user.username ? { ...player, ready: newIsReady } : player
        );
        setPlayers(updatedPlayers);
        sendMessage({ type: 'playerReady', matchId, token: user.token, ready: newIsReady });
    };

    const handleStart = () => {
        if (players.every(player => player.ready)) {
            sendMessage({ type: 'startGame', matchId, token: user.token, config: gameConfig });
        } else {
            // Show an error message that not all players are ready
        }
    };

    const handleExit = () => {
        navigate('/');
    };

    const handleCopyMatchId = () => {
        navigator.clipboard.writeText(matchId);
    };

    const handleConfigChange = (e) => {
        setGameConfig({ ...gameConfig, [e.target.name]: parseInt(e.target.value) });
        sendMessage({ type: 'updateConfig', matchId, token: user.token, config: { ...gameConfig, [e.target.name]: parseInt(e.target.value) } });
    };

    return (
        <StyledBox>
            <Typography variant="h4" gutterBottom align="center" sx={{ fontFamily: 'Explora, cursive', fontSize: '6rem', color: '#1E1E1E', marginBottom: '1rem' }}>
                LineUp
            </Typography>
            <Box sx={{ display: 'flex', alignItems: 'center', marginBottom: '2rem' }}>
                <Typography variant="body1" sx={{ fontWeight: 600, marginRight: '0.5rem' }}>
                    Match ID: {matchId}
                </Typography>
                <IconButton
                    onClick={handleCopyMatchId}
                    sx={{
                        '&:focus': {
                            outline: 'none',
                        },
                    }}
                >
                    <ContentCopyIcon />
                </IconButton>
            </Box>
            <Typography variant="h5" sx={{ fontWeight: 600, marginBottom: '1rem' }}>
                Board Settings
            </Typography>
            <Box sx={{ display: 'flex', alignItems: 'center', marginBottom: '1rem' }}>
                <Typography variant="body1" sx={{ width: '150px', textAlign: 'right', marginRight: '1rem' }}>Width:</Typography>
                <StyledTextField
                    type="number"
                    name="boardWidth"
                    value={gameConfig.boardWidth}
                    onChange={handleConfigChange}
                    inputProps={{ min: 3, max: 99 }}
                />
            </Box>
            <Box sx={{ display: 'flex', alignItems: 'center', marginBottom: '1rem' }}>
                <Typography variant="body1" sx={{ width: '150px', textAlign: 'right', marginRight: '1rem' }}>Height:</Typography>
                <StyledTextField
                    type="number"
                    name="boardHeight"
                    value={gameConfig.boardHeight}
                    onChange={handleConfigChange}
                    inputProps={{ min: 3, max: 99 }}
                />
            </Box>
            <Box sx={{ display: 'flex', alignItems: 'center', marginBottom: '1rem' }}>
                <Typography variant="body1" sx={{ width: '150px', textAlign: 'right', marginRight: '1rem' }}>Length To Win:</Typography>
                <StyledTextField
                    type="number"
                    name="winLength"
                    value={gameConfig.winLength}
                    onChange={handleConfigChange}
                    inputProps={{ min: 3, max: 19 }}
                />
            </Box>
            <Box sx={{ display: 'flex', alignItems: 'center', marginBottom: '1rem' }}>
                <EmojiEventsIcon sx={{ marginRight: '0.5rem' }} />
                <Typography variant="body1" sx={{ fontWeight: 600, marginRight: '0.5rem', width: '150px' }}>
                    Player1: {players[0].username}
                </Typography>
                {players[0].ready ? <CheckBoxIcon /> : <CheckBoxOutlineBlankIcon />}
            </Box>
            <Box sx={{ display: 'flex', alignItems: 'center', marginBottom: '2rem' }}>
                <Typography variant="body1" sx={{ fontWeight: 600, marginRight: '0.5rem', marginLeft: '1.5rem', width: '150px' }}>
                    Player2: {players[1].username || 'Waiting...'}
                </Typography>
                {players[1].ready ? <CheckBoxIcon /> : <CheckBoxOutlineBlankIcon />}
            </Box>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', width: '100%', position: 'absolute', bottom: '5%', padding: '0 25%' }}>
                <StyledButton onClick={handleExit} sx={{ color: '#B32D2D' }}>Exit</StyledButton>
                <StyledButton onClick={handleReady}>{isReady ? 'Cancel' : 'Ready'}</StyledButton>
                <StyledButton onClick={handleStart}>Start</StyledButton>
            </Box>
        </StyledBox>
    );
};

export default WaitingRoom;