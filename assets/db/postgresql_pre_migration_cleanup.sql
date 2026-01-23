-- ============================================================
-- PRE-MIGRATION CLEANUP SCRIPT - ALERTLY DATABASE
-- Generated: 2026-01-22
-- Database: railway @ metro.proxy.rlwy.net:48204
--
-- EJECUTAR ESTE SCRIPT ANTES DE postgresql_migration_optimization.sql
-- Este script identifica y opcionalmente limpia datos huérfanos
-- que impedirían la creación de Foreign Keys
-- ============================================================

\echo '=========================================='
\echo 'ALERTLY - Pre-Migration Cleanup Script'
\echo 'Identificando datos huérfanos...'
\echo '=========================================='

-- ============================================================
-- PASO 1: IDENTIFICAR DATOS HUÉRFANOS
-- ============================================================

\echo ''
\echo 'PASO 1: Identificando registros huérfanos...'
\echo ''

-- 1.1 account_cluster_saved -> incident_clusters
\echo '1. account_cluster_saved -> incident_clusters:'
SELECT
    'account_cluster_saved' AS table_name,
    COUNT(*) AS orphan_count,
    'Favoritos apuntando a clusters inexistentes' AS issue
FROM account_cluster_saved acs
LEFT JOIN incident_clusters ic ON acs.incl_id = ic.incl_id
WHERE ic.incl_id IS NULL
HAVING COUNT(*) > 0;

-- Detalle de registros huérfanos
SELECT
    acs.acs_id,
    acs.account_id,
    acs.incl_id AS missing_cluster_id,
    acs.created_at
FROM account_cluster_saved acs
LEFT JOIN incident_clusters ic ON acs.incl_id = ic.incl_id
WHERE ic.incl_id IS NULL
ORDER BY acs.created_at DESC
LIMIT 10;

\echo ''

-- 1.2 account_history -> incident_clusters
\echo '2. account_history -> incident_clusters:'
SELECT
    'account_history' AS table_name,
    COUNT(*) AS orphan_count,
    'Historial con referencias a clusters inexistentes' AS issue
FROM account_history ah
LEFT JOIN incident_clusters ic ON ah.incl_id = ic.incl_id
WHERE ah.incl_id IS NOT NULL AND ic.incl_id IS NULL
HAVING COUNT(*) > 0;

-- Detalle de registros huérfanos
SELECT
    ah.his_id,
    ah.account_id,
    ah.incl_id AS missing_cluster_id,
    ah.created_at
FROM account_history ah
LEFT JOIN incident_clusters ic ON ah.incl_id = ic.incl_id
WHERE ah.incl_id IS NOT NULL AND ic.incl_id IS NULL
ORDER BY ah.created_at DESC
LIMIT 10;

\echo ''

-- 1.3 notification_deliveries -> notifications
\echo '3. notification_deliveries -> notifications (CRÍTICO):'
SELECT
    'notification_deliveries' AS table_name,
    COUNT(*) AS orphan_count,
    'Entregas de notificaciones sin notificación padre' AS issue
FROM notification_deliveries nd
LEFT JOIN notifications n ON nd.noti_id = n.noti_id
WHERE n.noti_id IS NULL
HAVING COUNT(*) > 0;

-- Detalle de registros huérfanos
SELECT
    nd.node_id,
    nd.noti_id AS missing_notification_id,
    nd.to_account_id,
    nd.created_at
FROM notification_deliveries nd
LEFT JOIN notifications n ON nd.noti_id = n.noti_id
WHERE n.noti_id IS NULL
ORDER BY nd.created_at DESC
LIMIT 10;

\echo ''

-- 1.4 incident_clusters -> incident_subcategories
\echo '4. incident_clusters -> incident_subcategories:'
SELECT
    'incident_clusters' AS table_name,
    COUNT(*) AS orphan_count,
    'Clusters con subcategoría inválida' AS issue
FROM incident_clusters ic
LEFT JOIN incident_subcategories isc ON ic.insu_id = isc.insu_id
WHERE isc.insu_id IS NULL
HAVING COUNT(*) > 0;

-- Detalle de registros huérfanos
SELECT
    ic.incl_id,
    ic.insu_id AS invalid_subcategory_id,
    ic.category_code,
    ic.subcategory_code,
    ic.created_at
FROM incident_clusters ic
LEFT JOIN incident_subcategories isc ON ic.insu_id = isc.insu_id
WHERE isc.insu_id IS NULL
ORDER BY ic.created_at DESC
LIMIT 10;

\echo ''

-- 1.5 incident_reports -> incident_subcategories
\echo '5. incident_reports -> incident_subcategories:'
SELECT
    'incident_reports' AS table_name,
    COUNT(*) AS orphan_count,
    'Reportes con subcategoría inválida' AS issue
FROM incident_reports ir
LEFT JOIN incident_subcategories isc ON ir.insu_id = isc.insu_id
WHERE isc.insu_id IS NULL
HAVING COUNT(*) > 0;

-- Detalle de registros huérfanos
SELECT
    ir.inre_id,
    ir.insu_id AS invalid_subcategory_id,
    ir.incl_id,
    ir.account_id,
    ir.created_at
FROM incident_reports ir
LEFT JOIN incident_subcategories isc ON ir.insu_id = isc.insu_id
WHERE isc.insu_id IS NULL
ORDER BY ir.created_at DESC
LIMIT 10;

\echo ''

-- 1.6 incident_flags -> incident_reports
\echo '6. incident_flags -> incident_reports:'
SELECT
    'incident_flags' AS table_name,
    COUNT(*) AS orphan_count,
    'Flags de reportes eliminados' AS issue
FROM incident_flags if_tbl
LEFT JOIN incident_reports ir ON if_tbl.inre_id = ir.inre_id
WHERE ir.inre_id IS NULL
HAVING COUNT(*) > 0;

-- Detalle de registros huérfanos
SELECT
    if_tbl.infl_id,
    if_tbl.inre_id AS missing_report_id,
    if_tbl.account_id,
    if_tbl.created_at
FROM incident_flags if_tbl
LEFT JOIN incident_reports ir ON if_tbl.inre_id = ir.inre_id
WHERE ir.inre_id IS NULL
ORDER BY if_tbl.created_at DESC
LIMIT 10;

\echo ''

-- 1.7 incident_logs -> incident_reports
\echo '7. incident_logs -> incident_reports:'
SELECT
    'incident_logs' AS table_name,
    COUNT(*) AS orphan_count,
    'Logs de reportes eliminados' AS issue
FROM incident_logs il
LEFT JOIN incident_reports ir ON il.inre_id = ir.inre_id
WHERE ir.inre_id IS NULL
HAVING COUNT(*) > 0;

-- Detalle de registros huérfanos
SELECT
    il.inlo_id,
    il.inre_id AS missing_report_id,
    il.account_id,
    il.action,
    il.created_at
FROM incident_logs il
LEFT JOIN incident_reports ir ON il.inre_id = ir.inre_id
WHERE ir.inre_id IS NULL
ORDER BY il.created_at DESC
LIMIT 10;

\echo ''
\echo '=========================================='
\echo 'RESUMEN DE DATOS HUÉRFANOS'
\echo '=========================================='

-- Resumen consolidado
SELECT
    orphan_table,
    orphan_count,
    CASE
        WHEN orphan_count = 0 THEN '✓ OK'
        WHEN orphan_count < 10 THEN '⚠ WARNING'
        ELSE '✗ CRÍTICO'
    END AS status,
    recommendation
FROM (
    SELECT 'account_cluster_saved' AS orphan_table,
           COUNT(*) AS orphan_count,
           'Eliminar favoritos huérfanos' AS recommendation
    FROM account_cluster_saved acs
    LEFT JOIN incident_clusters ic ON acs.incl_id = ic.incl_id
    WHERE ic.incl_id IS NULL

    UNION ALL

    SELECT 'account_history',
           COUNT(*),
           'Eliminar historial huérfano'
    FROM account_history ah
    LEFT JOIN incident_clusters ic ON ah.incl_id = ic.incl_id
    WHERE ah.incl_id IS NOT NULL AND ic.incl_id IS NULL

    UNION ALL

    SELECT 'notification_deliveries',
           COUNT(*),
           'Eliminar entregas huérfanas (CRÍTICO)'
    FROM notification_deliveries nd
    LEFT JOIN notifications n ON nd.noti_id = n.noti_id
    WHERE n.noti_id IS NULL

    UNION ALL

    SELECT 'incident_clusters',
           COUNT(*),
           'Corregir insu_id o eliminar cluster'
    FROM incident_clusters ic
    LEFT JOIN incident_subcategories isc ON ic.insu_id = isc.insu_id
    WHERE isc.insu_id IS NULL

    UNION ALL

    SELECT 'incident_reports',
           COUNT(*),
           'Corregir insu_id o eliminar reporte'
    FROM incident_reports ir
    LEFT JOIN incident_subcategories isc ON ir.insu_id = isc.insu_id
    WHERE isc.insu_id IS NULL

    UNION ALL

    SELECT 'incident_flags',
           COUNT(*),
           'Eliminar flags huérfanos'
    FROM incident_flags if_tbl
    LEFT JOIN incident_reports ir ON if_tbl.inre_id = ir.inre_id
    WHERE ir.inre_id IS NULL

    UNION ALL

    SELECT 'incident_logs',
           COUNT(*),
           'Eliminar logs huérfanos'
    FROM incident_logs il
    LEFT JOIN incident_reports ir ON il.inre_id = ir.inre_id
    WHERE ir.inre_id IS NULL
) summary
ORDER BY orphan_count DESC;

\echo ''
\echo '=========================================='
\echo 'PASO 2: VALIDAR DATOS NULL EN COORDENADAS'
\echo '=========================================='

-- 2.1 Verificar NULLs en incident_clusters
\echo ''
\echo 'Validando coordenadas en incident_clusters:'
SELECT
    COUNT(*) AS total_rows,
    COUNT(CASE WHEN center_latitude IS NULL THEN 1 END) AS null_latitude,
    COUNT(CASE WHEN center_longitude IS NULL THEN 1 END) AS null_longitude,
    COUNT(CASE WHEN center_latitude IS NULL OR center_longitude IS NULL THEN 1 END) AS total_nulls
FROM incident_clusters;

-- Detalle de registros con coordenadas NULL
SELECT
    incl_id,
    center_latitude,
    center_longitude,
    category_code,
    subcategory_code,
    created_at
FROM incident_clusters
WHERE center_latitude IS NULL
   OR center_longitude IS NULL
LIMIT 10;

\echo ''

-- 2.2 Verificar NULLs en account_favorite_locations
\echo 'Validando coordenadas en account_favorite_locations:'
SELECT
    COUNT(*) AS total_rows,
    COUNT(CASE WHEN latitude IS NULL THEN 1 END) AS null_latitude,
    COUNT(CASE WHEN longitude IS NULL THEN 1 END) AS null_longitude,
    COUNT(CASE WHEN latitude IS NULL OR longitude IS NULL THEN 1 END) AS total_nulls
FROM account_favorite_locations;

-- Detalle de registros con coordenadas NULL
SELECT
    afl_id,
    account_id,
    latitude,
    longitude,
    title,
    city
FROM account_favorite_locations
WHERE latitude IS NULL
   OR longitude IS NULL
LIMIT 10;

\echo ''

-- 2.3 Verificar NULLs en incident_reports
\echo 'Validando coordenadas en incident_reports:'
SELECT
    COUNT(*) AS total_rows,
    COUNT(CASE WHEN latitude IS NULL THEN 1 END) AS null_latitude,
    COUNT(CASE WHEN longitude IS NULL THEN 1 END) AS null_longitude,
    COUNT(CASE WHEN latitude IS NULL OR longitude IS NULL THEN 1 END) AS total_nulls
FROM incident_reports;

\echo ''
\echo '=========================================='
\echo 'PASO 3: VALIDAR CREATED_AT NULL'
\echo '=========================================='

-- 3.1 incident_clusters
\echo ''
\echo 'incident_clusters:'
SELECT COUNT(*) AS rows_with_null_created_at
FROM incident_clusters
WHERE created_at IS NULL;

-- 3.2 incident_reports
\echo 'incident_reports:'
SELECT COUNT(*) AS rows_with_null_created_at
FROM incident_reports
WHERE created_at IS NULL;

-- 3.3 account
\echo 'account:'
SELECT COUNT(*) AS rows_with_null_created_at
FROM account
WHERE created_at IS NULL;

-- 3.4 notifications
\echo 'notifications:'
SELECT COUNT(*) AS rows_with_null_created_at
FROM notifications
WHERE created_at IS NULL;

\echo ''
\echo '=========================================='
\echo 'CONCLUSIÓN DEL ANÁLISIS'
\echo '=========================================='
\echo ''
\echo 'Si el resumen muestra registros huérfanos (orphan_count > 0),'
\echo 'ejecutar los comandos de limpieza del PASO 4 a continuación.'
\echo ''
\echo 'Si NO hay datos huérfanos, proceder directamente con:'
\echo 'postgresql_migration_optimization.sql'
\echo ''

-- ============================================================
-- PASO 4: COMANDOS DE LIMPIEZA (COMENTADOS POR SEGURIDAD)
-- ============================================================
-- DESCOMENTAR Y EJECUTAR SOLO DESPUÉS DE REVISAR LOS RESULTADOS DEL ANÁLISIS

\echo ''
\echo '=========================================='
\echo 'PASO 4: LIMPIEZA DE DATOS HUÉRFANOS'
\echo '=========================================='
\echo ''
\echo 'ADVERTENCIA: Los siguientes comandos están COMENTADOS.'
\echo 'Revisar los resultados del análisis antes de ejecutar.'
\echo 'Descomentar solo las líneas necesarias.'
\echo ''

/*
-- ============================================================
-- 4.1 ELIMINAR FAVORITOS HUÉRFANOS
-- ============================================================
BEGIN;

-- Backup de datos antes de eliminar
CREATE TEMP TABLE backup_account_cluster_saved AS
SELECT acs.*
FROM account_cluster_saved acs
LEFT JOIN incident_clusters ic ON acs.incl_id = ic.incl_id
WHERE ic.incl_id IS NULL;

-- Eliminar registros huérfanos
DELETE FROM account_cluster_saved
WHERE acs_id IN (
    SELECT acs.acs_id
    FROM account_cluster_saved acs
    LEFT JOIN incident_clusters ic ON acs.incl_id = ic.incl_id
    WHERE ic.incl_id IS NULL
);

-- Verificar eliminación
SELECT COUNT(*) AS deleted_rows FROM backup_account_cluster_saved;

COMMIT;

-- ============================================================
-- 4.2 ELIMINAR HISTORIAL HUÉRFANO
-- ============================================================
BEGIN;

-- Backup
CREATE TEMP TABLE backup_account_history AS
SELECT ah.*
FROM account_history ah
LEFT JOIN incident_clusters ic ON ah.incl_id = ic.incl_id
WHERE ah.incl_id IS NOT NULL AND ic.incl_id IS NULL;

-- Eliminar
DELETE FROM account_history
WHERE his_id IN (
    SELECT ah.his_id
    FROM account_history ah
    LEFT JOIN incident_clusters ic ON ah.incl_id = ic.incl_id
    WHERE ah.incl_id IS NOT NULL AND ic.incl_id IS NULL
);

-- Verificar
SELECT COUNT(*) AS deleted_rows FROM backup_account_history;

COMMIT;

-- ============================================================
-- 4.3 ELIMINAR ENTREGAS DE NOTIFICACIONES HUÉRFANAS (CRÍTICO)
-- ============================================================
BEGIN;

-- Backup
CREATE TEMP TABLE backup_notification_deliveries AS
SELECT nd.*
FROM notification_deliveries nd
LEFT JOIN notifications n ON nd.noti_id = n.noti_id
WHERE n.noti_id IS NULL;

-- Eliminar
DELETE FROM notification_deliveries
WHERE node_id IN (
    SELECT nd.node_id
    FROM notification_deliveries nd
    LEFT JOIN notifications n ON nd.noti_id = n.noti_id
    WHERE n.noti_id IS NULL
);

-- Verificar
SELECT COUNT(*) AS deleted_rows FROM backup_notification_deliveries;

COMMIT;

-- ============================================================
-- 4.4 CORREGIR O ELIMINAR CLUSTERS CON SUBCATEGORÍA INVÁLIDA
-- ============================================================
-- OPCIÓN A: Asignar subcategoría por defecto (más seguro)
BEGIN;

-- Ver distribución de categorías
SELECT category_code, COUNT(*) AS count
FROM incident_clusters ic
LEFT JOIN incident_subcategories isc ON ic.insu_id = isc.insu_id
WHERE isc.insu_id IS NULL
GROUP BY category_code;

-- Asignar subcategoría "other" de cada categoría (ajustar según tu DB)
-- UPDATE incident_clusters
-- SET insu_id = (
--     SELECT isc.insu_id
--     FROM incident_subcategories isc
--     WHERE isc.code = 'other'
--     AND isc.inca_id = (
--         SELECT inca_id FROM incident_categories WHERE code = incident_clusters.category_code
--     )
--     LIMIT 1
-- )
-- WHERE insu_id IN (
--     SELECT ic.insu_id
--     FROM incident_clusters ic
--     LEFT JOIN incident_subcategories isc ON ic.insu_id = isc.insu_id
--     WHERE isc.insu_id IS NULL
-- );

-- OPCIÓN B: Eliminar clusters con datos inválidos (más riesgoso)
-- DELETE FROM incident_clusters
-- WHERE incl_id IN (
--     SELECT ic.incl_id
--     FROM incident_clusters ic
--     LEFT JOIN incident_subcategories isc ON ic.insu_id = isc.insu_id
--     WHERE isc.insu_id IS NULL
-- );

COMMIT;

-- ============================================================
-- 4.5 CORREGIR O ELIMINAR REPORTES CON SUBCATEGORÍA INVÁLIDA
-- ============================================================
-- Similar al paso 4.4, ajustar según necesidad

-- ============================================================
-- 4.6 ELIMINAR FLAGS HUÉRFANOS
-- ============================================================
BEGIN;

DELETE FROM incident_flags
WHERE infl_id IN (
    SELECT if_tbl.infl_id
    FROM incident_flags if_tbl
    LEFT JOIN incident_reports ir ON if_tbl.inre_id = ir.inre_id
    WHERE ir.inre_id IS NULL
);

COMMIT;

-- ============================================================
-- 4.7 ELIMINAR LOGS HUÉRFANOS
-- ============================================================
BEGIN;

DELETE FROM incident_logs
WHERE inlo_id IN (
    SELECT il.inlo_id
    FROM incident_logs il
    LEFT JOIN incident_reports ir ON il.inre_id = ir.inre_id
    WHERE ir.inre_id IS NULL
);

COMMIT;

-- ============================================================
-- 4.8 ACTUALIZAR CREATED_AT NULL CON DEFAULTS
-- ============================================================
BEGIN;

-- incident_clusters
UPDATE incident_clusters
SET created_at = NOW()
WHERE created_at IS NULL;

-- incident_reports
UPDATE incident_reports
SET created_at = NOW()
WHERE created_at IS NULL;

-- account
UPDATE account
SET created_at = NOW()
WHERE created_at IS NULL;

-- notifications
UPDATE notifications
SET created_at = NOW()
WHERE created_at IS NULL;

COMMIT;
*/

\echo ''
\echo '=========================================='
\echo 'FIN DEL PRE-MIGRATION CLEANUP'
\echo '=========================================='
\echo ''
\echo 'Próximos pasos:'
\echo '1. Revisar los resultados del análisis de datos huérfanos'
\echo '2. Si hay datos huérfanos, descomentar y ejecutar PASO 4'
\echo '3. Hacer backup completo de la base de datos'
\echo '4. Ejecutar postgresql_migration_optimization.sql'
\echo ''
