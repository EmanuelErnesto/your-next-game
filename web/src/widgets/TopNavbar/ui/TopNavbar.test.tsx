import { expect, test, mock, describe, beforeEach } from 'bun:test';
import { render, screen, fireEvent } from '@testing-library/react';
import { TopNavbar } from './TopNavbar';
import React from 'react';

let mockAuth: any = null;

mock.module('@shared/auth', () => ({
  useAuth: () => mockAuth,
}));

describe('TopNavbar Component', () => {
  beforeEach(() => {
    mockAuth = {
      session: null,
      signOut: mock(() => {}),
      signIn: mock(() => {}),
      isLoading: false,
    };
  });

  test('renders top navbar brand and icons', () => {
    render(<TopNavbar />);
    expect(screen.getByText('Your Next Game')).toBeTruthy();
    // It should render bell, settings and logout buttons
    expect(screen.getByAltText('Your Next Game')).toBeTruthy();
  });

  test('displays user info badge when user is logged in', () => {
    mockAuth.session = {
      name: 'Ana Souza',
      image: 'https://example.com/ana.jpg',
    };

    render(<TopNavbar />);
    expect(screen.getByText('Ana Souza')).toBeTruthy();
    const avatar = screen.getByAltText('Ana Souza') as HTMLImageElement;
    expect(avatar.src).toBe('https://example.com/ana.jpg');
  });

  test('calls signOut when logout button is clicked', () => {
    render(<TopNavbar />);
    
    // Find logout button by icon tag / button type
    const buttons = screen.getAllByRole('button');
    // The last button is logout
    const logoutButton = buttons[buttons.length - 1];

    fireEvent.click(logoutButton);
    expect(mockAuth.signOut).toHaveBeenCalled();
  });
});
