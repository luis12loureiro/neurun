import * as jspb from 'google-protobuf'

import * as task_pb from './task_pb'; // proto import: "task.proto"


export class CreateWorkflowRequest extends jspb.Message {
  getName(): string;
  setName(value: string): CreateWorkflowRequest;

  getDescription(): string;
  setDescription(value: string): CreateWorkflowRequest;
  hasDescription(): boolean;
  clearDescription(): CreateWorkflowRequest;

  getTasksList(): Array<task_pb.CreateTaskRequest>;
  setTasksList(value: Array<task_pb.CreateTaskRequest>): CreateWorkflowRequest;
  clearTasksList(): CreateWorkflowRequest;
  addTasks(value?: task_pb.CreateTaskRequest, index?: number): task_pb.CreateTaskRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateWorkflowRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateWorkflowRequest): CreateWorkflowRequest.AsObject;
  static serializeBinaryToWriter(message: CreateWorkflowRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateWorkflowRequest;
  static deserializeBinaryFromReader(message: CreateWorkflowRequest, reader: jspb.BinaryReader): CreateWorkflowRequest;
}

export namespace CreateWorkflowRequest {
  export type AsObject = {
    name: string;
    description?: string;
    tasksList: Array<task_pb.CreateTaskRequest.AsObject>;
  };

  export enum DescriptionCase {
    _DESCRIPTION_NOT_SET = 0,
    DESCRIPTION = 2,
  }
}

export class GetWorkflowRequest extends jspb.Message {
  getId(): string;
  setId(value: string): GetWorkflowRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetWorkflowRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetWorkflowRequest): GetWorkflowRequest.AsObject;
  static serializeBinaryToWriter(message: GetWorkflowRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetWorkflowRequest;
  static deserializeBinaryFromReader(message: GetWorkflowRequest, reader: jspb.BinaryReader): GetWorkflowRequest;
}

export namespace GetWorkflowRequest {
  export type AsObject = {
    id: string;
  };
}

export class WorkflowResponse extends jspb.Message {
  getId(): string;
  setId(value: string): WorkflowResponse;

  getName(): string;
  setName(value: string): WorkflowResponse;

  getDescription(): string;
  setDescription(value: string): WorkflowResponse;

  getStatus(): string;
  setStatus(value: string): WorkflowResponse;

  getTasksList(): Array<task_pb.Task>;
  setTasksList(value: Array<task_pb.Task>): WorkflowResponse;
  clearTasksList(): WorkflowResponse;
  addTasks(value?: task_pb.Task, index?: number): task_pb.Task;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WorkflowResponse.AsObject;
  static toObject(includeInstance: boolean, msg: WorkflowResponse): WorkflowResponse.AsObject;
  static serializeBinaryToWriter(message: WorkflowResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WorkflowResponse;
  static deserializeBinaryFromReader(message: WorkflowResponse, reader: jspb.BinaryReader): WorkflowResponse;
}

export namespace WorkflowResponse {
  export type AsObject = {
    id: string;
    name: string;
    description: string;
    status: string;
    tasksList: Array<task_pb.Task.AsObject>;
  };
}

export class ExecuteWorkflowRequest extends jspb.Message {
  getId(): string;
  setId(value: string): ExecuteWorkflowRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExecuteWorkflowRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ExecuteWorkflowRequest): ExecuteWorkflowRequest.AsObject;
  static serializeBinaryToWriter(message: ExecuteWorkflowRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExecuteWorkflowRequest;
  static deserializeBinaryFromReader(message: ExecuteWorkflowRequest, reader: jspb.BinaryReader): ExecuteWorkflowRequest;
}

export namespace ExecuteWorkflowRequest {
  export type AsObject = {
    id: string;
  };
}

export class ExecuteWorkflowResponse extends jspb.Message {
  getWorkflowid(): string;
  setWorkflowid(value: string): ExecuteWorkflowResponse;

  getTaskid(): string;
  setTaskid(value: string): ExecuteWorkflowResponse;

  getTaskresult(): string;
  setTaskresult(value: string): ExecuteWorkflowResponse;

  getWorkflowstatus(): WorkflowStatus;
  setWorkflowstatus(value: WorkflowStatus): ExecuteWorkflowResponse;

  getTotaltasks(): number;
  setTotaltasks(value: number): ExecuteWorkflowResponse;

  getExecutedtasks(): number;
  setExecutedtasks(value: number): ExecuteWorkflowResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExecuteWorkflowResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ExecuteWorkflowResponse): ExecuteWorkflowResponse.AsObject;
  static serializeBinaryToWriter(message: ExecuteWorkflowResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExecuteWorkflowResponse;
  static deserializeBinaryFromReader(message: ExecuteWorkflowResponse, reader: jspb.BinaryReader): ExecuteWorkflowResponse;
}

export namespace ExecuteWorkflowResponse {
  export type AsObject = {
    workflowid: string;
    taskid: string;
    taskresult: string;
    workflowstatus: WorkflowStatus;
    totaltasks: number;
    executedtasks: number;
  };
}

export enum WorkflowStatus {
  WORKFLOW_STATUS_UNSPECIFIED = 0,
  WORKFLOW_STATUS_IDLE = 1,
  WORKFLOW_STATUS_RUNNING = 2,
  WORKFLOW_STATUS_COMPLETED = 3,
  WORKFLOW_STATUS_FAILED = 4,
}
