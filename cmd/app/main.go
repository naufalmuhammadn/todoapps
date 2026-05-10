package main

import (
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"todoapps/internal/repository"
	"todoapps/internal/service"
)

func main() {
	connStr := "user=admin password=secret dbname=todo_db sslmode=disable host=localhost"
	repo, err := repository.NewPostgresRepo(connStr)
	if err != nil {
		log.Fatal(err)
	}

	h := &service.TaskHandler{Repo: repo}
	mux := http.NewServeMux()

	mux.HandleFunc("POST /tasks", h.CreateTask)
	mux.HandleFunc("GET /tasks", h.GetAllTasks)
	mux.HandleFunc("PUT /tasks/{id}", h.UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", h.DeleteTask)
	mux.HandleFunc("GET /tasks/{id}", h.GetTask)

	log.Println("Server running on :8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		return
	}
}
