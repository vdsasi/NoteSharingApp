import React, { useEffect, useState } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import {
  getNoteById,
  deleteNoteById,
  getNoteVersions,
  restoreNoteVersion,
} from "../services/NoteService";
import styles from "./ViewNote.module.css";

const ViewNote = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [note, setNote] = useState(null);
  const [versions, setVersions] = useState([]);

  useEffect(() => {
    getNoteById(id)
      .then((res) => setNote(res.data))
      .catch((err) => console.error("Failed to fetch note", err));
  }, [id]);

  useEffect(() => {
    getNoteVersions(id)
      .then((res) => setVersions(res.data))
      .catch((err) => console.error("Failed to fetch versions", err));
  }, [id]);

  const handleDelete = () => {
    if (window.confirm("Are you sure you want to delete this note?")) {
      deleteNoteById(id)
        .then(() => {
          alert("Note deleted");
          navigate("/dashboard");
        })
        .catch((err) => alert("Failed to delete note: " + err.message));
    }
  };

  const handleRestoreVersion = (versionId) => {
    if (window.confirm("Restore to this version?")) {
      restoreNoteVersion(id, versionId)
        .then(() => {
          alert("Version restored");
          window.location.reload();
        })
        .catch((err) => alert("Failed to restore version: " + err.message));
    }
  };

  if (!note) return <p>Loading...</p>;

  return (
    <div className={styles.container}>
      <Link to="/dashboard" className={styles.backLink}>
        &larr; Back to Notes
      </Link>

      <div className={styles.header}>
        <h1 className={styles.title}>{note.title}</h1>
        <div className={styles.actions}>
          <button
            onClick={() => navigate(`/edit/${id}`)}
            className={`${styles.editBtn} ${styles.button}`}
          >
            Edit
          </button>
          <button
            onClick={handleDelete}
            className={`${styles.deleteBtn} ${styles.button}`}
          >
            Delete
          </button>
        </div>
      </div>

      <p className={styles.timestamp}>
        Updated: {new Date(note.updatedAt).toLocaleString()}
      </p>

      <div
        className={styles.content}
        dangerouslySetInnerHTML={{ __html: note.content }}
      />

      {note.tags?.length > 0 && (
        <div className={styles.tags}>
          {note.tags
            .filter((tag) => tag)
            .map((tag, i) => {
              const tagName =
                typeof tag === "string"
                  ? tag
                  : typeof tag.name === "string"
                  ? tag.name
                  : "Tag";
              return (
                <span key={tagName + i} className={styles.tag}>
                  {tagName}
                </span>
              );
            })}
        </div>
      )}

      <div className={styles.versionSection}>
        <h2 className={styles.title}>Version History</h2>
        {versions.length === 0 ? (
          <p className={styles.timestamp}>No versions available</p>
        ) : (
          versions.map((version, index) => (
            <div key={index} className={styles.versionItem}>
              <p className={styles.subtitle}>{version.title}</p>
              <p className={styles.timestamp}>
                Saved at: {new Date(version.versionedAt).toLocaleString()}
              </p>
              <div
                className={styles.content}
                dangerouslySetInnerHTML={{ __html: version.content }}
              />
              <button
                onClick={() => handleRestoreVersion(index + 1)}
                className={styles.restoreBtn}
              >
                Restore this Version
              </button>
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default ViewNote;
