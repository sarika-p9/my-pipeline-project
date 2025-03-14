import React, { useState, useEffect } from "react";
import { 
  AppBar, Toolbar, Typography, Container, Box} from "@mui/material";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import Sidebar from "../pages/Sidebar";
import Topbar from "./Topbar";


const isTokenExpired = () => {
  const token = localStorage.getItem("token");
  if (!token) return true;
  try {
    const payload = JSON.parse(atob(token.split(".")[1]));
    return Date.now() >= payload.exp * 1000;
  } catch {
    return true;
  }
};

const Dashboard = () => {
  const [user, setUser] = useState({ name: "", role: "", email: "" });
  const navigate = useNavigate();
  
  useEffect(() => {
    if (isTokenExpired()) {
      console.warn("Token expired. Logging out...");
      localStorage.clear();
      navigate("/login");
      return;
    }
    fetchUserProfile();
  }, []);

  const authAxios = axios.create({
    baseURL: "http://localhost:8080",
    headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
  });

  const logoutUser = () => {
    localStorage.clear();
    navigate("/login");
  };

  const fetchUserProfile = async () => {
    try {
      const response = await authAxios.get(`/user/${localStorage.getItem("user_id")}`);
      if (response.data) {
        setUser(response.data);
        localStorage.setItem("user_name", response.data.name);
        localStorage.setItem("user_role", response.data.role);
      }
    } catch (error) {
      console.error("Failed to fetch user profile", error);
      logoutUser();
    }
  };

  return (
    <Box
    sx={{
      width: "100vw",
      height: "100vh",
      display: "flex",
      justifyContent: "center",
      alignItems: "center",
      backgroundColor: "#42A5F5",
    }}
  >
    <Box sx={{ display: "flex", paddingTop: 70, paddingLeft: 3 }}>
      <Topbar />
      <Sidebar />
      <Container sx={{ maxWidth: "md" }}>
        <Box
          sx={{
            backgroundColor: "#FFFFFF", 
            textAlign: "center",
            mt: 5,
            p: 4,
            boxShadow: 3,
            borderRadius: 2,
          }}
        >
          <Typography variant="h4" sx={{ color: "#082B6F" }}> 
            Welcome, {user.name || "User"}
          </Typography>
        </Box>
      </Container>
    </Box>
  </Box>
  
  );
};

export default Dashboard;




