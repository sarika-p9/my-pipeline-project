import React from "react";
import { useNavigate } from "react-router-dom";
import { Card, CardContent, Typography } from "@mui/material";

const Dashboard = () => {
  const navigate = useNavigate();

  const handleNavigation = (path) => {
    navigate(path);
  };

  return (
    <div style={{ display: "grid", gap: "20px", padding: "20px" }}>
      <Card onClick={() => handleNavigation("/user-details")} style={cardStyle}>
        <CardContent>
          <Typography variant="h5">User Details</Typography>
        </CardContent>
      </Card>

      <Card onClick={() => handleNavigation("/create-pipeline")} style={cardStyle}>
        <CardContent>
          <Typography variant="h5">Create Pipeline</Typography>
        </CardContent>
      </Card>

      <Card onClick={() => handleNavigation("/execute-pipeline")} style={cardStyle}>
        <CardContent>
          <Typography variant="h5">Start/Execute Pipeline</Typography>
        </CardContent>
      </Card>

      <Card onClick={() => handleNavigation("/get-status")} style={cardStyle}>
        <CardContent>
          <Typography variant="h5">Get Pipeline Status</Typography>
        </CardContent>
      </Card>

      <Card onClick={() => handleNavigation("/cancel-pipeline")} style={cardStyle}>
        <CardContent>
          <Typography variant="h5">Cancel Pipeline</Typography>
        </CardContent>
      </Card>
    </div>
  );
};

const cardStyle = {
  cursor: "pointer",
  padding: "20px",
  textAlign: "center",
  background: "lightblue",
};

export default Dashboard;
