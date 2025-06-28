import React, { useEffect, useState } from "react";
import { fetchNotes } from "../services/NoteService";
import { useNavigate } from "react-router-dom";
import NoteCard from "../components/NoteCard";
import styles from "./Dashboard.module.css";

const Dashboard = () => {
  const [notes, setNotes] = useState([]);
  const navigate = useNavigate();

  useEffect(() => {
    fetchNotes()
      .then((res) => {
        const sortedNotes = res.data.sort((a, b) =>
          a.pinned === b.pinned ? 0 : a.pinned ? -1 : 1
        );
        setNotes(sortedNotes);
      })
      .catch((err) => console.error("Error fetching notes:", err));
  }, []);

  const handlePinToggle = (noteId) => {
    fetchNotes()
      .then((res) => setNotes(res.data))
      .catch((err) => console.error("Error fetching notes:", err));
  };

  const handleNoteDeleted = (deletedId) => {
    setNotes(notes.filter((note) => note.id !== deletedId));
  };

  return (
    <div className={styles.pageWrapper}>
      <button
        className={styles.createButton}
        onClick={() => navigate("/create")}
      >
        + Create New Note
      </button>

      <h1 className={styles.header}>Your Notes</h1>

      <div className={styles.notesGrid}>
        {notes.length > 0 ? (
          notes.map((note) => (
            <NoteCard
              key={note.id}
              note={note}
              onPinToggle={handlePinToggle}
              onDelete={handleNoteDeleted}
            />
          ))
        ) : (
          <p>No notes found.</p>
        )}
      </div>
    </div>
  );
};

export default Dashboard;
