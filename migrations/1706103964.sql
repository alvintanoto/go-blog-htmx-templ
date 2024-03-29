CREATE TABLE blog_user(
    id varchar(40) NOT NULL PRIMARY KEY,
    username varchar(25) NOT NULL,
    email varchar(100) NOT NULL,
    password varchar(100) NOT NULL,
    last_logged_in timestamp with time zone default NULL,
    created_at timestamp with time zone not null default NOW(),
    created_by varchar(40),
    updated_at timestamp with time zone not null default NOW(),
    updated_by varchar(40),
    is_deleted boolean not null default false
);

CREATE UNIQUE INDEX constraint_unique_username ON blog_user (username) WHERE is_deleted=false;