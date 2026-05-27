import { expect, test, describe } from 'bun:test';
import { render, screen } from '@testing-library/react';
import { HomeFooter } from './HomeFooter';
import React from 'react';

describe('HomeFooter Component', () => {
  test('renders footer brand and links successfully', () => {
    render(<HomeFooter />);
    
    expect(screen.getByText('YNG')).toBeTruthy();
    expect(screen.getByText('© 2026 Your Next Game. Todos os direitos reservados.')).toBeTruthy();
    expect(screen.getByText('Política de Privacidade')).toBeTruthy();
    expect(screen.getByText('Termos de Serviço')).toBeTruthy();
    expect(screen.getByText('Contato')).toBeTruthy();
  });
});
