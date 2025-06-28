import { useState } from "react";
import AuthService from "../services/AuthService";
import styles from "./AuthForm.module.css";
import { useNavigate } from "react-router-dom";

function Register() {
  const [name, setName] = useState("");
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const userData = { name, username, email, password };
      await AuthService.register(userData);
      alert("Registration successful! Please login.");
      navigate("/login");
    } catch (error) {
      alert("Registration failed!");
    }
  };

  return (
    <div className={styles.pageWrapper}>
      <div className={styles.container}>
        <h2 className={styles.title}>Register</h2>
        <form onSubmit={handleSubmit}>
          <div className={styles.formGroup}>
            <label className={styles.label}>Name:</label>
            <input
              className={styles.input}
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
            />
          </div>

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
            <label className={styles.label}>Email:</label>
            <input
              className={styles.input}
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
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
            Register
          </button>
        </form>
      </div>
    </div>
  );
}

export default Register;
