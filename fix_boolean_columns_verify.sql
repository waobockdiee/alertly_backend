-- =====================================================
-- VERIFICACIÓN: Migración de Columnas Booleanas
-- Fecha: 2026-01-18
-- =====================================================

\echo ''
\echo '=================================================='
\echo 'VERIFICACIÓN DE MIGRACIÓN'
\echo '=================================================='
\echo ''

-- =====================================================
-- 1. VERIFICAR TIPOS DE DATOS EN account
-- =====================================================

\echo '1. Tipos de datos en tabla account:'
\echo ''

SELECT
    column_name,
    data_type,
    column_default,
    is_nullable
FROM information_schema.columns
WHERE table_name = 'account'
  AND column_name IN ('is_private_profile', 'receive_notifications', 'has_finished_tutorial', 'has_watch_new_incident_tutorial')
ORDER BY column_name;

\echo ''

-- =====================================================
-- 2. VERIFICAR TIPOS DE DATOS EN account_favorite_locations
-- =====================================================

\echo '2. Tipos de datos en tabla account_favorite_locations:'
\echo ''

SELECT
    column_name,
    data_type,
    column_default,
    is_nullable
FROM information_schema.columns
WHERE table_name = 'account_favorite_locations'
  AND column_name IN (
      'crime', 'traffic_accident', 'medical_emergency', 'fire_incident',
      'vandalism', 'suspicious_activity', 'infrastructure_issues', 'extreme_weather',
      'community_events', 'dangerous_wildlife_sighting', 'positive_actions', 'lost_pet'
  )
ORDER BY column_name;

\echo ''

-- =====================================================
-- 3. CONTAR COLUMNAS CON TIPO INCORRECTO
-- =====================================================

\echo '3. Verificación de tipos (debe ser 0 para ambas tablas):'
\echo ''

SELECT
    'account' AS tabla,
    COUNT(*) AS columnas_con_tipo_incorrecto
FROM information_schema.columns
WHERE table_name = 'account'
  AND column_name IN ('is_private_profile', 'receive_notifications', 'has_finished_tutorial', 'has_watch_new_incident_tutorial')
  AND data_type != 'smallint'

UNION ALL

SELECT
    'account_favorite_locations' AS tabla,
    COUNT(*) AS columnas_con_tipo_incorrecto
FROM information_schema.columns
WHERE table_name = 'account_favorite_locations'
  AND column_name IN (
      'crime', 'traffic_accident', 'medical_emergency', 'fire_incident',
      'vandalism', 'suspicious_activity', 'infrastructure_issues', 'extreme_weather',
      'community_events', 'dangerous_wildlife_sighting', 'positive_actions', 'lost_pet'
  )
  AND data_type != 'smallint';

\echo ''

-- =====================================================
-- 4. VERIFICAR VALORES (DEBEN SER 0 O 1)
-- =====================================================

\echo '4. Verificación de valores inválidos en account (debe ser 0):'
\echo ''

SELECT COUNT(*) AS registros_con_valores_invalidos
FROM account
WHERE is_private_profile NOT IN (0, 1)
   OR receive_notifications NOT IN (0, 1)
   OR has_finished_tutorial NOT IN (0, 1)
   OR has_watch_new_incident_tutorial NOT IN (0, 1);

\echo ''
\echo '5. Verificación de valores inválidos en account_favorite_locations (debe ser 0):'
\echo ''

SELECT COUNT(*) AS registros_con_valores_invalidos
FROM account_favorite_locations
WHERE crime NOT IN (0, 1)
   OR traffic_accident NOT IN (0, 1)
   OR medical_emergency NOT IN (0, 1)
   OR fire_incident NOT IN (0, 1)
   OR vandalism NOT IN (0, 1)
   OR suspicious_activity NOT IN (0, 1)
   OR infrastructure_issues NOT IN (0, 1)
   OR extreme_weather NOT IN (0, 1)
   OR community_events NOT IN (0, 1)
   OR dangerous_wildlife_sighting NOT IN (0, 1)
   OR positive_actions NOT IN (0, 1)
   OR lost_pet NOT IN (0, 1);

\echo ''

-- =====================================================
-- 6. MUESTRA DE DATOS MIGRADOS
-- =====================================================

\echo '6. Muestra de datos en account (primeros 5 registros):'
\echo ''

SELECT
    account_id,
    is_premium,
    receive_notifications,
    is_private_profile,
    has_finished_tutorial,
    has_watch_new_incident_tutorial
FROM account
LIMIT 5;

\echo ''
\echo '7. Muestra de datos en account_favorite_locations (primeros 3 registros):'
\echo ''

SELECT
    afl_id,
    account_id,
    crime,
    traffic_accident,
    medical_emergency,
    fire_incident,
    vandalism,
    suspicious_activity
FROM account_favorite_locations
LIMIT 3;

\echo ''

-- =====================================================
-- 7. TEST DE QUERIES NUMÉRICAS (COMO EN GO)
-- =====================================================

\echo '8. Test de queries numéricas (compatibilidad con código Go):'
\echo ''

SELECT
    'is_premium = 1' AS query_test,
    COUNT(*) AS resultado
FROM account
WHERE is_premium = 1

UNION ALL

SELECT
    'receive_notifications = 1' AS query_test,
    COUNT(*) AS resultado
FROM account
WHERE receive_notifications = 1

UNION ALL

SELECT
    'is_private_profile = 0' AS query_test,
    COUNT(*) AS resultado
FROM account
WHERE is_private_profile = 0

UNION ALL

SELECT
    'crime = 1' AS query_test,
    COUNT(*) AS resultado
FROM account_favorite_locations
WHERE crime = 1

UNION ALL

SELECT
    'traffic_accident = 0' AS query_test,
    COUNT(*) AS resultado
FROM account_favorite_locations
WHERE traffic_accident = 0;

\echo ''

-- =====================================================
-- 8. DISTRIBUCIÓN DE VALORES
-- =====================================================

\echo '9. Distribución de valores en columnas migradas de account:'
\echo ''

SELECT
    'is_private_profile' AS columna,
    is_private_profile AS valor,
    COUNT(*) AS cantidad
FROM account
GROUP BY is_private_profile
ORDER BY valor

UNION ALL

SELECT
    'receive_notifications' AS columna,
    receive_notifications AS valor,
    COUNT(*) AS cantidad
FROM account
GROUP BY receive_notifications
ORDER BY valor

UNION ALL

SELECT
    'has_finished_tutorial' AS columna,
    has_finished_tutorial AS valor,
    COUNT(*) AS cantidad
FROM account
GROUP BY has_finished_tutorial
ORDER BY valor

UNION ALL

SELECT
    'has_watch_new_incident_tutorial' AS columna,
    has_watch_new_incident_tutorial AS valor,
    COUNT(*) AS cantidad
FROM account
GROUP BY has_watch_new_incident_tutorial
ORDER BY valor;

\echo ''
\echo '10. Distribución de valores en columnas migradas de account_favorite_locations:'
\echo ''

SELECT
    'crime' AS columna,
    crime AS valor,
    COUNT(*) AS cantidad
FROM account_favorite_locations
GROUP BY crime
ORDER BY valor

UNION ALL

SELECT
    'traffic_accident' AS columna,
    traffic_accident AS valor,
    COUNT(*) AS cantidad
FROM account_favorite_locations
GROUP BY traffic_accident
ORDER BY valor

UNION ALL

SELECT
    'medical_emergency' AS columna,
    medical_emergency AS valor,
    COUNT(*) AS cantidad
FROM account_favorite_locations
GROUP BY medical_emergency
ORDER BY valor;

\echo ''
\echo '=================================================='
\echo 'VERIFICACIÓN COMPLETADA'
\echo '=================================================='
\echo ''
\echo 'Revisa los resultados:'
\echo '  - Todas las columnas deben ser tipo SMALLINT'
\echo '  - Todos los valores deben ser 0 o 1'
\echo '  - Las queries numéricas deben funcionar'
\echo ''
\echo '=================================================='
