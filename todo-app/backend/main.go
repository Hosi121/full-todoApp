package main

import (
    "log"
    "net/http"
    "backend/db"
    "backend/handlers"
)

func main() {
    db.InitDB()

    http.HandleFunc("/register", handlers.RegisterHandler)
    http.HandleFunc("/login", handlers.LoginHandler)
    http.HandleFunc("/tasks", handlers.CreateTaskHandler)
    http.HandleFunc("/tasks/list", handlers.GetTasksHandler)
    http.HandleFunc("/tasks/update", handlers.UpdateTaskHandler)
    http.HandleFunc("/tasks/delete", handlers.DeleteTaskHandler)
    http.HandleFunc("/tags", handlers.CreateTagHandler)
    http.HandleFunc("/tasks/add-tag", handlers.AddTagToTaskHandler)
    http.HandleFunc("/tasks/filter", handlers.FilterTasksHandler)

    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Could not start server: %s\n", err.Error())
    }
}

