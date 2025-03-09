import React from "react";
import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import Dashboard from "./components/Dashboard";
import CreatePipeline from "./pages/CreatePipeline";
import UserProfile from "./pages/UserProfile";
const isAuthenticated = () => {
  return localStorage.getItem("token") !== null; 
};
const ProtectedRoute = ({ element }) => {
  return isAuthenticated() ? element : <Navigate to="/login" />;
};
const AppRoutes = () => (
  <Router>
    <Routes>
      <Route path="/dashboard" element={<ProtectedRoute element={<Dashboard />} />} />
      <Route path="/create-pipeline" element={<ProtectedRoute element={<CreatePipeline />} />} />
      <Route path="/user-profile" element={<ProtectedRoute element={<UserProfile />} />} />
      <Route path="*" element={<Navigate to={isAuthenticated() ? "/dashboard" : "/login"} />} />
    </Routes>
  </Router>
);
export default AppRoutes;
