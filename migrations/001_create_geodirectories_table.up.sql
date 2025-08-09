-- First, create the custom ENUM types required by PostgreSQL
-- Only create if it doesn't exist to avoid conflicts
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'geo_type') THEN
        CREATE TYPE geo_type AS ENUM (
            'CONTINENT',
            'SUBCONTINENT',
            'COUNTRY',
            'STATE',
            'PROVINCE',
            'REGENCY',
            'CITY',
            'DISTRICT',
            'SUBDISTRICT',
            'VILLAGE'
        );
    END IF;
END
$$;

-- Now, create the main table
-- Main table for storing hierarchical geodirectory data
CREATE TABLE IF NOT EXISTS "tm_geodirectories" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for each geodirectory entry
    "name" VARCHAR(255) NOT NULL,                    -- Name of the location
    "type" geo_type NOT NULL,                        -- Type of location (e.g., COUNTRY, CITY, etc.)
    "code" VARCHAR(10) DEFAULT NULL,                 -- Optional code for the location
    "postal_code" VARCHAR(255) DEFAULT NULL,         -- Optional postal code
    "longitude" VARCHAR(255) DEFAULT NULL,           -- Optional longitude coordinate
    "latitude" VARCHAR(255) DEFAULT NULL,            -- Optional latitude coordinate
    "record_left" INTEGER NULL,                      -- Nested set model: left value
    "record_right" INTEGER NULL,                     -- Nested set model: right value
    "record_ordering" INTEGER,                  -- Optional order data geodirectory in the same level
    "record_depth" INTEGER DEFAULT 0,                -- Depth in the hierarchy (root = 0)
    "parent_id" UUID DEFAULT NULL,                   -- Reference to parent geodirectory
    "created_at" TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL, -- Creation timestamp
    "updated_at" TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL  -- Last update timestamp
);

-- Create the indexes, as PostgreSQL uses CREATE INDEX instead of inline KEY statements
CREATE INDEX IF NOT EXISTS tm_geodirectories_record_left_index ON "tm_geodirectories" ("record_left");
CREATE INDEX IF NOT EXISTS tm_geodirectories_record_right_index ON "tm_geodirectories" ("record_right");
CREATE INDEX IF NOT EXISTS tm_geodirectories_parent_id_index ON "tm_geodirectories" ("parent_id");
CREATE INDEX IF NOT EXISTS tm_geodirectories_record_ordering_index ON "tm_geodirectories" ("record_ordering");
CREATE INDEX IF NOT EXISTS tm_geodirectories_record_depth_index ON "tm_geodirectories" ("record_depth");
CREATE INDEX IF NOT EXISTS tm_geodirectories_name_index ON "tm_geodirectories" ("name");
CREATE INDEX IF NOT EXISTS tm_geodirectories_code_index ON "tm_geodirectories" ("code");
CREATE INDEX IF NOT EXISTS tm_geodirectories_postal_code_index ON "tm_geodirectories" ("postal_code");
CREATE INDEX IF NOT EXISTS tm_geodirectories_longitude_index ON "tm_geodirectories" ("longitude");
CREATE INDEX IF NOT EXISTS tm_geodirectories_latitude_index ON "tm_geodirectories" ("latitude");
