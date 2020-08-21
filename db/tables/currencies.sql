CREATE TABLE currencies (
    code char(3) NOT NULL,
    title varchar(150) NULL,
	CONSTRAINT pk_currencies PRIMARY KEY (code)
);
