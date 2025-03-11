const BASE_URL = "http://localhost:5000"; // Replace with your actual API URL

const apiService = {
  createPipeline: async (pipelineName, stages) => {
    try {
      const response = await fetch(`${BASE_URL}/create-pipeline`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ pipeline_name: pipelineName, stages }),
      });
      return response.json();
    } catch (error) {
      console.error("Error creating pipeline:", error);
    }
  },

  executePipeline: async (pipelineId) => {
    try {
      const response = await fetch(`${BASE_URL}/execute-pipeline`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ pipelineId }),
      });
      return response.json();
    } catch (error) {
      console.error("Error executing pipeline:", error);
    }
  },

  getPipelineStatus: async (pipelineId) => {
    try {
      const response = await fetch(`${BASE_URL}/get-status/${pipelineId}`);
      return response.json();
    } catch (error) {
      console.error("Error fetching pipeline status:", error);
    }
  },

  cancelPipeline: async (pipelineId) => {
    try {
      const response = await fetch(`${BASE_URL}/cancel-pipeline`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ pipelineId }),
      });
      return response.json();
    } catch (error) {
      console.error("Error canceling pipeline:", error);
    }
  },
};

export default apiService;
