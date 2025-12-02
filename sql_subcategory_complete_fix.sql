-- ============================================
-- SCRIPT COMPLETO DE CORRECCIÓN - SUBCATEGORÍAS DEL BOT
-- Fecha: 2025-11-27
-- Propósito: Corregir TODOS los códigos de subcategoría inválidos
-- ============================================

-- PASO 1: VERIFICAR CÓDIGOS INVÁLIDOS ANTES DEL FIX
SELECT
    'CÓDIGOS INVÁLIDOS ENCONTRADOS' as status,
    subcategory_code,
    category_code,
    COUNT(*) as total_clusters
FROM incident_clusters
WHERE account_id = 1
  AND is_active = 1
  AND subcategory_code NOT IN (
    SELECT code FROM incident_subcategories
  )
GROUP BY subcategory_code, category_code
ORDER BY total_clusters DESC;

-- PASO 2: CREAR BACKUPS
CREATE TABLE IF NOT EXISTS incident_clusters_backup_20251127 AS
SELECT * FROM incident_clusters WHERE account_id = 1;

CREATE TABLE IF NOT EXISTS incident_reports_backup_20251127 AS
SELECT * FROM incident_reports WHERE account_id = 1;

-- PASO 3: INICIAR TRANSACCIÓN
START TRANSACTION;

-- ============================================
-- FIX 1: vehicle_collision → single_vehicle_accident
-- ============================================
UPDATE incident_clusters
SET subcategory_code = 'single_vehicle_accident'
WHERE account_id = 1 AND subcategory_code = 'vehicle_collision';

UPDATE incident_reports
SET subcategory_code = 'single_vehicle_accident'
WHERE account_id = 1 AND subcategory_code = 'vehicle_collision';

-- ============================================
-- FIX 2: other_fire_incident → fire_incident
-- ============================================
UPDATE incident_clusters
SET subcategory_code = 'fire_incident'
WHERE account_id = 1 AND subcategory_code = 'other_fire_incident';

UPDATE incident_reports
SET subcategory_code = 'fire_incident'
WHERE account_id = 1 AND subcategory_code = 'other_fire_incident';

-- ============================================
-- FIX 3: multi_vehicle_ileup → traffic_accident
-- ============================================
UPDATE incident_clusters
SET subcategory_code = 'traffic_accident'
WHERE account_id = 1 AND subcategory_code = 'multi_vehicle_ileup';

UPDATE incident_reports
SET subcategory_code = 'traffic_accident'
WHERE account_id = 1 AND subcategory_code = 'multi_vehicle_ileup';

-- ============================================
-- FIX 4: overdose_poisoning → medical_emergency
-- ============================================
UPDATE incident_clusters
SET subcategory_code = 'medical_emergency'
WHERE account_id = 1 AND subcategory_code = 'overdose_poisoning';

UPDATE incident_reports
SET subcategory_code = 'medical_emergency'
WHERE account_id = 1 AND subcategory_code = 'overdose_poisoning';

-- ============================================
-- FIX 5: stroke → medical_emergency
-- ============================================
UPDATE incident_clusters
SET subcategory_code = 'medical_emergency'
WHERE account_id = 1 AND subcategory_code = 'stroke';

UPDATE incident_reports
SET subcategory_code = 'medical_emergency'
WHERE account_id = 1 AND subcategory_code = 'stroke';

-- PASO 4: VERIFICAR RESULTADOS
SELECT
    '✅ CLUSTERS CORREGIDOS' as tabla,
    subcategory_code,
    category_code,
    COUNT(*) as total
FROM incident_clusters
WHERE account_id = 1 AND is_active = 1
GROUP BY subcategory_code, category_code

UNION ALL

SELECT
    '✅ REPORTS CORREGIDOS',
    subcategory_code,
    category_code,
    COUNT(*)
FROM incident_reports
WHERE account_id = 1 AND is_active = 1
GROUP BY subcategory_code, category_code
ORDER BY tabla, total DESC;

-- PASO 5: VERIFICACIÓN FINAL - NO DEBE RETORNAR NADA
SELECT
    'CODIGO NO EXISTE EN DB' as problema,
    bot_codes.code as codigo_invalido,
    bot_codes.categoria,
    bot_codes.total_incidentes
FROM (
    SELECT DISTINCT subcategory_code as code, category_code as categoria, COUNT(*) as total_incidentes
    FROM incident_clusters
    WHERE account_id = 1
      AND is_active = 1
    GROUP BY subcategory_code, category_code
) bot_codes
LEFT JOIN incident_subcategories isub ON bot_codes.code = isub.code
WHERE isub.code IS NULL
ORDER BY bot_codes.total_incidentes DESC;

-- PASO 6: COMMIT O ROLLBACK
-- Si la verificación está OK (retorna vacío), ejecutar:
COMMIT;

-- Si hay problemas, ejecutar:
-- ROLLBACK;

-- ============================================
-- RESUMEN DE CAMBIOS APLICADOS
-- ============================================
/*
1. vehicle_collision → single_vehicle_accident (2 clusters, 3 reports)
2. other_fire_incident → fire_incident (2 clusters)
3. multi_vehicle_ileup → traffic_accident (1 cluster)
4. overdose_poisoning → medical_emergency (1 cluster)
5. stroke → medical_emergency (1 cluster)

TOTAL: 7 clusters + 3 reports corregidos
*/
