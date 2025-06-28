import { useEffect, useState } from 'react';
import AuthService from '../services/AuthService';
import styles from './Profile.module.css';

const Profile = () => {
  const [profile, setProfile] = useState(null);

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const data = await AuthService.getProfile();
        setProfile(data);
      } catch (error) {
        console.error('Failed to fetch profile:', error);
      }
    };
    fetchUser();
  }, []);

  if (!profile) return <p className={styles.loading}>Loading profile...</p>;

  return (
    <div className={styles.container}>
      <h2 className={styles.heading}>User Profile</h2>
      <div className={styles.info}>
        <p><strong>Username:</strong> {profile.username}</p>
        <p><strong>Email:</strong> {profile.email}</p>
      </div>
    </div>
  );
};

export default Profile;
