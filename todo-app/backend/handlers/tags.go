package handlers

import (
    "encoding/json"
    "net/http"

    "backend/db"
    "backend/models"
)

func CreateTagHandler(w http.ResponseWriter, r *http.Request) {
    var tag models.Tag
    err := json.NewDecoder(r.Body).Decode(&tag)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    result, err := db.DB.Exec("INSERT INTO tags (name) VALUES (?)", tag.Name)
    if err != nil {
        http.Error(w, "Could not create tag", http.StatusInternalServerError)
        return
    }

    tagID, err := result.LastInsertId()
    if err != nil {
        http.Error(w, "Could not retrieve tag ID", http.StatusInternalServerError)
        return
    }

    tag.ID = int(tagID)
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(tag)
}

func AddTagToTaskHandler(w http.ResponseWriter, r *http.Request) {
    var payload struct {
        TaskID int `json:"task_id"`
        TagID  int `json:"tag_id"`
    }
    err := json.NewDecoder(r.Body).Decode(&payload)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    _, err = db.DB.Exec("INSERT INTO task_tags (task_id, tag_id) VALUES (?, ?)", payload.TaskID, payload.TagID)
    if err != nil {
        http.Error(w, "Could not add tag to task", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

