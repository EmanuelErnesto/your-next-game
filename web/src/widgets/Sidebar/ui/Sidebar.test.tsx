import { expect, test, mock, describe, beforeEach } from 'bun:test';
import { render, screen } from '@testing-library/react';
import { Sidebar } from './Sidebar';
import React from 'react';

let mockAuth: any = {
  session: null,
  isLoading: false,
};

mock.module('@shared/auth', () => ({
  useAuth: () => mockAuth,
}));

describe('Sidebar Component', () => {
  beforeEach(() => {
    mockAuth = {
      session: null,
      isLoading: false,
    };
  });

  test('renders user info when provided', () => {
    mockAuth.session = {
      name: 'João Silva',
      image: 'https://example.com/avatar.jpg',
    };
    
    render(<Sidebar />);
    
    expect(screen.getByText('João Silva')).toBeTruthy();
    const avatar = screen.getByAltText('João Silva') as HTMLImageElement;
    expect(avatar.src).toBe('https://example.com/avatar.jpg');

    // Filters should be visible
    expect(screen.getByText('Filtros')).toBeTruthy();
    expect(screen.getByText('Gênero')).toBeTruthy();
    expect(screen.getByText('Plataforma')).toBeTruthy();
    expect(screen.getByText('Status')).toBeTruthy();
  });

  test('falls back to default info when user is not provided', () => {
    render(<Sidebar />);
    
    expect(screen.getByText('Convidado')).toBeTruthy();
    const avatar = screen.getByAltText('Convidado') as HTMLImageElement;
    expect(avatar.src).toBe('https://i.pravatar.cc/150?img=11');
  });
});
