import React, { useState } from "react";
import { Box, Button, TextField, Typography } from "@mui/material";
import apiService from "../services/apiService";

const ExecutePipeline = () => {
  const [pipelineId, setPipelineId] = useState("");

  const handleSubmit = async () => {
    const response = await apiService.executePipeline(pipelineId);
    alert(response.message);
  };

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h5">Execute Pipeline</Typography>
      <TextField
        label="Pipeline ID"
        fullWidth
        value={pipelineId}
        onChange={(e) => setPipelineId(e.target.value)}
        sx={{ my: 2 }}
      />
      <Button variant="contained" color="primary" onClick={handleSubmit}>
        Execute
      </Button>
    </Box>
  );
};

export default ExecutePipeline;
