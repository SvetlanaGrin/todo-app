CREATE TABLE  IF NOT EXISTS users
( 
    id   SERIAL primary key,
    name    varchar(255) NOT NULL ,
    username varchar(255) UNIQUE NOT NULL ,
    password_hash varchar(255) NOT NULL
);


CREATE TABLE  IF NOT EXISTS  todo_lists
(
    id      SERIAL PRIMARY KEY,
    title    varchar(255) not null,
    description  varchar(255)
);


CREATE TABLE  IF NOT EXISTS users_lists
(
    id      SERIAL PRIMARY KEY,
    user_id INTEGER,
    list_id INTEGER,
    foreign key (user_id)   references  users(id) on delete cascade ,
    foreign key (list_id)   references  todo_lists(id) on delete cascade

);

CREATE TABLE  IF NOT EXISTS todo_items
(   
     id      SERIAL PRIMARY KEY,
    title    varchar(255) not null,
    description  varchar(255),
    done  boolean not null default false
);

CREATE TABLE  IF NOT EXISTS lists_items
(
    id      SERIAL PRIMARY KEY,
    item_id INTEGER,
    list_id INTEGER,
    foreign key (item_id)   references  todo_items(id) on delete cascade,
    foreign key (list_id)   references  todo_lists(id) on delete cascade
);
