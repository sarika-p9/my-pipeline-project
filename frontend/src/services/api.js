const API_BASE_URL = "http://localhost:5000"; // Change this if different

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
    const response = await fetch(`${API_BASE_URL}/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });
    return await response.json();
  } catch (error) {
    return { error: "Failed to connect to server" };
  }
};
