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
CREATE TABLE IF NOT EXISTS "tm_geodirectories" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" VARCHAR(255) NOT NULL,
    "type" geo_type NOT NULL,
    "code" VARCHAR(10) DEFAULT NULL,
    "postal_code" VARCHAR(255) DEFAULT NULL,
    "longitude" VARCHAR(255) DEFAULT NULL,
    "latitude" VARCHAR(255) DEFAULT NULL,
    "record_left" INTEGER NULL,
    "record_right" INTEGER NULL,
    "record_ordering" INTEGER NULL ,
    "parent_id" UUID DEFAULT NULL,
    "created_at" TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL,
    "updated_at" TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL
);

-- Create the indexes, as PostgreSQL uses CREATE INDEX instead of inline KEY statements
CREATE INDEX IF NOT EXISTS tm_geodirectories_record_left_index ON "tm_geodirectories" ("record_left");
CREATE INDEX IF NOT EXISTS tm_geodirectories_record_right_index ON "tm_geodirectories" ("record_right");
CREATE INDEX IF NOT EXISTS tm_geodirectories_parent_id_index ON "tm_geodirectories" ("parent_id");
CREATE INDEX IF NOT EXISTS tm_geodirectories_record_ordering_index ON "tm_geodirectories" ("record_ordering");
CREATE INDEX IF NOT EXISTS tm_geodirectories_name_index ON "tm_geodirectories" ("name");
CREATE INDEX IF NOT EXISTS tm_geodirectories_code_index ON "tm_geodirectories" ("code");
CREATE INDEX IF NOT EXISTS tm_geodirectories_postal_code_index ON "tm_geodirectories" ("postal_code");
CREATE INDEX IF NOT EXISTS tm_geodirectories_longitude_index ON "tm_geodirectories" ("longitude");
CREATE INDEX IF NOT EXISTS tm_geodirectories_latitude_index ON "tm_geodirectories" ("latitude");
