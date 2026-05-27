import 'server-only';
import { z } from 'zod';

export const envSchema = z.object({
  AUTH_SECRET: z.string(),
  NEXTAUTH_URL: z.string().url().default('http://localhost:3001'),
  IGDB_CLIENT_ID: z.string().optional(),
  IGDB_CLIENT_SECRET: z.string().optional(),
  BACKEND_BASE_URL: z.string().url(),
  BACKEND_TIMEOUT_MS: z.coerce.number().int().positive().default(4000),
  LOG_LEVEL: z.enum(['debug', 'info', 'warn', 'error']).default('info'),
  NODE_ENV: z.enum(['development', 'production', 'test']).default('development'),
  NEXT_PUBLIC_OTEL_ENDPOINT: z.string().url().default('http://localhost:4318/v1/logs'),
});

const _env = envSchema.safeParse(process.env);

if (!_env.success) {
  console.error('❌ Configuração inválida de variáveis de ambiente:', z.treeifyError(_env.error));
  throw new Error('Variáveis de ambiente inválidas');
}

export const env = _env.data;
