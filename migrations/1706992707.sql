CREATE TABLE posts(
    id varchar(40) NOT NULL PRIMARY KEY,
    content varchar(25) NOT NULL,
    reply_count int NOT NULL DEFAULT 0,
    like_count int NOT NULL DEFAULT 0,
    dislike_count int NOT NULL default 0,
    impressions int NOT NULL default 0,
    save_count int NOT NULL default 0,
    visibility int NOT NULL default 1,
    reply_to varchar(40) DEFAULT NULL, 
    previous_version_id varchar(40) DEFAULT NULL, 
    created_at timestamp with time zone not null default NOW(),
    created_by varchar(25),
    updated_at timestamp with time zone not null default NOW(),
    updated_by varchar(25),
    is_deleted boolean not null default false
);