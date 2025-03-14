import React, { useState, useEffect, useMemo } from "react";
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
import Topbar from "../components/Topbar";

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
  const [pipelineStages, setPipelineStages] = useState([]);
  const [isParallel, setIsParallel] = useState(true);
  const navigate = useNavigate();
  const [selectedPipelineStages, setSelectedPipelineStages] = useState([]);
  const [selectedPipelineId, setSelectedPipelineId] = useState(null);
  const [openStageModal, setOpenStageModal] = useState(false);
  const user_id = getUserIdFromToken();
  const [loading, setLoading] = useState(false);
  const [pipelineName, setPipelineName] = useState("");
  const [numStages, setNumStages] = useState(1);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [stageNames, setStageNames] = useState([]);
  const [stages, setStages] = useState([]);
  const [openDialog, setOpenDialog] = useState(false);
  const [pipelineStatus, setPipelineStatus] = useState({});
  const [socket, setSocket] = useState(null);


  
  useEffect(() => {
    if (isTokenExpired()) {
      console.warn("Token expired. Logging out...");
      localStorage.clear();
      navigate("/login");
      return;
    }
    fetchUserPipelines();
  }, []);


  useEffect(() => {
    const ws = new WebSocket("ws://127.0.0.1:8080/ws");

    ws.onopen = () => console.log("✅ WebSocket Connected");

    ws.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            console.log("🔄 WebSocket Data Received:", data);

            // ✅ Ensure we only update stages for the selected pipeline
            if (data.pipelineId === selectedPipelineId) {
                setSelectedPipelineStages((prevStages) =>
                    prevStages.map((stage) =>
                        stage.StageID === data.stageId
                            ? { ...stage, Status: data.status }
                            : stage
                    )
                );

                // ✅ Check if all stages are completed, then update the pipeline status
                setTimeout(() => {
                    setSelectedPipelineStages((prevStages) => {
                        const allCompleted = prevStages.every(stage => stage.Status === "Completed");
                        if (allCompleted) {
                            fetchAndUpdatePipelineStatus();
                        }
                        return prevStages;
                    });
                }, 500);
            }
        } catch (error) {
            console.error("❌ WebSocket Error:", event.data, error);
        }
    };

    ws.onerror = (error) => console.error("❌ WebSocket Error:", error);
    ws.onclose = () => console.log("🔴 WebSocket Disconnected");

    setSocket(ws);
    return () => ws.close();
}, [selectedPipelineId]);  // ✅ Add dependency to ensure it updates correctly



  const setupWebSocket = (pipelineId) => {
    const ws = new WebSocket("ws://127.0.0.1:8080/ws");

    ws.onopen = () => console.log("✅ WebSocket Connected");

    ws.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            console.log("🔄 WebSocket Data Received:", data);

            setSelectedPipelineStages((prevStages) => 
                prevStages.map((stage) => 
                    stage.StageName === data.stage_name ? { ...stage, Status: data.status } : stage
                )
            );

            // ✅ **Check if All Stages are Completed, Then Update Pipeline**
            setTimeout(() => {
                setSelectedPipelineStages((prevStages) => {
                    const allCompleted = prevStages.every(stage => stage.Status === "Completed");
                    if (allCompleted) {
                        fetchAndUpdatePipelineStatus(); // ✅ Update pipeline status
                    }
                    return prevStages;
                });
            }, 500);

        } catch (error) {
            console.error("❌ WebSocket Error:", event.data, error);
        }
    };

    ws.onerror = (error) => console.error("❌ WebSocket Error:", error);
    ws.onclose = () => console.log("🔴 WebSocket Disconnected");
};




const fetchAndUpdatePipelineStatus = async () => {
  try {
      const response = await authAxios.get(`/pipelines?user_id=${user_id}`);
      if (Array.isArray(response.data)) {
          setPipelines(response.data);
      }
  } catch (error) {
      console.error("❌ Failed to update pipeline status:", error);
  }
};














  

  const authAxios = useMemo(() => {
    return axios.create({
      baseURL: "http://localhost:8080",
      headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
    });
  }, []);


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
    console.log("numStages Value:", numStages);
    console.log("Type of numStages:", typeof numStages);
  
    const stageCount = parseInt(numStages, 10); // Explicitly convert
  
    if (!stageCount || stageCount <= 0) {  
      alert("Invalid number of stages! Please enter a valid number.");
      return;
    }
  
    const payload = {
      name: pipelineName,
      stages: stageCount,  // ✅ Ensure it's included correctly
      is_parallel: isParallel ?? true, 
      user_id: getUserIdFromToken() || "default-user-id", 
      stage_names: stageNames || [], 
    };
  
    console.log("🚀 Sending Payload:", JSON.stringify(payload, null, 2)); // Debugging log
  
    try {
      const response = await authAxios.post("/createpipelines", payload);
      console.log("✅ Pipeline Created:", response.data);
  
      alert(`Pipeline Created: ${response.data.pipeline_id}`);
      fetchUserPipelines();
    } catch (error) {
      console.error("❌ Failed to create pipeline", error.response?.data || error);
      alert(`Error: ${error.response?.data?.error || "Unknown error"}`);
    }
  };

  const handlePipelineAction = async (pipelineId, status) => {
    try {
        console.log("🚀 Starting pipeline:", pipelineId, "Current Status:", status);

        const payload = {
            user_id: user_id,
            input: { raw_material: "Steel", quantity: 100 },
            is_parallel: isParallel,
        };

        console.log("📤 Sending Request:", JSON.stringify(payload, null, 2));

        if (status === "Running") {
            await authAxios.post(`/pipelines/${pipelineId}/cancel`, payload);
        } else if (status === "Completed") {
            alert("Completed pipelines cannot be started again.");
            return;
        } else {
            const response = await authAxios.post(`/pipelines/${pipelineId}/start`, payload);
            console.log("✅ Response:", response.data);
        }

        setPipelines((prevPipelines) =>
            prevPipelines.map((pipeline) =>
                pipeline.PipelineID === pipelineId ? { ...pipeline, Status: "Running" } : pipeline
            )
        );

        await fetchPipelineStages(pipelineId);

        if (socket && socket.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify({ pipelineId, action: "start" }));
        }

        setTimeout(fetchUserPipelines, 1000);
    } catch (error) {
        console.error("❌ Failed to update pipeline status:", error.response?.data || error);
    }
};


const handleShowStages = async (pipelineID) => {
  await fetchPipelineStages(pipelineID);
  setOpenStageModal(true);
};


  
  

  const handleSaveStageNames = () => {
    setPipelineStages(stageNames.map((name, index) => ({
      name,
      stage_number: index + 1
    }))); 
    setDialogOpen(false); 
  };
  


  const fetchPipelineStages = async (pipelineId) => {
    try {
        console.log(`📥 Fetching stages for pipeline: ${pipelineId}`);

        const response = await authAxios.get(`/pipelines/${pipelineId}/stages`);

        if (Array.isArray(response.data)) {
            const initializedStages = response.data.map(stage => ({
                StageID: stage.StageID,
                StageName: stage.StageName || "Unknown",
                Status: stage.Status || "Pending",
            }));

            setSelectedPipelineStages(initializedStages);
            setSelectedPipelineId(pipelineId);
            setOpenStageModal(true);

            // ✅ Start WebSocket tracking for this pipeline
            if (socket && socket.readyState === WebSocket.OPEN) {
                socket.send(JSON.stringify({ pipelineId, action: "track" }));
            }

            // ✅ **Poll for real-time updates**
            const intervalId = setInterval(async () => {
                const updatedResponse = await authAxios.get(`/pipelines/${pipelineId}/stages`);
                if (Array.isArray(updatedResponse.data)) {
                    setSelectedPipelineStages(updatedResponse.data);
                }
            }, 5000);  // ✅ Poll every 5 seconds

            return () => clearInterval(intervalId);
        } else {
            console.error("❌ Unexpected response format:", response.data);
        }
    } catch (error) {
        console.error("❌ Failed to fetch pipeline stages:", error);
    }
};




    

  const handleStageNameChange = (index, value) => {
    setStageNames((prev) => {
      const updatedStages = [...prev];
      updatedStages[index] = value;
      return updatedStages;
    });
  };

  const handleStageDialogOpen = () => {
    setStageNames(new Array(numStages).fill("")); // Ensures correct number of input fields
    setDialogOpen(true);
  };
  
  const handleDeletePipeline = async (pipelineId) => {
    if (!window.confirm("Are you sure you want to delete this pipeline?")) return;
  
    try {
      await axios.delete(`http://localhost:8080/api/pipelines/${pipelineId}`);
      
      setPipelines(pipelines.filter(pipeline => pipeline.PipelineID !== pipelineId));
    } catch (error) {
      console.error("Error deleting pipeline:", error);
      alert("Failed to delete pipeline.");
    }
  };
  
  if (loading) return <CircularProgress />;
  
  return (
    <Box sx={{ display: "flex" }}>
      <Topbar />
      <Sidebar />
      <Container maxWidth="md">
   
        <Box sx={{ mt: 5, p: 3, boxShadow: 3, borderRadius: 2 }}>
          <Typography variant="h5" sx={{ mb: 2 }}>Your Pipelines</Typography>
          {pipelines.length > 0 ? (
            <Table sx={{ minWidth: 500, border: "1px solid #ccc", borderRadius: "8px" }}>
              <TableHead>
                <TableRow sx={{ backgroundColor: "#f5f5f5" }}>
                  <TableCell><strong>Pipeline Name</strong></TableCell>
                  <TableCell><strong>Pipeline ID</strong></TableCell>
                  <TableCell><strong>Status</strong></TableCell>
                  <TableCell><strong>Actions</strong></TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {pipelines.map((pipeline) => (
                  <TableRow key={pipeline.PipelineID}>
                    <TableCell>{pipeline.PipelineName}</TableCell>
                    <TableCell>{pipeline.PipelineID}</TableCell>
                    <TableCell>
                    <Typography
  sx={{
    fontWeight: "bold",
    color: pipelineStatus[pipeline.PipelineID] === "Running" ? "green" : "gray",
  }}
>
  {pipelineStatus[pipeline.PipelineID] || pipeline.Status}
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
                          console.log("🔎 Show Stages clicked for pipeline:", pipeline.PipelineID);
                          fetchPipelineStages(pipeline.PipelineID);
                          setOpenStageModal(true);  // ✅ Ensure modal opens every time
                      }}
                  >
                        Show Stages
                      </Button>
                      )}
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
            <TableCell><strong>Stage Name</strong></TableCell>
            <TableCell><strong>Status</strong></TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
        {selectedPipelineStages.map((stage, index) => (
    <TableRow key={stage.StageID || `stage-${index}`}>


              <TableCell>{stage.StageID}</TableCell>
              <TableCell>{stage.StageName}</TableCell>
              <TableCell>
                <Typography sx={{ fontWeight: "bold", color: 
                  stage.Status === "Running" ? "blue" :
                  stage.Status === "Completed" ? "green" : 
                  "gray"
                }}>
                  {stage.Status}
                </Typography>
              </TableCell>
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



  
        <Box sx={{ mt: 5, p: 3, boxShadow: 3, borderRadius: 2 }}>
          <Typography variant="h5" sx={{ mb: 2 }}>Create New Pipeline</Typography>
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
          <Box sx={{ display: "flex", alignItems: "center", gap: 3, mt: 3 }}>
            <Typography variant="h6" sx={{ minWidth: 150, textAlign: "right" }}>Number of Stages:</Typography>
            <TextField
              select
              value={numStages}
              onChange={(e) => setNumStages(Number(e.target.value))}
              sx={{ width: 80 }}
            >
              {[...Array(10).keys()].map((num) => (
                <MenuItem key={num + 1} value={num + 1}>{num + 1}</MenuItem>
              ))}
            </TextField>
            <Button variant="contained" onClick={handleStageDialogOpen}>Enter Stage Names</Button>
          </Box>
          <Box sx={{ display: "flex", justifyContent: "center", mt: 3 }}>
            <Button
              variant="contained"
              color="secondary"
              sx={{ px: 3, py: 1.2 }}
              onClick={handleCreatePipeline}
              disabled={!pipelineName.trim() || loading || pipelineStages.some(stage => !stage.name.trim())}
            >
              {loading ? "Creating..." : "Create Pipeline"}
            </Button>
          </Box>
        </Box>
  
        <Dialog open={dialogOpen} onClose={() => setDialogOpen(false)}>
          <DialogTitle>Enter Stage Names</DialogTitle>
          <DialogContent>
            {stageNames.map((stage, index) => (
              <TextField
                key={index}
                label={`Stage ${index + 1} Name`}
                variant="outlined"
                value={stageNames[index] || ""}
                onChange={(e) => handleStageNameChange(index, e.target.value)}
                fullWidth
                sx={{ mt: 2 }}
              />
            ))}
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setDialogOpen(false)}>Close</Button>
            <Button variant="contained" color="primary" onClick={handleSaveStageNames}>
              Save
            </Button>
          </DialogActions>
        </Dialog>
      </Container>
    </Box>
  );
  
};

export default CreatePipeline;




