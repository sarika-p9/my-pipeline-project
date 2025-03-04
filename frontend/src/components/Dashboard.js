import React from "react";
import { Box, Drawer, Card, Typography } from "@mui/material";
import { useNavigate } from "react-router-dom";

const Dashboard = () => {
  const navigate = useNavigate();

  return (
    <Box sx={{ display: "flex", height: "100vh" }}>
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
        <Box sx={{ display: "flex", flexDirection: "column", gap: 1.2, mt: 2, width: "90%" }}>
          {[
            { label: "User Details", color: "#082B6F", route: "/user-details" },
            { label: "Create Pipeline", color: "#0D47A1", route: "/create-pipeline" },
            { label: "Start Execution", color: "#1565C0", route: "/execute-pipeline" },
            { label: "Get Status", color: "#1E88E5", route: "/get-status" },
            { label: "Cancel Execution", color: "#42A5F5", route: "/cancel-pipeline" },
          ].map((item, index) => (
            <Card
              key={index}
              sx={{
                backgroundColor: item.color,
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
              <Typography variant="subtitle2" sx={{ fontSize: "14px", textAlign: "center", width: "100%" }}>
                {item.label}
              </Typography>
            </Card>
          ))}
        </Box>

        <Box sx={{ p: 1, mb: 2, width: "90%" }}>
          <Card
            sx={{
              backgroundColor: "#8E24AA",
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
            onClick={() => navigate("/login")}
          >
            <Typography variant="subtitle2" sx={{ fontSize: "14px", textAlign: "center", width: "100%" }}>
              Logout
            </Typography>
          </Card>
        </Box>
      </Drawer>
    </Box>
  );
};

export default Dashboard;
