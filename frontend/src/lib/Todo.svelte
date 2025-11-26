<script lang="ts">
  import type { TodoItem } from "./types";

  const { 
    title, 
    description, 
    completed, 
    id, 
    priority,
    dueDate,
    tags,
    onToggle, 
    onDelete,
    onEdit 
  }: TodoItem & {
    onToggle: (id: number) => void;
    onDelete: (id: number) => void;
    onEdit: (id: number) => void;
  } = $props();

  // Check if todo is overdue
  const isOverdue = $derived(() => {
    if (!dueDate || completed) return false;
    return new Date(dueDate) < new Date();
  });

  // Format due date for display
  const formattedDueDate = $derived(() => {
    if (!dueDate) return '';
    const date = new Date(dueDate);
    const today = new Date();
    const tomorrow = new Date(today);
    tomorrow.setDate(tomorrow.getDate() + 1);
    
    if (date.toDateString() === today.toDateString()) {
      return 'Due Today';
    } else if (date.toDateString() === tomorrow.toDateString()) {
      return 'Due Tomorrow';
    } else {
      return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
    }
  });

  const priorityColor = $derived(() => {
    switch (priority) {
      case 'high': return '#e53e3e';
      case 'medium': return '#dd6b20';
      case 'low': return '#38a169';
      default: return '#4a5568';
    }
  });
</script>

<div class="todo" class:completed={completed} class:overdue={isOverdue()}>
  <div class="todo-main">
    <div class="todo-details">
      <div class="todo-header">
        <p class="todo-title" class:completed-text={completed}>{title}</p>
        <span class="priority-badge" style="background-color: {priorityColor()}">
          {priority.toUpperCase()}
        </span>
      </div>
      <p class="todo-description">{description}</p>
      
      <div class="todo-meta">
        {#if dueDate}
          <span class="due-date" class:overdue-text={isOverdue()}>
            📅 {formattedDueDate()}
          </span>
        {/if}
        
        {#if tags.length > 0}
          <div class="tags">
            {#each tags as tag}
              <span class="tag">{tag}</span>
            {/each}
          </div>
        {/if}
      </div>
    </div>
    
    <div class="todo-actions">
      <button 
        class="toggle-btn" 
        class:completed-btn={completed}
        onclick={() => onToggle(id)}
        title={completed ? 'Mark as incomplete' : 'Mark as complete'}
      >
        {completed ? '✓' : '○'}
      </button>
      <button 
        class="edit-btn"
        onclick={() => onEdit(id)}
        title="Edit todo"
      >
        ✎
      </button>
      <button 
        class="delete-btn"
        onclick={() => onDelete(id)}
        title="Delete todo"
      >
        ✕
      </button>
    </div>
  </div>
</div>

<style>
  .todo {
    background-color: #303540;
    border: 1px solid #464e61;
    border-radius: 15px;
    margin: 10px 0px;
    transition: all 0.3s ease;
    overflow: hidden;
  }

  .todo.completed {
    opacity: 0.7;
    background-color: #2a2e38;
  }

  .todo.overdue {
    border-left: 4px solid #e53e3e;
  }

  .todo-main {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0px 20px;
  }

  .todo-details {
    display: flex;
    flex-direction: column;
    align-items: start;
    padding: 15px 0px;
    flex-grow: 1;
  }

  .todo-header {
    display: flex;
    align-items: center;
    gap: 12px;
    width: 100%;
    margin-bottom: 8px;
  }

  .todo-title {
    font-size: 28px;
    font-weight: bold;
    margin: 0px;
    transition: all 0.3s ease;
  }

  .todo-title.completed-text {
    text-decoration: line-through;
    color: #888;
  }

  .priority-badge {
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 12px;
    font-weight: bold;
    color: white;
  }

  .todo-description {
    margin: 0px 0px 12px 0px;
    font-size: 20px;
    color: #ddd;
    text-align: left;
  }

  .todo-meta {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    align-items: center;
  }

  .due-date {
    font-size: 16px;
    color: #a0aec0;
    display: flex;
    align-items: center;
    gap: 5px;
  }

  .due-date.overdue-text {
    color: #fc8181;
    font-weight: bold;
  }

  .tags {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
  }

  .tag {
    background-color: #4a5568;
    color: #e2e8f0;
    padding: 4px 10px;
    border-radius: 12px;
    font-size: 14px;
  }

  .todo-actions {
    display: flex;
    gap: 10px;
  }

  .toggle-btn,
  .edit-btn,
  .delete-btn {
    min-width: 45px;
    height: 45px;
    font-size: 24px;
    padding: 0;
    border-radius: 50%;
    transition: all 0.2s ease;
  }

  .toggle-btn {
    background-color: #4a5568;
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

  .edit-btn {
    background-color: #3182ce;
    border-color: #63b3ed;
  }

  .edit-btn:hover {
    background-color: #2c5282;
  }

  .delete-btn {
    background-color: #e53e3e;
    border-color: #fc8181;
  }

  .delete-btn:hover {
    background-color: #c53030;
  }
</style>