import React, { useState } from "react";
import { TextField, Button, Typography, Container, Box, Link } from "@mui/material";
import { useNavigate } from "react-router-dom";
import axios from "axios";

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
    try {
      const response = await axios.post("http://localhost:8080/login", { email, password });
      const data = response.data;
      localStorage.setItem("token", data.token); // Store token
      localStorage.setItem("user_id", data.user_id); // Store user ID
      setMessage("Login successful!");
      setTimeout(() => navigate("/dashboard"), 1000); // Redirect to Dashboard after 1s
    } catch (error) {
      setError("Login failed. Please check your credentials.");
    }
  };

  return (
    <Container maxWidth="sm">
      <Box sx={{ textAlign: "center", mt: 5, p: 3, boxShadow: 3, borderRadius: 2 }}>
        <Typography variant="h4" gutterBottom>Login</Typography>
        <form onSubmit={handleLogin}>
          <TextField label="Email" fullWidth margin="normal" value={email} onChange={(e) => setEmail(e.target.value)} required />
          <TextField label="Password" type="password" fullWidth margin="normal" value={password} onChange={(e) => setPassword(e.target.value)} required />
          <Button variant="contained" color="primary" fullWidth type="submit" sx={{ mt: 2 }}>Login</Button>
        </form>
        {message && <Typography sx={{ mt: 2, color: "green" }}>{message}</Typography>}
        {error && <Typography sx={{ mt: 2, color: "red" }}>{error}</Typography>}
        <Typography sx={{ mt: 2 }}>
          Don't have an account? <Link onClick={() => navigate("/register")} sx={{ cursor: "pointer" }}>Register here</Link>
        </Typography>
      </Box>
    </Container>
  );
};

export default Login;