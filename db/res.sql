DROP TABLE IF EXISTS operations CASCADE;
DROP TABLE IF EXISTS operation_types CASCADE;
DROP TABLE IF EXISTS portfolios CASCADE;
DROP TABLE IF EXISTS prices CASCADE;
DROP TABLE IF EXISTS securities CASCADE;
DROP TABLE IF EXISTS securities_types CASCADE;
DROP TABLE IF EXISTS user_sync CASCADE;
DROP TABLE IF EXISTS sync_providers CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS currencies CASCADE;
DROP TABLE IF EXISTS user_roles CASCADE;
DROP TABLE IF EXISTS settings CASCADE;
DROP FUNCTION IF EXISTS pseudo_encrypt_24 CASCADE;
DROP SEQUENCE IF EXISTS operations_id_seq CASCADE;
DROP SEQUENCE IF EXISTS portfolios_id_seq CASCADE;
DROP SEQUENCE IF EXISTS users_id_seq CASCADE;
DROP SEQUENCE IF EXISTS securities_id_seq CASCADE;
CREATE FUNCTION pseudo_encrypt_24(VALUE int) returns int AS $$
DECLARE
l1 int;
l2 int;
r1 int;
r2 int;
i int:=0;
BEGIN
  l1:= (VALUE >> 12) & (4096-1);
  r1:= VALUE & (4096-1);
  WHILE i < 3 LOOP
    l2 := r1;
    r2 := l1 # ((((1366 * r1 + 150889) % 714025) / 714025.0) * (4096-1))::int;
  l1 := l2;
  r1 := r2;
  i := i + 1;
  END LOOP;
  RETURN ((l1 << 12) + r1);
END;
$$ LANGUAGE plpgsql strict immutable;
CREATE TABLE user_roles (
    id smallint NOT NULL,
    id_name varchar(15) NOT NULL,
    title varchar(50) NOT NULL,
	CONSTRAINT pk_user_roles PRIMARY KEY (id)
);
CREATE SEQUENCE users_id_seq;

CREATE TABLE users (
    id integer DEFAULT pseudo_encrypt_24(CAST (nextval('users_id_seq') AS integer)),
    login varchar(50) NOT NULL,
    hash varchar(150) NOT NULL,
    role_id smallint NOT NULL,
    email varchar(100) NULL,
    g_sync_state varchar(24) NULL,
    g_sync_token json NULL,
	CONSTRAINT pk_users PRIMARY KEY (id),
    CONSTRAINT fk_users_user_roles FOREIGN KEY(role_id) REFERENCES user_roles(id)
);

CREATE UNIQUE INDEX pk_users_login ON users(login);
CREATE TABLE sync_providers (
    id smallint NOT NULL,
    name varchar(20) NOT NULL,
    description varchar(50) NOT NULL,
	CONSTRAINT pk_sync_providers PRIMARY KEY (id)
);
CREATE TABLE user_sync (
    uid integer NOT NULL,
    provider_id smallint NOT NULL,
    last_sync timestamp NOT NULL,
	CONSTRAINT pk_user_sync PRIMARY KEY (uid,provider_id),
    CONSTRAINT fk_user_sync_users FOREIGN KEY(uid) REFERENCES users(id),
    CONSTRAINT fk_user_sync_sync_providers FOREIGN KEY(provider_id) REFERENCES sync_providers(id)
);
CREATE SEQUENCE portfolios_id_seq;

CREATE TABLE portfolios (
    id integer DEFAULT pseudo_encrypt_24(CAST (nextval('portfolios_id_seq') AS integer)),
    uid integer NOT NULL,
    name varchar(50) NOT NULL,
    title varchar(150) NULL,
	CONSTRAINT pk_portfolios PRIMARY KEY (id),
    CONSTRAINT fk_portfolios_users FOREIGN KEY(uid) REFERENCES users(id)
);
CREATE TABLE currencies (
    code char(3) NOT NULL,
    title varchar(150) NULL,
	CONSTRAINT pk_currencies PRIMARY KEY (code)
);
CREATE TABLE securities_types (
	id smallint NOT NULL,
	id_name varchar(15) NOT NULL,
	title varchar(50) NOT NULL,
	CONSTRAINT pk_securities_types PRIMARY KEY (id)
);
CREATE SEQUENCE securities_id_seq;

CREATE TABLE securities (
	id integer DEFAULT pseudo_encrypt_24(CAST (nextval('securities_id_seq') AS integer)),
	isin varchar(12) NOT NULL,
	ticker varchar(12) NOT NULL,
	figi varchar(12) NOT NULL,
	currency char(3) NOT NULL,
	asset_type smallint NOT NULL,
	title varchar(100) NOT NULL,
	price_upd_time date NULL,
	CONSTRAINT pk_securities_id PRIMARY KEY (id),
    CONSTRAINT fk_securities_currency FOREIGN KEY(currency) REFERENCES currencies(code),
    CONSTRAINT fk_securities_securities_types FOREIGN KEY(asset_type) REFERENCES securities_types(id)
);

CREATE UNIQUE INDEX pk_securities_isin ON securities(isin);
CREATE TABLE operation_types (
    id smallint NOT NULL,
    name varchar(30) NOT NULL,
    title varchar(150) NOT NULL,
	CONSTRAINT pk_operation_types PRIMARY KEY (id)
);
CREATE SEQUENCE operations_id_seq;

CREATE TABLE operations (
    id integer DEFAULT pseudo_encrypt_24(CAST (nextval('operations_id_seq') AS integer)),
    pid integer NOT NULL,
    sid integer NOT NULL,
	time timestamp NOT NULL,
    op_id smallint NOT NULL,
	vol integer NOT NULL,
	price numeric(20,6) NOT NULL,
	CONSTRAINT pk_operations PRIMARY KEY (id),
    CONSTRAINT fk_operations_portfolios FOREIGN KEY(pid) REFERENCES portfolios(id),
    CONSTRAINT fk_operations_operation_types FOREIGN KEY(op_id) REFERENCES operation_types(id),
    CONSTRAINT fk_operations_securities FOREIGN KEY(sid) REFERENCES securities(id)
);

CREATE INDEX ix_operations ON operations(pid, sid, time);
CREATE TABLE prices (
    sid integer NOT NULL,
	date date NOT NULL,
	vol integer NOT NULL,
	price numeric(20,6) NOT NULL,
	CONSTRAINT pk_prices PRIMARY KEY (sid, date),
    CONSTRAINT fk_prices_securities FOREIGN KEY(sid) REFERENCES securities(id)
);
CREATE TABLE settings (
	settings jsonb NOT NULL
);
INSERT INTO currencies VALUES
    ('EUR','Евро'),
    ('USD','Доллар США'),
    ('RUB','Российский рубль');

INSERT INTO operation_types VALUES
    (1,'buy','Покупка'),
    (2,'sell','Продажа'),
    (3,'brokerageFee','Комиссия брокера'),
    (4,'exchangeFee','Комиссия биржи'),
    (5,'payIn','Ввод средств'),
    (6,'payOut','Вывод средств'),
    (7,'coupon','Выплата купона'),
    (8,'accruedInterestBuy','НКД при покупке'),
    (9,'accruedInterestSell','НКД при продаже'),
    (10,'buyback','Выкуп ценной бумаги');

INSERT INTO securities_types VALUES
    (10,'Stock','Акции'),
    (20,'Bond','Облигации'),
    (31,'EtfStock','ETF на акции'),
    (32,'EtfBond','ETF на облигации'),
    (34,'EtfMixed','Смешанный ETF'),
    (35,'EtfGold','ETF на золото'),
    (36,'EtfCurrency','ETF на аналог кэша'),
    (60,'Currency','Кэш');

INSERT INTO user_roles VALUES
    (1,'admin','Администратор'),
    (2,'user','Пользователь');

INSERT INTO sync_providers VALUES
    (1,'sber','Сбербанк'),
    (2,'tcs','Тинькофф'),
    (3,'vtb','ВТБ');

INSERT INTO settings (settings) VALUES (jsonb_build_object('ver', 1));

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO test;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO test;
GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO test;
