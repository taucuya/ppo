-- Вставляем пользователей с новыми ролями
-- Вставляем пользователей с хешированными паролями (bcrypt)
INSERT INTO "user" (name, date_of_birth, mail, password, phone, address, status, role) VALUES
-- Обычные пользователи (пароль: user123)
('Анна Смирнова', '1990-05-15', 'anna.smirnova@example.com', '$2a$10$2sOJWDJ/OJielZc2GINfNeyHplhp0r4KjImsXdjZWmLOgRKXbaEDC', '+79161234567', 'ул. Ленина, 10', 'active', 'обычный пользователь'),
('Елена Кузнецова', '1985-08-20', 'elena.kuznetsova@example.com', '$2a$10$2sOJWDJ/OJielZc2GINfNeyHplhp0r4KjImsXdjZWmLOgRKXbaEDC', '+79162345678', 'ул. Пушкина, 5', 'active', 'обычный пользователь'),
('Мария Иванова', '1995-02-10', 'maria.ivanova@example.com', '$2a$10$2sOJWDJ/OJielZc2GINfNeyHplhp0r4KjImsXdjZWmLOgRKXbaEDC', '+79163456789', 'пр. Мира, 15', 'active', 'обычный пользователь'),
('Ольга Петрова', '1988-11-30', 'olga.petrova@example.com', '$2a$10$2sOJWDJ/OJielZc2GINfNeyHplhp0r4KjImsXdjZWmLOgRKXbaEDC', '+79164567890', 'ул. Гагарина, 3', 'active', 'обычный пользователь'),
('Ирина Соколова', '1992-07-25', 'irina.sokolova@example.com', '$2a$10$2sOJWDJ/OJielZc2GINfNeyHplhp0r4KjImsXdjZWmLOgRKXbaEDC', '+79165678901', 'ул. Кирова, 7', 'active', 'обычный пользователь'),
('Наталья Волкова', '1993-04-18', 'natalia.volkova@example.com', '$2a$10$2sOJWDJ/OJielZc2GINfNeyHplhp0r4KjImsXdjZWmLOgRKXbaEDC', '+79166789012', 'ул. Советская, 22', 'active', 'обычный пользователь'),
('Виктория Федорова', '1987-09-05', 'victoria.fedorova@example.com', '$2a$10$2sOJWDJ/OJielZc2GINfNeyHplhp0r4KjImsXdjZWmLOgRKXbaEDC', '+79167890123', 'ул. Садовая, 14', 'active', 'обычный пользователь'),
('Юлия Морозова', '1991-12-12', 'julia.morozova@example.com', '$2a$10$2sOJWDJ/OJielZc2GINfNeyHplhp0r4KjImsXdjZWmLOgRKXbaEDC', '+79168901234', 'ул. Лесная, 8', 'active', 'обычный пользователь'),
('Александра Николаева', '1989-06-28', 'alexandra.nikolaeva@example.com', '$2a$10$2sOJWDJ/OJielZc2GINfNeyHplhp0r4KjImsXdjZWmLOgRKXbaEDC', '+79169012345', 'пр. Победы, 17', 'active', 'обычный пользователь'),
('Екатерина Павлова', '1994-03-08', 'ekaterina.pavlova@example.com', '$2a$10$2sOJWDJ/OJielZc2GINfNeyHplhp0r4KjImsXdjZWmLOgRKXbaEDC', '+79160123456', 'ул. Цветочная, 9', 'active', 'обычный пользователь'),

-- Администратор (пароль: admin123)
('Алексей Администратор', '1980-01-01', 'admin@cosmetics.ru', '$2a$10$2sOJWDJ/OJielZc2GINfNeyHplhp0r4KjImsXdjZWmLOgRKXbaEDC', '+79150000001', 'ул. Центральная, 1', 'active', 'админ'),

-- Работники склада (пароль: worker123)
('Иван Складской', '1985-03-15', 'worker1@cosmetics.ru', '$2a$10$2sOJWDJ/OJielZc2GINfNeyHplhp0r4KjImsXdjZWmLOgRKXbaEDC', '+79150000002', 'ул. Складская, 5', 'active', 'работник склада'),
('Петр Складской', '1988-07-20', 'worker2@cosmetics.ru', '$2a$10$2sOJWDJ/OJielZc2GINfNeyHplhp0r4KjImsXdjZWmLOgRKXbaEDC', '+79150000003', 'ул. Складская, 10', 'active', 'работник склада'),
('Сергей Складской', '1990-11-30', 'worker3@cosmetics.ru', '$2a$10$2sOJWDJ/OJielZc2GINfNeyHplhp0r4KjImsXdjZWmLOgRKXbaEDC', '+79150000004', 'ул. Складская, 15', 'active', 'работник склада');

-- Вставляем бренды с новыми категориями
INSERT INTO brand (name, description, price_category) VALUES
('L''Oreal', 'Французский косметический бренд премиум-класса', 'люкс'),
('Estée Lauder', 'Американский люксовый бренд косметики', 'люкс'),
('Maybelline', 'Бюджетный бренд декоративной косметики', 'бюджет'),
('Nivea', 'Немецкий бренд по уходу за кожей', 'средний'),
('Garnier', 'Французский бренд косметики для ухода', 'средний'),
('Clinique', 'Американский бренд гипоаллергенной косметики', 'люкс'),
('Revlon', 'Международный бренд декоративной косметики', 'средний'),
('The Ordinary', 'Канадский бренд уходовой косметики', 'средний'),
('La Roche-Posay', 'Французская аптечная косметика', 'люкс'),
('NYX', 'Бренд профессиональной декоративной косметики', 'средний');

INSERT INTO product (name, description, price, category, amount, id_brand, pic_link, art) VALUES
('Тональный крем True Match', 'Тональный крем с SPF 30, 30 мл', 2500.00, 'декоративная', 15, (SELECT id FROM brand WHERE name = 'L''Oreal'), '/images/foundation.jpg', 'LOR-TM-001'),
('Тушь для ресниц Lash Sensational', 'Объемная тушь для ресниц', 1200.00, 'декоративная', 50, (SELECT id FROM brand WHERE name = 'Maybelline'), '/images/mascara.jpg', 'MAY-LS-002'),
('Увлажняющий крем Soft', 'Крем для лица с гиалуроновой кислотой, 50 мл', 1800.00, 'уход', 30, (SELECT id FROM brand WHERE name = 'Nivea'), '/images/moisturizer.jpg', 'NIV-SF-003'),
('Помада Super Lustrous', 'Стойкая матовая помада, 4.5 г', 1500.00, 'декоративная', 25, (SELECT id FROM brand WHERE name = 'Revlon'), '/images/lipstick.jpg', 'REV-SL-004'),
('Сыворотка Vitamin C', 'Сыворотка с витамином С, 30 мл', 3500.00, 'уход', 20, (SELECT id FROM brand WHERE name = 'The Ordinary'), '/images/serum.jpg', 'ORD-VC-005'),
('Очищающий гель Pure Active', 'Гель для умывания для проблемной кожи', 900.00, 'уход', 100, (SELECT id FROM brand WHERE name = 'Garnier'), '/images/cleanser.jpg', 'GAR-PA-006'),
('Тени для век Ultimate', 'Палетка теней, 12 оттенков', 2800.00, 'декоративная', 35, (SELECT id FROM brand WHERE name = 'NYX'), '/images/eyeshadow.jpg', 'NYX-UL-007'),
('Солнцезащитный крем Anthelios', 'Крем с SPF 50, 50 мл', 2200.00, 'уход', 18, (SELECT id FROM brand WHERE name = 'La Roche-Posay'), '/images/sunscreen.jpg', 'LRP-AN-008'),
('Духи Beautiful', 'Цветочный аромат, 50 мл', 6500.00, 'парфюмерия', 40, (SELECT id FROM brand WHERE name = 'Estée Lauder'), '/images/perfume.jpg', 'EST-BF-009'),
('Крем для рук Deep Comfort', 'Интенсивный уход за сухой кожей рук', 800.00, 'уход', 22, (SELECT id FROM brand WHERE name = 'Clinique'), '/images/handcream.jpg', 'CLI-DC-010');


-- Вставляем 4 работника (1 админ и 3 работника склада)
INSERT INTO worker (id_user, job_title)
SELECT id, 
  CASE 
    WHEN mail = 'admin@cosmetics.ru' THEN 'admin'
    ELSE 'worker'
  END
FROM "user" 
WHERE role IN ('админ', 'работник склада');

-- Вставляем заказы с новыми статусами
WITH ranked_users AS (
  SELECT id, address, row_number() OVER () as rn FROM "user" WHERE role = 'обычный пользователь' LIMIT 10
)
INSERT INTO "order" (date, id_user, address, status, price)
SELECT 
  now() - (random() * 90 || ' days')::interval,
  u.id,
  u.address,
  CASE (u.rn % 5)
    WHEN 0 THEN 'некорректный'
    WHEN 1 THEN 'непринятый'
    WHEN 2 THEN 'принятый'
    WHEN 3 THEN 'собранный'
    WHEN 4 THEN 'отданный'
  END,
  (random() * 10000 + 1000)::numeric(10,2)
FROM ranked_users u;

WITH target_orders AS (
  SELECT id FROM "order" WHERE status IN ('принятый', 'собранный', 'отданный')
),
target_workers AS (
  SELECT id, row_number() OVER () as rn FROM worker WHERE job_title = 'работник склада'
),
numbered_orders AS (
  SELECT id, row_number() OVER () as rn FROM target_orders
)

INSERT INTO order_worker (id_order, id_worker)
SELECT 
  o.id,
  w.id
FROM numbered_orders o
JOIN target_workers w ON (o.rn % (SELECT COUNT(*) FROM target_workers)) + 1 = w.rn;

INSERT INTO basket (id_user, date) 
SELECT id, now() - (random() * 30 || ' days')::interval 
FROM "user" LIMIT 10;

WITH basket_list AS (
  SELECT id, row_number() OVER () as rn FROM basket LIMIT 10
),
product_list AS (
  SELECT id, row_number() OVER () as rn FROM product LIMIT 10
),
series AS (
  SELECT generate_series(1, 10) as i
)
INSERT INTO basket_item (id_basket, id_product, amount)
SELECT 
  b.id,
  p.id,
  (random() * 5 + 1)::int
FROM series s
JOIN basket_list b ON b.rn = s.i
JOIN product_list p ON p.rn = s.i;

-- Вставляем элементы заказа
WITH order_list AS (
  SELECT id, row_number() OVER () as rn FROM "order" LIMIT 10
),
product_list AS (
  SELECT id, row_number() OVER () as rn FROM product LIMIT 10
),
series AS (
  SELECT generate_series(1, 10) as i
)
INSERT INTO order_item (id_order, id_product, amount)
SELECT 
  o.id,
  p.id,
  (random() * 3 + 1)::int
FROM series s
JOIN order_list o ON o.rn = s.i
JOIN product_list p ON p.rn = s.i;

-- Вставляем отзывы с валидными рейтингами
WITH product_list AS (
  SELECT id, row_number() OVER () as rn FROM product LIMIT 10
),
user_list AS (
  SELECT id, row_number() OVER () as rn FROM "user" LIMIT 10
),
series AS (
  SELECT generate_series(1, 10) as i
)
INSERT INTO review (id_product, id_user, rating, r_text, date)
SELECT 
  p.id,
  u.id,
  (random() * 4 + 1)::int,
  CASE (s.i % 5)
    WHEN 0 THEN 'Отличный продукт, всем рекомендую!'
    WHEN 1 THEN 'Хорошее качество за свои деньги'
    WHEN 2 THEN 'Не совсем то, что ожидала'
    WHEN 3 THEN 'Пока не поняла, как правильно использовать'
    WHEN 4 THEN 'Лучшее, что я покупала за последнее время'
  END,
  now() - (random() * 180 || ' days')::interval
FROM series s
JOIN product_list p ON p.rn = s.i
JOIN user_list u ON u.rn = s.i;

-- Вставляем избранные списки для первых 5 пользователей
INSERT INTO favourites (id_user)
SELECT id
FROM "user"
WHERE role = 'обычный пользователь'
LIMIT 5;

-- Для каждого избранного списка добавим по 2 случайных продукта
WITH favourites_list AS (
    SELECT id, row_number() OVER () as rn FROM favourites
),
product_list AS (
    SELECT id, row_number() OVER () as rn FROM product
),
series AS (
    SELECT generate_series(1, 10) as i
)
INSERT INTO favourites_item (id_favourites, id_product)
SELECT 
    f.id,
    p.id
FROM series s
JOIN favourites_list f ON f.rn = ((s.i - 1) / 2) + 1  -- по 2 товара на избранное
JOIN product_list p ON p.rn = s.i;
