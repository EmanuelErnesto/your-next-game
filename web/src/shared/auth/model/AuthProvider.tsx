'use client';

import { useRouter } from 'next/navigation';
import type { ReactNode } from 'react';
import { createContext, useCallback, useContext, useEffect, useState } from 'react';

import type { UserSession } from '@shared/backend/types';
import { clearAuthSession, fetchAuthSession } from '../api';

type AuthStatus = 'loading' | 'authenticated' | 'unauthenticated';

type AuthContextValue = {
  session: UserSession | null;
  status: AuthStatus;
  isAuthenticated: boolean;
  isLoading: boolean;
  refreshSession: () => Promise<void>;
  signIn: () => void;
  signOut: () => Promise<void>;
};

type AuthProviderProps = {
  children: ReactNode;
  backendBaseUrl: string;
  initialSession: UserSession | null;
};

const AUTH_REVALIDATE_MS = 5 * 60 * 1000;

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children, backendBaseUrl, initialSession }: AuthProviderProps) {
  const router = useRouter();
  const [session, setSession] = useState<UserSession | null>(initialSession);
  const [status, setStatus] = useState<AuthStatus>(initialSession ? 'authenticated' : 'loading');

  const refreshSession = useCallback(async () => {
    try {
      const nextSession = await fetchAuthSession(backendBaseUrl);
      setSession(nextSession);
      setStatus(nextSession ? 'authenticated' : 'unauthenticated');
    } catch {
      setSession(null);
      setStatus('unauthenticated');
    }
  }, [backendBaseUrl]);

  useEffect(() => {
    void refreshSession();
  }, [refreshSession]);

  useEffect(() => {
    const handleFocus = () => {
      void refreshSession();
    };

    window.addEventListener('focus', handleFocus);
    const intervalId = window.setInterval(() => {
      void refreshSession();
    }, AUTH_REVALIDATE_MS);

    return () => {
      window.removeEventListener('focus', handleFocus);
      window.clearInterval(intervalId);
    };
  }, [refreshSession]);

  const signIn = useCallback(() => {
    window.location.assign(new URL('/v1/auth/steam/login', backendBaseUrl).toString());
  }, [backendBaseUrl]);

  const signOut = useCallback(async () => {
    await clearAuthSession(backendBaseUrl);
    setSession(null);
    setStatus('unauthenticated');
    router.replace('/');
  }, [backendBaseUrl, router]);

  const value: AuthContextValue = {
    session,
    status,
    isAuthenticated: session !== null,
    isLoading: status === 'loading',
    refreshSession,
    signIn,
    signOut,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }

  return context;
}
