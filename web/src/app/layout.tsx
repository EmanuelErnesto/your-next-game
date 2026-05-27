import type { Metadata } from 'next';
import { Geist, Geist_Mono } from 'next/font/google';
import './globals.css';
import { env } from '@shared/config/env';
import { getSession } from '@shared/lib/auth';
import { Providers } from './providers';

const geistSans = Geist({
  variable: '--font-geist-sans',
  subsets: ['latin'],
});

const geistMono = Geist_Mono({
  variable: '--font-geist-mono',
  subsets: ['latin'],
});

export const metadata: Metadata = {
  title: 'Your Next Game',
  description: 'Pare de navegar. Comece a jogar.',
  icons: {
    icon: '/icon.png',
    shortcut: '/favicon.ico',
    apple: '/icon.png',
  },
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const initialSession = await getSession();

  return (
    <html lang="pt-BR" data-theme="dark">
      <body className={`${geistSans.variable} ${geistMono.variable} antialiased`}>
        <Providers backendBaseUrl={env.BACKEND_BASE_URL} initialSession={initialSession}>
          {children}
        </Providers>
      </body>
    </html>
  );
}
