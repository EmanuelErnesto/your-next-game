import { expect, test, describe } from 'bun:test';
import { render, screen, fireEvent } from '@testing-library/react';
import { GameCard } from './GameCard';
import type { Game } from '@entities/game';
import React from 'react';

import { makeGame } from '../../../tests/factories';

const mockGame = makeGame({
  id: 'steam-12345',
  title: 'Portal 2',
  coverUrl: 'https://steamcdn-a.akamaihd.net/steam/apps/620/library_600x900_2x.jpg',
  status: 'Jogado',
  steamAppId: '620',
});

describe('GameCard Component', () => {
  test('renders game card with title, status and correct link', () => {
    render(<GameCard game={mockGame} />);
    
    expect(screen.getByText('Portal 2')).toBeTruthy();
    expect(screen.getByText('Jogado')).toBeTruthy();

    const link = screen.getByRole('link') as HTMLAnchorElement;
    expect(link.getAttribute('href')).toBe('/games/details/steam-12345');
  });

  test('renders extremely long titles correctly', () => {
    const longTitleGame = {
      ...mockGame,
      title: 'A Very Long Game Title That Should Test The Truncation And Styling Of The Title Text inside the card',
    };
    render(<GameCard game={longTitleGame} />);
    expect(screen.getByText(longTitleGame.title)).toBeTruthy();
  });

  test('applies correct status color classes based on played status', () => {
    const { rerender } = render(<GameCard game={mockGame} />);
    let statusLabel = screen.getByText('Jogado');
    expect(statusLabel.className).toContain('text-cyan-400');

    const unplayedGame: Game = {
      ...mockGame,
      status: 'Não Jogado',
    };
    rerender(<GameCard game={unplayedGame} />);
    statusLabel = screen.getByText('Não Jogado');
    expect(statusLabel.className).toContain('text-gray-400');
  });

  test('cycles through fallback image sources on image load error', () => {
    render(<GameCard game={mockGame} />);
    const img = screen.getByRole('img') as HTMLImageElement;
    
    // First source should be the original coverUrl
    expect(img.src).toBe(mockGame.coverUrl);

    // Trigger first error (should fall back to Steam header image)
    fireEvent.error(img);
    expect(img.src).toBe('https://cdn.akamai.steamstatic.com/steam/apps/12345/header.jpg');

    // Trigger second error (should fall back to Steam capsule image)
    fireEvent.error(img);
    expect(img.src).toBe('https://cdn.akamai.steamstatic.com/steam/apps/12345/capsule_616x353.jpg');

    // Trigger third error (should fall back to a placeholder image)
    fireEvent.error(img);
    expect(img.src).toContain('https://via.placeholder.com/600x900');
    expect(img.src).toContain(encodeURIComponent('Portal 2'));
  });
});
