CREATE EXTENSION if not exists "uuid-ossp";

create table if not exists notes (
    owner_uuid uuid not null,
    title text not null,
    content text not null
);