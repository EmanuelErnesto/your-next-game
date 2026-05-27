'use server';

import type { Game } from '@entities/game';
import { env } from '@shared/config/env';
import { cookies } from 'next/headers';

interface VerdictResponse {
  gameId: string;
  justification: string;
}

export async function getVerdictAction(
  timeAvailable: number,
  mood: string,
): Promise<VerdictResponse | null> {
  const cookieStore = await cookies();
  const cookieHeader = cookieStore
    .getAll()
    .map((c) => `${c.name}=${c.value}`)
    .join('; ');

  try {
    const url = new URL('/v1/recommendations', env.BACKEND_BASE_URL).toString();
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'content-type': 'application/json',
        ...(cookieHeader ? { cookie: cookieHeader } : {}),
      },
      body: JSON.stringify({ timeAvailable, mood }),
      cache: 'no-store',
    });

    if (!response.ok) {
      console.error(`Recommendation error: ${response.status}`);
      return null;
    }

    return await response.json();
  } catch (error) {
    console.error('Failed to get recommendation', error);
    return null;
  }
}
