INSERT INTO students (name, presence_status, check_in_at) VALUES ('Enzo', 'present', now());

INSERT INTO posts (title, description, pedagogical_note) VALUES
('Atividade Sensorial', 'Brincadeira com texturas variadas', 'O Enzo respondeu bem aos estímulos.'),
('Hora da Leitura', 'Leitura de história com fantoches', 'A atenção se manteve durante 15 minutos.'),
('Música e Movimento', 'Sessão de música com ritmo suave', 'O Enzo sincronizou alguns movimentos.');

INSERT INTO timeline_events (title, description, occurred_at) VALUES
('Chegada na escola', 'Enzo chegou animado e pronto para o dia.', now() - interval '5 hour'),
('Aula de Artes', 'Atividade de pintura livre.', now() - interval '4 hour'),
('Lanche da manhã', 'Momento de socialização e alimentação.', now() - interval '3 hour'),
('Atividade Física', 'Jogos leves com supervisão.', now() - interval '2 hour'),
('Relaxamento', 'Tempo de descanso com meditação guiada.', now() - interval '1 hour');

INSERT INTO events (title, description, location, starts_at, ends_at) VALUES
('Reunião com pais', 'Encontro sobre progresso e metas.', 'Sala 3', now() + interval '1 day', now() + interval '1 day 1 hour'),
('Sessão de Terapia', 'Acompanhamento com terapeuta ocupacional.', 'Clínica B', now() + interval '2 day', now() + interval '2 day 1 hour'),
('Oficina Pedagógica', 'Atividade em grupo para desenvolvimento socioemocional.', 'Sala 5', now() + interval '3 day', now() + interval '3 day 2 hours');
