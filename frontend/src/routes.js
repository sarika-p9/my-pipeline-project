import React from "react";
import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import Login from "./components/Login";
import Register from "./components/Register";
import Dashboard from "./components/Dashboard";
import CreatePipeline from "./components/CreatePipeline";
import ExecutePipeline from "./components/ExecutePipeline";
import GetStatus from "./components/GetStatus";
import CancelPipeline from "./components/CancelPipeline";
import UserDetails from "./components/UserDetails";

// Function to check if user is authenticated
const isAuthenticated = () => {
  return localStorage.getItem("token") !== null; // Returns true if token exists
};

// Protected Route Wrapper
const ProtectedRoute = ({ element }) => {
  return isAuthenticated() ? element : <Navigate to="/login" />;
};

const AppRoutes = () => (
  <Router>
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />
      <Route path="/dashboard" element={<ProtectedRoute element={<Dashboard />} />} />
      <Route path="/create-pipeline" element={<ProtectedRoute element={<CreatePipeline />} />} />
      <Route path="/execute-pipeline" element={<ProtectedRoute element={<ExecutePipeline />} />} />
      <Route path="/get-status" element={<ProtectedRoute element={<GetStatus />} />} />
      <Route path="/cancel-pipeline" element={<ProtectedRoute element={<CancelPipeline />} />} />
      <Route path="/user-details" element={<ProtectedRoute element={<UserDetails />} />} />
      <Route path="*" element={<Navigate to={isAuthenticated() ? "/dashboard" : "/login"} />} />
    </Routes>
  </Router>
);

export default AppRoutes;
