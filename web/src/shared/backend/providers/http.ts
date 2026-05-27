import type { Game } from '@entities/game';
import { env } from '@shared/config/env';
import { cookies } from 'next/headers';
import { logger } from '@shared/lib/logger';
import type { BackendProvider, UserSession } from '../types';

type BackendAuthUser = {
  id: string;
  name: string;
  image: string | null;
  email: string | null;
  steamId?: string | null;
};

interface RequestOptions {
  method?: 'GET' | 'POST' | 'PATCH';
  body?: unknown;
}

function buildUrl(path: string) {
  return new URL(path, env.BACKEND_BASE_URL).toString();
}

async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), env.BACKEND_TIMEOUT_MS);

  const cookieStore = await cookies();
  const cookieHeader = cookieStore
    .getAll()
    .map((c) => `${c.name}=${c.value}`)
    .join('; ');

  try {
    const response = await fetch(buildUrl(path), {
      method: options.method ?? 'GET',
      headers: {
        'content-type': 'application/json',
        ...(cookieHeader ? { cookie: cookieHeader } : {}),
      },
      body: options.body ? JSON.stringify(options.body) : undefined,
      signal: controller.signal,
      cache: 'no-store',
    });

    if (!response.ok) {
      if (response.status === 401) {
        throw new Error('Unauthorized');
      }
      throw new Error(`Backend request failed: ${response.status}`);
    }

    return (await response.json()) as T;
  } finally {
    clearTimeout(timeout);
  }
}

export const httpBackendProvider: BackendProvider = {
  async getGames(userId: string) {
    return request<Game[]>(`/v1/users/${encodeURIComponent(userId)}/games`);
  },
  async getGameById(id: string, userId: string) {
    return request<Game | null>(
      `/v1/users/${encodeURIComponent(userId)}/games/${encodeURIComponent(id)}`,
    );
  },
  async getSteamGames(steamId: string) {
    return request<Game[]>(
      `/v1/integrations/steam/users/${encodeURIComponent(steamId)}/games`,
    );
  },
  async getGameStoreDetails(appId: string) {
    return request<Partial<Game> | null>(`/v1/catalog/steam/apps/${encodeURIComponent(appId)}`);
  },
  async getMe() {
    try {
      const user = await request<BackendAuthUser>('/v1/auth/me');
      return {
        id: user.id,
        steamId: user.steamId ?? null,
        name: user.name,
        image: user.image,
      };
    } catch (e) {
      if (e instanceof Error && e.message === 'Unauthorized') {
        return null;
      }
      logger.error('Failed to fetch session from backend', e instanceof Error ? e : new Error(String(e)));
      return null;
    }
  },
};
