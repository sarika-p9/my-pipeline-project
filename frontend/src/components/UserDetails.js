import { useEffect, useState } from "react";

const UserDetails = () => {
  const [user, setUser] = useState(null);
  const token = localStorage.getItem("token");

  useEffect(() => {
    const fetchUser = async () => {
      if (!token) return;

      try {
        const response = await fetch("http://localhost:8080/user", {
          method: "GET",
          headers: { Authorization: `Bearer ${token}` },
        });

        if (!response.ok) {
          throw new Error("Failed to fetch user data");
        }

        const userData = await response.json();
        setUser(userData);
      } catch (err) {
        console.error(err);
      }
    };

    fetchUser();
  }, [token]);

  return (
    <div style={{ 
      backgroundColor: "#42A5F5", 
      minHeight: "98vh", 
      width: "98vw", 
      margin: "0", 
      padding: "0",
      display: "flex",
      justifyContent: "center",
      alignItems: "center"
    }}>
      <div style={{ textAlign: "center" }}>
        <h2>User Details</h2>
        {user ? (
          <div>
            <p><strong>User ID:</strong> {user.UserID}</p>
            <p><strong>Email:</strong> {user.Email}</p>
            <p><strong>Role:</strong> {user.Role}</p>
          </div>
        ) : (
          <p>Loading user data...</p>
        )}
      </div>
    </div>
  );

};

export default UserDetails;
