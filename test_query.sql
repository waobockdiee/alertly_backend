-- Test Query para verificar que la sintaxis SQL es correcta
-- Simula la query que se construye en getclustersbylocation

USE alertly;

-- Test 1: Query básica sin categorías
SELECT
    t1.incl_id, t1.center_latitude, t1.center_longitude, t1.insu_id, t1.category_code, t1.subcategory_code
FROM incident_clusters t1
WHERE t1.center_latitude BETWEEN 40.0 AND 50.0
  AND t1.center_longitude BETWEEN -120.0 AND -110.0
  AND DATE(t1.start_time) <= '2025-08-08'
  AND DATE(t1.end_time) >= '2025-08-07'
  AND (0 = 0 OR t1.insu_id = 0)
  AND t1.is_active = 1
ORDER BY t1.created_at DESC
LIMIT 100;

-- Test 2: Query con categorías específicas
SELECT
    t1.incl_id, t1.center_latitude, t1.center_longitude, t1.insu_id, t1.category_code, t1.subcategory_code
FROM incident_clusters t1
WHERE t1.center_latitude BETWEEN 40.0 AND 50.0
  AND t1.center_longitude BETWEEN -120.0 AND -110.0
  AND DATE(t1.start_time) <= '2025-08-08'
  AND DATE(t1.end_time) >= '2025-08-07'
  AND (0 = 0 OR t1.insu_id = 0)
  AND t1.is_active = 1
  AND t1.category_code IN ('crime', 'traffic_accident', 'medical_emergency')
ORDER BY t1.created_at DESC
LIMIT 100;

-- Test 3: Query con insu_id específico
SELECT
    t1.incl_id, t1.center_latitude, t1.center_longitude, t1.insu_id, t1.category_code, t1.subcategory_code
FROM incident_clusters t1
WHERE t1.center_latitude BETWEEN 40.0 AND 50.0
  AND t1.center_longitude BETWEEN -120.0 AND -110.0
  AND DATE(t1.start_time) <= '2025-08-08'
  AND DATE(t1.end_time) >= '2025-08-07'
  AND (1 = 0 OR t1.insu_id = 1)
  AND t1.is_active = 1
ORDER BY t1.created_at DESC
LIMIT 100;

-- Verificar que las tablas existen y tienen datos
SELECT COUNT(*) as total_clusters FROM incident_clusters WHERE is_active = 1;
SELECT COUNT(*) as total_categories FROM incident_categories;
SELECT DISTINCT category_code FROM incident_clusters WHERE category_code IS NOT NULL LIMIT 10;
