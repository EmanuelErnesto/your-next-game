import type { Game } from '@entities/game';
import { getGameById, getGameStoreDetails } from '@shared/backend';
import { getSession } from '@shared/lib/auth';
import { ArrowLeft, Calendar, Clock, Code, ExternalLink, Gamepad2, Star } from 'lucide-react';
import Image from 'next/image';
import Link from 'next/link';
import { notFound, redirect } from 'next/navigation';

type GameDetailsPageProps = {
  params: Promise<{ gameId: string }>;
};

export default async function GameDetailsPage({ params }: GameDetailsPageProps) {
  const { gameId } = await params;
  const session = await getSession();
  if (!session) {
    redirect('/');
  }

  const localGame = await getGameById(gameId, session.id);

  const appId = localGame?.steamAppId || gameId;
  const storeData = await getGameStoreDetails(appId);

  if (!localGame && !storeData) {
    notFound();
  }

  const game: Game = {
    id: gameId,
    title: storeData?.title || localGame?.title || 'Jogo Desconhecido',
    coverUrl: storeData?.coverUrl || localGame?.coverUrl || '',
    status: localGame?.status || 'Não Jogado',
    genre: storeData?.genre || localGame?.genre || 'Steam',
    platform: localGame?.platform || 'Steam',
    hoursPlayed: localGame?.hoursPlayed || 0,
    description: storeData?.description || localGame?.description || '',
    developer: storeData?.developer || localGame?.developer || '',
    releaseDate: storeData?.releaseDate || localGame?.releaseDate || '',
    rating: localGame?.rating || null,
    steamAppId: appId,
  };

  const isPlayed = game.status === 'Jogado';
  const statusColor = isPlayed
    ? 'text-cyan-400 border-cyan-400 bg-cyan-400/10'
    : 'text-gray-400 border-gray-600 bg-gray-600/10';
  const steamUrl = `https://store.steampowered.com/app/${appId}`;

  return (
    <div className="h-full overflow-y-auto">
      <div className="relative w-full h-72 md:h-96 overflow-hidden">
        <Image
          src={game.coverUrl}
          alt={game.title}
          fill
          className="object-cover object-top blur-sm scale-105"
          priority
        />
        <div className="absolute inset-0 bg-linear-to-t from-background-surface via-background-surface/80 to-transparent"></div>

        <div className="absolute top-6 left-6 z-10">
          <Link
            href="/games"
            className="btn btn-ghost btn-sm gap-2 text-gray-300 hover:text-white backdrop-blur-sm bg-black/20"
          >
            <ArrowLeft size={16} />
            Voltar à Biblioteca
          </Link>
        </div>
      </div>

      <div className="relative -mt-32 z-10 px-8 pb-12 max-w-5xl mx-auto">
        <div className="flex flex-col md:flex-row gap-8">
          <div className="w-48 md:w-56 shrink-0">
            <div className="relative aspect-2/3 rounded-xl overflow-hidden border-2 border-[#2a475e] shadow-[0_0_30px_rgba(102,192,244,0.2)]">
              <Image src={game.coverUrl} alt={game.title} fill className="object-cover" />
            </div>
          </div>

          <div className="flex flex-col gap-4 flex-1 min-w-0 pt-4">
            <div className="flex flex-wrap items-center gap-3">
              <span
                className={`text-xs px-3 py-1 rounded-full border font-semibold uppercase ${statusColor}`}
              >
                {game.status}
              </span>
              <span className="badge badge-outline badge-sm text-primary border-primary/50">
                {game.genre}
              </span>
              <span className="badge badge-outline badge-sm text-gray-400 border-gray-600">
                {game.platform}
              </span>
            </div>

            <h1 className="text-3xl md:text-4xl font-extrabold text-white tracking-tight">
              {game.title}
            </h1>

            {game.developer && (
              <p className="text-gray-400 text-sm">
                por <span className="text-primary font-medium">{game.developer}</span>
              </p>
            )}

            <div className="flex flex-wrap gap-3 mt-2">
              <a
                href={steamUrl}
                target="_blank"
                rel="noopener noreferrer"
                className="btn btn-primary btn-sm gap-2"
              >
                <ExternalLink size={14} />
                Abrir na Steam
              </a>
              <button className="btn btn-outline btn-sm gap-2 border-[#2a475e] text-gray-300 hover:bg-[#2a475e]">
                <Gamepad2 size={14} />
                Marcar como Jogando
              </button>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mt-10">
          <StatCard
            icon={<Clock size={20} className="text-primary" />}
            label="Horas Jogadas"
            value={game.hoursPlayed > 0 ? `${game.hoursPlayed.toFixed(1)}h` : '—'}
          />
          <StatCard
            icon={<Star size={20} className="text-yellow-400" />}
            label="Nota Pessoal"
            value={game.rating ? `${game.rating.toFixed(1)} / 5` : '—'}
          />
          <StatCard
            icon={<Calendar size={20} className="text-primary" />}
            label="Lançamento"
            value={game.releaseDate || '—'}
          />
          <StatCard
            icon={<Code size={20} className="text-primary" />}
            label="Desenvolvedor"
            value={game.developer || '—'}
          />
        </div>

        {game.description && (
          <div className="mt-10 bg-background-base/60 border border-[#2a475e]/50 rounded-2xl p-6">
            <h2 className="text-lg font-bold text-white mb-4">Sobre o Jogo</h2>
            <div
              className="text-gray-400 leading-relaxed text-sm overflow-hidden"
              dangerouslySetInnerHTML={{ __html: game.description }}
            />
          </div>
        )}
      </div>
    </div>
  );
}

function StatCard({ icon, label, value }: { icon: React.ReactNode; label: string; value: string }) {
  return (
    <div className="bg-background-base/60 border border-[#2a475e]/50 rounded-xl p-4 flex flex-col gap-2">
      <div className="flex items-center gap-2 text-gray-400 text-xs font-medium uppercase tracking-wide">
        {icon}
        {label}
      </div>
      <span className="text-white font-bold text-lg truncate">{value}</span>
    </div>
  );
}
