-- ============================================
-- SCRIPT DE CORRECCIÓN DE SUBCATEGORÍAS DEL BOT
-- Fecha: 2025-11-26
-- Propósito: Corregir códigos de subcategoría inválidos
-- ============================================

-- PASO 1: VERIFICAR CÓDIGOS INVÁLIDOS
-- Ejecutar primero para ver qué hay que corregir
SELECT
    'CÓDIGOS INVÁLIDOS ENCONTRADOS' as status,
    subcategory_code,
    category_code,
    COUNT(*) as total_clusters
FROM incident_clusters
WHERE account_id = 1
  AND is_active = 1
  AND subcategory_code NOT IN (
    'theft', 'robbery', 'assault', 'homicide', 'fraud',
    'vehicle_collision', 'pedestrian_nvolvement', 'hit_and_run', 'multi_vehicle_ileup',
    'cardiac_arrest', 'stroke', 'trauma_Injury', 'overdose_poisoning', 'other_medical_emergency',
    'residential_fire', 'wildfire', 'vehicle_fire', 'other_fire_incident',
    'graffiti', 'vehicle_vandalism', 'public_property_damage',
    'suspicious_person', 'suspicious_vehicle', 'unusual_behavior', 'other_suspicious_activity',
    'road_damage_potholes', 'streetlight_traffic_signal_failure', 'sidewalk_pathway_damage',
    'public_utility_issues', 'structural_damage'
  )
GROUP BY subcategory_code, category_code
ORDER BY total_clusters DESC;

-- PASO 2: CREAR BACKUPS
CREATE TABLE IF NOT EXISTS incident_clusters_backup_20251126 AS
SELECT * FROM incident_clusters WHERE account_id = 1;

CREATE TABLE IF NOT EXISTS incident_reports_backup_20251126 AS
SELECT * FROM incident_reports WHERE account_id = 1;

-- PASO 3: INICIAR TRANSACCIÓN
START TRANSACTION;

-- PASO 4: CORREGIR incident_clusters

-- Corregir traffic_accident genérico → vehicle_collision
UPDATE incident_clusters
SET subcategory_code = 'vehicle_collision'
WHERE account_id = 1
  AND category_code = 'traffic_accident'
  AND subcategory_code = 'traffic_accident';

-- Corregir crime genérico → theft
UPDATE incident_clusters
SET subcategory_code = 'theft'
WHERE account_id = 1
  AND category_code = 'crime'
  AND subcategory_code = 'crime';

-- Corregir building_fire → residential_fire
UPDATE incident_clusters
SET subcategory_code = 'residential_fire'
WHERE account_id = 1
  AND category_code = 'fire_incident'
  AND subcategory_code = 'building_fire';

-- Corregir fire_incident genérico → other_fire_incident
UPDATE incident_clusters
SET subcategory_code = 'other_fire_incident'
WHERE account_id = 1
  AND category_code = 'fire_incident'
  AND subcategory_code = 'fire_incident';

-- Corregir medical_emergency genérico → other_medical_emergency
UPDATE incident_clusters
SET subcategory_code = 'other_medical_emergency'
WHERE account_id = 1
  AND category_code = 'medical_emergency'
  AND subcategory_code = 'medical_emergency';

-- Corregir utility_issues → public_utility_issues
UPDATE incident_clusters
SET subcategory_code = 'public_utility_issues'
WHERE account_id = 1
  AND category_code = 'infrastructure_issues'
  AND subcategory_code = 'utility_issues';

-- PASO 5: CORREGIR incident_reports

-- Corregir traffic_accident genérico → vehicle_collision
UPDATE incident_reports
SET subcategory_code = 'vehicle_collision'
WHERE account_id = 1
  AND category_code = 'traffic_accident'
  AND subcategory_code = 'traffic_accident';

-- Corregir crime genérico → theft
UPDATE incident_reports
SET subcategory_code = 'theft'
WHERE account_id = 1
  AND category_code = 'crime'
  AND subcategory_code = 'crime';

-- Corregir building_fire → residential_fire
UPDATE incident_reports
SET subcategory_code = 'residential_fire'
WHERE account_id = 1
  AND category_code = 'fire_incident'
  AND subcategory_code = 'building_fire';

-- Corregir fire_incident genérico → other_fire_incident
UPDATE incident_reports
SET subcategory_code = 'other_fire_incident'
WHERE account_id = 1
  AND category_code = 'fire_incident'
  AND subcategory_code = 'fire_incident';

-- Corregir medical_emergency genérico → other_medical_emergency
UPDATE incident_reports
SET subcategory_code = 'other_medical_emergency'
WHERE account_id = 1
  AND category_code = 'medical_emergency'
  AND subcategory_code = 'medical_emergency';

-- Corregir utility_issues → public_utility_issues
UPDATE incident_reports
SET subcategory_code = 'public_utility_issues'
WHERE account_id = 1
  AND category_code = 'infrastructure_issues'
  AND subcategory_code = 'utility_issues';

-- PASO 6: VERIFICAR RESULTADOS
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

-- PASO 7: VERIFICAR QUE NO QUEDAN CÓDIGOS INVÁLIDOS
SELECT
    CASE
        WHEN COUNT(*) = 0 THEN '✅ NO HAY CÓDIGOS INVÁLIDOS - TODO CORRECTO'
        ELSE '❌ TODAVÍA HAY CÓDIGOS INVÁLIDOS'
    END as verificacion_final,
    COUNT(*) as total_invalidos
FROM incident_clusters
WHERE account_id = 1
  AND is_active = 1
  AND subcategory_code NOT IN (
    'theft', 'robbery', 'assault', 'homicide', 'fraud',
    'vehicle_collision', 'pedestrian_nvolvement', 'hit_and_run', 'multi_vehicle_ileup',
    'cardiac_arrest', 'stroke', 'trauma_Injury', 'overdose_poisoning', 'other_medical_emergency',
    'residential_fire', 'wildfire', 'vehicle_fire', 'other_fire_incident',
    'graffiti', 'vehicle_vandalism', 'public_property_damage',
    'suspicious_person', 'suspicious_vehicle', 'unusual_behavior', 'other_suspicious_activity',
    'road_damage_potholes', 'streetlight_traffic_signal_failure', 'sidewalk_pathway_damage',
    'public_utility_issues', 'structural_damage'
  );

-- PASO 8: COMMIT O ROLLBACK
-- Si la verificación está OK, ejecutar:
COMMIT;

-- Si hay problemas, ejecutar:
-- ROLLBACK;
