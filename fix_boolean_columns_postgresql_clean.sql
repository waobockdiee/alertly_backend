-- =====================================================
-- MIGRACIÓN: Estandarizar Columnas Booleanas en PostgreSQL
-- Fecha: 2026-01-18
-- Objetivo: Hacer que PostgreSQL se comporte IGUAL que MySQL
-- =====================================================

\timing on

BEGIN;

-- =====================================================
-- FASE 1: CREAR BACKUPS
-- =====================================================

\echo 'Creando backups...'

DROP TABLE IF EXISTS account_backup_20260118;
DROP TABLE IF EXISTS account_favorite_locations_backup_20260118;

CREATE TABLE account_backup_20260118 AS SELECT * FROM account;
CREATE TABLE account_favorite_locations_backup_20260118 AS SELECT * FROM account_favorite_locations;

\echo 'Backups creados exitosamente'

-- =====================================================
-- FASE 2: MIGRAR TABLA account
-- =====================================================

\echo 'Migrando tabla account...'

-- 2.1: Convertir is_private_profile (BOOLEAN -> SMALLINT)
ALTER TABLE account
  ALTER COLUMN is_private_profile TYPE SMALLINT
  USING (CASE WHEN is_private_profile THEN 1 ELSE 0 END);

ALTER TABLE account
  ALTER COLUMN is_private_profile SET DEFAULT 0;

ALTER TABLE account
  ALTER COLUMN is_private_profile SET NOT NULL;

-- 2.2: Convertir receive_notifications (BOOLEAN -> SMALLINT)
ALTER TABLE account
  ALTER COLUMN receive_notifications TYPE SMALLINT
  USING (CASE WHEN receive_notifications THEN 1 ELSE 0 END);

ALTER TABLE account
  ALTER COLUMN receive_notifications SET DEFAULT 1;

ALTER TABLE account
  ALTER COLUMN receive_notifications SET NOT NULL;

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

\echo 'Tabla account migrada'

-- =====================================================
-- FASE 3: MIGRAR TABLA account_favorite_locations
-- =====================================================

\echo 'Migrando tabla account_favorite_locations...'

-- Convertir todas las columnas de categorías (BOOLEAN -> SMALLINT)
ALTER TABLE account_favorite_locations
  ALTER COLUMN crime TYPE SMALLINT USING (CASE WHEN crime THEN 1 ELSE 0 END),
  ALTER COLUMN crime SET DEFAULT 1,
  ALTER COLUMN crime SET NOT NULL;

ALTER TABLE account_favorite_locations
  ALTER COLUMN traffic_accident TYPE SMALLINT USING (CASE WHEN traffic_accident THEN 1 ELSE 0 END),
  ALTER COLUMN traffic_accident SET DEFAULT 1,
  ALTER COLUMN traffic_accident SET NOT NULL;

ALTER TABLE account_favorite_locations
  ALTER COLUMN medical_emergency TYPE SMALLINT USING (CASE WHEN medical_emergency THEN 1 ELSE 0 END),
  ALTER COLUMN medical_emergency SET DEFAULT 1,
  ALTER COLUMN medical_emergency SET NOT NULL;

ALTER TABLE account_favorite_locations
  ALTER COLUMN fire_incident TYPE SMALLINT USING (CASE WHEN fire_incident THEN 1 ELSE 0 END),
  ALTER COLUMN fire_incident SET DEFAULT 1,
  ALTER COLUMN fire_incident SET NOT NULL;

ALTER TABLE account_favorite_locations
  ALTER COLUMN vandalism TYPE SMALLINT USING (CASE WHEN vandalism THEN 1 ELSE 0 END),
  ALTER COLUMN vandalism SET DEFAULT 1,
  ALTER COLUMN vandalism SET NOT NULL;

ALTER TABLE account_favorite_locations
  ALTER COLUMN suspicious_activity TYPE SMALLINT USING (CASE WHEN suspicious_activity THEN 1 ELSE 0 END),
  ALTER COLUMN suspicious_activity SET DEFAULT 1,
  ALTER COLUMN suspicious_activity SET NOT NULL;

ALTER TABLE account_favorite_locations
  ALTER COLUMN infrastructure_issues TYPE SMALLINT USING (CASE WHEN infrastructure_issues THEN 1 ELSE 0 END),
  ALTER COLUMN infrastructure_issues SET DEFAULT 1,
  ALTER COLUMN infrastructure_issues SET NOT NULL;

ALTER TABLE account_favorite_locations
  ALTER COLUMN extreme_weather TYPE SMALLINT USING (CASE WHEN extreme_weather THEN 1 ELSE 0 END),
  ALTER COLUMN extreme_weather SET DEFAULT 1,
  ALTER COLUMN extreme_weather SET NOT NULL;

ALTER TABLE account_favorite_locations
  ALTER COLUMN community_events TYPE SMALLINT USING (CASE WHEN community_events THEN 1 ELSE 0 END),
  ALTER COLUMN community_events SET DEFAULT 1,
  ALTER COLUMN community_events SET NOT NULL;

ALTER TABLE account_favorite_locations
  ALTER COLUMN dangerous_wildlife_sighting TYPE SMALLINT USING (CASE WHEN dangerous_wildlife_sighting THEN 1 ELSE 0 END),
  ALTER COLUMN dangerous_wildlife_sighting SET DEFAULT 1,
  ALTER COLUMN dangerous_wildlife_sighting SET NOT NULL;

ALTER TABLE account_favorite_locations
  ALTER COLUMN positive_actions TYPE SMALLINT USING (CASE WHEN positive_actions THEN 1 ELSE 0 END),
  ALTER COLUMN positive_actions SET DEFAULT 1,
  ALTER COLUMN positive_actions SET NOT NULL;

ALTER TABLE account_favorite_locations
  ALTER COLUMN lost_pet TYPE SMALLINT USING (CASE WHEN lost_pet THEN 1 ELSE 0 END),
  ALTER COLUMN lost_pet SET DEFAULT 1,
  ALTER COLUMN lost_pet SET NOT NULL;

\echo 'Tabla account_favorite_locations migrada'

-- =====================================================
-- COMMIT
-- =====================================================

COMMIT;

\echo ''
\echo '=================================================='
\echo 'MIGRACIÓN COMPLETADA EXITOSAMENTE'
\echo '=================================================='
\echo 'Total de columnas migradas: 16'
\echo '  - account: 4 columnas'
\echo '  - account_favorite_locations: 12 columnas'
\echo ''
\echo 'Ejecuta fix_boolean_columns_verify.sql para verificar'
\echo '=================================================='
