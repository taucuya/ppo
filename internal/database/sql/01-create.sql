create extension if not exists "uuid-ossp";

DROP TABLE IF EXISTS order_item CASCADE;
DROP TABLE IF EXISTS "order" CASCADE;
DROP TABLE IF EXISTS review CASCADE;
DROP TABLE IF EXISTS basket_item CASCADE;
DROP TABLE IF EXISTS basket CASCADE;
DROP TABLE IF EXISTS worker CASCADE;
DROP TABLE IF EXISTS product CASCADE;
DROP TABLE IF EXISTS brand CASCADE;
DROP TABLE IF EXISTS "user" CASCADE;
DROP TABLE IF EXISTS token CASCADE;

create table if not exists "user" (
    id uuid primary key default uuid_generate_v4(),
    name varchar(255),
    date_of_birth date,
    mail varchar(255),
    password varchar(255),
    phone varchar(20),
    address text,
    status varchar(50),
    role varchar(50)
);

CREATE TABLE IF NOT EXISTS brand (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    name varchar(255) NOT NULL,
    description text,
    price_category varchar(50),
    CONSTRAINT brand_pkey PRIMARY KEY (id)
);

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

create table if not exists basket (
    id uuid primary key default uuid_generate_v4(),
    id_user uuid,
    date timestamp without time zone
);

create table if not exists basket_item (
    id uuid primary key default uuid_generate_v4(),
    id_basket uuid,
    id_product uuid,
    amount int
);

create table if not exists worker (
    id uuid primary key default uuid_generate_v4(),
    id_user uuid,
    job_title varchar(100)
);

create table if not exists "order" (
    id uuid primary key default uuid_generate_v4(),
    date timestamp without time zone,
    id_user uuid,
    address text,
    status varchar(50),
    price decimal(10,2),
    id_worker uuid
);

create table if not exists order_item (
    id uuid primary key default uuid_generate_v4(),
    id_order uuid,
    id_product uuid,
    amount int
);

create table if not exists review (
    id uuid primary key default uuid_generate_v4(),
    id_product uuid,
    id_user uuid,
    rating int,
    r_text text,
    date timestamp without time zone
);

create table if not exists token (
    id uuid primary key default uuid_generate_v4(),
    rtoken text
);