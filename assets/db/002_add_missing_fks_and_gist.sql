-- ============================================================
-- Alertly: Add Missing Foreign Keys & GiST Spatial Indexes
-- Post-migration MySQL → PostgreSQL
-- Date: 2026-01-23
-- Status: EXECUTED ✅
-- ============================================================

-- ============================================================
-- PART 1: Foreign Keys (EXECUTED)
-- ============================================================

-- 1. incident_clusters.insu_id → incident_subcategories.insu_id ✅
-- ALTER TABLE incident_clusters
-- ADD CONSTRAINT fk_clusters_subcategory
-- FOREIGN KEY (insu_id) REFERENCES incident_subcategories(insu_id)
-- ON DELETE RESTRICT ON UPDATE CASCADE;

-- 2. notification_deliveries.noti_id → notifications.noti_id ✅
-- ALTER TABLE notification_deliveries
-- ADD CONSTRAINT fk_deliveries_notification
-- FOREIGN KEY (noti_id) REFERENCES notifications(noti_id)
-- ON DELETE CASCADE ON UPDATE CASCADE;

-- 3. account_cluster_saved.incl_id → incident_clusters.incl_id ✅
-- ALTER TABLE account_cluster_saved
-- ADD CONSTRAINT fk_saved_cluster
-- FOREIGN KEY (incl_id) REFERENCES incident_clusters(incl_id)
-- ON DELETE CASCADE ON UPDATE CASCADE;

-- 4. account_history.incl_id → incident_clusters.incl_id ✅
-- ALTER TABLE account_history
-- ADD CONSTRAINT fk_history_cluster
-- FOREIGN KEY (incl_id) REFERENCES incident_clusters(incl_id)
-- ON DELETE SET NULL ON UPDATE CASCADE;

-- ============================================================
-- PART 2: PostGIS Spatial Columns & GiST Indexes (EXECUTED)
-- ============================================================

-- incident_clusters.center_location ✅
-- ALTER TABLE incident_clusters
-- ADD COLUMN center_location GEOGRAPHY(POINT, 4326);

-- UPDATE incident_clusters
-- SET center_location = ST_SetSRID(ST_MakePoint(center_longitude, center_latitude), 4326)::geography
-- WHERE center_location IS NULL;
-- Result: 1059 rows updated

-- CREATE INDEX idx_clusters_location_gist ON incident_clusters USING GIST (center_location); ✅

-- account_favorite_locations.location ✅
-- ALTER TABLE account_favorite_locations
-- ADD COLUMN location GEOGRAPHY(POINT, 4326);

-- UPDATE account_favorite_locations
-- SET location = ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)::geography
-- WHERE location IS NULL;
-- Result: 7 rows updated

-- CREATE INDEX idx_favorite_locations_gist ON account_favorite_locations USING GIST (location); ✅

-- ============================================================
-- PART 3: Go Code Changes (UPDATED)
-- ============================================================
--
-- Files modified to use ST_DWithin with GiST indexes:
--
-- 1. internal/newincident/repository.go
--    - CheckAndGetIfClusterExist: ST_DistanceSphere → ST_DWithin
--    - SaveCluster: Added center_location column
--    - UpdateClusterAsTrue/AsFalse: Added center_location update
--    - UpdateClusterLocation: Added center_location update
--
-- 2. internal/getclusterbyradius/repository.go
--    - GetClustersByRadius: ST_DistanceSphere → ST_DWithin
--
-- 3. internal/cronjobs/cjnewcluster/repository.go
--    - GetDeviceTokensForNewCluster: ST_DistanceSphere → ST_DWithin
--
-- 4. internal/getincidentsasreels/repository.go
--    - GetReel: ST_DistanceSphere → ST_DWithin
--
-- 5. internal/cronjobs/cjbot_creator/repository.go
--    - CheckClusterExists: ST_DistanceSphere → ST_DWithin
--    - CreateCluster: Added center_location column
--
-- ============================================================
-- VERIFICATION QUERIES
-- ============================================================

-- Verify GiST indexes exist:
-- SELECT indexname, indexdef FROM pg_indexes
-- WHERE indexname LIKE '%gist%';

-- Verify geography columns:
-- SELECT column_name, data_type
-- FROM information_schema.columns
-- WHERE table_name IN ('incident_clusters', 'account_favorite_locations')
-- AND data_type = 'USER-DEFINED';

-- Verify all data migrated:
-- SELECT COUNT(*), COUNT(center_location) FROM incident_clusters;
-- SELECT COUNT(*), COUNT(location) FROM account_favorite_locations;
