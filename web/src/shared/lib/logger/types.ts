export type LogLevel = 'debug' | 'info' | 'warn' | 'error' | 'fatal' | 'trace';

export type LogContext = Record<string, unknown>;

export type RequestContext = {
  method: string;
  path: string;
  ip?: string;
  userAgent?: string;
  referer?: string;
};

export type LogEvent = {
  correlationId: string;
  request?: RequestContext;
  context?: LogContext;
  durationMs?: number;
  error?: {
    name: string;
    message: string;
    stack?: string;
  };
};
