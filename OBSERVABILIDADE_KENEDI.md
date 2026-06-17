# UP Espaco - Banco de Dados e Observabilidade

Este pacote entrega a parte de PostgreSQL, Prometheus e Grafana.

## Como rodar

Na raiz do projeto `up-space-aps`, execute:

```bash
docker compose up --build
```

Servicos principais:

- Backend: http://localhost:8000
- Metricas do backend: http://localhost:8000/metrics
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3001
- Login Grafana: `admin`
- Senha Grafana: `admin`

## Banco de dados

O PostgreSQL usa o volume Docker `postgres_data`, entao os dados continuam existindo apos `docker compose down`.

Quando o volume esta vazio, o arquivo `backend/internal/adapters/database/init/001_schema_and_seed.sql` cria as tabelas, visoes e dados iniciais automaticamente.

As migrations tambem foram mantidas em `backend/internal/adapters/database/migrations` para execucao manual quando necessario.

### Entidades do enunciado

O backend atual consulta as tabelas `students`, `timeline_events`, `posts` e `events`. Para nao quebrar essas consultas, o esquema foi ampliado e foram criadas visoes com os nomes pedidos no plano:

- `timeline_items`, baseada em `timeline_events`
- `activity_posts`, baseada em `posts`
- `calendar_events`, baseada em `events`

## Prometheus

O arquivo `prometheus.yml` coleta metricas a cada 15 segundos dos seguintes alvos:

- `backend:8000/metrics`, com requisicoes HTTP e latencia
- `postgres-exporter:9187`, com metricas do PostgreSQL
- `cadvisor:8080`, com CPU e memoria dos containers

## Grafana

O Grafana usa o volume `grafana_data` para persistencia e tambem recebe provisioning automatico em `grafana/provisioning`.

O dashboard `UP Espaco - Observabilidade` ja nasce configurado com:

- Requisicoes HTTP por segundo
- Grafico de pizza por status HTTP
- CPU do container do backend
- Memoria RAM do container do backend
- Saude e conexoes do PostgreSQL
- Latencia HTTP p95

## Roteiro rapido para a apresentacao

1. Abrir o Grafana em http://localhost:3001.
2. Entrar no dashboard `UP Espaco - Observabilidade`.
3. Acessar algumas rotas do backend, como `/api/student`, `/api/posts`, `/api/timeline` e `/api/events`.
4. Mostrar o painel de requisicoes mudando em tempo real.
5. Explicar que o Prometheus coleta os dados a cada 15 segundos e que o Grafana apenas visualiza essas series temporais.
