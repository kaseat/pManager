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
    (31,'ErfStock','ETF на акции'),
    (32,'EtfBond','ETF на облигации'),
    (34,'EtfMixed','Смешанный ETF'),
    (35,'EtfGold','ETF на золото'),
    (36,'EtfCash','ETF на аналог кэша');