# Up - Espaço — Documentação Técnica Completa

> Documentação de arquitetura, código, infraestrutura e operação do sistema **Up - Espaço**, uma plataforma de gestão e comunicação escolar (estilo "diário digital" / agenda escolar) para creches/escolas de educação infantil.

---

## Sumário

1. [Visão Geral do Projeto](#1-visão-geral-do-projeto)
2. [Arquitetura da Solução](#2-arquitetura-da-solução)
3. [Tecnologias Utilizadas](#3-tecnologias-utilizadas)
4. [Estrutura de Diretórios](#4-estrutura-de-diretórios)
5. [Análise Detalhada do Backend](#5-análise-detalhada-do-backend)
6. [Análise Detalhada do Frontend](#6-análise-detalhada-do-frontend)
7. [Fluxo de Execução da Aplicação](#7-fluxo-de-execução-da-aplicação)
8. [Banco de Dados](#8-banco-de-dados)
9. [APIs (Referência de Endpoints)](#9-apis-referência-de-endpoints)
10. [Segurança](#10-segurança)
11. [Configuração e Infraestrutura](#11-configuração-e-infraestrutura)
12. [Observabilidade (Prometheus + Grafana)](#12-observabilidade-prometheus--grafana)
13. [Dependências do Projeto](#13-dependências-do-projeto)
14. [Regras de Negócio](#14-regras-de-negócio)
15. [Resumo Executivo](#15-resumo-executivo)

---

## 1. Visão Geral do Projeto

### 1.1 Objetivo principal

O **Up - Espaço** é um sistema web (API REST em Go + SPA em React) que conecta a equipe pedagógica de uma escola/creche aos responsáveis (pais/mães/cuidadores) dos alunos. Ele centraliza:

- Registro de presença/falta diária dos alunos.
- Feed de atividades pedagógicas com fotos (estilo rede social).
- Linha do tempo (timeline) de eventos do dia de cada criança.
- Marcos de desenvolvimento (jornada) da criança.
- Agenda de eventos da escola com confirmação de presença (RSVP).
- Comunicados/avisos institucionais com controle de leitura.
- Cadastro de alunos, turmas, responsáveis autorizados e dados de saúde.

### 1.2 Problema que resolve

Substitui a comunicação fragmentada entre escola e família (grupos de WhatsApp, agendas de papel, recados soltos) por um canal único, estruturado e auditável, com controle de quem pode ver o quê (visibilidade privada vs. turma) e com métricas operacionais (presença, engajamento no feed).

### 1.3 Usuários-alvo

| Papel (`role`) | Quem é | O que pode fazer |
|---|---|---|
| `profissional` | Professores, equipe pedagógica, direção | Cadastrar/editar alunos, turmas, posts, eventos, comunicados, marcar presença, ver todos os alunos |
| `responsavel` | Pais, mães, responsáveis legais | Ver apenas os filhos vinculados à própria conta, curtir/comentar/salvar posts, confirmar presença em eventos, ler comunicados |

### 1.4 Fluxo geral de funcionamento

1. O usuário cria conta (`/api/auth/register`) escolhendo o papel (`profissional` ou `responsavel`) e faz login (`/api/auth/login`), recebendo um **JWT**.
2. Se for `profissional`, vê a lista de todos os alunos e pode cadastrar novos, vinculando-os a uma turma, um professor e um responsável.
3. Se for `responsavel`, vê apenas os alunos vinculados ao seu `user_id` (campo `guardian_user_id` na tabela `students`).
4. Ao selecionar uma criança, o app carrega: dados do aluno, timeline do dia e próximo evento da agenda.
5. A partir daí, a navegação acontece entre as abas: **Início**, **Atividades** (feed), **Jornada** (marcos), **Agenda** (eventos) e **Perfil**.
6. Toda a comunicação frontend ↔ backend é via HTTP/JSON, autenticada por `Authorization: Bearer <token>`.

### 1.5 Principais módulos e responsabilidades

| Módulo | Responsabilidade |
|---|---|
| `backend-up` | API REST em Go, regras de negócio, persistência em Postgres, métricas Prometheus |
| `frontend-up` | SPA em React (Vite + Tailwind), consome a API, UI mobile-first |
| `grafana/` + `prometheus.yml` | Stack de observabilidade (dashboards e coleta de métricas) |
| `docker-compose.yml` | Orquestração de todo o ambiente (db, backend, frontend, monitoramento) |

---

## 2. Arquitetura da Solução

### 2.1 Estilo arquitetural

O **backend** segue uma variação de **Clean Architecture / Arquitetura Hexagonal (Ports & Adapters)**, organizada em 3 camadas concêntricas:

```
domain        →  regras e contratos puros (entidades, interfaces de repositório, erros)
application   →  casos de uso (orquestram regras de negócio) + DTOs + ports
adapters      →  implementações concretas (HTTP handlers, Postgres, JWT, bcrypt)
```

Não há frameworks "mágicos" de DI: a montagem das dependências é feita manualmente em `cmd/api/main.go` (Composition Root), o que é uma escolha consciente para manter o projeto simples e explícito.

O **frontend** é um **monolito de componentes React** sem roteador dedicado (a "rota" é controlada por estado local, `active` em `App.jsx`) — arquitetura adequada para o tamanho atual do app, mas listada como ponto de atenção na seção 15.

### 2.2 Separação de responsabilidades (backend)

```
domain/entities         → "o que é" um Student, User, Post, Event... (structs puras)
domain/repositories      → contratos ("interfaces") de persistência, sem implementação
domain/errors             → erros de negócio reaproveitáveis

application/usecases     → 1 caso de uso = 1 ação de negócio (ex: CreatePostUseCase)
application/dto          → formato de entrada/saída das requisições HTTP
application/ports        → interfaces intermediárias entre usecases e repositórios concretos

adapters/http/handlers   → traduzem HTTP ↔ usecases (parse de request, status code, JSON)
adapters/http/middleware → autenticação JWT, CORS, métricas Prometheus
adapters/http/routes     → registro de rotas (gorilla/mux) + Swagger
adapters/database/postgres → implementação real das interfaces de domain/repositories
adapters/security        → JWT e hashing de senha (bcrypt)
adapters/repositories     → container que agrega todos os repositórios postgres
```

A regra de dependência é sempre **de fora para dentro**: `adapters` depende de `application`, que depende de `domain`. O `domain` não depende de nada externo — não importa nem o driver do Postgres.

### 2.3 Fluxo de uma requisição HTTP (exemplo: criar um post)

```
Cliente (React)
   │  POST /api/posts?student_id=5  { title, description, ... }  + Authorization: Bearer <token>
   ▼
middleware.WithMetrics            → mede latência/contagem (Prometheus)
   ▼
middleware.WithCORS               → valida origem
   ▼
mux.Router                        → casa a rota com PostsHandler.Create
   ▼
middleware.RequireAuth            → valida o JWT, injeta userID/role no contexto
   ▼
PostsHandler.Create               → decodifica JSON em dto.CreatePostRequest, monta entities.Post
   ▼
CreatePostUseCase.Execute         → valida campos obrigatórios e tamanho da imagem (validateImage)
   ▼
ports.PostPort (PostRepository)   → INSERT INTO posts ... RETURNING id, created_at, updated_at
   ▼
PostsHandler                      → serializa o Post criado como JSON, HTTP 201
   ▼
Cliente (React)                   → atualiza o estado local (setPosts) e re-renderiza
```

### 2.4 Fluxo de autenticação

```
1. POST /api/auth/login {email, password}
2. LoginUseCase busca o usuário por email (UserRepository.GetByEmail)
3. security.CheckPassword compara bcrypt(hash salvo) com a senha enviada
4. Se válido: security.GenerateToken cria um JWT HS256 com Subject=userID, Role=role, TTL=7 dias
5. Backend devolve { token, user }
6. Frontend salva o token em localStorage (api.js → setToken)
7. Toda requisição subsequente envia Authorization: Bearer <token>
8. middleware.RequireAuth (ou OptionalAuth) decodifica o token e injeta userID/role no context.Context
9. Handlers/usecases leem esses valores via middleware.UserIDFromContext / RoleFromContext
```

### 2.5 Comunicação entre módulos (frontend)

```
Componentes React (Login, Feed, Agenda, ...)
        │  chamam funções exportadas
        ▼
frontend-up/src/api.js   (camada única de acesso à API)
        │  fetch() com Authorization header
        ▼
Backend REST API (Go)
```

Não existe um state-manager global (Redux/Zustand/Context API de domínio): o estado vive em `useState`/`useCallback` dentro de `App.jsx`, que repassa dados e callbacks via props para as páginas filhas.

---

## 3. Tecnologias Utilizadas

### Backend

| Tecnologia | O que é | Função no projeto | Por que foi escolhida |
|---|---|---|---|
| **Go 1.22** | Linguagem compilada, estaticamente tipada | Implementa toda a API REST | Performance, concorrência nativa, binário único fácil de deployar, tipagem forte reduz bugs |
| **gorilla/mux** | Router HTTP para Go | Define rotas, parâmetros de URL (`{id}`), métodos HTTP | Mais flexível que o `net/http` puro (suporta `{id:[0-9]+}`, subrotas), leve e maduro |
| **lib/pq** | Driver Postgres puro-Go | Conexão com o banco via `database/sql` | Driver estável e amplamente usado, compatível com `database/sql` padrão |
| **golang-jwt/jwt v5** | Biblioteca de JSON Web Tokens | Gera e valida os tokens de autenticação | Implementação de referência em Go para JWT, suporta claims customizadas |
| **golang.org/x/crypto (bcrypt)** | Hashing de senhas | `HashPassword` / `CheckPassword` | Algoritmo padrão de mercado para senha, resistente a brute-force (custo ajustável) |
| **prometheus/client_golang** | Cliente Prometheus para Go | Exposição de métricas em `/metrics` | Padrão de mercado para métricas em apps cloud-native |
| **database/sql + embed** | Pacotes nativos do Go | Migrations embutidas no binário (`//go:embed *.up.sql`) | Elimina a necessidade de copiar arquivos `.sql` separadamente no deploy |

### Frontend

| Tecnologia | O que é | Função no projeto | Por que foi escolhida |
|---|---|---|---|
| **React 18** | Biblioteca de UI baseada em componentes | Toda a interface do usuário | Ecossistema maduro, componentização, hooks |
| **Vite 6** | Build tool / dev server | Compilação, HMR, bundling de produção | Extremamente rápido em dev, configuração mínima |
| **Tailwind CSS 4** (`@tailwindcss/vite`) | Framework CSS utility-first | Estilização de toda a UI | Produtividade alta, consistência visual, sem CSS custom espalhado |
| **lucide-react** | Biblioteca de ícones SVG | Ícones em toda a interface | Conjunto consistente, leve, tree-shakeable |
| **JSX puro (sem TypeScript)** | — | Componentes `.jsx` | Escolha de simplicidade/velocidade de desenvolvimento — ver recomendações na seção 15 |

### Infraestrutura / Observabilidade

| Tecnologia | Função |
|---|---|
| **Docker / Docker Compose** | Orquestra `db`, `backend`, `frontend`, `prometheus`, `grafana`, `postgres-exporter`, `cadvisor` em containers isolados |
| **PostgreSQL 15 (alpine)** | Banco de dados relacional principal |
| **Prometheus** | Coleta métricas do backend (`/metrics`), do Postgres (via exporter) e dos containers (via cAdvisor) |
| **Grafana** | Dashboards de observabilidade (requisições/s, latência p95, CPU/RAM do backend, saúde do Postgres) |
| **postgres_exporter** | Exporta métricas internas do Postgres pro Prometheus |
| **cAdvisor** | Exporta métricas de uso de recursos (CPU/RAM) por container |

---

## 4. Estrutura de Diretórios

```text
up-espaco/
├── docker-compose.yml         # orquestração de todos os serviços
├── prometheus.yml              # config de scrape do Prometheus
├── README.md                   # este documento
│
├── backend-up/                 # API REST em Go
│   ├── Dockerfile
│   ├── go.mod / go.sum
│   ├── docs/                   # swagger.json (OpenAPI) servido em /swagger
│   ├── cmd/
│   │   └── api/
│   │       └── main.go         # composition root: monta tudo e inicia o servidor
│   └── internal/
│       ├── domain/
│       │   ├── entities/       # structs de domínio (Student, User, Post, Event...)
│       │   ├── repositories/   # interfaces de persistência (contratos)
│       │   └── errors/         # erros de domínio reaproveitáveis
│       ├── application/
│       │   ├── usecases/       # 1 arquivo/grupo = 1 ação de negócio
│       │   ├── dto/            # requests/responses HTTP
│       │   └── ports/          # interfaces usadas pelos usecases (StudentPort, PostPort...)
│       ├── adapters/
│       │   ├── http/
│       │   │   ├── handlers/   # um handler por agregado (auth, student, post, event...)
│       │   │   ├── middleware/ # auth (JWT), cors, metrics
│       │   │   └── routes/     # registro de rotas + swagger
│       │   ├── database/
│       │   │   ├── postgres/   # implementação dos repositórios em SQL puro
│       │   │   └── migrations/ # arquivos *.up.sql / *.down.sql + runner embutido
│       │   ├── repositories/   # RepositoryContainer (agrega todos os repos)
│       │   └── security/       # jwt.go, password.go
│       └── config/
│           └── config.go       # leitura de variáveis de ambiente
│
├── frontend-up/                 # SPA em React
│   ├── Dockerfile
│   ├── package.json
│   ├── vite.config.js
│   ├── index.html
│   ├── assets/ , public/        # logo, avatar padrão
│   └── src/
│       ├── main.jsx             # entry point (ReactDOM.createRoot)
│       ├── App.jsx              # componente raiz: sessão, layout, navegação
│       ├── api.js               # única camada de acesso à API (fetch wrapper)
│       ├── Login.jsx / CriarConta.jsx / EsqueceuSenha.jsx   # autenticação
│       ├── Feed.jsx             # feed de atividades (posts, curtidas, comentários)
│       ├── Jornada.jsx          # marcos de desenvolvimento
│       ├── Agenda.jsx           # eventos da escola + RSVP
│       ├── Comunicados.jsx      # avisos institucionais
│       ├── Perfil.jsx           # perfil do usuário logado
│       ├── Crianca.jsx          # perfil completo do aluno (saúde, escola, responsáveis)
│       ├── Turmas.jsx           # gestão de turmas + chamada de presença
│       ├── DashboardAdmin.jsx   # painel administrativo (visão profissional)
│       └── SearchableSelect.jsx # combobox de busca reutilizável
│
└── grafana/
    └── provisioning/
        ├── alerting/contact-points.yml
        ├── dashboards/up-espaco-observability.json + up-espaco.yml
        ├── datasources/prometheus.yml
        ├── notifiers/ , plugins/
```

### Por que essa organização existe

- **`domain` isolado**: garante que as regras de negócio não fiquem acopladas a Postgres, HTTP ou qualquer biblioteca externa — facilita testes unitários e troca futura de banco/framework.
- **`usecases` 1:1 com ações de negócio**: cada operação (criar post, marcar presença, fazer RSVP) é uma unidade testável isoladamente, sem depender do HTTP.
- **`ports` vs `repositories`**: existem dois níveis de interface porque `ports.PostPort`, por exemplo, expõe nomes pensados para o usecase (`LikePost`, `CreatePost`) enquanto `repositories.PostRepository` expõe nomes mais "CRUD" (`Create`, `IncrementLikes`). A implementação concreta (`postgres.PostRepository`) satisfaz as duas interfaces.
- **Frontend "flat"**: como o app tem poucas telas, todos os componentes de página ficam direto em `src/`, sem subpastas `pages/`/`components/` — organização adequada ao tamanho atual, mas que merece revisão se o projeto crescer (ver seção 15).

---

## 5. Análise Detalhada do Backend

### 5.1 `cmd/api/main.go` — Composition Root

**Responsabilidade**: ponto de entrada do binário. Não contém lógica de negócio — apenas monta e conecta as peças.

Fluxo interno:
1. `config.Load()` lê variáveis de ambiente (porta, DSN do Postgres, segredo JWT, CORS).
2. Abre a conexão `sql.Open("postgres", ...)` e faz `PingContext` (falha rápido se o banco não estiver acessível).
3. `migrations.Run(ctx, db)` aplica todas as migrations pendentes.
4. `postgres.NewDB(db)` envolve a conexão; `repositories.NewRepositoryContainer(store)` instancia todos os repositórios concretos.
5. `routes.NewRouter(repos, jwtSecret, tokenTTL)` monta todos os handlers/usecases e registra as rotas.
6. O handler final é envolvido por `middleware.WithMetrics` e `middleware.WithCORS`.
7. `http.Server{}` inicia, com timeouts de leitura/escrita de 15s.

### 5.2 `domain/entities`

Cada arquivo define uma struct simples com tags `json:"..."` (ver tabela completa na seção 8 — modelagem de dados). Não têm métodos de comportamento (são *Plain Old Go Objects*); toda regra fica nos usecases.

### 5.3 `domain/repositories` (contratos)

11 interfaces, uma por agregado: `AnnouncementRepository`, `AttendanceRepository`, `CommentRepository`, `EventRepository`, `GuardianRepository`, `MilestoneRepository`, `PostRepository`, `StudentRepository`, `TimelineEventRepository`, `TurmaRepository`, `UserRepository`. Cada uma define os métodos de persistência que os usecases precisam, sem nenhuma referência a SQL.

### 5.4 `application/usecases` (40+ arquivos)

Padrão consistente em **todos** os usecases:

```go
type XxxUseCase struct { repo <Porta ou Repositório> }
func NewXxxUseCase(repo ...) *XxxUseCase { return &XxxUseCase{repo: repo} }
func (u *XxxUseCase) Execute(ctx context.Context, ...) (resultado, error) { ... }
```

Exemplos de regras de negócio encapsuladas:

| Usecase | Regra |
|---|---|
| `CreateAnnouncementUseCase` | Título e corpo obrigatórios; prioridade deve ser `Urgente`, `Importante` ou `Informativo` |
| `SaveAttendanceUseCase` | Status deve ser `present` ou `absent`; aluno e data obrigatórios |
| `CreateEventUseCase` / `UpdateEventUseCase` | Data de término não pode ser antes da de início |
| `RSVPEventUseCase` | Não permite confirmar presença em evento que já passou (`EndsAt < now()`) |
| `CreateMilestoneUseCase` | Categoria deve ser uma das 4 válidas (`Motor`, `Linguagem`, `Social`, `Cognitivo`) |
| `CreatePostUseCase` / `CreateStudentUseCase` / `UpdateProfileUseCase` | Validam tamanho de imagem em base64 (máx. ~4MB) via `validateImage` |
| `RegisterUserUseCase` | Impede e-mail duplicado; exige `role` válido; faz hash da senha com bcrypt antes de salvar |
| `LoginUseCase` | Erro genérico (`ErrInvalidCredentials`) tanto para e-mail inexistente quanto senha errada — evita enumeração de contas |
| `ForgotPasswordUseCase` | Sempre responde sucesso, mesmo se o e-mail não existir (placeholder — não há envio de e-mail implementado ainda) |

### 5.5 `adapters/http/handlers`

10 handlers, um por agregado HTTP. Todos seguem o mesmo padrão: decodificar JSON → montar entidade → chamar usecase → responder JSON com `writeJSON`/`writeError`. Funções auxiliares compartilhadas (definidas em `student_handler.go` mas usadas por todos, pois estão no mesmo pacote `handlers`):

- `parseIDParam(r)`: extrai `{id}` da URL via `mux.Vars`.
- `parseStudentIDQuery(r)`: extrai e valida `?student_id=` da query string.
- `requireProfissional(r)` / `userIDPointer(r)`: helpers de autorização baseados no contexto.

### 5.6 `adapters/http/middleware`

| Middleware | Função |
|---|---|
| `RequireAuth(secret)` | Bloqueia com 401 se não houver `Authorization: Bearer <token>` válido |
| `OptionalAuth(secret)` | Tenta decodificar o token, mas nunca bloqueia (usado em rotas públicas que mudam de comportamento se logado, ex: `/api/events`) |
| `WithCORS(origins)` | Reflete a origem permitida, trata `OPTIONS` (preflight) |
| `WithMetrics` | Envolve cada requisição contando (`requests_total`), medindo duração (`request_duration_seconds`) e contabilizando requisições simultâneas (`active_requests`) |

### 5.7 `adapters/database/postgres`

Implementação SQL pura (sem ORM) usando `database/sql` + `lib/pq`. Padrões notáveis:

- Cada repositório tem uma função `scanXxx` privada para converter `*sql.Row`/`*sql.Rows` em entidade, tratando campos `NULL` (`sql.NullString`, `sql.NullTime`, `sql.NullInt64`).
- `var _ repositories.XxxRepository = (*XxxRepository)(nil)` no final de cada arquivo garante, em tempo de compilação, que a struct implementa a interface.
- `PostRepository` implementa **duas** interfaces (`repositories.PostRepository` e, via métodos-ponte como `ListPosts`/`CreatePost`/`LikePost`, a interface `ports.PostPort`).
- `StudentRepository.Create` gera o `enrollment_code` (matrícula) **depois** do INSERT, pois o código (`#ano-id`) depende do ID gerado pelo banco.

### 5.8 `adapters/security`

- `jwt.go`: `GenerateToken` (HS256, claim customizada `Role` + `RegisteredClaims` padrão) e `ParseToken`.
- `password.go`: wrappers diretos do `bcrypt.GenerateFromPassword` / `bcrypt.CompareHashAndPassword`.

### 5.9 `internal/config/config.go`

Lê 5 variáveis de ambiente (`APP_PORT`, `DATABASE_URL`, `JWT_SECRET`, `CORS_ORIGIN`) com fallback de desenvolvimento hardcoded, e fixa `TokenTTL = 7 * 24h`.

---

## 6. Análise Detalhada do Frontend

### 6.1 Hierarquia de componentes

```
main.jsx
└── App.jsx                      (estado global: sessão, lista de filhos, navegação)
    ├── Login.jsx / CriarConta.jsx / EsqueceuSenha.jsx   (não autenticado)
    └── (autenticado)
        ├── Sidebar / BottomNav / Topbar      (definidos dentro de App.jsx)
        ├── ChildPicker → CreateChildForm     (profissional sem filho selecionado)
        ├── DashboardAdmin.jsx                (painel administrativo)
        ├── Feed.jsx                          (aba "Atividades")
        │     ├── Post, PostActions, PostDetail, Comment, NewPostForm
        ├── Jornada.jsx                       (aba "Jornada")
        │     ├── ProgressHeader, MilestoneCard, NewMilestoneForm
        ├── Agenda.jsx                        (aba "Agenda")
        │     ├── EventCard, NewEventForm
        ├── Comunicados.jsx
        │     ├── ComunicadoCard, NewAnnouncementForm
        ├── Perfil.jsx → Crianca.jsx → GuardianForm
        └── Turmas.jsx → ChamadaView (chamada de presença)
```

`SearchableSelect.jsx` é o único componente realmente "genérico/reutilizável", usado em `CreateChildForm` e `Crianca.jsx` para seleção de turma/professor/responsável com busca.

### 6.2 Gerenciamento de estado

Não há Redux, Zustand ou Context API de domínio. Cada página gerencia seu próprio estado local via `useState`, e `App.jsx` mantém o estado "global" mínimo necessário:

```js
user, students, selectedChildId, selectedStudent, timeline, nextEvent
```

Esse estado é repassado às páginas via props. Efeitos colaterais (chamadas à API) usam `useEffect` + `useCallback` para evitar loops de re-fetch.

### 6.3 `api.js` — camada de acesso à API

Único ponto de comunicação HTTP do frontend. Função central:

```js
async function request(path, { method = 'GET', body, auth = true } = {}) {
  // monta headers, injeta Authorization se auth=true e houver token
  // faz fetch, trata 204 (sem corpo), parseia JSON, lança Error(data.error) se !res.ok
}
```

Todas as ~40 funções exportadas (`login`, `getStudentById`, `createPost`, `rsvpEvent`, ...) são one-liners que chamam `request()` com o path/método certos. Também expõe `getToken`/`setToken` (persistência em `localStorage`, chave `up_espaco_token`) e `fileToBase64` (conversão de upload de imagem para data URI, usada em fotos de post/aluno/avatar).

### 6.4 Estratégia de renderização e rotas

Não há React Router. A "rota" é o estado `active` (string: `'Início'`, `'Atividades'`, `'Jornada'`, `'Agenda'`, `'Comunicados'`, `'Perfil'`, `'Crianca'`) controlado em `App.jsx`, e a função `renderMain()` decide qual componente mostrar via `if`/`return` sequenciais. É um *client-side router artesanal* — simples e suficiente para o número atual de telas.

### 6.5 Responsividade

Layout mobile-first com Tailwind: `Sidebar` (desktop, `hidden md:flex`) e `BottomNav` (mobile, `md:hidden`) coexistem, alternando conforme o breakpoint.

---

## 7. Fluxo de Execução da Aplicação

### Backend (subida do processo)

1. `main()` carrega config → abre conexão Postgres → `PingContext` (timeout 10s) → roda migrations → monta repositórios → monta router → inicia `http.Server` na porta configurada (padrão `8000`).
2. Se qualquer etapa de inicialização falhar (DB inacessível, migration quebrada), o processo encerra com `log.Fatalf` — *fail fast*.

### Frontend (carregamento inicial)

1. `main.jsx` monta `<App />` na div `#root`.
2. `App` roda um `useEffect` de "restaurar sessão": se existir token salvo, chama `getMe()`; se o token for inválido/expirado, limpa-o e mostra a tela de login.
3. Após autenticado, outro `useEffect` carrega a lista de alunos certa pro papel do usuário (`listStudents` para `profissional`, `getMyChildren` para `responsavel`).
4. Ao selecionar uma criança, `loadChildData` dispara em paralelo (`Promise.all`): dados do aluno, timeline do dia e lista de eventos (filtrando o próximo evento futuro no cliente).
5. Erros de rede são capturados localmente em cada componente (`try/catch` + `setError`), exibidos como banners vermelhos — não há um *error boundary* global.

### Tratamento de erros (padrão ponta a ponta)

```
Backend: usecase retorna error → handler responde {"error": "mensagem"} com status HTTP apropriado
Frontend: api.js detecta res.ok === false → throw new Error(data.error) → componente captura no catch e mostra na tela
```

---

## 8. Banco de Dados

PostgreSQL 15, schema evoluído por **18 migrations sequenciais** (`000001` a `000018`), aplicadas automaticamente no boot do backend e rastreadas na tabela `schema_migrations`.

### 8.1 Tabelas e relacionamentos

```
users (id) ──┬─< students.guardian_user_id   (responsável → filhos)
             ├─< students.teacher_user_id     (professor → alunos)
             ├─< attendance.marked_by_user_id
             ├─< comments.user_id
             ├─< announcement_reads.user_id
             └─< event_rsvps.user_id

turmas (id) ──< students.turma_id

students (id) ─┬─< student_guardians.student_id   (responsáveis autorizados a buscar)
                ├─< posts.student_id
                ├─< timeline_events.student_id
                ├─< milestones.student_id
                └─< attendance.student_id

posts (id) ──< comments.post_id

events (id) ──< event_rsvps.event_id

announcements (id) ──< announcement_reads.announcement_id
```

### 8.2 Tabela `users`

| Coluna | Tipo | Observação |
|---|---|---|
| id | SERIAL PK | |
| name, email | TEXT | `email` é `UNIQUE` |
| password_hash | TEXT | bcrypt |
| role | TEXT | `CHECK (role IN ('profissional','responsavel'))` |
| phone, address, avatar_url | TEXT | default `''` |
| created_at, updated_at | TIMESTAMPTZ | |

### 8.3 Tabela `students` (consolidada após migrations 0001, 0007, 0014)

| Coluna | Tipo | Observação |
|---|---|---|
| id | SERIAL PK | |
| name | TEXT | |
| presence_status | TEXT | default `'absent'` |
| check_in_at | TIMESTAMPTZ NULL | preenchido ao marcar presença |
| guardian_user_id | INT → users(id) | `ON DELETE SET NULL` |
| teacher_user_id | INT → users(id) | `ON DELETE SET NULL` |
| turma_id | INT → turmas(id) | `ON DELETE SET NULL` |
| photo_url, group_name, teacher_name, enrollment_code, blood_type, restrictions, medications | TEXT | default `''` |
| birth_date | DATE NULL | |
| allergies | TEXT[] | default `'{}'` |
| created_at, updated_at | TIMESTAMPTZ | |

### 8.4 Tabela `attendance` (migration 0016)

| Coluna | Tipo | Observação |
|---|---|---|
| id | SERIAL PK | |
| student_id | INT → students(id) | `ON DELETE CASCADE` |
| date | DATE | |
| status | TEXT | `present` / `absent` |
| marked_by_user_id | INT → users(id) NULL | |
| **UNIQUE(student_id, date)** | | garante 1 registro por aluno/dia (usado no `ON CONFLICT DO UPDATE`) |
| índice `idx_attendance_student_date` | | otimiza consultas por aluno+data |

### 8.5 Demais tabelas (resumo)

| Tabela | Campos-chave | Observação |
|---|---|---|
| `posts` | student_id, title, description, pedagogical_note, image_url, likes, bookmarks, visibility | `visibility`: `private` ou `turma` (mural compartilhado) |
| `comments` | post_id, user_id (nullable), author_name, avatar_url, text | autor pode ser anônimo/externo (campos denormalizados) |
| `timeline_events` | student_id, title, description, occurred_at | eventos do "dia a dia" |
| `milestones` | student_id, title, category, description, achieved_at, done | categoria restrita a 4 valores no usecase |
| `events` | title, description, location, starts_at, ends_at, rsvp_count | agenda da escola |
| `event_rsvps` | event_id, user_id, **UNIQUE(event_id, user_id)** | confirmação de presença |
| `announcements` | title, sender, priority, preview, body, attachment_name | `priority`: Urgente/Importante/Informativo |
| `announcement_reads` | announcement_id, user_id | controla quem já leu |
| `turmas` | name | salas/grupos da escola |
| `student_guardians` | student_id, name, relation, phone, authorized | contatos autorizados a buscar a criança (distinto de `guardian_user_id`, que é a conta com login) |

### 8.6 Migrations e seeds

- Runner customizado em `migrations/migrate.go`: embute os `.up.sql` no binário (`//go:embed`), cria a tabela `schema_migrations` se não existir, aplica em ordem alfabética dentro de uma transação por arquivo.
- Existem 3 migrations de seed (`000005`, `000012`, `000018`) que populam dados de demonstração/apresentação.
- Cada `.up.sql` tem seu `.down.sql` correspondente (rollback), embora o runner atual só execute `up` (não há comando de rollback exposto).

---

## 9. APIs (Referência de Endpoints)

> Base URL: `http://localhost:8000` (dev). Autenticação via header `Authorization: Bearer <token>` (omitido nas tabelas quando a rota aceita acesso anônimo).

### Autenticação (`/api/auth`, `/api/me`, `/api/users`)

| Método | Rota | Auth | Descrição |
|---|---|---|---|
| POST | `/api/auth/register` | não | Cria conta. Body: `{name, email, password, role}` |
| POST | `/api/auth/login` | não | Login. Body: `{email, password}` → `{token, user}` |
| POST | `/api/auth/forgot-password` | não | Recuperação de senha (placeholder, sempre 200) |
| GET | `/api/me` | sim | Dados do usuário logado |
| PUT | `/api/me` | sim | Atualiza nome/telefone/endereço/avatar |
| GET | `/api/me/children` | sim | Filhos vinculados (papel `responsavel`) |
| GET | `/api/users?role=` | sim (profissional) | Lista usuários por papel |

### Alunos (`/api/student`, `/api/students`)

| Método | Rota | Auth | Descrição |
|---|---|---|---|
| GET | `/api/student` | não | Aluno "ativo" (modo demo/single-student) |
| GET | `/api/students` | não | Lista todos os alunos |
| POST | `/api/students` | sim (profissional) | Cadastra aluno |
| GET | `/api/students/{id}` | não | Detalhe do aluno + responsáveis |
| PUT | `/api/students/{id}` | sim (profissional) | Atualiza aluno |
| DELETE | `/api/students/{id}` | sim (profissional) | Remove aluno |
| PATCH | `/api/students/{id}/presence` | não* | Marca presença/falta. Body: `{status}` |
| GET/POST | `/api/students/{id}/guardians` | GET não / POST sim | Lista/cria responsáveis autorizados |
| PUT/DELETE | `/api/guardians/{id}` | sim | Edita/remove responsável |
| POST/GET | `/api/students/{id}/attendance` | sim | Registra/lista histórico de presença |

\* `UpdatePresence` não passa por `requireAuth` na definição de rota, mas idealmente deveria — ver seção 10 (pontos de atenção de segurança).

### Turmas (`/api/turmas`)

| Método | Rota | Auth | Descrição |
|---|---|---|---|
| GET | `/api/turmas` | não | Lista turmas |
| POST | `/api/turmas` | sim (profissional) | Cria turma |
| PUT/DELETE | `/api/turmas/{id}` | sim (profissional) | Edita/remove turma |
| GET | `/api/turmas/{id}/attendance?date=` | sim (profissional) | Presença da turma num dia |

### Timeline (`/api/timeline`)

| Método | Rota | Auth | Descrição |
|---|---|---|---|
| GET | `/api/timeline?student_id=` | não | Eventos de hoje do aluno |
| POST | `/api/timeline?student_id=` | sim | Cria evento |
| GET/PUT/DELETE | `/api/timeline/{id}` | GET não / outros sim | Detalhe/edita/remove |

### Posts e comentários (`/api/posts`, `/api/comments`)

| Método | Rota | Auth | Descrição |
|---|---|---|---|
| GET | `/api/posts?student_id=` | não | Lista posts do aluno (+ posts de turma) |
| POST | `/api/posts?student_id=` | sim | Cria post |
| GET/PUT/DELETE | `/api/posts/{id}` | GET não / outros sim | Detalhe/edita/remove |
| POST | `/api/posts/{id}/like` \| `/unlike` | não | Curtir/descurtir |
| POST | `/api/posts/{id}/bookmark` \| `/unbookmark` | não | Salvar/remover dos favoritos |
| GET | `/api/posts/{id}/comments` | não | Lista comentários |
| POST | `/api/posts/{id}/comments` | sim | Cria comentário |
| DELETE | `/api/comments/{id}` | sim | Remove comentário |

### Eventos (`/api/events`)

| Método | Rota | Auth | Descrição |
|---|---|---|---|
| GET | `/api/events` | opcional | Lista eventos (marca RSVP se logado) |
| POST | `/api/events` | sim (profissional) | Cria evento |
| GET/PUT/DELETE | `/api/events/{id}` | GET não / outros sim | Detalhe/edita/remove |
| POST | `/api/events/{id}/rsvp` | sim | Confirma/cancela presença |

### Comunicados (`/api/announcements`)

| Método | Rota | Auth | Descrição |
|---|---|---|---|
| GET | `/api/announcements` | opcional | Lista (marca lidos se logado) |
| POST | `/api/announcements` | sim (profissional) | Cria |
| GET/PUT/DELETE | `/api/announcements/{id}` | opcional/sim/sim | Detalhe/edita/remove |
| POST | `/api/announcements/{id}/read` | sim | Marca como lido |

### Marcos (`/api/milestones`)

| Método | Rota | Auth | Descrição |
|---|---|---|---|
| GET | `/api/milestones?student_id=` | não | Lista marcos do aluno |
| POST | `/api/milestones` | sim | Cria (body inclui `student_id`) |
| PUT/DELETE | `/api/milestones/{id}` | sim | Edita/remove |

### Observabilidade e documentação

| Rota | Descrição |
|---|---|
| `GET /metrics` | Métricas Prometheus (`requests_total`, `request_duration_seconds`, `active_requests`) |
| `GET /swagger` | UI do Swagger (carregado via CDN) |
| `GET /swagger.json` | Especificação OpenAPI estática (`backend-up/docs/swagger.json`) |

### Formato de erro padrão

```json
{ "error": "mensagem legível em português" }
```

Códigos HTTP usados: `400` (validação), `401` (sem token/token inválido), `403` (papel sem permissão), `404` (não encontrado), `409` (conflito, ex: e-mail em uso), `500` (erro interno).

---

## 10. Segurança

### 10.1 Autenticação

- JWT assinado com **HS256**, segredo único (`JWT_SECRET`, env var).
- Claims: `Subject` (userID), `Role`, `IssuedAt`, `ExpiresAt` (TTL fixo de **7 dias**, `internal/config/config.go`).
- Token enviado pelo cliente como `Authorization: Bearer <token>`, extraído em `middleware.extractToken`.

### 10.2 Autorização

- Baseada em **papel** (`role`), não em permissões granulares — checagens do tipo `if role != "profissional"` espalhadas nos handlers (ex: `requireProfissional`, checagens inline em `student_handler.go`, `attendance_handler.go`).
- Não há verificação de propriedade fina em todas as rotas (ex: qualquer usuário autenticado pode comentar em qualquer post; a posse de "ser o pai daquele aluno" não é checada em todos os endpoints de leitura, já que vários são intencionalmente públicos/`OptionalAuth` para simplificar a navegação do feed).

### 10.3 Armazenamento de credenciais

- Senhas nunca armazenadas em texto puro: `bcrypt.GenerateFromPassword` com custo padrão (`bcrypt.DefaultCost`).
- `entities.User.PasswordHash` tem a tag `json:"-"` — nunca é serializado nas respostas da API.

### 10.4 Proteção contra enumeração

- `LoginUseCase` devolve sempre `ErrInvalidCredentials` (não diferencia "email não existe" de "senha errada").
- `ForgotPasswordUseCase` sempre responde sucesso, independente do e-mail existir.

### 10.5 CORS

- `middleware.WithCORS` lê uma lista de origens permitidas separada por vírgula (`CORS_ORIGIN`), reflete a origem exata se houver match (ou usa `*` apenas se explicitamente configurado), trata `OPTIONS` (preflight) com `204`.

### 10.6 Limites de payload

- `validateImage` (em `application/usecases/validation.go`) rejeita imagens base64 acima de ~4MB antes de persistir — mitigação simples contra payloads excessivos no banco.

### 10.7 Pontos de atenção (gaps conhecidos)

- `JWT_SECRET` e credenciais do Postgres têm **valores padrão hardcoded** em `config.go` e `docker-compose.yml` (`up-espaco-dev-secret-change-me`, `postgres/postgres`) — adequado para desenvolvimento, **deve ser sobrescrito em produção** via variáveis de ambiente reais/secrets manager.
- Algumas rotas sensíveis (ex: `PATCH /api/students/{id}/presence`) não exigem autenticação no roteamento (`routes.go` usa `HandleFunc` em vez de `requireAuth(...)`), embora o aluno marcado seja identificado só pelo ID na URL.
- Não há rate limiting nem proteção contra brute-force no login.
- Não há CSRF token (mitigado parcialmente pelo uso de Bearer token em vez de cookies de sessão).

---

## 11. Configuração e Infraestrutura

### 11.1 Variáveis de ambiente (backend)

| Variável | Default (dev) | Descrição |
|---|---|---|
| `APP_PORT` | `8000` | Porta HTTP do servidor |
| `DATABASE_URL` | `postgres://postgres:postgres@db:5432/up_espaco?sslmode=disable` | DSN de conexão Postgres |
| `JWT_SECRET` | `up-espaco-dev-secret-change-me` | Segredo de assinatura do JWT |
| `CORS_ORIGIN` | `*` | Lista de origens permitidas (separadas por vírgula) |

### 11.2 Variáveis de ambiente (frontend)

| Variável | Default | Descrição |
|---|---|---|
| `VITE_API_URL` | `http://localhost:8000` | Base URL da API consumida pelo `api.js` |

### 11.3 Docker / Docker Compose

`docker-compose.yml` define 7 serviços:

| Serviço | Imagem/Build | Porta exposta (localhost) |
|---|---|---|
| `db` | `postgres:15-alpine` | 5432 |
| `backend` | build de `backend-up/Dockerfile` | 8000 |
| `frontend` | build de `frontend-up/Dockerfile` | 5173 |
| `postgres-exporter` | `prometheuscommunity/postgres-exporter` | 9187 |
| `cadvisor` | `gcr.io/cadvisor/cadvisor` | 8080 |
| `prometheus` | `prom/prometheus` | 9090 |
| `grafana` | `grafana/grafana` (admin/admin) | 3000 |

Volumes nomeados: `db_data`, `prometheus_data`, `grafana_data` (persistência entre restarts).

O backend monta as migrations como volume (`./backend-up/internal/adapters/database/migrations:/migrations`) no container do `db` apenas como referência — a aplicação real das migrations acontece dentro do binário Go via `//go:embed`, não por esse mount.

### 11.4 Dockerfiles

- **Backend**: `golang:1.22-alpine` → `go mod download` → `go build -o /app/backend ./cmd/api` → binário executado direto (sem multi-stage build, então a imagem final inclui o toolchain do Go).
- **Frontend**: `node:20-alpine` → `npm install` → roda em modo dev (`npm run dev -- --host 0.0.0.0`) — ou seja, a imagem Docker atual serve o **Vite dev server**, não um build de produção estático.

### 11.5 CI/CD

Não há pipeline de CI/CD configurado no repositório (sem `.github/workflows`, `.gitlab-ci.yml` etc.) — deploy e testes são manuais atualmente.

---

## 12. Observabilidade (Prometheus + Grafana)

### 12.1 Métricas expostas pelo backend (`/metrics`)

| Métrica | Tipo | Labels | Significado |
|---|---|---|---|
| `requests_total` | Counter | `method`, `path`, `status` | total de requisições HTTP |
| `request_duration_seconds` | Histogram | `method`, `path` | latência das requisições |
| `active_requests` | Gauge | — | requisições em andamento no momento |

### 12.2 Scrape (`prometheus.yml`)

3 jobs: `up-espaco-backend` (`backend:8000/metrics`), `up-espaco-postgres` (`postgres-exporter:9187`), `up-espaco-containers` (`cadvisor:8080`), intervalo de 15s.

### 12.3 Dashboard Grafana (`up-espaco-observability.json`)

6 painéis pré-configurados:
1. Requisições HTTP por segundo (por método/path).
2. Proporção de status HTTP (pizza).
3. CPU do container backend.
4. Memória RAM do container backend.
5. Saúde do PostgreSQL (`pg_up`, conexões ativas).
6. Latência HTTP p95 por path.

Datasource provisionado automaticamente (`grafana/provisioning/datasources/prometheus.yml`), apontando para `http://prometheus:9090`.

---

## 13. Dependências do Projeto

### Backend (`go.mod`)

| Dependência | Finalidade | Necessidade real |
|---|---|---|
| `github.com/golang-jwt/jwt/v5` | Autenticação JWT | Essencial |
| `github.com/gorilla/mux` | Roteamento HTTP | Essencial (rotas com parâmetros e métodos) |
| `github.com/lib/pq` | Driver Postgres | Essencial |
| `github.com/prometheus/client_golang` | Métricas | Importante para observabilidade em produção |
| `golang.org/x/crypto` | bcrypt | Essencial para segurança de senhas |
| (indiretas: `beorn7/perks`, `cespare/xxhash`, `golang/protobuf`, `prometheus/client_model`, `prometheus/common`, `prometheus/procfs`, `golang.org/x/sys`, `google.golang.org/protobuf`) | Transitivas do client Prometheus | Não usadas diretamente pelo código |

### Frontend (`package.json`)

| Dependência | Finalidade | Necessidade real |
|---|---|---|
| `react` / `react-dom` | UI | Essencial |
| `lucide-react` | Ícones | Usado extensivamente em toda a UI — essencial visualmente, mas substituível |
| `vite` | Build/dev server | Essencial |
| `@vitejs/plugin-react` | Suporte JSX/Fast Refresh no Vite | Essencial |
| `tailwindcss` + `@tailwindcss/vite` | Estilização | Essencial (toda a UI usa classes Tailwind) |

Nenhuma dependência de roteamento, gerenciamento de estado, formulários ou requisições HTTP foi adicionada — o projeto optou por implementar essas necessidades manualmente (`fetch` nativo em `api.js`, navegação por estado em `App.jsx`), mantendo o bundle final pequeno.

---

## 14. Regras de Negócio

### Explícitas (validadas em usecases)

- E-mail de usuário é único; papel deve ser `profissional` ou `responsavel`.
- Status de presença só pode ser `present` ou `absent`.
- Prioridade de comunicado só pode ser `Urgente`, `Importante` ou `Informativo`.
- Categoria de marco só pode ser `Motor`, `Linguagem`, `Social` ou `Cognitivo`.
- Evento: data de término deve ser posterior à de início.
- RSVP: não é possível confirmar presença em evento já encerrado.
- Imagens (post, aluno, perfil) limitadas a ~4MB em base64.
- Post/aluno: título e campos descritivos obrigatórios.

### Implícitas (na modelagem/queries)

- Um responsável só vê os alunos vinculados ao seu `user_id` (`ListByGuardian`), nunca a lista completa.
- Posts com `visibility = 'turma'` aparecem para todos os alunos da mesma turma (`turma_id`), não só para o aluno-autor — implementado via subquery em `PostRepository.List`.
- Presença é única por aluno+dia (`UNIQUE(student_id, date)`); marcar de novo **atualiza** o registro existente em vez de duplicar (`ON CONFLICT ... DO UPDATE`).
- Código de matrícula (`enrollment_code`) é gerado automaticamente no formato `#<ano>-<id com 4 dígitos>` e nunca é alterado depois de criado.
- RSVP é uma alternância (toggle): chamar a rota duas vezes confirma e depois cancela a presença, ajustando `rsvp_count` atomicamente dentro de uma transação.
- Exclusão de aluno, turma, evento, etc. propaga via `ON DELETE CASCADE`/`SET NULL` conforme a tabela (ex: excluir um aluno remove seus posts/timeline/attendance em cascata; excluir uma turma só desvincula os alunos, sem excluí-los).

### Fluxos críticos

1. **Cadastro de aluno por profissional**: pode vincular um responsável existente por ID ou buscando por e-mail (`CreateStudentRequest.GuardianEmail`); resolve turma e professor por ID, validando que o "professor" selecionado realmente tem `role = profissional`.
2. **Chamada de presença por turma**: tela `Turmas.jsx` → `ChamadaView` lista os alunos da turma filtrando localmente (`students.filter(s => s.turma_id === turma.id)`) e faz toggle otimista de presença, revertendo em caso de erro de rede.

---

## 15. Resumo Executivo

### Resumo geral

O **Up - Espaço** é uma aplicação full-stack (Go + React + Postgres) que digitaliza a comunicação escola-família em creches/escolas infantis, com feed de atividades, agenda, comunicados, controle de presença e acompanhamento de desenvolvimento infantil. A arquitetura do backend segue princípios de Clean Architecture com separação clara entre domínio, casos de uso e adapters, sem frameworks de injeção de dependência — tudo montado explicitamente. O frontend é uma SPA React simples, sem roteador ou state manager externos, adequada ao escopo atual do produto.

### Principais tecnologias

Go 1.22 + gorilla/mux + PostgreSQL 15 (backend) · React 18 + Vite 6 + Tailwind 4 (frontend) · JWT + bcrypt (segurança) · Prometheus + Grafana + cAdvisor (observabilidade) · Docker Compose (orquestração).

### Principais módulos

Autenticação/Perfil · Alunos & Turmas · Presença/Chamada · Feed de Posts (curtidas/comentários/favoritos) · Timeline diária · Jornada (marcos) · Agenda/Eventos (RSVP) · Comunicados · Observabilidade.

### Fluxo principal

Login → seleção de aluno (ou lista completa, se profissional) → navegação entre Início/Atividades/Jornada/Agenda/Perfil, com toda persistência via API REST autenticada por JWT.

### Pontos fortes

- Separação de camadas consistente e didática no backend — fácil de localizar onde uma regra de negócio vive.
- Modelagem de dados coerente, com migrations incrementais bem documentadas pelo próprio nome do arquivo.
- Observabilidade já embutida desde o início (métricas + dashboards), algo raro em projetos deste porte.
- Validações de negócio centralizadas nos usecases, reaproveitáveis independente do transporte (HTTP).
- Frontend simples e direto, sem overengineering para o tamanho atual do produto.

### Pontos de atenção

- Segredos de produção (`JWT_SECRET`, senha do Postgres) ainda com defaults de desenvolvimento no código/compose — precisam ser injetados via secrets reais antes de qualquer deploy público.
- Algumas rotas de escrita sensíveis (ex.: marcar presença) não exigem autenticação no nível de roteamento.
- Ausência de testes automatizados (unitários/integração) no repositório.
- Ausência de pipeline de CI/CD.
- Frontend sem TypeScript e sem roteador formal — viável hoje, mas limita a escalabilidade de manutenção se o número de telas crescer.
- Dockerfile do frontend roda o dev server do Vite em produção (não há build estático + servidor como Nginx).
- Autorização é só por papel (profissional/responsável), sem checagem fina de propriedade em todas as rotas (ex.: comentar em post de aluno de outra família).

### Recomendações futuras

1. Extrair segredos para variáveis de ambiente reais/secret manager antes de qualquer deploy fora do ambiente de desenvolvimento.
2. Adicionar testes unitários nos `usecases` (são a camada mais fácil de testar isoladamente, por já não depender de HTTP/SQL).
3. Adicionar um multi-stage build no Dockerfile do backend (build em `golang:alpine`, runtime em `scratch`/`alpine` puro) para reduzir o tamanho da imagem final.
4. Criar um build de produção real do frontend (`vite build` + Nginx/Caddy) em vez de servir o dev server.
5. Introduzir checagem de propriedade (ex.: confirmar que o `responsavel` autenticado realmente está vinculado ao aluno antes de permitir certas leituras/escritas).
6. Configurar um pipeline básico de CI (lint + build + testes) antes de cada merge.
7. Avaliar a adoção de React Router se o número de telas/rotas continuar crescendo, e TypeScript se a equipe de frontend crescer.
