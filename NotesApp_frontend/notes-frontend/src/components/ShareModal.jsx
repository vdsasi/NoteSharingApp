import React from "react";
import styles from "./ShareModal.module.css";
import { searchUsers } from "../services/UserService";

const ShareModal = ({ isOpen, onClose, collaborators, onAdd, onRemove, isOwner, loading, error }) => {
  const [username, setUsername] = React.useState("");
  const [userResults, setUserResults] = React.useState([]);
  const [searching, setSearching] = React.useState(false);
  const [searchError, setSearchError] = React.useState("");

  React.useEffect(() => {
    if (!username || username.length < 2) {
      setUserResults([]);
      setSearchError("");
      return;
    }
    setSearching(true);
    setSearchError("");
    searchUsers(username)
      .then(users => setUserResults(users || []))
      .catch(() => {
        setUserResults([]);
        setSearchError("No users found");
      })
      .finally(() => setSearching(false));
  }, [username]);

  if (!isOpen) return null;

  return (
    <div className={styles.modalOverlay}>
      <div className={styles.modalContent}>
        <h2>Share Note</h2>
        {isOwner && (
          <div className={styles.addRow}>
            <input
              type="text"
              placeholder="Search username"
              value={username}
              onChange={e => setUsername(e.target.value)}
              className={styles.input}
            />
            {searching && <span className={styles.searching}>Searching...</span>}
            {username && userResults.length > 0 && (
              <div className={styles.userDropdown}>
                {userResults.map(user => (
                  <div key={user.id} className={styles.userOption} onClick={() => { onAdd(user.username); setUsername(""); setUserResults([]); }}>
                    {user.username} ({user.email})
                  </div>
                ))}
              </div>
            )}
            {searchError && <div className={styles.error}>{searchError}</div>}
          </div>
        )}
        <div className={styles.collabList}>
          <h4>Collaborators:</h4>
          {collaborators.length === 0 && <div>No collaborators yet.</div>}
          {collaborators.map((c) => (
            <div key={c.id} className={styles.collabItem}>
              <span>{c.username} ({c.email})</span>
              {isOwner && (
                <button onClick={() => onRemove(c.username)} className={styles.removeBtn} disabled={loading}>
                  Remove
                </button>
              )}
            </div>
          ))}
        </div>
        {error && <div className={styles.error}>{error}</div>}
        <button onClick={onClose} className={styles.closeBtn}>Close</button>
      </div>
    </div>
  );
};

export default ShareModal;
