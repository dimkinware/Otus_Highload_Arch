CREATE TABLE users(
    id UUID not null,
    first_name VARCHAR(255),
    second_name VARCHAR(255),
    birth_date bigint,
    gender smallint,
    bio VARCHAR(255),
    city VARCHAR(255),
    pwd_hash VARCHAR(255)
);

CREATE TABLE tokens(
    user_id varchar(36) not null,
    token varchar(36) not null,
    expired_time bigint
);

ALTER TABLE users ADD PRIMARY KEY (id);
CREATE INDEX idx_first_name_start_letter ON users (first_name varchar_pattern_ops);
CREATE INDEX idx_second_name_start_letter ON users (second_name varchar_pattern_ops);

-- MIGRATION 1

ALTER TABLE tokens ALTER COLUMN user_id TYPE uuid USING user_id::uuid;
ALTER TABLE tokens ADD CONSTRAINT fk_tokens_user_id FOREIGN KEY (user_id) REFERENCES users(id);

CREATE TABLE friends(
    user_id_f1 UUID,
    user_id_f2 UUID,
    FOREIGN KEY (user_id_f1) REFERENCES users(id),
    FOREIGN KEY (user_id_f1) REFERENCES users(id),
    UNIQUE (user_id_f1, user_id_f2)
);

CREATE TABLE posts(
    id UUID not null,
    author_user_id UUID, 
    post_text TEXT,
    PRIMARY KEY (id),
    FOREIGN KEY (author_user_id) REFERENCES users(id),
);
