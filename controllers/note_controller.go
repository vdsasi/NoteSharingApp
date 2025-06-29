package controllers

import (
    "net/http"
    "notes-app/models"
    "notes-app/services"
    "notes-app/config"
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type NoteController struct {
    noteService *services.NoteService
}

func NewNoteController(noteService *services.NoteService) *NoteController {
    return &NoteController{noteService: noteService}
}

func (nc *NoteController) GetAll(c *gin.Context) {
    user := c.MustGet("user").(*models.User)
    
    notes, err := nc.noteService.GetUserNotes(user.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, notes)
}

func (nc *NoteController) Create(c *gin.Context) {
    user := c.MustGet("user").(*models.User)
    
    var req models.NoteRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    note, err := nc.noteService.CreateNote(user.ID, req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, note)
}

func (nc *NoteController) GetNote(c *gin.Context) {
    user := c.MustGet("user").(*models.User)
    
    noteID, err := primitive.ObjectIDFromHex(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
        return
    }

    note, err := nc.noteService.GetNote(noteID, user.ID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, note)
}

func (nc *NoteController) Update(c *gin.Context) {
    user := c.MustGet("user").(*models.User)
    
    noteID, err := primitive.ObjectIDFromHex(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
        return
    }

    var req models.NoteRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    note, err := nc.noteService.UpdateNote(noteID, user.ID, req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, note)
}

func (nc *NoteController) Delete(c *gin.Context) {
    user := c.MustGet("user").(*models.User)
    
    noteID, err := primitive.ObjectIDFromHex(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
        return
    }

    err = nc.noteService.DeleteNote(noteID, user.ID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}

func (nc *NoteController) GetTrashed(c *gin.Context) {
    user := c.MustGet("user").(*models.User)
    
    notes, err := nc.noteService.GetTrashedNotes(user.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, notes)
}

func (nc *NoteController) Restore(c *gin.Context) {
    user := c.MustGet("user").(*models.User)
    
    noteID, err := primitive.ObjectIDFromHex(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
        return
    }

    note, err := nc.noteService.RestoreNote(noteID, user.ID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, note)
}

func (nc *NoteController) TogglePin(c *gin.Context) {
    user := c.MustGet("user").(*models.User)
    
    noteID, err := primitive.ObjectIDFromHex(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
        return
    }

    note, err := nc.noteService.TogglePin(noteID, user.ID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, note)
}

func (nc *NoteController) GetHistory(c *gin.Context) {
    user := c.MustGet("user").(*models.User)
    
    noteID, err := primitive.ObjectIDFromHex(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
        return
    }

    versions, err := nc.noteService.GetVersionHistory(noteID, user.ID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, versions)
}

func (nc *NoteController) RestoreVersion(c *gin.Context) {
    user := c.MustGet("user").(*models.User)
    
    noteID, err := primitive.ObjectIDFromHex(c.Param("noteId"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
        return
    }

    versionID, err := primitive.ObjectIDFromHex(c.Param("versionId"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version ID"})
        return
    }

    note, err := nc.noteService.RestoreVersion(noteID, versionID, user.ID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, note)
}

func (nc *NoteController) FilterByTag(c *gin.Context) {
    user := c.MustGet("user").(*models.User)
    
    tag := c.Query("tag")
    if tag == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Tag parameter is required"})
        return
    }

    notes, err := nc.noteService.GetNotesByTag(tag, user.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, notes)
}

func (nc *NoteController) AutoSave(c *gin.Context) {
    user := c.MustGet("user").(*models.User)
    
    noteID, err := primitive.ObjectIDFromHex(c.Param("noteId"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
        return
    }

    var req models.NoteRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    note, err := nc.noteService.AutoSaveNote(noteID, user.ID, req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, note)
}

// ShareNote adds a collaborator (only owner)
func (nc *NoteController) ShareNote(c *gin.Context) {
    user := c.MustGet("user").(*models.User)
    noteID, err := primitive.ObjectIDFromHex(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
        return
    }
    var req models.AddCollaboratorRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    // Find collaborator user by username
    var collab models.User
    err = config.DB.Collection("users").FindOne(c, bson.M{"username": req.Username}).Decode(&collab)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Collaborator not found"})
        return
    }
    if collab.ID == user.ID {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot add yourself as collaborator"})
        return
    }
    err = nc.noteService.AddCollaborator(noteID, user.ID, collab.ID)
    if err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Collaborator added"})
}

// RemoveCollaborator removes a collaborator (only owner)
func (nc *NoteController) RemoveCollaborator(c *gin.Context) {
    user := c.MustGet("user").(*models.User)
    noteID, err := primitive.ObjectIDFromHex(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
        return
    }
    var req models.RemoveCollaboratorRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    // Find collaborator user by username
    var collab models.User
    err = config.DB.Collection("users").FindOne(c, bson.M{"username": req.Username}).Decode(&collab)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Collaborator not found"})
        return
    }
    err = nc.noteService.RemoveCollaborator(noteID, user.ID, collab.ID)
    if err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Collaborator removed"})
}

// ListCollaborators returns the list of collaborators for a note
func (nc *NoteController) ListCollaborators(c *gin.Context) {
    user := c.MustGet("user").(*models.User)
    noteID, err := primitive.ObjectIDFromHex(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
        return
    }
    collabs, err := nc.noteService.ListCollaborators(noteID, user.ID)
    if err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, collabs)
}
