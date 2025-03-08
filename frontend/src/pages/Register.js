// import React, { useState } from "react";
// import { TextField, Button, Typography, Container, Box, Link } from "@mui/material";
// import { useNavigate } from "react-router-dom";
// import axios from "axios";

// const Register = () => {
//   const [email, setEmail] = useState("");
//   const [password, setPassword] = useState("");
//   const [message, setMessage] = useState("");
//   const navigate = useNavigate();

//   const handleRegister = async () => {
//     setMessage("");
//     try {
//       const response = await axios.post("http://localhost:8080/register", { email, password });
//       setMessage("User Registered Successfully!!\nCheck your email and confirm authentication.");
//     } catch (error) {
//       setMessage("Registration failed. Please try again.");
//     }
//   };

//   return (
//     <Container maxWidth="sm">
//       <Box sx={{ textAlign: "center", mt: 5, p: 3, boxShadow: 3, borderRadius: 2 }}>
//         <Typography variant="h4" gutterBottom>Register</Typography>
//         <TextField label="Email" fullWidth margin="normal" value={email} onChange={(e) => setEmail(e.target.value)} />
//         <TextField label="Password" type="password" fullWidth margin="normal" value={password} onChange={(e) => setPassword(e.target.value)} />
//         <Button variant="contained" color="primary" fullWidth onClick={handleRegister} sx={{ mt: 2 }}>Register</Button>
//         {message && <Typography sx={{ mt: 2, color: "green" }}>{message}</Typography>}
//         <Typography sx={{ mt: 2 }}>
//           Already have an account? <Link onClick={() => navigate("/login")} sx={{ cursor: "pointer" }}>Login here</Link>
//         </Typography>
//       </Box>
//     </Container>
//   );
// };

// export default Register;