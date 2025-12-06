import { Component, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { WorkflowService, TaskResult } from './workflow.service';

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
  workflowName = signal('Fan-In Workflow with LOG tasks');
  workflowId = 'fan-in-test';
  
  tasks = signal<Task[]>([
    { id: 'start', name: 'Start', type: 'start', status: 'pending', x: 335, y: 80, connections: ['task-a', 'task-b', 'task-c'], isStart: true },
    { id: 'task-a', name: 'Task A', type: 'log', status: 'pending', x: 150, y: 200, connections: ['task-d'] },
    { id: 'task-b', name: 'Task B', type: 'log', status: 'pending', x: 300, y: 200, connections: ['task-d'] },
    { id: 'task-c', name: 'Task C', type: 'log', status: 'pending', x: 450, y: 200, connections: ['task-d'] },
    { id: 'task-d', name: 'Task D', type: 'log', status: 'pending', x: 300, y: 340, connections: ['task-e'] },
    { id: 'task-e', name: 'Task E', type: 'log', status: 'pending', x: 300, y: 480, connections: [] }
  ]);

  executionResults = signal<TaskResult[]>([]);
  isExecuting = signal(false);

  private dragState: DragState = {
    isDragging: false,
    taskId: null,
    offsetX: 0,
    offsetY: 0
  };

  constructor(private workflowService: WorkflowService) {}

  ngOnInit(): void {}

  executeWorkflow(): void {
    if (this.isExecuting()) return;

    this.isExecuting.set(true);
    this.executionResults.set([]);

    // Reset all tasks to pending
    this.tasks.update(tasks =>
      tasks.map(task => ({ ...task, status: task.isStart ? task.status : 'pending' as const }))
    );

    this.workflowService.executeWorkflow(this.workflowId).subscribe({
      next: (result: TaskResult) => {
        console.log('Task execution result:', result);
        
        // Add result to the list
        this.executionResults.update(results => [...results, result]);

        // Update task status based on result
        this.updateTaskStatus(result.taskId, 'completed');
      },
      error: (error) => {
        console.error('Workflow execution error:', error);
        this.isExecuting.set(false);
      },
      complete: () => {
        this.isExecuting.set(false);
      }
    });
  }

  private updateTaskStatus(taskId: string, status: 'pending' | 'running' | 'completed' | 'failed'): void {
    this.tasks.update(tasks =>
      tasks.map(task =>
        task.id === taskId ? { ...task, status } : task
      )
    );
  }

  onStartClick(): void {
    if (!this.dragState.isDragging) {
      this.executeWorkflow();
    }
  }

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
