import * as jspb from 'google-protobuf'

import * as google_protobuf_duration_pb from 'google-protobuf/google/protobuf/duration_pb'; // proto import: "google/protobuf/duration.proto"


export class CreateTaskRequest extends jspb.Message {
  getName(): string;
  setName(value: string): CreateTaskRequest;

  getType(): TaskType;
  setType(value: TaskType): CreateTaskRequest;

  getRetries(): number;
  setRetries(value: number): CreateTaskRequest;

  getRetrydelay(): google_protobuf_duration_pb.Duration | undefined;
  setRetrydelay(value?: google_protobuf_duration_pb.Duration): CreateTaskRequest;
  hasRetrydelay(): boolean;
  clearRetrydelay(): CreateTaskRequest;

  getCondition(): string;
  setCondition(value: string): CreateTaskRequest;
  hasCondition(): boolean;
  clearCondition(): CreateTaskRequest;

  getLogpayload(): LogPayload | undefined;
  setLogpayload(value?: LogPayload): CreateTaskRequest;
  hasLogpayload(): boolean;
  clearLogpayload(): CreateTaskRequest;

  getHttppayload(): HTTPPayload | undefined;
  setHttppayload(value?: HTTPPayload): CreateTaskRequest;
  hasHttppayload(): boolean;
  clearHttppayload(): CreateTaskRequest;

  getNextList(): Array<CreateTaskRequest>;
  setNextList(value: Array<CreateTaskRequest>): CreateTaskRequest;
  clearNextList(): CreateTaskRequest;
  addNext(value?: CreateTaskRequest, index?: number): CreateTaskRequest;

  getPayloadCase(): CreateTaskRequest.PayloadCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateTaskRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateTaskRequest): CreateTaskRequest.AsObject;
  static serializeBinaryToWriter(message: CreateTaskRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateTaskRequest;
  static deserializeBinaryFromReader(message: CreateTaskRequest, reader: jspb.BinaryReader): CreateTaskRequest;
}

export namespace CreateTaskRequest {
  export type AsObject = {
    name: string;
    type: TaskType;
    retries: number;
    retrydelay?: google_protobuf_duration_pb.Duration.AsObject;
    condition?: string;
    logpayload?: LogPayload.AsObject;
    httppayload?: HTTPPayload.AsObject;
    nextList: Array<CreateTaskRequest.AsObject>;
  };

  export enum PayloadCase {
    PAYLOAD_NOT_SET = 0,
    LOGPAYLOAD = 6,
    HTTPPAYLOAD = 7,
  }

  export enum ConditionCase {
    _CONDITION_NOT_SET = 0,
    CONDITION = 5,
  }
}

export class Task extends jspb.Message {
  getId(): string;
  setId(value: string): Task;
  hasId(): boolean;
  clearId(): Task;

  getName(): string;
  setName(value: string): Task;

  getType(): TaskType;
  setType(value: TaskType): Task;

  getStatus(): TaskStatus;
  setStatus(value: TaskStatus): Task;

  getRetries(): number;
  setRetries(value: number): Task;

  getRetrydelay(): google_protobuf_duration_pb.Duration | undefined;
  setRetrydelay(value?: google_protobuf_duration_pb.Duration): Task;
  hasRetrydelay(): boolean;
  clearRetrydelay(): Task;

  getCondition(): string;
  setCondition(value: string): Task;
  hasCondition(): boolean;
  clearCondition(): Task;

  getLogpayload(): LogPayload | undefined;
  setLogpayload(value?: LogPayload): Task;
  hasLogpayload(): boolean;
  clearLogpayload(): Task;

  getHttppayload(): HTTPPayload | undefined;
  setHttppayload(value?: HTTPPayload): Task;
  hasHttppayload(): boolean;
  clearHttppayload(): Task;

  getNextList(): Array<Task>;
  setNextList(value: Array<Task>): Task;
  clearNextList(): Task;
  addNext(value?: Task, index?: number): Task;

  getPayloadCase(): Task.PayloadCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Task.AsObject;
  static toObject(includeInstance: boolean, msg: Task): Task.AsObject;
  static serializeBinaryToWriter(message: Task, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Task;
  static deserializeBinaryFromReader(message: Task, reader: jspb.BinaryReader): Task;
}

export namespace Task {
  export type AsObject = {
    id?: string;
    name: string;
    type: TaskType;
    status: TaskStatus;
    retries: number;
    retrydelay?: google_protobuf_duration_pb.Duration.AsObject;
    condition?: string;
    logpayload?: LogPayload.AsObject;
    httppayload?: HTTPPayload.AsObject;
    nextList: Array<Task.AsObject>;
  };

  export enum PayloadCase {
    PAYLOAD_NOT_SET = 0,
    LOGPAYLOAD = 8,
    HTTPPAYLOAD = 9,
  }

  export enum IdCase {
    _ID_NOT_SET = 0,
    ID = 1,
  }

  export enum ConditionCase {
    _CONDITION_NOT_SET = 0,
    CONDITION = 7,
  }
}

export class LogPayload extends jspb.Message {
  getMessage(): string;
  setMessage(value: string): LogPayload;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LogPayload.AsObject;
  static toObject(includeInstance: boolean, msg: LogPayload): LogPayload.AsObject;
  static serializeBinaryToWriter(message: LogPayload, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LogPayload;
  static deserializeBinaryFromReader(message: LogPayload, reader: jspb.BinaryReader): LogPayload;
}

export namespace LogPayload {
  export type AsObject = {
    message: string;
  };
}

export class HTTPPayload extends jspb.Message {
  getUrl(): string;
  setUrl(value: string): HTTPPayload;

  getMethod(): string;
  setMethod(value: string): HTTPPayload;

  getBody(): Uint8Array | string;
  getBody_asU8(): Uint8Array;
  getBody_asB64(): string;
  setBody(value: Uint8Array | string): HTTPPayload;

  getHeadersMap(): jspb.Map<string, string>;
  clearHeadersMap(): HTTPPayload;

  getQueryparamsMap(): jspb.Map<string, string>;
  clearQueryparamsMap(): HTTPPayload;

  getTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setTimeout(value?: google_protobuf_duration_pb.Duration): HTTPPayload;
  hasTimeout(): boolean;
  clearTimeout(): HTTPPayload;

  getAuth(): HTTPAuth | undefined;
  setAuth(value?: HTTPAuth): HTTPPayload;
  hasAuth(): boolean;
  clearAuth(): HTTPPayload;

  getFollowredirects(): boolean;
  setFollowredirects(value: boolean): HTTPPayload;

  getVerifyssl(): boolean;
  setVerifyssl(value: boolean): HTTPPayload;

  getExpectedstatuscode(): number;
  setExpectedstatuscode(value: number): HTTPPayload;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HTTPPayload.AsObject;
  static toObject(includeInstance: boolean, msg: HTTPPayload): HTTPPayload.AsObject;
  static serializeBinaryToWriter(message: HTTPPayload, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HTTPPayload;
  static deserializeBinaryFromReader(message: HTTPPayload, reader: jspb.BinaryReader): HTTPPayload;
}

export namespace HTTPPayload {
  export type AsObject = {
    url: string;
    method: string;
    body: Uint8Array | string;
    headersMap: Array<[string, string]>;
    queryparamsMap: Array<[string, string]>;
    timeout?: google_protobuf_duration_pb.Duration.AsObject;
    auth?: HTTPAuth.AsObject;
    followredirects: boolean;
    verifyssl: boolean;
    expectedstatuscode: number;
  };
}

export class HTTPAuth extends jspb.Message {
  getBasic(): HTTPBasicAuth | undefined;
  setBasic(value?: HTTPBasicAuth): HTTPAuth;
  hasBasic(): boolean;
  clearBasic(): HTTPAuth;

  getBearer(): HTTPBearerAuth | undefined;
  setBearer(value?: HTTPBearerAuth): HTTPAuth;
  hasBearer(): boolean;
  clearBearer(): HTTPAuth;

  getApikey(): HTTPApiKeyAuth | undefined;
  setApikey(value?: HTTPApiKeyAuth): HTTPAuth;
  hasApikey(): boolean;
  clearApikey(): HTTPAuth;

  getAuthTypeCase(): HTTPAuth.AuthTypeCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HTTPAuth.AsObject;
  static toObject(includeInstance: boolean, msg: HTTPAuth): HTTPAuth.AsObject;
  static serializeBinaryToWriter(message: HTTPAuth, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HTTPAuth;
  static deserializeBinaryFromReader(message: HTTPAuth, reader: jspb.BinaryReader): HTTPAuth;
}

export namespace HTTPAuth {
  export type AsObject = {
    basic?: HTTPBasicAuth.AsObject;
    bearer?: HTTPBearerAuth.AsObject;
    apikey?: HTTPApiKeyAuth.AsObject;
  };

  export enum AuthTypeCase {
    AUTH_TYPE_NOT_SET = 0,
    BASIC = 1,
    BEARER = 2,
    APIKEY = 3,
  }
}

export class HTTPBasicAuth extends jspb.Message {
  getUsername(): string;
  setUsername(value: string): HTTPBasicAuth;

  getPassword(): string;
  setPassword(value: string): HTTPBasicAuth;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HTTPBasicAuth.AsObject;
  static toObject(includeInstance: boolean, msg: HTTPBasicAuth): HTTPBasicAuth.AsObject;
  static serializeBinaryToWriter(message: HTTPBasicAuth, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HTTPBasicAuth;
  static deserializeBinaryFromReader(message: HTTPBasicAuth, reader: jspb.BinaryReader): HTTPBasicAuth;
}

export namespace HTTPBasicAuth {
  export type AsObject = {
    username: string;
    password: string;
  };
}

export class HTTPBearerAuth extends jspb.Message {
  getToken(): string;
  setToken(value: string): HTTPBearerAuth;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HTTPBearerAuth.AsObject;
  static toObject(includeInstance: boolean, msg: HTTPBearerAuth): HTTPBearerAuth.AsObject;
  static serializeBinaryToWriter(message: HTTPBearerAuth, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HTTPBearerAuth;
  static deserializeBinaryFromReader(message: HTTPBearerAuth, reader: jspb.BinaryReader): HTTPBearerAuth;
}

export namespace HTTPBearerAuth {
  export type AsObject = {
    token: string;
  };
}

export class HTTPApiKeyAuth extends jspb.Message {
  getKey(): string;
  setKey(value: string): HTTPApiKeyAuth;

  getValue(): string;
  setValue(value: string): HTTPApiKeyAuth;

  getLocation(): HTTPApiKeyLocation;
  setLocation(value: HTTPApiKeyLocation): HTTPApiKeyAuth;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HTTPApiKeyAuth.AsObject;
  static toObject(includeInstance: boolean, msg: HTTPApiKeyAuth): HTTPApiKeyAuth.AsObject;
  static serializeBinaryToWriter(message: HTTPApiKeyAuth, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HTTPApiKeyAuth;
  static deserializeBinaryFromReader(message: HTTPApiKeyAuth, reader: jspb.BinaryReader): HTTPApiKeyAuth;
}

export namespace HTTPApiKeyAuth {
  export type AsObject = {
    key: string;
    value: string;
    location: HTTPApiKeyLocation;
  };
}

export enum TaskType {
  TASK_TYPE_UNSPECIFIED = 0,
  TASK_TYPE_LOG = 1,
  TASK_TYPE_HTTP = 2,
}
export enum TaskStatus {
  TASK_STATUS_UNSPECIFIED = 0,
  TASK_STATUS_PENDING = 1,
  TASK_STATUS_RUNNING = 2,
  TASK_STATUS_COMPLETED = 3,
  TASK_STATUS_FAILED = 4,
}
export enum HTTPApiKeyLocation {
  HTTP_API_KEY_LOCATION_UNSPECIFIED = 0,
  HTTP_API_KEY_LOCATION_HEADER = 1,
  HTTP_API_KEY_LOCATION_QUERY = 2,
}
