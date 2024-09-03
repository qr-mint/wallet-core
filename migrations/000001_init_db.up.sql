BEGIN TRANSACTION;

CREATE TABLE users
(
    id         BIGINT                                         NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    type       VARCHAR(20) CHECK (users.type IN ('telegram')) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP            NOT NULL
);

CREATE TABLE mnemonics
(
    id   BIGINT       NOT NULl PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    hash TEXT         NOT NULL,
    name VARCHAR(255) NOT NULL
);
CREATE UNIQUE INDEX unique_mnemonics_hash_user_id ON mnemonics (hash);

CREATE TABLE users_mnemonics
(
    id          BIGINT NOT NULl PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id     BIGINT NOT NULL,
    mnemonic_id BIGINT NOT NULL,
    CONSTRAINT fk_mnemonic_id FOREIGN KEY (mnemonic_id) REFERENCES mnemonics (id) ON DELETE CASCADE,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX unique_users_mnemonics_user_id_mnemonic_id ON users_mnemonics (user_id, mnemonic_id);

CREATE TABLE telegram_users
(
    id                 BIGINT        NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    telegram_id        BIGINT        NOT NULL,
    telegram_bot_token VARCHAR(255)  NOT NULL,
    user_id            BIGINT UNIQUE NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX unique_telegram_users_telegram_id_telegram_bot_id ON telegram_users (telegram_id);

CREATE TABLE profiles
(
    id       BIGINT                                            NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    language VARCHAR(10) CHECK (language IN ('ru', 'en'))      NOT NULL,
    type     VARCHAR(20) CHECK (profiles.type IN ('telegram')) NOT NULL,
    user_id  BIGINT UNIQUE                                     NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE telegram_user_profiles
(
    id           BIGINT        NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    username     VARCHAR(255)  NOT NULL,
    first_name   VARCHAR(255) DEFAULT NULL,
    last_name    VARCHAR(255) DEFAULT NULL,
    image_source VARCHAR(255)  NOT NULL,
    profile_id   BIGINT UNIQUE NOT NULL,
    CONSTRAINT fk_profile_id FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE
);

CREATE TABLE coins
(
    id           BIGINT                     NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    network      VARCHAR(255)               NOT NULL CHECK (network IN ('ton', 'trc20')),
    name         VARCHAR(255)               NOT NULL CHECK (name IN ('ton', 'tether', 'tron')),
    symbol       VARCHAR(255)               NOT NULL,
    caption      VARCHAR(255)               NOT NULL,
    address      VARCHAR(255) DEFAULT NULL,
    decimals     INT                        NOT NULL,
    is_token     BOOLEAN      DEFAULT FALSE NOT NULL,
    image_source VARCHAR(255)               NOT NULL,
    is_default   BOOLEAN      DEFAULT TRUE  NOT NULL
);
CREATE UNIQUE INDEX unique_coins_coin_network ON coins (name, network);

INSERT INTO coins (network, name, symbol, caption, address, decimals, is_token, image_source, is_default)
VALUES ('ton', 'ton', 'TON', 'Toncoin', NULL, 9, FALSE, 'images/coins/ton.png', TRUE),
       ('trc20', 'tether', 'USDT', 'Tether USDt', 'TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t', 6, TRUE,
        'images/coins/usdt.png', TRUE),
       ('trc20', 'tron', 'TRON', 'TRON NETWORK', NULL, 6, FALSE, 'images/coins/tron.png', TRUE);

CREATE TABLE coin_prices
(
    id            BIGINT       NOT NULl PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    coin_id       BIGINT       NOT NULL,
    price         BIGINT       NOT NULL,
    fiat_currency VARCHAR(255) NOT NULL CHECK (fiat_currency IN ('usd', 'rub', 'eur')),
    date          TIMESTAMP    NOT NULL,
    CONSTRAINT fk_coin_id FOREIGN KEY (coin_id) REFERENCES coins (id)
);
CREATE UNIQUE INDEX unique_coin_prices_date_coin_id_fiat_currency ON coin_prices (date, coin_id, fiat_currency);

CREATE TABLE wallet_addresses
(
    id             BIGINT       NOT NULl PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    address        TEXT         NOT NULL,
    network        VARCHAR(255) NOT NULL CHECK (network IN ('ton', 'trc20')),
    mnemonic_id    BIGINT       NOT NULL,
    CONSTRAINT fk_mnemonic_id FOREIGN KEY (mnemonic_id) REFERENCES mnemonics (id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX unique_wallet_addresses_address_network ON wallet_addresses (address, network);
CREATE UNIQUE INDEX unique_wallet_addresses_mnemonic_id_network ON wallet_addresses (mnemonic_id, network);

CREATE TABLE wallet_address_coins
(
    id         BIGINT  NOT NULl PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    amount     BIGINT  NOT NULL,
    address    TEXT DEFAULT NULL,
    is_visible BOOLEAN NOT NULL,
    coin_id    BIGINT  NOT NULL,
    address_id BIGINT  NOT NULL,
    CONSTRAINT fk_coin_id FOREIGN KEY (coin_id) REFERENCES coins (id) ON DELETE CASCADE,
    CONSTRAINT fk_address_id FOREIGN KEY (address_id) REFERENCES wallet_addresses (id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX unique_wallet_coins_address_id_coin_id ON wallet_address_coins (address_id, coin_id);

CREATE TABLE notifications
(
    id         BIGINT                                                 NOT NULl PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    type       VARCHAR(20) CHECK (notifications.type IN ('telegram')) NOT NULL,
    expires_at TIMESTAMP                                              NOT NULL
);

CREATE TABLE notification_translations
(
    id              BIGINT                                       NOT NULl PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    language        VARCHAR(10) CHECK (language IN ('ru', 'en')) NOT NULL,
    text            TEXT                                         NOT NULL,
    image_path      TEXT DEFAULT NULL,
    notification_id BIGINT                                       NOT NULL,
    CONSTRAINT fk_notification_id FOREIGN KEY (notification_id) REFERENCES notifications (id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX unique_notification_translations_notification_id_language ON notification_translations (notification_id, language);

CREATE TABLE processed_notifications
(
    id              BIGINT NOT NULl PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id         BIGINT NOT NULL,
    notification_id BIGINT NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_notification_id FOREIGN KEY (notification_id) REFERENCES notifications (id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX unique_processed_notifications_user_id_notification_id ON processed_notifications (user_id, notification_id);

INSERT INTO coin_prices (coin_id, price, fiat_currency, date)
VALUES ((SELECT id FROM coins WHERE name = 'ton'), 653573810610, 'rub', now()),
       ((SELECT id FROM coins WHERE name = 'ton'), 7426800000, 'usd', now()),
       ((SELECT id FROM coins WHERE name = 'ton'), 6883589077, 'eur', now()),
       ((SELECT id FROM coins WHERE name = 'tether'), 88416813, 'rub', now()),
       ((SELECT id FROM coins WHERE name = 'tether'), 999138, 'usd', now()),
       ((SELECT id FROM coins WHERE name = 'tether'), 925985, 'eur', now()),
       ((SELECT id FROM coins WHERE name = 'tron'), 11371216, 'rub', now()),
       ((SELECT id FROM coins WHERE name = 'tron'), 128504, 'usd', now()),
       ((SELECT id FROM coins WHERE name = 'tron'), 119092, 'eur', now());

COMMIT;