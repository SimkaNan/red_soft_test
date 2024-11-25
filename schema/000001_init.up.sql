CREATE TABLE IF NOT EXISTS users(
    id serial unique not null,
    name varchar(50) not null,
    surname varchar(50) not null,
    middle_name varchar(50) not null,
    age integer not null,
    nation varchar(50) not null,
    gender varchar(10) not null
);

CREATE TABLE IF NOT EXISTS emails(
    id serial unique not null,
    email varchar(50) unique not null,
    user_id integer references users(id) on delete cascade
);

CREATE TABLE IF NOT EXISTS friendships(
  user_id1 integer references users(id) on delete cascade,
  user_id2 integer references users(id) on delete cascade,
  primary key (user_id1, user_id2)
);