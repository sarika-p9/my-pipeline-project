import React from "react";
import AuthPage from "./components/authpage";
import AppRoutes from "./routes";

const isAuthenticated = () => {
  return localStorage.getItem("token") !== null; // Check if user is logged in
};

function App() {
  return isAuthenticated() ? <AppRoutes /> : <AuthPage />;
}

export default App;
