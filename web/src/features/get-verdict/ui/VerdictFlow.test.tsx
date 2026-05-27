import { expect, test, mock, describe, beforeEach } from 'bun:test';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { VerdictFlow } from './VerdictFlow';
import { useVerdictStore } from '../model/verdictStore';
import type { Game } from '@entities/game';
import React from 'react';

import { makeGame } from '../../../tests/factories';

const mockGames: Game[] = [
  makeGame({ id: 'game-1', title: 'Witcher 3', hoursPlayed: 120, rating: 98 }),
  makeGame({ id: 'game-2', title: 'Celeste', status: 'Não Jogado', hoursPlayed: 20, rating: 95 }),
];

const originalConsultAI = useVerdictStore.getState().consultAI;
const originalAcceptChallenge = useVerdictStore.getState().acceptChallenge;
const originalChangeMind = useVerdictStore.getState().changeMind;

describe('VerdictFlow Component', () => {
  beforeEach(() => {
    // Reset Zustand store and restore original actions
    useVerdictStore.setState({
      step: 'CLOSED',
      timeAvailable: 2,
      mood: null,
      genres: [],
      selectedGame: null,
      justification: null,
      consultAI: originalConsultAI,
      acceptChallenge: originalAcceptChallenge,
      changeMind: originalChangeMind,
    });
  });

  test('renders closed state (only trigger button visible)', () => {
    render(<VerdictFlow games={mockGames} />);
    
    // Obter Veredito button should be visible
    const triggerButton = screen.getByText('Obter Veredito');
    expect(triggerButton).toBeTruthy();

    // Modal content should not be present
    expect(screen.queryByText('Nova recomendação de Jogo')).toBeNull();
  });

  test('renders CONFIG_VIBE state', () => {
    useVerdictStore.setState({ step: 'CONFIG_VIBE' });
    render(<VerdictFlow games={mockGames} />);

    expect(screen.getByText('Nova recomendação de Jogo')).toBeTruthy();
    expect(screen.getByText('Tempo Disponível (horas)')).toBeTruthy();
    
    // Mood options
    expect(screen.getByText('Energético')).toBeTruthy();
    expect(screen.getByText('Relaxante')).toBeTruthy();
    expect(screen.getByText('Desafiador')).toBeTruthy();
    expect(screen.getByText('Narrativo')).toBeTruthy();

    // Button should be disabled since mood is null
    const submitButton = screen.getByText('Gerar Recomendação') as HTMLButtonElement;
    expect(submitButton.disabled).toBe(true);
  });

  test('enables button and triggers consultAI when mood is selected', () => {
    const mockConsultAI = mock(async () => {});
    
    useVerdictStore.setState({
      step: 'CONFIG_VIBE',
      mood: 'Energético',
      consultAI: mockConsultAI,
    });

    render(<VerdictFlow games={mockGames} />);

    const submitButton = screen.getByText('Gerar Recomendação') as HTMLButtonElement;
    expect(submitButton.disabled).toBe(false);

    fireEvent.click(submitButton);
    expect(mockConsultAI).toHaveBeenCalled();
  });

  test('renders LOADING state', () => {
    useVerdictStore.setState({ step: 'LOADING' });
    render(<VerdictFlow games={mockGames} />);

    expect(screen.getByText('A IA está analisando sua biblioteca...')).toBeTruthy();
    expect(screen.getByText('Isso pode levar alguns segundos.')).toBeTruthy();
  });

  test('renders RESULT state with selected game details', () => {
    useVerdictStore.setState({
      step: 'RESULT',
      selectedGame: mockGames[0],
      justification: 'Porque é uma obra de arte.',
    });

    render(<VerdictFlow games={mockGames} />);

    expect(screen.getByText('Seu jogo ideal para hoje é...')).toBeTruthy();
    expect(screen.getByText('Witcher 3')).toBeTruthy();
    expect(screen.getByText('Porque é uma obra de arte.')).toBeTruthy();

    // Challenge actions
    expect(screen.getByText('✓ Aceitar Desafio')).toBeTruthy();
    expect(screen.getByText('✕ Mudar de Ideia')).toBeTruthy();
  });

  test('handles action buttons in RESULT step', () => {
    const mockAcceptChallenge = mock(() => {});
    const mockChangeMind = mock(() => {});

    useVerdictStore.setState({
      step: 'RESULT',
      selectedGame: mockGames[0],
      justification: 'Porque é uma obra de arte.',
      acceptChallenge: mockAcceptChallenge,
      changeMind: mockChangeMind,
    });

    render(<VerdictFlow games={mockGames} />);

    const acceptButton = screen.getByText('✓ Aceitar Desafio');
    const changeButton = screen.getByText('✕ Mudar de Ideia');

    fireEvent.click(acceptButton);
    expect(mockAcceptChallenge).toHaveBeenCalled();

    fireEvent.click(changeButton);
    expect(mockChangeMind).toHaveBeenCalled();
  });
});
