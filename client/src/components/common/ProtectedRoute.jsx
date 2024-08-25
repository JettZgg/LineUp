// src/components/common/ProtectedRoute.jsx
import React, { useState, useEffect } from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from './AuthContext';

const ProtectedRoute = ({ children }) => {
    const { user } = useAuth();
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const checkUser = async () => {
            await new Promise(resolve => setTimeout(resolve, 100)); // Small delay to ensure AuthContext has initialized
            setLoading(false);
        };
        checkUser();
    }, []);

    if (loading) {
        return null; // or a loading spinner
    }

    if (!user) {
        return <Navigate to="/login" replace />;
    }

    return children;
};

export default ProtectedRoute;