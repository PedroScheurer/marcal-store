-- Corrige hash bcrypt inválido (59 chars) do admin após V3.
-- Gera hash compatível com BCryptPasswordEncoder do Spring, igual ao fluxo de signup.
CREATE EXTENSION IF NOT EXISTS pgcrypto;

UPDATE tb_user
SET password = crypt('admin123', gen_salt('bf', 10))
WHERE email = 'admin@admin.dev';
