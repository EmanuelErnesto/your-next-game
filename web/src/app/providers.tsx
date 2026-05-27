'use client';

import { ReactNode } from 'react';
import { AuthProvider } from '@shared/auth';
import type { UserSession } from '@shared/backend/types';

type ProvidersProps = {
  children: ReactNode;
  backendBaseUrl: string;
  initialSession: UserSession | null;
};

export function Providers({ children, backendBaseUrl, initialSession }: ProvidersProps) {
  return (
    <AuthProvider backendBaseUrl={backendBaseUrl} initialSession={initialSession}>
      {children}
    </AuthProvider>
  );
}
