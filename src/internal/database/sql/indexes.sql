create table if not exists product (
    id uuid primary key default uuid_generate_v4(),
    name varchar(255),
    description text,
    price decimal(10,2),
    category varchar(100),
    amount int,
    id_brand uuid,
    pic_link text,
    art varchar(50)
);

insert into product (name, description, price, category, amount, id_brand, pic_link, art)
select 
    'product ' || i,
    'description for product ' || i,
    round((random() * 1000)::numeric, 2),
    case when i % 2 = 0 then 'electronics' else 'home' end,
    (random() * 100)::int,
    uuid_generate_v4(),
    'https://image.url/' || i,
    'art' || (100000 + i)
from generate_series(1, 1000000) as s(i);

\timing

create index idx_product_art_btree on product(art);
create index idx_product_price_btree on product(price);

-- B-Tree

explain analyze select * from product where art = 'art123456';

discard plans;

explain analyze select * from product where price = 499.9;

discard plans;


drop index if exists idx_product_art_btree;
drop index if exists idx_product_price_btree;

create index idx_product_art_hash on product using hash(art);
create index idx_product_price_hash on product using hash(price);

-- Hash
explain analyze select * from product where art = 'art123456';

discard plans;

explain analyze select * from product where price = 499.9;

discard plans;

drop index if exists idx_product_art_hash;
drop index if exists idx_product_price_hash;

create index idx_product_art_gin on product using gin(to_tsvector('english', art));
create index idx_product_price_gin on product using gin(to_tsvector('english', price::text));

-- GIN
explain analyze select * from product 
where to_tsvector('english', art) @@ plainto_tsquery('english', 'art123456');

discard plans;

explain analyze select * from product 
where to_tsvector('english', price::text) @@ plainto_tsquery('english', '499.9');

discard plans;

drop index if exists idx_product_art_gin;
drop index if exists idx_product_price_gin;
