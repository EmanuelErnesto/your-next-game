'use client';

import type { Game } from '@entities/game';
import Image from 'next/image';
import Link from 'next/link';
import { useState } from 'react';

interface GameCardProps {
  game: Game;
}

export function GameCard({ game }: GameCardProps) {
  const [imgSrc, setImgSrc] = useState(game.coverUrl);
  const isPlayed = game.status === 'Jogado';
  const statusColor = isPlayed ? 'text-cyan-400 border-cyan-400' : 'text-gray-400 border-gray-600';

  const handleError = () => {
    if (imgSrc.includes('library_600x900_2x.jpg')) {
      const appId = game.id.replace('steam-', '');
      setImgSrc(`https://cdn.akamai.steamstatic.com/steam/apps/${appId}/header.jpg`);
    } else if (imgSrc.includes('header.jpg')) {
      const appId = game.id.replace('steam-', '');
      setImgSrc(`https://cdn.akamai.steamstatic.com/steam/apps/${appId}/capsule_616x353.jpg`);
    } else {
      setImgSrc(
        `https://via.placeholder.com/600x900/1a1625/606060?text=${encodeURIComponent(game.title)}`,
      );
    }
  };

  return (
    <Link
      href={`/games/details/${game.id}`}
      className="bg-background-surface rounded-xl overflow-hidden hover:scale-[1.02] transition-transform duration-300 border border-[#2C2940] flex flex-col group cursor-pointer shadow-lg hover:shadow-primary/20 hover:border-primary/50 relative"
    >
      <div className="relative aspect-2/3 w-full bg-[#1a1625] overflow-hidden">
        <Image
          src={imgSrc}
          alt={game.title}
          fill
          className="object-cover"
          sizes="(max-width: 768px) 50vw, (max-width: 1200px) 33vw, 250px"
          onError={handleError}
          unoptimized={true}
        />
        <div className="absolute inset-0 bg-linear-to-t from-background-base/90 to-transparent opacity-0 group-hover:opacity-100 transition-opacity flex items-end p-4">
          <span className="btn btn-primary btn-sm w-full">Ver Detalhes</span>
        </div>
      </div>
      <div className="p-4 flex flex-col gap-2">
        <h4 className="font-bold text-gray-100 truncate text-sm">{game.title}</h4>
        <div className="flex items-center gap-2">
          <span className="text-xs text-gray-500 uppercase tracking-wider font-medium">status</span>
          <span
            className={`text-[10px] px-2 py-0.5 rounded-full border ${statusColor} bg-black/40 uppercase font-bold`}
          >
            {game.status}
          </span>
        </div>
      </div>
    </Link>
  );
}
