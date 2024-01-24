CREATE TABLE blog_user(
    id varchar(40) NOT NULL PRIMARY KEY,
    username varchar(25) NOT NULL,
    email varchar(100) NOT NULL,
    password varchar(100) NOT NULL,
    last_logged_in timestamp with time zone default NULL,
    created_at timestamp with time zone not null default NOW(),
    created_by varchar(25),
    updated_at timestamp with time zone not null default NOW(),
    updated_by varchar(25),
    is_deleted boolean not null default false
)