export type TodoItem = {
  id: number;
  title: string;
  description: string;
  completed: boolean;
  priority: 'low' | 'medium' | 'high';
  dueDate: string;
  tags: string[];
  createdAt: string;
  updatedAt: string;
};

export type FilterType = 'all' | 'active' | 'completed';

export type SortType = 'created' | 'dueDate' | 'priority' | 'title';

export type Stats = {
  total: number;
  active: number;
  completed: number;
  highPriority: number;
  overdue: number;
  byPriority: {
    high: number;
    medium: number;
    low: number;
  };
};