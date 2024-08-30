CREATE EXTENSION if not exists "uuid-ossp";
CREATE EXTENSION if not exists "pgcrypto";



create table if not exists users (

    Id serial primary key,
    Uuid uuid default uuid_generate_v4() unique,
    Name varchar(30) not null,
    Login varchar(30) unique not null,
    Password varchar(200) not null,
    Email varchar(256) not null

);

create table if not exists refresh_sessions (

    Id serial primary key,
    Token varchar(300),
    Ip varchar(22) not null,
    time_life bigint,
    Owner_uuid uuid not null references users(uuid) on delete cascade

);