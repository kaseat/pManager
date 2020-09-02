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
