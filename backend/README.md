# UP Espaço Backend

Backend REST para o portal UP Espaço, implementado em Go com arquitetura hexagonal (Ports & Adapters).

---

## Visão geral

O backend e o banco de dados são normalmente executados em containers (Docker). Opcionalmente, um desenvolvedor pode usar `devbox` para um ambiente local reproduzível sem depender do Docker.

---

## Pré-requisitos

- Docker & Docker Compose (recomendado)
- Devbox (opcional)
- Ou Go 1.22 e PostgreSQL local caso queira rodar sem Docker

---

## Rápido (modo recomendado) — Docker

1. Copie o exemplo de variáveis de ambiente se quiser personalizar:

```bash
cp .env.example .env
```

2. Subir os serviços (backend + postgres):

```bash
docker compose up --build -d
```

3. As migrations não são aplicadas automaticamente pelo Compose — aplicar as migrations dentro do container Postgres (uma vez):

```bash
# Executar as migrations incluídas na imagem do Postgres
docker exec -i backend-db-1 psql -U postgres -d up_espaco -f /migrations/000001_create_students.up.sql
docker exec -i backend-db-1 psql -U postgres -d up_espaco -f /migrations/000002_create_timeline_events.up.sql
docker exec -i backend-db-1 psql -U postgres -d up_espaco -f /migrations/000003_create_posts.up.sql
docker exec -i backend-db-1 psql -U postgres -d up_espaco -f /migrations/000004_create_events.up.sql
docker exec -i backend-db-1 psql -U postgres -d up_espaco -f /migrations/000005_seed_initial_data.up.sql
```

4. Verifique que a API está rodando:

```bash
curl http://localhost:8000/api/student
```

5. Para ver logs (em tempo real):

```bash
docker compose logs -f backend
```

6. Parar e remover os containers:

```bash
docker compose down
```

---

## Rodando sem Docker (opcional, com Devbox)

Use `devbox` para garantir que a máquina do dev tenha as dependências corretas sem instalar manualmente.

1. Abra o shell do devbox:

```bash
devbox shell
```

2. Configure variáveis de ambiente localmente (por exemplo, no terminal do devbox):

```bash
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/up_espaco?sslmode=disable
export APP_PORT=8000
```

> Observação: você precisará de um PostgreSQL acessível em `DATABASE_URL`. Pode ser uma instância local instalada ou remota.

3. Aplicar migrations (opção 1 — usando `psql`):

```bash
psql "$DATABASE_URL" -f internal/adapters/database/migrations/000001_create_students.up.sql
psql "$DATABASE_URL" -f internal/adapters/database/migrations/000002_create_timeline_events.up.sql
psql "$DATABASE_URL" -f internal/adapters/database/migrations/000003_create_posts.up.sql
psql "$DATABASE_URL" -f internal/adapters/database/migrations/000004_create_events.up.sql
psql "$DATABASE_URL" -f internal/adapters/database/migrations/000005_seed_initial_data.up.sql
```

4. Build e run do backend (dentro do `devbox shell`):

```bash
# build
go build -o backend ./cmd/api
# executar
./backend
```

Ou usar `go run` durante desenvolvimento:

```bash
go run ./cmd/api
```

5. Alternativa: usar `golang-migrate` (instalado no `devbox.json`) para aplicar todas as migrations de uma vez:

```bash
migrate -path internal/adapters/database/migrations -database "$DATABASE_URL" up
```

---

## Arquivos importantes

- Migrations: `internal/adapters/database/migrations`
- Código principal: `cmd/api/main.go`
- Handlers: `internal/adapters/http/handlers`
- Repositórios Postgres: `internal/adapters/database/postgres`
- Use cases: `internal/application/usecases`

---

## Variáveis de ambiente

Ver ` .env.example` para valores de exemplo:

```
APP_PORT=8000
DATABASE_URL=postgres://postgres:postgres@db:5432/up_espaco?sslmode=disable
```

---

## Observações e dicas

- Para desenvolvimento iterativo, é comum usar um bind-mount para sobrescrever o código no container e rodar `go run` dentro do container, mas o `docker-compose.yml` deste repositório foi ajustado para executar o binário compilado dentro da imagem (evita problemas com diferenças de ambiente). Se preferir montar o código localmente para desenvolvimento rápido, restaure o volume `.:/app` no `docker-compose.yml` e ajuste o `command` para `go run ./cmd/api`.

- CORS: o backend permite origem `http://localhost:3000`.

---

## Swagger (documentação)

Adicionado um endpoint simples de documentação OpenAPI/Swagger UI.

- UI: `http://localhost:8000/swagger`
- JSON OpenAPI: `http://localhost:8000/swagger.json`

O arquivo `docs/swagger.json` contém uma especificação mínima e foi montado no container via volume (`./docs:/app/docs`) no `docker-compose.yml`. Isso permite editar `docs/swagger.json` localmente sem rebuildar a imagem — basta recriar o serviço:

```bash
docker compose up -d --force-recreate backend
```

Observação: o Swagger UI carrega a biblioteca via CDN, então é necessário acesso à internet para a interface.
