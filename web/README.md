# Your Next Game - Frontend 💻

Esta é a aplicação frontend desenvolvida com as tecnologias mais modernas do ecossistema React/Next.js. Ela serve como a interface principal de gerenciamento de backlog para os usuários.

## 🚀 Tecnologias Utilizadas

- **Next.js (App Router)** com React Server Components (RSC).
- **Tailwind CSS + Daisy UI** para estilização flexível, moderna e rica.
- **Zustand** para gerenciamento de estado global interativo (Client State).
- **Zod** para validação estrita de contratos nas bordas da aplicação.
- **OpenTelemetry & Pino** para telemetria robusta e logs estruturados em formato JSON.
- **Biome** como linter e formatador de código ultrarrápido.

---

## 📐 Diretrizes Arquiteturais (Feature-Sliced Design)

Seguimos estritamente o padrão **Feature-Sliced Design (FSD)** adaptado para o App Router do Next.js, mantendo uma clara separação de conceitos com dependência unidirecional:

1. **`app/`**: Roteamento nativo do Next.js, layouts globais, providers e composição de páginas.
2. **`widgets/`**: Blocos independentes de UI compostos por várias features (ex: Header, Sidebar).
3. **`features/`**: Casos de uso e interações diretas do usuário. Contém lógica de negócio, stores Zustand e ações do cliente.
4. **`entities/`**: Contratos e definições puras de tipos. Proibido conter lógica, stores ou componentes de UI. Usa exclusivamente `type` para tipagens fechadas.
5. **`shared/`**: Utilitários, constantes, chamadas de infraestrutura (como `env.ts` para variáveis de ambiente) e o Design System baseado em **Atomic Design** (`shared/ui/atoms`, `shared/ui/molecules`, etc.).

---

## ⚙️ Variáveis de Ambiente

As variáveis de ambiente são lidas e validadas estritamente via `@shared/config/env.ts`.
Configure seu arquivo `.env.local` dentro da pasta `web/` com as seguintes definições:

```env
# URL de Comunicação com a API Backend Go
BACKEND_BASE_URL=http://localhost:8080

# Segredos de Autenticação do NextAuth
AUTH_SECRET=P79ER6Uij/I3QAiH+HuVXGTKnay4n0lGI8jB6oeHrZI=
NEXTAUTH_URL=http://localhost:3001
AUTH_TRUST_HOST=true

# Outros tempos limites
BACKEND_TIMEOUT_MS=1200
```

---

## 🛠️ Executando a Aplicação Localmente

Navegue até a pasta `web/` ou utilize os comandos globais do Makefile na raiz.

### 1. Instalar Dependências

```bash
bun install
```

### 2. Executar Servidor de Desenvolvimento

Inicia o Next.js na porta padrão customizada `3001`:

```bash
bun dev
```

Abra [http://localhost:3001](http://localhost:3001) para interagir com o app.

---

## 🧪 Testes e Qualidade de Código

### Testes Unitários

Para rodar os testes unitários isolados utilizando o executor do `bun`:

```bash
bun test
# ou use o script utilitário
./scripts/run-tests.sh
```

### Linter & Formatador (Biome)

Para auditar a formatação e qualidade do código:

```bash
bun run lint
```

Para aplicar correções automáticas:

```bash
bun run lint:fix
bun run format:fix
```
