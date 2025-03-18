import React, { useState, useEffect } from "react";
import { Box, Button, Typography, Dialog, DialogActions, DialogContent, DialogTitle, TextField, Select, MenuItem, FormControl, InputLabel } from "@mui/material";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import Topbar from "../components/Topbar";
import Sidebar from "./Sidebar";

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

const getUserIdFromToken = () => {
  const token = localStorage.getItem("token");
  if (!token) return null;
  try {
    const payload = JSON.parse(atob(token.split(".")[1]));
    return payload.sub;
  } catch {
    return null;
  }
};

const UserProfile = () => {
  const [user, setUser] = useState({ name: "", role: "", email: "" });
  const [editUser, setEditUser] = useState({ name: "", role: "" });
  const [profileOpen, setProfileOpen] = useState(false);
  const navigate = useNavigate();
  const user_id = getUserIdFromToken();

  useEffect(() => {
    if (isTokenExpired()) {
      console.warn("Token expired. Logging out...");
      setTimeout(() => {
        localStorage.clear();
        navigate("/login");
      }, 500);
      return;
    }
    fetchUserProfile();
  }, []);

  const authAxios = axios.create({
    baseURL: "http://localhost:30002",
    headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
  });

  const logoutUser = () => {
    localStorage.clear();
    navigate("/login");
  };

  const fetchUserProfile = async () => {
    if (!user_id) {
      console.error("User ID not found in token.");
      logoutUser();
      return;
    }

    try {
      const response = await authAxios.get(`/user/${user_id}`);
      if (response.data) {
        setUser(response.data);
        setEditUser(response.data);
      }
    } catch (error) {
      console.error("Failed to fetch user profile", error);
      logoutUser();
    }
  };

  const handleProfileSave = async () => {
    if (!user_id) {
      console.error("User ID not found. Cannot update profile.");
      return;
    }

    try {
      console.log("Updating profile for user_id:", user_id);
      await authAxios.put(`/user/${user_id}`, {
        name: editUser.name,
        role: editUser.role,
      });

      console.log("Profile updated successfully.");
      setUser(editUser);
      setProfileOpen(false);
    } catch (error) {
      console.error("Failed to update profile", error);
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
  
       <Topbar />
       <Sidebar />  
      
       <Box
    sx={{
      backgroundColor: "#FFFFFF",
      padding: 4,
      borderRadius: 2,
      boxShadow: "0px 4px 10px rgba(0, 0, 0, 0.1)",
      textAlign: "center",
      minWidth: "300px",
      maxWidth: "80vw",
      paddingBottom: "50px", 
    }}
  >
        <Typography variant="h4" gutterBottom>
          User Profile
        </Typography>
        <Typography variant="body1">Name: {user.name}</Typography>
        <Typography variant="body1">Role: {user.role}</Typography>
        <Typography variant="body1">Email: {user.email}</Typography>

        <Button variant="contained" color="primary" onClick={() => setProfileOpen(true)} sx={{ mt: 2 }}>
          Edit Profile
        </Button>
      </Box>

      <Dialog open={profileOpen} onClose={() => setProfileOpen(false)}>
        <DialogTitle>Edit Profile</DialogTitle>
        <DialogContent>
          <TextField
            label="Name"
            fullWidth
            margin="normal"
            value={editUser.name}
            onChange={(e) => setEditUser({ ...editUser, name: e.target.value })}
          />
          <FormControl fullWidth margin="normal">
            <InputLabel shrink={true} id="role-label">
              Role
            </InputLabel>
            <Select
              labelId="role-label"
              value={editUser.role}
              onChange={(e) => setEditUser({ ...editUser, role: e.target.value })}
              displayEmpty
            >
              <MenuItem value="" disabled>Select Role</MenuItem>
              <MenuItem value="super_admin">Super Admin</MenuItem>
              <MenuItem value="admin">Admin</MenuItem>
              <MenuItem value="manager">Manager</MenuItem>
              <MenuItem value="worker">Worker</MenuItem>
            </Select>
          </FormControl>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setProfileOpen(false)}>Cancel</Button>
          <Button variant="contained" color="primary" onClick={handleProfileSave}>
            Save Changes
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default UserProfile;
