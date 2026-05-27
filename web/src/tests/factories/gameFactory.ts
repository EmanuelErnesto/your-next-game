import type { Game } from '@entities/game';

export function makeGame(overrides: Partial<Game> = {}): Game {
  return {
    id: 'steam-12345',
    title: 'Portal 2',
    coverUrl: 'https://steamcdn-a.akamaihd.net/steam/apps/620/library_600x900_2x.jpg',
    status: 'Jogado',
    genre: 'Puzzle',
    platform: 'PC',
    description: 'Portal puzzle game',
    developer: 'Valve',
    releaseDate: '2011-04-18',
    hoursPlayed: 45,
    rating: 99,
    steamAppId: '620',
    ...overrides,
  };
}
