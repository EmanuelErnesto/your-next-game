# Your Next Game 🎮

Bem-vindo ao **Your Next Game**, um ecossistema projetado para ajudar os usuários a gerenciar seu backlog de jogos e descobrir novas recomendações. Este repositório é um monorepo que encapsula tanto a interface do usuário quanto o servidor de serviços.

## 📁 Estrutura do Repositório

O projeto é organizado de maneira modular e limpa:

```
your-next-game/
├── backend/                  # Servidor de API de alta performance desenvolvido em Go
├── web/                      # Frontend Next.js (App Router, Tailwind CSS, DaisyUI, Zustand)
├── docs/                     # Documentação do projeto, guias de migração e arquitetura
├── Makefile                  # Orquestração central de tarefas e automação local
├── docker-compose.yaml       # Serviços de infraestrutura (PostgreSQL, OTel, Grafana stack)
└── otel-collector-config.yaml # Configurações de coleta de telemetria
```

---

## 🛠️ Pré-requisitos

Para executar e desenvolver no projeto localmente, você precisará de:

- [Bun](https://bun.sh/) (para gerenciar dependências e executar o frontend Next.js)
- [Go](https://go.dev/) 1.21+ (para rodar o backend e ferramentas de migração)
- [Docker](https://www.docker.com/) & Docker Compose (para subir banco de dados e serviços de monitoramento)

---

## ⚡ Início Rápido

O arquivo `Makefile` na raiz atua como o painel de controle principal. Siga os passos abaixo para preparar o ambiente:

### 1. Configurar e Instalar Dependências

Instale todas as dependências do frontend e do backend, inicie os containers locais de banco de dados e execute as migrações automáticas:

```bash
make setup
```

### 2. Executar em Modo de Desenvolvimento

Inicie tanto o frontend Next.js quanto o backend Go concorrentemente em um único terminal:

```bash
make dev
```

* O Frontend estará disponível em: [http://localhost:3001](http://localhost:3001)
* O Backend de desenvolvimento rodará em: [http://localhost:8080](http://localhost:8080)

---

## 📌 Principais Comandos de Automação

| Comando | Descrição |
| :--- | :--- |
| `make setup` | Instala dependências, sobe containers e aplica migrations do banco |
| `make dev` | Sobe o banco e roda concorrentemente o frontend e backend |
| `make frontend-dev` | Roda apenas o frontend Next.js em modo de desenvolvimento |
| `make backend-dev` | Roda apenas o backend Go com reinicialização em tempo de execução |
| `make migrate` | Aplica as migrations do banco de dados PostgreSQL |
| `make up-containers` | Inicia os serviços Docker (Postgres, OTel, Loki, Grafana) |
| `make down-containers` | Encerra e remove todos os containers locais ativos |

---

## 📘 Documentação Detalhada

Cada subdiretório principal possui seu próprio arquivo de orientações detalhadas:

* **Frontend**: Consulte o [web/README.md](file:///home/emanuel/projects/personal/your-next-game/web/README.md) para detalhes de estilização, padrões de design (Feature-Sliced Design), roteamento e testes unitários.
* **Backend**: Consulte o [backend/README.md](file:///home/emanuel/projects/personal/your-next-game/backend/README.md) para endpoints disponíveis, arquitetura em Go, migrations de banco e testes de integração.
* **Migrações e Testes**: Guias adicionais de design de infraestrutura estão em [docs/](file:///home/emanuel/projects/personal/your-next-game/docs).
