import 'server-only';

import type { Game } from '@entities/game';
import { httpBackendProvider } from './providers/http';
import type { UserSession } from './types';

export async function getGames(userId: string): Promise<Game[]> {
  return httpBackendProvider.getGames(userId);
}

export async function getGameById(id: string, userId: string): Promise<Game | null> {
  return httpBackendProvider.getGameById(id, userId);
}

export async function getSteamGames(steamId: string): Promise<Game[]> {
  return httpBackendProvider.getSteamGames(steamId);
}

export async function getGameStoreDetails(appId: string): Promise<Partial<Game> | null> {
  return httpBackendProvider.getGameStoreDetails(appId);
}

export async function getMe(): Promise<UserSession | null> {
  return httpBackendProvider.getMe();
}
