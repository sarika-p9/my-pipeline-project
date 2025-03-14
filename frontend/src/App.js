import React from "react";
import AuthPage from "./components/authpage";
import AppRoutes from "./routes";

const isAuthenticated = () => {
  return localStorage.getItem("token") !== null; 
};

function App() {
  return isAuthenticated() ? <AppRoutes /> : <AuthPage />;
}

export default App;
