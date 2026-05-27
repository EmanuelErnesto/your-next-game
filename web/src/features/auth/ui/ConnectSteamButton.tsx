'use client';

import Image from 'next/image';
import { useAuth } from '@shared/auth';

type ConnectSteamButtonProps = {
  className?: string;
};

export function ConnectSteamButton({ className }: ConnectSteamButtonProps) {
  const { signIn } = useAuth();

  return (
    <button onClick={signIn} className={`btn btn-primary font-bold ${className ?? ''}`}>
      <Image src="/steam-icon.svg" alt="Steam Logo" width={24} height={24} className="mr-2" />
      Conectar Biblioteca Steam
    </button>
  );
}
