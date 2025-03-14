import React from "react";
import { AppBar, Toolbar, Typography } from "@mui/material";

const Topbar = () => {
  return (
    <AppBar position="fixed" sx={{ backgroundColor: "#082B6F", zIndex: 1300 }}>
      <Toolbar>
        <Typography variant="h6" sx={{ flexGrow: 1, color: "white", textAlign: "center", fontWeight: "bold" }}>
          Distributed Manufacturing Pipeline Simulation System
        </Typography>
      </Toolbar>
    </AppBar>
  );
};

export default Topbar;
