import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Dashboard from "./components/Dashboard";
import UserDetails from "./pages/UserDetails";
import CreatePipeline from "./pages/CreatePipeline";
import ExecutePipeline from "./pages/ExecutePipeline";
import GetStatus from "./pages/GetStatus";
import CancelPipeline from "./pages/CancelPipeline";
import Login from "./pages/Login";
import Register from "./pages/Register";

const App = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route path="/user-details" element={<UserDetails />} />
        <Route path="/create-pipeline" element={<CreatePipeline />} />
        <Route path="/execute-pipeline" element={<ExecutePipeline />} />
        <Route path="/get-status" element={<GetStatus />} />
        <Route path="/cancel-pipeline" element={<CancelPipeline />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/dashboard" element={<Dashboard />} />
      </Routes>
    </Router>
  );
};

export default App;
