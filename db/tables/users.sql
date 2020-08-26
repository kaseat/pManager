CREATE SEQUENCE users_id_seq;

CREATE TABLE users (
    id integer DEFAULT pseudo_encrypt_24(CAST (nextval('users_id_seq') AS integer)),
    login varchar(50) NOT NULL,
    hash varchar(150) NOT NULL,
    role_id smallint NOT NULL,
    email varchar(100) NULL,
    g_sync_state varchar(24) NULL,
    g_sync_token json NULL,
	CONSTRAINT pk_users PRIMARY KEY (id),
    CONSTRAINT fk_users_user_roles FOREIGN KEY(role_id) REFERENCES user_roles(id)
);

CREATE UNIQUE INDEX pk_users_login ON users(login);
