import React, { useState } from "react";
import apiService from "../services/apiService";

const CancelPipeline = () => {
  const [pipelineId, setPipelineId] = useState("");

  const handleCancel = async () => {
    await apiService.cancelPipeline(pipelineId);
  };

  return (
    <div>
      <h2>Cancel Pipeline</h2>
      <input
        type="text"
        placeholder="Enter Pipeline ID"
        value={pipelineId}
        onChange={(e) => setPipelineId(e.target.value)}
      />
      <button onClick={handleCancel}>Cancel</button>
    </div>
  );
};

export default CancelPipeline;
