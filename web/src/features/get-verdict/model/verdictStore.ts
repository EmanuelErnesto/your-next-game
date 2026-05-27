import type { Game } from '@entities/game';
import { create } from 'zustand';
import { getVerdictAction } from './getVerdictAction';

type Step = 'CLOSED' | 'CONFIG_VIBE' | 'RESULT' | 'LOADING';
type Mood = 'Energético' | 'Relaxante' | 'Desafiador' | 'Narrativo';

interface VerdictState {
  step: Step;
  timeAvailable: number;
  mood: Mood | null;
  genres: string[];
  selectedGame: Game | null;
  justification: string | null;

  openVibeConfig: () => void;
  close: () => void;
  setTimeAvailable: (val: number) => void;
  setMood: (mood: Mood) => void;
  toggleGenre: (genre: string) => void;
  consultAI: (games: Game[]) => Promise<void>;
  acceptChallenge: () => void;
  changeMind: () => void;
}

export const useVerdictStore = create<VerdictState>((set, get) => ({
  step: 'CLOSED',
  timeAvailable: 2,
  mood: null,
  genres: [],
  selectedGame: null,
  justification: null,

  openVibeConfig: () =>
    set({
      step: 'CONFIG_VIBE',
      timeAvailable: 2,
      mood: null,
      genres: [],
      selectedGame: null,
      justification: null,
    }),
  close: () => set({ step: 'CLOSED' }),

  setTimeAvailable: (time) => set({ timeAvailable: time }),
  setMood: (mood) => set({ mood }),
  toggleGenre: (genre) =>
    set((state) => ({
      genres: state.genres.includes(genre)
        ? state.genres.filter((g) => g !== genre)
        : [...state.genres, genre],
    })),

  consultAI: async (games: Game[]) => {
    const { mood, timeAvailable } = get();
    if (!mood) return;

    set({ step: 'LOADING' });

    const verdict = await getVerdictAction(timeAvailable, mood);

    if (verdict) {
      const chosen = games.find((g) => g.id === verdict.gameId) || games[0];
      set({ step: 'RESULT', selectedGame: chosen, justification: verdict.justification });
    } else {
      console.error('Failed to get recommendation from AI');
      set({ step: 'CONFIG_VIBE' });
    }
  },

  acceptChallenge: () => {
    set({ step: 'CLOSED' });
  },

  changeMind: () => {
    set({ step: 'CONFIG_VIBE' });
  },
}));
