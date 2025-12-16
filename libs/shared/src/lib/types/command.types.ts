/**
 * Command types for device control system
 */

export enum CommandType {
  SCRIPT = 0,
  URL = 1,
}

export interface Command {
  id: number;
  commandType: CommandType;
  data: string;
  timestamp: number;
  triggeredBy?: string;
  backendCommandId?: string;
}

export interface CommandPayload {
  commandType: CommandType;
  data: string;
  backendCommandId?: string;
}

export interface CommandResponse {
  success: boolean;
  message?: string;
  command?: Command;
}

export interface CommandEvent {
  commandId: number;
  timestamp: number;
  commandType: CommandType;
  backendCommandId: string;
  blockNumber?: number;
  transactionHash?: string;
}
