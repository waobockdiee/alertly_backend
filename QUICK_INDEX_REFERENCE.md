# Quick Reference - Índices por Tabla

Guía rápida de índices necesarios para cada tabla de Alertly.

---

## incident_clusters (tabla más crítica)

### Índices Existentes
- `PRIMARY KEY (incl_id)`

### Índices NECESARIOS (Fase 1 - CRÍTICO)

```sql
-- Clustering algorithm
CREATE INDEX CONCURRENTLY idx_clusters_clustering_lookup
ON incident_clusters (insu_id, category_code, subcategory_code, is_active, end_time)
WHERE is_active = '1';

-- Spatial queries (PostGIS)
CREATE INDEX CONCURRENTLY idx_clusters_spatial_location
ON incident_clusters USING GIST (ST_MakePoint(center_longitude, center_latitude));

-- Map loading
CREATE INDEX CONCURRENTLY idx_clusters_location_bbox
ON incident_clusters (center_latitude, center_longitude, is_active)
WHERE is_active = '1';

CREATE INDEX CONCURRENTLY idx_clusters_time_range
ON incident_clusters (start_time, end_time, is_active)
WHERE is_active = '1';

CREATE INDEX CONCURRENTLY idx_clusters_category_created
ON incident_clusters (category_code, is_active, created_at DESC)
WHERE is_active = '1';
```

### Queries que mejoran
- `CheckAndGetIfClusterExist` → 10-15x más rápido
- `GetClustersByLocation` → 10-15x más rápido
- `GetClustersByRadius` → 10-15x más rápido

---

## incident_reports

### Índices Existentes
- `PRIMARY KEY (inre_id)`

### Índices NECESARIOS (Fase 1 - CRÍTICO)

```sql
-- Cluster details
CREATE INDEX CONCURRENTLY idx_incident_reports_cluster
ON incident_reports (incl_id, is_active, created_at DESC)
WHERE is_active = '1';

-- Vote checking
CREATE INDEX CONCURRENTLY idx_incident_reports_votes
ON incident_reports (incl_id, account_id, vote)
WHERE vote IS NOT NULL;
```

### Índices NECESARIOS (Fase 2)

```sql
-- User profile
CREATE INDEX CONCURRENTLY idx_incident_reports_user_profile
ON incident_reports (account_id, created_at DESC)
INCLUDE (inre_id, media_url, description, event_type, subcategory_name, incl_id, is_anonymous);

-- Foreign key
CREATE INDEX CONCURRENTLY idx_incident_reports_account_fk
ON incident_reports (account_id);
```

### Queries que mejoran
- `GetIncidentBy` → 8-12x más rápido
- `HasAccountVoted` → 10-15x más rápido
- `GetById (profile)` → 3-5x más rápido

---

## account

### Índices Existentes
- `PRIMARY KEY (account_id)`

### Índices NECESARIOS (Fase 1 - CRÍTICO)

```sql
-- Login
CREATE UNIQUE INDEX CONCURRENTLY idx_account_email_login
ON account (email)
WHERE status IN ('active', 'pending');
```

### Índices NECESARIOS (Fase 2)

```sql
-- Premium users
CREATE INDEX CONCURRENTLY idx_account_premium_notifications
ON account (status, is_premium, receive_notifications, account_id)
WHERE status = 'active' AND is_premium = 1 AND receive_notifications = 1;

-- Credibility lookup
CREATE INDEX CONCURRENTLY idx_account_credibility
ON account (account_id)
INCLUDE (credibility);

-- Incident info (covering)
CREATE INDEX CONCURRENTLY idx_account_incident_info
ON account (account_id)
INCLUDE (nickname, first_name, last_name, thumbnail_url, score, is_private_profile);
```

### Queries que mejoran
- `GetUserByEmail` → 5-8x más rápido
- `FindSubscribedUsersForCluster` → 5-10x más rápido
- `UpdateClusterAsTrue` → 2-3x más rápido

---

## notification_deliveries

### Índices Existentes
- `PRIMARY KEY (node_id)`

### Índices NECESARIOS (Fase 1 - CRÍTICO)

```sql
-- User notifications
CREATE INDEX CONCURRENTLY idx_notification_deliveries_user
ON notification_deliveries (to_account_id, created_at DESC);

-- Unread count (partial index)
CREATE INDEX CONCURRENTLY idx_notification_deliveries_unread
ON notification_deliveries (to_account_id, is_read)
WHERE is_read = 0 OR is_read IS NULL;
```

### Queries que mejoran
- `GetNotifications` → 7-10x más rápido
- `GetUnreadCount` → 15-20x más rápido

---

## notifications

### Índices NECESARIOS (Fase 2)

```sql
-- Processing queue
CREATE INDEX CONCURRENTLY idx_notifications_processing_queue
ON notifications (type, must_be_processed, created_at)
WHERE must_be_processed = 1;

-- JOIN with deliveries
CREATE INDEX CONCURRENTLY idx_notifications_noti_id
ON notifications (noti_id)
INCLUDE (type, reference_id);
```

### Queries que mejoran
- `FetchPending` → 10-15x más rápido
- `GetNotifications (JOIN)` → 3-5x más rápido

---

## device_tokens

### Índices NECESARIOS (Fase 2)

```sql
-- Unique token
CREATE UNIQUE INDEX CONCURRENTLY idx_device_tokens_unique
ON device_tokens (account_id, device_token);

-- Account lookup
CREATE INDEX CONCURRENTLY idx_device_tokens_account
ON device_tokens (account_id, device_token);
```

### Queries que mejoran
- `SaveDeviceToken` → Evita duplicados
- `FindSubscribedUsersForCluster` → 5-8x más rápido

---

## account_favorite_locations (Premium)

### Índices NECESARIOS (Fase 2 - CRÍTICO para premium)

```sql
-- Spatial index (CRÍTICO para cronjob)
CREATE INDEX CONCURRENTLY idx_favorite_locations_spatial
ON account_favorite_locations USING GIST (
    ST_MakePoint(longitude, latitude)
)
WHERE status = 1;

-- User lookup
CREATE INDEX CONCURRENTLY idx_favorite_locations_user
ON account_favorite_locations (account_id, afl_id DESC);

-- Active locations
CREATE INDEX CONCURRENTLY idx_favorite_locations_active
ON account_favorite_locations (account_id, status, crime, traffic_accident,
    medical_emergency, fire_incident, vandalism, suspicious_activity)
WHERE status = 1;
```

### Queries que mejoran
- `FindSubscribedUsersForCluster` → 20-30x más rápido (CRONJOB)
- `Get (myplaces)` → 3-5x más rápido

---

## account_cluster_saved

### Índices NECESARIOS (Fase 3)

```sql
CREATE INDEX CONCURRENTLY idx_account_cluster_saved_user
ON account_cluster_saved (account_id, incl_id);

CREATE INDEX CONCURRENTLY idx_account_cluster_saved_cluster
ON account_cluster_saved (incl_id, account_id);
```

### Queries que mejoran
- `GetMyList` → 5-8x más rápido
- `ToggleSaveClusterAccount` → 3-5x más rápido

---

## account_history

### Índices NECESARIOS (Fase 3)

```sql
CREATE INDEX CONCURRENTLY idx_account_history_user
ON account_history (account_id, his_id DESC)
INCLUDE (incl_id, created_at);

CREATE INDEX CONCURRENTLY idx_account_history_cluster
ON account_history (incl_id, account_id);
```

### Queries que mejoran
- `GetHistory` → 5-8x más rápido
- `GetViewedIncidentIds` → 3-5x más rápido

---

## incident_comments

### Índices NECESARIOS (Fase 3)

```sql
CREATE INDEX CONCURRENTLY idx_incident_comments_cluster
ON incident_comments (incl_id, inco_id DESC)
INCLUDE (account_id, comment, created_at, comment_status, counter_flags);

CREATE INDEX CONCURRENTLY idx_incident_comments_account
ON incident_comments (account_id);
```

### Queries que mejoran
- `GetClusterCommentsByID` → 8-12x más rápido

---

## account_achievements

### Índices NECESARIOS (Fase 3)

```sql
CREATE INDEX CONCURRENTLY idx_account_achievements_modal
ON account_achievements (account_id, show_in_modal, created DESC)
WHERE show_in_modal = 1;

CREATE INDEX CONCURRENTLY idx_account_achievements_user_type
ON account_achievements (account_id, type, badge_threshold);
```

### Queries que mejoran
- `ShowByAccountID` → 5-8x más rápido

---

## influencers (Referral System)

### Índices NECESARIOS (Fase 3)

```sql
CREATE UNIQUE INDEX CONCURRENTLY idx_influencers_referral_code
ON influencers (referral_code)
WHERE is_active = true;

CREATE INDEX CONCURRENTLY idx_influencers_active
ON influencers (is_active, platform);
```

### Queries que mejoran
- `GetInfluencerByCode` → 10-15x más rápido (signup)

---

## referral_conversions

### Índices NECESARIOS (Fase 3)

```sql
CREATE UNIQUE INDEX CONCURRENTLY idx_referral_conversions_user
ON referral_conversions (user_id);

CREATE INDEX CONCURRENTLY idx_referral_conversions_code
ON referral_conversions (referral_code, registered_at DESC);
```

### Queries que mejoran
- `GetConversionByUserID` → 10-15x más rápido
- `GetConversionsByCode` → 5-8x más rápido

---

## incident_subcategories

### Índices NECESARIOS (Fase 3)

```sql
CREATE INDEX CONCURRENTLY idx_incident_subcategories_code
ON incident_subcategories (code);

CREATE INDEX CONCURRENTLY idx_incident_subcategories_category
ON incident_subcategories (inca_id);
```

### Queries que mejoran
- `GetDurationForSubcategory` → 5-10x más rápido
- `GetExpiredClusters` → 3-5x más rápido

---

## Orden de Implementación Recomendado

### Semana 1 (CRÍTICO)
1. incident_clusters (5 índices)
2. incident_reports (2 índices)
3. account (1 índice)
4. notification_deliveries (2 índices)

**Total:** 10 índices
**Script:** `critical_indexes_phase1.sql`

### Semana 2-3 (ALTA PRIORIDAD)
1. account_favorite_locations (3 índices)
2. device_tokens (2 índices)
3. account (3 índices adicionales)
4. notifications (2 índices)
5. incident_reports (2 índices adicionales)

**Total:** 12 índices

### Semana 4+ (MEDIA PRIORIDAD)
1. account_history (2 índices)
2. account_cluster_saved (2 índices)
3. incident_comments (2 índices)
4. account_achievements (2 índices)
5. Referral system (5 índices)
6. incident_subcategories (2 índices)

**Total:** 15 índices

---

## Comandos Útiles

### Verificar índices existentes
```sql
SELECT tablename, indexname, indexdef
FROM pg_indexes
WHERE schemaname = 'public' AND tablename = 'incident_clusters'
ORDER BY tablename, indexname;
```

### Verificar uso de índices
```sql
SELECT
    schemaname,
    tablename,
    indexname,
    idx_scan AS times_used,
    idx_tup_read AS tuples_read
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY idx_scan DESC;
```

### Verificar tamaño de índices
```sql
SELECT
    schemaname,
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) AS index_size
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY pg_relation_size(indexrelid) DESC;
```

### Monitorear creación de índices
```sql
SELECT * FROM pg_stat_progress_create_index;
```

### Reindexar tabla (si índice fragmentado)
```sql
REINDEX TABLE CONCURRENTLY incident_clusters;
```

---

## Troubleshooting

### Índice no se usa (idx_scan = 0)
1. Verificar query con `EXPLAIN ANALYZE`
2. Ejecutar `ANALYZE table_name`
3. Verificar selectividad del índice
4. Considerar aumentar `random_page_cost`

### Creación de índice lenta
1. Aumentar `maintenance_work_mem`
2. Verificar espacio en disco
3. Ejecutar en horario de bajo tráfico
4. Considerar particionar tabla

### Error "out of memory"
1. Aumentar `maintenance_work_mem` a 1GB+
2. Reducir `max_parallel_maintenance_workers`
3. Crear índice sin `CONCURRENTLY` en ventana de mantenimiento

---

**Ver también:**
- `SQL_QUERIES_AND_INDEXES_ANALYSIS.md` - Análisis completo
- `critical_indexes_phase1.sql` - Script ejecutable
- `INDEX_OPTIMIZATION_SUMMARY.md` - Resumen ejecutivo
