-- Create the trigger function for updating updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create API keys table
CREATE TABLE IF NOT EXISTS tm_api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    key VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    expires_at TIMESTAMP WITH TIME ZONE,
    last_used_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_tm_api_keys_key ON tm_api_keys(key) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tm_api_keys_is_active ON tm_api_keys(is_active) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tm_api_keys_expires_at ON tm_api_keys(expires_at) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tm_api_keys_deleted_at ON tm_api_keys(deleted_at);

-- Add trigger for updating updated_at timestamp (if not exists)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.triggers WHERE trigger_name = 'update_tm_api_keys_updated_at') THEN
        CREATE TRIGGER update_tm_api_keys_updated_at
            BEFORE UPDATE ON tm_api_keys
            FOR EACH ROW
            EXECUTE FUNCTION update_updated_at_column();
    END IF;
END
$$;
