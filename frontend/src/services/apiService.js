const API_BASE_URL = "http://localhost:5000"; // Adjust this if needed

// User Authentication
export const registerUser = async (email, password) => {
  try {
    const response = await fetch(`${API_BASE_URL}/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });
    return await response.json();
  } catch (error) {
    return { error: "Failed to connect to server" };
  }
};

export const loginUser = async (email, password) => {
  try {
    console.log("Sending Login Request:", { email, password });

    const response = await fetch(`${API_BASE_URL}/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });

    const data = await response.json();
    console.log("Login Response:", data);

    if (response.ok && data.token) {
      localStorage.setItem("token", data.token);
    } else {
      console.error("Login Failed:", data.error || "Unknown error");
    }

    return data;
  } catch (error) {
    console.error("Login Error:", error);
    return { error: "Failed to connect to server" };
  }
};

// Fetch User Token (Assuming it's stored in localStorage)
const getToken = () => localStorage.getItem("token");

// Pipelines API
export const createPipeline = async (stages) => {
  try {
    const response = await fetch(`${API_BASE_URL}/pipelines`, {
      method: "POST",
      headers: { 
        "Content-Type": "application/json",
        Authorization: `Bearer ${getToken()}`,
      },
      body: JSON.stringify({ stages, is_parallel: true }),
    });
    return await response.json();
  } catch (error) {
    return { error: "Failed to create pipeline" };
  }
};

export const executePipeline = async (pipelineId, userId) => {
  try {
    const response = await fetch(`${API_BASE_URL}/pipelines/${pipelineId}/start`, {
      method: "POST",
      headers: { 
        "Content-Type": "application/json",
        Authorization: `Bearer ${getToken()}`,
      },
      body: JSON.stringify({ input: {}, is_parallel: true, user_id: userId }),
    });
    return await response.json();
  } catch (error) {
    return { error: "Failed to execute pipeline" };
  }
};

export const getPipelineStatus = async (pipelineId) => {
  try {
    const response = await fetch(`${API_BASE_URL}/pipelines/${pipelineId}/status`, {
      method: "GET",
      headers: { Authorization: `Bearer ${getToken()}` },
    });
    return await response.json();
  } catch (error) {
    return { error: "Failed to get pipeline status" };
  }
};

export const cancelPipeline = async (pipelineId) => {
  try {
    const response = await fetch(`${API_BASE_URL}/pipelines/${pipelineId}/cancel`, {
      method: "POST",
      headers: { 
        "Content-Type": "application/json",
        Authorization: `Bearer ${getToken()}`,
      },
      body: JSON.stringify({}),
    });
    return await response.json();
  } catch (error) {
    return { error: "Failed to cancel pipeline" };
  }
};
