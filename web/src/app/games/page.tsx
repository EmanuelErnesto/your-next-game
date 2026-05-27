import type { Game } from '@entities/game';
import { VerdictFlow } from '@features/get-verdict/ui/VerdictFlow';
import { getGames, getSteamGames } from '@shared/backend';
import { getSession } from '@shared/lib/auth';
import { validateImageUrl } from '@shared/lib/image/validateImage';
import { GameGrid } from '@widgets/GameGrid';
import { redirect } from 'next/navigation';
import { logger } from '@/shared/lib/logger';

export const dynamic = 'force-dynamic';

export default async function DashboardPage() {
  const session = await getSession();

  if (!session) {
    redirect('/');
  }

  const [localGames = [], steamGames = []] = await Promise.all([
    getGames(session.id).catch((e) => {
      logger.error('Erro ao carregar games locais:', e as Error);
      return [];
    }),
    session.steamId
      ? getSteamGames(session.steamId).catch((e) => {
          logger.error('Erro ao carregar games da Steam:', e as Error);
          return [];
        })
      : Promise.resolve([]),
  ]);

  const gameMap = new Map<string, Game>();

  localGames.forEach((g) => {
    const uniqueGame = { ...g, id: `local-${g.id}` };
    gameMap.set(uniqueGame.id, uniqueGame);
  });

  steamGames.forEach((g) => {
    const isDuplicate = localGames.some(
      (lg) => lg.steamAppId === g.steamAppId || lg.title.toLowerCase() === g.title.toLowerCase(),
    );

    if (!isDuplicate) {
      const uniqueGame = { ...g, id: `steam-${g.id}` };
      gameMap.set(uniqueGame.id, uniqueGame);
    }
  });

  const gamesList = Array.from(gameMap.values()).sort((a, b) => a.title.localeCompare(b.title));

  const validatedGames = await Promise.all(
    gamesList.map(async (game) => {
      const isValid = await validateImageUrl(game.coverUrl);
      if (isValid) return game;

      if (game.id.startsWith('steam-')) {
        const appId = game.id.replace('steam-', '');
        const headerUrl = `https://cdn.akamai.steamstatic.com/steam/apps/${appId}/header.jpg`;
        if (await validateImageUrl(headerUrl)) return { ...game, coverUrl: headerUrl };
        return {
          ...game,
          coverUrl: `https://cdn.akamai.steamstatic.com/steam/apps/${appId}/capsule_616x353.jpg`,
        };
      }

      return {
        ...game,
        coverUrl: `https://via.placeholder.com/600x900/1a1625/606060?text=${encodeURIComponent(game.title)}`,
      };
    }),
  );

  logger.info('DashboardPage summary:', {
    local: localGames.length,
    steam: steamGames.length,
    total: validatedGames.length,
  });

  return (
    <div className="flex flex-col gap-8">
      <GameGrid games={validatedGames} />
      <VerdictFlow games={validatedGames} />
    </div>
  );
}
