import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { togglePinNote, deleteNoteById } from "../services/NoteService";
import styles from "./NoteCard.module.css";

const NoteCard = ({ note, onPinToggle, onDelete }) => {
  const navigate = useNavigate();
  const [isPinned, setIsPinned] = useState(note.pinned);

  useEffect(() => {
    setIsPinned(note.pinned);
  }, [note.pinned]);

  const handleDelete = () => {
    if (window.confirm("Are you sure you want to delete this note?")) {
      deleteNoteById(note.id)
        .then(() => {
          alert("Note deleted");
          if (onDelete) onDelete(note.id);
        })
        .catch((err) => alert("Failed to delete note: " + err.message));
    }
  };

  const handlePinClick = () => {
    togglePinNote(note.id)
      .then(() => {
        setIsPinned(!isPinned);
        if (onPinToggle) onPinToggle(note.id);
      })
      .catch((err) => {
        alert("Failed to toggle pin status");
        console.error(err);
      });
  };

  const createSnippet = (html) => {
    const div = document.createElement("div");
    div.innerHTML = html;
    const text = div.textContent || div.innerText || "";
    return text.length > 100 ? text.substring(0, 100) + "..." : text;
  };

  return (
    <div className={styles.card} onClick={() => navigate(`/view/${note.id}`)}>
      <h2 className={styles.title}>{note.title}</h2>
      <p className={styles.snippet}>{createSnippet(note.content)}</p>
      <p className={styles.updated}>
        Updated: {new Date(note.updatedAt).toLocaleString()}
      </p>

      <div className={styles.tags}>
        {Array.isArray(note.tags) && note.tags.length > 0 ? (
          note.tags.map((tag, index) => {
            if (!tag) return null;
            const tagName =
              typeof tag === "string" ? tag : tag.name || "Unknown";
            const tagKey =
              typeof tag === "string" ? tag + index : tag.id || tagName + index;

            return (
              <span key={tagKey} className={styles.tag}>
                {tagName}
              </span>
            );
          })
        ) : (
          <span className={styles.tag}>No tags</span>
        )}
      </div>

      <div className={styles.buttons}>
        <button
          className={styles.button}
          onClick={(e) => {
            e.stopPropagation();
            handlePinClick();
          }}
        >
          {isPinned ? "Unpin" : "Pin"}
        </button>

        <button
          className={styles.button}
          onClick={(e) => {
            e.stopPropagation();
            navigate(`/edit/${note.id}`);
          }}
        >
          Edit
        </button>

        <button
          className={styles.button}
          onClick={(e) => {
            e.stopPropagation();
            handleDelete();
          }}
        >
          Delete
        </button>
      </div>
    </div>
  );
};

export default NoteCard;
