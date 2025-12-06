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
  template: `
    <div class="workflow-editor-container">
      <div class="editor-header">
        <h2>{{ workflowName() }}</h2>
        <p class="subtitle">{{ workflowDescription() }}</p>
      </div>
      <div class="canvas-container"
           (mousedown)="onCanvasMouseDown($event)"
           (mousemove)="onCanvasMouseMove($event)"
           (mouseup)="onCanvasMouseUp()"
           (mouseleave)="onCanvasMouseUp()">
        
        <!-- SVG for connections -->
        <svg class="connections-svg">
          @for (task of tasks(); track task.id) {
            @for (targetId of task.connections; track targetId) {
              <line 
                [attr.x1]="task.x + 90" 
                [attr.y1]="task.y + 110"
                [attr.x2]="getTask(targetId)!.x + 90"
                [attr.y2]="getTask(targetId)!.y + 10"
                class="connection-line" />
            }
          }
        </svg>

        <!-- Task nodes -->
        @for (task of tasks(); track task.id) {
          <div class="task-node"
               [class.pending]="task.status === 'pending'"
               [class.running]="task.status === 'running'"
               [class.completed]="task.status === 'completed'"
               [class.failed]="task.status === 'failed'"
               [style.left.px]="task.x"
               [style.top.px]="task.y"
               (mousedown)="onTaskMouseDown($event, task.id)">
            <div class="task-icon">{{ task.name.charAt(task.name.length - 1) }}</div>
            <div class="task-content">
              <h3>{{ task.name }}</h3>
              <span class="task-type">{{ task.type }}</span>
            </div>
          </div>
        }
      </div>
    </div>
  `,
  styles: [`
    .workflow-editor-container {
      width: 100%;
      height: 100vh;
      display: flex;
      flex-direction: column;
      background: #fafafa;
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
    }

    .editor-header {
      padding: 20px 32px;
      background: #ffffff;
      border-bottom: 1px solid #e5e7eb;
      box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
      z-index: 10;

      h2 {
        font-size: 1.5rem;
        font-weight: 600;
        margin: 0 0 4px 0;
        color: #111827;
        letter-spacing: -0.01em;
      }

      .subtitle {
        font-size: 0.875rem;
        color: #6b7280;
        margin: 0;
        font-weight: 400;
      }
    }

    .canvas-container {
      flex: 1;
      background: linear-gradient(0deg, #f9fafb 0%, #ffffff 100%);
      position: relative;
      overflow: hidden;
      cursor: default;
    }

    .connections-svg {
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      pointer-events: none;
      z-index: 1;
    }

    .connection-line {
      stroke: #9ca3af;
      stroke-width: 2;
      fill: none;
      opacity: 0.6;
    }

    .task-node {
      position: absolute;
      width: 180px;
      background: #ffffff;
      border: 1px solid #e5e7eb;
      border-radius: 8px;
      box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
      padding: 16px;
      cursor: grab;
      user-select: none;
      z-index: 2;
      transition: box-shadow 0.2s ease;
      will-change: transform;

      &:active {
        cursor: grabbing;
        box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
        transform: translateY(-2px);
      }

      &:hover {
        box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
        border-color: #d1d5db;
      }

      &.pending {
        border-left: 3px solid #9ca3af;
      }

      &.running {
        border-left: 3px solid #3b82f6;
        background: #eff6ff;
      }

      &.completed {
        border-left: 3px solid #10b981;
        background: #f0fdf4;
      }

      &.failed {
        border-left: 3px solid #ef4444;
        background: #fef2f2;
      }
    }

    .task-icon {
      font-size: 1.25rem;
      font-weight: 600;
      color: #ffffff;
      width: 36px;
      height: 36px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      border-radius: 6px;
      margin: 0 0 12px 0;
    }

    .task-node.running .task-icon {
      background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
    }

    .task-node.completed .task-icon {
      background: linear-gradient(135deg, #10b981 0%, #059669 100%);
    }

    .task-node.failed .task-icon {
      background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
    }

    .task-content {
      h3 {
        font-size: 0.9rem;
        font-weight: 600;
        margin: 0 0 6px 0;
        color: #111827;
        line-height: 1.3;
      }

      .task-type {
        display: inline-block;
        font-size: 0.75rem;
        color: #6b7280;
        background: #f3f4f6;
        padding: 2px 8px;
        border-radius: 4px;
        font-weight: 500;
      }
    }
  `]
})
export class WorkflowCanvasComponent implements OnInit {
  workflowName = signal('Fan-In Test Workflow');
  workflowDescription = signal('Drag nodes to rearrange • 3 root tasks → 1 shared task');

  tasks = signal<Task[]>([
    { id: 'task-a', name: 'Root Task A', type: 'log', status: 'pending', x: 100, y: 100, connections: ['task-d'] },
    { id: 'task-b', name: 'Root Task B', type: 'log', status: 'pending', x: 400, y: 100, connections: ['task-d'] },
    { id: 'task-c', name: 'Root Task C', type: 'log', status: 'pending', x: 700, y: 100, connections: ['task-d'] },
    { id: 'task-d', name: 'Shared Task D', type: 'log', status: 'running', x: 400, y: 300, connections: ['task-e'] },
    { id: 'task-e', name: 'Task E', type: 'log', status: 'completed', x: 400, y: 500, connections: [] }
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
