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
            await joinMatch(matchId);
            navigate(`/match/${matchId}`);
        } catch (error) {
            console.error('Failed to join match:', error);
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