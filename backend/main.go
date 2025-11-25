package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Todo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

var (
	todos   []Todo
	todosMu sync.RWMutex
)

func main() {
	// Initialize with some sample data
	todos = []Todo{
		{Title: "Welcome", Description: "This is your first todo item"},
		{Title: "Learn Go", Description: "Complete the backend implementation"},
		{Title: "Learn Svelte", Description: "Complete the frontend form"},
	}

	http.HandleFunc("/", ToDoListHandler)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func ToDoListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case "GET":
		handleGetTodos(w, r)
	case "POST":
		handleCreateTodo(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetTodos(w http.ResponseWriter, r *http.Request) {
	todosMu.RLock()
	defer todosMu.RUnlock()

	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if todo.Title == "" || todo.Description == "" {
		http.Error(w, "Title and description are required", http.StatusBadRequest)
		return
	}

	todosMu.Lock()
	todos = append(todos, todo)
	todosMu.Unlock()

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}