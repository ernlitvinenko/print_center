create table profile
(
    id          serial
        constraint profile_pk
            primary key,
    first_name  varchar(16)                               not null,
    last_name   varchar(32)                               not null,
    father_name varchar(32),
    email       varchar(320)                              not null,
    phone_dgt   bigint                                    not null,
    password    varchar(60) default ''::character varying not null
);


create table role
(
    id   smallserial
        constraint role_pk
            primary key,
    name varchar(32) not null
);


create table profile_role
(
    id         serial
        constraint profile_role_pk
            primary key,
    profile_id integer not null
        constraint profile_role_profile_id_fk
            references profile,
    role_id    integer not null
        constraint profile_role_role_id_fk
            references role
);


create table nomenclature
(
    id        serial
        constraint nomenclature_pk
            primary key,
    gost_id   varchar(8),
    image_url varchar(256) not null
);


create table material
(
    id   serial
        constraint material_pk
            primary key,
    name varchar(32) not null
);


create table size
(
    id     serial
        constraint size_pk
            primary key,
    name   varchar(32) not null,
    width  integer     not null,
    height integer     not null
);


create table status
(
    id   serial
        constraint status_pk
            primary key,
    name varchar(16) not null
);


create table counterparties
(
    id               serial
        constraint counterparties_pk
            primary key,
    is_individual    boolean default false not null,
    unp              integer,
    name             varchar(256)          not null,
    address          varchar(256),
    email            varchar(320),
    phone_number_dgt bigint,
    contact_name     varchar(256)
);


create table "order"
(
    id                serial
        constraint order_pk
            primary key,
    date_from         timestamp with time zone default now() not null,
    date_till         timestamp with time zone               not null,
    manager_id        integer                                not null
        constraint order_profile_id_fk
            references profile,
    counterparties_id integer                                not null
        constraint order_counterparties_id_fk
            references counterparties,
    status_id         smallint                 default 1     not null
        constraint order_status_id_fk
            references status,
    priority          smallint                 default 0     not null
);


create table order_items
(
    id              serial
        constraint order_items_pk
            primary key,
    nomenclature_id integer not null
        constraint order_items_nomenclature_id_fk
            references nomenclature,
    order_id        integer not null
        constraint order_items_order_id_fk
            references "order",
    size_id         integer not null
        constraint order_items_size_id_fk
            references size,
    material_id     integer not null
        constraint order_items___fk
            references material,
    planning_count  integer not null,
    total_count     integer not null
);


create table order_comments
(
    id         serial
        constraint order_comments_pk
            primary key,
    order_id   integer not null
        constraint order_comments_order_id_fk
            references "order",
    text       text,
    profile_id integer not null
        constraint order_comments_profile_id_fk
            references profile
);

create table files
(
    id  serial
        constraint files_pk
            primary key,
    url text not null
);

create table order_items_file
(
    id            serial
        constraint order_items_file_pk
            primary key,
    file_id       integer not null
        constraint order_items_file_files_id_fk
            references files,
    order_item_id integer not null
        constraint order_items_file_order_items_id_fk
            references order_items
);
