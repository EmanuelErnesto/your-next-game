'use client';

import type { Game } from '@entities/game';
import { Button } from '@shared/ui/atoms/Button';
import { Modal } from '@shared/ui/molecules/Modal';
import { Battery, Brain, Clock, Moon, Sparkles, Theater } from 'lucide-react';
import Image from 'next/image';
import { useVerdictStore } from '../model/verdictStore';

interface VerdictFlowProps {
  games: Game[];
}

export function VerdictFlow({ games }: VerdictFlowProps) {
  const {
    step,
    mood,
    timeAvailable,
    genres,
    openVibeConfig,
    close,
    setTimeAvailable,
    setMood,
    toggleGenre,
    consultAI,
    selectedGame,
    acceptChallenge,
    changeMind,
    justification,
  } = useVerdictStore();

  return (
    <>
      <div className="fixed bottom-8 right-8 z-50">
        <Button
          variant="primary"
          size="lg"
          onClick={openVibeConfig}
          className="shadow-xl shadow-primary/20 rounded-full pr-6 pl-5 gap-3 border border-primary/50 text-white font-semibold transform hover:scale-105 transition-all bg-background-base"
        >
          <Sparkles size={20} className="text-cyan-400" />
          Obter Veredito
        </Button>
      </div>

      <Modal isOpen={step !== 'CLOSED'} onClose={close}>
        {step === 'CONFIG_VIBE' && (
          <div className="flex flex-col gap-6">
            <h2 className="text-2xl font-bold">Nova recomendação de Jogo</h2>

            <div className="flex flex-col gap-2">
              <label className="text-sm font-semibold flex items-center justify-between">
                <span>Tempo Disponível (horas)</span>
              </label>
              <div className="flex items-center gap-4 text-gray-400">
                <Clock size={16} />
                <input
                  type="number"
                  min={1}
                  max={24}
                  value={timeAvailable}
                  onChange={(e) => {
                    let val = Number(e.target.value);
                    if (val > 24) val = 24;
                    if (val < 1 && e.target.value !== '') val = 1;
                    setTimeAvailable(val);
                  }}
                  className="input input-bordered input-sm input-primary w-full max-w-[100px] text-center"
                />
              </div>
            </div>

            <div className="flex flex-col gap-3">
              <label className="text-sm font-semibold">Humor</label>
              <div className="grid grid-cols-2 gap-3">
                <MoodButton
                  label="Energético"
                  icon={<Battery className="text-green-400" size={18} />}
                  isSelected={mood === 'Energético'}
                  onClick={() => setMood('Energético')}
                />
                <MoodButton
                  label="Relaxante"
                  icon={<Moon className="text-yellow-400" size={18} />}
                  isSelected={mood === 'Relaxante'}
                  onClick={() => setMood('Relaxante')}
                />
                <MoodButton
                  label="Desafiador"
                  icon={<Brain className="text-pink-400" size={18} />}
                  isSelected={mood === 'Desafiador'}
                  onClick={() => setMood('Desafiador')}
                />
                <MoodButton
                  label="Narrativo"
                  icon={<Theater className="text-orange-400" size={18} />}
                  isSelected={mood === 'Narrativo'}
                  onClick={() => setMood('Narrativo')}
                />
              </div>
            </div>

            <Button
              variant="outline"
              className="mt-4 border-cyan-400/50 text-cyan-400 hover:bg-cyan-400/10 hover:border-cyan-400 font-bold"
              onClick={() => consultAI(games)}
              disabled={!mood}
            >
              Gerar Recomendação
            </Button>
          </div>
        )}

        {step === 'LOADING' && (
          <div className="flex flex-col items-center justify-center gap-6 py-12">
            <Sparkles size={48} className="text-cyan-400 animate-pulse" />
            <h3 className="text-xl font-bold text-white text-center">
              A IA está analisando sua biblioteca...
            </h3>
            <p className="text-gray-400 text-sm text-center">Isso pode levar alguns segundos.</p>
          </div>
        )}

        {step === 'RESULT' && selectedGame && (
          <div className="flex flex-col items-center gap-6 text-center">
            <div className="relative w-40 aspect-[2/3] rounded-xl overflow-hidden shadow-[0_0_20px_rgba(138,43,226,0.6)] border border-primary/50">
              <Image
                src={selectedGame.coverUrl}
                alt={selectedGame.title}
                fill
                className="object-cover"
              />
            </div>

            <div className="flex flex-col items-center gap-1">
              <span className="text-sm text-gray-400">Seu jogo ideal para hoje é...</span>
              <h2 className="text-3xl font-extrabold text-white">{selectedGame.title}</h2>
            </div>

            <div className="bg-[#1A1829] border border-[#2C2940] rounded-xl p-4 text-left w-full text-sm text-gray-300">
              <h4 className="font-bold text-white mb-2 text-sm">Justificativa da IA</h4>
              <p>
                {justification ||
                  'Nossa IA encontrou este jogo perfeito para você com base no seu humor!'}
              </p>
            </div>

            <div className="flex gap-4 w-full">
              <Button
                variant="primary"
                className="flex-1 bg-green-600 hover:bg-green-700 border-none text-white"
                onClick={acceptChallenge}
              >
                ✓ Aceitar Desafio
              </Button>
              <Button
                variant="ghost"
                className="flex-1 border border-gray-600"
                onClick={changeMind}
              >
                ✕ Mudar de Ideia
              </Button>
            </div>
          </div>
        )}
      </Modal>
    </>
  );
}

function MoodButton({
  label,
  icon,
  isSelected,
  onClick,
}: {
  label: string;
  icon: React.ReactNode;
  isSelected: boolean;
  onClick: () => void;
}) {
  return (
    <button
      className={`flex items-center gap-3 p-3 rounded-xl border transition-all ${
        isSelected
          ? 'border-cyan-400 bg-cyan-400/10 shadow-[0_0_10px_rgba(34,211,238,0.2)]'
          : 'border-[#2C2940] bg-background-base hover:border-gray-500'
      }`}
      onClick={onClick}
    >
      {icon}
      <span className="text-sm font-medium">{label}</span>
    </button>
  );
}
