CREATE TABLE prices (
    isin varchar(12) NOT NULL,
	date date NOT NULL,
	vol integer NOT NULL,
	price numeric(20,6) NOT NULL,
	CONSTRAINT pk_prices PRIMARY KEY (isin, date),
    CONSTRAINT fk_prices_securities FOREIGN KEY(isin) REFERENCES securities(isin)
);
