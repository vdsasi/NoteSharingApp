import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import DOMPurify from "dompurify";
import ReactQuill from "react-quill";
import "react-quill/dist/quill.snow.css";
import styles from "./CreateNote.module.css";

const CreateNote = () => {
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [tags, setTags] = useState("");
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();

    const sanitizedContent = DOMPurify.sanitize(content);

    try {
      const response = await fetch("http://localhost:8080/api/notes", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include", // Send cookies with request
        body: JSON.stringify({
          title,
          content: sanitizedContent,
          tags: tags
            .split(",")
            .map((t) => t.trim())
            .filter(Boolean),
        }),
      });

      if (response.ok) {
        navigate("/dashboard");
      } else {
        console.error("Failed to create note");
      }
    } catch (error) {
      console.error("Error:", error);
    }
  };

  return (
    <div className={styles.container}>
      <h2 className={styles.heading}>Create New Note</h2>
      <form onSubmit={handleSubmit} className={styles.form}>
        <input
          type="text"
          className={styles.input}
          placeholder="Title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          required
        />

        <ReactQuill
          value={content}
          onChange={setContent}
          className={styles.quill}
          theme="snow"
          placeholder="Write your note here..."
        />

        <input
          type="text"
          className={styles.tagInput}
          placeholder="Tags (comma-separated)"
          value={tags}
          onChange={(e) => setTags(e.target.value)}
        />

        <button type="submit" className={styles.submitButton}>
          Create Note
        </button>
      </form>
    </div>
  );
};

export default CreateNote;
