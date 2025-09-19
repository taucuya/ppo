create or replace function accept_order_trigger()
returns trigger as $$
declare
    found_worker_id uuid;
begin
    select id into found_worker_id from worker where id_user = new.id_worker;

    new.id_worker := found_worker_id;
    
    update "order" set status = 'принятый' where id = new.id_order;
    
    return new;
end;
$$ language plpgsql;

create trigger accept_order_trigger
before insert on order_worker
for each row
execute function accept_order_trigger();