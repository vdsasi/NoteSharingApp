import React from "react";
import { useNavigate } from "react-router-dom";
import UserMenu from "./UserMenu";
import styles from "./Header.module.css";

const Header = () => {
  const navigate = useNavigate();

  return (
    <header className={styles.header}>
      <div className={styles.left} onClick={() => navigate("/dashboard")}>
        <span className={styles.title}>My Notes</span>
      </div>
      <div className={styles.right}>
        <UserMenu />
      </div>
    </header>
  );
};

export default Header;
