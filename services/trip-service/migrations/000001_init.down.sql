-- Drop indexes first
DROP INDEX IF EXISTS idx_ride_fares_user_id;
DROP INDEX IF EXISTS idx_trips_status;
DROP INDEX IF EXISTS idx_trips_user_id;

-- Drop tables
DROP TABLE IF EXISTS ride_fares;
DROP TABLE IF EXISTS trips;