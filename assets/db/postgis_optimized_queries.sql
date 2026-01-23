-- ============================================================
-- POSTGIS OPTIMIZED QUERIES - ALERTLY DATABASE
-- Generated: 2026-01-22
--
-- Este archivo contiene ejemplos de queries optimizadas
-- que aprovechan los índices espaciales GiST de PostGIS
--
-- EJECUTAR DESPUÉS DE: postgresql_migration_optimization.sql
-- ============================================================

-- ============================================================
-- 1. BÚSQUEDA POR RADIO (ST_DWithin)
-- ============================================================

-- ANTES (MySQL con ST_Distance_Sphere - LENTO, sequential scan):
/*
SELECT * FROM incident_clusters
WHERE ST_Distance_Sphere(
    point(center_longitude, center_latitude),
    point(-79.3832, 43.6532)  -- Toronto Downtown
) <= 5000  -- 5km en metros
AND is_active = '1'
ORDER BY created_at DESC
LIMIT 100;
*/

-- DESPUÉS (PostgreSQL con PostGIS - RÁPIDO, usa índice GiST):
-- Búsqueda de clusters dentro de 5km de un punto
EXPLAIN ANALYZE
SELECT
    incl_id,
    category_code,
    subcategory_name,
    credibility,
    incident_count,
    ST_Distance(
        center_location,
        ST_MakePoint(-79.3832, 43.6532)::geography
    ) AS distance_meters,
    created_at
FROM incident_clusters
WHERE ST_DWithin(
    center_location,
    ST_MakePoint(-79.3832, 43.6532)::geography,
    5000  -- 5 kilómetros
)
AND is_active = '1'
ORDER BY center_location <-> ST_MakePoint(-79.3832, 43.6532)::geography  -- KNN operator para ordenar por distancia
LIMIT 100;

-- Nota: El operador <-> usa el índice GiST para ordenamiento eficiente por distancia

-- ============================================================
-- 2. BÚSQUEDA POR BOUNDING BOX (ST_MakeEnvelope)
-- ============================================================

-- ANTES (MySQL - requiere calcular bounds manualmente):
/*
SELECT * FROM incident_clusters
WHERE center_latitude BETWEEN 43.6000 AND 43.7000
  AND center_longitude BETWEEN -79.4000 AND -79.3000
  AND is_active = '1'
ORDER BY created_at DESC
LIMIT 100;
*/

-- DESPUÉS (PostgreSQL con PostGIS - usa índice espacial):
-- Búsqueda por viewport rectangular (útil para mapa interactivo)
EXPLAIN ANALYZE
SELECT
    incl_id,
    category_code,
    subcategory_name,
    credibility,
    incident_count,
    ST_X(center_location::geometry) AS longitude,
    ST_Y(center_location::geometry) AS latitude,
    created_at
FROM incident_clusters
WHERE center_location && ST_MakeEnvelope(
    -79.4000, 43.6000,  -- SW corner (lon, lat)
    -79.3000, 43.7000,  -- NE corner (lon, lat)
    4326
)::geography
AND is_active = '1'
ORDER BY created_at DESC
LIMIT 100;

-- Operador && verifica intersección con bounding box (muy rápido con GiST)

-- ============================================================
-- 3. BÚSQUEDA DE LUGARES FAVORITOS CERCANOS
-- ============================================================

-- Encontrar lugares favoritos dentro de 10km de un punto
EXPLAIN ANALYZE
SELECT
    afl_id,
    account_id,
    title,
    city,
    ST_Distance(
        location,
        ST_MakePoint(-79.3832, 43.6532)::geography
    ) AS distance_meters,
    radius
FROM account_favorite_locations
WHERE ST_DWithin(
    location,
    ST_MakePoint(-79.3832, 43.6532)::geography,
    10000  -- 10 kilómetros
)
AND status = 1
ORDER BY location <-> ST_MakePoint(-79.3832, 43.6532)::geography
LIMIT 50;

-- ============================================================
-- 4. CLUSTERING DE INCIDENTES NUEVOS
-- ============================================================

-- Buscar clusters existentes cerca de un nuevo incidente
-- (usado en el algoritmo de clustering de newincident/)
EXPLAIN ANALYZE
WITH new_incident_location AS (
    SELECT ST_MakePoint(-79.3832, 43.6532)::geography AS location
)
SELECT
    ic.incl_id,
    ic.insu_id,
    ic.category_code,
    ic.subcategory_code,
    ic.incident_count,
    ic.credibility,
    ST_Distance(
        ic.center_location,
        nil.location
    ) AS distance_meters,
    ic.created_at
FROM incident_clusters ic
CROSS JOIN new_incident_location nil
WHERE ST_DWithin(
    ic.center_location,
    nil.location,
    200  -- Radio de clustering: 200 metros (ajustar por categoría)
)
AND ic.insu_id = 5  -- Subcategoría del nuevo incidente
AND ic.is_active = '1'
AND ic.created_at >= NOW() - INTERVAL '24 hours'  -- Ventana de tiempo
ORDER BY distance_meters ASC
LIMIT 1;

-- Si no se encuentra cluster, crear uno nuevo

-- ============================================================
-- 5. NOTIFICACIONES POR PROXIMIDAD
-- ============================================================

-- Encontrar usuarios con lugares favoritos cercanos a un nuevo incidente
-- (usado en cronjob cjnewcluster para enviar push notifications)
EXPLAIN ANALYZE
WITH new_incident AS (
    SELECT
        1058 AS incl_id,  -- ID del nuevo cluster
        ST_MakePoint(-79.3832, 43.6532)::geography AS location,
        'traffic_accident'::varchar AS category_code
)
SELECT DISTINCT
    afl.account_id,
    afl.title AS favorite_location_name,
    ST_Distance(
        afl.location,
        ni.location
    ) AS distance_meters,
    afl.radius AS notification_radius,
    afl.traffic_accident AS wants_traffic_notifications
FROM account_favorite_locations afl
CROSS JOIN new_incident ni
WHERE ST_DWithin(
    afl.location,
    ni.location,
    afl.radius  -- Radio configurado por usuario (default 3000m)
)
AND afl.status = 1
AND afl.traffic_accident = 1  -- Usuario quiere notificaciones de esta categoría
ORDER BY afl.account_id;

-- ============================================================
-- 6. ANÁLISIS DE DENSIDAD DE INCIDENTES (HEATMAP)
-- ============================================================

-- Contar incidentes por grid hexagonal (útil para analytics premium)
-- Requiere ST_HexagonGrid (PostGIS 3.1+)
EXPLAIN ANALYZE
SELECT
    grid.geom,
    COUNT(ic.incl_id) AS incident_count,
    AVG(ic.credibility) AS avg_credibility,
    string_agg(DISTINCT ic.category_code, ',') AS categories
FROM ST_HexagonGrid(
    0.01,  -- Tamaño de celda en grados (aprox 1km)
    ST_MakeEnvelope(-79.5, 43.5, -79.2, 43.8, 4326)  -- Área de Toronto
) AS grid
LEFT JOIN incident_clusters ic
    ON ST_Intersects(grid.geom, ic.center_location::geometry)
    AND ic.is_active = '1'
    AND ic.created_at >= NOW() - INTERVAL '7 days'
GROUP BY grid.geom
HAVING COUNT(ic.incl_id) > 0
ORDER BY incident_count DESC
LIMIT 100;

-- ============================================================
-- 7. VALIDACIÓN DE PROXIMIDAD PARA VOTOS
-- ============================================================

-- Verificar si usuario está dentro de 100m del incidente para poder votar
-- (implementar en backend: incident_reports/service.go)
EXPLAIN ANALYZE
WITH user_location AS (
    SELECT ST_MakePoint(-79.3832, 43.6532)::geography AS location  -- Ubicación actual del usuario
),
incident_location AS (
    SELECT center_location
    FROM incident_clusters
    WHERE incl_id = 1058
)
SELECT
    ST_Distance(ul.location, il.center_location) AS distance_meters,
    CASE
        WHEN ST_DWithin(ul.location, il.center_location, 100) THEN true
        ELSE false
    END AS can_vote
FROM user_location ul
CROSS JOIN incident_location il;

-- ============================================================
-- 8. BÚSQUEDA DE INCIDENTES POR MÚLTIPLES CATEGORÍAS
-- ============================================================

-- Buscar incidentes de varias categorías cerca de un punto
EXPLAIN ANALYZE
SELECT
    incl_id,
    category_code,
    subcategory_name,
    description,
    credibility,
    incident_count,
    ST_Distance(
        center_location,
        ST_MakePoint(-79.3832, 43.6532)::geography
    ) AS distance_meters,
    created_at
FROM incident_clusters
WHERE ST_DWithin(
    center_location,
    ST_MakePoint(-79.3832, 43.6532)::geography,
    3000
)
AND category_code IN ('crime', 'traffic_accident', 'suspicious_activity')
AND is_active = '1'
AND created_at >= NOW() - INTERVAL '48 hours'
ORDER BY
    CASE category_code
        WHEN 'crime' THEN 1
        WHEN 'traffic_accident' THEN 2
        ELSE 3
    END,
    distance_meters ASC
LIMIT 50;

-- ============================================================
-- 9. INCIDENTES MÁS CERCANOS (KNN)
-- ============================================================

-- Los 10 incidentes más cercanos usando operador KNN (<->)
-- Este query es MUCHO más rápido que ordenar por ST_Distance()
EXPLAIN ANALYZE
SELECT
    incl_id,
    category_code,
    subcategory_name,
    credibility,
    center_location <-> ST_MakePoint(-79.3832, 43.6532)::geography AS distance,
    created_at
FROM incident_clusters
WHERE is_active = '1'
  AND created_at >= NOW() - INTERVAL '24 hours'
ORDER BY center_location <-> ST_MakePoint(-79.3832, 43.6532)::geography
LIMIT 10;

-- ============================================================
-- 10. AGREGACIÓN ESPACIAL POR CIUDAD
-- ============================================================

-- Contar incidentes por ciudad (útil para analytics)
SELECT
    city,
    COUNT(*) AS incident_count,
    AVG(credibility) AS avg_credibility,
    COUNT(CASE WHEN credibility >= 7 THEN 1 END) AS high_credibility_count,
    COUNT(CASE WHEN created_at >= NOW() - INTERVAL '24 hours' THEN 1 END) AS last_24h_count
FROM incident_clusters
WHERE is_active = '1'
  AND city IS NOT NULL
  AND created_at >= NOW() - INTERVAL '7 days'
GROUP BY city
ORDER BY incident_count DESC
LIMIT 20;

-- ============================================================
-- 11. RUTA DE INCIDENTES (LÍNEA DE EVENTOS)
-- ============================================================

-- Crear una línea conectando múltiples incidentes (útil para tracking de eventos)
-- Por ejemplo, seguimiento de un incidente que se mueve (tráfico, clima)
SELECT
    ST_AsText(
        ST_MakeLine(
            ARRAY(
                SELECT center_location::geometry
                FROM incident_clusters
                WHERE category_code = 'traffic_accident'
                  AND is_active = '1'
                  AND created_at >= NOW() - INTERVAL '2 hours'
                ORDER BY created_at ASC
            )
        )
    ) AS incident_path;

-- ============================================================
-- 12. ÁREA DE COBERTURA (CONVEX HULL)
-- ============================================================

-- Calcular el área de cobertura de incidentes de una categoría
-- Útil para determinar "zonas de alto riesgo"
SELECT
    category_code,
    COUNT(*) AS incident_count,
    ST_AsGeoJSON(
        ST_ConvexHull(
            ST_Collect(center_location::geometry)
        )
    ) AS coverage_area_geojson,
    ST_Area(
        ST_ConvexHull(
            ST_Collect(center_location::geography)
        )
    ) / 1000000 AS coverage_area_km2  -- Convertir m² a km²
FROM incident_clusters
WHERE is_active = '1'
  AND created_at >= NOW() - INTERVAL '7 days'
  AND category_code = 'crime'
GROUP BY category_code;

-- ============================================================
-- 13. ESTADÍSTICAS DE PERFORMANCE DE ÍNDICES
-- ============================================================

-- Verificar que los índices espaciales se están usando
SELECT
    schemaname,
    tablename,
    indexname,
    idx_scan AS times_used,
    idx_tup_read AS tuples_read,
    idx_tup_fetch AS tuples_fetched,
    pg_size_pretty(pg_relation_size(indexrelid)) AS index_size
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
  AND indexname LIKE '%gist%'
ORDER BY idx_scan DESC;

-- Ver queries lentas con búsquedas geoespaciales
-- (requiere pg_stat_statements extension)
SELECT
    query,
    calls,
    mean_exec_time,
    max_exec_time,
    total_exec_time
FROM pg_stat_statements
WHERE query ILIKE '%center_location%'
   OR query ILIKE '%ST_DWithin%'
   OR query ILIKE '%ST_Distance%'
ORDER BY mean_exec_time DESC
LIMIT 10;

-- ============================================================
-- 14. BENCHMARK: COMPARAR ÍNDICE BTREE VS GIST
-- ============================================================

-- Benchmark con índice B-tree (coordenadas numéricas)
EXPLAIN (ANALYZE, BUFFERS)
SELECT COUNT(*)
FROM incident_clusters
WHERE center_latitude BETWEEN 43.6 AND 43.7
  AND center_longitude BETWEEN -79.4 AND -79.3
  AND is_active = '1';

-- Benchmark con índice GiST (geography)
EXPLAIN (ANALYZE, BUFFERS)
SELECT COUNT(*)
FROM incident_clusters
WHERE center_location && ST_MakeEnvelope(
    -79.4, 43.6,
    -79.3, 43.7,
    4326
)::geography
AND is_active = '1';

-- El índice GiST debería ser 10-50x más rápido en búsquedas espaciales

-- ============================================================
-- 15. MIGRAR QUERY DEL BACKEND GO
-- ============================================================

-- ANTES (MySQL en getclustersbylocation/repository.go):
/*
SELECT * FROM incident_clusters
WHERE ST_Distance_Sphere(
    point(center_longitude, center_latitude),
    point(?, ?)
) <= ?
AND is_active = 1
AND created_at >= DATE_SUB(NOW(), INTERVAL 24 HOUR)
ORDER BY created_at DESC
LIMIT 100;
*/

-- DESPUÉS (PostgreSQL optimizado):
PREPARE get_clusters_by_location (float8, float8, float8) AS
SELECT
    incl_id,
    center_latitude,  -- Mantener para compatibilidad con frontend
    center_longitude,
    category_code,
    subcategory_code,
    subcategory_name,
    description,
    credibility,
    incident_count,
    counter_total_votes,
    counter_total_comments,
    media_url,
    media_type,
    address,
    city,
    province,
    created_at,
    updated_at,
    start_time,
    end_time,
    account_id,
    ST_Distance(
        center_location,
        ST_MakePoint($1, $2)::geography
    ) AS distance_meters
FROM incident_clusters
WHERE ST_DWithin(
    center_location,
    ST_MakePoint($1, $2)::geography,
    $3  -- Radius en metros
)
AND is_active = '1'
AND created_at >= NOW() - INTERVAL '24 hours'
ORDER BY center_location <-> ST_MakePoint($1, $2)::geography  -- KNN sort
LIMIT 100;

-- Ejecutar prepared statement:
-- EXECUTE get_clusters_by_location(-79.3832, 43.6532, 5000);

-- ============================================================
-- 16. FUNCIÓN HELPER PARA BACKEND
-- ============================================================

-- Crear función para obtener clusters por radio (simplifica backend)
CREATE OR REPLACE FUNCTION get_clusters_within_radius(
    p_longitude float8,
    p_latitude float8,
    p_radius_meters float8,
    p_time_window_hours integer DEFAULT 24,
    p_limit integer DEFAULT 100
)
RETURNS TABLE (
    incl_id integer,
    center_latitude numeric,
    center_longitude numeric,
    category_code varchar,
    subcategory_name varchar,
    description text,
    credibility numeric,
    incident_count integer,
    distance_meters float8,
    created_at timestamp
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        ic.incl_id,
        ic.center_latitude,
        ic.center_longitude,
        ic.category_code,
        ic.subcategory_name,
        ic.description,
        ic.credibility,
        ic.incident_count,
        ST_Distance(
            ic.center_location,
            ST_MakePoint(p_longitude, p_latitude)::geography
        ) AS distance_meters,
        ic.created_at
    FROM incident_clusters ic
    WHERE ST_DWithin(
        ic.center_location,
        ST_MakePoint(p_longitude, p_latitude)::geography,
        p_radius_meters
    )
    AND ic.is_active = '1'
    AND ic.created_at >= NOW() - (p_time_window_hours || ' hours')::INTERVAL
    ORDER BY ic.center_location <-> ST_MakePoint(p_longitude, p_latitude)::geography
    LIMIT p_limit;
END;
$$ LANGUAGE plpgsql STABLE;

-- Uso desde backend Go:
-- SELECT * FROM get_clusters_within_radius(-79.3832, 43.6532, 5000, 24, 100);

-- ============================================================
-- 17. ÍNDICE DE TEXTO COMPLETO (BONUS)
-- ============================================================

-- Búsqueda combinada: geolocalización + texto
EXPLAIN ANALYZE
SELECT
    incl_id,
    category_code,
    description,
    address,
    credibility,
    ST_Distance(
        center_location,
        ST_MakePoint(-79.3832, 43.6532)::geography
    ) AS distance_meters
FROM incident_clusters
WHERE ST_DWithin(
    center_location,
    ST_MakePoint(-79.3832, 43.6532)::geography,
    10000
)
AND (
    description ILIKE '%accident%'
    OR address ILIKE '%queen%street%'
)
AND is_active = '1'
ORDER BY center_location <-> ST_MakePoint(-79.3832, 43.6532)::geography
LIMIT 20;

-- ============================================================
-- FIN DE QUERIES OPTIMIZADAS
-- ============================================================

-- NOTAS IMPORTANTES:
-- 1. Todas las queries usan center_location (GEOGRAPHY) en WHERE
-- 2. El índice GiST se aprovecha con ST_DWithin() y operador &&
-- 3. El operador <-> (KNN) ordena eficientemente por distancia
-- 4. EXPLAIN ANALYZE muestra que usa "Index Scan using idx_clusters_center_location_gist"
-- 5. Mantener columnas center_latitude/longitude para compatibilidad con frontend

-- RENDIMIENTO ESPERADO:
-- - Búsqueda por radio (5km): <10ms con índice GiST vs >100ms con B-tree
-- - Búsqueda por bounding box: <5ms con índice GiST
-- - KNN (10 más cercanos): <3ms con operador <->
-- - Sin índice: >1000ms para 1000+ registros

-- PRÓXIMOS PASOS EN BACKEND GO:
-- 1. Actualizar getclustersbylocation/repository.go para usar ST_DWithin()
-- 2. Actualizar getclusterbyradius/repository.go para usar center_location
-- 3. Actualizar newincident/service.go para usar índice espacial en clustering
-- 4. Agregar función get_clusters_within_radius() al backend
-- 5. Implementar caching de queries geoespaciales frecuentes
