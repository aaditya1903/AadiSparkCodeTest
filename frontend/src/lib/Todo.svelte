<script lang="ts">
  import type { TodoItem } from "./types";

  const { title, description, completed, id, onToggle, onDelete }: TodoItem & {
    onToggle: (id: number) => void;
    onDelete: (id: number) => void;
  } = $props();
</script>

<div class="todo" class:completed={completed}>
  <div class="todo-details">
    <p class="todo-title" class:completed-text={completed}>{title}</p>
    <p class="todo-description">{description}</p>
  </div>
  <div class="todo-actions">
    <button 
      class="toggle-btn" 
      class:completed-btn={completed}
      onclick={() => onToggle(id)}
    >
      {completed ? '✓' : '○'}
    </button>
    <button 
      class="delete-btn"
      onclick={() => onDelete(id)}
    >
      ✕
    </button>
  </div>
</div>

<style>
  .todo {
    display: flex;
    justify-content: space-between;
    align-items: center;

    padding: 0px 20px;

    background-color: #303540;

    border-style: solid;
    border-color: #464e61;
    border-width: 1px;
    border-radius: 15px;

    margin: 10px 0px;
    transition: all 0.3s ease;
  }

  .todo.completed {
    opacity: 0.7;
    background-color: #2a2e38;
  }

  .todo-details {
    display: flex;
    flex-direction: column;
    align-items: start;

    padding: 15px 0px;
    flex-grow: 1;
  }

  .todo-title {
    font-size: 28px;
    font-weight: bold;

    margin: 0px;
    margin-bottom: 5px;
    transition: all 0.3s ease;
  }

  .todo-title.completed-text {
    text-decoration: line-through;
    color: #888;
  }

  .todo-description {
    display: flex;
    margin: 0px;

    font-size: 20px;
    color: #ddd;

    text-align: left;
  }

  .todo-actions {
    display: flex;
    gap: 10px;
  }

  .toggle-btn {
    background-color: #4a5568;
    min-width: 45px;
    height: 45px;
    font-size: 24px;
    padding: 0;
    border-radius: 50%;
    transition: all 0.2s ease;
  }

  .toggle-btn:hover {
    background-color: #5a6578;
  }

  .toggle-btn.completed-btn {
    background-color: #48bb78;
    border-color: #68d391;
  }

  .toggle-btn.completed-btn:hover {
    background-color: #38a169;
  }

  .delete-btn {
    background-color: #e53e3e;
    border-color: #fc8181;
    min-width: 45px;
    height: 45px;
    font-size: 24px;
    padding: 0;
    border-radius: 50%;
    transition: all 0.2s ease;
  }

  .delete-btn:hover {
    background-color: #c53030;
  }
</style>