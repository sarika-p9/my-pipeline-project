import React, { useState } from "react";
import { Box, Button, TextField, Typography } from "@mui/material";
import apiService from "../services/apiService";

const GetStatus = () => {
  const [pipelineId, setPipelineId] = useState("");

  const handleSubmit = async () => {
    const response = await apiService.getStatus(pipelineId);
    alert(response.message);
  };

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h5">Get Pipeline Status</Typography>
      <TextField
        label="Pipeline ID"
        fullWidth
        value={pipelineId}
        onChange={(e) => setPipelineId(e.target.value)}
        sx={{ my: 2 }}
      />
      <Button variant="contained" color="primary" onClick={handleSubmit}>
        Get Status
      </Button>
    </Box>
  );
};

export default GetStatus;
