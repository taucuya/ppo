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
-- insert into brand (name, description, price_category) values
-- ('nike', 'американский производитель спортивной одежды', 'premium'),
-- ('adidas', 'немецкий производитель спортивной одежды', 'premium'),
-- ('puma', 'немецкий производитель спортивной одежды', 'medium'),
-- ('reebok', 'производитель спортивной одежды', 'medium'),
-- ('under armour', 'американский бренд спортивной одежды', 'premium'),
-- ('new balance', 'американский производитель обуви', 'medium'),
-- ('asics', 'японский производитель спортивной обуви', 'medium'),
-- ('fila', 'итальянский спортивный бренд', 'economy'),
-- ('decathlon', 'французский производитель спортивных товаров', 'economy'),
-- ('columbia', 'американский производитель outdoor-одежды', 'premium');

-- -- вставляем продукты
-- insert into product (name, description, price, category, amount, id_brand, pic_link, art) values
-- ('кроссовки air max', 'кроссовки для бега с технологией air', 12000.00, 'обувь', 15, (select id from brand where name = 'nike'), '/images/airmax.jpg', 'art001'),
-- ('футболка спортивная', 'дышащая футболка для тренировок', 2500.00, 'одежда', 50, (select id from brand where name = 'adidas'), '/images/t-shirt.jpg', 'art002'),
-- ('шорты беговые', 'легкие шорты с карманом для телефона', 3500.00, 'одежда', 30, (select id from brand where name = 'puma'), '/images/shorts.jpg', 'art003'),
-- ('кепка бейсболка', 'кепка с защитой от солнца', 1800.00, 'аксессуары', 25, (select id from brand where name = 'reebok'), '/images/cap.jpg', 'art004'),
-- ('рюкзак спортивный', 'рюкзак с отделением для ноутбука', 4500.00, 'аксессуары', 20, (select id from brand where name = 'under armour'), '/images/backpack.jpg', 'art005'),
-- ('носки спортивные', 'носки с анатомической поддержкой', 800.00, 'аксессуары', 100, (select id from brand where name = 'new balance'), '/images/socks.jpg', 'art006'),
-- ('леггинсы женские', 'леггинсы с высокой талией', 4200.00, 'одежда', 35, (select id from brand where name = 'asics'), '/images/leggings.jpg', 'art007'),
-- ('куртка ветровка', 'ветровка с мембраной', 7500.00, 'одежда', 18, (select id from brand where name = 'fila'), '/images/jacket.jpg', 'art008'),
-- ('перчатки для бега', 'перчатки для бега в холодную погоду', 2200.00, 'аксессуары', 40, (select id from brand where name = 'decathlon'), '/images/gloves.jpg', 'art009'),
-- ('очки спортивные', 'очки с защитой от ультрафиолета', 3800.00, 'аксессуары', 22, (select id from brand where name = 'columbia'), '/images/glasses.jpg', 'art010');

-- -- вставляем корзины
-- insert into basket (id_user, date) 
-- select id, now() - (random() * 30 || ' days')::interval 
-- from "user" limit 10;

-- -- вставляем работников
-- insert into worker (id_user, job_title) 
-- select id, 
-- case (row_number() over ()) 
--   when 1 then 'менеджер'
--   when 2 then 'администратор'
--   when 3 then 'курьер'
--   when 4 then 'маркетолог'
--   when 5 then 'аналитик'
--   when 6 then 'директор'
--   when 7 then 'бухгалтер'
--   when 8 then 'hr-специалист'
--   when 9 then 'логист'
--   when 10 then 'контент-менеджер'
-- end
-- from "user" limit 10;

-- -- вставляем заказы
-- insert into "order" (date, id_user, address, status, price, id_worker)
-- select 
--   now() - (random() * 90 || ' days')::interval,
--   id,
--   address,
--   case (row_number() over () % 4)
--     when 0 then 'completed'
--     when 1 then 'processing'
--     when 2 then 'shipped'
--     when 3 then 'cancelled'
--   end,
--   (random() * 10000 + 1000)::numeric(10,2),
--   (select id from worker limit 1 offset (row_number() over () - 1) % 10)
-- from "user" limit 10;

-- -- вставляем элементы корзины
-- insert into basket_item (id_basket, id_product, amount)
-- select 
--   (select id from basket limit 1 offset (row_number() over () - 1) % 10),
--   (select id from product limit 1 offset (row_number() over () - 1) % 10),
--   (random() * 5 + 1)::int
-- from generate_series(1, 10);

-- -- вставляем элементы заказа
-- insert into order_item (id_order, id_product, amount)
-- select 
--   (select id from "order" limit 1 offset (row_number() over () - 1) % 10),
--   (select id from product limit 1 offset (row_number() over () - 1) % 10),
--   (random() * 3 + 1)::int
-- from generate_series(1, 10);

-- -- вставляем отзывы
-- insert into review (id_product, id_user, rating, r_text, date)
-- select 
--   (select id from product limit 1 offset (row_number() over () - 1) % 10),
--   (select id from "user" limit 1 offset (row_number() over () - 1) % 10),
--   (random() * 4 + 1)::int,
--   case (row_number() over () % 5)
--     when 0 then 'отличный товар, всем рекомендую!'
--     when 1 then 'хорошее качество за свои деньги'
--     when 2 then 'не совсем то, что ожидал'
--     when 3 then 'пока не понял, как пользоваться'
--     when 4 then 'лучшее, что я покупал за последнее время'
--   end,
--   now() - (random() * 180 || ' days')::interval
-- from generate_series(1, 10);

-- -- вставляем пустые токены
-- insert into token (rtoken) values 
-- (null), (null), (null), (null), (null), 
-- (null), (null), (null), (null), (null);