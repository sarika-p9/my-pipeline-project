import React, { useState } from "react";
import apiService from "../services/apiService";

const GetStatus = () => {
  const [pipelineId, setPipelineId] = useState("");
  const [status, setStatus] = useState("");

  const handleGetStatus = async () => {
    const response = await apiService.getPipelineStatus(pipelineId);
    setStatus(response);
  };

  return (
    <div>
      <h2>Get Pipeline Status</h2>
      <input
        type="text"
        placeholder="Enter Pipeline ID"
        value={pipelineId}
        onChange={(e) => setPipelineId(e.target.value)}
      />
      <button onClick={handleGetStatus}>Get Status</button>
      {status && <p>Status: {status}</p>}
    </div>
  );
};

export default GetStatus;
