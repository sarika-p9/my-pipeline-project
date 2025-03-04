import React, { useState } from "react";
import { Box, Button, TextField, Typography } from "@mui/material";
import apiService from "../services/apiService";

const CancelPipeline = () => {
  const [pipelineId, setPipelineId] = useState("");

  const handleSubmit = async () => {
    const response = await apiService.cancelPipeline(pipelineId);
    alert(response.message);
  };

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h5">Cancel Pipeline</Typography>
      <TextField
        label="Pipeline ID"
        fullWidth
        value={pipelineId}
        onChange={(e) => setPipelineId(e.target.value)}
        sx={{ my: 2 }}
      />
      <Button variant="contained" color="primary" onClick={handleSubmit}>
        Cancel
      </Button>
    </Box>
  );
};

export default CancelPipeline;
