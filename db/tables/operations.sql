CREATE TABLE operations (
    pid integer NOT NULL,
    isin char(12) NOT NULL,
	time timestamp NOT NULL,
    op_id smallint NOT NULL,
	vol integer NOT NULL,
	price numeric(20,6) NOT NULL,
	CONSTRAINT pk_operations PRIMARY KEY (pid, isin, time),
    CONSTRAINT fk_operations_users FOREIGN KEY(pid) REFERENCES users(id),
    CONSTRAINT fk_operations_operation_types FOREIGN KEY(op_id) REFERENCES operation_types(id),
    CONSTRAINT fk_operations_securities FOREIGN KEY(isin) REFERENCES securities(isin)
);
