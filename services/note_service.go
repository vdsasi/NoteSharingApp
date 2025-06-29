package services

import ( 
    "context"
    "errors"
    "time"
    "notes-app/config"
    "notes-app/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type NoteService struct{}

func NewNoteService() *NoteService {
    return &NoteService{}
}

func (s *NoteService) CreateNote(userID primitive.ObjectID, req models.NoteRequest) (models.NoteResponse, error) {
    ctx := context.Background()
    
    note := models.Note{
        ID:              primitive.NewObjectID(),
        Title:           req.Title,
        Content:         req.Content,
        UserID:          userID,
        Tags:            req.Tags,
        AutoSaveEnabled: req.AutoSaveEnabled,
        CreatedAt:       time.Now(),
        UpdatedAt:       time.Now(),
    }

    _, err := config.DB.Collection("notes").InsertOne(ctx, note)
    if err != nil {
        return models.NoteResponse{}, err
    }

    return s.noteToResponse(note), nil
}

func (s *NoteService) GetUserNotes(userID primitive.ObjectID) ([]models.NoteResponse, error) {
    ctx := context.Background()
    
    filter := bson.M{"userId": userID, "trashed": false}
    opts := options.Find().SetSort(bson.D{{"pinned", -1}, {"updatedAt", -1}})
    
    cursor, err := config.DB.Collection("notes").Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var notes []models.Note
    if err = cursor.All(ctx, &notes); err != nil {
        return nil, err
    }

    var responses []models.NoteResponse
    for _, note := range notes {
        responses = append(responses, s.noteToResponse(note))
    }

    return responses, nil
}

func (s *NoteService) GetNote(noteID, userID primitive.ObjectID) (models.NoteResponse, error) {
    ctx := context.Background()
    
    var note models.Note
    err := config.DB.Collection("notes").FindOne(ctx, bson.M{
        "_id": noteID,
        "$or": []bson.M{
            {"userId": userID},
            {"collaborators": userID},
        },
    }).Decode(&note)
    
    if err != nil {
        return models.NoteResponse{}, errors.New("note not found or access denied")
    }

    return s.noteToResponse(note), nil
}

func (s *NoteService) UpdateNote(noteID, userID primitive.ObjectID, req models.NoteRequest) (models.NoteResponse, error) {
    ctx := context.Background()
    
    // Check if user is owner or collaborator
    var currentNote models.Note
    err := config.DB.Collection("notes").FindOne(ctx, bson.M{
        "_id": noteID,
        "$or": []bson.M{
            {"userId": userID},
            {"collaborators": userID},
        },
    }).Decode(&currentNote)
    
    if err != nil {
        return models.NoteResponse{}, errors.New("note not found or access denied")
    }

    // Save version
    version := models.NoteVersion{
        ID:          primitive.NewObjectID(),
        NoteID:      noteID,
        Title:       currentNote.Title,
        Content:     currentNote.Content,
        VersionedAt: time.Now(),
    }
    config.DB.Collection("note_versions").InsertOne(ctx, version)

    // Update note
    update := bson.M{
        "$set": bson.M{
            "title":           req.Title,
            "content":         req.Content,
            "tags":            req.Tags,
            "autoSaveEnabled": req.AutoSaveEnabled,
            "updatedAt":       time.Now(),
        },
    }

    result, err := config.DB.Collection("notes").UpdateOne(ctx, bson.M{
        "_id": noteID,
        "$or": []bson.M{
            {"userId": userID},
            {"collaborators": userID},
        },
    }, update)

    if err != nil || result.MatchedCount == 0 {
        return models.NoteResponse{}, errors.New("failed to update note or access denied")
    }

    // Return updated note
    return s.GetNote(noteID, userID)
}

func (s *NoteService) DeleteNote(noteID, userID primitive.ObjectID) error {
    ctx := context.Background()
    
    result, err := config.DB.Collection("notes").UpdateOne(ctx, bson.M{
        "_id": noteID,
        "userId": userID,
    }, bson.M{
        "$set": bson.M{"trashed": true, "updatedAt": time.Now()},
    })

    if err != nil || result.MatchedCount == 0 {
        return errors.New("note not found")
    }

    return nil
}

func (s *NoteService) GetTrashedNotes(userID primitive.ObjectID) ([]models.NoteResponse, error) {
    ctx := context.Background()
    
    filter := bson.M{"userId": userID, "trashed": true}
    opts := options.Find().SetSort(bson.D{{"updatedAt", -1}})
    
    cursor, err := config.DB.Collection("notes").Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var notes []models.Note
    if err = cursor.All(ctx, &notes); err != nil {
        return nil, err
    }

    var responses []models.NoteResponse
    for _, note := range notes {
        responses = append(responses, s.noteToResponse(note))
    }

    return responses, nil
}

func (s *NoteService) RestoreNote(noteID, userID primitive.ObjectID) (models.NoteResponse, error) {
    ctx := context.Background()
    
    result, err := config.DB.Collection("notes").UpdateOne(ctx, bson.M{
        "_id": noteID,
        "userId": userID,
    }, bson.M{
        "$set": bson.M{"trashed": false, "updatedAt": time.Now()},
    })

    if err != nil || result.MatchedCount == 0 {
        return models.NoteResponse{}, errors.New("note not found")
    }

    return s.GetNote(noteID, userID)
}

func (s *NoteService) TogglePin(noteID, userID primitive.ObjectID) (models.NoteResponse, error) {
    ctx := context.Background()
    
    var note models.Note
    err := config.DB.Collection("notes").FindOne(ctx, bson.M{
        "_id": noteID,
        "userId": userID,
    }).Decode(&note)
    
    if err != nil {
        return models.NoteResponse{}, errors.New("note not found")
    }

    newPinnedState := !note.Pinned
    
    _, err = config.DB.Collection("notes").UpdateOne(ctx, bson.M{
        "_id": noteID,
        "userId": userID,
    }, bson.M{
        "$set": bson.M{"pinned": newPinnedState, "updatedAt": time.Now()},
    })

    if err != nil {
        return models.NoteResponse{}, err
    }

    return s.GetNote(noteID, userID)
}

func (s *NoteService) GetVersionHistory(noteID, userID primitive.ObjectID) ([]models.NoteVersionResponse, error) {
    ctx := context.Background()
    
    // Verify note ownership
    var note models.Note
    err := config.DB.Collection("notes").FindOne(ctx, bson.M{
        "_id": noteID,
        "userId": userID,
    }).Decode(&note)
    
    if err != nil {
        return nil, errors.New("note not found")
    }

    // Get versions
    opts := options.Find().SetSort(bson.D{{"versionedAt", -1}})
    cursor, err := config.DB.Collection("note_versions").Find(ctx, bson.M{"noteId": noteID}, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var versions []models.NoteVersion
    if err = cursor.All(ctx, &versions); err != nil {
        return nil, err
    }

    var responses []models.NoteVersionResponse
    for _, version := range versions {
        responses = append(responses, models.NoteVersionResponse{
            ID:          version.ID.Hex(),
            Title:       version.Title,
            Content:     version.Content,
            VersionedAt: version.VersionedAt,
        })
    }

    return responses, nil
}

func (s *NoteService) RestoreVersion(noteID, versionID, userID primitive.ObjectID) (models.NoteResponse, error) {
    ctx := context.Background()
    
    // Verify note ownership
    var note models.Note
    err := config.DB.Collection("notes").FindOne(ctx, bson.M{
        "_id": noteID,
        "userId": userID,
    }).Decode(&note)
    
    if err != nil {
        return models.NoteResponse{}, errors.New("note not found")
    }

    // Get version
    var version models.NoteVersion
    err = config.DB.Collection("note_versions").FindOne(ctx, bson.M{
        "_id": versionID,
        "noteId": noteID,
    }).Decode(&version)
    
    if err != nil {
        return models.NoteResponse{}, errors.New("version not found")
    }

    // Create current version before restoring
    currentVersion := models.NoteVersion{
        ID:          primitive.NewObjectID(),
        NoteID:      noteID,
        Title:       note.Title,
        Content:     note.Content,
        VersionedAt: time.Now(),
    }
    config.DB.Collection("note_versions").InsertOne(ctx, currentVersion)

    // Restore version
    _, err = config.DB.Collection("notes").UpdateOne(ctx, bson.M{
        "_id": noteID,
    }, bson.M{
        "$set": bson.M{
            "title":     version.Title,
            "content":   version.Content,
            "updatedAt": time.Now(),
        },
    })

    if err != nil {
        return models.NoteResponse{}, err
    }

    return s.GetNote(noteID, userID)
}

func (s *NoteService) GetNotesByTag(tag string, userID primitive.ObjectID) ([]models.NoteResponse, error) {
    ctx := context.Background()
    
    filter := bson.M{
        "userId": userID,
        "trashed": false,
        "tags": bson.M{"$in": []string{tag}},
    }
    opts := options.Find().SetSort(bson.D{{"updatedAt", -1}})
    
    cursor, err := config.DB.Collection("notes").Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var notes []models.Note
    if err = cursor.All(ctx, &notes); err != nil {
        return nil, err
    }

    var responses []models.NoteResponse
    for _, note := range notes {
        responses = append(responses, s.noteToResponse(note))
    }

    return responses, nil
}

func (s *NoteService) AutoSaveNote(noteID, userID primitive.ObjectID, req models.NoteRequest) (models.NoteResponse, error) {
    ctx := context.Background()
    
    // Update without creating version for autosave
    update := bson.M{
        "$set": bson.M{
            "title":     req.Title,
            "content":   req.Content,
            "tags":      req.Tags,
            "updatedAt": time.Now(),
        },
    }

    result, err := config.DB.Collection("notes").UpdateOne(ctx, bson.M{
        "_id": noteID,
        "userId": userID,
    }, update)

    if err != nil || result.MatchedCount == 0 {
        return models.NoteResponse{}, errors.New("failed to autosave note")
    }

    return s.GetNote(noteID, userID)
}

func (s *NoteService) AddCollaborator(noteID, ownerID, collaboratorID primitive.ObjectID) error {
    ctx := context.Background()
    // Only owner can add
    result := config.DB.Collection("notes").FindOneAndUpdate(ctx, bson.M{
        "_id": noteID,
        "userId": ownerID,
    }, bson.M{
        "$addToSet": bson.M{"collaborators": collaboratorID},
    })
    if result.Err() != nil {
        return errors.New("not found or not owner")
    }
    return nil
}

// RemoveCollaborator removes a collaborator from a note (only owner can do this)
func (s *NoteService) RemoveCollaborator(noteID, ownerID, collaboratorID primitive.ObjectID) error {
    ctx := context.Background()
    result := config.DB.Collection("notes").FindOneAndUpdate(ctx, bson.M{
        "_id": noteID,
        "userId": ownerID,
    }, bson.M{
        "$pull": bson.M{"collaborators": collaboratorID},
    })
    if result.Err() != nil {
        return errors.New("not found or not owner")
    }
    return nil
}

// ListCollaborators returns the list of collaborators for a note
func (s *NoteService) ListCollaborators(noteID, userID primitive.ObjectID) ([]models.UserProfileDto, error) {
    ctx := context.Background()
    var note models.Note
    err := config.DB.Collection("notes").FindOne(ctx, bson.M{
        "_id": noteID,
        "$or": []bson.M{
            {"userId": userID},
            {"collaborators": userID},
        },
    }).Decode(&note)
    if err != nil {
        return nil, errors.New("note not found or access denied")
    }
    if len(note.Collaborators) == 0 {
        return []models.UserProfileDto{}, nil
    }
    // Fetch user info for each collaborator
    var users []models.User
    cursor, err := config.DB.Collection("users").Find(ctx, bson.M{"_id": bson.M{"$in": note.Collaborators}})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    if err = cursor.All(ctx, &users); err != nil {
        return nil, err
    }
    var result []models.UserProfileDto
    for _, u := range users {
        result = append(result, models.UserProfileDto{
            ID: u.ID.Hex(),
            Username: u.Username,
            Email: u.Email,
        })
    }
    return result, nil
}

func (s *NoteService) noteToResponse(note models.Note) models.NoteResponse {
    return models.NoteResponse{
        ID:              note.ID.Hex(),
        Title:           note.Title,
        Content:         note.Content,
        Pinned:          note.Pinned,
        Trashed:         note.Trashed,
        AutoSaveEnabled: note.AutoSaveEnabled,
        Tags:            note.Tags,
        CreatedAt:       note.CreatedAt,
        UpdatedAt:       note.UpdatedAt,
        UserID:          note.UserID.Hex(), // Add this line
    }
}
