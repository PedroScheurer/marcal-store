-- Senha do admin: password
-- Hash BCrypt (strength 10), compatível com BCryptPasswordEncoder do Spring Security.
UPDATE tb_user
SET password = '$2a$10$EixZaYVK1fsbw1ZfbX3OXePaWxn96p36WQoeG6Lruj3vjPDwqM5.'
WHERE email = 'admin@admin.dev';
