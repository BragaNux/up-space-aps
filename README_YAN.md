# 💻 Plano de Atividades: Yan (Backend & API Integration)

Olá, Yan! Como desenvolvedor responsável pelo **Backend & Integração de API** do projeto acadêmico **UP Espaço**, seu principal objetivo será construir e conectar o servidor que gerencia os dados do nosso portal de desenvolvimento infantil.

Abaixo estão listadas as suas **7 Responsabilidades Clínicas & Técnicas** na equipe:

---

## 📋 Suas 7 Responsabilidades no Projeto

### 1. ⚙️ Desenvolvimento da API RESTful
Você deverá programar o servidor backend utilizando a tecnologia/linguagem acordada com o grupo (ex. Node.js com Express, Python/FastAPI ou C#/.NET). O servidor deve expor endpoints funcionais que atendam ao frontend:
- `GET /api/student`: Retorna os dados do Enzo (o aluno ativo).
- `POST /api/student/presence`: Altera o status de presença (Faltoso / Presente).
- `GET /api/timeline`: Retorna a linha do tempo do dia.
- `GET /api/posts`: Retorna a lista de diários e atividades dos terapeutas.
- `POST /api/posts`: Registra um novo diário de atividade.
- `POST /api/posts/:id/like` e `/bookmark`: Manipula as interações dos pais nos posts.
- `GET /api/events` e `POST /api/events/:id/rsvp`: Controla a agenda clínica e RSVP de eventos.

### 2. 🗄️ Integração com o PostgreSQL
Configurar o ORM (como Prisma, Sequelize, Mongoose) ou drivers SQL nativos no backend para se comunicar diretamente com o contêiner `db` (PostgreSQL) configurado no `docker-compose.yml`.

### 3. 📐 Modelagem de Entidades no Código
Mapear e codificar as entidades e tabelas no backend correspondendo às interfaces do TypeScript definidas pelo frontend (consulte o arquivo [src/types.ts](file:///c:/Users/Braga/Desktop/frontend_up/src/types.ts)).

### 4. 🧪 Regras de Negócio e Validações
Garantir que o backend trate corretamente as regras antes de salvar no banco de dados (ex: impedir o cadastro de posts sem título ou sem nota pedagógica, atualizar horários de check-in automaticamente ao dar presença no dia).

### 5. 🛡️ Configuração de CORS e Segurança
Configurar as permissões de CORS (Cross-Origin Resource Sharing) no backend para aceitar requisições originárias da porta do frontend (`http://localhost:3000`), evitando bloqueios de navegadores.

### 6. 📊 Endpoint de Métricas para Observabilidade (Requisito Bônus)
Criar e expor um endpoint público `/metrics` no backend utilizando bibliotecas padrões (ex. `prom-client` para Node.js ou `prometheus-client` para Python). O Prometheus utilizará esse endpoint para raspar métricas de performance (latência, CPU, chamadas por segundo) para desenharmos no Grafana.

### 🤝 7. Teste de Integração Ponta a Ponta
Trabalhar junto com o PO (Braga) para trocar a variável de ambiente `VITE_API_URL` para o endereço da sua API e validar se a transição dos dados simulados locais para a sua API real está acontecendo de forma limpa e sem bugs na UI.

---

## 🛠️ Como rodar o seu ambiente localmente
1. O backend deve ler a variável de conexão `DATABASE_URL` (configurada no `docker-compose.yml`).
2. Garanta que a sua porta mapeada no container seja a `8000`.
3. Siga o fluxo detalhado no [README.md](file:///c:/Users/Braga/Desktop/frontend_up/README.md) principal para subir o banco de dados.
