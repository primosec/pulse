# Pulse — Sistema de Monitoramento e Alertas em Tempo Real

> **Status:** Em desenvolvimento | **Objetivo:** Portfolio técnico com foco em backend Go, concorrência e sistemas real-time.

---

## O que é o Pulse

Pulse é um sistema de monitoramento de serviços e endpoints que verifica periodicamente a disponibilidade de URLs e APIs, calcula métricas de uptime e notifica o usuário em tempo real via WebSocket quando um serviço cai.

Pense num Uptime Kuma ou Pingdom simplificado, construído do zero com foco em demonstrar domínio de Go, concorrência, sistemas distribuídos e integração full-stack.

---

## Por que este projeto importa para o portfolio

- Não é um CRUD genérico — tem domínio de negócio real e reconhecível
- Demonstra concorrência em Go de forma prática (worker pool, goroutines, channels)
- WebSocket em produção é raro em portfolios de backend
- Cobre preocupações de produção: rate limiting, graceful shutdown, logging estruturado
- Full-stack coerente: Go backend + Vue 3 frontend modernos
- Infraestrutura completa: Docker, CI, migrations, documentação de API

---

## Tech Stack

### Backend

| Tecnologia | Papel | Por quê |
|---|---|---|
| **Go 1.23+** | Linguagem principal | Performance, concorrência nativa, ecossistema crescente |
| **Gin** | HTTP framework | Mais usado em Go, familiar para recrutadores |
| **PostgreSQL** | Banco de dados principal | Relacional, ideal para histórico de checks e métricas |
| **Redis** | Cache + filas + rate limiting | Volatilidade desejada, baixíssima latência |
| **sqlc** | Geração de código SQL | Queries tipadas, sem ORM overhead, mais idiomático em Go do que GORM |
| **Asynq** | Background jobs | Biblioteca madura sobre Redis, elimina worker pool manual desnecessário |
| **coder/websocket** | WebSocket | Fork ativo do gorilla/websocket (arquivado em 2023) |
| **JWT (golang-jwt)** | Autenticação | Stateless, padrão de mercado |
| **Viper** | Configuração | Lê env vars, arquivos .env, .yaml — flexível para Docker |
| **Zerolog** | Logging estruturado | Zero alocação, JSON output, melhor performance que zap em muitos casos |

### Frontend

| Tecnologia | Papel |
|---|---|
| **Vue 3** (Composition API + `<script setup>`) | Framework principal |
| **Vite** | Build tool |
| **Pinia** | Gerenciamento de estado |
| **Vue Router** | Navegação SPA |
| **Tailwind CSS** | Estilização utilitária |
| **ApexCharts** | Gráficos de uptime e histórico |

### DevOps & Ferramentas

| Tecnologia | Papel |
|---|---|
| **Docker + Docker Compose** | Containerização e orquestração local |
| **GitHub Actions** | CI (build, lint, testes) |
| **Air** | Hot reload em desenvolvimento |
| **Make** | Scripts de automação (run, migrate, test, build) |
| **Swagger (swaggo)** | Documentação da API gerada automaticamente |

---

## Arquitetura

### Visão Geral

```
┌─────────────┐     HTTP/WS      ┌─────────────────────────────────┐
│  Vue 3 SPA  │ ◄──────────────► │         Gin HTTP Server          │
└─────────────┘                  └────────────┬────────────────────┘
                                              │
                          ┌───────────────────┼───────────────────┐
                          ▼                   ▼                   ▼
                    ┌──────────┐       ┌──────────┐       ┌──────────┐
                    │ Service  │       │ WebSocket│       │  Worker  │
                    │  Layer   │       │ Manager  │       │  Pool    │
                    └────┬─────┘       └────┬─────┘       └────┬─────┘
                         │                  │                   │
                    ┌────▼─────┐            │            ┌──────▼──────┐
                    │Repository│            │            │    Asynq    │
                    │  Layer   │            │            │  (Redis Q)  │
                    └────┬─────┘            │            └──────┬──────┘
                         │                  │                   │
                    ┌────▼──────────────────▼───────────────────▼─────┐
                    │              PostgreSQL + Redis                   │
                    └──────────────────────────────────────────────────┘
```

### Estrutura de Pastas

```
pulse-go/
├── cmd/
│   └── pulse/
│       └── main.go              # Entry point — inicializa deps, roda o servidor
│
├── internal/
│   ├── api/
│   │   ├── handler/             # HTTP handlers (um arquivo por recurso)
│   │   │   ├── auth.go
│   │   │   ├── monitor.go
│   │   │   └── check.go
│   │   ├── middleware/          # Auth JWT, rate limit, CORS, request ID
│   │   │   ├── auth.go
│   │   │   └── ratelimit.go
│   │   └── router.go            # Registro de todas as rotas Gin
│   │
│   ├── domain/                  # Entidades de negócio — structs puras, sem deps externas
│   │   ├── monitor.go
│   │   ├── check.go
│   │   └── user.go
│   │
│   ├── dto/                     # Request/Response structs — separa API da domain
│   │   ├── monitor_dto.go
│   │   ├── auth_dto.go
│   │   └── check_dto.go
│   │
│   ├── repository/              # Acesso a dados — implementa interfaces definidas em service
│   │   ├── monitor_repo.go
│   │   ├── check_repo.go
│   │   ├── user_repo.go
│   │   └── redis_repo.go
│   │
│   ├── service/                 # Regras de negócio — orquestra domain + repository
│   │   ├── monitor_service.go
│   │   ├── auth_service.go
│   │   └── check_service.go
│   │
│   ├── worker/                  # Background jobs via Asynq
│   │   ├── dispatcher.go
│   │   ├── health_check.go
│   │   └── scheduler.go
│   │
│   ├── websocket/               # Gerenciamento de conexões WebSocket
│   │   ├── hub.go
│   │   ├── client.go
│   │   └── message.go
│   │
│   └── config/
│       └── config.go
│
├── pkg/
│   ├── apperror/                # Tipos de erro customizados com código HTTP
│   │   └── errors.go
│   ├── httputil/                # Helper para respostas JSON padronizadas
│   │   └── response.go
│   └── validator/
│       └── validator.go
│
├── db/
│   ├── migrations/              # Arquivos SQL numerados
│   │   ├── 001_create_users.sql
│   │   ├── 002_create_monitors.sql
│   │   └── 003_create_check_results.sql
│   └── queries/                 # Queries SQL para o sqlc
│       ├── monitor.sql
│       ├── check.sql
│       └── user.sql
│
├── frontend/                    # Vue 3 SPA
│   ├── src/
│   │   ├── components/
│   │   ├── views/
│   │   ├── stores/
│   │   ├── router/
│   │   └── api/
│   └── vite.config.ts
│
├── .github/
│   └── workflows/
│       └── ci.yml
│
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── .env.example
├── sqlc.yaml
└── PULSE.md
```

### Padrões e Decisões de Design

**Repository Pattern**
Cada recurso (Monitor, User, CheckResult) tem uma interface de repositório definida na camada de service. A implementação concreta fica em `repository/`. Isso permite trocar PostgreSQL por outra coisa sem alterar regras de negócio.

```go
// service define a interface
type MonitorRepository interface {
    Create(ctx context.Context, m *domain.Monitor) error
    FindByID(ctx context.Context, id uuid.UUID) (*domain.Monitor, error)
    FindActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Monitor, error)
    UpdateStatus(ctx context.Context, id uuid.UUID, status domain.MonitorStatus) error
    Delete(ctx context.Context, id uuid.UUID) error
}

// repository implementa
type PostgresMonitorRepository struct {
    queries *db.Queries // gerado pelo sqlc
}
```

**Separação Domain / DTO**
Domain entities (`internal/domain/`) representam o estado interno do sistema. DTOs (`internal/dto/`) representam o contrato externo da API. Nunca retorne uma entity de domínio diretamente num handler — faça a conversão no handler ou no service.

**Worker Pool com Asynq**
O scheduler periodicamente escaneia monitores ativos e enfileira uma task `health_check` por monitor no Redis. O Asynq server consome essas tasks com N workers concorrentes. Se o servidor reiniciar, as tasks pendentes ainda estão na fila.

```
Scheduler (ticker a cada 10s)
  └── Para cada monitor ativo:
        └── Enfileira task "health_check:{monitorID}" no Redis (se não existe já)

Asynq Worker (N goroutines)
  └── Consome "health_check:{monitorID}"
        ├── Executa HTTP GET com timeout
        ├── Salva CheckResult no PostgreSQL
        ├── Atualiza status do monitor
        └── Publica evento no WebSocket Hub se status mudou
```

**WebSocket Hub**
Um único Hub goroutine centraliza todas as conexões WebSocket. Clients se registram/desregistram via channels. O Hub recebe eventos dos workers e faz broadcast para os clientes corretos (por userID).

```go
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
        case client := <-h.unregister:
        case message := <-h.broadcast:
        }
    }
}
```

---

## Modelo de Dados

### Tabela: `users`
```sql
id          UUID PRIMARY KEY DEFAULT gen_random_uuid()
email       TEXT UNIQUE NOT NULL
password    TEXT NOT NULL          -- bcrypt hash
name        TEXT NOT NULL
created_at  TIMESTAMPTZ DEFAULT NOW()
```

### Tabela: `monitors`
```sql
id           UUID PRIMARY KEY DEFAULT gen_random_uuid()
user_id      UUID REFERENCES users(id) ON DELETE CASCADE
name         TEXT NOT NULL
url          TEXT NOT NULL
type         TEXT NOT NULL DEFAULT 'http'    -- http | tcp | ping
interval     INT NOT NULL DEFAULT 60         -- segundos
timeout      INT NOT NULL DEFAULT 10         -- segundos
status       TEXT NOT NULL DEFAULT 'pending' -- pending | up | down
uptime_pct   FLOAT DEFAULT 0
last_checked TIMESTAMPTZ
is_active    BOOLEAN DEFAULT TRUE
created_at   TIMESTAMPTZ DEFAULT NOW()
```

### Tabela: `check_results`
```sql
id          UUID PRIMARY KEY DEFAULT gen_random_uuid()
monitor_id  UUID REFERENCES monitors(id) ON DELETE CASCADE
status      TEXT NOT NULL    -- up | down | timeout
status_code INT              -- HTTP status code (nullable para TCP/ping)
latency_ms  INT              -- tempo de resposta em ms
error_msg   TEXT             -- mensagem de erro se houver
checked_at  TIMESTAMPTZ DEFAULT NOW()
```

---

## API — Endpoints

### Autenticação
```
POST /api/v1/auth/register    Body: { email, password, name }
POST /api/v1/auth/login       Body: { email, password } → { token }
```

### Monitores (requer JWT)
```
GET    /api/v1/monitors              Lista monitores do usuário
POST   /api/v1/monitors              Cria monitor
GET    /api/v1/monitors/:id          Detalhe do monitor
PUT    /api/v1/monitors/:id          Atualiza monitor
DELETE /api/v1/monitors/:id          Remove monitor
POST   /api/v1/monitors/:id/pause    Pausa/retoma monitor
```

### Checks
```
GET /api/v1/monitors/:id/checks      Histórico de checks (paginado)
```

### WebSocket
```
WS /ws?token=<jwt>                   Conexão autenticada para eventos em tempo real
```

### Healthcheck da própria aplicação
```
GET /health                          { status: "ok", version: "...", uptime: "..." }
```

---

## Eventos WebSocket

```json
// Monitor mudou de status
{
  "type": "status_change",
  "data": {
    "monitor_id": "uuid",
    "monitor_name": "Meu Site",
    "old_status": "up",
    "new_status": "down",
    "checked_at": "2026-05-02T10:00:00Z"
  }
}

// Resultado de check em tempo real
{
  "type": "check_result",
  "data": {
    "monitor_id": "uuid",
    "status": "up",
    "latency_ms": 142,
    "checked_at": "2026-05-02T10:00:00Z"
  }
}
```

## O que este projeto demonstra

| Habilidade | Como demonstrada |
|---|---|
| Concorrência em Go | Worker pool via Asynq, WebSocket Hub goroutine, channels |
| Design de API REST | Versionamento, DTOs, error handling padronizado |
| Sistemas real-time | WebSocket bidirecional com autenticação |
| Background processing | Fila de jobs com retry, persistência, deduplicação |
| Banco de dados | Migrations, queries otimizadas, índices corretos |
| Cache | Redis para rate limiting e status recente |
| Segurança | JWT stateless, bcrypt, rate limiting, input validation |
| Qualidade de código | Testes, lint, logging estruturado, graceful shutdown |
| DevOps | Docker multistage, Compose, CI com GitHub Actions |
| Full-stack | Backend Go + Frontend Vue 3 completamente integrados |

---

*Pulse — Desenvolvido por [PrimoSec] | Maio 2026*
