-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for all tables with updated_at column
CREATE TRIGGER update_tm_countries_updated_at BEFORE UPDATE ON tm_countries FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tm_provinces_updated_at BEFORE UPDATE ON tm_provinces FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tm_cities_updated_at BEFORE UPDATE ON tm_cities FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tm_districts_updated_at BEFORE UPDATE ON tm_districts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tm_villages_updated_at BEFORE UPDATE ON tm_villages FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tm_banks_updated_at BEFORE UPDATE ON tm_banks FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tm_currencies_updated_at BEFORE UPDATE ON tm_currencies FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tm_languages_updated_at BEFORE UPDATE ON tm_languages FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); 