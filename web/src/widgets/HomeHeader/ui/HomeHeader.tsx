"use client";

import { Gamepad2, Menu } from 'lucide-react';
import Image from 'next/image';

import { useAuth } from '@shared/auth';

export function HomeHeader() {
  const { session, signIn, signOut, isLoading } = useAuth();

  return (
    <header className="fixed top-0 left-0 right-0 z-50 transition-all duration-300">
      <div className="absolute inset-0 bg-background-base/80 backdrop-blur-md border-b border-white/5"></div>

      <div className="container mx-auto px-4 lg:px-8 relative">
        <div className="flex items-center justify-between h-20">
          {/* Logo */}
          <div className="flex items-center gap-3">
            <div className="relative w-10 h-10 rounded-xl overflow-hidden shadow-lg shadow-primary/20 flex items-center justify-center bg-gradient-to-br from-primary to-secondary">
              <Gamepad2 className="text-white w-6 h-6" />
            </div>
            <span className="font-bold text-xl tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-white to-gray-400">
              Your Next Game
            </span>
          </div>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex items-center gap-8">
            <a
              href="#features"
              className="text-sm font-medium text-gray-300 hover:text-white transition-colors"
            >
              Features
            </a>
            <a
              href="#how-it-works"
              className="text-sm font-medium text-gray-300 hover:text-white transition-colors"
            >
              Como Funciona
            </a>
            <a
              href="#faq"
              className="text-sm font-medium text-gray-300 hover:text-white transition-colors"
            >
              FAQ
            </a>
          </nav>

          {/* User / CTA */}
          <div className="hidden md:flex items-center gap-4">
            {isLoading ? null : session ? (
              <div className="flex items-center gap-4">
                <a
                  href="/dashboard"
                  className="text-sm font-medium text-gray-300 hover:text-white transition-colors"
                >
                  Dashboard
                </a>
                <div className="flex items-center gap-3 px-4 py-2 bg-white/5 rounded-full border border-white/10 hover:bg-white/10 transition-colors cursor-pointer group">
                  <div className="w-8 h-8 rounded-full bg-gradient-to-br from-primary to-secondary overflow-hidden">
                    {session.image ? (
                      <Image
                        src={session.image}
                        alt={session.name || 'User'}
                        width={32}
                        height={32}
                      />
                    ) : (
                      <span className="w-full h-full flex items-center justify-center text-white font-bold text-sm">
                        {session.name?.charAt(0) || 'U'}
                      </span>
                    )}
                  </div>
                  <span className="text-sm font-semibold text-gray-200">
                    {session.name}
                  </span>
                </div>
                <button onClick={() => void signOut()} className="btn btn-ghost btn-sm text-gray-400 hover:text-white">
                  Sair
                </button>
              </div>
            ) : (
              <button className="btn btn-primary rounded-full px-6" onClick={signIn}>
                <Gamepad2 className="w-4 h-4 mr-2" />
                Entrar com Steam
              </button>
            )}
          </div>

          {/* Mobile Menu Button */}
          <button className="md:hidden btn btn-ghost btn-circle text-gray-300 hover:text-white">
            <Menu className="w-6 h-6" />
          </button>
        </div>
      </div>
    </header>
  );
}
