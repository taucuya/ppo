-- вставляем пользователей
INSERT INTO "user" (name, date_of_birth, mail, password, phone, address, status, role) VALUES
('Иван Иванов', '1990-05-15', 'ivan@example.com', 'password123', '+79161234567', 'ул. Ленина, 10', 'active', 'customer'),
('Петр Петров', '1985-08-20', 'petr@example.com', 'password456', '+79162345678', 'ул. Пушкина, 5', 'active', 'customer'),
('Анна Сидорова', '1995-02-10', 'anna@example.com', 'password789', '+79163456789', 'пр. Мира, 15', 'active', 'customer'),
('Мария Кузнецова', '1988-11-30', 'maria@example.com', 'password012', '+79164567890', 'ул. Гагарина, 3', 'active', 'customer'),
('Алексей Смирнов', '1992-07-25', 'alex@example.com', 'password345', '+79165678901', 'ул. Кирова, 7', 'active', 'customer'),
('Елена Волкова', '1993-04-18', 'elena@example.com', 'password678', '+79166789012', 'ул. Советская, 22', 'active', 'customer'),
('Дмитрий Федоров', '1987-09-05', 'dmitry@example.com', 'password901', '+79167890123', 'ул. Садовая, 14', 'active', 'customer'),
('Ольга Морозова', '1991-12-12', 'olga@example.com', 'password234', '+79168901234', 'ул. Лесная, 8', 'active', 'customer'),
('Сергей Николаев', '1989-06-28', 'sergey@example.com', 'password567', '+79169012345', 'пр. Победы, 17', 'active', 'customer'),
('Наталья Павлова', '1994-03-08', 'natalya@example.com', 'password890', '+79160123456', 'ул. Цветочная, 9', 'active', 'customer');

-- вставляем бренды
insert into brand (name, description, price_category) values
('nike', 'американский производитель спортивной одежды', 'premium'),
('adidas', 'немецкий производитель спортивной одежды', 'premium'),
('puma', 'немецкий производитель спортивной одежды', 'medium'),
('reebok', 'производитель спортивной одежды', 'medium'),
('under armour', 'американский бренд спортивной одежды', 'premium'),
('new balance', 'американский производитель обуви', 'medium'),
('asics', 'японский производитель спортивной обуви', 'medium'),
('fila', 'итальянский спортивный бренд', 'economy'),
('decathlon', 'французский производитель спортивных товаров', 'economy'),
('columbia', 'американский производитель outdoor-одежды', 'premium');

-- вставляем продукты
insert into product (name, description, price, category, amount, id_brand, pic_link, art) values
('кроссовки air max', 'кроссовки для бега с технологией air', 12000.00, 'обувь', 15, (select id from brand where name = 'nike'), '/images/airmax.jpg', 'art001'),
('футболка спортивная', 'дышащая футболка для тренировок', 2500.00, 'одежда', 50, (select id from brand where name = 'adidas'), '/images/t-shirt.jpg', 'art002'),
('шорты беговые', 'легкие шорты с карманом для телефона', 3500.00, 'одежда', 30, (select id from brand where name = 'puma'), '/images/shorts.jpg', 'art003'),
('кепка бейсболка', 'кепка с защитой от солнца', 1800.00, 'аксессуары', 25, (select id from brand where name = 'reebok'), '/images/cap.jpg', 'art004'),
('рюкзак спортивный', 'рюкзак с отделением для ноутбука', 4500.00, 'аксессуары', 20, (select id from brand where name = 'under armour'), '/images/backpack.jpg', 'art005'),
('носки спортивные', 'носки с анатомической поддержкой', 800.00, 'аксессуары', 100, (select id from brand where name = 'new balance'), '/images/socks.jpg', 'art006'),
('леггинсы женские', 'леггинсы с высокой талией', 4200.00, 'одежда', 35, (select id from brand where name = 'asics'), '/images/leggings.jpg', 'art007'),
('куртка ветровка', 'ветровка с мембраной', 7500.00, 'одежда', 18, (select id from brand where name = 'fila'), '/images/jacket.jpg', 'art008'),
('перчатки для бега', 'перчатки для бега в холодную погоду', 2200.00, 'аксессуары', 40, (select id from brand where name = 'decathlon'), '/images/gloves.jpg', 'art009'),
('очки спортивные', 'очки с защитой от ультрафиолета', 3800.00, 'аксессуары', 22, (select id from brand where name = 'columbia'), '/images/glasses.jpg', 'art010');

-- вставляем корзины
insert into basket (id_user, date) 
select id, now() - (random() * 30 || ' days')::interval 
from "user" limit 10;

-- вставляем работников
WITH ranked_users AS (
  SELECT id, row_number() OVER () as rn FROM "user" LIMIT 10
)
INSERT INTO worker (id_user, job_title)
SELECT id,
  CASE rn
    WHEN 1 THEN 'менеджер'
    WHEN 2 THEN 'администратор'
    WHEN 3 THEN 'курьер'
    WHEN 4 THEN 'маркетолог'
    WHEN 5 THEN 'аналитик'
    WHEN 6 THEN 'директор'
    WHEN 7 THEN 'бухгалтер'
    WHEN 8 THEN 'hr-специалист'
    WHEN 9 THEN 'логист'
    WHEN 10 THEN 'контент-менеджер'
  END
FROM ranked_users;


-- вставляем заказы
WITH ranked_users AS (
  SELECT id, address, row_number() OVER () as rn FROM "user" LIMIT 10
),
worker_list AS (
  SELECT id, row_number() OVER () as rn FROM worker
)
INSERT INTO "order" (date, id_user, address, status, price, id_worker)
SELECT 
  now() - (random() * 90 || ' days')::interval,
  u.id,
  u.address,
  CASE (u.rn % 4)
    WHEN 0 THEN 'completed'
    WHEN 1 THEN 'processing'
    WHEN 2 THEN 'shipped'
    WHEN 3 THEN 'cancelled'
  END,
  (random() * 10000 + 1000)::numeric(10,2),
  w.id
FROM ranked_users u
JOIN worker_list w ON w.rn = ((u.rn - 1) % 10 + 1);


-- вставляем элементы корзины
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


-- вставляем элементы заказа
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


-- вставляем отзывы
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
    WHEN 0 THEN 'отличный товар, всем рекомендую!'
    WHEN 1 THEN 'хорошее качество за свои деньги'
    WHEN 2 THEN 'не совсем то, что ожидал'
    WHEN 3 THEN 'пока не понял, как пользоваться'
    WHEN 4 THEN 'лучшее, что я покупал за последнее время'
  END,
  now() - (random() * 180 || ' days')::interval
FROM series s
JOIN product_list p ON p.rn = s.i
JOIN user_list u ON u.rn = s.i;


-- вставляем пустые токены
insert into token (rtoken) values 
(null), (null), (null), (null), (null), 
(null), (null), (null), (null), (null);