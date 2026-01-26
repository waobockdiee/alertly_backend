package cjbot_creator

import (
	"alertly/internal/common"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	BOT_USER_ID = 16 // System Bot account ID - "Alertly Official" (official@alertly.ca)
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// ============================================
// HASH DEDUPLICATION METHODS
// ============================================

// GenerateIncidentHash creates SHA256 hash for deduplication
// Uses only source + externalID (NOT timestamp) to prevent duplicates
// when TFS/TPS updates dispatch_time for the same incident
func GenerateIncidentHash(source, externalID string, timestamp time.Time) string {
	// ✅ FIX: Remove timestamp from hash to prevent duplicate reports
	// TFS/TPS use unique IDs that don't change during incident lifetime
	data := fmt.Sprintf("%s:%s", source, externalID)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// CheckHashExists verifies if incident was already processed
func (r *Repository) CheckHashExists(hash string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM bot_incident_hashes WHERE hash = $1)`
	err := r.db.QueryRow(query, hash).Scan(&exists)
	return exists, err
}

// SaveIncidentHash stores hash to prevent duplicates
func (r *Repository) SaveIncidentHash(hash BotIncidentHash) error {
	query := `
		INSERT INTO bot_incident_hashes (hash, source, external_id, category_code, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (hash) DO NOTHING
	`
	_, err := r.db.Exec(query, hash.Hash, hash.Source, hash.ExternalID, hash.CategoryCode, hash.CreatedAt, hash.ExpiresAt)
	return err
}

// CleanupExpiredHashes removes old hashes (called by separate cleanup cronjob)
func (r *Repository) CleanupExpiredHashes() (int64, error) {
	query := `DELETE FROM bot_incident_hashes WHERE expires_at IS NOT NULL AND expires_at < NOW()`
	result, err := r.db.Exec(query)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// ============================================
// GEOCODING CACHE METHODS
// ============================================

// NormalizeAddress cleans address string for consistent caching
func NormalizeAddress(address string) string {
	addr := strings.ToUpper(strings.TrimSpace(address))
	// Remove multiple spaces
	addr = strings.Join(strings.Fields(addr), " ")
	// Add "TORONTO, ON" if not present
	if !strings.Contains(addr, "TORONTO") {
		addr = addr + ", TORONTO, ON"
	}
	return addr
}

// GenerateAddressHash creates SHA256 hash of normalized address
func GenerateAddressHash(normalizedAddress string) string {
	hash := sha256.Sum256([]byte(normalizedAddress))
	return fmt.Sprintf("%x", hash)
}

// GetCachedGeocode retrieves cached coordinates for an address
func (r *Repository) GetCachedGeocode(addressHash string) (*GeocodingCache, error) {
	var cache GeocodingCache
	query := `
		SELECT address_hash, original_address, normalized_address, latitude, longitude, source, created_at, last_used_at
		FROM geocoding_cache
		WHERE address_hash = $1
	`
	err := r.db.QueryRow(query, addressHash).Scan(
		&cache.AddressHash,
		&cache.OriginalAddress,
		&cache.NormalizedAddress,
		&cache.Latitude,
		&cache.Longitude,
		&cache.Source,
		&cache.CreatedAt,
		&cache.LastUsedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Update last_used_at timestamp
	go r.touchGeocodeCache(addressHash)

	return &cache, nil
}

// SaveGeocodeCache stores geocoding result
func (r *Repository) SaveGeocodeCache(cache GeocodingCache) error {
	query := `
		INSERT INTO geocoding_cache (address_hash, original_address, normalized_address, latitude, longitude, source, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (address_hash) DO UPDATE SET last_used_at = NOW()
	`
	_, err := r.db.Exec(query,
		cache.AddressHash,
		cache.OriginalAddress,
		cache.NormalizedAddress,
		cache.Latitude,
		cache.Longitude,
		cache.Source,
		cache.CreatedAt,
	)
	return err
}

// touchGeocodeCache updates last_used_at timestamp (async)
func (r *Repository) touchGeocodeCache(addressHash string) {
	query := `UPDATE geocoding_cache SET last_used_at = NOW() WHERE address_hash = $1`
	_, err := r.db.Exec(query, addressHash)
	if err != nil {
		log.Printf("WARNING: Failed to update geocode cache timestamp: %v", err)
	}
}

// CleanupOldGeocodeCache removes unused cache entries older than 30 days
func (r *Repository) CleanupOldGeocodeCache() (int64, error) {
	query := `DELETE FROM geocoding_cache WHERE last_used_at < NOW() - INTERVAL '30 days'`
	result, err := r.db.Exec(query)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// ============================================
// INCIDENT PERSISTENCE METHODS
// ============================================

// GetSubcategoryID retrieves the insu_id for a given subcategory code
func (r *Repository) GetSubcategoryID(subcategoryCode string) (int64, error) {
	var insuID int64
	query := `SELECT insu_id FROM incident_subcategories WHERE code = $1 LIMIT 1`
	err := r.db.QueryRow(query, subcategoryCode).Scan(&insuID)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("subcategory not found: %s", subcategoryCode)
	}
	return insuID, err
}

// SaveBotIncident inserts a new incident_report from bot
// This method handles the complete flow: find/create cluster + create report
func (r *Repository) SaveBotIncident(incident NormalizedIncident) (int64, error) {
	// Get subcategory ID
	insuID, err := r.GetSubcategoryID(incident.SubcategoryCode)
	if err != nil {
		return 0, fmt.Errorf("getting subcategory ID: %w", err)
	}

	// Step 1: Try to find existing cluster nearby
	radius := r.getClusteringRadius(incident.CategoryCode)
	clusterID, err := r.CheckClusterExists(incident.CategoryCode, incident.Latitude, incident.Longitude, int(radius), 24)
	if err != nil {
		return 0, fmt.Errorf("checking cluster: %w", err)
	}

	// Step 2: If no cluster found, create one
	isNewCluster := clusterID == 0
	if isNewCluster {
		clusterID, err = r.CreateCluster(incident, insuID)
		if err != nil {
			return 0, fmt.Errorf("creating cluster: %w", err)
		}
	}

	// Step 3: Insert incident_report with cluster ID
	query := `
		INSERT INTO incident_reports (
			account_id,
			incl_id,
			insu_id,
			latitude,
			longitude,
			media_url,
			description,
			address,
			city,
			province,
			postal_code,
			subcategory_name,
			subcategory_code,
			category_code,
			event_type,
			vote,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, 1, NOW()) RETURNING inre_id
	`

	var reportID int64
	err = r.db.QueryRow(query,
		BOT_USER_ID,
		clusterID,
		insuID,
		incident.Latitude,
		incident.Longitude,
		incident.ImageURL,
		incident.Description,
		incident.Address,
		incident.City,
		incident.Province,
		incident.PostalCode,
		incident.SubcategoryCode, // Using code as name
		incident.SubcategoryCode,
		incident.CategoryCode,
		incident.EventType,
	).Scan(&reportID)
	if err != nil {
		return 0, err
	}

	// 🔔 Create notification when joining existing cluster (async)
	if !isNewCluster {
		go func(accountID int64, inclID int64, reportID int64) {
			if err := common.SaveNotification(r.db, "new_incident_cluster", accountID, inclID); err != nil {
				log.Printf("⚠️ [Bot] Error saving notification for cluster update %d: %v\n", inclID, err)
			} else {
				log.Printf("✅ [Bot] Notification created for cluster update %d (report %d)\n", inclID, reportID)
			}
		}(BOT_USER_ID, clusterID, reportID)
	}

	return reportID, nil
}

// getClusteringRadius returns appropriate radius for category (in meters)
func (r *Repository) getClusteringRadius(categoryCode string) float64 {
	switch categoryCode {
	case "traffic_accident":
		return 100 // 100m for traffic accidents
	case "crime":
		return 200 // 200m for crime
	case "fire_incident":
		return 150 // 150m for fire incidents
	case "medical_emergency":
		return 150 // 150m for medical
	case "infrastructure_issues":
		return 500 // 500m for infrastructure
	case "extreme_weather":
		return 10000 // 10km for weather
	default:
		return 200 // Default 200m
	}
}

// CheckClusterExists verifies if a cluster exists for this category and location
func (r *Repository) CheckClusterExists(categoryCode string, lat, lng float64, radiusMeters int, hoursBack int) (int64, error) {
	var clusterID int64
	query := `
		SELECT incl_id FROM incident_clusters
		WHERE category_code = $1
		AND ST_DWithin(center_location, ST_MakePoint($2, $3)::geography, $4)
		AND created_at >= NOW() - INTERVAL '1 hour' * $5
		AND is_active = '1'
		ORDER BY created_at DESC
		LIMIT 1
	`
	err := r.db.QueryRow(query, categoryCode, lng, lat, radiusMeters, hoursBack).Scan(&clusterID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return clusterID, err
}

// CreateCluster creates a new incident cluster
func (r *Repository) CreateCluster(incident NormalizedIncident, insuID int64) (int64, error) {
	// Get duration for subcategory (default 24h)
	duration := r.getSubcategoryDuration(incident.SubcategoryCode)
	endTime := time.Now().Add(time.Duration(duration) * time.Hour)

	// ✅ FIX: Explicit type casts for coordinates
	// $3::decimal and $4::decimal for DECIMAL columns, then ::float8 for ST_MakePoint
	// Prevents "inconsistent types deduced for parameter" PostgreSQL error
	query := `
		INSERT INTO incident_clusters (
			account_id,
			insu_id,
			center_latitude,
			center_longitude,
			center_location,
			media_url,
			media_type,
			event_type,
			description,
			address,
			city,
			province,
			postal_code,
			subcategory_name,
			subcategory_code,
			category_code,
			is_active,
			created_at,
			start_time,
			end_time
		) VALUES ($1, $2, $3::decimal, $4::decimal, ST_SetSRID(ST_MakePoint($4::float8, $3::float8), 4326)::geography, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, 1, NOW(), NOW(), $16) RETURNING incl_id
	`

	var clusterID int64
	err := r.db.QueryRow(query,
		BOT_USER_ID,
		insuID,
		incident.Latitude,
		incident.Longitude,
		incident.ImageURL,
		"image", // Always image for bot reports
		incident.EventType,
		incident.Description,
		incident.Address,
		incident.City,
		incident.Province,
		incident.PostalCode,
		incident.SubcategoryCode, // Using code as name for now
		incident.SubcategoryCode,
		incident.CategoryCode,
		endTime,
	).Scan(&clusterID)
	if err != nil {
		return 0, err
	}

	// 🔔 Create notification for new cluster (async to not block bot processing)
	go func(accountID int64, clusterID int64) {
		if err := common.SaveNotification(r.db, "new_cluster", accountID, clusterID); err != nil {
			log.Printf("⚠️ [Bot] Error saving notification for cluster %d: %v\n", clusterID, err)
		} else {
			log.Printf("✅ [Bot] Notification created for cluster %d\n", clusterID)
		}
	}(BOT_USER_ID, clusterID)

	return clusterID, nil
}

// getSubcategoryDuration retrieves duration for a subcategory (defaults to 24h)
func (r *Repository) getSubcategoryDuration(subcategoryCode string) int {
	var duration int
	query := `SELECT default_duration_hours FROM incident_subcategories WHERE code = $1`
	err := r.db.QueryRow(query, subcategoryCode).Scan(&duration)
	if err != nil || duration < 24 {
		return 24 // Default minimum
	}
	return duration
}
