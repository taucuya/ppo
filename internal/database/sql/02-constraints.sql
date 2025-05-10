-- USER
alter table "user"
alter column "name" set not null,
alter column "mail" set not null,
alter column "password" set not null,
add constraint "user_mail_format" check (mail ~* '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'),
add constraint "user_phone_format" check (phone ~* '^(\+7|8)[\s\-()]?\d{3}[\s\-()]?\d{3}[\s\-()]?\d{2}[\s\-()]?\d{2}$'),
add constraint "user_mail_unique" unique (mail),
add constraint "user_phone_unique" unique (phone);

-- PRODUCT
alter table "product"
alter column "name" set not null,
alter column "price" set not null,
alter column "amount" set not null,
alter column "art" set not null,
add constraint "product_art_unique" unique (art),
add constraint "fk_product_brand" foreign key ("id_brand") references "brand"("id") on delete set null;

-- BASKET
alter table "basket"
alter column "date" set default current_timestamp,
add constraint "fk_basket_user" foreign key ("id_user") references "user"("id") on delete cascade;

-- BASKET-ITEM
alter table "basket_item"
add constraint "fk_basket_item_basket" foreign key ("id_basket") references "basket"("id") on delete cascade,
add constraint "fk_basket_item_product" foreign key ("id_product") references "product"("id") on delete cascade;

-- WORKER
alter table "worker"
add constraint "fk_worker_user" foreign key ("id_user") references "user"("id") on delete cascade,
add constraint "worker_user_unique" unique ("id_user");

-- ORDER
alter table "order"
alter column "date" set default current_timestamp,
alter column "price" set not null,
add constraint "fk_order_user" foreign key ("id_user") references "user"("id") on delete cascade;

--ORDER-ITEM
alter table "order_item"
add constraint "fk_order_item_order" foreign key ("id_order") references "order"("id") on delete cascade,
add constraint "fk_order_item_product" foreign key ("id_product") references "product"("id") on delete cascade;

-- REVIEW
alter table "review"
alter column "date" set default current_timestamp,
add constraint "fk_review_user" foreign key ("id_user") references "user"("id") on delete cascade,
add constraint "fk_review_product" foreign key ("id_product") references "product"("id") on delete cascade,
add constraint "review_rating_check" check ("rating" between 1 and 5);