import { expect, test, mock, describe, beforeEach } from 'bun:test';
import { render, screen, fireEvent } from '@testing-library/react';
import { ConnectSteamButton } from './ConnectSteamButton';
import React from 'react';

let mockAuth: any = null;

mock.module('@shared/auth', () => ({
  useAuth: () => mockAuth,
}));

describe('ConnectSteamButton Component', () => {
  beforeEach(() => {
    mockAuth = {
      signIn: mock(() => {}),
      signOut: mock(() => {}),
      session: null,
      isLoading: false,
    };
  });

  test('renders button with Steam connect text', () => {
    render(<ConnectSteamButton />);
    const button = screen.getByText('Conectar Biblioteca Steam');
    expect(button).toBeTruthy();
  });

  test('redirects to steam login api when clicked', () => {
    render(<ConnectSteamButton />);
    const button = screen.getByText('Conectar Biblioteca Steam');
    fireEvent.click(button);
    expect(mockAuth.signIn).toHaveBeenCalled();
  });
});
