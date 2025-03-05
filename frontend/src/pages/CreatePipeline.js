import React, { useState } from "react";
import apiService from "../services/apiService";

const CreatePipeline = () => {
  const [stages, setStages] = useState("");

  const handleCreate = async () => {
    await apiService.createPipeline(stages);
  };

  return (
    <div>
      <h2>Create Pipeline</h2>
      <input
        type="number"
        placeholder="Enter number of stages"
        value={stages}
        onChange={(e) => setStages(e.target.value)}
      />
      <button onClick={handleCreate}>Create</button>
    </div>
  );
};

export default CreatePipeline;
