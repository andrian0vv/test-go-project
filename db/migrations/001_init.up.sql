create table users
(
    id         serial primary key,
    first_name varchar(100) not null,
    last_name  varchar(100) not null,
    full_name  varchar(210) not null,
    age        int          not null,
    is_married boolean   default false,
    password   varchar(255) not null,
    created_at timestamp default current_timestamp
);

create table products
(
    id          serial primary key,
    description text not null,
    price       int  not null,
    tags        varchar[],
    quantity    int  not null,
    check (quantity >= 0),
    created_at  timestamp default current_timestamp
);

create table orders
(
    id         serial primary key,
    user_id    bigint not null references users (id),
    created_at timestamp default current_timestamp
);

create table order_products
(
    id         serial primary key,
    price      int    not null,
    order_id   bigint not null references orders (id),
    product_id bigint not null references products (id),
    quantity   int    not null
);

create index orders_user_id_idx on orders (user_id);
create index order_products_order_id_idx on order_products (order_id);
create index order_products_product_id_idx on order_products (product_id);

-- For tests
insert into products (description, price, quantity) values ('lorem', 100, 50);
insert into products (description, price, quantity) values ('ipsum', 1000, 10);
insert into products (description, price, quantity) values ('test', 70, 5);

insert into users (first_name, last_name, full_name, age, is_married, password)
values ('ivan', 'ivanov', 'ivan ivanov', 21, false, 'password');