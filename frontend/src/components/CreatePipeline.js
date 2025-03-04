import React, { useState } from "react";
import { Box, Button, TextField, Typography } from "@mui/material";
import apiService from "../services/apiService";

const CreatePipeline = () => {
  const [stages, setStages] = useState("");

  const handleSubmit = async () => {
    const response = await apiService.createPipeline(stages);
    alert(response.message);
  };

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h5">Create Pipeline</Typography>
      <TextField
        label="Number of Stages"
        type="number"
        fullWidth
        value={stages}
        onChange={(e) => setStages(e.target.value)}
        sx={{ my: 2 }}
      />
      <Button variant="contained" color="primary" onClick={handleSubmit}>
        Create
      </Button>
    </Box>
  );
};

export default CreatePipeline;
