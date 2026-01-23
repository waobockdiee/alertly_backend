-- ============================================================
-- POSTGRESQL OPTIMIZATION SCRIPT - ALERTLY DATABASE
-- Generated: 2026-01-22
-- Database: railway @ metro.proxy.rlwy.net:48204
-- Version: PostgreSQL 16.9 with PostGIS 3.7.0dev
--
-- ADVERTENCIA: Este script debe ejecutarse en ORDEN
-- Hacer BACKUP completo antes de ejecutar
-- ============================================================

\echo '=========================================='
\echo 'ALERTLY - PostgreSQL Optimization Script'
\echo 'Starting migration...'
\echo '=========================================='

-- Iniciar transacción
BEGIN;

-- ============================================================
-- PASO 1: VALIDAR DATOS ANTES DE AGREGAR FOREIGN KEYS
-- ============================================================
\echo 'PASO 1: Validando datos huérfanos...'

DO $$
DECLARE
    orphan_count INTEGER;
    total_orphans INTEGER := 0;
BEGIN
    -- Verificar account_cluster_saved -> incident_clusters
    SELECT COUNT(*) INTO orphan_count
    FROM account_cluster_saved acs
    LEFT JOIN incident_clusters ic ON acs.incl_id = ic.incl_id
    WHERE ic.incl_id IS NULL;

    IF orphan_count > 0 THEN
        RAISE WARNING 'account_cluster_saved: % registros huérfanos', orphan_count;
        total_orphans := total_orphans + orphan_count;
    END IF;

    -- Verificar account_history -> incident_clusters
    SELECT COUNT(*) INTO orphan_count
    FROM account_history ah
    LEFT JOIN incident_clusters ic ON ah.incl_id = ic.incl_id
    WHERE ah.incl_id IS NOT NULL AND ic.incl_id IS NULL;

    IF orphan_count > 0 THEN
        RAISE WARNING 'account_history: % registros huérfanos', orphan_count;
        total_orphans := total_orphans + orphan_count;
    END IF;

    -- Verificar notification_deliveries -> notifications
    SELECT COUNT(*) INTO orphan_count
    FROM notification_deliveries nd
    LEFT JOIN notifications n ON nd.noti_id = n.noti_id
    WHERE n.noti_id IS NULL;

    IF orphan_count > 0 THEN
        RAISE WARNING 'notification_deliveries: % registros huérfanos', orphan_count;
        total_orphans := total_orphans + orphan_count;
    END IF;

    -- Verificar incident_clusters -> incident_subcategories
    SELECT COUNT(*) INTO orphan_count
    FROM incident_clusters ic
    LEFT JOIN incident_subcategories isc ON ic.insu_id = isc.insu_id
    WHERE isc.insu_id IS NULL;

    IF orphan_count > 0 THEN
        RAISE WARNING 'incident_clusters: % registros con insu_id inválido', orphan_count;
        total_orphans := total_orphans + orphan_count;
    END IF;

    -- Verificar incident_reports -> incident_subcategories
    SELECT COUNT(*) INTO orphan_count
    FROM incident_reports ir
    LEFT JOIN incident_subcategories isc ON ir.insu_id = isc.insu_id
    WHERE isc.insu_id IS NULL;

    IF orphan_count > 0 THEN
        RAISE WARNING 'incident_reports: % registros con insu_id inválido', orphan_count;
        total_orphans := total_orphans + orphan_count;
    END IF;

    -- Verificar incident_flags -> incident_reports
    SELECT COUNT(*) INTO orphan_count
    FROM incident_flags if_tbl
    LEFT JOIN incident_reports ir ON if_tbl.inre_id = ir.inre_id
    WHERE ir.inre_id IS NULL;

    IF orphan_count > 0 THEN
        RAISE WARNING 'incident_flags: % registros huérfanos', orphan_count;
        total_orphans := total_orphans + orphan_count;
    END IF;

    -- Verificar incident_logs -> incident_reports
    SELECT COUNT(*) INTO orphan_count
    FROM incident_logs il
    LEFT JOIN incident_reports ir ON il.inre_id = ir.inre_id
    WHERE ir.inre_id IS NULL;

    IF orphan_count > 0 THEN
        RAISE WARNING 'incident_logs: % registros huérfanos', orphan_count;
        total_orphans := total_orphans + orphan_count;
    END IF;

    IF total_orphans > 0 THEN
        RAISE EXCEPTION 'Se encontraron % registros huérfanos. Limpiar datos antes de continuar.', total_orphans;
    ELSE
        RAISE NOTICE 'Validación OK: No se encontraron datos huérfanos';
    END IF;
END $$;

SAVEPOINT after_validation;

-- ============================================================
-- PASO 2: AGREGAR FOREIGN KEYS FALTANTES
-- ============================================================
\echo 'PASO 2: Creando Foreign Keys...'

-- FK: account_cluster_saved -> incident_clusters
ALTER TABLE account_cluster_saved
DROP CONSTRAINT IF EXISTS fk_account_cluster_saved_incident_clusters;

ALTER TABLE account_cluster_saved
ADD CONSTRAINT fk_account_cluster_saved_incident_clusters
FOREIGN KEY (incl_id)
REFERENCES incident_clusters(incl_id)
ON DELETE CASCADE
ON UPDATE CASCADE;

\echo '  ✓ FK: account_cluster_saved.incl_id -> incident_clusters.incl_id'

-- FK: account_history -> incident_clusters
ALTER TABLE account_history
DROP CONSTRAINT IF EXISTS fk_account_history_incident_clusters;

ALTER TABLE account_history
ADD CONSTRAINT fk_account_history_incident_clusters
FOREIGN KEY (incl_id)
REFERENCES incident_clusters(incl_id)
ON DELETE CASCADE
ON UPDATE CASCADE;

\echo '  ✓ FK: account_history.incl_id -> incident_clusters.incl_id'

-- FK: notification_deliveries -> notifications (CRÍTICO)
ALTER TABLE notification_deliveries
DROP CONSTRAINT IF EXISTS fk_notification_deliveries_notifications;

ALTER TABLE notification_deliveries
ADD CONSTRAINT fk_notification_deliveries_notifications
FOREIGN KEY (noti_id)
REFERENCES notifications(noti_id)
ON DELETE CASCADE
ON UPDATE CASCADE;

\echo '  ✓ FK: notification_deliveries.noti_id -> notifications.noti_id'

-- FK: incident_clusters -> incident_subcategories
ALTER TABLE incident_clusters
DROP CONSTRAINT IF EXISTS fk_incident_clusters_subcategories;

ALTER TABLE incident_clusters
ADD CONSTRAINT fk_incident_clusters_subcategories
FOREIGN KEY (insu_id)
REFERENCES incident_subcategories(insu_id)
ON DELETE RESTRICT
ON UPDATE CASCADE;

\echo '  ✓ FK: incident_clusters.insu_id -> incident_subcategories.insu_id'

-- FK: incident_reports -> incident_subcategories
ALTER TABLE incident_reports
DROP CONSTRAINT IF EXISTS fk_incident_reports_subcategories;

ALTER TABLE incident_reports
ADD CONSTRAINT fk_incident_reports_subcategories
FOREIGN KEY (insu_id)
REFERENCES incident_subcategories(insu_id)
ON DELETE RESTRICT
ON UPDATE CASCADE;

\echo '  ✓ FK: incident_reports.insu_id -> incident_subcategories.insu_id'

-- FK: incident_flags -> incident_reports
ALTER TABLE incident_flags
DROP CONSTRAINT IF EXISTS fk_incident_flags_incident_reports;

ALTER TABLE incident_flags
ADD CONSTRAINT fk_incident_flags_incident_reports
FOREIGN KEY (inre_id)
REFERENCES incident_reports(inre_id)
ON DELETE CASCADE
ON UPDATE CASCADE;

\echo '  ✓ FK: incident_flags.inre_id -> incident_reports.inre_id'

-- FK: incident_logs -> incident_reports
ALTER TABLE incident_logs
DROP CONSTRAINT IF EXISTS fk_incident_logs_incident_reports;

ALTER TABLE incident_logs
ADD CONSTRAINT fk_incident_logs_incident_reports
FOREIGN KEY (inre_id)
REFERENCES incident_reports(inre_id)
ON DELETE CASCADE
ON UPDATE CASCADE;

\echo '  ✓ FK: incident_logs.inre_id -> incident_reports.inre_id'

SAVEPOINT after_foreign_keys;

-- ============================================================
-- PASO 3: AGREGAR COLUMNAS GEOGRAPHY PARA POSTGIS
-- ============================================================
\echo 'PASO 3: Creando columnas GEOGRAPHY para PostGIS...'

-- incident_clusters: Agregar columna center_location
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'incident_clusters'
        AND column_name = 'center_location'
    ) THEN
        ALTER TABLE incident_clusters
        ADD COLUMN center_location GEOGRAPHY(POINT, 4326);

        RAISE NOTICE 'Columna center_location creada';
    ELSE
        RAISE NOTICE 'Columna center_location ya existe';
    END IF;
END $$;

-- Poblar center_location desde coordenadas existentes
UPDATE incident_clusters
SET center_location = ST_SetSRID(
    ST_MakePoint(center_longitude, center_latitude),
    4326
)::geography
WHERE center_latitude IS NOT NULL
  AND center_longitude IS NOT NULL
  AND center_location IS NULL;

\echo '  ✓ incident_clusters.center_location poblada'

-- account_favorite_locations: Agregar columna location
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'account_favorite_locations'
        AND column_name = 'location'
    ) THEN
        ALTER TABLE account_favorite_locations
        ADD COLUMN location GEOGRAPHY(POINT, 4326);

        RAISE NOTICE 'Columna location creada';
    ELSE
        RAISE NOTICE 'Columna location ya existe';
    END IF;
END $$;

-- Poblar location desde coordenadas existentes
UPDATE account_favorite_locations
SET location = ST_SetSRID(
    ST_MakePoint(longitude, latitude),
    4326
)::geography
WHERE latitude IS NOT NULL
  AND longitude IS NOT NULL
  AND location IS NULL;

\echo '  ✓ account_favorite_locations.location poblada'

SAVEPOINT after_geography_columns;

-- ============================================================
-- PASO 4: CREAR ÍNDICES ESPACIALES (GiST)
-- ============================================================
\echo 'PASO 4: Creando índices espaciales GiST...'

-- Índice espacial para incident_clusters
DROP INDEX IF EXISTS idx_clusters_center_location_gist;
CREATE INDEX idx_clusters_center_location_gist
ON incident_clusters
USING GIST (center_location);

\echo '  ✓ idx_clusters_center_location_gist creado'

-- Índice espacial para account_favorite_locations
DROP INDEX IF EXISTS idx_favorite_locations_gist;
CREATE INDEX idx_favorite_locations_gist
ON account_favorite_locations
USING GIST (location);

\echo '  ✓ idx_favorite_locations_gist creado'

SAVEPOINT after_spatial_indexes;

-- ============================================================
-- PASO 5: CREAR ÍNDICES DE PERFORMANCE FALTANTES
-- ============================================================
\echo 'PASO 5: Creando índices de performance...'

-- Índice compuesto para geolocalización con time window
-- NOTA: Después de implementar GEOGRAPHY, este puede ser deprecado
DROP INDEX IF EXISTS idx_clusters_location_time;
CREATE INDEX idx_clusters_location_time
ON incident_clusters (center_latitude, center_longitude, start_time, end_time, is_active)
WHERE is_active = '1';

\echo '  ✓ idx_clusters_location_time'

-- Índice para clustering de incidentes por subcategoría
DROP INDEX IF EXISTS idx_clusters_insu_created;
CREATE INDEX idx_clusters_insu_created
ON incident_clusters (insu_id, created_at DESC, is_active)
WHERE is_active = '1';

\echo '  ✓ idx_clusters_insu_created'

-- Índice para validación de votos duplicados
DROP INDEX IF EXISTS idx_reports_incl_account;
CREATE INDEX idx_reports_incl_account
ON incident_reports (incl_id, account_id, is_active)
WHERE is_active = '1';

\echo '  ✓ idx_reports_incl_account'

-- Índice para historial de reportes por usuario
DROP INDEX IF EXISTS idx_reports_account_created;
CREATE INDEX idx_reports_account_created
ON incident_reports (account_id, created_at DESC);

\echo '  ✓ idx_reports_account_created'

-- Índice para lugares favoritos por usuario
DROP INDEX IF EXISTS idx_favorite_locations_account;
CREATE INDEX idx_favorite_locations_account
ON account_favorite_locations (account_id, status)
WHERE status = 1;

\echo '  ✓ idx_favorite_locations_account'

-- Índice para activación de cuenta
DROP INDEX IF EXISTS idx_account_activation;
CREATE INDEX idx_account_activation
ON account (activation_code, status)
WHERE activation_code IS NOT NULL;

\echo '  ✓ idx_account_activation'

-- Índice para categorías por código
DROP INDEX IF EXISTS idx_categories_code;
CREATE INDEX idx_categories_code
ON incident_categories (code);

\echo '  ✓ idx_categories_code'

-- Índice para subcategorías
DROP INDEX IF EXISTS idx_subcategories_category;
CREATE INDEX idx_subcategories_category
ON incident_subcategories (inca_id, code);

\echo '  ✓ idx_subcategories_category'

-- Índice para clusters guardados
DROP INDEX IF EXISTS idx_cluster_saved_account;
CREATE INDEX idx_cluster_saved_account
ON account_cluster_saved (account_id, created_at DESC);

\echo '  ✓ idx_cluster_saved_account'

-- Índice para comentarios
DROP INDEX IF EXISTS idx_comments_cluster_created;
CREATE INDEX idx_comments_cluster_created
ON incident_comments (incl_id, created_at DESC);

\echo '  ✓ idx_comments_cluster_created'

-- ============================================================
-- ÍNDICES ADICIONALES RECOMENDADOS
-- ============================================================
\echo 'PASO 5b: Creando índices adicionales recomendados...'

-- Índice para notificaciones no procesadas (cronjobs)
DROP INDEX IF EXISTS idx_notifications_must_process;
CREATE INDEX idx_notifications_must_process
ON notifications (must_be_processed, created_at DESC)
WHERE must_be_processed = 1;

\echo '  ✓ idx_notifications_must_process'

-- Índice para device tokens activos
DROP INDEX IF EXISTS idx_device_tokens_active;
CREATE INDEX idx_device_tokens_active
ON device_tokens (device_token, updated_at DESC)
WHERE account_id IS NOT NULL;

\echo '  ✓ idx_device_tokens_active'

-- Índice para cuenta por nickname
DROP INDEX IF EXISTS idx_account_nickname_status;
CREATE INDEX idx_account_nickname_status
ON account (nickname, status)
WHERE status = 1;

\echo '  ✓ idx_account_nickname_status'

SAVEPOINT after_performance_indexes;

-- ============================================================
-- PASO 6: AGREGAR CONSTRAINTS NOT NULL
-- ============================================================
\echo 'PASO 6: Agregando constraints NOT NULL...'

-- Verificar que no hay valores NULL en columnas críticas
DO $$
DECLARE
    null_count INTEGER;
BEGIN
    -- Verificar incident_clusters
    SELECT COUNT(*) INTO null_count
    FROM incident_clusters
    WHERE center_latitude IS NULL
       OR center_longitude IS NULL
       OR center_location IS NULL;

    IF null_count > 0 THEN
        RAISE EXCEPTION 'incident_clusters tiene % registros con coordenadas NULL', null_count;
    END IF;

    -- Verificar created_at
    SELECT COUNT(*) INTO null_count
    FROM incident_clusters
    WHERE created_at IS NULL;

    IF null_count > 0 THEN
        RAISE EXCEPTION 'incident_clusters tiene % registros con created_at NULL', null_count;
    END IF;

    RAISE NOTICE 'Validación NOT NULL OK';
END $$;

-- Agregar NOT NULL constraints
ALTER TABLE incident_clusters
    ALTER COLUMN center_latitude SET NOT NULL,
    ALTER COLUMN center_longitude SET NOT NULL,
    ALTER COLUMN center_location SET NOT NULL,
    ALTER COLUMN created_at SET NOT NULL,
    ALTER COLUMN created_at SET DEFAULT NOW();

\echo '  ✓ incident_clusters: NOT NULL constraints agregados'

-- incident_reports
ALTER TABLE incident_reports
    ALTER COLUMN created_at SET NOT NULL,
    ALTER COLUMN created_at SET DEFAULT NOW();

\echo '  ✓ incident_reports: NOT NULL en created_at'

-- account
ALTER TABLE account
    ALTER COLUMN created_at SET NOT NULL,
    ALTER COLUMN created_at SET DEFAULT NOW();

\echo '  ✓ account: NOT NULL en created_at'

-- notifications
ALTER TABLE notifications
    ALTER COLUMN created_at SET NOT NULL,
    ALTER COLUMN created_at SET DEFAULT NOW();

\echo '  ✓ notifications: NOT NULL en created_at'

-- account_favorite_locations
ALTER TABLE account_favorite_locations
    ALTER COLUMN location SET NOT NULL;

\echo '  ✓ account_favorite_locations: NOT NULL en location'

SAVEPOINT after_not_null_constraints;

-- ============================================================
-- PASO 7: CONFIGURAR AUTOVACUUM
-- ============================================================
\echo 'PASO 7: Configurando autovacuum...'

-- Configurar autovacuum más agresivo para notifications
ALTER TABLE notifications SET (
    autovacuum_vacuum_scale_factor = 0.05,
    autovacuum_analyze_scale_factor = 0.05
);

\echo '  ✓ autovacuum configurado para notifications'

-- Configurar para incident_clusters (tabla crítica)
ALTER TABLE incident_clusters SET (
    autovacuum_vacuum_scale_factor = 0.10,
    autovacuum_analyze_scale_factor = 0.10
);

\echo '  ✓ autovacuum configurado para incident_clusters'

SAVEPOINT after_autovacuum;

-- ============================================================
-- PASO 8: ACTUALIZAR ESTADÍSTICAS (ANALYZE)
-- ============================================================
\echo 'PASO 8: Actualizando estadísticas...'

ANALYZE incident_clusters;
\echo '  ✓ ANALYZE incident_clusters'

ANALYZE incident_reports;
\echo '  ✓ ANALYZE incident_reports'

ANALYZE account;
\echo '  ✓ ANALYZE account'

ANALYZE notifications;
\echo '  ✓ ANALYZE notifications'

ANALYZE notification_deliveries;
\echo '  ✓ ANALYZE notification_deliveries'

ANALYZE device_tokens;
\echo '  ✓ ANALYZE device_tokens'

ANALYZE account_favorite_locations;
\echo '  ✓ ANALYZE account_favorite_locations'

ANALYZE account_cluster_saved;
\echo '  ✓ ANALYZE account_cluster_saved'

ANALYZE incident_comments;
\echo '  ✓ ANALYZE incident_comments'

SAVEPOINT after_analyze;

-- ============================================================
-- PASO 9: VALIDACIÓN POST-MIGRACIÓN
-- ============================================================
\echo ''
\echo '=========================================='
\echo 'VALIDACIÓN POST-MIGRACIÓN'
\echo '=========================================='

-- Contar foreign keys
\echo 'Foreign Keys creadas:'
SELECT
    tc.table_name,
    tc.constraint_name,
    kcu.column_name,
    ccu.table_name AS foreign_table_name
FROM information_schema.table_constraints AS tc
JOIN information_schema.key_column_usage AS kcu
    ON tc.constraint_name = kcu.constraint_name
JOIN information_schema.constraint_column_usage AS ccu
    ON ccu.constraint_name = tc.constraint_name
WHERE tc.constraint_type = 'FOREIGN KEY'
    AND tc.table_schema = 'public'
    AND tc.table_name IN (
        'account_cluster_saved',
        'account_history',
        'notification_deliveries',
        'incident_clusters',
        'incident_reports',
        'incident_flags',
        'incident_logs'
    )
ORDER BY tc.table_name;

-- Contar índices espaciales
\echo ''
\echo 'Índices espaciales creados:'
SELECT tablename, indexname
FROM pg_indexes
WHERE schemaname = 'public'
    AND indexdef LIKE '%gist%'
ORDER BY tablename;

-- Verificar columnas GEOGRAPHY pobladas
\echo ''
\echo 'Validación de columnas GEOGRAPHY:'
SELECT
    'incident_clusters' AS table_name,
    COUNT(*) AS total_rows,
    COUNT(center_location) AS rows_with_location,
    COUNT(*) - COUNT(center_location) AS rows_missing_location
FROM incident_clusters
UNION ALL
SELECT
    'account_favorite_locations',
    COUNT(*),
    COUNT(location),
    COUNT(*) - COUNT(location)
FROM account_favorite_locations;

-- Estadísticas generales
\echo ''
\echo 'Estadísticas de base de datos:'
SELECT
    'Foreign Keys' AS check_type,
    COUNT(*) AS count
FROM information_schema.table_constraints
WHERE constraint_type = 'FOREIGN KEY' AND table_schema = 'public'
UNION ALL
SELECT 'Total Indexes', COUNT(*)
FROM pg_indexes
WHERE schemaname = 'public'
UNION ALL
SELECT 'Spatial Indexes (GiST)', COUNT(*)
FROM pg_indexes
WHERE schemaname = 'public' AND indexdef LIKE '%gist%'
UNION ALL
SELECT 'Tables', COUNT(*)
FROM pg_tables
WHERE schemaname = 'public';

-- ============================================================
-- COMMIT TRANSACCIÓN
-- ============================================================
COMMIT;

\echo ''
\echo '=========================================='
\echo 'MIGRACIÓN COMPLETADA EXITOSAMENTE'
\echo '=========================================='
\echo ''
\echo 'Próximos pasos:'
\echo '1. Ejecutar VACUUM FULL en horario de bajo tráfico:'
\echo '   VACUUM FULL ANALYZE notifications;'
\echo '   VACUUM ANALYZE account;'
\echo ''
\echo '2. Actualizar backend Go para usar columnas GEOGRAPHY:'
\echo '   - Cambiar queries ST_Distance_Sphere() por ST_DWithin()'
\echo '   - Usar center_location en lugar de center_latitude/longitude'
\echo ''
\echo '3. Monitorear performance:'
\echo '   - Verificar explain plans de queries geoespaciales'
\echo '   - Revisar pg_stat_user_indexes para índices no usados'
\echo '   - Monitorear dead_rows en pg_stat_user_tables'
\echo ''
\echo '4. Configurar monitoring:'
\echo '   - Instalar pg_stat_statements extension'
\echo '   - Configurar alertas para autovacuum failures'
\echo '   - Monitorear tamaño de índices y tablas'
\echo ''
