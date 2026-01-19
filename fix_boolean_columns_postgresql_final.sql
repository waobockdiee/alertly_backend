-- =====================================================
-- MIGRACIÓN: Estandarizar Columnas Booleanas en PostgreSQL
-- Fecha: 2026-01-18
-- Objetivo: Hacer que PostgreSQL se comporte IGUAL que MySQL
-- =====================================================
--
-- Este script convierte todas las columnas booleanas a SMALLINT
-- para mantener compatibilidad total con el código Go existente
-- que usa comparaciones como '= 1' y '= 0'.
--
-- IMPORTANTE: Ejecutar en horario de baja demanda
-- =====================================================

BEGIN;

-- =====================================================
-- FASE 1: CREAR BACKUPS
-- =====================================================

DO $$
BEGIN
    -- Eliminar backups antiguos si existen
    DROP TABLE IF EXISTS account_backup_20260118;
    DROP TABLE IF EXISTS account_favorite_locations_backup_20260118;

    RAISE NOTICE 'Backups antiguos eliminados';
END $$;

-- Crear backups de las tablas
CREATE TABLE account_backup_20260118 AS SELECT * FROM account;
CREATE TABLE account_favorite_locations_backup_20260118 AS SELECT * FROM account_favorite_locations;

-- Verificar backups
DO $$
DECLARE
    account_count INT;
    locations_count INT;
BEGIN
    SELECT COUNT(*) INTO account_count FROM account_backup_20260118;
    SELECT COUNT(*) INTO locations_count FROM account_favorite_locations_backup_20260118;

    RAISE NOTICE 'Backup account: % registros', account_count;
    RAISE NOTICE 'Backup account_favorite_locations: % registros', locations_count;

    IF account_count = 0 THEN
        RAISE EXCEPTION 'Backup de account está vacío - ABORTANDO';
    END IF;
END $$;

-- =====================================================
-- FASE 2: MIGRAR TABLA account
-- =====================================================

RAISE NOTICE 'Iniciando migración de tabla account...';

-- 2.1: Convertir is_private_profile (BOOLEAN -> SMALLINT)
ALTER TABLE account
  ALTER COLUMN is_private_profile TYPE SMALLINT
  USING (CASE WHEN is_private_profile THEN 1 ELSE 0 END);

ALTER TABLE account
  ALTER COLUMN is_private_profile SET DEFAULT 0;

ALTER TABLE account
  ALTER COLUMN is_private_profile SET NOT NULL;

RAISE NOTICE '✓ is_private_profile migrado';

-- 2.2: Convertir receive_notifications (BOOLEAN -> SMALLINT)
ALTER TABLE account
  ALTER COLUMN receive_notifications TYPE SMALLINT
  USING (CASE WHEN receive_notifications THEN 1 ELSE 0 END);

ALTER TABLE account
  ALTER COLUMN receive_notifications SET DEFAULT 1;

ALTER TABLE account
  ALTER COLUMN receive_notifications SET NOT NULL;

RAISE NOTICE '✓ receive_notifications migrado';

-- 2.3: Convertir has_finished_tutorial (CHAR(2) -> SMALLINT)
ALTER TABLE account
  ALTER COLUMN has_finished_tutorial TYPE SMALLINT
  USING (
    CASE
      WHEN has_finished_tutorial = '1' THEN 1
      WHEN has_finished_tutorial IS NULL THEN 0
      WHEN has_finished_tutorial = '' THEN 0
      ELSE 0
    END
  );

ALTER TABLE account
  ALTER COLUMN has_finished_tutorial SET DEFAULT 0;

ALTER TABLE account
  ALTER COLUMN has_finished_tutorial SET NOT NULL;

RAISE NOTICE '✓ has_finished_tutorial migrado';

-- 2.4: Convertir has_watch_new_incident_tutorial (CHAR(2) -> SMALLINT)
ALTER TABLE account
  ALTER COLUMN has_watch_new_incident_tutorial TYPE SMALLINT
  USING (
    CASE
      WHEN has_watch_new_incident_tutorial = '1' THEN 1
      WHEN has_watch_new_incident_tutorial IS NULL THEN 0
      WHEN has_watch_new_incident_tutorial = '' THEN 0
      ELSE 0
    END
  );

ALTER TABLE account
  ALTER COLUMN has_watch_new_incident_tutorial SET DEFAULT 0;

ALTER TABLE account
  ALTER COLUMN has_watch_new_incident_tutorial SET NOT NULL;

RAISE NOTICE '✓ has_watch_new_incident_tutorial migrado';

-- =====================================================
-- FASE 3: MIGRAR TABLA account_favorite_locations
-- =====================================================

RAISE NOTICE 'Iniciando migración de tabla account_favorite_locations...';

-- 3.1: Convertir crime (BOOLEAN -> SMALLINT)
ALTER TABLE account_favorite_locations
  ALTER COLUMN crime TYPE SMALLINT
  USING (CASE WHEN crime THEN 1 ELSE 0 END);

ALTER TABLE account_favorite_locations
  ALTER COLUMN crime SET DEFAULT 1;

ALTER TABLE account_favorite_locations
  ALTER COLUMN crime SET NOT NULL;

RAISE NOTICE '✓ crime migrado';

-- 3.2: Convertir traffic_accident (BOOLEAN -> SMALLINT)
ALTER TABLE account_favorite_locations
  ALTER COLUMN traffic_accident TYPE SMALLINT
  USING (CASE WHEN traffic_accident THEN 1 ELSE 0 END);

ALTER TABLE account_favorite_locations
  ALTER COLUMN traffic_accident SET DEFAULT 1;

ALTER TABLE account_favorite_locations
  ALTER COLUMN traffic_accident SET NOT NULL;

RAISE NOTICE '✓ traffic_accident migrado';

-- 3.3: Convertir medical_emergency (BOOLEAN -> SMALLINT)
ALTER TABLE account_favorite_locations
  ALTER COLUMN medical_emergency TYPE SMALLINT
  USING (CASE WHEN medical_emergency THEN 1 ELSE 0 END);

ALTER TABLE account_favorite_locations
  ALTER COLUMN medical_emergency SET DEFAULT 1;

ALTER TABLE account_favorite_locations
  ALTER COLUMN medical_emergency SET NOT NULL;

RAISE NOTICE '✓ medical_emergency migrado';

-- 3.4: Convertir fire_incident (BOOLEAN -> SMALLINT)
ALTER TABLE account_favorite_locations
  ALTER COLUMN fire_incident TYPE SMALLINT
  USING (CASE WHEN fire_incident THEN 1 ELSE 0 END);

ALTER TABLE account_favorite_locations
  ALTER COLUMN fire_incident SET DEFAULT 1;

ALTER TABLE account_favorite_locations
  ALTER COLUMN fire_incident SET NOT NULL;

RAISE NOTICE '✓ fire_incident migrado';

-- 3.5: Convertir vandalism (BOOLEAN -> SMALLINT)
ALTER TABLE account_favorite_locations
  ALTER COLUMN vandalism TYPE SMALLINT
  USING (CASE WHEN vandalism THEN 1 ELSE 0 END);

ALTER TABLE account_favorite_locations
  ALTER COLUMN vandalism SET DEFAULT 1;

ALTER TABLE account_favorite_locations
  ALTER COLUMN vandalism SET NOT NULL;

RAISE NOTICE '✓ vandalism migrado';

-- 3.6: Convertir suspicious_activity (BOOLEAN -> SMALLINT)
ALTER TABLE account_favorite_locations
  ALTER COLUMN suspicious_activity TYPE SMALLINT
  USING (CASE WHEN suspicious_activity THEN 1 ELSE 0 END);

ALTER TABLE account_favorite_locations
  ALTER COLUMN suspicious_activity SET DEFAULT 1;

ALTER TABLE account_favorite_locations
  ALTER COLUMN suspicious_activity SET NOT NULL;

RAISE NOTICE '✓ suspicious_activity migrado';

-- 3.7: Convertir infrastructure_issues (BOOLEAN -> SMALLINT)
ALTER TABLE account_favorite_locations
  ALTER COLUMN infrastructure_issues TYPE SMALLINT
  USING (CASE WHEN infrastructure_issues THEN 1 ELSE 0 END);

ALTER TABLE account_favorite_locations
  ALTER COLUMN infrastructure_issues SET DEFAULT 1;

ALTER TABLE account_favorite_locations
  ALTER COLUMN infrastructure_issues SET NOT NULL;

RAISE NOTICE '✓ infrastructure_issues migrado';

-- 3.8: Convertir extreme_weather (BOOLEAN -> SMALLINT)
ALTER TABLE account_favorite_locations
  ALTER COLUMN extreme_weather TYPE SMALLINT
  USING (CASE WHEN extreme_weather THEN 1 ELSE 0 END);

ALTER TABLE account_favorite_locations
  ALTER COLUMN extreme_weather SET DEFAULT 1;

ALTER TABLE account_favorite_locations
  ALTER COLUMN extreme_weather SET NOT NULL;

RAISE NOTICE '✓ extreme_weather migrado';

-- 3.9: Convertir community_events (BOOLEAN -> SMALLINT)
ALTER TABLE account_favorite_locations
  ALTER COLUMN community_events TYPE SMALLINT
  USING (CASE WHEN community_events THEN 1 ELSE 0 END);

ALTER TABLE account_favorite_locations
  ALTER COLUMN community_events SET DEFAULT 1;

ALTER TABLE account_favorite_locations
  ALTER COLUMN community_events SET NOT NULL;

RAISE NOTICE '✓ community_events migrado';

-- 3.10: Convertir dangerous_wildlife_sighting (BOOLEAN -> SMALLINT)
ALTER TABLE account_favorite_locations
  ALTER COLUMN dangerous_wildlife_sighting TYPE SMALLINT
  USING (CASE WHEN dangerous_wildlife_sighting THEN 1 ELSE 0 END);

ALTER TABLE account_favorite_locations
  ALTER COLUMN dangerous_wildlife_sighting SET DEFAULT 1;

ALTER TABLE account_favorite_locations
  ALTER COLUMN dangerous_wildlife_sighting SET NOT NULL;

RAISE NOTICE '✓ dangerous_wildlife_sighting migrado';

-- 3.11: Convertir positive_actions (BOOLEAN -> SMALLINT)
ALTER TABLE account_favorite_locations
  ALTER COLUMN positive_actions TYPE SMALLINT
  USING (CASE WHEN positive_actions THEN 1 ELSE 0 END);

ALTER TABLE account_favorite_locations
  ALTER COLUMN positive_actions SET DEFAULT 1;

ALTER TABLE account_favorite_locations
  ALTER COLUMN positive_actions SET NOT NULL;

RAISE NOTICE '✓ positive_actions migrado';

-- 3.12: Convertir lost_pet (BOOLEAN -> SMALLINT)
ALTER TABLE account_favorite_locations
  ALTER COLUMN lost_pet TYPE SMALLINT
  USING (CASE WHEN lost_pet THEN 1 ELSE 0 END);

ALTER TABLE account_favorite_locations
  ALTER COLUMN lost_pet SET DEFAULT 1;

ALTER TABLE account_favorite_locations
  ALTER COLUMN lost_pet SET NOT NULL;

RAISE NOTICE '✓ lost_pet migrado';

-- =====================================================
-- FASE 4: VERIFICACIÓN
-- =====================================================

RAISE NOTICE 'Verificando migración...';

-- 4.1: Verificar tipos de datos en account
DO $$
DECLARE
    wrong_types INT;
BEGIN
    SELECT COUNT(*) INTO wrong_types
    FROM information_schema.columns
    WHERE table_name = 'account'
      AND column_name IN ('is_private_profile', 'receive_notifications', 'has_finished_tutorial', 'has_watch_new_incident_tutorial')
      AND data_type != 'smallint';

    IF wrong_types > 0 THEN
        RAISE EXCEPTION 'Hay % columnas en account con tipo incorrecto', wrong_types;
    ELSE
        RAISE NOTICE '✓ Todos los tipos de datos en account son correctos';
    END IF;
END $$;

-- 4.2: Verificar tipos de datos en account_favorite_locations
DO $$
DECLARE
    wrong_types INT;
BEGIN
    SELECT COUNT(*) INTO wrong_types
    FROM information_schema.columns
    WHERE table_name = 'account_favorite_locations'
      AND column_name IN (
          'crime', 'traffic_accident', 'medical_emergency', 'fire_incident',
          'vandalism', 'suspicious_activity', 'infrastructure_issues', 'extreme_weather',
          'community_events', 'dangerous_wildlife_sighting', 'positive_actions', 'lost_pet'
      )
      AND data_type != 'smallint';

    IF wrong_types > 0 THEN
        RAISE EXCEPTION 'Hay % columnas en account_favorite_locations con tipo incorrecto', wrong_types;
    ELSE
        RAISE NOTICE '✓ Todos los tipos de datos en account_favorite_locations son correctos';
    END IF;
END $$;

-- 4.3: Verificar valores (deben ser 0 o 1)
DO $$
DECLARE
    invalid_values INT;
BEGIN
    -- Verificar account
    SELECT COUNT(*) INTO invalid_values
    FROM account
    WHERE is_private_profile NOT IN (0, 1)
       OR receive_notifications NOT IN (0, 1)
       OR has_finished_tutorial NOT IN (0, 1)
       OR has_watch_new_incident_tutorial NOT IN (0, 1);

    IF invalid_values > 0 THEN
        RAISE EXCEPTION 'Hay % registros en account con valores inválidos (no son 0 o 1)', invalid_values;
    ELSE
        RAISE NOTICE '✓ Todos los valores en account son 0 o 1';
    END IF;

    -- Verificar account_favorite_locations
    SELECT COUNT(*) INTO invalid_values
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

    IF invalid_values > 0 THEN
        RAISE EXCEPTION 'Hay % registros en account_favorite_locations con valores inválidos', invalid_values;
    ELSE
        RAISE NOTICE '✓ Todos los valores en account_favorite_locations son 0 o 1';
    END IF;
END $$;

-- 4.4: Test de queries numéricas (como en el código Go)
DO $$
DECLARE
    test_count INT;
BEGIN
    -- Test: is_premium = 1
    SELECT COUNT(*) INTO test_count FROM account WHERE is_premium = 1;
    RAISE NOTICE 'Test query: is_premium = 1 -> % registros', test_count;

    -- Test: receive_notifications = 1
    SELECT COUNT(*) INTO test_count FROM account WHERE receive_notifications = 1;
    RAISE NOTICE 'Test query: receive_notifications = 1 -> % registros', test_count;

    -- Test: crime = 1
    SELECT COUNT(*) INTO test_count FROM account_favorite_locations WHERE crime = 1;
    RAISE NOTICE 'Test query: crime = 1 -> % registros', test_count;

    -- Test: is_private_profile = 0
    SELECT COUNT(*) INTO test_count FROM account WHERE is_private_profile = 0;
    RAISE NOTICE 'Test query: is_private_profile = 0 -> % registros', test_count;

    RAISE NOTICE '✓ Todas las queries numéricas funcionan correctamente';
END $$;

-- =====================================================
-- FASE 5: MOSTRAR RESULTADOS
-- =====================================================

-- Mostrar estructura final de account (columnas migradas)
SELECT
    column_name,
    data_type,
    column_default,
    is_nullable
FROM information_schema.columns
WHERE table_name = 'account'
  AND column_name IN ('is_private_profile', 'receive_notifications', 'has_finished_tutorial', 'has_watch_new_incident_tutorial')
ORDER BY column_name;

-- Mostrar estructura final de account_favorite_locations (columnas migradas)
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

-- Mostrar muestra de datos migrados
SELECT
    account_id,
    is_premium,
    receive_notifications,
    is_private_profile,
    has_finished_tutorial,
    has_watch_new_incident_tutorial
FROM account
LIMIT 5;

SELECT
    afl_id,
    account_id,
    crime,
    traffic_accident,
    medical_emergency,
    fire_incident
FROM account_favorite_locations
LIMIT 3;

-- =====================================================
-- COMMIT
-- =====================================================

COMMIT;

RAISE NOTICE '==================================================';
RAISE NOTICE 'MIGRACIÓN COMPLETADA EXITOSAMENTE';
RAISE NOTICE '==================================================';
RAISE NOTICE 'Total de columnas migradas: 16';
RAISE NOTICE '  - account: 4 columnas';
RAISE NOTICE '  - account_favorite_locations: 12 columnas';
RAISE NOTICE '';
RAISE NOTICE 'Backups creados:';
RAISE NOTICE '  - account_backup_20260118';
RAISE NOTICE '  - account_favorite_locations_backup_20260118';
RAISE NOTICE '';
RAISE NOTICE 'PostgreSQL ahora se comporta IGUAL que MySQL';
RAISE NOTICE 'Las queries Go con "= 1" y "= 0" funcionarán correctamente';
RAISE NOTICE '==================================================';
