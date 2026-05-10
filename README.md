# Todo CRUD API

## Problem Description

A simple RESTful CRUD API for managing tasks, built with Go and PostgreSQL. No frameworks — just the standard library and a Postgres driver.

## Dev Environment Setup

- Go 1.22+
- PostgreSQL

### Database Setup

Create a PostgreSQL database and the tasks table:

```sql
CREATE DATABASE todo_db;

CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    task TEXT NOT NULL,
    done BOOLEAN DEFAULT FALSE
);
```

## How to Run

```bash
go run ./cmd/app/main.go
```

The server starts on `http://localhost:8080`.

> By default it connects with `user=admin password=secret dbname=todo_db host=localhost`. Update the connection string in `main.go` to match your setup.

## How to Build

```bash
go build -o todo-app ./cmd/app/main.go
./todo-app
```

## How to Test

```bash
go test -v ./...
```

## Usage Examples

### Create a task

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"task": "Buy groceries", "done": false}'
```

**Response:** `201 Created`

### Get all tasks

```bash
curl http://localhost:8080/tasks
```

**Response:**

```json
[
  {"id": 1, "task": "Buy groceries", "done": false}
]
```

### Get a single task

```bash
curl http://localhost:8080/tasks/1
```

**Response:**

```json
{"id": 1, "task": "Buy groceries", "done": false}
```

### Update a task

```bash
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"task": "Buy groceries", "done": true}'
```

**Response:** `200 OK`

### Delete a task

```bash
curl -X DELETE http://localhost:8080/tasks/1
```

**Response:** `200 OK`
