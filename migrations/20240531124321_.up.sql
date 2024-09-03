BEGIN TRANSACTION;

CREATE TABLE blockchain_transactions
(
    id           BIGINT                              NOT NULl PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    hash         VARCHAR(255)                        NOT NULL,
    amount       BIGINT                              NOT NULL,
    address      VARCHAR(255)                        NOT NULL,
    address_to   VARCHAR(255)                        NOT NULL,
    address_from VARCHAR(255)                        NOT NULL,
    status       VARCHAR(255)                        NOT NULL CHECK (blockchain_transactions.status IN ('confirmed', 'failed')),
    type         VARCHAR(255)                        NOT NULL CHECK ( type IN ('in', 'out') ),
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    coin_id      BIGINT                              NOT NULL,
    CONSTRAINT fk_coin_id FOREIGN KEY (coin_id) REFERENCES coins (id)
);
CREATE UNIQUE INDEX unique_blockchain_transactions_hash_coin_id ON blockchain_transactions (hash, coin_id);

COMMIT;