import type { Game } from '@entities/game';

export type UserSession = {
  id: string;
  steamId: string;
  name: string;
  image: string | null;
};

export type BackendProvider = {
  getGames(userId: string): Promise<Game[]>;
  getGameById(id: string, userId: string): Promise<Game | null>;
  getSteamGames(steamId: string): Promise<Game[]>;
  getGameStoreDetails(appId: string): Promise<Partial<Game> | null>;
  getMe(): Promise<UserSession | null>;
};
