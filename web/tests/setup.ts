import './env';
import { mock, afterEach } from 'bun:test';
import React from 'react';

// Safely require cleanup at top level after DOM setup to prevent hoisting and beforeAll errors
const { cleanup } = require('@testing-library/react');

afterEach(() => {
  cleanup();
  if (globalThis.document && globalThis.document.body) {
    globalThis.document.body.innerHTML = '';
  }
});

// Mock server-only env file
mock.module('@shared/config/env', () => ({
  env: {
    NODE_ENV: 'test',
    LOG_LEVEL: 'error',
    BACKEND_BASE_URL: 'http://localhost:8080',
    BACKEND_TIMEOUT_MS: 4000,
    AUTH_SECRET: 'test-secret',
    NEXTAUTH_URL: 'http://localhost:3001',
  },
}));

// Mock logger to avoid test stdout pollution
mock.module('@shared/lib/logger', () => ({
  logger: {
    debug: () => {},
    info: () => {},
    warn: () => {},
    error: () => {},
    flush: () => {},
  },
  addContext: () => {},
  getCorrelationId: () => 'test-correlation-id',
  runWithContext: (id: string, fn: () => any) => fn(),
}));

// Mock next/headers
mock.module('next/headers', () => ({
  headers: async () => new Map(),
  cookies: async () => ({
    getAll: () => [],
    get: () => undefined,
  }),
}));

// Mock next/image
mock.module('next/image', () => ({
  __esModule: true,
  default: (props: any) => {
    const { src, alt, fill, ...rest } = props;
    return React.createElement('img', { src, alt, ...rest });
  },
}));

// Mock next/link
mock.module('next/link', () => ({
  __esModule: true,
  default: (props: any) => {
    const { href, children, ...rest } = props;
    return React.createElement('a', { href, ...rest }, children);
  },
}));
