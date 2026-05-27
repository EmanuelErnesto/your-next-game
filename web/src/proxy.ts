import type { NextRequest } from 'next/server';
import { NextResponse } from 'next/server';

// No Next.js Middleware/Proxy, não podemos importar 'server-only' modules.
// Faremos um log estruturado simples direto no console para o proxy
function logProxyRequest(correlationId: string, request: NextRequest, durationMs: number) {
  console.log(JSON.stringify({
    level: 'info',
    message: 'Canonical Request Log',
    correlationId,
    request: {
      method: request.method,
      path: request.nextUrl.pathname,
      userAgent: request.headers.get('user-agent') || undefined,
      ip: request.headers.get('x-forwarded-for')?.split(',')[0] || undefined,
    },
    durationMs,
  }));
}

export function proxy(request: NextRequest) {
  const correlationId = request.headers.get('x-correlation-id') || `req-${crypto.randomUUID().slice(0, 8)}`;

  const requestHeaders = new Headers(request.headers);
  requestHeaders.set('x-correlation-id', correlationId);

  // Armazenamos metadados adicionais úteis nos headers de invocation para getRequestContext do logger
  requestHeaders.set('x-invoke-path', request.nextUrl.pathname);
  requestHeaders.set('x-invoke-method', request.method);

  const startTime = Date.now();

  const response = NextResponse.next({
    request: {
      headers: requestHeaders,
    },
  });

  response.headers.set('x-correlation-id', correlationId);

  // Log Canonical Request Line para a borda (Edge / Proxy)
  logProxyRequest(correlationId, request, Date.now() - startTime);

  return response;
}

export const config = {
  matcher: ['/((?!_next/static|_next/image|favicon.ico).*)'],
};
