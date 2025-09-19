create role guest login;
create role client login;
create role worker login;
create role administrator login superuser;

grant select on brand, product, review to guest;

-- client
grant insert(name, date_of_birth, mail, password, phone, address) 
on "user" to client;
grant usage, select on sequence "user_id_seq" to client;

grant select on brand, product to client;

grant select, insert, delete on basket_item to client;
grant usage, select on sequence basket_item_id_seq to client;

grant insert on "order" to client;
grant usage, select on sequence order_id_seq to client;

grant select on order_item to client;
grant usage, select on sequence order_item_id_seq to client;

grant select, insert on review to client;
grant usage, select on sequence review_id_seq to client;

-- worker
grant select, update(status) on "order" to worker;

grant select on order_item to worker;

grant insert, select on order_worker to worker;

-- admin
grant all privileges on all tables in schema public to administrator;
grant all privileges on all sequences in schema public to administrator;
