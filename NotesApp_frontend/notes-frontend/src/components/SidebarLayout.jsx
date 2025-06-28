import React from "react";
import { Link, useLocation } from "react-router-dom";
import Header from "./Header";
import styles from "./SidebarLayout.module.css";

const SidebarLayout = ({ children }) => {
  const location = useLocation();

  const isActive = (path) =>
    location.pathname.startsWith(path) ? styles.active : "";

  return (
    <div className={styles.wrapper}>
      <Header />
      <div className={styles.container}>
        {/* Sidebar */}
        <div className={styles.sidebar}>
          <nav className={styles.nav}>
            <Link to="/dashboard" className={`${styles.link} ${isActive("/dashboard")}`}>
              Dashboard
            </Link>
            <Link to="/trash" className={`${styles.link} ${isActive("/trash")}`}>
              Trash
            </Link>
            <Link to="/logout" className={`${styles.link} ${styles.logout}`}>
              Logout
            </Link>
          </nav>
        </div>

        {/* Main content */}
        <div className={styles.mainContent}>
          {children}
        </div>
      </div>
    </div>
  );
};

export default SidebarLayout;
