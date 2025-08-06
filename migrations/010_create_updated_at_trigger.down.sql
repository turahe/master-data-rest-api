-- Drop triggers
DROP TRIGGER IF EXISTS update_tm_countries_updated_at ON tm_countries;
DROP TRIGGER IF EXISTS update_tm_provinces_updated_at ON tm_provinces;
DROP TRIGGER IF EXISTS update_tm_cities_updated_at ON tm_cities;
DROP TRIGGER IF EXISTS update_tm_districts_updated_at ON tm_districts;
DROP TRIGGER IF EXISTS update_tm_villages_updated_at ON tm_villages;
DROP TRIGGER IF EXISTS update_tm_banks_updated_at ON tm_banks;
DROP TRIGGER IF EXISTS update_tm_currencies_updated_at ON tm_currencies;
DROP TRIGGER IF EXISTS update_tm_languages_updated_at ON tm_languages;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column(); 