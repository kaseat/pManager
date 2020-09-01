CREATE TABLE sync_providers (
    id smallint NOT NULL,
    name varchar(20) NOT NULL,
    description varchar(50) NOT NULL,
	CONSTRAINT pk_sync_providers PRIMARY KEY (id)
);
