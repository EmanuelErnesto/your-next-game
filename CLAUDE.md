# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

- `bun run dev` — Start dev server
- `bun run build` — Production build
- `bun run lint` — Biome check (formatter + linter + assist)
- `bun run lint:fix` — Biome safe fixes
- `bun run format:fix` — Biome format only
- `bun run format` — Biome all fixes including unsafe

No test framework is configured.

## Architecture

**Feature-Sliced Design (FSD) with Next.js App Router.** Unidirectional dependency flow:

```
app/ → widgets/ → features/ → entities/ → shared/
```

- **`app/`** — Next.js routes, layouts, API routes, providers. Server Components by default. Data fetching happens here; pass data down via props.
- **`widgets/`** — Composite UI blocks (GameGrid, Sidebar, HomeHeader, etc.)
- **`features/`** — Business logic + client-side UI. Zustand stores and interactive components live here. These use `'use client'`.
- **`entities/`** — Type-only layer. Contains entity `type` aliases. No executable logic, no UI.
- **`shared/`** — Cross-cutting: API clients (`shared/api/`), config (`shared/config/`), design system (`shared/ui/` with Atomic Design: atoms/molecules/organisms), and utilities (`shared/lib/`).

Path aliases: `@app/*`, `@widgets/*`, `@features/*`, `@entities/*`, `@shared/*`.

Every slice must have a public `index.ts` barrel export. Import from slice roots only (e.g., `@features/auth`), never from internal paths.

## Stack

- **Runtime:** Bun
- **Framework:** Next.js 16 (App Router) + React 19 + RSC
- **Database:** External Go backend (your-next-game-backend)
- **Auth:** NextAuth v5 beta with custom Steam OpenID provider
- **State:** Zustand (client), RSC + fetch (server)
- **UI:** Tailwind CSS v4 + DaisyUI v5 + Lucide icons
- **Validation:** Zod v4
- **Logging:** Pino structured logger at `@shared/lib/logger` — `console.log` is forbidden

## Critical Rules

- **Next.js 16 has breaking changes** — read the relevant guide in `node_modules/next/dist/docs/` before writing any code. Heed deprecation notices.
- **Use `type`, never `interface` or `class`** for DTOs and entities.
- **Server Components by default** — only add `'use client'` where interactivity is needed.
- **No `console.log`** — use `@shared/lib/logger`.
- **No `process.env` scattered in code** — import from `@shared/config/env` exclusively.
- **Zod validation required** on all Server Action inputs.
- **No code comments** unless truly necessary.
- **All UI text is in Brazilian Portuguese (pt-BR).**

## Formatting

Biome: single quotes, trailing commas (all), semicolons, 2-space indent, 100 char line width.
