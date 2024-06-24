package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "time"

    "golang.org/x/crypto/bcrypt"
    "backend/db"
)

type Credentials struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    enableCors(&w) // 追加
    var creds Credentials
    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Could not hash password", http.StatusInternalServerError)
        return
    }

    _, err = db.DB.Exec("INSERT INTO users (username, password_hash, created_at) VALUES (?, ?, ?)", creds.Username, string(hashedPassword), time.Now())
    if err != nil {
        http.Error(w, "Could not register user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    enableCors(&w) // 追加
    var creds Credentials
    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    var storedCreds Credentials
    row := db.DB.QueryRow("SELECT username, password_hash FROM users WHERE username = ?", creds.Username)
    err = row.Scan(&storedCreds.Username, &storedCreds.Password)
    if err == sql.ErrNoRows || bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)) != nil {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    } else if err != nil {
        http.Error(w, "Could not fetch user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func enableCors(w *http.ResponseWriter) {
    (*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    (*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
