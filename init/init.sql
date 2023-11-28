create table cluster_amount_tasks
(
    uuid       varchar(191) not null,
    type       tinyint      not null,
    operators  varchar(200) not null,
    amount     varchar(100) not null,
    txhash     varchar(100) null,
    errorMsg   longtext     null,
    created_at datetime(3)  null,
    updated_at datetime(3)  null,
    primary key (uuid)
);

create index idx_cluster_amount_tasks_type
    on cluster_amount_tasks (type);

create index idx_cluster_amount_tasks_updated_at
    on cluster_amount_tasks (updated_at);

create table encrypts
(
    link               varchar(70) not null,
    encrypted_password text        not null,
    primary key (link)
);

create table fee_recipient_tasks
(
    uuid          varchar(191) not null,
    fee_recipient varchar(100) null,
    txhash        varchar(100) null,
    errorMsg      longtext     null,
    created_at    datetime(3)  null,
    updated_at    datetime(3)  null,
    primary key (uuid)
);

create index idx_fee_recipient_tasks_updated_at
    on fee_recipient_tasks (updated_at);

create table validator_infos
(
    id          bigint unsigned auto_increment
        primary key,
    publicKey   varchar(100)      not null,
    keystore    longtext          not null,
    operators   longtext          not null,
    state       tinyint default 0 not null,
    pendingTime bigint  default 0 not null,
    created_at  datetime(3)       null,
    updated_at  datetime(3)       null,
    constraint idx_validator_infos_public_key
        unique (publicKey)
);

create table cluster_validator_tasks
(
    uuid         varchar(191)    not null,
    type         tinyint         not null,
    validator_id bigint unsigned null,
    operators    varchar(200)    not null,
    txhash       varchar(100)    null,
    errorMsg     longtext        null,
    created_at   datetime(3)     null,
    updated_at   datetime(3)     null,
    primary key (uuid),
    constraint fk_validator_infos_tasks
        foreign key (validator_id) references validator_infos (id)
);

create index idx_cluster_validator_tasks_type
    on cluster_validator_tasks (type);

create index idx_cluster_validator_tasks_updated_at
    on cluster_validator_tasks (updated_at);

create index idx_validator_infos_state
    on validator_infos (state);

create index idx_validator_infos_updated_at
    on validator_infos (updated_at);

