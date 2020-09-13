CREATE TABLE exchange (
    id smallint NOT NULL,
    code varchar(10) NOT NULL,
    title varchar(150) NULL,
	CONSTRAINT pk_exchange PRIMARY KEY (id)
);
