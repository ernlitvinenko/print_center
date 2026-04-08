-- we don't know how to generate root <with-no-name> (class Root) :(

comment on database postgres is 'default administrative connection database';

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

comment on table profile is 'Таблица профилей пользователей';

comment on column profile.id is 'Идентификатор пользователя';

comment on column profile.first_name is 'Имя';

comment on column profile.last_name is 'Фамилия';

comment on column profile.father_name is 'Отчество (При наличии)';

comment on column profile.email is 'эл. почта';

comment on column profile.phone_dgt is 'Номер телефона (только цифры)';

comment on column profile.password is 'Хешированный пароль пользователя';

alter table profile
    owner to ernlitvinenko;

create table role
(
    id   smallserial
        constraint role_pk
            primary key,
    name varchar(32) not null
);

comment on table role is 'Таблица ролей пользователей';

comment on column role.id is 'Идентификатор роли';

comment on column role.name is 'Наименование роли';

alter table role
    owner to ernlitvinenko;

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

comment on table profile_role is 'Кросс-таблица пользователей и их ролей';

comment on column profile_role.id is 'Идентификатор записи';

comment on column profile_role.profile_id is 'Идентификатор профиля пользователя';

comment on column profile_role.role_id is 'Идентификатор роли';

alter table profile_role
    owner to ernlitvinenko;

create table nomenclature
(
    id        serial
        constraint nomenclature_pk
            primary key,
    gost_id   varchar(8),
    image_url varchar(256) not null
);

comment on table nomenclature is 'Номенклатура Наклеек';

comment on column nomenclature.id is 'Идентификатор наклейки';

comment on column nomenclature.gost_id is 'Номер по ГОСТу (при наличии)';

comment on column nomenclature.image_url is 'Изображение';

alter table nomenclature
    owner to ernlitvinenko;

create table material
(
    id   serial
        constraint material_pk
            primary key,
    name varchar(32) not null
);

comment on table material is 'Материалы производства типографии';

comment on column material.id is 'Идентификатор материала';

comment on column material.name is 'Наименование материала';

alter table material
    owner to ernlitvinenko;

create table size
(
    id     serial
        constraint size_pk
            primary key,
    name   varchar(32) not null,
    width  integer     not null,
    height integer     not null
);

comment on table size is 'Размеры наклеек';

comment on column size.name is 'Наименование размера (напрмер 100x100)';

comment on column size.width is 'Ширина (мм)';

comment on column size.height is 'Высота (мм)';

alter table size
    owner to ernlitvinenko;

create table status
(
    id   serial
        constraint status_pk
            primary key,
    name varchar(16) not null
);

comment on table status is 'Статусы готовности';

comment on column status.id is 'Идентификатор статуса';

comment on column status.name is 'Наименование статуса';

alter table status
    owner to ernlitvinenko;

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

comment on table counterparties is 'Контрагенты';

comment on column counterparties.id is 'Идентификатор контрагента';

comment on column counterparties.is_individual is 'Флаг, является ли контрагент физическим лицом';

comment on column counterparties.unp is 'УНП - компании';

comment on column counterparties.name is 'Полное наименование организации (юрлицо) / имя контрагента (физлицо)';

comment on column counterparties.address is 'Адрес организации';

comment on column counterparties.email is 'Электронная почта контрагента';

comment on column counterparties.phone_number_dgt is 'Контактный номер телефона контрагента';

comment on column counterparties.contact_name is 'Контактное лицо контрагента';

alter table counterparties
    owner to ernlitvinenko;

create table "order"
(
    id                serial
        constraint order_pk
            primary key,
    date_from         timestamp with time zone default now() not null,
    date_till         timestamp with time zone               not null,
    manager_id        integer
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

comment on table "order" is 'Заказы';

comment on column "order".id is 'Идентификатор заказа';

comment on column "order".date_from is 'Дата/время создания заказа';

comment on column "order".date_till is 'Планируемое время завершения заказа';

comment on column "order".manager_id is 'Ответственный менеджер';

comment on column "order".counterparties_id is 'Идентификатор контрагента';

comment on column "order".status_id is 'Идентификатор статуса готовности';

comment on column "order".priority is 'Приоритет срочности выполнения заказа от 0 до 3 (0 - стандартный приоритет, 3 - срочный приоритет (HIGH, впереди в очереди))';

alter table "order"
    owner to ernlitvinenko;

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

comment on table order_items is 'Составляющие заказа';

comment on column order_items.id is 'Идентификатор составляющей';

comment on column order_items.nomenclature_id is 'Идентификатор номенклатуры';

comment on column order_items.order_id is 'Идентификатор заказа';

comment on column order_items.size_id is 'Идентификатор размера';

comment on column order_items.material_id is 'Идентификатор материала';

comment on column order_items.planning_count is 'Количество заказанных деталей';

comment on column order_items.total_count is 'Количество произведенных деталей';

alter table order_items
    owner to ernlitvinenko;

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

comment on table order_comments is 'Коментарии к заказам';

comment on column order_comments.id is 'Идентификатор комментария';

comment on column order_comments.order_id is 'Идентификатор заказа';

comment on column order_comments.text is 'Текст комментария';

comment on column order_comments.profile_id is 'Идентификатор пользователя оставивший коментарий';

alter table order_comments
    owner to ernlitvinenko;

create table files
(
    id  serial
        constraint files_pk
            primary key,
    url text not null
);

comment on table files is 'Файлы';

comment on column files.id is 'Идентификатор файла';

comment on column files.url is 'Ссылка на файл';

alter table files
    owner to ernlitvinenko;

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

comment on table order_items_file is 'Файлы на наклейки';

alter table order_items_file
    owner to ernlitvinenko;

