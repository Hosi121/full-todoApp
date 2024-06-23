package handlers

import (
    "encoding/json"
    "net/http"

    "todo-app/db"
    "todo-app/models"
)

func FilterTasksHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("user_id")
    tagID := r.URL.Query().Get("tag_id")
    date := r.URL.Query().Get("date")

    query := "SELECT tasks.id, tasks.user_id, tasks.title, tasks.description, tasks.priority, tasks.created_at, tasks.updated_at FROM tasks"
    var args []interface{}

    if tagID != "" {
        query += " JOIN task_tags ON tasks.id = task_tags.task_id WHERE task_tags.tag_id = ?"
        args = append(args, tagID)
    } else {
        query += " WHERE 1=1"
    }

    if userID != "" {
        query += " AND tasks.user_id = ?"
        args = append(args, userID)
    }

    if date != "" {
        query += " AND DATE(tasks.created_at) = ?"
        args = append(args, date)
    }

    rows, err := db.DB.Query(query, args...)
    if err != nil {
        http.Error(w, "Could not fetch tasks", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    tasks := []models.Task{}
    for rows.Next() {
        var task models.Task
        if err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Priority, &task.CreatedAt, &task.UpdatedAt); err != nil {
            http.Error(w, "Could not scan task", http.StatusInternalServerError)
            return
        }
        tasks = append(tasks, task)
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(tasks)
}

