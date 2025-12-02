-- =====================================================
-- Bot Seeder Tables Migration
-- Description: Tables for data-seeder cronjob system
-- =====================================================

-- Table for deduplication of bot-scraped incidents
CREATE TABLE IF NOT EXISTS bot_incident_hashes (
    hash VARCHAR(64) PRIMARY KEY COMMENT 'SHA256 hash of source+external_id+timestamp',
    source VARCHAR(50) NOT NULL COMMENT 'Source system: tps, tfs, ttc, hydro, weather',
    external_id VARCHAR(255) NOT NULL COMMENT 'ID from source system',
    category_code VARCHAR(50) COMMENT 'Alertly category code',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NULL COMMENT 'TTL for automatic cleanup',

    INDEX idx_source_external (source, external_id),
    INDEX idx_expires (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Deduplication table for bot-generated incidents';

-- Table for geocoding cache to avoid hitting Nominatim repeatedly
CREATE TABLE IF NOT EXISTS geocoding_cache (
    address_hash VARCHAR(64) PRIMARY KEY COMMENT 'SHA256 hash of normalized address',
    original_address VARCHAR(500) NOT NULL COMMENT 'Original address string from source',
    normalized_address VARCHAR(500) NOT NULL COMMENT 'Cleaned address for geocoding',
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    source VARCHAR(50) DEFAULT 'nominatim' COMMENT 'Geocoding provider used',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_created (created_at),
    INDEX idx_last_used (last_used_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Cache for geocoded addresses to reduce API calls';

-- Add index on incident_reports for bot user tracking (optional, for analytics)
-- This assumes user_id=1 is the System Bot account
CREATE INDEX IF NOT EXISTS idx_bot_reports ON incident_reports(user_id, created_at)
WHERE user_id = 1;
