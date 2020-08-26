CREATE SEQUENCE operations_id_seq;

CREATE TABLE operations (
    id integer DEFAULT pseudo_encrypt_24(CAST (nextval('operations_id_seq') AS integer)),
    pid integer NOT NULL,
    isin varchar(12) NOT NULL,
	time timestamp NOT NULL,
    op_id smallint NOT NULL,
	vol integer NOT NULL,
	price numeric(20,6) NOT NULL,
	CONSTRAINT pk_operations PRIMARY KEY (id),
    CONSTRAINT fk_operations_portfolios FOREIGN KEY(pid) REFERENCES portfolios(id),
    CONSTRAINT fk_operations_operation_types FOREIGN KEY(op_id) REFERENCES operation_types(id),
    CONSTRAINT fk_operations_securities FOREIGN KEY(isin) REFERENCES securities(isin)
);

CREATE INDEX ix_operations ON operations(pid, isin, time);
