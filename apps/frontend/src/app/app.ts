import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { WorkflowCanvasComponent } from './workflow-canvas/workflow-canvas.component';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, WorkflowCanvasComponent],
  templateUrl: './app.html',
  styleUrl: './app.scss'
})
export class App {}
