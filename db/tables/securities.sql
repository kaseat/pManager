CREATE TABLE securities (
	isin varchar(12) NOT NULL,
	ticker varchar(12) NOT NULL,
	figi varchar(12) NOT NULL,
	currency char(3) NOT NULL,
	asset_type smallint NOT NULL,
	title varchar(100) NOT NULL,
	CONSTRAINT pk_securities PRIMARY KEY (isin),
    CONSTRAINT fk_securities_currency FOREIGN KEY(currency) REFERENCES currencies(code),
    CONSTRAINT fk_securities_securities_types FOREIGN KEY(asset_type) REFERENCES securities_types(id)
);
