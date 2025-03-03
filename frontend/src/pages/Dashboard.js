import React from 'react';
import { Container, Typography, Button } from '@mui/material';
import { useNavigate } from 'react-router-dom';

const Dashboard = () => {
    const navigate = useNavigate();

    const handleLogout = () => {
        localStorage.removeItem('token');
        navigate('/login');
    };

    return (
        <Container maxWidth="md">
            <Typography variant="h4" mt={5}>Welcome to Dashboard</Typography>
            <Button variant="contained" color="secondary" onClick={handleLogout}>Logout</Button>
        </Container>
    );
};

export default Dashboard;
