BEGIN TRANSACTION;

CREATE TABLE exchanges
(
    id                   BIGINT       NOT NULl PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    external_id          VARCHAR(255) NOT NULL,
    support_link         VARCHAR(255) NOT NULL,
    address_coin_id_from BIGINT       NOT NULL,
    address_coin_id_to   BIGINT       NOT NULL,
    mnemonic_id          BIGINT       NOT NULL,
    CONSTRAINT fk_address_coin_id_from FOREIGN KEY (address_coin_id_from) REFERENCES wallet_address_coins (id) ON DELETE CASCADE,
    CONSTRAINT fk_address_coin_id_to FOREIGN KEY (address_coin_id_to) REFERENCES wallet_address_coins (id) ON DELETE CASCADE,
    CONSTRAINT fk_mnemonic_id FOREIGN KEY (mnemonic_id) REFERENCES mnemonics (id) ON DELETE CASCADE
);

COMMIT;