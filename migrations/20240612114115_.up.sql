BEGIN TRANSACTION;

CREATE TABLE wallet_address_nfts
(
    id                     BIGINT       NOT NULl PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    address                VARCHAR(255) NOT NULL,
    price                  BIGINT       NOT NULL,
    token_symbol           VARCHAR(255) NOT NULL,
    index                  BIGINT       NOT NULL,
    collection_address     VARCHAR(255) NOT NULL,
    collection_name        VARCHAR(255) NOT NULL,
    collection_description VARCHAR(255) NOT NULL,
    previews_urls          JSONB        NOT NULL,
    address_id             BIGINT       NOT NULL,
    CONSTRAINT fk_address_id FOREIGN KEY (address_id) REFERENCES wallet_addresses (id) ON DELETE CASCADE
);

COMMIT;