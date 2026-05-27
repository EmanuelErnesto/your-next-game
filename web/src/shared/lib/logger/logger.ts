import { AsyncLocalStorage } from 'node:async_hooks';
import { env } from '@shared/config/env';
import { headers } from 'next/headers';
import pino from 'pino';
import type { LogContext, LogEvent, LogLevel, RequestContext } from './types';

const pinoInstance = pino({
  level: env.LOG_LEVEL,
  transport:
    env.NODE_ENV === 'development'
      ? {
          target: 'pino-pretty',
          options: {
            colorize: true,
          },
        }
      : undefined,
});

type LoggerContext = {
  correlationId: string;
  data: LogContext;
  startTime: number;
};

const storage = new AsyncLocalStorage<LoggerContext>();

export async function getRequestContext(): Promise<RequestContext | undefined> {
  try {
    const headersList = await headers();

    const path = headersList.get('x-invoke-path') || '';
    const method = headersList.get('x-invoke-method') || 'GET';
    const ip =
      headersList.get('x-forwarded-for')?.split(',')[0] ||
      headersList.get('x-real-ip') ||
      undefined;
    const userAgent = headersList.get('user-agent') || undefined;
    const referer = headersList.get('referer') || undefined;

    if (!path && !userAgent) return undefined;

    return {
      method,
      path,
      ip,
      userAgent,
      referer,
    };
  } catch {
    return undefined;
  }
}

export async function getCorrelationId(): Promise<string> {
  const store = storage.getStore();
  if (store) return store.correlationId;

  try {
    const headersList = await headers();
    return headersList.get('x-correlation-id') || `req-${crypto.randomUUID().slice(0, 8)}`;
  } catch {
    return 'system';
  }
}

export function runWithContext<T>(correlationId: string, fn: () => T): T {
  return storage.run({ correlationId, data: {}, startTime: Date.now() }, fn);
}

export function addContext(key: string, value: unknown): void {
  const store = storage.getStore();
  if (store) {
    store.data[key] = value;
  }
}

// Provedor de logs simplificado para OpenTelemetry
async function sendToOTel(level: LogLevel, message: string, payload: LogEvent) {
  // Apenas envia logs de info para o OTel se NODE_ENV for development
  if (level === 'info' && env.NODE_ENV !== 'development') {
    return;
  }

  try {
    await fetch(env.NEXT_PUBLIC_OTEL_ENDPOINT, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        resourceLogs: [{
          resource: {
            attributes: [
              { key: 'service.name', value: { stringValue: 'your-next-game-frontend' } },
              { key: 'deployment.environment', value: { stringValue: env.NODE_ENV } }
            ]
          },
          scopeLogs: [{
            logRecords: [{
              timeUnixNano: String(Date.now() * 1000000),
              severityText: level.toUpperCase(),
              severityNumber: level === 'error' ? 17 : level === 'warn' ? 13 : 9,
              body: { stringValue: message },
              attributes: [
                { key: 'correlation_id', value: { stringValue: payload.correlationId } },
                { key: 'context', value: { stringValue: JSON.stringify(payload.context || {}) } },
                { key: 'request', value: { stringValue: JSON.stringify(payload.request || {}) } },
                { key: 'error', value: { stringValue: JSON.stringify(payload.error || {}) } }
              ]
            }]
          }]
        }]
      })
    });
  } catch {
    // Falha silenciosa no envio do OTel para não quebrar a aplicação
  }
}

async function emit(
  level: LogLevel,
  message: string,
  extra?: LogContext,
  error?: Error,
): Promise<void> {
  const store = storage.getStore();
  const correlationId = store?.correlationId || (await getCorrelationId());
  const request = await getRequestContext();

  const payload: LogEvent = {
    correlationId,
    context: {
      ...store?.data,
      ...extra,
    },
    request,
  };

  if (error) {
    payload.error = {
      name: error.name,
      message: error.message,
      stack: error.stack,
    };
  }

  // Regra do Usuário: Todos os logs de warn/error vão para o OTel. Em produção (ou browser runtime), 
  // evitamos imprimi-los no console do browser (terminal do browser).
  const isBrowser = typeof window !== 'undefined';
  const shouldPrintConsole = !(env.NODE_ENV === 'production' && (level === 'warn' || level === 'error') && isBrowser);

  if (shouldPrintConsole) {
    pinoInstance[level](payload, message);
  }

  // Exportar para o OpenTelemetry Collector
  await sendToOTel(level, message, payload);
}

export const logger = {
  debug: (message: string, extra?: LogContext) => emit('debug', message, extra),
  info: (message: string, extra?: LogContext) => emit('info', message, extra),
  warn: (message: string, extra?: LogContext) => emit('warn', message, extra),
  error: (message: string, error?: Error, extra?: LogContext) =>
    emit('error', message, extra, error),

  flush: async (message = 'Wide Event (Canonical Log Line)') => {
    const store = storage.getStore();
    const duration = store ? Date.now() - store.startTime : 0;

    await emit('info', message, { durationMs: duration });
  },
};
