<script lang="ts">
  import Todo from "./lib/Todo.svelte";
  import type { TodoItem, FilterType, SortType, Stats } from "./lib/types";

  let todos: TodoItem[] = $state([]);
  let filter: FilterType = $state('all');
  let sortBy: SortType = $state('created');
  let searchQuery: string = $state('');
  let stats: Stats | null = $state(null);
  let showAddForm: boolean = $state(false);
  let editingTodo: TodoItem | null = $state(null);

  // Form fields
  let formTitle: string = $state('');
  let formDescription: string = $state('');
  let formPriority: 'low' | 'medium' | 'high' = $state('medium');
  let formDueDate: string = $state('');
  let formTags: string = $state('');

  const filteredAndSortedTodos = $derived(() => {
    let filtered = todos;

    // Apply filter
    if (filter === 'active') filtered = todos.filter(t => !t.completed);
    if (filter === 'completed') filtered = todos.filter(t => t.completed);

    // Apply search
    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(t => 
        t.title.toLowerCase().includes(query) ||
        t.description.toLowerCase().includes(query) ||
        t.tags.some(tag => tag.toLowerCase().includes(query))
      );
    }

    // Apply sort
    const sorted = [...filtered];
    switch (sortBy) {
      case 'created':
        sorted.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime());
        break;
      case 'dueDate':
        sorted.sort((a, b) => {
          if (!a.dueDate) return 1;
          if (!b.dueDate) return -1;
          return new Date(a.dueDate).getTime() - new Date(b.dueDate).getTime();
        });
        break;
      case 'priority':
        const priorityOrder = { high: 0, medium: 1, low: 2 };
        sorted.sort((a, b) => priorityOrder[a.priority] - priorityOrder[b.priority]);
        break;
      case 'title':
        sorted.sort((a, b) => a.title.localeCompare(b.title));
        break;
    }

    return sorted;
  });

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

  async function fetchStats() {
    try {
      const response = await fetch("http://localhost:8080/stats");
      if (response.status !== 200) {
        console.error("Error fetching stats. Response status not 200");
        return;
      }
      stats = await response.json();
    } catch (e) {
      console.error("Could not fetch stats", e);
    }
  }

  async function handleSubmit(event: SubmitEvent) {
    event.preventDefault();

    if (!formTitle || !formDescription) {
      alert("Please fill in both title and description");
      return;
    }

    const tagsArray = formTags
      .split(',')
      .map(tag => tag.trim())
      .filter(tag => tag.length > 0);

    const todoData = {
      title: formTitle,
      description: formDescription,
      priority: formPriority,
      dueDate: formDueDate || '',
      tags: tagsArray,
    };

    try {
      if (editingTodo) {
        // Update existing todo
        const response = await fetch(`http://localhost:8080/update/?id=${editingTodo.id}`, {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(todoData),
        });

        if (response.status !== 200) {
          console.error("Error updating todo. Response status not 200");
          return;
        }
      } else {
        // Create new todo
        const response = await fetch("http://localhost:8080/", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(todoData),
        });

        if (response.status !== 200) {
          console.error("Error creating todo. Response status not 200");
          return;
        }
      }

      // Clear form and refresh
      resetForm();
      await fetchTodos();
      await fetchStats();
    } catch (e) {
      console.error("Could not connect to server. Ensure it is running.", e);
    }
  }

  function resetForm() {
    formTitle = '';
    formDescription = '';
    formPriority = 'medium';
    formDueDate = '';
    formTags = '';
    editingTodo = null;
    showAddForm = false;
  }

  function openAddForm() {
    resetForm();
    showAddForm = true;
  }

  function openEditForm(id: number) {
    const todo = todos.find(t => t.id === id);
    if (!todo) return;

    editingTodo = todo;
    formTitle = todo.title;
    formDescription = todo.description;
    formPriority = todo.priority;
    formDueDate = todo.dueDate ? todo.dueDate.split('T')[0] : '';
    formTags = todo.tags.join(', ');
    showAddForm = true;
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
      await fetchStats();
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
      await fetchStats();
    } catch (e) {
      console.error("Could not connect to server. Ensure it is running.", e);
    }
  }

  // Initially fetch todos and stats on page load
  $effect(() => {
    fetchTodos();
    fetchStats();
  });
</script>

<main class="app">
  <header class="app-header">
    <h1>TODO App</h1>
    
    {#if stats}
      <div class="stats">
        <div class="stat-card">
          <div class="stat-value">{stats.total}</div>
          <div class="stat-label">Total</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">{stats.active}</div>
          <div class="stat-label">Active</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">{stats.completed}</div>
          <div class="stat-label">Completed</div>
        </div>
        <div class="stat-card warning">
          <div class="stat-value">{stats.overdue}</div>
          <div class="stat-label">Overdue</div>
        </div>
      </div>
    {/if}
  </header>

  <div class="controls">
    <div class="search-box">
      <input 
        type="text" 
        placeholder="🔍 Search todos..." 
        bind:value={searchQuery}
        class="search-input"
      />
    </div>

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

    <div class="sort-controls">
      <label for="sort">Sort by:</label>
      <select id="sort" bind:value={sortBy}>
        <option value="created">Created Date</option>
        <option value="dueDate">Due Date</option>
        <option value="priority">Priority</option>
        <option value="title">Title</option>
      </select>
    </div>
  </div>

  <div class="todo-list">
    {#each filteredAndSortedTodos() as todo (todo.id)}
      <Todo 
        {...todo} 
        onToggle={toggleTodo}
        onDelete={deleteTodo}
        onEdit={openEditForm}
      />
    {:else}
      <p class="empty-message">
        {searchQuery ? 'No todos match your search' : 'No todos to display'}
      </p>
    {/each}
  </div>

  <button class="fab" onclick={openAddForm} title="Add new todo">
    +
  </button>

  {#if showAddForm}
    <div class="modal-overlay" onclick={resetForm}></div>
    <div class="modal">
      <div class="modal-header">
        <h2>{editingTodo ? 'Edit Todo' : 'Add New Todo'}</h2>
        <button class="close-btn" onclick={resetForm}>✕</button>
      </div>
      
      <form class="todo-form" onsubmit={handleSubmit}>
        <div class="form-group">
          <label for="title">Title *</label>
          <input 
            id="title"
            placeholder="Enter todo title" 
            bind:value={formTitle}
            required 
          />
        </div>

        <div class="form-group">
          <label for="description">Description *</label>
          <textarea
            id="description"
            placeholder="Enter todo description" 
            bind:value={formDescription}
            required
            rows="4"
          ></textarea>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label for="priority">Priority</label>
            <select id="priority" bind:value={formPriority}>
              <option value="low">Low</option>
              <option value="medium">Medium</option>
              <option value="high">High</option>
            </select>
          </div>

          <div class="form-group">
            <label for="dueDate">Due Date</label>
            <input 
              id="dueDate"
              type="date"
              bind:value={formDueDate}
            />
          </div>
        </div>

        <div class="form-group">
          <label for="tags">Tags (comma-separated)</label>
          <input 
            id="tags"
            placeholder="work, important, urgent" 
            bind:value={formTags}
          />
        </div>

        <div class="form-actions">
          <button type="button" class="cancel-btn" onclick={resetForm}>
            Cancel
          </button>
          <button type="submit" class="submit-btn">
            {editingTodo ? 'Update Todo' : 'Add Todo'}
          </button>
        </div>
      </form>
    </div>
  {/if}
</main>

<style>
  .app {
    color: white;
    background-color: #282c34;
    text-align: center;
    font-size: 24px;
    min-height: 100vh;
    padding: 20px;
    padding-bottom: 100px;
  }

  .app-header {
    font-size: calc(10px + 4vmin);
    margin-top: 30px;
  }

  .stats {
    display: flex;
    justify-content: center;
    gap: 20px;
    margin-top: 30px;
    flex-wrap: wrap;
  }

  .stat-card {
    background-color: #303540;
    padding: 20px 30px;
    border-radius: 12px;
    border: 1px solid #464e61;
    min-width: 120px;
  }

  .stat-card.warning {
    border-color: #e53e3e;
  }

  .stat-value {
    font-size: 36px;
    font-weight: bold;
    margin-bottom: 5px;
  }

  .stat-label {
    font-size: 16px;
    color: #a0aec0;
  }

  .controls {
    margin: 40px auto;
    max-width: 800px;
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .search-box {
    width: 100%;
  }

  .search-input {
    width: 100%;
    font-size: 18px;
    padding: 12px 20px;
  }

  .filters {
    display: flex;
    justify-content: center;
    gap: 15px;
    flex-wrap: wrap;
  }

  .filter-btn {
    background-color: #303540;
    border-color: #464e61;
    font-size: 18px;
    padding: 10px 24px;
  }

  .filter-btn.active {
    background-color: #326ac7;
    border-color: #779ad4;
  }

  .sort-controls {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 10px;
    font-size: 18px;
  }

  select {
    background-color: #303540;
    color: white;
    border: 1px solid #464e61;
    border-radius: 5px;
    padding: 8px 12px;
    font-size: 16px;
    cursor: pointer;
  }

  .todo-list {
    margin: 30px auto;
    max-width: 900px;
  }

  .empty-message {
    color: #888;
    font-size: 20px;
    margin-top: 60px;
  }

  .fab {
    position: fixed;
    bottom: 30px;
    right: 30px;
    width: 70px;
    height: 70px;
    border-radius: 50%;
    background-color: #326ac7;
    color: white;
    font-size: 40px;
    border: none;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
    cursor: pointer;
    transition: all 0.3s ease;
    z-index: 100;
  }

  .fab:hover {
    background-color: #4a83e0;
    transform: scale(1.1);
  }

  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.7);
    z-index: 200;
  }

  .modal {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background-color: #1e2127;
    border: 1px solid #464e61;
    border-radius: 15px;
    padding: 30px;
    max-width: 600px;
    width: 90%;
    max-height: 90vh;
    overflow-y: auto;
    z-index: 300;
    box-shadow: 0 10px 40px rgba(0, 0, 0, 0.5);
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 25px;
  }

  .modal-header h2 {
    margin: 0;
    font-size: 28px;
  }

  .close-btn {
    background-color: transparent;
    border: none;
    color: white;
    font-size: 30px;
    cursor: pointer;
    padding: 0;
    width: 40px;
    height: 40px;
  }

  .close-btn:hover {
    color: #e53e3e;
    background-color: transparent;
  }

  .todo-form {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .form-group label {
    font-size: 16px;
    font-weight: bold;
    color: #a0aec0;
  }

  .form-group input,
  .form-group textarea,
  .form-group select {
    width: 100%;
    font-size: 18px;
    padding: 12px;
    background-color: #303540;
    border: 1px solid #464e61;
    border-radius: 8px;
    color: white;
  }

  .form-group textarea {
    resize: vertical;
    font-family: inherit;
  }

  .form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 20px;
  }

  .form-actions {
    display: flex;
    gap: 15px;
    justify-content: flex-end;
    margin-top: 10px;
  }

  .cancel-btn {
    background-color: #4a5568;
    border-color: #718096;
  }

  .cancel-btn:hover {
    background-color: #5a6678;
  }

  .submit-btn {
    background-color: #48bb78;
    border-color: #68d391;
  }

  .submit-btn:hover {
    background-color: #38a169;
  }

  @media (max-width: 768px) {
    .stats {
      gap: 10px;
    }

    .stat-card {
      min-width: 90px;
      padding: 15px 20px;
    }

    .stat-value {
      font-size: 28px;
    }

    .form-row {
      grid-template-columns: 1fr;
    }

    .todo-list {
      margin: 30px 20px;
    }
  }
</style>