package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "time"

    "todo-app/db"
    "todo-app/models"
)

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
    var task models.Task
    err := json.NewDecoder(r.Body).Decode(&task)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    result, err := db.DB.Exec("INSERT INTO tasks (user_id, title, description, priority) VALUES (?, ?, ?, ?)", task.UserID, task.Title, task.Description, task.Priority)
    if err != nil {
        http.Error(w, "Could not create task", http.StatusInternalServerError)
        return
    }

    taskID, err := result.LastInsertId()
    if err != nil {
        http.Error(w, "Could not retrieve task ID", http.StatusInternalServerError)
        return
    }

    task.ID = int(taskID)
    task.CreatedAt = time.Now()
    task.UpdatedAt = time.Now()

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(task)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
    var task models.Task
    err := json.NewDecoder(r.Body).Decode(&task)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    _, err = db.DB.Exec("UPDATE tasks SET title = ?, description = ?, priority = ?, updated_at = NOW() WHERE id = ?", task.Title, task.Description, task.Priority, task.ID)
    if err != nil {
        http.Error(w, "Could not update task", http.StatusInternalServerError)
        return
    }

    task.UpdatedAt = time.Now()

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(task)
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("user_id")
    rows, err := db.DB.Query("SELECT id, user_id, title, description, priority, created_at, updated_at FROM tasks WHERE user_id = ?", userID)
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

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
    var task models.Task
    err := json.NewDecoder(r.Body).Decode(&task)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    _, err = db.DB.Exec("DELETE FROM tasks WHERE id = ?", task.ID)
    if err != nil {
        http.Error(w, "Could not delete task", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

