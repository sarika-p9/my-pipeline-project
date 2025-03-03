import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { registerUser } from "../services/api"; // Import API function
import { Button, TextField, Typography, Container, Box, Alert } from "@mui/material";

const Register = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();

  const handleRegister = async (e) => {
    e.preventDefault();
    setMessage("");
    setError("");

    const result = await registerUser(email, password);
    if (result.error) {
      setError(result.error);
    } else {
      setMessage("Registration successful! Check your email to verify.");
      setTimeout(() => navigate("/login"), 2000); // Redirect after 2s
    }
  };

  return (
    <Container maxWidth="sm">
      <Box sx={{ mt: 5, p: 4, boxShadow: 3, borderRadius: 2, bgcolor: "white" }}>
        <Typography variant="h4" textAlign="center" mb={2}>Register</Typography>
        {message && <Alert severity="success">{message}</Alert>}
        {error && <Alert severity="error">{error}</Alert>}
        <form onSubmit={handleRegister}>
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
            Register
          </Button>
        </form>
        <Typography textAlign="center" mt={2}>
          Already registered? <Button onClick={() => navigate("/login")}>Login</Button>
        </Typography>
      </Box>
    </Container>
  );
};

export default Register;
