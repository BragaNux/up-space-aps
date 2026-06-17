-- 1. Insert Turmas
INSERT INTO turmas (id, name, created_at, updated_at) VALUES
(1, 'Turma das Borboletas (3-4 anos)', now(), now()),
(2, 'Turma dos Leões (4-5 anos)', now(), now()),
(3, 'Turma dos Golfinhos (5-6 anos)', now(), now())
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name;
SELECT setval('turmas_id_seq', COALESCE((SELECT MAX(id) FROM turmas), 1));

-- 2. Insert Users (senha123)
-- Hash bcrypt de "senha123": $2a$10$eDDuus0hIwc38PanTSeAxu.v7cIGIYlWx4wNG.CgQlfvqczSntOdK
INSERT INTO users (id, name, email, password_hash, role, phone, address, avatar_url, created_at, updated_at) VALUES
(1, 'Ana Souza (Professora)', 'ana@upespaco.com.br', '$2a$10$eDDuus0hIwc38PanTSeAxu.v7cIGIYlWx4wNG.CgQlfvqczSntOdK', 'profissional', '(11) 97777-6666', 'Rua Educativa, 100', 'https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=150', now(), now()),
(2, 'Carlos Lima (Psicopedagogo)', 'carlos@upespaco.com.br', '$2a$10$eDDuus0hIwc38PanTSeAxu.v7cIGIYlWx4wNG.CgQlfvqczSntOdK', 'profissional', '(11) 98888-7777', 'Av. Pedagógica, 200', 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150', now(), now()),
(6, 'Beatriz Rocha (Psicomotricista)', 'beatriz@upespaco.com.br', '$2a$10$eDDuus0hIwc38PanTSeAxu.v7cIGIYlWx4wNG.CgQlfvqczSntOdK', 'profissional', '(11) 97777-1111', 'Av. das Acácia, 50', 'https://images.unsplash.com/photo-1580489944761-15a19d654956?w=150', now(), now()),
(7, 'Marcos Silva (Ed. Física)', 'marcos@upespaco.com.br', '$2a$10$eDDuus0hIwc38PanTSeAxu.v7cIGIYlWx4wNG.CgQlfvqczSntOdK', 'profissional', '(11) 98888-2222', 'Rua Esportiva, 12', 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=150', now(), now()),

(3, 'Mariana Costa (Mãe)', 'mariana@gmail.com', '$2a$10$eDDuus0hIwc38PanTSeAxu.v7cIGIYlWx4wNG.CgQlfvqczSntOdK', 'responsavel', '(11) 98765-4321', 'Rua das Flores, 123', 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=150', now(), now()),
(4, 'Pedro Rocha (Pai)', 'pedro@gmail.com', '$2a$10$eDDuus0hIwc38PanTSeAxu.v7cIGIYlWx4wNG.CgQlfvqczSntOdK', 'responsavel', '(11) 99887-7665', 'Av. Paulista, 1000', 'https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?w=150', now(), now()),
(5, 'Luciana Mello (Mãe)', 'luciana@gmail.com', '$2a$10$eDDuus0hIwc38PanTSeAxu.v7cIGIYlWx4wNG.CgQlfvqczSntOdK', 'responsavel', '(11) 95555-4444', 'Alameda Santos, 500', 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=150', now(), now()),
(8, 'Roberto Santos (Pai)', 'roberto@gmail.com', '$2a$10$eDDuus0hIwc38PanTSeAxu.v7cIGIYlWx4wNG.CgQlfvqczSntOdK', 'responsavel', '(11) 95555-8888', 'Av. Interlagos, 1200', 'https://images.unsplash.com/photo-1492562080023-ab3db95bfbce?w=150', now(), now()),
(9, 'Fernanda Souza (Mãe)', 'fernanda@gmail.com', '$2a$10$eDDuus0hIwc38PanTSeAxu.v7cIGIYlWx4wNG.CgQlfvqczSntOdK', 'responsavel', '(11) 96666-9999', 'Av. Sto Amaro, 800', 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=150', now(), now()),
(10, 'Juliana Ribeiro (Mãe)', 'juliana@gmail.com', '$2a$10$eDDuus0hIwc38PanTSeAxu.v7cIGIYlWx4wNG.CgQlfvqczSntOdK', 'responsavel', '(11) 97777-0000', 'Rua Bela Cintra, 45', 'https://images.unsplash.com/photo-1524504388940-b1c1722653e1?w=150', now(), now()),
(11, 'Ricardo Oliveira (Pai)', 'ricardo@gmail.com', '$2a$10$eDDuus0hIwc38PanTSeAxu.v7cIGIYlWx4wNG.CgQlfvqczSntOdK', 'responsavel', '(11) 98888-1111', 'Av. Rebouças, 1500', 'https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?w=150', now(), now()),
(12, 'Camila Martins (Mãe)', 'camila@gmail.com', '$2a$10$eDDuus0hIwc38PanTSeAxu.v7cIGIYlWx4wNG.CgQlfvqczSntOdK', 'responsavel', '(11) 99999-2222', 'Rua Augusta, 300', 'https://images.unsplash.com/photo-1531746020798-e6953c6e8e04?w=150', now(), now())
ON CONFLICT (id) DO UPDATE SET 
    name = EXCLUDED.name, email = EXCLUDED.email, password_hash = EXCLUDED.password_hash, 
    role = EXCLUDED.role, phone = EXCLUDED.phone, address = EXCLUDED.address, avatar_url = EXCLUDED.avatar_url;
SELECT setval('users_id_seq', COALESCE((SELECT MAX(id) FROM users), 1));

-- 3. Insert Students (Turma das Borboletas = 1, Turma dos Leões = 2, Turma dos Golfinhos = 3)
INSERT INTO students (id, name, presence_status, check_in_at, guardian_user_id, photo_url, turma_id, group_name, teacher_user_id, teacher_name, birth_date, enrollment_code, blood_type, allergies, restrictions, medications, created_at, updated_at) VALUES
(1, 'Alice Costa', 'present', now(), 3, 'https://images.unsplash.com/photo-1503919545889-aef636e10ad4?w=200', 1, 'Turma das Borboletas (3-4 anos)', 1, 'Ana Souza (Professora)', '2021-04-12', '#2026-0001', 'A+', ARRAY['Amendoim', 'Poeira'], 'Sem açúcar refinado', 'Nenhum', now(), now()),
(2, 'Bernardo Costa', 'present', now(), 3, 'https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?w=200', 2, 'Turma dos Leões (4-5 anos)', 2, 'Carlos Lima (Psicopedagogo)', '2020-09-23', '#2026-0002', 'O-', ARRAY[]::TEXT[], 'Sem lactose', 'Nenhum', now(), now()),
(3, 'Davi Rocha', 'present', now(), 4, 'https://images.unsplash.com/photo-1471286174890-9c112ffca5b4?w=200', 1, 'Turma das Borboletas (3-4 anos)', 1, 'Ana Souza (Professora)', '2021-08-05', '#2026-0003', 'B-', ARRAY['Glúten'], 'Dieta celíaca', 'Nenhum', now(), now()),
(4, 'Enzo Mello', 'absent', NULL, 5, 'https://images.unsplash.com/photo-1519238263530-99bdd11df2ea?w=200', 3, 'Turma dos Golfinhos (5-6 anos)', 1, 'Ana Souza (Professora)', '2019-12-15', '#2026-0004', 'AB+', ARRAY[]::TEXT[], 'Nenhuma', 'Nenhum', now(), now()),

(5, 'Gabriel Santos', 'present', now(), 8, 'https://images.unsplash.com/photo-1513551573338-107695245848?w=200', 1, 'Turma das Borboletas (3-4 anos)', 1, 'Ana Souza (Professora)', '2021-02-18', '#2026-0005', 'O+', ARRAY[]::TEXT[], 'Nenhuma', 'Nenhum', now(), now()),
(6, 'Laura Souza', 'present', now(), 9, 'https://images.unsplash.com/photo-1516627145497-ae6968895b74?w=200', 1, 'Turma das Borboletas (3-4 anos)', 1, 'Ana Souza (Professora)', '2021-06-05', '#2026-0006', 'A-', ARRAY['Picada de abelha'], 'Nenhuma', 'Anti-histamínico em SOS', now(), now()),
(7, 'Manuela Ribeiro', 'present', now(), 10, 'https://images.unsplash.com/photo-1518887570146-0612132dd618?w=200', 1, 'Turma das Borboletas (3-4 anos)', 1, 'Ana Souza (Professora)', '2021-09-30', '#2026-0007', 'B+', ARRAY[]::TEXT[], 'Nenhuma', 'Nenhum', now(), now()),
(8, 'Pedro Oliveira', 'present', now(), 11, 'https://images.unsplash.com/photo-1502082553048-f009c37129b9?w=200', 1, 'Turma das Borboletas (3-4 anos)', 1, 'Ana Souza (Professora)', '2021-11-12', '#2026-0008', 'O+', ARRAY['Corantes artificiais'], 'Evitar refrigerantes', 'Nenhum', now(), now()),
(9, 'Sofia Martins', 'present', now(), 12, 'https://images.unsplash.com/photo-1484186139897-d5fc6b908812?w=200', 1, 'Turma das Borboletas (3-4 anos)', 1, 'Ana Souza (Professora)', '2021-01-22', '#2026-0009', 'A+', ARRAY[]::TEXT[], 'Nenhuma', 'Nenhum', now(), now()),

(10, 'Lucas Lima', 'present', now(), 8, 'https://images.unsplash.com/photo-1536640788329-e21c8d41cfc8?w=200', 2, 'Turma dos Leões (4-5 anos)', 2, 'Carlos Lima (Psicopedagogo)', '2020-05-14', '#2026-0010', 'B-', ARRAY['Pelo de gato'], 'Nenhuma', 'Nenhum', now(), now()),
(11, 'Beatriz Costa', 'present', now(), 9, 'https://images.unsplash.com/photo-1544816155-12df9643f363?w=200', 2, 'Turma dos Leões (4-5 anos)', 2, 'Carlos Lima (Psicopedagogo)', '2020-07-28', '#2026-0011', 'O-', ARRAY[]::TEXT[], 'Sem açúcar refinado', 'Nenhum', now(), now()),
(12, 'Henrique Rocha', 'present', now(), 4, 'https://images.unsplash.com/photo-1485546246426-74dc88dec4d9?w=200', 3, 'Turma dos Golfinhos (5-6 anos)', 1, 'Ana Souza (Professora)', '2019-10-09', '#2026-0012', 'A-', ARRAY['Mofo', 'Poeira'], 'Nenhuma', 'Bombinha de asma em SOS', now(), now())
ON CONFLICT (id) DO UPDATE SET 
    name = EXCLUDED.name, presence_status = EXCLUDED.presence_status, check_in_at = EXCLUDED.check_in_at, 
    guardian_user_id = EXCLUDED.guardian_user_id, photo_url = EXCLUDED.photo_url, turma_id = EXCLUDED.turma_id, 
    group_name = EXCLUDED.group_name, teacher_user_id = EXCLUDED.teacher_user_id, teacher_name = EXCLUDED.teacher_name, 
    birth_date = EXCLUDED.birth_date, enrollment_code = EXCLUDED.enrollment_code, blood_type = EXCLUDED.blood_type, 
    allergies = EXCLUDED.allergies, restrictions = EXCLUDED.restrictions, medications = EXCLUDED.medications;
SELECT setval('students_id_seq', COALESCE((SELECT MAX(id) FROM students), 1));

-- 4. Insert Student Guardians
INSERT INTO student_guardians (id, student_id, name, relation, phone, avatar_url, authorized, created_at) VALUES
(1, 1, 'Vovó Regina', 'Avó', '(11) 91111-2222', 'https://images.unsplash.com/photo-1551836022-d5d88e9218df?w=150', true, now()),
(2, 1, 'Tio Marcos', 'Tio', '(11) 92222-3333', 'https://images.unsplash.com/photo-1492562080023-ab3db95bfbce?w=150', false, now()),
(3, 3, 'Vovô Geraldo', 'Avô', '(11) 93333-4444', 'https://images.unsplash.com/photo-1484186139897-d5fc6b908812?w=150', true, now()),
(4, 5, 'Vovó Neusa', 'Avó', '(11) 94444-5555', 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=150', true, now()),
(5, 6, 'Tia Paula', 'Tia', '(11) 95555-6666', 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=150', true, now()),
(6, 9, 'Tio Fábio', 'Tio', '(11) 96666-7777', 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150', false, now())
ON CONFLICT (id) DO UPDATE SET 
    student_id = EXCLUDED.student_id, name = EXCLUDED.name, relation = EXCLUDED.relation, 
    phone = EXCLUDED.phone, avatar_url = EXCLUDED.avatar_url, authorized = EXCLUDED.authorized;
SELECT setval('student_guardians_id_seq', COALESCE((SELECT MAX(id) FROM student_guardians), 1));

-- 5. Insert Manual Posts (Initial ones)
INSERT INTO posts (id, student_id, title, description, pedagogical_note, image_url, likes, bookmarks, visibility, created_at, updated_at) VALUES
(1, 1, 'Primeiro dia de aula', 'Alice se adaptou muito bem e fez vários amigos durante as brincadeiras de integração no parquinho.', 'Excelente desenvolvimento social e adaptativo inicial.', '', 2, 1, 'private', '2026-02-03 10:00:00-03', '2026-02-03 10:00:00-03'),
(2, 1, 'Pintura com Guache', 'Hoje fizemos uma atividade de pintura artística usando tintas coloridas! As crianças exploraram texturas com as mãos.', 'Estímulo à coordenação motora fina e expressão visual.', 'https://images.unsplash.com/photo-1513364776144-60967b0f800f?w=600', 5, 2, 'turma', '2026-06-16 14:00:00-03', '2026-06-16 14:00:00-03'),
(3, 3, 'Teatro de Dedoches', 'Davi se divertiu muito participando do teatro de dedoches e interpretando o personagem do jacaré.', 'Desenvolvimento de fala, linguagem e narrativa lúdica.', 'https://images.unsplash.com/photo-1560421683-6856ea585c78?w=600', 3, 0, 'turma', '2026-06-17 09:00:00-03', '2026-06-17 09:00:00-03'),
(4, 2, 'Atividade de Blocos', 'Bernardo montou uma torre enorme usando blocos de madeira e testou as leis do equilíbrio.', 'Raciocínio espacial e equilíbrio geométrico.', '', 4, 1, 'private', '2026-06-17 11:00:00-03', '2026-06-17 11:00:00-03')
ON CONFLICT (id) DO UPDATE SET 
    student_id = EXCLUDED.student_id, title = EXCLUDED.title, description = EXCLUDED.description, 
    pedagogical_note = EXCLUDED.pedagogical_note, image_url = EXCLUDED.image_url, likes = EXCLUDED.likes, 
    bookmarks = EXCLUDED.bookmarks, visibility = EXCLUDED.visibility, created_at = EXCLUDED.created_at, updated_at = EXCLUDED.updated_at;
SELECT setval('posts_id_seq', COALESCE((SELECT MAX(id) FROM posts), 1));

-- 6. DYNAMIC GENERATION OF 60 STUDENTS & GUARDIANS (30 in Turma 1, 20 in Turma 2, 10 in Turma 3)
DO $$
DECLARE
    first_names text[] := ARRAY['Gabriel', 'Lucas', 'Matheus', 'Pedro', 'Enzo', 'João', 'Guilherme', 'Nicolas', 'Rafael', 'Thiago', 'Gustavo', 'Felipe', 'Arthur', 'Samuel', 'Daniel', 'Leonardo', 'Bruno', 'Eduardo', 'Henrique', 'Murilo', 'Sophia', 'Alice', 'Júlia', 'Isabella', 'Manuela', 'Laura', 'Luiza', 'Valentina', 'Giovanna', 'Maria', 'Beatriz', 'Heloísa', 'Lívia', 'Lara', 'Mariana', 'Yasmin', 'Gabriela', 'Eduarda', 'Lorena', 'Rafaela'];
    last_names text[] := ARRAY['Silva', 'Santos', 'Oliveira', 'Souza', 'Rodrigues', 'Ferreira', 'Alves', 'Pereira', 'Lima', 'Gomes', 'Costa', 'Ribeiro', 'Martins', 'Carvalho', 'Almeida', 'Lopes', 'Soares', 'Fernandes', 'Vieira', 'Barbosa'];
    parent_first_names text[] := ARRAY['Carlos', 'Marcos', 'Roberto', 'Ricardo', 'Fernando', 'Rodrigo', 'Felipe', 'André', 'Júlio', 'Adriano', 'Fernanda', 'Juliana', 'Camila', 'Aline', 'Patrícia', 'Bruna', 'Amanda', 'Renata', 'Larissa', 'Mariana'];
    
    avatar_urls text[] := ARRAY[
        'https://images.unsplash.com/photo-1503919545889-aef636e10ad4?w=200',
        'https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?w=200',
        'https://images.unsplash.com/photo-1519238263530-99bdd11df2ea?w=200',
        'https://images.unsplash.com/photo-1471286174890-9c112ffca5b4?w=200',
        'https://images.unsplash.com/photo-1484186139897-d5fc6b908812?w=200',
        'https://images.unsplash.com/photo-1513551573338-107695245848?w=200',
        'https://images.unsplash.com/photo-1516627145497-ae6968895b74?w=200',
        'https://images.unsplash.com/photo-1518887570146-0612132dd618?w=200',
        'https://images.unsplash.com/photo-1502082553048-f009c37129b9?w=200',
        'https://images.unsplash.com/photo-1485546246426-74dc88dec4d9?w=200'
    ];
    
    parent_avatar_urls text[] := ARRAY[
        'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=150',
        'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150',
        'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=150',
        'https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=150',
        'https://images.unsplash.com/photo-1524504388940-b1c1722653e1?w=150',
        'https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=150'
    ];

    stud_idx int := 100;
    parent_user_idx int := 100;
    guardian_idx int := 100;
    
    turma_id int;
    i int;
    lim int;
    
    p_first text;
    p_last text;
    p_fullname text;
    p_email text;
    
    s_first text;
    s_last text;
    s_fullname text;
    
    g_relation text;
    teacher_id int;
    teacher_name text;
    group_name text;
BEGIN
    FOR turma_id IN 1..3 LOOP
        IF turma_id = 1 THEN
            lim := 30;
            teacher_id := 1;
            teacher_name := 'Ana Souza (Professora)';
            group_name := 'Turma das Borboletas (3-4 anos)';
        ELSIF turma_id = 2 THEN
            lim := 20;
            teacher_id := 2;
            teacher_name := 'Carlos Lima (Psicopedagogo)';
            group_name := 'Turma dos Leões (4-5 anos)';
        ELSE
            lim := 10;
            teacher_id := 1;
            teacher_name := 'Ana Souza (Professora)';
            group_name := 'Turma dos Golfinhos (5-6 anos)';
        END IF;

        FOR i IN 1..lim LOOP
            stud_idx := stud_idx + 1;
            parent_user_idx := parent_user_idx + 1;
            guardian_idx := guardian_idx + 1;

            p_first := parent_first_names[1 + ((parent_user_idx * 7) % array_length(parent_first_names, 1))];
            p_last := last_names[1 + ((parent_user_idx * 11) % array_length(last_names, 1))];
            p_fullname := p_first || ' ' || p_last;
            p_email := lower(p_first) || '.' || lower(p_last) || parent_user_idx || '@gmail.com';

            INSERT INTO users (id, name, email, password_hash, role, phone, address, avatar_url, created_at, updated_at)
            VALUES (parent_user_idx, p_fullname, p_email, '$2a$10$eDDuus0hIwc38PanTSeAxu.v7cIGIYlWx4wNG.CgQlfvqczSntOdK', 'responsavel', '(11) 9' || (10000 + parent_user_idx) || '-8888', 'Av. das Nacoes, ' || parent_user_idx, parent_avatar_urls[1 + (parent_user_idx % array_length(parent_avatar_urls, 1))], now(), now())
            ON CONFLICT (id) DO NOTHING;

            s_first := first_names[1 + ((stud_idx * 3) % array_length(first_names, 1))];
            s_last := p_last;
            s_fullname := s_first || ' ' || s_last;

            INSERT INTO students (id, name, presence_status, check_in_at, guardian_user_id, photo_url, turma_id, group_name, teacher_user_id, teacher_name, birth_date, enrollment_code, blood_type, allergies, restrictions, medications, created_at, updated_at)
            VALUES (stud_idx, s_fullname, 'present', now(), parent_user_idx, avatar_urls[1 + (stud_idx % array_length(avatar_urls, 1))], turma_id, group_name, teacher_id, teacher_name, '2021-01-01'::date + (stud_idx * 3), '#2026-' || stud_idx, 'A+', ARRAY[]::TEXT[], 'Nenhuma', 'Nenhum', now(), now())
            ON CONFLICT (id) DO NOTHING;

            IF parent_user_idx % 2 = 0 THEN
                g_relation := 'Avó';
            ELSE
                g_relation := 'Tio';
            END IF;

            INSERT INTO student_guardians (id, student_id, name, relation, phone, avatar_url, authorized, created_at)
            VALUES (guardian_idx, stud_idx, 'Parente de ' || s_first, g_relation, '(11) 9' || (20000 + guardian_idx) || '-9999', parent_avatar_urls[1 + ((guardian_idx * 2) % array_length(parent_avatar_urls, 1))], true, now())
            ON CONFLICT (id) DO NOTHING;
        END LOOP;
    END LOOP;
END $$;

SELECT setval('users_id_seq', COALESCE((SELECT MAX(id) FROM users), 1));
SELECT setval('students_id_seq', COALESCE((SELECT MAX(id) FROM students), 1));
SELECT setval('student_guardians_id_seq', COALESCE((SELECT MAX(id) FROM student_guardians), 1));

-- 7. DYNAMIC GENERATION OF 70+ HISTORICAL POSTS (Activities from 2026-02-02 to 2026-06-15)
DO $$
DECLARE
    titles text[] := ARRAY[
        'Atividade de Pintura Coletiva',
        'Oficina de Jardinagem e Plantio',
        'Circuito de Obstáculos Psicomotores',
        'Roda de Leitura de Histórias',
        'Atividade de Reciclagem Criativa',
        'Teatro de Fantoches e Expressão',
        'Brincadeiras na Caixa de Areia',
        'Desenho Livre com Giz de Cera',
        'Modelagem com Argila Natural',
        'Oficina de Culinária Saudável',
        'Jogos Cooperativos de Tabuleiro',
        'Brincadeiras Rítmicas e Musicais',
        'Caça ao Tesouro Pedagógica',
        'Dobraduras de Papel e Origami',
        'Exploração Sensorial com Tintas'
    ];
    
    descriptions text[] := ARRAY[
        'Hoje as crianças se reuniram para pintar um grande mural usando tintas guache coloridas e pincéis. O trabalho em equipe foi lindo!',
        'No jardim da escola, preparamos a terra e plantamos pequenas sementes de flores. Aprendemos sobre o crescimento das plantas.',
        'Desenvolvemos um circuito físico no pátio com túneis, cones e cordas para exercitar o equilíbrio e agilidade de todos.',
        'Sentamos em círculo embaixo da árvore para ler uma bela fábula infantil. O interesse e engajamento foi total.',
        'Reutilizamos garrafas plásticas e rolos de papel para criar divertidos brinquedos e carrinhos ecológicos.',
        'A turma encenou uma pequena peça com fantoches de feltro, trabalhando a fala e a imaginação de forma lúdica.',
        'Construímos castelos e pontes gigantes na caixa de areia do parque, testando o equilíbrio e cooperação espacial.',
        'Utilizamos papel pardo no chão para criar desenhos gigantes com giz colorido, expressando muita criatividade.',
        'Exploramos a maleabilidade da argila molhada hoje, modelando pequenos animais e objetos da natureza.',
        'Preparamos uma deliciosa salada de frutas com banana, maçã e morango. As crianças ajudaram na higienização.',
        'Jogamos jogos de tabuleiro cooperativos simples, aprendendo sobre regras, paciência e a respeitar a vez dos amigos.',
        'Brincamos de bater palmas e acompanhar o ritmo com chocalhos artesanais que confeccionamos em sala.',
        'Seguimos pistas espalhadas pelo jardim para encontrar o baú de livros dourado, trabalhando raciocínio e orientação.',
        'Fizemos dobraduras simples de aviões, barcos e animais de papel para exercitar a coordenação fina de precisão.',
        'Experimentamos misturar tintas coloridas com creme corporal para sentir texturas e criar novas cores em papel cartão.'
    ];

    ped_notes text[] := ARRAY[
        'Promove a socialização, colaboração e a expressão livre de cores.',
        'Estimula a educação ambiental, paciência e a percepção do ciclo de vida das plantas.',
        'Exercita o equilíbrio, orientação espacial e coordenação motora grossa.',
        'Desenvolve a escuta ativa, vocabulário e o gosto pela literatura infantil.',
        'Desperta a consciência ecológica e estimula a criatividade criadora.',
        'Desenvolve a oralidade, dramatização e expressão das emoções.',
        'Estimula a coordenação visomotora, textura e a diversão compartilhada no brincar.',
        'Fortalece os músculos das mãos e braços, auxiliando no desenho livre.',
        'Desenvolve a motricidade fina das mãos e a percepção tridimensional.',
        'Desenvolve a autonomia prática e noções básicas de alimentação saudável.',
        'Ensina a convivência de grupo, cooperação e habilidades de tolerância.',
        'Desenvolve a audição ativa, percepção rítmica e sintonia corporal.',
        'Estimula a resolução de problemas lógicos e a coordenação motora global.',
        'Exercita a precisão dos movimentos, paciência e a percepção de formas.',
        'Trabalha a estimulação sensorial tátil e a criatividade através das misturas.'
    ];

    image_urls text[] := ARRAY[
        'https://images.unsplash.com/photo-1513364776144-60967b0f800f?w=600',
        'https://images.unsplash.com/photo-1466692476868-aef1dfb1e735?w=600',
        'https://images.unsplash.com/photo-1502086223501-7ea6ecd79368?w=600',
        'https://images.unsplash.com/photo-1506880018603-83d5b814b5a6?w=600',
        'https://images.unsplash.com/photo-1587654780291-39c9404d746b?w=600',
        'https://images.unsplash.com/photo-1560421683-6856ea585c78?w=600',
        'https://images.unsplash.com/photo-1502082553048-f009c37129b9?w=600',
        'https://images.unsplash.com/photo-1516627145497-ae6968895b74?w=600',
        'https://images.unsplash.com/photo-1518887570146-0612132dd618?w=200',
        'https://images.unsplash.com/photo-1498837167922-ddd27525d352?w=600',
        'https://images.unsplash.com/photo-1530273673054-127c427fd744?w=600',
        'https://images.unsplash.com/photo-1511671782779-c97d3d27a1d4?w=600',
        'https://images.unsplash.com/photo-1473448912268-2022ce9509d8?w=600',
        'https://images.unsplash.com/photo-1544816155-12df9643f363?w=600',
        'https://images.unsplash.com/photo-1536640788329-e21c8d41cfc8?w=600'
    ];

    post_idx int := 100;
    curr_date DATE;
    s_id int;
    template_idx int;
    vis text;
    post_likes int;
    post_bookmarks int;
BEGIN
    FOR curr_date IN SELECT generate_series('2026-02-02'::date, '2026-06-15'::date, '2 days'::interval)::date LOOP
        post_idx := post_idx + 1;
        template_idx := 1 + (post_idx % array_length(titles, 1));
        
        -- Pick a student id from generated students (101..160)
        s_id := 101 + (post_idx % 60);

        IF post_idx % 3 = 0 THEN
            vis := 'private';
        ELSE
            vis := 'turma';
        END IF;

        post_likes := (post_idx * 3) % 15;
        post_bookmarks := (post_idx * 2) % 6;

        INSERT INTO posts (id, student_id, title, description, pedagogical_note, image_url, likes, bookmarks, visibility, created_at, updated_at)
        VALUES (post_idx, s_id, titles[template_idx], descriptions[template_idx], ped_notes[template_idx], image_urls[template_idx], post_likes, post_bookmarks, vis, curr_date + time '14:00:00', curr_date + time '14:00:00')
        ON CONFLICT (id) DO NOTHING;
    END LOOP;
END $$;

SELECT setval('posts_id_seq', COALESCE((SELECT MAX(id) FROM posts), 1));

-- 8. DYNAMIC GENERATION OF COMMENTS ON THE HISTORICAL POSTS
DO $$
DECLARE
    comment_texts text[] := ARRAY[
        'Que atividade maravilhosa! Meu filho adorou participar disso.',
        'Muito bom ver o progresso da turma nesta dinâmica. Parabéns professora!',
        'Excelente trabalho pedagógico. Obrigado pelo cuidado e carinho de sempre.',
        'Lindo de ver! Ele chegou em casa contando todos os detalhes super empolgado.',
        'Que capricho! Adoro acompanhar as fotinhas por aqui.',
        'Muito importante esse estímulo na infância. Obrigado equipe!',
        'Ele adorou! Disse que quer fazer de novo amanhã kkkk.',
        'Que iniciativa fantástica da escola. Parabéns a todos os envolvidos!'
    ];
    
    c_idx int := 100;
    p_id int;
    s_id int;
    g_id int;
    g_name text;
    g_avatar text;
    comm_text text;
    p_date timestamp;
BEGIN
    -- Loop over all generated posts (ids 101 to 167)
    FOR p_id IN 101..167 LOOP
        IF EXISTS (SELECT 1 FROM posts WHERE id = p_id) THEN
            SELECT student_id, created_at INTO s_id, p_date FROM posts WHERE id = p_id;
            
            SELECT guardian_user_id INTO g_id FROM students WHERE id = s_id;
            SELECT name, avatar_url INTO g_name, g_avatar FROM users WHERE id = g_id;

            c_idx := c_idx + 1;
            comm_text := comment_texts[1 + (c_idx % array_length(comment_texts, 1))];

            INSERT INTO comments (id, post_id, user_id, author_name, avatar_url, text, created_at)
            VALUES (c_idx, p_id, g_id, g_name, g_avatar, comm_text, p_date + interval '4 hours')
            ON CONFLICT (id) DO NOTHING;

            IF c_idx % 2 = 0 THEN
                c_idx := c_idx + 1;
                INSERT INTO comments (id, post_id, user_id, author_name, avatar_url, text, created_at)
                VALUES (c_idx, p_id, 1, 'Ana Souza (Professora)', 'https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=150', 'Obrigado pelo retorno! Ficamos muito felizes com o desenvolvimento deles.', p_date + interval '5 hours')
                ON CONFLICT (id) DO NOTHING;
            END IF;
        END IF;
    END LOOP;
END $$;

SELECT setval('comments_id_seq', COALESCE((SELECT MAX(id) FROM comments), 1));

-- 9. Timeline Events (Manual / static ones)
INSERT INTO timeline_events (id, student_id, title, description, occurred_at, created_at) VALUES
(1, 1, 'Chegada na Escola', 'Alice entrou animada na sala de aula e cumprimentou os colegas com abraços.', '2026-06-16 08:00:00-03', '2026-06-16 08:00:00-03'),
(2, 1, 'Hora do Lanche', 'Comeu toda a porção de frutas (banana e maçã) oferecida hoje no lanche coletivo.', '2026-06-16 10:15:00-03', '2026-06-16 10:15:00-03'),
(3, 1, 'Atividade no Jardim', 'Brincou com terra e plantou uma sementinha de girassol no vaso reciclado.', '2026-06-16 14:30:00-03', '2026-06-16 14:30:00-03'),
(4, 3, 'Chegada na Escola', 'Davi chegou com um pouco de sono, mas logo correu para a mesa de blocos.', '2026-06-16 08:10:00-03', '2026-06-16 08:10:00-03'),
(5, 3, 'Hora do Lanche', 'Comeu a salada de frutas sem nenhuma restrição.', '2026-06-16 10:20:00-03', '2026-06-16 10:20:00-03'),
(6, 1, 'Chegada na Escola', 'Alice chegou rindo muito, correu diretamente para abraçar a professora Ana.', '2026-06-17 07:55:00-03', '2026-06-17 07:55:00-03'),
(7, 1, 'Colação da Manhã', 'Tomou todo o suco de uva integral e comeu biscoito de polvilho.', '2026-06-17 09:30:00-03', '2026-06-17 09:30:00-03'),
(8, 1, 'Brincadeira Livre', 'Alice montou quebra-cabeça de animais de madeira junto com o Gabriel.', '2026-06-17 11:15:00-03', '2026-06-17 11:15:00-03'),
(9, 3, 'Chegada na Escola', 'Davi chegou super agitado mostrando sua blusa de dinossauro.', '2026-06-17 08:05:00-03', '2026-06-17 08:05:00-03'),
(10, 3, 'Colação da Manhã', 'Comeu melão e tomou chá de camomila.', '2026-06-17 09:35:00-03', '2026-06-17 09:35:00-03'),
(11, 2, 'Chegada na Escola', 'Bernardo chegou calmo e ajudou o tio Carlos a organizar a mesa de argila.', '2026-06-17 08:15:00-03', '2026-06-17 08:15:00-03')
ON CONFLICT (id) DO UPDATE SET 
    student_id = EXCLUDED.student_id, title = EXCLUDED.title, description = EXCLUDED.description, 
    occurred_at = EXCLUDED.occurred_at;
SELECT setval('timeline_events_id_seq', COALESCE((SELECT MAX(id) FROM timeline_events), 1));

-- 10. Milestones
INSERT INTO milestones (id, student_id, title, category, description, achieved_at, done, created_at) VALUES
(1, 1, 'Falar frases completas', 'Linguagem', 'Consegue expressar ideias complexas e formular frases com sujeito e verbo estruturados.', '2026-05-10 00:00:00-03', true, '2026-05-10 00:00:00-03'),
(2, 1, 'Compartilhar brinquedos', 'Social', 'Demonstra empatia e divide brinquedos espontaneamente durante as dinâmicas em grupo.', '2026-06-02 00:00:00-03', true, '2026-06-02 00:00:00-03'),
(3, 1, 'Amarrar os sapatos sozinho', 'Motor', 'Em progresso no aprendizado do laço duplo dos calçados escolares.', NULL, false, '2026-06-07 00:00:00-03'),
(4, 3, 'Segurar o giz de cera corretamente', 'Motor', 'Usa a pinça fina para pintar desenhos delimitados de forma coordenada.', '2026-06-12 00:00:00-03', true, '2026-06-12 00:00:00-03')
ON CONFLICT (id) DO UPDATE SET 
    student_id = EXCLUDED.student_id, title = EXCLUDED.title, category = EXCLUDED.category, 
    description = EXCLUDED.description, achieved_at = EXCLUDED.achieved_at, done = EXCLUDED.done;
SELECT setval('milestones_id_seq', COALESCE((SELECT MAX(id) FROM milestones), 1));

-- 11. Events
INSERT INTO events (id, title, description, location, starts_at, ends_at, rsvp_count, created_at) VALUES
(1, 'Festa Junina do Up Espaço', 'Venha comemorar conosco com muitas brincadeiras, comidas típicas e dança junina!', 'Pátio Central da Escola', now() + interval '8 days', now() + interval '8 days 4 hours', 12, now()),
(2, 'Reunião de Pais e Mestres', 'Conversa individual sobre o desenvolvimento escolar e psicomotor do primeiro semestre.', 'Salas de Aula Individuais', now() + interval '13 days', now() + interval '13 days 4 hours', 4, now()),
(3, 'Abertura do Ano Letivo 2026', 'Boas-vindas a todos os alunos e pais para o início das atividades escolares.', 'Auditório Principal', '2026-02-02 08:00:00-03', '2026-02-02 12:00:00-03', 55, '2026-01-20 09:00:00-03'),
(4, 'Carnaval no Up Espaço', 'Festa de carnaval infantil com fantasias, marchinhas e desfile de blocos.', 'Pátio Recreativo', '2026-02-13 13:30:00-03', '2026-02-13 17:30:00-03', 48, '2026-02-05 10:00:00-03'),
(5, 'Palestra sobre Nutrição Infantil', 'Encontro com nutricionista convidada para debater hábitos alimentares e lanches saudáveis.', 'Sala Multiuso', '2026-03-18 19:00:00-03', '2026-03-18 21:00:00-03', 35, '2026-03-01 08:30:00-03'),
(6, 'Caça aos Ovos de Páscoa', 'Atividade lúdica de caça aos ovos de chocolate e brincadeiras com o Coelhinho.', 'Jardins da Escola', '2026-04-02 09:00:00-03', '2026-04-02 11:30:00-03', 62, '2026-03-20 14:00:00-03'),
(7, 'Oficina do Dia das Mães', 'Oficina especial para as mães e filhos confeccionarem sabonetes artesanais juntos.', 'Oficina de Artes', '2026-05-08 14:00:00-03', '2026-05-08 16:30:00-03', 50, '2026-04-25 09:00:00-03'),
(8, 'Semana do Meio Ambiente: Plantio', 'Atividade coletiva de plantio de mudas frutíferas e palestra sobre reciclagem.', 'Horta e Bosque da Escola', '2026-06-05 09:30:00-03', '2026-06-05 12:00:00-03', 42, '2026-05-25 11:00:00-03')
ON CONFLICT (id) DO UPDATE SET 
    title = EXCLUDED.title, description = EXCLUDED.description, location = EXCLUDED.location, 
    starts_at = EXCLUDED.starts_at, ends_at = EXCLUDED.ends_at, rsvp_count = EXCLUDED.rsvp_count;
SELECT setval('events_id_seq', COALESCE((SELECT MAX(id) FROM events), 1));

-- 12. Announcements
INSERT INTO announcements (id, title, sender, priority, preview, body, attachment_name, created_at) VALUES
(1, 'Campanha de Vacinação Infantil', 'Direção Up Espaço', 'Importante', 'Vacinação contra a gripe acontecerá na próxima quarta-feira.', 'Solicitamos aos pais que tragam a caderneta de vacinação da criança no dia da campanha. A equipe de saúde da prefeitura estará presente a partir das 9h no pátio.', 'manual_vacinacao.pdf', now() - interval '1 days'),
(2, 'Aviso de Feriado', 'Secretaria', 'Informativo', 'Não haverá aula no dia 09 de Julho devido ao feriado paulista.', 'Comunicamos que a escola estará fechada na próxima terça-feira (09/07). Retornaremos as atividades normais na quarta-feira (10/07). Bom descanso a todos.', NULL, now() - interval '2 hours')
ON CONFLICT (id) DO UPDATE SET 
    title = EXCLUDED.title, sender = EXCLUDED.sender, priority = EXCLUDED.priority, 
    preview = EXCLUDED.preview, body = EXCLUDED.body, attachment_name = EXCLUDED.attachment_name;
SELECT setval('announcements_id_seq', COALESCE((SELECT MAX(id) FROM announcements), 1));

-- 13. DYNAMIC GENERATION OF 4 MONTHS OF ATTENDANCE (Weekday records for all students)
DO $$
DECLARE
    curr_date DATE;
    stud_id INT;
    stat TEXT;
    marked_by INT;
BEGIN
    FOR curr_date IN SELECT generate_series('2026-02-02'::date, '2026-06-17'::date, '1 day'::interval)::date LOOP
        -- Skip weekends (Saturday = 6, Sunday = 0)
        IF extract(dow from curr_date) IN (0, 6) THEN
            CONTINUE;
        END IF;

        FOR stud_id IN SELECT id FROM students LOOP
            -- 90% chance present, 10% chance absent
            IF random() < 0.10 THEN
                stat := 'absent';
            ELSE
                stat := 'present';
            END IF;

            SELECT teacher_user_id INTO marked_by FROM students WHERE id = stud_id;
            IF marked_by IS NULL THEN
                marked_by := 1;
            END IF;

            INSERT INTO attendance (student_id, date, status, marked_by_user_id, created_at)
            VALUES (stud_id, curr_date, stat, marked_by, curr_date + time '08:00:00' + (random() * interval '30 minutes'))
            ON CONFLICT (student_id, date) DO UPDATE SET status = EXCLUDED.status, marked_by_user_id = EXCLUDED.marked_by_user_id;
        END LOOP;
    END LOOP;
END $$;

SELECT setval('attendance_id_seq', COALESCE((SELECT MAX(id) FROM attendance), 1));
