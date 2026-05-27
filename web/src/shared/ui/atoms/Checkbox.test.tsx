import { expect, test, mock, describe } from 'bun:test';
import { render, screen, fireEvent } from '@testing-library/react';
import { Checkbox } from './Checkbox';
import React from 'react';

describe('Checkbox Component', () => {
  test('renders with the provided label', () => {
    render(<Checkbox label="Aceito os termos" />);
    const labelText = screen.getByText('Aceito os termos');
    expect(labelText).toBeTruthy();
  });

  test('displays checked state properly', () => {
    const { rerender } = render(<Checkbox label="Check me" checked={true} readOnly />);
    const checkbox = screen.getByRole('checkbox') as HTMLInputElement;
    expect(checkbox.checked).toBe(true);

    rerender(<Checkbox label="Check me" checked={false} readOnly />);
    expect(checkbox.checked).toBe(false);
  });

  test('triggers onChange handler when clicked', () => {
    const handleChange = mock(() => {});
    render(<Checkbox label="Click me" onChange={handleChange} />);
    const checkbox = screen.getByRole('checkbox');
    fireEvent.click(checkbox);
    expect(handleChange).toHaveBeenCalled();
  });
});
