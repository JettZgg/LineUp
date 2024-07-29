// src/components/Game/CreateMatch.jsx
import React, { useState } from 'react';
import { createMatch } from '../../services/api';
import { useNavigate } from 'react-router-dom';

const CreateMatch = () => {
    const [boardWidth, setBoardWidth] = useState(10);
    const [boardHeight, setBoardHeight] = useState(10);
    const [winLength, setWinLength] = useState(5);
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const response = await createMatch({ boardWidth, boardHeight, winLength });
            navigate(`/match/${response.data.id}`);
        } catch (error) {
            console.error('Failed to create match:', error);
        }
    };

    return (
        <form onSubmit={handleSubmit}>
            <input
                type="number"
                value={boardWidth}
                onChange={(e) => setBoardWidth(parseInt(e.target.value))}
                placeholder="Board Width"
                min="3"
                max="99"
                required
            />
            <input
                type="number"
                value={boardHeight}
                onChange={(e) => setBoardHeight(parseInt(e.target.value))}
                placeholder="Board Height"
                min="3"
                max="99"
                required
            />
            <input
                type="number"
                value={winLength}
                onChange={(e) => setWinLength(parseInt(e.target.value))}
                placeholder="Win Length"
                min="3"
                max="19"
                required
            />
            <button type="submit">Create Match</button>
        </form>
    );
};

export default CreateMatch;