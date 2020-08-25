CREATE SEQUENCE users_id_seq;

CREATE TABLE users (
    id integer DEFAULT pseudo_encrypt_24(CAST (nextval('users_id_seq') AS integer)),
    hash varchar(150) NOT NULL,
    login varchar(50) NOT NULL,
    state varchar(24) NULL,
    token json NULL,
	CONSTRAINT pk_users PRIMARY KEY (id)
);

CREATE UNIQUE INDEX pk_users_login ON users(login);
