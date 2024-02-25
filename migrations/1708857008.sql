CREATE TABLE ref_user_config(
    key varchar(50) not null primary key,
    description varchar(200)
);

CREATE TABLE rel_user_config(
    id bigserial primary key, 
    user_id varchar(40) not null,
    config_id varchar(50) not null,
    value varchar(200) not null default '',
    CONSTRAINT fk_rel_user_config_user
        FOREIGN KEY(user_id)
            REFERENCES blog_user(id)
            ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_rel_user_config_config
        FOREIGN KEY(config_id)
            REFERENCES ref_user_config(key)
            ON DELETE RESTRICT ON UPDATE RESTRICT
);

INSERT INTO public.ref_user_config ("key")
	VALUES ('USER_THEME');

ALTER TABLE public.rel_user_config
ADD CONSTRAINT unique_user_config UNIQUE (user_id, config_id);
