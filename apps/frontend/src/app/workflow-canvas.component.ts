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
      background: #f5f3f0;
    }

    .editor-header {
      padding: 24px 40px;
      border-bottom: 3px solid #2a2925;
      background: #f5f3f0;
      z-index: 10;

      h2 {
        font-size: 1.75rem;
        font-weight: 700;
        margin: 0 0 8px 0;
        color: #2a2925;
      }

      .subtitle {
        font-size: 0.9rem;
        color: #6b6860;
        margin: 0;
      }
    }

    .canvas-container {
      flex: 1;
      background: #e8e5e0;
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
      stroke: #a8a89a;
      stroke-width: 3;
      fill: none;
    }

    .task-node {
      position: absolute;
      width: 180px;
      min-height: 110px;
      background: #f5f3f0;
      border: 3px solid #2a2925;
      box-shadow: 6px 6px 0 rgba(42, 41, 37, 0.15);
      padding: 16px;
      cursor: grab;
      user-select: none;
      z-index: 2;
      transition: box-shadow 0.2s ease;

      &:active {
        cursor: grabbing;
        box-shadow: 8px 8px 0 rgba(42, 41, 37, 0.25);
      }

      &:hover {
        box-shadow: 8px 8px 0 rgba(42, 41, 37, 0.25);
      }

      &.pending {
        border-color: #a8a89a;
      }

      &.running {
        border-color: #b4a582;
        background: #faf7f2;
      }

      &.completed {
        border-color: #8a9f7f;
        background: #f0f4ed;
      }

      &.failed {
        border-color: #9d7070;
        background: #f4ede8;
      }
    }

    .task-icon {
      font-size: 1.5rem;
      font-weight: 700;
      color: #2a2925;
      width: 40px;
      height: 40px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #d1cbc5;
      border: 2px solid #2a2925;
      margin: 0 auto 12px;
    }

    .task-content {
      text-align: center;

      h3 {
        font-size: 0.9rem;
        font-weight: 600;
        margin: 0 0 4px 0;
        color: #2a2925;
      }

      .task-type {
        font-size: 0.75rem;
        color: #6b6860;
        font-weight: 400;
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

    const newX = event.clientX - this.dragState.offsetX;
    const newY = event.clientY - this.dragState.offsetY;

    this.tasks.update(tasks => 
      tasks.map(task => 
        task.id === this.dragState.taskId 
          ? { ...task, x: newX, y: newY }
          : task
      )
    );
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
