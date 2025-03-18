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
  const [selectedPipelineName, setSelectedPipelineName] = useState("");



  
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
    if (!selectedPipelineId) return;

    const ws = new WebSocket("ws://localhost:30002/ws");

    ws.onopen = () => console.log("✅ WebSocket Connected");

    ws.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            console.log("🔄 WebSocket Data Received:", data);

            // ✅ Ensure updates are only for the selected pipeline
            if (data.pipelineId !== selectedPipelineId) return;

            setSelectedPipelineStages((prevStages) =>
                prevStages.map((stage) =>
                    stage.StageID === data.stageId
                        ? { ...stage, Status: data.status }
                        : stage
                ).sort((a, b) => a.StageID.localeCompare(b.StageID))
            );

            setTimeout(async () => {
                const allCompleted = selectedPipelineStages.every(
                    (stage) => stage.Status === "Completed"
                );

                if (allCompleted) {
                    console.log("✅ All Stages Completed! Updating Pipeline Status...");
                    await fetchAndUpdatePipelineStatus();
                    setTimeout(() => {
                        setOpenStageModal(false); // ✅ Auto-close modal if needed
                    }, 1000);
                }
            }, 500);
        } catch (error) {
            console.error("❌ WebSocket Error:", event.data, error);
        }
    };

    ws.onerror = (error) => console.error("❌ WebSocket Error:", error);
    ws.onclose = () => console.log("🔴 WebSocket Disconnected");

    setSocket(ws);
    return () => ws.close();
}, [selectedPipelineId]);


const closeStageModal = () => {
  console.log("❌ Closing Stage Modal & Stopping Polling...");
  setOpenStageModal(false);
  setSelectedPipelineId(null);  // ✅ Prevents unnecessary polling
};




useEffect(() => {
  let intervalId;

  if (openStageModal && selectedPipelineId) {
      intervalId = setInterval(async () => {
          const updatedResponse = await authAxios.get(`/pipelines/${selectedPipelineId}/stages`);
          if (Array.isArray(updatedResponse.data)) {
              setSelectedPipelineStages(updatedResponse.data);

              const allCompleted = updatedResponse.data.every(stage => stage.Status === "Completed");
              if (allCompleted) {
                  console.log("✅ All Stages Completed! Stopping Polling.");
                  clearInterval(intervalId);
              }
          }
      }, 1000);
  }

  // ✅ Stop polling when modal is closed
  return () => {
      clearInterval(intervalId);
      setSelectedPipelineId(null);
  };
}, [openStageModal, selectedPipelineId]);



const fetchAndUpdatePipelineStatus = async () => {
  try {
      const response = await authAxios.get(`/pipelines?user_id=${user_id}`);
      if (Array.isArray(response.data)) {
          setPipelines(response.data);
          console.log("✅ Updated Pipelines:", response.data);
      }
  } catch (error) {
      console.error("❌ Failed to update pipeline status:", error);
  }
};



  const authAxios = useMemo(() => {
    return axios.create({
      baseURL: "http://localhost:30002",
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
  
    const stageCount = parseInt(numStages, 10); 
  
    if (!stageCount || stageCount <= 0) {  
      alert("Invalid number of stages! Please enter a valid number.");
      return;
    }
  
    const payload = {
      name: pipelineName,
      stages: stageCount, 
      is_parallel: isParallel ?? true, 
      user_id: getUserIdFromToken() || "default-user-id", 
      stage_names: stageNames || [], 
    };
  
    console.log("🚀 Sending Payload:", JSON.stringify(payload, null, 2));
  
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
            await fetchPipelineStages(pipelineId);
            console.log("✅ Response:", response.data);
        }

        setTimeout(fetchUserPipelines, 1000); // ✅ Refresh pipeline list
        await fetchPipelineStages(pipelineId); // ✅ Fetch stages immediately

        if (socket && socket.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify({ pipelineId, action: "start" }));
        }
    } catch (error) {
        console.error("❌ Failed to update pipeline status:", error.response?.data || error);
    }
};



  const handleSaveStageNames = () => {
    setPipelineStages(stageNames.map((name, index) => ({
      name,
      stage_number: index + 1
    }))); 
    setDialogOpen(false); 
  };
  

  const fetchPipelineStages = async (pipelineId,  pipelineName) => {
    try {
        console.log(`📥 Fetching stages for pipeline: ${pipelineId}`);

        const response = await authAxios.get(`/pipelines/${pipelineId}/stages`);

        if (Array.isArray(response.data)) {
            setSelectedPipelineStages(response.data);
            setSelectedPipelineId(pipelineId);
            setSelectedPipelineName(pipelineName); 
            setOpenStageModal(true);

            if (socket && socket.readyState === WebSocket.OPEN) {
                socket.send(JSON.stringify({ pipelineId, action: "track" }));
            }

            // ✅ Stop polling if the pipeline is already completed
            const isCompleted = response.data.every(stage => stage.Status === "Completed");
            if (isCompleted) {
                console.log("✅ Pipeline already completed, stopping further updates.");
                await fetchAndUpdatePipelineStatus();
                return;
            }

            // ✅ Poll for updates ONLY for this pipeline
            const intervalId = setInterval(async () => {
                const updatedResponse = await authAxios.get(`/pipelines/${pipelineId}/stages`);
                if (Array.isArray(updatedResponse.data)) {
                    setSelectedPipelineStages(updatedResponse.data);

                    // ✅ Stop polling if all stages are completed
                    const allCompleted = updatedResponse.data.every(stage => stage.Status === "Completed");
                    if (allCompleted) {
                        console.log("✅ All Stages Completed! Stopping Polling.");
                        clearInterval(intervalId);
                        await fetchAndUpdatePipelineStatus();
                    }
                }
            }, 5000);

            return () => clearInterval(intervalId);  // ✅ Cleanup when modal closes
        } else {
            console.error("❌ Unexpected response format:", response.data);
        }
    } catch (error) {
        console.error("❌ Failed to fetch pipeline stages:", error);
    }
};





const handleShowStages = async (pipelineID) => {
  setSelectedPipelineStages([]); // ✅ Clear previous stages before fetching
  await fetchPipelineStages(pipelineID);
  setOpenStageModal(true);
};


  const handleStageNameChange = (index, value) => {
    setStageNames((prev) => {
      const updatedStages = [...prev];
      updatedStages[index] = value;
      return updatedStages;
    });
  };

  const handleStageDialogOpen = () => {
    setStageNames(new Array(numStages).fill("")); 
    setDialogOpen(true);
  };
  
  const handleDeletePipeline = async (pipelineId) => {
    if (!window.confirm("Are you sure you want to delete this pipeline?")) return;
  
    try {
      await axios.delete(`http://localhost:30002/api/pipelines/${pipelineId}`);
      
      setPipelines(pipelines.filter(pipeline => pipeline.PipelineID !== pipelineId));
    } catch (error) {
      console.error("Error deleting pipeline:", error);
      alert("Failed to delete pipeline.");
    }
  };
  
  if (loading) return <CircularProgress />;
  
  return (
    <Box sx={{ paddingTop: 7,display: "flex" }}>
      <Topbar />
      <Sidebar />
      <Container maxWidth="md">
   
        <Box sx={{ mt: 5, p: 3, boxShadow: 3, borderRadius: 2 }}>
          <Typography variant="h5"  sx={{ fontWeight: "bold", color: "black", mb: 2 }}>Your Pipelines</Typography>
          {pipelines.length > 0 ? (
            <Table sx={{ minWidth: 500, border: "1px solid #ccc", borderRadius: "8px" }}>
              <TableHead>
                <TableRow sx={{ backgroundColor: "#f5f5f5" }}>
                  <TableCell><strong>Pipeline Name</strong></TableCell>
                  {/* <TableCell><strong>Pipeline ID</strong></TableCell> */}
                  <TableCell><strong>Status</strong></TableCell>
                  <TableCell><strong>Actions</strong></TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {pipelines.map((pipeline) => (
                  <TableRow key={pipeline.PipelineID}>
                    <TableCell>{pipeline.PipelineName}</TableCell>
                    {/* <TableCell>{pipeline.PipelineID}</TableCell> */}
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
                    <Box
    sx={{
      display: "flex",
      flexWrap: "wrap",
      gap: 2, 
      justifyContent: "center",
      "@media (max-width: 600px)": {
        flexDirection: "column", 
        alignItems: "center",
      },
    }}
  >
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
                          fetchPipelineStages(pipeline.PipelineID, pipeline.PipelineName);
                          setOpenStageModal(true); 
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
                        </Box>
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
        <DialogTitle>{selectedPipelineName ? `${selectedPipelineName} - Stages` : "Pipeline Stages"}</DialogTitle>
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
          <Typography variant="h5" sx={{ fontWeight: "bold", mb: 2 }}>Create New Pipeline</Typography>
          <Box sx={{ display: "flex", alignItems: "center", gap: 3, mt: 2 }}>
            <Typography variant="h6" sx={{ minWidth: 150, fontWeight: "bold", textAlign: "right" }}>Pipeline Name:</Typography>
            <TextField
              variant="outlined"
              value={pipelineName}
              onChange={(e) => setPipelineName(e.target.value)}
              placeholder="Enter pipeline name"
              sx={{ flexGrow: 1 }}
            />
          </Box>
          <Box sx={{ display: "flex", alignItems: "center", gap: 3, mt: 3 }}>
            <Typography variant="h6" sx={{ fontWeight: "bold", minWidth: 150, textAlign: "right" }}>No. of Stages:</Typography>
            <TextField
              select
              value={numStages}
              onChange={(e) => setNumStages(Number(e.target.value))}
              sx={{ width: 80 }}
            >
              {[...Array(20).keys()].map((num) => (
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



