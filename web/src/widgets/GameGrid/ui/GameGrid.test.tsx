import { expect, test, describe } from 'bun:test';
import { render, screen } from '@testing-library/react';
import { GameGrid } from './GameGrid';
import type { Game } from '@entities/game';
import React from 'react';

import { makeGame } from '../../../tests/factories';

const mockGames: Game[] = [
  makeGame({ id: 'game-1', title: 'Witcher 3', hoursPlayed: 100, rating: 98 }),
  makeGame({ id: 'game-2', title: 'Celeste', status: 'Não Jogado', hoursPlayed: 20, rating: 95 }),
];

describe('GameGrid Component', () => {
  test('renders grid with games and correct game count label', () => {
    render(<GameGrid games={mockGames} />);
    
    expect(screen.getByText('Sua Biblioteca')).toBeTruthy();
    expect(screen.getByText('(2 Jogos)')).toBeTruthy();

    expect(screen.getByText('Witcher 3')).toBeTruthy();
    expect(screen.getByText('Celeste')).toBeTruthy();
  });

  test('renders correctly with an empty games array', () => {
    render(<GameGrid games={[]} />);
    
    expect(screen.getByText('Sua Biblioteca')).toBeTruthy();
    expect(screen.getByText('(0 Jogos)')).toBeTruthy();
    
    // No cards should be rendered
    expect(screen.queryByRole('img')).toBeNull();
  });
});
