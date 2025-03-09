import React, { useState, useEffect } from "react";
import { 
  AppBar, Toolbar, Typography, Button, Container, Box, Dialog, DialogTitle, DialogContent, DialogActions,
  Menu, MenuItem, ToggleButton, ToggleButtonGroup, IconButton, Table, TableHead, TableBody, TableRow, TableCell, TextField 
} from "@mui/material";
import AddIcon from "@mui/icons-material/Add";
import RemoveIcon from "@mui/icons-material/Remove";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import Sidebar from "../pages/Sidebar";
import CircularProgress from "@mui/material/CircularProgress";



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

const CreatePipeline = () => {
  const [user, setUser] = useState({ name: "", role: "", email: "" });
  const [pipelines, setPipelines] = useState([]);
  const [profileOpen, setProfileOpen] = useState(false);
  const [pipelineStages, setPipelineStages] = useState(1);
  const [isParallel, setIsParallel] = useState(false);
  const [anchorEl, setAnchorEl] = useState(null);
  const navigate = useNavigate();
  const [stagesDialogOpen, setStagesDialogOpen] = useState(false);
  const [selectedPipelineStages, setSelectedPipelineStages] = useState([]);
  const [selectedPipelineId, setSelectedPipelineId] = useState(null);
  const [openStageModal, setOpenStageModal] = useState(false);
  const user_id = getUserIdFromToken();
  const [loading, setLoading] = useState(false);
  const [pipelineName, setPipelineName] = useState("");


  
  useEffect(() => {
    if (isTokenExpired()) {
      console.warn("Token expired. Logging out...");
      localStorage.clear();
      navigate("/login");
      return;
    }
    fetchUserProfile();
    fetchUserPipelines();
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

  const fetchUserPipelines = async () => {
    try {
      const response = await authAxios.get(`/pipelines?user_id=${user_id}`);
      console.log("Fetched Pipelines:", response.data);
      if (Array.isArray(response.data)) {
        setPipelines(response.data);
      } else {
        console.error("Unexpected response format:", response.data);
      }
    } catch (error) {
      console.error("Failed to fetch pipelines", error);
    }
  };

  const handleCreatePipeline = async () => {
    try {
      await authAxios.post("/createpipelines", {
        name: pipelineName,  // ✅ Added Name
        stages: pipelineStages,
        is_parallel: isParallel,
        user_id: user_id,
      });
      fetchUserPipelines();
    } catch (error) {
      console.error("Failed to create pipeline", error);
    }
  };
  

  const handlePipelineAction = async (pipelineId, status) => {
    try {
      if (status === "Running") {
        await authAxios.post(`/pipelines/${pipelineId}/cancel`, {
          user_id: user_id,
          is_parallel: isParallel,
        });
      } else if (status === "Completed") {
        alert("Completed pipelines cannot be started again.");
        return;
      } else {
        await authAxios.post(`/pipelines/${pipelineId}/start`, {
          user_id: user_id,
          input: { raw_material: "Steel", quantity: 100 },
          is_parallel: isParallel,
        });
      }
  
      setPipelines((prevPipelines) =>
        prevPipelines.map((pipeline) =>
          pipeline.PipelineID === pipelineId ? { ...pipeline, Status: "Running" } : pipeline
        )
      );
      setTimeout(fetchUserPipelines, 1000); 
    } catch (error) {
      console.error("Failed to update pipeline status", error);
    }
  };


  const fetchPipelineStages = async (pipelineId) => {
    try {
      console.log(`Fetching stages for pipeline: ${pipelineId}`); 
  
      const response = await authAxios.get(`/pipelines/${pipelineId}/stages`);
      
      console.log("Stages Data:", response.data); 
  
      if (Array.isArray(response.data)) {
        setSelectedPipelineStages(response.data);
        setOpenStageModal(true); 
      } else {
        console.error("Unexpected response format:", response.data);
      }
    } catch (error) {
      console.error("Failed to fetch pipeline stages:", error);
      logoutUser();
    }
  };

  const handleProfileSave = async () => {
    try {
      await authAxios.put(`/user/${user_id}`, {
        name: user.name,
        role: user.role,
      });
      setProfileOpen(false);
    } catch (error) {
      console.error("Failed to update profile", error);
    }
  };
  
  const handleDeletePipeline = async (pipelineId) => {
    if (!window.confirm("Are you sure you want to delete this pipeline?")) return;
  
    try {
      await axios.delete(`http://localhost:8080/api/pipelines/${pipelineId}`);
      
      // ✅ Ensure correct key: Filter out deleted pipeline
      setPipelines(pipelines.filter(pipeline => pipeline.PipelineID !== pipelineId));
    } catch (error) {
      console.error("Error deleting pipeline:", error);
      alert("Failed to delete pipeline.");
    }
  };
  
  // Show loader while data is being fetched
  if (loading) return <CircularProgress />;
  


  return (
    <Box sx={{ display: "flex" }}> 
    <Sidebar />  
    <Container maxWidth="md">
    
      <Box sx={{ mt: 5, p: 3, boxShadow: 3, borderRadius: 2 }}>
        <Typography variant="h5" sx={{ mb: 2 }}>Your Pipelines</Typography>
        {pipelines.length > 0 ? (
          <Table sx={{ minWidth: 500, border: "1px solid #ccc", borderRadius: "8px" }}>
            <TableHead>
              <TableRow sx={{ backgroundColor: "#f5f5f5" }}>
                <TableCell><strong>Pipeline ID</strong></TableCell>
                <TableCell><strong>Pipeline Name</strong></TableCell>
                <TableCell><strong>Status</strong></TableCell>
                <TableCell><strong>Actions</strong></TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {pipelines.map((pipeline) => (
                <TableRow key={pipeline.PipelineID}>
                  <TableCell>{pipeline.PipelineID}</TableCell>
                  <TableCell>{pipeline.PipelineName}</TableCell>  
                  <TableCell>
                    <Typography sx={{ fontWeight: "bold", color: pipeline.Status === "Running" ? "green" : "gray" }}>
                      {pipeline.Status}
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Button
                      variant="contained"
                      color={pipeline.Status === "Running" ? "error" : "primary"}
                      onClick={() => handlePipelineAction(pipeline.PipelineID, pipeline.Status)}
                    >
                      {pipeline.Status === "Running" ? "Cancel Pipeline" : "Start Pipeline"}
                    </Button>
                    {pipeline.Status !== "Created" && (
                      <Button
                      variant="outlined"
                      sx={{ ml: 2 }}
                      onClick={() => {
                        console.log("Show Stages button clicked for pipeline:", pipeline.PipelineID); // ✅ Debug log
                        fetchPipelineStages(pipeline.PipelineID);
                      }}
                    >
                      Show Stages
                    </Button>
                    )}
                    {/* Delete Button (Only if pipeline is NOT running) */}
        {pipeline.Status !== "Running" && (
          <Button
            variant="contained"
            color="secondary"
            sx={{ ml: 2 }}
            onClick={() => handleDeletePipeline(pipeline.PipelineID)}
          >
            Delete
          </Button>
        )}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        ) : (
          <Typography>No pipelines created.</Typography>
        )}
      </Box>
      <Dialog open={openStageModal} onClose={() => setOpenStageModal(false)}>
  <DialogTitle>Pipeline Stages</DialogTitle>
  <DialogContent>
    {selectedPipelineStages.length > 0 ? (
      <Table>
        <TableHead>
          <TableRow>
            <TableCell><strong>Stage ID</strong></TableCell>
            <TableCell><strong>Status</strong></TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {selectedPipelineStages.map((stage) => (
            <TableRow key={stage.StageID}>
              <TableCell>{stage.StageID}</TableCell>
              <TableCell>{stage.Status}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    ) : (
      <Typography>No stages found for this pipeline.</Typography>
    )}
  </DialogContent>
  <DialogActions>
    <Button onClick={() => setOpenStageModal(false)}>Close</Button>
  </DialogActions>
</Dialog>


<Dialog open={profileOpen} onClose={() => setProfileOpen(false)}>
        <DialogTitle>Edit Profile</DialogTitle>
        <DialogContent>
          <TextField label="Name" fullWidth margin="normal" value={user.name} onChange={(e) => setUser({ ...user, name: e.target.value })} />
          <TextField label="Role" fullWidth margin="normal" value={user.role} onChange={(e) => setUser({ ...user, role: e.target.value })} />
          <TextField label="Email" fullWidth margin="normal" value={user.email} disabled />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setProfileOpen(false)}>Cancel</Button>
          <Button onClick={handleProfileSave} color="primary">Save</Button>
        </DialogActions>
      </Dialog>



      <Box sx={{ mt: 5, p: 3, boxShadow: 3, borderRadius: 2 }}>
      <Typography variant="h5" sx={{ mb: 2 }}>Create New Pipeline</Typography>

      {/* Pipeline Name Input */}
      <Box sx={{ display: "flex", alignItems: "center", gap: 3, mt: 2 }}>
        <Typography variant="h6" sx={{ minWidth: 150, textAlign: "right" }}>Pipeline Name:</Typography>
        <TextField
          variant="outlined"
          value={pipelineName}
          onChange={(e) => setPipelineName(e.target.value)}
          placeholder="Enter pipeline name"
          sx={{ flexGrow: 1 }}
        />
      </Box>

      {/* Dropdown for Number of Stages */}
      <Box sx={{ display: "flex", alignItems: "center", gap: 3, mt: 3 }}>
        <Typography variant="h6" sx={{ minWidth: 150, textAlign: "right" }}>Number of Stages:</Typography>
        <TextField
          select
          value={pipelineStages}
          onChange={(e) => setPipelineStages(e.target.value)}
          sx={{ width: 80 }}
        >
          {[...Array(10).keys()].map((num) => (
            <MenuItem key={num + 1} value={num + 1}>
              {num + 1}
            </MenuItem>
          ))}
        </TextField>
      </Box>

      {/* Create Pipeline Button */}
      <Box sx={{ display: "flex", justifyContent: "center", mt: 3 }}>
        <Button
          variant="contained"
          color="secondary"
          sx={{ px: 3, py: 1.2 }}
          onClick={handleCreatePipeline}
          disabled={!pipelineName.trim() || loading}
        >
          {loading ? "Creating..." : "Create Pipeline"}
        </Button>
      </Box>
    </Box>
    </Container>
    </Box>
  );
};

export default CreatePipeline;




