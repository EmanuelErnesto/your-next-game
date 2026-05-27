import type { Game } from '@entities/game';
import { GameCard } from './GameCard';

interface GameGridProps {
  games: Game[];
}

export function GameGrid({ games }: GameGridProps) {
  return (
    <div className="p-8 relative pb-32">
      <div className="mb-6">
        <h2 className="text-2xl font-bold flex items-center gap-2">
          Sua Biblioteca{' '}
          <span className="text-gray-500 text-lg font-normal">({games.length} Jogos)</span>
        </h2>
      </div>

      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-6">
        {games.map((game) => (
          <GameCard key={game.id} game={game} />
        ))}
      </div>
    </div>
  );
}
