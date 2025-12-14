/**
 * Application constants
 */

export const APP_CONFIG = {
  DEFAULT_POLLING_INTERVAL: 5000, // 5 seconds
  MAX_RETRY_ATTEMPTS: 3,
  COMMAND_TIMEOUT: 30000, // 30 seconds
} as const;

export const CONTRACT_EVENTS = {
  COMMAND_TRIGGERED: 'CommandTriggered',
  ADMIN_UPDATED: 'AdminUpdated',
} as const;

export const ERROR_MESSAGES = {
  CONTRACT_NOT_FOUND: 'Smart contract not found',
  INVALID_COMMAND: 'Invalid command',
  UNAUTHORIZED: 'Unauthorized access',
  NETWORK_ERROR: 'Network error occurred',
  EXECUTION_FAILED: 'Command execution failed',
} as const;
