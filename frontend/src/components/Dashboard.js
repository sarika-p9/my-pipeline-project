import React from "react";
import { Container, Typography, Box, Button } from "@mui/material";
import { useNavigate } from "react-router-dom";

const Dashboard = () => {
  const navigate = useNavigate();

  const handleLogout = () => {
    navigate("/login");
  };

  return (
    <Container maxWidth="sm">
      <Box sx={{ textAlign: "center", mt: 5, p: 3, boxShadow: 3, borderRadius: 2 }}>
        <Typography variant="h4">Welcome!</Typography>
        <Button variant="contained" color="secondary" sx={{ mt: 3 }} onClick={handleLogout}>Logout</Button>
      </Box>
    </Container>
  );
};

export default Dashboard;
