import { useState } from 'react';
import AuthService from '../services/AuthService';
import styles from './ChangePassword.module.css';

const ChangePassword = () => {
  const [formData, setFormData] = useState({
    oldPassword: '',
    newPassword: '',
    confirmPassword: ''
  });
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
    setMessage('');
    setError('');
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (formData.newPassword !== formData.confirmPassword) {
      setError('New passwords do not match.');
      return;
    }

    try {
      await AuthService.changePassword(formData.oldPassword, formData.newPassword);
      setMessage('Password changed successfully!');
      setFormData({ oldPassword: '', newPassword: '', confirmPassword: '' });
    } catch (error) {
      console.error(error);
      setError('Failed to change password.');
    }
  };

  return (
    <div className={styles.container}>
      <h2 className={styles.heading}>Change Password</h2>
      <form onSubmit={handleSubmit} className={styles.form}>
        <input
          type="password"
          name="oldPassword"
          placeholder="Old Password"
          value={formData.oldPassword}
          onChange={handleChange}
          className={styles.input}
          required
        />
        <input
          type="password"
          name="newPassword"
          placeholder="New Password"
          value={formData.newPassword}
          onChange={handleChange}
          className={styles.input}
          required
        />
        <input
          type="password"
          name="confirmPassword"
          placeholder="Confirm New Password"
          value={formData.confirmPassword}
          onChange={handleChange}
          className={styles.input}
          required
        />
        <button type="submit" className={styles.button}>
          Change Password
        </button>
      </form>
      {message && <p className={styles.success}>{message}</p>}
      {error && <p className={styles.error}>{error}</p>}
    </div>
  );
};

export default ChangePassword;
