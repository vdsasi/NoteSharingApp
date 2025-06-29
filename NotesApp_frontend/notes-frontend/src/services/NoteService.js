import axios from '../api/axios';
const BASE_URL = 'http://localhost:8080/api/notes';

export const fetchNotes = () => axios.get(BASE_URL);
export const togglePinNote = (noteId) => axios.post(`${BASE_URL}/${noteId}/pin`);
export const getNoteById = (id) => axios.get(`${BASE_URL}/${id}`);
export const updateNote = (id, data) => axios.put(`${BASE_URL}/${id}`, data);
export const deleteNoteById = (id) => axios.delete(`${BASE_URL}/${id}`);
export const getTrashedNotes = () => axios.get(`${BASE_URL}/trash`);
export const restoreNoteById = (id) => axios.post(`${BASE_URL}/${id}/restore`);
export const getNoteVersions = (noteId) => axios.get(`${BASE_URL}/${noteId}/versions`);
export const restoreNoteVersion = (noteId, versionId) =>
  axios.post(`${BASE_URL}/version-restore/${noteId}/${versionId}`);
export const getCollaborators = (noteId) => axios.get(`${BASE_URL}/${noteId}/collaborators`);
export const addCollaborator = (noteId, username) => axios.post(`${BASE_URL}/${noteId}/share`, { username });
export const removeCollaborator = (noteId, username) => axios.delete(`${BASE_URL}/${noteId}/share`, { data: { username } });

