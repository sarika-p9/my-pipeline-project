import React, { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route, Link, useNavigate } from "react-router-dom";
import { TextField, Button, Select, MenuItem, Typography, Container, Box, AppBar, Toolbar } from "@mui/material";
import axios from "axios";
import Dashboard from "./Dashboard";
import Topbar from "./Topbar";

const decodeToken = (token) => {
  try {
    const payload = JSON.parse(atob(token.split(".")[1]));
    return payload;
  } catch (error) {
    console.error("Failed to decode token:", error);
    return null;
  }
};

const isTokenExpired = (token) => {
  const payload = decodeToken(token);
  if (!payload || !payload.exp) return true;
  return Date.now() >= payload.exp * 1000;
};

const checkSession = (navigate) => {
  const token = localStorage.getItem("token");
  if (!token || isTokenExpired(token)) {
    console.warn("Token expired. Logging out...");
    localStorage.clear();
    navigate("/login");
  }
};


const AuthLayout = ({ children, title }) => {
  return (
    <Container sx={{maxWidth: "sm", paddingTop: 8}}>
      <Box sx={{ textAlign: "center", mt: 5, p: 4, boxShadow: 3, borderRadius: 2, paddingTop: 4 }}>
        <Typography variant="h4" sx={{ mb: 2 }}>{title}</Typography>
        {children}
      </Box>
    </Container>
  );
};

const RegisterPage = ({ apiType }) => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");
  const navigate = useNavigate();

  const handleRegister = async () => {
    setMessage("");
    if (apiType === "rest") {
      try {
        await axios.post("http://localhost:30002/register", { email, password });
        setMessage("Registration successful! Please check your email to verify.");
        window.open("https://mail.google.com", "_blank");
      } catch {
        setMessage("Registration failed. Please try again.");
      }
    } else {
      console.log(`Run this gRPC command manually: grpcurl -plaintext -d '{"email": "${email}", "password": "${password}"}' localhost:50051 auth.AuthService/Register`);
      setMessage("Open your email and click the link to authenticate.");
    }
  };

  return (
    <AuthLayout title="Register">
      <TextField label="Email" fullWidth margin="normal" value={email} onChange={(e) => setEmail(e.target.value)} />
      <TextField label="Password" type="password" fullWidth margin="normal" value={password} onChange={(e) => setPassword(e.target.value)} />
      <Button variant="contained" color="primary" fullWidth onClick={handleRegister} sx={{ mt: 2 }}>Register</Button>
      {message && <Typography sx={{ mt: 2, color: "green" }}>{message}</Typography>}
      <Typography sx={{ mt: 2 }}>Already registered? <Link to="/login">Login</Link></Typography>
    </AuthLayout>
  );
};

const LoginPage = ({ apiType }) => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    checkSession(navigate);  
  }, []);

  const handleLogin = async () => {
    setMessage("");
    if (apiType === "rest") {
      try {

        const response = await axios.post("http://localhost:30002/login", { email, password });
  
        const { token } = response.data;
        if (!token) throw new Error("Token not received");

        localStorage.setItem("token", token);

        const payload = JSON.parse(atob(token.split(".")[1]));
        localStorage.setItem("user_id", payload.sub);
  
        navigate("/dashboard");
        setTimeout(() => {
          window.location.reload();
        }, 100);
      } catch (error) {
        console.error("Login failed:", error);
        setMessage("Login failed. Please check your credentials.");
      }
    } else {
      console.log(`Run this gRPC command manually: grpcurl -plaintext -d '{"email": "${email}", "password": "${password}"}' localhost:50051 auth.AuthService/Login`);
      setMessage("Check console for gRPC login command.");
    }
  };
  

  return (
    <AuthLayout title="Login">
      <TextField label="Email" fullWidth margin="normal" value={email} onChange={(e) => setEmail(e.target.value)} />
      <TextField label="Password" type="password" fullWidth margin="normal" value={password} onChange={(e) => setPassword(e.target.value)} />
      <Button variant="contained" color="primary" fullWidth onClick={handleLogin} sx={{ mt: 2 }}>Login</Button>
      {message && <Typography sx={{ mt: 2, color: "red" }}>{message}</Typography>}
      <Typography sx={{ mt: 2 }}>Do not have an account? <Link to="/register">Register</Link></Typography>
    </AuthLayout>
  );
};

const App = () => {
  const [apiType, setApiType] = useState("rest");

  return (
    <Router>
      <Topbar />
      <Routes>
        <Route path="/register" element={<RegisterPage apiType={apiType} />} />
        <Route path="/login" element={<LoginPage apiType={apiType} />} />
        <Route path="/dashboard/*" element={<Dashboard />} />
        <Route path="/" element={<RegisterPage apiType={apiType} />} />
      </Routes>
    </Router>
  );
};

export default App;