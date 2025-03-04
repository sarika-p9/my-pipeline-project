import React from "react";
import { Navigate } from "react-router-dom";

const ProtectedRoute = ({ element }) => {
  const token = localStorage.getItem("token");
  console.log("Checking authentication. Token found:", token); // Debugging log

  return token ? element : <Navigate to="/login" />;
};

export default ProtectedRoute;
