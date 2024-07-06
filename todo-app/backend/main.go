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
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
    var creds struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Could not hash password", http.StatusInternalServerError)
        return
    }

    _, err = db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", creds.Username, hash)
    if err != nil {
        http.Error(w, "Could not register user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    var creds struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    var storedHash string
    err = db.QueryRow("SELECT password_hash FROM users WHERE username = ?", creds.Username).Scan(&storedHash)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        } else {
            http.Error(w, "Could not fetch user", http.StatusInternalServerError)
        }
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(creds.Password))
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
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

    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Could not start server: %s\n", err.Error())
    }
}
