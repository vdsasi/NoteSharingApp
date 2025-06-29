import React, { useState, useEffect } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import ReactQuill from "react-quill";
import "react-quill/dist/quill.snow.css";
import { getNoteById, updateNote, getCollaborators, addCollaborator, removeCollaborator } from "../services/NoteService";
import AuthService from "../services/AuthService";
import ShareModal from "../components/ShareModal";
import styles from "./EditNote.module.css";

const EditNote = () => {
  const { id } = useParams();
  const navigate = useNavigate();

  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [tags, setTags] = useState([]);
  const [newTag, setNewTag] = useState("");
  const [showShare, setShowShare] = useState(false);
  const [collaborators, setCollaborators] = useState([]);
  const [isOwner, setIsOwner] = useState(false);
  const [shareLoading, setShareLoading] = useState(false);
  const [shareError, setShareError] = useState("");
  const [currentUser, setCurrentUser] = useState(null);
  const [notFound, setNotFound] = useState(false);

  useEffect(() => {
    getNoteById(id)
      .then((res) => {
        const note = res.data;
        setTitle(note.title);
        setContent(note.content);

        if (note.tags && Array.isArray(note.tags)) {
          const safeTags = note.tags
            .filter(
              (t) => t && (typeof t === "object" || typeof t === "string")
            )
            .map((t) => (typeof t === "string" ? t : t.name || "Unnamed Tag"));
          setTags(safeTags);
        } else if (note.tagNames && Array.isArray(note.tagNames)) {
          setTags(note.tagNames);
        } else {
          setTags([]);
        }
        setIsOwner(note.userId === currentUser?.id);
        setNotFound(false);
      })
      .catch((err) => {
        if (err?.response?.status === 404 || err?.response?.status === 403) {
          setNotFound(true);
        } else {
          console.error("Failed to fetch note:", err);
        }
      });
    // Fetch collaborators
    getCollaborators(id).then(res => setCollaborators(res.data)).catch(() => setCollaborators([]));
  }, [id, currentUser]);

  useEffect(() => {
    AuthService.getProfile().then(setCurrentUser).catch(() => setCurrentUser(null));
  }, []);

  const handleAddTag = () => {
    const trimmedTag = newTag.trim();
    if (trimmedTag && !tags.includes(trimmedTag)) {
      setTags([...tags, trimmedTag]);
      setNewTag("");
    }
  };

  const handleRemoveTag = (tagToRemove) => {
    setTags(tags.filter((tag) => tag !== tagToRemove));
  };

  const handleUpdateNote = () => {
    const data = { title, content, tags };
    updateNote(id, data)
      .then(() => {
        navigate("/dashboard");
      })
      .catch((err) => {
        console.error("Failed to update note:", err);
      });
  };

  const handleAddCollaborator = async (username) => {
    setShareLoading(true); setShareError("");
    try {
      await addCollaborator(id, username);
      const res = await getCollaborators(id);
      setCollaborators(res.data);
    } catch (e) {
      setShareError(e?.response?.data?.error || "Failed to add collaborator");
    }
    setShareLoading(false);
  };
  const handleRemoveCollaborator = async (username) => {
    setShareLoading(true); setShareError("");
    try {
      await removeCollaborator(id, username);
      const res = await getCollaborators(id);
      setCollaborators(res.data);
    } catch (e) {
      setShareError(e?.response?.data?.error || "Failed to remove collaborator");
    }
    setShareLoading(false);
  };

  if (notFound) {
    return (
      <div className={styles.container}>
        <h2 className={styles.heading}>Access Denied</h2>
        <p>You do not have access to this note or it does not exist.</p>
        <Link to="/dashboard" className={styles.backLink}>&larr; Back to Dashboard</Link>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <Link to="/dashboard" className={styles.backLink}>
        &larr; Back
      </Link>

      <h1 className={styles.heading}>Edit Note</h1>

      <input
        type="text"
        placeholder="Title"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        className={styles.input}
      />

      <ReactQuill
        theme="snow"
        value={content}
        onChange={setContent}
        className={styles.quillEditor}
      />

      <div className={styles.tagsContainer}>
        <h3>Tags</h3>
        <div className={styles.tagList}>
          {tags.map((tag, index) => (
            <div key={`${tag}-${index}`} className={styles.tag}>
              <span>{tag}</span>
              <button
                type="button"
                onClick={() => handleRemoveTag(tag)}
                className={styles.removeBtn}
              >
                &times;
              </button>
            </div>
          ))}
        </div>
        <div className={styles.tagInputRow}>
          <input
            type="text"
            placeholder="Add tag"
            value={newTag}
            onChange={(e) => setNewTag(e.target.value)}
            className={styles.tagInput}
          />
          <button onClick={handleAddTag} className={styles.addTagBtn}>
            Add
          </button>
        </div>
      </div>

      <button onClick={handleUpdateNote} className={styles.updateBtn}>
        Update Note
      </button>
      <button onClick={() => setShowShare(true)} className={styles.updateBtn} style={{marginBottom: 12}}>
        Share
      </button>
      <ShareModal
        isOpen={showShare}
        onClose={() => setShowShare(false)}
        collaborators={collaborators}
        onAdd={handleAddCollaborator}
        onRemove={handleRemoveCollaborator}
        isOwner={isOwner}
        loading={shareLoading}
        error={shareError}
      />
    </div>
  );
};

export default EditNote;
