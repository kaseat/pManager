CREATE TABLE user_sync (
    uid integer NOT NULL,
    provider_id smallint NOT NULL,
    last_sync timestamp NOT NULL,
	CONSTRAINT pk_user_sync PRIMARY KEY (uid,provider_id),
    CONSTRAINT fk_user_sync_users FOREIGN KEY(uid) REFERENCES users(id),
    CONSTRAINT fk_user_sync_sync_providers FOREIGN KEY(provider_id) REFERENCES sync_providers(id)
);
