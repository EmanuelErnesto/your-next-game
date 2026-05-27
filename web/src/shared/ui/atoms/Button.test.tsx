import { expect, test, mock, describe } from 'bun:test';
import { render, screen, fireEvent } from '@testing-library/react';
import { Button } from './Button';
import React from 'react';

describe('Button Component', () => {
  test('renders successfully with a label', () => {
    render(<Button label="Clique aqui" />);
    const button = screen.getByText('Clique aqui');
    expect(button).toBeTruthy();
    expect(button.className).toContain('btn-primary');
    expect(button.className).toContain('btn-md');
  });

  test('renders successfully with children', () => {
    render(<Button>Texto Filho</Button>);
    const button = screen.getByText('Texto Filho');
    expect(button).toBeTruthy();
  });

  test('applies custom variant and size classes', () => {
    render(<Button label="Test" variant="secondary" size="sm" isFullWidth />);
    const button = screen.getByText('Test');
    expect(button.className).toContain('btn-secondary');
    expect(button.className).toContain('btn-sm');
    expect(button.className).toContain('w-full');
  });

  test('renders icon when provided', () => {
    const testIcon = <span data-testid="test-icon">icon</span>;
    render(<Button label="Test" icon={testIcon} />);
    const icon = screen.getByTestId('test-icon');
    expect(icon).toBeTruthy();
  });

  test('handles click events', () => {
    const handleClick = mock(() => {});
    render(<Button label="Click Me" onClick={handleClick} />);
    const button = screen.getByText('Click Me');
    fireEvent.click(button);
    expect(handleClick).toHaveBeenCalled();
  });

  test('respects disabled state', () => {
    const handleClick = mock(() => {});
    render(<Button label="Disabled" onClick={handleClick} disabled />);
    const button = screen.getByText('Disabled');
    fireEvent.click(button);
    expect(handleClick).not.toHaveBeenCalled();
  });
});
