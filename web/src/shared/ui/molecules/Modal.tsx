import type { ReactNode } from 'react';

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title?: string;
  children: ReactNode;
}

export function Modal({ isOpen, onClose, title, children }: ModalProps) {
  if (!isOpen) return null;

  return (
    <div className="modal modal-open bg-black/60 backdrop-blur-sm" role="dialog">
      <div className="modal-box border border-primary/30 shadow-[0_0_15px_rgba(138,43,226,0.3)] bg-background-surface">
        <form method="dialog">
          <button
            className="btn btn-sm btn-circle btn-ghost absolute right-2 top-2"
            onClick={onClose}
          >
            ✕
          </button>
        </form>
        {title && <h3 className="font-bold text-lg mb-4">{title}</h3>}
        {children}
      </div>
      <div className="modal-backdrop" onClick={onClose}>
        <button className="cursor-default">close</button>
      </div>
    </div>
  );
}
