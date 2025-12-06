import { Component, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';

interface Task {
  id: string;
  name: string;
  type: string;
  status: 'pending' | 'running' | 'completed' | 'failed';
  x: number;
  y: number;
  connections: string[];
  isStart?: boolean;
}

interface DragState {
  isDragging: boolean;
  taskId: string | null;
  offsetX: number;
  offsetY: number;
}

@Component({
  selector: 'app-workflow-canvas',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './workflow-canvas.component.html',
  styleUrl: './workflow-canvas.component.scss'
})
export class WorkflowCanvasComponent implements OnInit {
  workflowName = signal('Fan-In Workflow Canvas');
  tasks = signal<Task[]>([
    { id: 'start', name: 'Start', type: 'start', status: 'pending', x: 335, y: 20, connections: ['task-a', 'task-b', 'task-c'], isStart: true },
    { id: 'task-a', name: 'Task A', type: 'log', status: 'pending', x: 120, y: 140, connections: ['task-d'] },
    { id: 'task-b', name: 'Task B', type: 'log', status: 'pending', x: 300, y: 140, connections: ['task-d'] },
    { id: 'task-c', name: 'Task C', type: 'log', status: 'pending', x: 480, y: 140, connections: ['task-d'] },
    { id: 'task-d', name: 'Task D', type: 'log', status: 'running', x: 300, y: 280, connections: ['task-e'] },
    { id: 'task-e', name: 'Task E', type: 'log', status: 'completed', x: 300, y: 420, connections: [] }
  ]);

  private dragState: DragState = {
    isDragging: false,
    taskId: null,
    offsetX: 0,
    offsetY: 0
  };

  constructor() {}

  ngOnInit(): void {}

  getTask(id: string): Task | undefined {
    return this.tasks().find(t => t.id === id);
  }

  onTaskMouseDown(event: MouseEvent, taskId: string): void {
    event.stopPropagation();
    const task = this.getTask(taskId);
    if (!task) return;

    this.dragState = {
      isDragging: true,
      taskId,
      offsetX: event.clientX - task.x,
      offsetY: event.clientY - task.y
    };
  }

  onCanvasMouseDown(event: MouseEvent): void {
    // Prevent default canvas drag
  }

  onCanvasMouseMove(event: MouseEvent): void {
    if (!this.dragState.isDragging || !this.dragState.taskId) return;

    event.preventDefault();
    
    const newX = event.clientX - this.dragState.offsetX;
    const newY = event.clientY - this.dragState.offsetY;

    // Use requestAnimationFrame for smoother updates
    requestAnimationFrame(() => {
      this.tasks.update(tasks => 
        tasks.map(task => 
          task.id === this.dragState.taskId 
            ? { ...task, x: newX, y: newY }
            : task
        )
      );
    });
  }

  onCanvasMouseUp(): void {
    this.dragState = {
      isDragging: false,
      taskId: null,
      offsetX: 0,
      offsetY: 0
    };
  }
}
