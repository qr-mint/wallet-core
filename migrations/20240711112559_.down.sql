BEGIN TRANSACTION;

ALTER TABLE wallet_address_nfts
    ALTER COLUMN collection_description TYPE VARCHAR(255),
    ALTER COLUMN collection_name TYPE VARCHAR(255);

COMMIT;