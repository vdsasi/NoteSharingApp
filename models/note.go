package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Title           string             `bson:"title" json:"title"`
    Content         string             `bson:"content" json:"content"`
    Pinned          bool               `bson:"pinned" json:"pinned"`
    Trashed         bool               `bson:"trashed" json:"trashed"`
    AutoSaveEnabled bool               `bson:"autoSaveEnabled" json:"autoSaveEnabled"`
    UserID          primitive.ObjectID `bson:"userId" json:"userId"`
    Tags            []string           `bson:"tags" json:"tags"`
    CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
    UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type NoteRequest struct {
    Title           string   `json:"title"`
    Content         string   `json:"content"`
    Tags            []string `json:"tags"`
    AutoSaveEnabled bool     `json:"autoSaveEnabled"`
}

type NoteResponse struct {
    ID              string    `json:"id"`
    Title           string    `json:"title"`
    Content         string    `json:"content"`
    Pinned          bool      `json:"pinned"`
    Trashed         bool      `json:"trashed"`
    AutoSaveEnabled bool      `json:"autoSaveEnabled"`
    Tags            []string  `json:"tags"`
    CreatedAt       time.Time `json:"createdAt"`
    UpdatedAt       time.Time `json:"updatedAt"`
}

type NoteVersion struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    NoteID      primitive.ObjectID `bson:"noteId" json:"noteId"`
    Title       string             `bson:"title" json:"title"`
    Content     string             `bson:"content" json:"content"`
    VersionedAt time.Time          `bson:"versionedAt" json:"versionedAt"`
}

type NoteVersionResponse struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Content     string    `json:"content"`
    VersionedAt time.Time `json:"versionedAt"`
}
