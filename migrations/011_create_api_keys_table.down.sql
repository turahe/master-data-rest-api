-- Drop API keys table and related objects
DROP TRIGGER IF EXISTS update_tm_api_keys_updated_at ON tm_api_keys;
DROP INDEX IF EXISTS idx_tm_api_keys_deleted_at;
DROP INDEX IF EXISTS idx_tm_api_keys_expires_at;
DROP INDEX IF EXISTS idx_tm_api_keys_is_active;
DROP INDEX IF EXISTS idx_tm_api_keys_key;
DROP TABLE IF EXISTS tm_api_keys;
DROP FUNCTION IF EXISTS update_updated_at_column();
