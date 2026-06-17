CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    avatar TEXT NULL,
    classroom TEXT NULL,
    presence_status TEXT NOT NULL DEFAULT 'absent',
    check_in_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS timeline_events (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL,
    icon TEXT NULL,
    category TEXT NULL,
    color TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    author_name TEXT NOT NULL DEFAULT 'Equipe UP Espaco',
    author_avatar TEXT NULL,
    author_role TEXT NOT NULL DEFAULT 'Educador(a)',
    title TEXT NOT NULL,
    category TEXT NULL,
    description TEXT NOT NULL,
    image_url TEXT NULL,
    pedagogical_note TEXT NOT NULL,
    likes BIGINT NOT NULL DEFAULT 0,
    comments_count BIGINT NOT NULL DEFAULT 0,
    bookmarks BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    location TEXT NOT NULL,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    is_important BOOLEAN NOT NULL DEFAULT false,
    icon TEXT NULL,
    color TEXT NULL,
    rsvp_prompt TEXT NULL,
    is_rsvped BOOLEAN NOT NULL DEFAULT false,
    rsvp_count BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE OR REPLACE VIEW timeline_items AS
SELECT
    id,
    title,
    to_char(occurred_at AT TIME ZONE 'America/Sao_Paulo', 'HH24:MI') AS "time",
    description,
    icon,
    category,
    color
FROM timeline_events;

CREATE OR REPLACE VIEW activity_posts AS
SELECT
    id,
    author_name,
    author_avatar,
    author_role,
    to_char(created_at AT TIME ZONE 'America/Sao_Paulo', 'HH24:MI') AS "time",
    title,
    category,
    description,
    image_url AS "imageUrl",
    likes,
    comments_count AS "commentsCount"
FROM posts;

CREATE OR REPLACE VIEW calendar_events AS
SELECT
    id,
    title,
    to_char(starts_at AT TIME ZONE 'America/Sao_Paulo', 'DD/MM HH24:MI') AS "time",
    is_important AS "isImportant",
    icon,
    color,
    rsvp_prompt AS "rsvpPrompt",
    is_rsvped AS "isRSVPed"
FROM events;

INSERT INTO students (id, name, avatar, classroom, presence_status, check_in_at)
VALUES
    (1, 'Enzo', 'https://images.unsplash.com/photo-1503454537195-1dcabb73ffb9?auto=format&fit=crop&w=120&q=80', 'Sala Sensorial A', 'present', now() - interval '5 hours 20 minutes')
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    avatar = EXCLUDED.avatar,
    classroom = EXCLUDED.classroom,
    presence_status = EXCLUDED.presence_status,
    check_in_at = EXCLUDED.check_in_at,
    updated_at = now();

SELECT setval('students_id_seq', GREATEST((SELECT max(id) FROM students), 1), true);

INSERT INTO timeline_events (id, title, description, occurred_at, icon, category, color)
VALUES
    (1, 'Chegada acolhedora', 'Enzo chegou tranquilo, guardou a mochila e cumprimentou a equipe.', date_trunc('day', now()) + interval '7 hours 30 minutes', 'LogIn', 'Rotina', '#2563eb'),
    (2, 'Atividade sensorial', 'Exploracao de texturas com tecidos, espuma e pecas macias.', date_trunc('day', now()) + interval '8 hours 40 minutes', 'Sparkles', 'Terapia', '#16a34a'),
    (3, 'Lanche da manha', 'Participou do lanche com boa interacao e autonomia.', date_trunc('day', now()) + interval '9 hours 30 minutes', 'Apple', 'Alimentacao', '#f97316'),
    (4, 'Musica e movimento', 'Acompanhou o ritmo com palmas e pequenos deslocamentos guiados.', date_trunc('day', now()) + interval '10 hours 20 minutes', 'Music', 'Pedagogico', '#7c3aed'),
    (5, 'Relaxamento', 'Fez pausa com respiracao guiada antes da saida.', date_trunc('day', now()) + interval '11 hours 45 minutes', 'Moon', 'Bem-estar', '#0f766e')
ON CONFLICT (id) DO UPDATE SET
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    occurred_at = EXCLUDED.occurred_at,
    icon = EXCLUDED.icon,
    category = EXCLUDED.category,
    color = EXCLUDED.color;

SELECT setval('timeline_events_id_seq', GREATEST((SELECT max(id) FROM timeline_events), 1), true);

INSERT INTO posts (id, author_name, author_avatar, author_role, title, category, description, image_url, pedagogical_note, likes, comments_count, bookmarks, created_at, updated_at)
VALUES
    (1, 'Marina Lopes', 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?auto=format&fit=crop&w=120&q=80', 'Terapeuta Ocupacional', 'Atividade Sensorial', 'Terapia', 'Brincadeira com texturas variadas para estimular percepcao tatil e tolerancia a novos materiais.', 'https://images.unsplash.com/photo-1587654780291-39c9404d746b?auto=format&fit=crop&w=900&q=80', 'Enzo respondeu bem aos estimulos e aceitou repetir a exploracao por mais tempo.', 12, 3, 2, now() - interval '3 hours', now() - interval '3 hours'),
    (2, 'Rafael Nunes', 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?auto=format&fit=crop&w=120&q=80', 'Educador', 'Hora da Leitura', 'Pedagogico', 'Leitura de historia com fantoches, perguntas simples e apoio visual.', 'https://images.unsplash.com/photo-1519682337058-a94d519337bc?auto=format&fit=crop&w=900&q=80', 'A atencao se manteve durante 15 minutos, com respostas por apontamento.', 9, 2, 1, now() - interval '2 hours', now() - interval '2 hours'),
    (3, 'Bianca Rocha', 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?auto=format&fit=crop&w=120&q=80', 'Musicoterapeuta', 'Musica e Movimento', 'Expressao', 'Sessao de musica com ritmo suave, instrumentos pequenos e movimentos guiados.', 'https://images.unsplash.com/photo-1516280440614-37939bbacd81?auto=format&fit=crop&w=900&q=80', 'Enzo sincronizou alguns movimentos e buscou interacao com os colegas.', 15, 5, 4, now() - interval '1 hour', now() - interval '1 hour')
ON CONFLICT (id) DO UPDATE SET
    author_name = EXCLUDED.author_name,
    author_avatar = EXCLUDED.author_avatar,
    author_role = EXCLUDED.author_role,
    title = EXCLUDED.title,
    category = EXCLUDED.category,
    description = EXCLUDED.description,
    image_url = EXCLUDED.image_url,
    pedagogical_note = EXCLUDED.pedagogical_note,
    likes = EXCLUDED.likes,
    comments_count = EXCLUDED.comments_count,
    bookmarks = EXCLUDED.bookmarks,
    created_at = EXCLUDED.created_at,
    updated_at = EXCLUDED.updated_at;

SELECT setval('posts_id_seq', GREATEST((SELECT max(id) FROM posts), 1), true);

INSERT INTO events (id, title, description, location, starts_at, ends_at, is_important, icon, color, rsvp_prompt, is_rsvped, rsvp_count)
VALUES
    (1, 'Reuniao com pais', 'Encontro sobre progresso, metas da semana e combinados de rotina.', 'Sala 3', date_trunc('day', now()) + interval '1 day 14 hours', date_trunc('day', now()) + interval '1 day 15 hours', true, 'Users', '#dc2626', 'Confirmar presenca ate hoje?', false, 6),
    (2, 'Sessao de Terapia', 'Acompanhamento com terapeuta ocupacional.', 'Clinica B', date_trunc('day', now()) + interval '2 days 10 hours', date_trunc('day', now()) + interval '2 days 11 hours', false, 'HeartPulse', '#16a34a', 'Deseja reservar este horario?', false, 3),
    (3, 'Oficina Pedagogica', 'Atividade em grupo para desenvolvimento socioemocional.', 'Sala 5', date_trunc('day', now()) + interval '3 days 9 hours', date_trunc('day', now()) + interval '3 days 11 hours', false, 'Palette', '#2563eb', 'Confirmar participacao?', false, 8)
ON CONFLICT (id) DO UPDATE SET
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    location = EXCLUDED.location,
    starts_at = EXCLUDED.starts_at,
    ends_at = EXCLUDED.ends_at,
    is_important = EXCLUDED.is_important,
    icon = EXCLUDED.icon,
    color = EXCLUDED.color,
    rsvp_prompt = EXCLUDED.rsvp_prompt,
    is_rsvped = EXCLUDED.is_rsvped,
    rsvp_count = EXCLUDED.rsvp_count;

SELECT setval('events_id_seq', GREATEST((SELECT max(id) FROM events), 1), true);
