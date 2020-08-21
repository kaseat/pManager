CREATE TABLE securities (
	isin char(12) NOT NULL,
	ticker char(10) NOT NULL,
	figi char(12) NOT NULL,
	title varchar(100) NOT NULL,
	CONSTRAINT pk_securities PRIMARY KEY (isin)
);
