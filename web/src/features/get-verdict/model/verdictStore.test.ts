import { expect, test, mock, describe, beforeEach } from 'bun:test';
import { useVerdictStore } from './verdictStore';
import type { Game } from '@entities/game';
import { makeGame } from '../../../tests/factories';

// Define a type that matches the imported function's signature
type GetVerdictActionType = (timeAvailable: number, mood: string) => Promise<{ gameId: string; justification: string } | null>;

let mockVerdictAction: GetVerdictActionType = async () => null;

mock.module('./getVerdictAction', () => ({
  getVerdictAction: (timeAvailable: number, mood: string) => mockVerdictAction(timeAvailable, mood),
}));

const mockGames: Game[] = [
  makeGame({ id: 'game-1', title: 'Witcher 3', hoursPlayed: 120, rating: 98 }),
  makeGame({ id: 'game-2', title: 'Celeste', status: 'Não Jogado', hoursPlayed: 20, rating: 95 }),
];

describe('verdictStore', () => {
  beforeEach(() => {
    // Reset state before each test and restore actions
    useVerdictStore.setState({
      step: 'CLOSED',
      timeAvailable: 2,
      mood: null,
      genres: [],
      selectedGame: null,
      justification: null,
      openVibeConfig: () =>
        useVerdictStore.setState({
          step: 'CONFIG_VIBE',
          timeAvailable: 2,
          mood: null,
          genres: [],
          selectedGame: null,
          justification: null,
        }),
      close: () => useVerdictStore.setState({ step: 'CLOSED' }),
      setTimeAvailable: (time) => useVerdictStore.setState({ timeAvailable: time }),
      setMood: (mood) => useVerdictStore.setState({ mood }),
      toggleGenre: (genre) =>
        useVerdictStore.setState((state) => ({
          genres: state.genres.includes(genre)
            ? state.genres.filter((g) => g !== genre)
            : [...state.genres, genre],
        })),
      consultAI: async (games: Game[]) => {
        const { mood, timeAvailable } = useVerdictStore.getState();
        if (!mood) return;

        useVerdictStore.setState({ step: 'LOADING' });

        const { getVerdictAction } = require('./getVerdictAction');
        const verdict = await getVerdictAction(timeAvailable, mood);

        if (verdict) {
          const chosen = games.find((g) => g.id === verdict.gameId) || games[0];
          useVerdictStore.setState({ step: 'RESULT', selectedGame: chosen, justification: verdict.justification });
        } else {
          useVerdictStore.setState({ step: 'CONFIG_VIBE' });
        }
      },
      acceptChallenge: () => useVerdictStore.setState({ step: 'CLOSED' }),
      changeMind: () => useVerdictStore.setState({ step: 'CONFIG_VIBE' }),
    });
  });

  test('initial state is correct', () => {
    const state = useVerdictStore.getState();
    expect(state.step).toBe('CLOSED');
    expect(state.timeAvailable).toBe(2);
    expect(state.mood).toBeNull();
    expect(state.genres).toEqual([]);
    expect(state.selectedGame).toBeNull();
    expect(state.justification).toBeNull();
  });

  test('openVibeConfig resets state and opens config', () => {
    // Set dirty state
    useVerdictStore.setState({
      step: 'RESULT',
      mood: 'Energético',
      genres: ['Action'],
      selectedGame: mockGames[0],
      justification: 'Looks fun',
    });

    useVerdictStore.getState().openVibeConfig();

    const state = useVerdictStore.getState();
    expect(state.step).toBe('CONFIG_VIBE');
    expect(state.mood).toBeNull();
    expect(state.genres).toEqual([]);
    expect(state.selectedGame).toBeNull();
    expect(state.justification).toBeNull();
  });

  test('setTimeAvailable and setMood update state', () => {
    const store = useVerdictStore.getState();
    store.setTimeAvailable(5);
    store.setMood('Relaxante');

    const state = useVerdictStore.getState();
    expect(state.timeAvailable).toBe(5);
    expect(state.mood).toBe('Relaxante');
  });

  test('toggleGenre works correctly', () => {
    const store = useVerdictStore.getState();
    store.toggleGenre('RPG');
    expect(useVerdictStore.getState().genres).toEqual(['RPG']);

    store.toggleGenre('Action');
    expect(useVerdictStore.getState().genres).toEqual(['RPG', 'Action']);

    store.toggleGenre('RPG');
    expect(useVerdictStore.getState().genres).toEqual(['Action']);
  });

  test('consultAI flow with successful response', async () => {
    mockVerdictAction = async () => ({
      gameId: 'game-2',
      justification: 'Short platformer for high energy vibes',
    });

    useVerdictStore.setState({ mood: 'Energético', timeAvailable: 2 });
    
    const promise = useVerdictStore.getState().consultAI(mockGames);
    
    // While loading
    expect(useVerdictStore.getState().step).toBe('LOADING');
    
    await promise;

    const state = useVerdictStore.getState();
    expect(state.step).toBe('RESULT');
    expect(state.selectedGame).toEqual(mockGames[1]);
    expect(state.justification).toBe('Short platformer for high energy vibes');
  });

  test('consultAI flow with failed response', async () => {
    mockVerdictAction = async () => null;

    useVerdictStore.setState({ mood: 'Energético', timeAvailable: 2 });
    
    await useVerdictStore.getState().consultAI(mockGames);

    const state = useVerdictStore.getState();
    expect(state.step).toBe('CONFIG_VIBE');
    expect(state.selectedGame).toBeNull();
  });

  test('acceptChallenge and changeMind transition step', () => {
    const store = useVerdictStore.getState();

    store.acceptChallenge();
    expect(useVerdictStore.getState().step).toBe('CLOSED');

    store.changeMind();
    expect(useVerdictStore.getState().step).toBe('CONFIG_VIBE');
  });
});
