BEGIN TRANSACTION;

DROP TABLE notifications CASCADE;
DROP TABLE notification_translations CASCADE;
DROP TABLE processed_notifications CASCADE;
DROP TABLE telegram_user_profiles CASCADE;
DROP TABLE profiles CASCADE;
DROP TABLE wallet_address_coins CASCADE;
DROP TABLE wallet_addresses CASCADE;
DROP TABLE users_mnemonics CASCADE;
DROP TABLE telegram_users CASCADE;
DROP TABLE users CASCADE;
DROP TABLE mnemonics CASCADE;
DROP TABLE coin_prices CASCADE;
DROP TABLE coins CASCADE;

COMMIT;