drop trigger if exists order_creation_trigger on "order";

create or replace function process_order_creation()
returns trigger as $$
declare
    basket_id uuid;
    item record;
    product_amount int;
    total_price decimal(10,2);
begin
    select id into basket_id from basket where id_user = new.id_user;
    if not found then
        raise exception 'корзина для пользователя % не найдена', new.id_user;
    end if;

    for item in select * from basket_item where id_basket = basket_id loop
        select amount into product_amount from product where id = item.id_product;
        if product_amount < item.amount then
            raise exception 'недостаточно товара на складе (id товара: %, доступно: %, требуется: %)', 
                            item.id_product, product_amount, item.amount;
        end if;
        
        update product set amount = amount - item.amount where id = item.id_product;
        
        insert into order_item (id_product, id_order, amount)
        values (item.id_product, new.id, item.amount);
    end loop;
    
    select sum(p.price * oi.amount) into total_price
    from order_item oi
    join product p on oi.id_product = p.id
    where oi.id_order = new.id;
    
    update "order" set price = total_price where id = new.id;
    
    delete from basket_item where id_basket = basket_id;
    delete from basket where id = basket_id;
    
    return new;
end;
$$ language plpgsql;

create trigger order_creation_trigger
after insert on "order"
for each row
execute function process_order_creation();
