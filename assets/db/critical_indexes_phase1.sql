-- =====================================================================
-- ALERTLY DATABASE - CRITICAL INDEXES (FASE 1)
-- =====================================================================
-- Descripción: Índices críticos para performance de endpoints más usados
-- Impacto esperado: 10-50x mejora en queries críticas
-- Tiempo estimado de ejecución: 5-10 minutos
-- Downtime: 0 (usando CONCURRENTLY)
-- =====================================================================
-- IMPORTANTE: Ejecutar en horario de bajo tráfico (madrugada)
-- Monitorear con: SELECT * FROM pg_stat_progress_create_index;
-- =====================================================================

-- Iniciar transacción (solo para verificaciones, los CREATE INDEX CONCURRENTLY no pueden estar en transacción)
BEGIN;

-- Verificar que PostGIS está instalado
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'postgis') THEN
        RAISE EXCEPTION 'PostGIS extension is not installed. Run: CREATE EXTENSION postgis;';
    END IF;
END
$$;

COMMIT;

-- =====================================================================
-- TABLA: incident_clusters (CRÍTICO - clustering algorithm)
-- =====================================================================

-- 1. Índice para clustering algorithm (CheckAndGetIfClusterExist)
-- Usado en: POST /newincident (CADA nuevo reporte)
-- Impacto: 10-15x más rápido
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_clusters_clustering_lookup
ON incident_clusters (insu_id, category_code, subcategory_code, is_active, end_time)
WHERE is_active = '1';

COMMENT ON INDEX idx_clusters_clustering_lookup IS
'Índice compuesto para algoritmo de clustering. Usado en cada nuevo reporte para buscar clusters existentes.';

-- 2. Índice espacial PostGIS para distancias (ST_DistanceSphere)
-- Usado en: POST /newincident, cronjob notifications
-- Impacto: 20-30x más rápido en queries espaciales
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_clusters_spatial_location
ON incident_clusters USING GIST (ST_MakePoint(center_longitude, center_latitude));

COMMENT ON INDEX idx_clusters_spatial_location IS
'Índice espacial GiST para cálculos de distancia (ST_DistanceSphere). Crítico para clustering y notificaciones.';

-- =====================================================================
-- TABLA: incident_clusters (CRÍTICO - map loading)
-- =====================================================================

-- 3. Índice para bounding box queries (mapa)
-- Usado en: GET /getclustersbylocation (ENDPOINT MÁS USADO)
-- Impacto: 10-15x más rápido
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_clusters_location_bbox
ON incident_clusters (center_latitude, center_longitude, is_active)
WHERE is_active = '1';

COMMENT ON INDEX idx_clusters_location_bbox IS
'Índice para búsquedas por bounding box (mapa). Usado en cada carga del mapa.';

-- 4. Índice para filtrado por tiempo (start_time/end_time)
-- Usado en: GET /getclustersbylocation, GET /getclusterbyradius
-- Impacto: 5-8x más rápido
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_clusters_time_range
ON incident_clusters (start_time, end_time, is_active)
WHERE is_active = '1';

COMMENT ON INDEX idx_clusters_time_range IS
'Índice para filtrado por rango de tiempo. Evita DATE() function scan.';

-- 5. Índice para categorías + ordering
-- Usado en: GET /getclustersbylocation (ORDER BY created_at)
-- Impacto: 3-5x más rápido
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_clusters_category_created
ON incident_clusters (category_code, is_active, created_at DESC)
WHERE is_active = '1';

COMMENT ON INDEX idx_clusters_category_created IS
'Índice para filtrado por categoría + ordenamiento por fecha. Usado en queries de mapa con categorías.';

-- =====================================================================
-- TABLA: account (CRÍTICO - login)
-- =====================================================================

-- 6. Índice único para login (email lookup)
-- Usado en: POST /auth/login (CADA LOGIN)
-- Impacto: 5-8x más rápido
CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_account_email_login
ON account (email)
WHERE status IN ('active', 'pending');

COMMENT ON INDEX idx_account_email_login IS
'Índice único para login. Búsqueda rápida de usuarios por email.';

-- =====================================================================
-- TABLA: incident_reports (CRÍTICO - cluster details)
-- =====================================================================

-- 7. Índice para lookup de reports por cluster
-- Usado en: GET /getclusterby/:id (CADA apertura de incidente)
-- Impacto: 8-12x más rápido
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_incident_reports_cluster
ON incident_reports (incl_id, is_active, created_at DESC)
WHERE is_active = '1';

COMMENT ON INDEX idx_incident_reports_cluster IS
'Índice para obtener reports de un cluster específico. Usado en detalles de incidente.';

-- 8. Índice para verificación de votos
-- Usado en: POST /newincident (verificar si ya votó)
-- Impacto: 10-15x más rápido
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_incident_reports_votes
ON incident_reports (incl_id, account_id, vote)
WHERE vote IS NOT NULL;

COMMENT ON INDEX idx_incident_reports_votes IS
'Índice para verificar si un usuario ya votó en un cluster. Previene votos duplicados.';

-- =====================================================================
-- TABLA: notification_deliveries (CRÍTICO - notification center)
-- =====================================================================

-- 9. Índice para notificaciones del usuario
-- Usado en: GET /notifications (CADA apertura de notification center)
-- Impacto: 7-10x más rápido
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_notification_deliveries_user
ON notification_deliveries (to_account_id, created_at DESC);
-- NOTA: Incluir columnas en PostgreSQL 11+ con INCLUDE
-- INCLUDE (node_id, is_read, noti_id, title, message)

COMMENT ON INDEX idx_notification_deliveries_user IS
'Índice para listar notificaciones de un usuario. Usado en notification center.';

-- 10. Índice parcial para contador de no leídas (badge)
-- Usado en: GET /notifications/unread-count (POLLING cada 30s)
-- Impacto: 15-20x más rápido
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_notification_deliveries_unread
ON notification_deliveries (to_account_id, is_read)
WHERE is_read = 0 OR is_read IS NULL;

COMMENT ON INDEX idx_notification_deliveries_unread IS
'Índice parcial para contar notificaciones no leídas. Usado en badge count.';

-- =====================================================================
-- VERIFICACIÓN POST-CREACIÓN
-- =====================================================================

-- Verificar que todos los índices se crearon correctamente
DO $$
DECLARE
    missing_indexes TEXT[] := ARRAY[]::TEXT[];
BEGIN
    -- Verificar cada índice
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_clusters_clustering_lookup') THEN
        missing_indexes := array_append(missing_indexes, 'idx_clusters_clustering_lookup');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_clusters_spatial_location') THEN
        missing_indexes := array_append(missing_indexes, 'idx_clusters_spatial_location');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_clusters_location_bbox') THEN
        missing_indexes := array_append(missing_indexes, 'idx_clusters_location_bbox');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_clusters_time_range') THEN
        missing_indexes := array_append(missing_indexes, 'idx_clusters_time_range');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_clusters_category_created') THEN
        missing_indexes := array_append(missing_indexes, 'idx_clusters_category_created');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_account_email_login') THEN
        missing_indexes := array_append(missing_indexes, 'idx_account_email_login');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_incident_reports_cluster') THEN
        missing_indexes := array_append(missing_indexes, 'idx_incident_reports_cluster');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_incident_reports_votes') THEN
        missing_indexes := array_append(missing_indexes, 'idx_incident_reports_votes');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_notification_deliveries_user') THEN
        missing_indexes := array_append(missing_indexes, 'idx_notification_deliveries_user');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_notification_deliveries_unread') THEN
        missing_indexes := array_append(missing_indexes, 'idx_notification_deliveries_unread');
    END IF;

    -- Reportar resultados
    IF array_length(missing_indexes, 1) > 0 THEN
        RAISE WARNING 'Los siguientes índices NO se crearon: %', missing_indexes;
    ELSE
        RAISE NOTICE '✅ Todos los índices de FASE 1 se crearon correctamente';
    END IF;
END
$$;

-- =====================================================================
-- MANTENIMIENTO POST-CREACIÓN
-- =====================================================================

-- Actualizar estadísticas de las tablas (para que el query planner use los nuevos índices)
ANALYZE incident_clusters;
ANALYZE incident_reports;
ANALYZE account;
ANALYZE notification_deliveries;

-- =====================================================================
-- QUERIES DE MONITOREO
-- =====================================================================

-- Verificar tamaño de índices creados
SELECT
    schemaname,
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) AS index_size
FROM pg_stat_user_indexes
WHERE indexname IN (
    'idx_clusters_clustering_lookup',
    'idx_clusters_spatial_location',
    'idx_clusters_location_bbox',
    'idx_clusters_time_range',
    'idx_clusters_category_created',
    'idx_account_email_login',
    'idx_incident_reports_cluster',
    'idx_incident_reports_votes',
    'idx_notification_deliveries_user',
    'idx_notification_deliveries_unread'
)
ORDER BY pg_relation_size(indexrelid) DESC;

-- Verificar uso de índices (ejecutar después de 24h en producción)
-- SELECT
--     schemaname,
--     tablename,
--     indexname,
--     idx_scan AS "times_used",
--     idx_tup_read AS "tuples_read",
--     idx_tup_fetch AS "tuples_fetched"
-- FROM pg_stat_user_indexes
-- WHERE indexname LIKE 'idx_clusters_%'
--    OR indexname LIKE 'idx_account_%'
--    OR indexname LIKE 'idx_incident_reports_%'
--    OR indexname LIKE 'idx_notification_deliveries_%'
-- ORDER BY idx_scan DESC;

-- =====================================================================
-- NOTAS FINALES
-- =====================================================================
/*
1. EJECUCIÓN:
   - Ejecutar en horario de bajo tráfico (2-5 AM)
   - CONCURRENTLY permite operaciones sin bloquear escrituras
   - Tiempo estimado: 5-10 minutos (depende del tamaño de datos)

2. MONITOREO DURANTE EJECUCIÓN:
   SELECT * FROM pg_stat_progress_create_index;

3. ROLLBACK (si es necesario):
   DROP INDEX CONCURRENTLY idx_clusters_clustering_lookup;
   DROP INDEX CONCURRENTLY idx_clusters_spatial_location;
   -- ... etc

4. VALIDACIÓN POST-DEPLOYMENT:
   - Ejecutar queries de monitoreo después de 24h
   - Verificar que idx_scan > 0 para cada índice
   - Comparar response times antes/después

5. PRÓXIMOS PASOS:
   - Si FASE 1 exitosa → ejecutar FASE 2 (critical_indexes_phase2.sql)
   - Si hay problemas → revisar logs de PostgreSQL
   - Considerar incrementar shared_buffers si hay memory issues

6. TROUBLESHOOTING:
   - Si índice falla: revisar espacio en disco (df -h)
   - Si índice lento: verificar maintenance_work_mem
   - Si queries lentas después: ejecutar VACUUM ANALYZE

7. CONTACT:
   - Ver SQL_QUERIES_AND_INDEXES_ANALYSIS.md para detalles completos
*/
