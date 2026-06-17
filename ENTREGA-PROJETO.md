# Up Espaço - Documento de Entrega do Projeto

> Este documento cobre os entregáveis 1 a 6 do trabalho (Descrição do Problema, Levantamento de Requisitos, Arquitetura, Design Patterns, Product Backlog e Implementação).

---

## Organização da Equipe

| Papel | Integrante | Responsabilidade principal |
|---|---|---|
| Product Owner + Desenvolvedor | **Brayan** | Visão de produto, priorização do backlog, distribuição de atividades, acompanhamento e validação das entregas; desenvolveu o **Frontend** (React) |
| Desenvolvedor | **Yan** | Desenvolveu o **Backend** (API REST em Go, regras de negócio, autenticação) |
| Desenvolvedor | **Kenedi** | Desenvolveu o **Banco de Dados** (modelagem, migrations) e a **Observabilidade** (Prometheus + Grafana) |

As responsabilidades acima foram a divisão principal, mas a equipe trabalhou de forma colaborativa, todos contribuíram em pontos fora da sua frente principal sempre que necessário (ex.: ajustes de integração frontend-backend, revisão de queries, configuração de métricas).

---

## 1. Descrição do Problema

### Contexto

Escolas e creches de educação infantil normalmente comunicam o dia a dia das crianças aos responsáveis por canais informais e fragmentados: grupos de WhatsApp, agendas de papel, recados verbais na porta da escola. Isso gera perda de informação, falta de histórico, dificuldade de auditoria (quem leu o quê, quem confirmou presença em um evento) e nenhum controle estruturado de presença, saúde ou desenvolvimento da criança.

### Público-alvo

- **Profissionais da escola** (professores, equipe pedagógica, direção): precisam registrar presença, publicar atividades, criar eventos e comunicados, e manter o cadastro dos alunos atualizado.
- **Responsáveis** (pais, mães, cuidadores legais): precisam acompanhar o dia da criança, ver fotos/atividades, confirmar presença em eventos e ler avisos da escola, tudo restrito apenas aos filhos vinculados à própria conta.

### Objetivo da aplicação

O **Up - Espaço** é uma plataforma web (API REST + SPA) que centraliza a comunicação escola-família em um único sistema, oferecendo:

- Feed de atividades pedagógicas (estilo rede social, com fotos, curtidas e comentários).
- Controle de presença/falta diária por aluno e por turma.
- Linha do tempo diária e marcos de desenvolvimento (jornada) de cada criança.
- Agenda de eventos da escola com confirmação de presença (RSVP).
- Comunicados institucionais com controle de leitura.
- Cadastro de alunos, turmas, responsáveis autorizados e dados de saúde.

---

## 2. Levantamento de Requisitos

### Requisitos Funcionais (RF)

| ID | Descrição | Status |
|---|---|---|
| **RF01** | O sistema deve permitir que um usuário se **cadastre e faça login** (e-mail/senha), recebendo um token de acesso (JWT) válido para uso nas demais funcionalidades. | ✅ Implementado |
| **RF02** | O sistema deve permitir que um profissional **cadastre, edite e remova alunos**, vinculando-os a uma turma, um professor e um responsável. | ✅ Implementado |
| **RF03** | O sistema deve permitir **registrar a presença ou falta** de um aluno em uma data, individualmente ou via chamada por turma. | ✅ Implementado |
| **RF04** | O sistema deve permitir **publicar atividades no feed** (com foto, descrição e nota pedagógica) e que os responsáveis **curtam, comentem e salvem** essas publicações. | ✅ Implementado |
| **RF05** | O sistema deve permitir a criação de **eventos na agenda da escola** e que os usuários **confirmem presença (RSVP)**, com contador de confirmados. | ✅ Implementado |

### Requisitos Não Funcionais (RNF)

| ID | Descrição | Status |
|---|---|---|
| **RNF01** | **Segurança**: as senhas dos usuários devem ser armazenadas com hash (bcrypt), e o acesso à API deve ser autenticado por token JWT assinado, nunca trafegando a senha em texto puro após o login. | ✅ Implementado |
| **RNF02** | **Desempenho/Monitoramento**: o sistema deve expor métricas de tempo de resposta e volume de requisições, permitindo identificar lentidão ou picos de uso em tempo real. | ✅ Implementado (Prometheus) |
| **RNF03** | **Usabilidade**: a interface deve ser responsiva, adaptando a navegação para desktop (menu lateral) e mobile (menu inferior), sem necessidade de telas diferentes por dispositivo. | ✅ Implementado |
| **RNF04** | **Observabilidade**: o sistema deve possuir dashboards visuais de monitoramento (uso de CPU/memória do backend, saúde do banco de dados, latência e taxa de erros) acessíveis sem necessidade de acessar logs manualmente. | ✅ Implementado (Grafana) |
| **RNF05** | **Portabilidade/Implantação**: a aplicação completa (frontend, backend, banco de dados e monitoramento) deve poder ser executada em qualquer ambiente com um único comando, via containers. | ✅ Implementado (Docker Compose) |

> Documentação técnica completa de cada requisito (campos, regras de validação, fluxo) está detalhada no [README.md](README.md), seções 5, 9 e 14.

---

## 3. Arquitetura

### Padrão escolhido: Arquitetura em Camadas (variação de Clean Architecture / Hexagonal)

O backend foi estruturado em **camadas concêntricas com inversão de dependência**, próxima do padrão *Ports & Adapters* (Arquitetura Hexagonal), uma evolução da clássica Arquitetura em Camadas:

```
domain        → entidades e contratos puros (sem dependências externas)
application   → casos de uso (regras de negócio) + DTOs + ports
adapters      → implementações concretas (HTTP, Postgres, JWT, bcrypt)
```

A regra de dependência é sempre de fora para dentro: `adapters` conhece `application`, que conhece `domain`; o inverso nunca acontece. O `domain` não importa nem o driver do banco de dados.

```
Cliente (React)
      │ HTTP/JSON
      ▼
adapters/http (handlers, middleware, routes)   ← camada de apresentação
      │
      ▼
application/usecases                            ← camada de regras de negócio
      │
      ▼
domain/repositories (interfaces)                ← contrato, sem implementação
      ▲
      │ implementa
adapters/database/postgres                       ← camada de persistência
```

O **frontend** segue uma arquitetura **Cliente-Servidor** simples: a SPA React consome a API REST via uma camada única de acesso HTTP (`api.js`), sem lógica de negócio replicada no cliente, toda regra de negócio (validações, permissões) vive no backend.

### Justificativa da escolha

- **Separação de responsabilidades clara**: cada camada tem um motivo único para mudar (regra de negócio muda em `usecases`; troca de banco mudaria só em `adapters/database`; troca de protocolo mudaria só em `adapters/http`).
- **Testabilidade**: como `domain` e `application` não dependem de HTTP nem de SQL, os casos de uso podem ser testados isoladamente, passando implementações falsas das interfaces de repositório.
- **Facilidade de manutenção em equipe**: como a equipe dividiu o trabalho por área (Backend, Frontend, Banco de Dados/Observabilidade), a separação em camadas e a definição de contratos (interfaces) permitiu que cada integrante trabalhasse de forma relativamente independente, sem conflitos constantes de código.
- **Simplicidade deliberada**: optamos por não usar um framework de injeção de dependência, a montagem das peças é feita manualmente no `main.go` (Composition Root), o que mantém o fluxo de inicialização fácil de seguir, ponto importante para um projeto acadêmico onde todo o time precisa entender o código de ponta a ponta.

---

## 4. Design Patterns

A equipe identificou e aplicou deliberadamente os seguintes padrões de projeto no código:

### 4.1 Repository Pattern

**Problema que resolve**: evitar que as regras de negócio (`usecases`) fiquem acopladas aos detalhes de acesso a dados (SQL, driver do banco). Sem esse padrão, qualquer mudança de banco de dados ou de biblioteca de acesso exigiria reescrever a lógica de negócio inteira.

**Onde foi aplicado**: em `backend-up/internal/domain/repositories/` (11 interfaces, ex.: `StudentRepository`, `PostRepository`, `EventRepository`) e suas implementações concretas em `backend-up/internal/adapters/database/postgres/` (ex.: `postgres.StudentRepository`). Os `usecases` dependem apenas da interface, nunca da implementação:

```go
// domain/repositories/student_repository.go - o contrato
type StudentRepository interface {
    GetByID(ctx context.Context, id int64) (*entities.Student, error)
    Create(ctx context.Context, student *entities.Student) error
    // ...
}

// adapters/database/postgres/student_repository.go - a implementação real
type StudentRepository struct{ db *DB }
func (r *StudentRepository) GetByID(ctx context.Context, id int64) (*entities.Student, error) {
    // SQL real aqui
}
var _ repositories.StudentRepository = (*StudentRepository)(nil) // garante em compile-time
```

**Por que foi escolhido**: é o padrão mais natural para a arquitetura em camadas adotada, permite trocar Postgres por outro banco no futuro sem tocar em nenhuma regra de negócio, e facilita testes unitários dos `usecases` com repositórios falsos (mocks).

### 4.2 Decorator Pattern (middlewares HTTP)

**Problema que resolve**: adicionar comportamentos transversais (autenticação, CORS, coleta de métricas) a um handler HTTP sem modificar o handler original e sem duplicar código em cada rota.

**Onde foi aplicado**: em `backend-up/internal/adapters/http/middleware/`. Cada middleware é uma função que recebe um `http.Handler` e devolve outro `http.Handler` "decorado":

```go
// cada middleware envolve o handler anterior, adicionando um comportamento
wrapped := middleware.WithMetrics(middleware.WithCORS(cfg.CORSOrigin)(handler))

// dentro de cada rota autenticada:
router.Handle("/api/students", requireAuth(http.HandlerFunc(studentHandler.Create)))
```

Cada chamada (`WithMetrics`, `WithCORS`, `RequireAuth`) envolve a anterior numa cadeia (`handler → CORS(handler) → Metrics(CORS(handler))`), exatamente a estrutura do padrão Decorator.

**Por que foi escolhido**: é o padrão idiomático em Go para HTTP (usado por praticamente todo framework web da linguagem), permite combinar middlewares livremente por rota (ex.: algumas rotas usam `RequireAuth`, outras `OptionalAuth`, outras nenhum) e mantém cada responsabilidade (autenticação, métricas, CORS) isolada em um arquivo próprio.

### 4.3 (Complementar) Factory Method - Construtores `NewXxx`

Como reforço além dos dois padrões principais, todo objeto do backend (usecases, handlers, repositórios) é criado por uma função construtora dedicada (`NewCreatePostUseCase`, `NewPostsHandler`, `NewPostRepository`...), que centraliza a montagem do objeto e suas dependências, em vez de instanciar structs diretamente com `&Struct{}` espalhado pelo código. Isso facilita a adição futura de novos parâmetros de inicialização sem quebrar quem já usa o construtor outros lugares.

---

## 5. Product Backlog

| ID | Descrição | Prioridade | Responsável |
|---|---|---|---|
| RF01 | Cadastro e login de usuário (JWT) | Alta | Yan (backend) / Brayan (frontend) |
| RF02 | Cadastro, edição e remoção de alunos | Alta | Yan (backend) / Brayan (frontend) |
| RF03 | Registro de presença/falta (individual e por chamada de turma) | Alta | Kenedi (banco de dados) / Yan (backend) |
| RF04 | Feed de atividades com curtidas, comentários e favoritos | Média | Brayan (frontend) / Yan (backend) |
| RF05 | Agenda de eventos com confirmação de presença (RSVP) | Média | Yan (backend) / Brayan (frontend) |
| RNF01 | Segurança: hash de senha (bcrypt) e autenticação JWT | Alta | Yan (backend) |
| RNF02 | Monitoramento de desempenho (métricas Prometheus) | Média | Kenedi |
| RNF03 | Interface responsiva (desktop e mobile) | Média | Brayan |
| RNF04 | Observabilidade com dashboards (Grafana) | Alta | Kenedi |
| RNF05 | Containerização completa (Docker Compose) | Média | Kenedi |
| - | Modelagem do banco de dados (18 migrations incrementais) | Alta | Kenedi |
| - | Documentação técnica da arquitetura e da API | Média | Brayan (PO) |

---

## 6. Implementação

### Status

A aplicação **compila e executa corretamente** via `docker-compose up -d --build` (backend Go, frontend React/Vite, PostgreSQL, Prometheus e Grafana sobem juntos) ou de forma nativa (backend com `go run ./cmd/api`, frontend com `npm run dev`, desde que haja um Postgres acessível).

Todos os requisitos funcionais e não funcionais listados na seção 2 estão implementados e demonstráveis na aplicação atual:

- **Arquitetura escolhida** (camadas / hexagonal) está aplicada de forma consistente em todo o backend, ver detalhamento de cada camada e arquivo em [README.md, seção 5](README.md#5-análise-detalhada-do-backend).
- **Design Patterns** (Repository e Decorator) estão aplicados nos pontos descritos na seção 4 deste documento, com exemplos de código reais do projeto.
- O **frontend** consome 100% das funcionalidades via API REST documentada em [README.md, seção 9](README.md#9-apis-referência-de-endpoints).
- O **bônus de Banco de Dados + Observabilidade** está coberto: PostgreSQL com 18 migrations versionadas (seção 8 do README) e stack completa de observabilidade com Prometheus + Grafana + cAdvisor + postgres_exporter, com dashboard pronto de 6 painéis (seção 12 do README).

### Como executar

```bash
# clonar o repositório e, na raiz do projeto:
docker compose up -d --build

# Frontend:    http://localhost:5173
# Backend:     http://localhost:8000  (Swagger em /swagger)
# Grafana:     http://localhost:3000  (usuário/senha: admin/admin)
# Prometheus:  http://localhost:9090
```

### Onde encontrar mais detalhes técnicos

Toda a análise profunda de código (arquivos, funções, fluxos, modelagem do banco e referência completa de endpoints) está no [README.md](README.md) deste repositório, que serve como documentação técnica de apoio a este documento de entrega.
