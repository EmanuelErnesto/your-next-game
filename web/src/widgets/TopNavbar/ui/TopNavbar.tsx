'use client';

import { Bell, LogOut, Menu, Settings } from 'lucide-react';
import Image from 'next/image';

import { useAuth } from '@shared/auth';

export function TopNavbar() {
  const { session, signOut, isLoading } = useAuth();

  return (
    <header className="h-16 border-b border-[#2C2940] bg-background-base flex items-center justify-between px-6 shrink-0">
      <div className="flex items-center gap-4">
        <button className="btn btn-ghost btn-sm btn-square text-gray-400 hover:text-white lg:hidden">
          <Menu size={20} />
        </button>
        <div className="flex items-center gap-3">
          <Image
            src="/logo.png"
            alt="Your Next Game"
            width={32}
            height={32}
            className="rounded shadow-lg shadow-primary/20"
          />
          <span className="font-bold text-lg tracking-tight">Your Next Game</span>
        </div>
      </div>

      <div className="flex items-center gap-2">
        {!isLoading && session && (
          <div className="flex items-center gap-3 px-3 py-1 bg-background-surface rounded-full border border-[#2C2940] mr-4">
            {session.image && (
              <Image
                src={session.image}
                alt={session.name || 'User'}
                width={24}
                height={24}
                className="rounded-full"
              />
            )}
            <span className="text-xs font-semibold text-gray-300">{session.name}</span>
          </div>
        )}

        <button className="btn btn-ghost btn-sm btn-square text-gray-400 hover:text-white">
          <Bell size={20} />
        </button>
        <button className="btn btn-ghost btn-sm btn-square text-gray-400 hover:text-white">
          <Settings size={20} />
        </button>
        <div className="w-px h-6 bg-[#2C2940] mx-2"></div>
        <button
          type="button"
          onClick={() => {
            void signOut();
          }}
          className="btn btn-ghost btn-sm btn-square text-gray-400 hover:text-white"
        >
          <LogOut size={20} />
        </button>
      </div>
    </header>
  );
}
