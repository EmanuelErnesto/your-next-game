import { expect, test, mock, describe, beforeEach, afterEach } from 'bun:test';
import { render, screen, fireEvent, cleanup } from '@testing-library/react';
import { HomeHeader } from './HomeHeader';
import React from 'react';

let mockAuth: any = null;

mock.module('@shared/auth', () => ({
  useAuth: () => mockAuth,
}));

describe('HomeHeader Component', () => {
  beforeEach(() => {
    mockAuth = {
      session: null,
      signIn: mock(() => {}),
      signOut: mock(() => {}),
      isLoading: false,
    };
    window.location.href = 'http://localhost:3001';
  });

  afterEach(() => {
    cleanup();
    document.body.innerHTML = '';
  });

  test('renders header with brand name and links', () => {
    render(<HomeHeader />);
    expect(screen.getByText('Your Next Game')).toBeTruthy();
    expect(screen.getByText('Como Funciona')).toBeTruthy();
    expect(screen.getByText('Features')).toBeTruthy();
  });

  test('displays Entrar com Steam button when unauthenticated', () => {
    render(<HomeHeader />);
    const loginButton = screen.getByText('Entrar com Steam');
    expect(loginButton).toBeTruthy();

    fireEvent.click(loginButton);
    expect(mockAuth.signIn).toHaveBeenCalled();
  });

  test('displays user info, Sair button, and Dashboard when authenticated', () => {
    mockAuth.session = {
      name: 'Gamer 123',
      image: 'https://steamcdn.com/avatar.jpg',
    };

    render(<HomeHeader />);

    // Brand and nav should still be there
    expect(screen.getByText('Your Next Game')).toBeTruthy();

    // User name and image
    expect(screen.getByText('Gamer 123')).toBeTruthy();
    const avatar = screen.getByAltText('Gamer 123') as HTMLImageElement;
    expect(avatar.src).toBe('https://steamcdn.com/avatar.jpg');

    // Dashboard and Sair buttons
    expect(screen.getByText('Dashboard')).toBeTruthy();
    const logoutButton = screen.getByText('Sair');
    expect(logoutButton).toBeTruthy();

    fireEvent.click(logoutButton);
    expect(mockAuth.signOut).toHaveBeenCalled();
  });
});
