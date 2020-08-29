CREATE SEQUENCE portfolios_id_seq;

CREATE TABLE portfolios (
    id integer DEFAULT pseudo_encrypt_24(CAST (nextval('portfolios_id_seq') AS integer)),
    uid integer NOT NULL,
    name varchar(50) NOT NULL,
    title varchar(150) NULL,
	CONSTRAINT pk_portfolios PRIMARY KEY (id),
    CONSTRAINT fk_portfolios_users FOREIGN KEY(uid) REFERENCES users(id)
);
