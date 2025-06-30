create table
    users (
        id serial primary key,
        email varchar(254) unique NOT NULL,
        password_hash varchar(72) NOT NULL,
        created_at date default now (),
        updated_at date
    );