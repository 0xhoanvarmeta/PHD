/**
 * API request/response types
 */

export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}

export interface PaginatedResponse<T> extends ApiResponse<T[]> {
  total: number;
  page: number;
  limit: number;
  hasMore: boolean;
}

export interface QueryParams {
  page?: number;
  limit?: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
  search?: string;
}

export interface ClientInfo {
  clientId: string;
  deviceName: string;
  os: string;
  osVersion: string;
  lastSeen: Date;
  status: 'online' | 'offline';
}

export interface CommandExecutionResult {
  commandId: number;
  clientId: string;
  success: boolean;
  executedAt: Date;
  error?: string;
  output?: string;
}
