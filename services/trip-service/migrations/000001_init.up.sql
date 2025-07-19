-- Create trips table
CREATE TABLE trips (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create ride_fares table
CREATE TABLE ride_fares (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    package_slug VARCHAR(50) NOT NULL,
    total_price_in_cents DECIMAL(12,2) NOT NULL,
    route JSONB -- Store OSRM API response as JSON
);

-- Create indexes for better performance
CREATE INDEX idx_trips_user_id ON trips(user_id);
CREATE INDEX idx_trips_status ON trips(status);
CREATE INDEX idx_ride_fares_user_id ON ride_fares(user_id);

-- Insert some sample data for testing
INSERT INTO trips (user_id, status) VALUES 
    ('user123', 'pending'),
    ('user456', 'completed'),
    ('user789', 'cancelled');

INSERT INTO ride_fares (user_id, package_slug, total_price_in_cents, route) VALUES 
    ('user123', 'sedan', 1250.00, '{}'),
    ('user456', 'suv', 1750.00, '{}'),
    ('user789', 'luxury', 3500.00, '{}');
