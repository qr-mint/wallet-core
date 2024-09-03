BEGIN TRANSACTION;

ALTER TABLE wallet_address_nfts
    ALTER COLUMN collection_description TYPE TEXT,
    ALTER COLUMN collection_name TYPE TEXT;

COMMIT;