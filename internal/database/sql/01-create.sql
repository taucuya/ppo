create extension if not exists "uuid-ossp";

drop table if exists order_item cascade;
drop table if exists "order" cascade;
drop table if exists review cascade;
drop table if exists basket_item cascade;
drop table if exists basket cascade;
drop table if exists favourites_item cascade;
drop table if exists favourites cascade;
drop table if exists worker cascade;
drop table if exists product cascade;
drop table if exists brand cascade;
drop table if exists "user" cascade;
drop table if exists token cascade;
drop table if exists order_worker cascade;

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

create table if not exists brand (
    id uuid not null default uuid_generate_v4(),
    name varchar(255) not null,
    description text,
    price_category varchar(50),
    constraint brand_pkey primary key (id)
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

create table if not exists favourites (
    id uuid primary key default uuid_generate_v4(),
    id_user uuid
);

create table if not exists favourites_item (
    id uuid primary key default uuid_generate_v4(),
    id_favourites uuid,
    id_product uuid
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
    price decimal(10,2)
);

create table if not exists order_worker (
    id_order uuid,
    id_worker uuid,
    primary key (id_order, id_worker),
    constraint fk_order foreign key (id_order) references "order"(id) on delete set null,
    constraint fk_worker foreign key (id_worker) references worker(id) on delete set null
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