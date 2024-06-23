package db

import (
    "database/sql"
    "log"
    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
    var err error
    DB, err = sql.Open("mysql", "user:password@tcp(db:3306)/todo_app")
    if err != nil {
        log.Fatalf("Could not connect to the database: %v", err)
    }

    _, err = DB.Exec(`CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        username VARCHAR(50) NOT NULL UNIQUE,
        password_hash VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`)
    if err != nil {
        log.Fatalf("Could not create users table: %v", err)
    }

    _, err = DB.Exec(`CREATE TABLE IF NOT EXISTS tasks (
        id INT AUTO_INCREMENT PRIMARY KEY,
        user_id INT NOT NULL,
        title VARCHAR(100) NOT NULL,
        description TEXT,
        priority INT DEFAULT 0,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id)
    )`)
    if err != nil {
        log.Fatalf("Could not create tasks table: %v", err)
    }

    _, err = DB.Exec(`CREATE TABLE IF NOT EXISTS tags (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(50) NOT NULL UNIQUE
    )`)
    if err != nil {
        log.Fatalf("Could not create tags table: %v", err)
    }

    _, err = DB.Exec(`CREATE TABLE IF NOT EXISTS task_tags (
        task_id INT NOT NULL,
        tag_id INT NOT NULL,
        FOREIGN KEY (task_id) REFERENCES tasks(id),
        FOREIGN KEY (tag_id) REFERENCES tags(id),
        PRIMARY KEY (task_id, tag_id)
    )`)
    if err != nil {
        log.Fatalf("Could not create task_tags table: %v", err)
    }
}

