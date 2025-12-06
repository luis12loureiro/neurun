import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { WorkflowServiceClient } from '../proto-gen/WorkflowServiceClientPb';
import * as workflow_pb from '../proto-gen/workflow_pb';
import * as task_pb from '../proto-gen/task_pb';

export interface TaskResult {
  workflowId: string;
  taskId: string;
  taskResult: string;
  workflowStatus: string;
  totalTasks: number;
  executedTasks: number;
}

@Injectable({
  providedIn: 'root'
})
export class WorkflowService {
  private client: WorkflowServiceClient;
  private readonly serviceUrl = 'http://localhost:50051'; // gRPC-Web proxy URL

  constructor() {
    this.client = new WorkflowServiceClient(this.serviceUrl, null, null);
  }

  /**
   * Execute a workflow and stream results
   */
  executeWorkflow(workflowId: string): Observable<TaskResult> {
    return new Observable<TaskResult>(observer => {
      const request = new workflow_pb.ExecuteWorkflowRequest();
      request.setId(workflowId);

      const stream = this.client.executeWorkflow(request, {});

      stream.on('data', (response: workflow_pb.ExecuteWorkflowResponse) => {
        const result: TaskResult = {
          workflowId: response.getWorkflowid(),
          taskId: response.getTaskid(),
          taskResult: response.getTaskresult(),
          workflowStatus: this.getWorkflowStatusString(response.getWorkflowstatus()),
          totalTasks: response.getTotaltasks(),
          executedTasks: response.getExecutedtasks()
        };
        observer.next(result);
      });

      stream.on('status', (status: any) => {
        if (status.code !== 0) {
          console.error('Stream error:', status);
          observer.error(new Error(status.details));
        }
      });

      stream.on('end', () => {
        observer.complete();
      });

      // Cleanup function
      return () => {
        stream.cancel();
      };
    });
  }

  /**
   * Get workflow details
   */
  getWorkflow(workflowId: string): Promise<workflow_pb.WorkflowResponse> {
    const request = new workflow_pb.GetWorkflowRequest();
    request.setId(workflowId);
    
    return new Promise((resolve, reject) => {
      this.client.getWorkflow(request, null, (err: any, response: workflow_pb.WorkflowResponse) => {
        if (err) {
          reject(err);
        } else {
          resolve(response);
        }
      });
    });
  }

  /**
   * Create a new workflow
   */
  createWorkflow(
    name: string,
    description: string,
    tasks: task_pb.CreateTaskRequest[]
  ): Promise<workflow_pb.WorkflowResponse> {
    const request = new workflow_pb.CreateWorkflowRequest();
    request.setName(name);
    request.setDescription(description);
    request.setTasksList(tasks);
    
    return new Promise((resolve, reject) => {
      this.client.createWorkflow(request, null, (err: any, response: workflow_pb.WorkflowResponse) => {
        if (err) {
          reject(err);
        } else {
          resolve(response);
        }
      });
    });
  }

  /**
   * Convert WorkflowStatus enum to string
   */
  private getWorkflowStatusString(status: workflow_pb.WorkflowStatus): string {
    switch (status) {
      case workflow_pb.WorkflowStatus.WORKFLOW_STATUS_IDLE:
        return 'idle';
      case workflow_pb.WorkflowStatus.WORKFLOW_STATUS_RUNNING:
        return 'running';
      case workflow_pb.WorkflowStatus.WORKFLOW_STATUS_COMPLETED:
        return 'completed';
      case workflow_pb.WorkflowStatus.WORKFLOW_STATUS_FAILED:
        return 'failed';
      default:
        return 'unknown';
    }
  }
}
