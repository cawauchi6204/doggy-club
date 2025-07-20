-- Database initialization script for DoggyClub
-- This script sets up the production database with proper permissions and extensions

-- Create extensions if they don't exist
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS "btree_gin";

-- Create custom types for the new schema
DO $$ BEGIN
    CREATE TYPE visibility AS ENUM ('public', 'private');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE detection_method AS ENUM ('gps', 'bluetooth');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE device_type AS ENUM ('ios', 'android');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE notification_type AS ENUM ('encounter', 'gift', 'like');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE subscription_status AS ENUM ('active', 'canceled');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Create application user if it doesn't exist
DO $$ BEGIN
    CREATE ROLE doggyclub_app WITH LOGIN PASSWORD 'app_password_change_me';
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Grant necessary permissions
GRANT CONNECT ON DATABASE doggyclub TO doggyclub_app;
GRANT USAGE ON SCHEMA public TO doggyclub_app;
GRANT CREATE ON SCHEMA public TO doggyclub_app;

-- Grant permissions on existing tables (will be created by GORM migrations)
-- These permissions will be applied to future tables as well
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO doggyclub_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE, SELECT ON SEQUENCES TO doggyclub_app;

-- Performance indexes for the new schema
-- These will be created by GORM as well, but we define them here for reference

-- User indexes
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_email ON users(email);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_username ON users(username);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_visibility ON users(visibility);

-- Dog indexes
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_dogs_user_id ON dogs(user_id);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_dogs_breed ON dogs(breed);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_dogs_created_at ON dogs(created_at);

-- Encounter indexes
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_encounters_dog1_id ON encounters(dog1_id);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_encounters_dog2_id ON encounters(dog2_id);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_encounters_timestamp ON encounters(timestamp);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_encounters_detection_method ON encounters(detection_method);

-- Device location indexes with GiST for geographic queries
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_device_locations_dog_id ON device_locations(dog_id);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_device_locations_updated_at ON device_locations(updated_at);

-- Post indexes
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_posts_dog_id ON posts(dog_id);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_posts_created_at ON posts(created_at);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_posts_content_gin ON posts USING GIN(to_tsvector('english', content));

-- Like indexes
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_likes_post_id ON likes(post_id);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_likes_dog_id ON likes(dog_id);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_likes_created_at ON likes(created_at);

-- Comment indexes
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_comments_post_id ON comments(post_id);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_comments_dog_id ON comments(dog_id);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_comments_created_at ON comments(created_at);

-- Hashtag indexes
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_hashtags_tag ON hashtags(tag);

-- Follower indexes
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_followers_follower_dog_id ON followers(follower_dog_id);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_followers_followed_dog_id ON followers(followed_dog_id);

-- Gift indexes
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_gifts_sender_dog_id ON gifts(sender_dog_id);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_gifts_receiver_dog_id ON gifts(receiver_dog_id);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_gifts_sent_at ON gifts(sent_at);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_gifts_gift_type ON gifts(gift_type);

-- Subscription indexes
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_subscription_plans_name ON subscription_plans(name);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_user_subscriptions_user_id ON user_subscriptions(user_id);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_user_subscriptions_plan_id ON user_subscriptions(plan_id);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_user_subscriptions_status ON user_subscriptions(status);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_user_subscriptions_end_date ON user_subscriptions(end_date);

-- Device token indexes
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_device_tokens_user_id ON device_tokens(user_id);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_device_tokens_device_type ON device_tokens(device_type);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_device_tokens_last_active ON device_tokens(last_active);

-- Notification indexes
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_notifications_type ON notifications(type);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_notifications_sent_at ON notifications(sent_at);

-- Create the GiST index for location-based queries
-- This is critical for performance of encounter detection
CREATE OR REPLACE FUNCTION create_location_indexes() RETURNS void AS $$
BEGIN
    -- Check if the device_locations table exists before creating index
    IF EXISTS (SELECT FROM pg_tables WHERE tablename = 'device_locations') THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_device_locations_location_gist ON device_locations USING GIST(location)';
    END IF;
    
    -- Same for encounters table
    IF EXISTS (SELECT FROM pg_tables WHERE tablename = 'encounters') THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_encounters_location_gist ON encounters USING GIST(location)';
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Create function for full-text search on posts
CREATE OR REPLACE FUNCTION create_fulltext_indexes() RETURNS void AS $$
BEGIN
    IF EXISTS (SELECT FROM pg_tables WHERE tablename = 'posts') THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_posts_content_fulltext ON posts USING GIN(to_tsvector(''english'', content))';
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Create function for updating updated_at timestamp (not used in new schema but kept for reference)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Database settings for performance
-- These settings should also be configured in postgresql.conf
-- ALTER SYSTEM SET shared_preload_libraries = 'pg_stat_statements';
-- ALTER SYSTEM SET track_activity_query_size = 2048;
-- ALTER SYSTEM SET pg_stat_statements.track = all;

-- Partitioning setup for large tables (posts)
-- CREATE TABLE posts_y2024m01 PARTITION OF posts
-- FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

-- Materialized views for analytics (will be created after tables exist)
/*
CREATE MATERIALIZED VIEW IF NOT EXISTS daily_dog_stats AS
SELECT 
    date_trunc('day', created_at) as date,
    count(*) as new_dogs,
    count(DISTINCT user_id) as new_dog_owners
FROM dogs 
GROUP BY date_trunc('day', created_at)
ORDER BY date;

CREATE INDEX IF NOT EXISTS idx_daily_dog_stats_date ON daily_dog_stats(date);

CREATE MATERIALIZED VIEW IF NOT EXISTS daily_encounter_stats AS
SELECT 
    date_trunc('day', timestamp) as date,
    count(*) as total_encounters,
    count(DISTINCT dog1_id) + count(DISTINCT dog2_id) as active_dogs,
    detection_method,
    count(*) FILTER (WHERE detection_method = 'gps') as gps_encounters,
    count(*) FILTER (WHERE detection_method = 'bluetooth') as bluetooth_encounters
FROM encounters 
GROUP BY date_trunc('day', timestamp), detection_method
ORDER BY date;

CREATE INDEX IF NOT EXISTS idx_daily_encounter_stats_date ON daily_encounter_stats(date);
*/

-- Grant access to materialized views
-- GRANT SELECT ON daily_dog_stats TO doggyclub_app;
-- GRANT SELECT ON daily_encounter_stats TO doggyclub_app;

COMMIT;