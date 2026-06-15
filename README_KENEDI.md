# 📊 Plano de Atividades: Kenedi (Database & Observability)

Olá, Kenedi! Como desenvolvedor responsável pelo **Banco de Dados (PostgreSQL) & Observabilidade (Prometheus + Grafana)** do projeto acadêmico **UP Espaço**, seu principal objetivo será estruturar o armazenamento de dados e a infraestrutura de monitoramento, essencial para garantirmos o ponto extra bônus do projeto!

Abaixo estão listadas as suas **7 Responsabilidades Clínicas & Técnicas** na equipe:

---

## 📋 Suas 7 Responsabilidades no Projeto

### 1. 🗄️ Esquema Relacional de Dados DDL
Modelar e criar as tabelas físicas no PostgreSQL correspondentes às entidades necessárias do sistema:
- `students` (id, name, avatar, classroom, isPresent, entranceTime)
- `timeline_items` (id, title, time, description, icon, category, color)
- `activity_posts` (id, author_name, author_avatar, author_role, time, title, category, description, imageUrl, likes, commentsCount)
- `calendar_events` (id, title, time, isImportant, icon, color, rsvpPrompt, isRSVPed)

### 2. ⚡ Script de Inicialização de Dados (Seeders)
Criar e manter scripts de população inicial do banco com dados realistas (usando como base as informações originais do arquivo [src/data.ts](src/data.ts)), para que o projeto inicie com posts e timelines de demonstração.

### 3. 🐳 Gestão de Volumes e Persistência do Postgres
Garantir que as tabelas e dados inseridos no PostgreSQL não sejam apagados ao derrubar os contêineres Docker. Isso será feito mantendo e ajustando o volume `postgres_data` declarado em nosso `docker-compose.yml`.

### 🎛️ 4. Configuração do Coletor Prometheus
Configurar e gerenciar o arquivo [prometheus.yml](prometheus.yml) para que o Prometheus colete dados de telemetria da máquina do backend e do banco de dados a cada 15 segundos.

### 📈 5. Criação do Painel Visual no Grafana
Após iniciar o contêiner do Grafana (porta `3001`), você deve acessar a interface web e criar um Dashboard de Monitoramento contendo gráficos dinâmicos:
- Quantidade de requisições HTTP tratadas.
- Gráfico de pizza mostrando a proporção de status HTTP (200 OK vs 500 Erros).
- Uso de memória ram e processamento do contêiner da aplicação.

### 💾 6. Exportação e Persistência dos Gráficos do Grafana
Garantir que os Dashboards criados no Grafana sejam persistidos no volume Docker `grafana_data` para que o professor e o grupo consigam visualizá-los imediatamente sem precisar reconfigurar as telas do Grafana do zero.

### 🤝 7. Apresentação do Bônus de Monitoramento
Liderar a explicação da seção de Observabilidade no slide e demonstrar ao vivo para o professor o painel do Grafana se mexendo em tempo real enquanto o PO (Brayan) navega pela aplicação e faz requisições.

---

## 🛠️ Como rodar o seu ambiente localmente
1. Certifique-se de ter o Docker e Docker Compose instalados na máquina.
2. Na raiz do projeto, execute o comando para subir todos os contêineres de observabilidade e banco:
   ```bash
   docker-compose up --build
   ```
3. Acesse o Prometheus em `http://localhost:9090` e o Grafana em `http://localhost:3001` (Credenciais padrão do Grafana: login `admin`, senha `admin`).
