package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	Priority    string    `json:"priority"` // "low", "medium", "high"
	DueDate     string    `json:"dueDate"`  // ISO 8601 format
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

var (
	todos    []Todo
	todosMu  sync.RWMutex
	nextID   int = 1
	nextIDMu sync.Mutex
)

func main() {
	// Initialize with sample data demonstrating all features
	todos = []Todo{
		{
			ID:          getNextID(),
			Title:       "Welcome to Enhanced TODO",
			Description: "This app now has priority levels, due dates, tags, and search!",
			Completed:   false,
			Priority:    "high",
			DueDate:     time.Now().AddDate(0, 0, 1).Format(time.RFC3339),
			Tags:        []string{"important", "demo"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          getNextID(),
			Title:       "Learn Go Backend",
			Description: "Complete the RESTful API with advanced features",
			Completed:   false,
			Priority:    "medium",
			DueDate:     time.Now().AddDate(0, 0, 3).Format(time.RFC3339),
			Tags:        []string{"learning", "backend"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          getNextID(),
			Title:       "Master Svelte Frontend",
			Description: "Build an interactive UI with filters and search",
			Completed:   true,
			Priority:    "low",
			DueDate:     time.Now().AddDate(0, 0, -1).Format(time.RFC3339),
			Tags:        []string{"learning", "frontend"},
			CreatedAt:   time.Now().Add(-48 * time.Hour),
			UpdatedAt:   time.Now(),
		},
	}

	http.HandleFunc("/", ToDoListHandler)
	http.HandleFunc("/toggle/", ToggleTodoHandler)
	http.HandleFunc("/delete/", DeleteTodoHandler)
	http.HandleFunc("/update/", UpdateTodoHandler)
	http.HandleFunc("/search", SearchTodosHandler)
	http.HandleFunc("/stats", StatsHandler)

	log.Println("🚀 Server starting on :8080")
	log.Println("📝 Endpoints available:")
	log.Println("   GET/POST  /           - List/Create todos")
	log.Println("   POST      /toggle/    - Toggle completion")
	log.Println("   DELETE    /delete/    - Delete todo")
	log.Println("   PUT       /update/    - Update todo")
	log.Println("   GET       /search     - Search todos")
	log.Println("   GET       /stats      - Get statistics")
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func getNextID() int {
	nextIDMu.Lock()
	defer nextIDMu.Unlock()
	id := nextID
	nextID++
	return id
}

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
}

func ToDoListHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

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

func ToggleTodoHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	todosMu.Lock()
	defer todosMu.Unlock()

	for i := range todos {
		if todos[i].ID == id {
			todos[i].Completed = !todos[i].Completed
			todos[i].UpdatedAt = time.Now()
			json.NewEncoder(w).Encode(todos[i])
			return
		}
	}

	http.Error(w, "Todo not found", http.StatusNotFound)
}

func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	todosMu.Lock()
	defer todosMu.Unlock()

	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Todo deleted successfully"})
			return
		}
	}

	http.Error(w, "Todo not found", http.StatusNotFound)
}

func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	var updates struct {
		Title       *string   `json:"title"`
		Description *string   `json:"description"`
		Priority    *string   `json:"priority"`
		DueDate     *string   `json:"dueDate"`
		Tags        *[]string `json:"tags"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	todosMu.Lock()
	defer todosMu.Unlock()

	for i := range todos {
		if todos[i].ID == id {
			if updates.Title != nil {
				todos[i].Title = *updates.Title
			}
			if updates.Description != nil {
				todos[i].Description = *updates.Description
			}
			if updates.Priority != nil {
				todos[i].Priority = *updates.Priority
			}
			if updates.DueDate != nil {
				todos[i].DueDate = *updates.DueDate
			}
			if updates.Tags != nil {
				todos[i].Tags = *updates.Tags
			}
			todos[i].UpdatedAt = time.Now()
			json.NewEncoder(w).Encode(todos[i])
			return
		}
	}

	http.Error(w, "Todo not found", http.StatusNotFound)
}

func SearchTodosHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := strings.ToLower(r.URL.Query().Get("q"))
	if query == "" {
		json.NewEncoder(w).Encode([]Todo{})
		return
	}

	todosMu.RLock()
	defer todosMu.RUnlock()

	var results []Todo
	for _, todo := range todos {
		if strings.Contains(strings.ToLower(todo.Title), query) ||
			strings.Contains(strings.ToLower(todo.Description), query) {
			results = append(results, todo)
		}
	}

	if results == nil {
		results = []Todo{}
	}

	json.NewEncoder(w).Encode(results)
}

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	todosMu.RLock()
	defer todosMu.RUnlock()

	stats := map[string]interface{}{
		"total":       len(todos),
		"active":      0,
		"completed":   0,
		"highPriority": 0,
		"overdue":     0,
		"byPriority": map[string]int{
			"high":   0,
			"medium": 0,
			"low":    0,
		},
	}

	now := time.Now()
	for _, todo := range todos {
		if todo.Completed {
			stats["completed"] = stats["completed"].(int) + 1
		} else {
			stats["active"] = stats["active"].(int) + 1
		}

		if todo.Priority == "high" {
			stats["highPriority"] = stats["highPriority"].(int) + 1
		}

		priorityMap := stats["byPriority"].(map[string]int)
		priorityMap[todo.Priority]++

		if !todo.Completed && todo.DueDate != "" {
			dueDate, err := time.Parse(time.RFC3339, todo.DueDate)
			if err == nil && dueDate.Before(now) {
				stats["overdue"] = stats["overdue"].(int) + 1
			}
		}
	}

	json.NewEncoder(w).Encode(stats)
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
	var input struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Priority    string   `json:"priority"`
		DueDate     string   `json:"dueDate"`
		Tags        []string `json:"tags"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if input.Title == "" || input.Description == "" {
		http.Error(w, "Title and description are required", http.StatusBadRequest)
		return
	}

	// Set default priority if not provided
	if input.Priority == "" {
		input.Priority = "medium"
	}

	// Validate priority
	if input.Priority != "low" && input.Priority != "medium" && input.Priority != "high" {
		http.Error(w, "Priority must be 'low', 'medium', or 'high'", http.StatusBadRequest)
		return
	}

	// Initialize empty tags array if nil
	if input.Tags == nil {
		input.Tags = []string{}
	}

	todo := Todo{
		ID:          getNextID(),
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
		DueDate:     input.DueDate,
		Tags:        input.Tags,
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
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