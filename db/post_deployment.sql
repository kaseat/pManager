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
