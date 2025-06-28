import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import AuthService from '../services/AuthService';
import styles from './UserMenu.module.css';

const UserMenu = () => {
  const [showMenu, setShowMenu] = useState(false);
  const [user, setUser] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const profile = await AuthService.getProfile();
        setUser(profile);
      } catch (error) {
        console.error('Failed to fetch user profile');
      }
    };
    fetchUser();
  }, []);

  const handleLogout = () => {
    localStorage.removeItem('token');
    navigate('/login');
  };

  return (
    <div className={styles.menuWrapper}>
      <button className={styles.toggleButton} onClick={() => setShowMenu(!showMenu)}>
        {user?.username || 'User'} ▾
      </button>

      {showMenu && (
        <div className={styles.dropdown}>
          <button className={styles.menuItem} onClick={() => navigate('/profile')}>Profile</button>
          <button className={styles.menuItem} onClick={() => navigate('/change-password')}>Change Password</button>
          <button className={`${styles.menuItem} ${styles.logout}`} onClick={handleLogout}>Logout</button>
        </div>
      )}
    </div>
  );
};

export default UserMenu;
