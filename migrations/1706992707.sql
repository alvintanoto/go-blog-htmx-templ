-- post --
CREATE TABLE posts(
    id varchar(40) NOT NULL PRIMARY KEY,
    user_id varchar(40) NOT NULL, 
    content varchar(150) NOT NULL,
    reply_count int NOT NULL DEFAULT 0,
    like_count int NOT NULL DEFAULT 0,
    dislike_count int NOT NULL default 0,
    impressions int NOT NULL default 0,
    save_count int NOT NULL default 0,
    visibility int NOT NULL default 1,
    reply_to varchar(40) DEFAULT NULL, 
    original_version_id varchar(40) DEFAULT NULL, 
    is_draft bool not null default false,
    posted_at timestamp with time zone default null,
    created_at timestamp with time zone not null default NOW(),
    created_by varchar(40),
    updated_at timestamp with time zone not null default NOW(),
    updated_by varchar(40),
    is_deleted boolean not null default false,
    CONSTRAINT fk_post_user
        FOREIGN KEY(user_id)
            REFERENCES blog_user(id)
);