import { expect, test, mock, describe } from 'bun:test';
import { render, screen, fireEvent } from '@testing-library/react';
import { Modal } from './Modal';
import React from 'react';

describe('Modal Component', () => {
  test('does not render when isOpen is false', () => {
    const handleClose = mock(() => {});
    render(
      <Modal isOpen={false} onClose={handleClose}>
        <div>Conteúdo do Modal</div>
      </Modal>,
    );
    const content = screen.queryByText('Conteúdo do Modal');
    expect(content).toBeNull();
  });

  test('renders title and children when isOpen is true', () => {
    const handleClose = mock(() => {});
    render(
      <Modal isOpen={true} onClose={handleClose} title="Título do Teste">
        <div>Conteúdo do Modal</div>
      </Modal>,
    );

    const title = screen.getByText('Título do Teste');
    const content = screen.getByText('Conteúdo do Modal');

    expect(title).toBeTruthy();
    expect(content).toBeTruthy();
  });

  test('calls onClose when close button is clicked', () => {
    const handleClose = mock(() => {});
    render(
      <Modal isOpen={true} onClose={handleClose}>
        <div>Conteúdo</div>
      </Modal>,
    );

    const closeButton = screen.getByText('✕');
    fireEvent.click(closeButton);

    expect(handleClose).toHaveBeenCalled();
  });

  test('calls onClose when backdrop is clicked', () => {
    const handleClose = mock(() => {});
    render(
      <Modal isOpen={true} onClose={handleClose}>
        <div>Conteúdo</div>
      </Modal>,
    );

    const backdrop = screen.getByText('close');
    fireEvent.click(backdrop);

    expect(handleClose).toHaveBeenCalled();
  });
});
