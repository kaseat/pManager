CREATE SEQUENCE securities_id_seq;

CREATE TABLE securities (
	id integer DEFAULT pseudo_encrypt_24(CAST (nextval('securities_id_seq') AS integer)),
	isin varchar(12) NOT NULL,
	ticker varchar(12) NOT NULL,
	figi varchar(12) NOT NULL,
	currency char(3) NOT NULL,
	exchange_id smallint NOT NULL,
	asset_type smallint NOT NULL,
	title varchar(100) NOT NULL,
	price_upd_time date NULL,
	CONSTRAINT pk_securities_id PRIMARY KEY (id),
    CONSTRAINT fk_securities_currency FOREIGN KEY(currency) REFERENCES currencies(code),
    CONSTRAINT fk_securities_securities_types FOREIGN KEY(asset_type) REFERENCES securities_types(id),
    CONSTRAINT fk_securities_exchange FOREIGN KEY(exchange_id) REFERENCES exchange(id)
);

CREATE UNIQUE INDEX pk_securities_isin ON securities(isin, currency);
