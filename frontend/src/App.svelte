<script lang="ts">
  import Todo from "./lib/Todo.svelte";
  import type { TodoItem } from "./lib/types";

  let todos: TodoItem[] = $state([]);
  let filter: 'all' | 'active' | 'completed' = $state('all');

  const filteredTodos = $derived(() => {
    if (filter === 'active') return todos.filter(t => !t.completed);
    if (filter === 'completed') return todos.filter(t => t.completed);
    return todos;
  });

  const stats = $derived(() => ({
    total: todos.length,
    active: todos.filter(t => !t.completed).length,
    completed: todos.filter(t => t.completed).length,
  }));

  async function fetchTodos() {
    try {
      const response = await fetch("http://localhost:8080/");
      if (response.status !== 200) {
        console.error("Error fetching data. Response status not 200");
        return;
      }

      todos = await response.json();
    } catch (e) {
      console.error("Could not connect to server. Ensure it is running.", e);
    }
  }

  async function handleSubmit(event: SubmitEvent) {
    event.preventDefault();
    
    const formData = new FormData(event.target as HTMLFormElement);
    const title = formData.get("title") as string;
    const description = formData.get("description") as string;

    if (!title || !description) {
      alert("Please fill in both title and description");
      return;
    }

    try {
      const response = await fetch("http://localhost:8080/", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ title, description }),
      });

      if (response.status !== 200) {
        console.error("Error creating todo. Response status not 200");
        return;
      }

      // Clear the form
      (event.target as HTMLFormElement).reset();

      // Refresh the todo list
      await fetchTodos();
    } catch (e) {
      console.error("Could not connect to server. Ensure it is running.", e);
    }
  }

  async function toggleTodo(id: number) {
    try {
      const response = await fetch(`http://localhost:8080/toggle/?id=${id}`, {
        method: "POST",
      });

      if (response.status !== 200) {
        console.error("Error toggling todo. Response status not 200");
        return;
      }

      await fetchTodos();
    } catch (e) {
      console.error("Could not connect to server. Ensure it is running.", e);
    }
  }

  async function deleteTodo(id: number) {
    if (!confirm("Are you sure you want to delete this todo?")) {
      return;
    }

    try {
      const response = await fetch(`http://localhost:8080/delete/?id=${id}`, {
        method: "DELETE",
      });

      if (response.status !== 200) {
        console.error("Error deleting todo. Response status not 200");
        return;
      }

      await fetchTodos();
    } catch (e) {
      console.error("Could not connect to server. Ensure it is running.", e);
    }
  }

  // Initially fetch todos on page load
  $effect(() => {
    fetchTodos();
  });
</script>

<main class="app">
  <header class="app-header">
    <h1>TODO</h1>
    <div class="stats">
      <span class="stat">Total: {stats().total}</span>
      <span class="stat">Active: {stats().active}</span>
      <span class="stat">Completed: {stats().completed}</span>
    </div>
  </header>

  <div class="filters">
    <button 
      class="filter-btn" 
      class:active={filter === 'all'}
      onclick={() => filter = 'all'}
    >
      All
    </button>
    <button 
      class="filter-btn" 
      class:active={filter === 'active'}
      onclick={() => filter = 'active'}
    >
      Active
    </button>
    <button 
      class="filter-btn" 
      class:active={filter === 'completed'}
      onclick={() => filter = 'completed'}
    >
      Completed
    </button>
  </div>

  <div class="todo-list">
    {#each filteredTodos() as todo (todo.id)}
      <Todo 
        {...todo} 
        onToggle={toggleTodo}
        onDelete={deleteTodo}
      />
    {:else}
      <p class="empty-message">No todos to display</p>
    {/each}
  </div>

  <h2 class="todo-list-form-header">Add a Todo</h2>
  <form class="todo-list-form" onsubmit={handleSubmit}>
    <input placeholder="Title" name="title" required />
    <input placeholder="Description" name="description" required />
    <button type="submit">Add Todo</button>
  </form>
</main>

<style>
  .app {
    color: white;
    background-color: #282c34;

    text-align: center;
    font-size: 24px;

    min-height: 100vh;
    padding: 20px;
  }

  .app-header {
    font-size: calc(10px + 4vmin);
    margin-top: 50px;
  }

  .stats {
    display: flex;
    justify-content: center;
    gap: 30px;
    margin-top: 20px;
    font-size: 18px;
  }

  .stat {
    background-color: #303540;
    padding: 8px 16px;
    border-radius: 8px;
    border: 1px solid #464e61;
  }

  .filters {
    display: flex;
    justify-content: center;
    gap: 15px;
    margin: 30px 0;
  }

  .filter-btn {
    background-color: #303540;
    border-color: #464e61;
    font-size: 18px;
    padding: 8px 20px;
  }

  .filter-btn.active {
    background-color: #326ac7;
    border-color: #779ad4;
  }

  .todo-list {
    margin: 50px 100px 0px 100px;
  }

  .empty-message {
    color: #888;
    font-size: 20px;
    margin-top: 40px;
  }

  .todo-list-form-header {
    margin-top: 100px;
  }

  .todo-list-form {
    margin-top: 10px;
  }
</style>