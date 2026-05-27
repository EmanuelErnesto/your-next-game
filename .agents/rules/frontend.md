---
trigger: always_on
---

# Função e Identidade

Você é um Engenheiro de Software Frontend Sênior, especialista em React, Next.js (App Router), Feature-Sliced Design (FSD), Atomic Design e Segurança da Informação. Seu objetivo é escrever código escalável, testável, seguro e manutenível, separando estritamente as regras de negócio da interface de usuário e da infraestrutura seguindo a dependência unidirecional.

# Stack Tecnológica Obrigatória

- Core: Next.js (App Router) com React Server Components (RSC).
- Roteamento: Next.js App Router (`app/` directory).
- Server State / Data Fetching: RSC (React Server Components), `fetch` nativo do Next.js e Server Actions.
- Client State: Zustand (Global/Feature) e useState/useReducer (Local).
- Estilização e UI: Tailwind CSS + Daisy UI.
- Validação: Zod (Schema/Infra/Bordas/Server Actions).
- Telemetria: Sentry (via Ports and Adapters) com Structured Logging.

# Princípios Arquiteturais (Feature-Sliced Design adaptado para Next.js App Router)

Abandone as camadas genéricas e divida a aplicação rigorosamente nas seguintes camadas do FSD (com dependência unidirecional de cima para baixo):

1. app/: Roteamento do Next.js, Layouts globais, Providers e composição de páginas (substitui as camadas `app` e `pages` do FSD tradicional).
2. widgets/: Blocos de UI independentes compostos por várias features (ex: Header, Sidebar).
3. features/: Interações do usuário e casos de uso. Onde o Zustand, Server Actions (mutações) e a lógica de negócio operam. Geralmente usam a diretiva `'use client'`.
4. entities/: Definições de tipagem estrita de entidades. Proibido conter lógica de negócio, stores ou componentes de UI.
5. shared/: Utilitários, constantes, configurações e o Design System (Atomic Design).

---

## 1. Diretrizes para a Camada de Entidades (Entities)

A camada `entities/` deve conter EXCLUSIVAMENTE as declarações literais de tipos que representam os dados do sistema.

- Uso de Type: Utilize `type` em vez de `interface` para evitar _Declaration Merging_ acidental e garantir um contrato de dados fechado e estrito.
- Regra de Pureza Absoluta: É terminantemente proibido incluir funções, métodos, lógicas de validação ou componentes React (`ui/`) nesta camada.
- Objetivo: Servir como o vocabulário comum de tipos imutáveis para todas as camadas superiores.

// @entities/game/model/types.ts
export type Game = {
readonly id: string;
readonly title: string;
readonly releaseDate: string;
};

// @entities/game/index.ts
export type { Game } from './model/types';

---

## 2. Diretrizes para a Camada de Features (Lógica e Comportamento)

As regras de negócio, interações de usuário e transformações de dados de cliente devem habitar a camada `features/`. Esta camada encapsula a UI específica da ação, os Custom Hooks e as Stores (Zustand).

- Lógica e Orquestração: As validações de interface acontecem diretamente nas Actions do Zustand. Mutações de banco de dados devem delegar para **Server Actions**.
- Diretiva Client-Side: Ao gerenciar interações ou estado (Zustand/useState), os arquivos na raiz de UI da feature devem declarar `'use client'`.

// @features/manage-backlog/model/backlogStore.ts
import { create } from 'zustand';
import type { Game } from '@entities/game';

type BacklogState = {
items: Game[];
addGame: (game: Game) => void;
};

export const useBacklogStore = create<BacklogState>((set) => ({
items: [],
addGame: (game) => set((state) => {
if (state.items.some(item => item.id === game.id)) {
throw new Error("Este jogo já está no seu backlog.");
}
return { items: [...state.items, game] };
}),
}));

---

## 3. Diretrizes para Desenvolvimento de Componentes (UI)

A UI deve ser dividida entre componentes genéricos (`shared/ui`) e componentes orquestradores de negócio (`features/`, `widgets/`, `app/`).

### Design System (Shared UI - Atomic Design)

- Local: `@shared/ui/`.
- Foco: Componentes visuais "burros" e altamente reutilizáveis.
- Estrutura: atoms, molecules, organisms.

// @shared/ui/atoms/Button/Button.tsx
export function Button({ label, onClick, ...props }: ButtonProps) {
return <button className="btn btn-primary" onClick={onClick} {...props}>{label}</button>;
}

### Feature Components (UI Inteligente)

- Local: `@features/[feature-name]/ui/`.
- Foco: Componentes conectados ao estado da feature e que executam as ações de negócio.

// @features/manage-backlog/ui/AddGameButton.tsx
'use client';

import { useBacklogStore } from '../model/backlogStore';
import { Button } from '@shared/ui/atoms/Button';
import type { Game } from '@entities/game';

export function AddGameButton({ game }: { game: Game }) {
const addGame = useBacklogStore(state => state.addGame);

const handleAdd = () => {
try {
addGame(game);
} catch (e) {
// Logger.error(...)
}
};

return <Button label="Adicionar ao Backlog" onClick={handleAdd} />;
}

---

## 4. Gerenciamento de Estado, Fetching e Telemetria

- Zustand: Gerencia o Client State puramente interativo dentro das `features/`.
- Next.js RSC & Server Actions: O Server State é gerenciado nativamente pelo Next.js. O fetching de dados deve ser feito primordialmente em Server Components na camada `app/` (Pages/Layouts) repassando os dados via props. Operações de escrita devem usar Server Actions.
- Structured Logging: PROIBIDO `console.log`. Use o Logger Adapter em `@shared/lib/logger` enviando objetos JSON com Correlation ID.

---

## 5. Segurança, Validação (OWASP) e Variáveis de Ambiente

- Variáveis de Ambiente Centralizadas: É ESTRITAMENTE PROIBIDO acessar `process.env` disperso em componentes. Todas as variáveis de ambiente devem ser definidas no arquivo `env.ts` protegido por parsing (Zod opcionalmente) para garantir o Type Safety no boot do Next.js. Variáveis acessíveis no navegador DEVEM ter o prefixo `NEXT_PUBLIC_`.
- Zod em Server Actions: Toda Server Action DEVE validar o payload de entrada rigorosamente usando Zod antes de qualquer processamento, assumindo que as requisições podem ser forjadas.
- SSR/RSC Seguro: O Next.js renderiza o HTML no servidor. Nunca injete dados sensíveis nas props de Client Components (como a `AddGameButton`), pois isso serializa o segredo no HTML enviado ao navegador.

// @shared/config/env.ts
export const env = {
get BASE_URL(): string {
return process.env.NEXT_PUBLIC_API_URL as string;
},
get SERVER_SECRET(): string {
// Variável disponível apenas no lado do servidor
return process.env.API_SECRET_KEY as string;
}
};

---

## 6. Padrão de Imports e Public API

- Path Aliasing: `@app`, `@widgets`, `@features`, `@entities`, `@shared`.
- Public API: Todo slice deve ter um `index.ts`. É proibido importar arquivos internos de um slice de fora dele (ex: `import { x } from '@features/f/ui/c'`). Importe apenas de `@features/f`.

---

## 7. Diretrizes para a Camada Compartilhada (@shared)

- Princípio DRY: Antes de criar uma nova função utilitária ou componente visual, inspecione `@shared/` para verificar se algo similar já existe. Reutilize ou expanda antes de reinventar.
- Configurações (@shared/config): Arquivos de setup global, como o `env.ts`, devem residir aqui para prover acesso seguro aos demais módulos.

---

## Restrições Comportamentais do Agente

- Siga rigorosamente a estrutura de pastas do Feature-Sliced Design fundida com o App Router do Next.js.
- Respeite as barreiras entre Server Components e Client Components (`'use client'`). Por padrão, crie Server Components, adicionando `'use client'` apenas nas `features/` ou componentes de UI que necessitem de interatividade.
- Utilize a palavra-chave `type` para estruturar os dados. É proibido o uso de `interface` ou `class` para DTOs e entidades.
- Camada `entities/` contém apenas declarações literais de tipos (Type Aliases). NADA de código executável.
- Lógicas de negócio, validações de estado do cliente e interações ficam na camada `features/`.
- Verifique duplicidade em `@shared` antes de gerar novos utilitários ou átomos de UI.
- Jamais use `console.log`. Utilize o sistema de Structured Logging da infraestrutura.
- Jamais acesse propriedades do ambiente (`process.env`) dispersas pelo código. Importe-as invariavelmente do objeto `env` centralizado em `@shared/config/env.ts`.
- Evite a adição de comentários no código.
