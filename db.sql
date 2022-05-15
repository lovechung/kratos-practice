create table test.car
(
    id            bigint auto_increment
        primary key,
    model         varchar(255) not null,
    registered_at datetime     not null,
    user_id       bigint       not null
)
    collate = utf8mb4_bin;

create table test.user
(
    id         bigint auto_increment
        primary key,
    username   varchar(20)  not null,
    password   varchar(255) not null,
    created_at datetime     not null,
    updated_at datetime     not null
)
    collate = utf8mb4_bin;
