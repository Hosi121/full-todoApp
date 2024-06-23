package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"

    "golang.org/x/crypto/bcrypt"
    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB() {
    var err error
    db, err = sql.Open("mysql", "user:password@tcp(db:3306)/todo_app")
    if err != nil {
        log.Fatalf("Could not connect to the database: %v", err)
    }

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        username VARCHAR(50) NOT NULL UNIQUE,
        password_hash VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`)
    if err != nil {
        log.Fatalf("Could not create users table: %v", err)
    }

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
        id INT AUTO_INCREMENT PRIMARY KEY,
        user_id INT NOT NULL,
        title VARCHAR(100) NOT NULL,
        description TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id)
    )`)
    if err != nil {
        log.Fatalf("Could not create tasks table: %v", err)
    }
}

type Task struct {
    ID          int       `json:"id"`
    UserID      int       `json:"user_id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Priority    int       `json:"priority"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
    var task Task
    err := json.NewDecoder(r.Body).Decode(&task)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    result, err := db.Exec("INSERT INTO tasks (user_id, title, description, priority) VALUES (?, ?, ?, ?)", task.UserID, task.Title, task.Description, task.Priority)
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
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
    var task Task
    err := json.NewDecoder(r.Body).Decode(&task)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    _, err = db.Exec("UPDATE tasks SET title = ?, description = ?, priority = ?, updated_at = NOW() WHERE id = ?", task.Title, task.Description, task.Priority, task.ID)
    if err != nil {
        http.Error(w, "Could not update task", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(task)
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("user_id")
    rows, err := db.Query("SELECT id, user_id, title, description, created_at, updated_at FROM tasks WHERE user_id = ?", userID)
    if err != nil {
        http.Error(w, "Could not fetch tasks", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    tasks := []Task{}
    for rows.Next() {
        var task Task
        if err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.CreatedAt, &task.UpdatedAt); err != nil {
            http.Error(w, "Could not scan task", http.StatusInternalServerError)
            return
        }
        tasks = append(tasks, task)
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(tasks)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
    var task Task
    err := json.NewDecoder(r.Body).Decode(&task)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    _, err = db.Exec("DELETE FROM tasks WHERE id = ?", task.ID)
    if err != nil {
        http.Error(w, "Could not delete task", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

type Tag struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func createTagHandler(w http.ResponseWriter, r *http.Request) {
    var tag Tag
    err := json.NewDecoder(r.Body).Decode(&tag)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    result, err := db.Exec("INSERT INTO tags (name) VALUES (?)", tag.Name)
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

func addTagToTaskHandler(w http.ResponseWriter, r *http.Request) {
    var payload struct {
        TaskID int `json:"task_id"`
        TagID  int `json:"tag_id"`
    }
    err := json.NewDecoder(r.Body).Decode(&payload)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    _, err = db.Exec("INSERT INTO task_tags (task_id, tag_id) VALUES (?, ?)", payload.TaskID, payload.TagID)
    if err != nil {
        http.Error(w, "Could not add tag to task", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func filterTasksHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("user_id")
    tagID := r.URL.Query().Get("tag_id")
    date := r.URL.Query().Get("date")

    query := "SELECT tasks.id, tasks.user_id, tasks.title, tasks.description, tasks.created_at, tasks.updated_at FROM tasks"
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

    rows, err := db.Query(query, args...)
    if err != nil {
        http.Error(w, "Could not fetch tasks", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    tasks := []Task{}
    for rows.Next() {
        var task Task
        if err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.CreatedAt, &task.UpdatedAt); err != nil {
            http.Error(w, "Could not scan task", http.StatusInternalServerError)
            return
        }
        tasks = append(tasks, task)
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(tasks)
}

func main() {
    initDB()

    http.HandleFunc("/register", registerHandler)
    http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/tasks", createTaskHandler)
    http.HandleFunc("/tasks/list", getTasksHandler)
    http.HandleFunc("/tasks/update", updateTaskHandler)
    http.HandleFunc("/tasks/delete", deleteTaskHandler)
    http.HandleFunc("/tags", createTagHandler)
    http.HandleFunc("/tasks/add-tag", addTagToTaskHandler)
    http.HandleFunc("/tasks/filter", filterTasksHandler)

    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Could not start server: %s\n", err.Error())
    }
}
