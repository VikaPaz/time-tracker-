CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


create table if not exists users
(
    id         uuid default uuid_generate_v4() primary key,
    name       varchar(100),
    surname    varchar(100),
    patronymic varchar(100),
    address    text,
    passport   varchar(100)
);

create table if not exists tasks
(
    id      uuid default uuid_generate_v4() primary key,
    task    text,
    user_id uuid references users on delete cascade
);

create table if not exists labor_time
(
    id      uuid default uuid_generate_v4() primary key,
    start   timestamp default (now() at time zone 'utc'),
    stop    timestamp,
    task_id uuid references tasks on delete cascade
);

-- drop table users, tasks, labor_time