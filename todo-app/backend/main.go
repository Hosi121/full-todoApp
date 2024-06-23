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

    result, err := db.Exec("INSERT INTO tasks (user_id, title, description) VALUES (?, ?, ?)", task.UserID, task.Title, task.Description)
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

func main() {
    initDB()

    http.HandleFunc("/register", registerHandler)
    http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/tasks", createTaskHandler)
    http.HandleFunc("/tasks/list", getTasksHandler)

    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Could not start server: %s\n", err.Error())
    }
}

