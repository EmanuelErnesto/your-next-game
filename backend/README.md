# Your Next Game - Backend ⚙️

Este é o servidor backend de alta performance desenvolvido em Go (Golang) para gerenciar dados de catálogo, integrações com a Steam, perfis de usuários e bancos de dados relacionais.

## 🚀 Tecnologias Utilizadas

- **Go (Golang)**: Linguagem compilada de altíssima performance e concorrência nativa.
- **Chi Router**: Roteador HTTP minimalista e extremamente performático.
- **PostgreSQL**: Banco de dados relacional para persistência segura e transacional.
- **Golang Migrations**: Gerenciamento evolutivo de schemas de banco de dados (`cmd/migrate`).
- **OpenTelemetry**: Coleta integrada de métricas e rastreamentos (traces) distribuídos.

---

## ⚙️ Variáveis de Ambiente

O backend é altamente configurável através de variáveis de ambiente. Configure um arquivo `.env` dentro da pasta `backend/` para execução local:

```env
# Chave Secreta de Integração da Steam
STEAM_API_KEY="A067221CC77DFDEABFE6CA9DA495804C"

# URL de Sucesso/Redirecionamento do Frontend
FRONTEND_SUCCESS_URL=http://localhost:3001/

# Porta HTTP do Servidor
HTTP_PORT=8080

# Conexão PostgreSQL (Substitua se necessário)
POSTGRES_DSN="postgres://postgres:postgres@localhost:5432/your_next_game?sslmode=disable"

# Repositório de biblioteca (usar 'postgres' ou 'memory' para testes rápidos)
LIBRARY_REPOSITORY=postgres
```

---

## 🛠️ Executando a API Localmente

### 1. Aplicar Migrações do Banco de Dados

Garanta que o container Postgres esteja ativo e execute as migrações automáticas:

```bash
POSTGRES_DSN="postgres://postgres:postgres@localhost:5432/your_next_game?sslmode=disable" go run ./cmd/migrate
```

### 2. Rodar o Servidor

Inicie a API HTTP localmente:

```bash
go run ./cmd/api/main.go
```

O backend estará ativo em: [http://localhost:8080](http://localhost:8080)

---

## 📡 Endpoints Disponíveis da API

A API expõe as seguintes rotas e recursos integrados:

### Verificação de Saúde e Diagnóstico
- `GET /healthz` - Healthcheck simples
- `GET /readyz` - Prontidão de dependências (Postgres)
- `GET /v1/version` - Retorna a versão de build ativa

### Integração Steam & Autenticação
- `GET /v1/auth/steam/login-url` - URL de login OpenID da Steam
- `POST /v1/auth/steam/exchange` - Troca credenciais e realiza login
- `GET /v1/integrations/steam/users/{steamId}/games` - Busca lista de jogos direto da API Steam

### Catálogo de Jogos do Usuário
- `GET /v1/users/{userId}/games` - Lista todos os jogos associados ao usuário
- `GET /v1/users/{userId}/games/{gameId}` - Detalhes de um jogo do usuário
- `GET /v1/catalog/steam/apps/{appId}` - Busca informações de catálogo da Steam

---

## 🧪 Testes de Integração

O backend conta com testes integrados que validam consultas e conexões reais com o PostgreSQL. Para executá-los:

```bash
POSTGRES_DSN="postgres://postgres:postgres@localhost:5432/your_next_game?sslmode=disable" go test ./tests/integration -v
```

## 📖 Contrato OpenAPI

O contrato descritivo das APIs está documentado e padronizado usando o formato OpenAPI em:

- `api/openapi/openapi.yaml`
