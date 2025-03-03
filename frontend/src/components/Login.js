import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { loginUser } from "../services/api"; // Import API function
import { Button, TextField, Typography, Container, Box, Alert } from "@mui/material";

const Login = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();

  const handleLogin = async (e) => {
    e.preventDefault();
    setMessage("");
    setError("");

    const result = await loginUser(email, password);
    if (result.error) {
      setError(result.error);
    } else {
      localStorage.setItem("token", result.token);
      setMessage("Login successful! Redirecting...");
      setTimeout(() => navigate("/dashboard"), 1500); // Redirect after 1.5s
    }
  };

  return (
    <Container maxWidth="sm">
      <Box sx={{ mt: 5, p: 4, boxShadow: 3, borderRadius: 2, bgcolor: "white" }}>
        <Typography variant="h4" textAlign="center" mb={2}>Login</Typography>
        {message && <Alert severity="success">{message}</Alert>}
        {error && <Alert severity="error">{error}</Alert>}
        <form onSubmit={handleLogin}>
          <TextField
            label="Email"
            fullWidth
            margin="normal"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
          <TextField
            label="Password"
            type="password"
            fullWidth
            margin="normal"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          <Button type="submit" variant="contained" color="primary" fullWidth sx={{ mt: 2 }}>
            Login
          </Button>
        </form>
        <Typography textAlign="center" mt={2}>
          Don't have an account? <Button onClick={() => navigate("/register")}>Register</Button>
        </Typography>
      </Box>
    </Container>
  );
};

export default Login;
