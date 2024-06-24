# Todo App

This is a simple Todo application built with Go, React, MySQL, and Docker. The application allows users to register, log in, and manage their tasks with the ability to add tags and set priorities.

## Project Structure

backend/
├── db/
│ ├── init.go
│ └── queries.go
├── handlers/
│ ├── auth.go
│ ├── tasks.go
│ ├── tags.go
│ └── filter.go
├── models/
│ ├── task.go
│ └── tag.go
├── go.mod
├── go.sum
└── main.go
frontend/
├── src/
│ ├── components/
│ │ ├── Register.tsx
│ │ ├── CreateTask.tsx
│ │ ├── TaskList.tsx
│ ├── setupTests.ts
│ ├── App.tsx
│ ├── index.tsx
├── test/
│ ├── Register.test.tsx
│ ├── CreateTask.test.tsx
│ ├── TaskList.test.tsx
├── jest.config.ts
├── package.json
├── tsconfig.json
docker-compose.yml


## Prerequisites

- Docker
- Docker Compose
- Node.js (for frontend development)
- Go (for backend development)

## Getting Started

### Clone the repository

```bash
git clone https://github.com/your-username/todo-app.git
cd todo-app
```

This will build the Docker images for the backend and frontend services and start the application.

## Runnnig Tests

To run the frontend tests, navigate to the frontend directory and run:

```bash
npm test
```

## Development

### Frontend
run:

```bash
npm install
npm run dev
```

### Backend
run:
```bash
go run main.go
```

