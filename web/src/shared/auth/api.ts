import type { UserSession } from '@shared/backend/types';

type AuthUserResponse = {
  id: string;
  steamId: string;
  name: string;
  image: string | null;
  email: string | null;
};

async function parseAuthResponse(response: Response): Promise<UserSession | null> {
  if (response.status === 401) {
    return null;
  }

  if (!response.ok) {
    throw new Error(`Auth request failed: ${response.status}`);
  }

  const user = (await response.json()) as AuthUserResponse;
  return {
    id: user.id,
    steamId: user.steamId,
    name: user.name,
    image: user.image,
  };
}

export async function fetchAuthSession(backendBaseUrl: string): Promise<UserSession | null> {
  const response = await fetch(new URL('/v1/auth/me', backendBaseUrl).toString(), {
    credentials: 'include',
    cache: 'no-store',
  });

  return parseAuthResponse(response);
}

export async function clearAuthSession(backendBaseUrl: string): Promise<void> {
  await fetch(new URL('/v1/auth/logout', backendBaseUrl).toString(), {
    credentials: 'include',
    cache: 'no-store',
    redirect: 'manual',
  });
}
