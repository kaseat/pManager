CREATE TABLE prices (
    sid integer NOT NULL,
	date date NOT NULL,
	vol integer NOT NULL,
	price numeric(20,6) NOT NULL,
	CONSTRAINT pk_prices PRIMARY KEY (sid, date),
    CONSTRAINT fk_prices_securities FOREIGN KEY(sid) REFERENCES securities(id)
);
