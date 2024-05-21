CREATE TABLE users(
    id varchar(36) not null,
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