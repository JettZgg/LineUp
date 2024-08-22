// src/components/Game/JoinMatch.jsx
import React, { useState } from 'react';
import { joinMatch } from '../../services/api';
import { useNavigate } from 'react-router-dom';

const JoinMatch = () => {
    const [matchId, setMatchId] = useState('');
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            console.log('Attempting to join match:', matchId);
            const response = await joinMatch(matchId);
            console.log('Join match response:', response);
            if (response && response.data) {
                navigate(`/match/${matchId}/waiting`);
            } else {
                console.error('Invalid response from server:', response);
            }
        } catch (error) {
            console.error('Failed to join match:', error);
            if (error.response) {
                console.error('Error response:', error.response.data);
            }
        }
    };

    return (
        <form onSubmit={handleSubmit}>
            <input
                type="text"
                value={matchId}
                onChange={(e) => setMatchId(e.target.value)}
                placeholder="Match ID"
                required
            />
            <button type="submit">Join Match</button>
        </form>
    );
};

export default JoinMatch;