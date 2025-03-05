import React, { useState } from "react";
import apiService from "../services/apiService";

const ExecutePipeline = () => {
  const [pipelineId, setPipelineId] = useState("");

  const handleExecute = async () => {
    await apiService.executePipeline(pipelineId);
  };

  return (
    <div>
      <h2>Execute Pipeline</h2>
      <input
        type="text"
        placeholder="Enter Pipeline ID"
        value={pipelineId}
        onChange={(e) => setPipelineId(e.target.value)}
      />
      <button onClick={handleExecute}>Execute</button>
    </div>
  );
};

export default ExecutePipeline;
