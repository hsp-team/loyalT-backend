-- Create Users
INSERT INTO users (id, name, email, password, created_at)
VALUES ('123e4567-e89b-12d3-a456-426614174000', 'Иван Петров', 'ivan@example.com',
        '$2a$10$q1oDgYM7AcqWshA3kmKLkuqHANh12lGWt0BBHBaVBSOKcmLmzebnO', '2025-03-03 14:04:34'),
       ('123e4567-e89b-12d3-a456-426614174001', 'Мария Сидорова', 'maria@example.com',
        '$2a$10$q1oDgYM7AcqWshA3kmKLkuqHANh12lGWt0BBHBaVBSOKcmLmzebnO', '2025-03-03 14:04:34'),
       ('123e4567-e89b-12d3-a456-426614174002', 'Алексей Иванов', 'alex@example.com',
        '$2a$10$q1oDgYM7AcqWshA3kmKLkuqHANh12lGWt0BBHBaVBSOKcmLmzebnO', '2025-03-03 14:04:34');

-- Create Businesses
INSERT INTO businesses (id, name, email, password, description, created_at)
VALUES ('123e4567-e89b-12d3-a456-426614174003', 'Кофейня "Утро"', 'coffee@morning.com',
        '$2a$10$q1oDgYM7AcqWshA3kmKLkuqHANh12lGWt0BBHBaVBSOKcmLmzebnO', 'Уютная кофейня в центре города',
        '2025-03-03 14:04:34'),
       ('123e4567-e89b-12d3-a456-426614174004', 'Пиццерия "Вкусно"', 'pizza@tasty.com',
        '$2a$10$q1oDgYM7AcqWshA3kmKLkuqHANh12lGWt0BBHBaVBSOKcmLmzebnO', 'Лучшая пицца в городе', '2025-03-03 14:04:34');

-- Create Coin Programs
INSERT INTO coin_programs (id, name, description, day_limit, card_color, business_coin_program, created_at)
VALUES ('123e4567-e89b-12d3-a456-426614174005', 'Утренние бонусы', 'Программа лояльности кофейни', 5, '#FFD700',
        '123e4567-e89b-12d3-a456-426614174003', '2025-03-03 14:04:34'),
       ('123e4567-e89b-12d3-a456-426614174006', 'Пицца-клуб', 'Бонусная программа пиццерии', 3, '#FF4500',
        '123e4567-e89b-12d3-a456-426614174004', '2025-03-03 14:04:34');

-- Create Rewards for Coffee Shop
INSERT INTO rewards (id, name, description, cost, image_url, coin_program_rewards, created_at)
VALUES ('123e4567-e89b-12d3-a456-426614174007', 'Капучино бесплатно', 'Капучино на выбор', 2,
        'https://images.unsplash.com/photo-1572442388796-11668a67e53d?w=800', '123e4567-e89b-12d3-a456-426614174005',
        '2025-03-03 14:04:34'),
       ('123e4567-e89b-12d3-a456-426614174008', 'Круассан в подарок', 'Свежий круассан', 1,
        'https://images.unsplash.com/photo-1555507036-ab1f4038808a?w=800', '123e4567-e89b-12d3-a456-426614174005',
        '2025-03-03 14:04:34'),
       ('123e4567-e89b-12d3-a456-426614174009', 'Кофе с собой', 'Любой кофе с собой', 2,
        'https://images.unsplash.com/photo-1497935586351-b67a49e012bf?w=800', '123e4567-e89b-12d3-a456-426614174005',
        '2025-03-03 14:04:34');

-- Create Rewards for Pizza Place
INSERT INTO rewards (id, name, description, cost, image_url, coin_program_rewards, created_at)
VALUES ('123e4567-e89b-12d3-a456-426614174010', 'Пицца 30см', 'Любая пицца 30см', 2,
        'https://images.unsplash.com/photo-1513104890138-7c749659a591?w=800', '123e4567-e89b-12d3-a456-426614174006',
        '2025-03-03 14:04:34'),
       ('123e4567-e89b-12d3-a456-426614174011', 'Напиток 0.5л', 'Любой напиток 0.5л', 1,
        'https://images.unsplash.com/photo-1581006852262-e4307cf6283a?w=800', '123e4567-e89b-12d3-a456-426614174006',
        '2025-03-03 14:04:34'),
       ('123e4567-e89b-12d3-a456-426614174012', 'Десерт дня', 'Десерт на выбор', 2,
        'https://images.unsplash.com/photo-1551024601-bec78aea704b?w=800', '123e4567-e89b-12d3-a456-426614174006',
        '2025-03-03 14:04:34');

-- Create Coin Program Participants
INSERT INTO coin_program_participants (id, balance, user_coin_programs, coin_program_participant_coin_program,
                                       created_at)
VALUES ('123e4567-e89b-12d3-a456-426614174013', 100, '123e4567-e89b-12d3-a456-426614174000',
        '123e4567-e89b-12d3-a456-426614174005', '2025-03-03 14:04:34'),
       ('123e4567-e89b-12d3-a456-426614174014', 50, '123e4567-e89b-12d3-a456-426614174001',
        '123e4567-e89b-12d3-a456-426614174005', '2025-03-03 14:04:34'),
       ('123e4567-e89b-12d3-a456-426614174015', 75, '123e4567-e89b-12d3-a456-426614174002',
        '123e4567-e89b-12d3-a456-426614174006', '2025-03-03 14:04:34');