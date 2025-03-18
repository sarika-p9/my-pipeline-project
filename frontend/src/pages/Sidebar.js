import React from "react";
import { Drawer, Box, Card, Typography } from "@mui/material";
import { useNavigate } from "react-router-dom";
import LogoutIcon from "@mui/icons-material/Logout";

const handleLogout = async () => {
  const token = localStorage.getItem("token");
  if (!token) {
    console.error("No token found");
    return;
  }

  try {
    const response = await fetch("http://localhost:30002/logout", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ token }), 
    });

    if (response.ok) {
      localStorage.removeItem("token"); 
      alert("Logout successful");
      window.location.href = "/login"; 
    } else {
      console.error("Logout failed");
    }
  } catch (error) {
    console.error("Error:", error);
  }
};


const Sidebar = () => {
  const navigate = useNavigate();

  return (
    <Box sx={{ display: "flex", height: "100vh"}}>
      <Drawer
        variant="permanent"
        sx={{
          width: 240,
          flexShrink: 0,
          "& .MuiDrawer-paper": {
            width: 240,
            backgroundColor: "white",
            padding: 2,
            display: "flex",
            flexDirection: "column",
            justifyContent: "space-between",
            alignItems: "center",
          },
        }}
      >
        <Box sx={{ display: "flex", flexDirection: "column", gap: 1.2, mt: 2, width: "90%", paddingTop: 8 }}>
          {[
            { label: "User Profile", color: "#0D47A1", route: "/user-profile" },
            { label: "Pipelines", color: "#1565C0", route: "/create-pipeline" },
           // { label: "Start Execution", color: "#1565C0", route: "/execute-pipeline" },
          //  { label: "Get Status", color: "#1E88E5", route: "/get-status" },
          ].map((item, index) => (
            <Card
              key={index}
              sx={{
                backgroundColor: item.color,
                fontWeight: "bold",
                color: "white",
                p: 1,
                borderRadius: 3,
                boxShadow: 3,
                width: "100%",
                height: 45,
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                cursor: "pointer",
              }}
              onClick={() => navigate(item.route)}
            >
              <Typography variant="subtitle2" sx={{ fontWeight: "bold", fontSize: "14px", textAlign: "center", width: "100%" }}>
                {item.label}
              </Typography>
            </Card>
          ))}
        </Box>
        <Box sx={{ width: "90%", mb: 3, paddingBottom: 1 }}>
          <Card
            sx={{
              backgroundColor: "#1E88E5",
              color: "white",
              p: 1,
              borderRadius: 3,
              boxShadow: 3,
              width: "100%",
              height: 45,
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              cursor: "pointer",
              gap: 1
            }}
            onClick={handleLogout}
          >
            <LogoutIcon sx={{ fontSize: 18 }} />  {/* Logout Icon */}
    <Typography variant="subtitle2" sx={{ fontWeight: "bold", fontSize: "14px", textAlign: "center" }}>
      Logout
    </Typography>
          </Card>
        </Box>
      </Drawer>
    </Box>
  );
};

export default Sidebar;
