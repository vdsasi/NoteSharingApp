import { useState } from "react";
import AuthService from "../services/AuthService";
import styles from "./AuthForm.module.css";
import { useNavigate } from "react-router-dom";

function Login() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await AuthService.login(username, password);
      const profile = await AuthService.getProfile();
      alert(`Welcome ${profile.username}!`);
      navigate("/dashboard"); // or wherever you want to go
    } catch (error) {
      alert("Invalid credentials!");
    }
  };

  return (
    <div className={styles.pageWrapper}>
      <div className={styles.container}>
        <h2 className={styles.title}>Login</h2>
        <form onSubmit={handleSubmit}>
          <div className={styles.formGroup}>
            <label className={styles.label}>Username:</label>
            <input
              className={styles.input}
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />
          </div>

          <div className={styles.formGroup}>
            <label className={styles.label}>Password:</label>
            <input
              className={styles.input}
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>

          <button className={styles.button} type="submit">
            Login
          </button>
        </form>
      </div>
    </div>
  );
}

export default Login;
