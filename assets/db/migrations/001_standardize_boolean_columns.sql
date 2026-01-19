-- Migration: Estandarizar todas las columnas booleanas a tipo BOOLEAN en PostgreSQL
-- Fecha: 2026-01-18
-- Descripción: Convierte SMALLINT/CHAR(2) a BOOLEAN para consistencia
-- Base de datos: PostgreSQL (Railway)
-- Impacto:
--   - Tabla account: 7 columnas (is_premium, can_update_*, has_finished_tutorial, has_watch_new_incident_tutorial)
--   - Tabla account_favorite_locations: 1 columna (status)
--
-- IMPORTANTE: Ejecutar en horario de bajo tráfico. Duración estimada: <1 minuto.

BEGIN;

-- ============================================================================
-- TABLA: account
-- ============================================================================

-- 1. is_premium: SMALLINT → BOOLEAN
-- Actualmente almacena 0/1, convertir a false/true
ALTER TABLE account
ALTER COLUMN is_premium TYPE BOOLEAN
USING CASE WHEN is_premium = 1 THEN true ELSE false END;

COMMENT ON COLUMN account.is_premium IS 'Premium subscription status (migrated from SMALLINT to BOOLEAN on 2026-01-18)';

-- 2. can_update_email: SMALLINT → BOOLEAN
ALTER TABLE account
ALTER COLUMN can_update_email TYPE BOOLEAN
USING CASE WHEN can_update_email = 1 THEN true ELSE false END;

-- 3. can_update_nickname: SMALLINT → BOOLEAN
ALTER TABLE account
ALTER COLUMN can_update_nickname TYPE BOOLEAN
USING CASE WHEN can_update_nickname = 1 THEN true ELSE false END;

-- 4. can_update_fullname: SMALLINT → BOOLEAN
ALTER TABLE account
ALTER COLUMN can_update_fullname TYPE BOOLEAN
USING CASE WHEN can_update_fullname = 1 THEN true ELSE false END;

-- 5. can_update_birthdate: SMALLINT → BOOLEAN
ALTER TABLE account
ALTER COLUMN can_update_birthdate TYPE BOOLEAN
USING CASE WHEN can_update_birthdate = 1 THEN true ELSE false END;

-- 6. has_finished_tutorial: CHAR(2) → BOOLEAN
-- Actualmente almacena '0 ', '1 ' (con espacios), convertir a false/true
ALTER TABLE account
ALTER COLUMN has_finished_tutorial TYPE BOOLEAN
USING CASE WHEN TRIM(has_finished_tutorial) = '1' THEN true ELSE false END;

COMMENT ON COLUMN account.has_finished_tutorial IS 'Tutorial completion status (migrated from CHAR(2) to BOOLEAN on 2026-01-18)';

-- 7. has_watch_new_incident_tutorial: CHAR(2) → BOOLEAN
ALTER TABLE account
ALTER COLUMN has_watch_new_incident_tutorial TYPE BOOLEAN
USING CASE WHEN TRIM(has_watch_new_incident_tutorial) = '1' THEN true ELSE false END;

COMMENT ON COLUMN account.has_watch_new_incident_tutorial IS 'New incident tutorial watched status (migrated from CHAR(2) to BOOLEAN on 2026-01-18)';

-- ============================================================================
-- TABLA: account_favorite_locations
-- ============================================================================

-- 8. status: SMALLINT → BOOLEAN
-- Actualmente almacena 0/1, convertir a false/true (1 = active, 0 = inactive)
ALTER TABLE account_favorite_locations
ALTER COLUMN status TYPE BOOLEAN
USING CASE WHEN status = 1 THEN true ELSE false END;

COMMENT ON COLUMN account_favorite_locations.status IS 'Location active status (migrated from SMALLINT to BOOLEAN on 2026-01-18)';

-- ============================================================================
-- VERIFICACIÓN DE CAMBIOS
-- ============================================================================

-- Verificar tipos de columnas migradas en account
DO $$
DECLARE
    v_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO v_count
    FROM information_schema.columns
    WHERE table_name = 'account'
      AND column_name IN (
          'is_premium', 'can_update_email', 'can_update_nickname',
          'can_update_fullname', 'can_update_birthdate',
          'has_finished_tutorial', 'has_watch_new_incident_tutorial'
      )
      AND data_type = 'boolean';

    IF v_count <> 7 THEN
        RAISE EXCEPTION 'MIGRATION FAILED: Expected 7 BOOLEAN columns in account, found %', v_count;
    END IF;

    RAISE NOTICE '✅ account: 7 columns successfully migrated to BOOLEAN';
END $$;

-- Verificar tipos de columnas migradas en account_favorite_locations
DO $$
DECLARE
    v_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO v_count
    FROM information_schema.columns
    WHERE table_name = 'account_favorite_locations'
      AND column_name = 'status'
      AND data_type = 'boolean';

    IF v_count <> 1 THEN
        RAISE EXCEPTION 'MIGRATION FAILED: Expected status column as BOOLEAN in account_favorite_locations';
    END IF;

    RAISE NOTICE '✅ account_favorite_locations: status column successfully migrated to BOOLEAN';
END $$;

-- Verificar que no hay valores NULL inesperados
DO $$
DECLARE
    v_null_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO v_null_count
    FROM account
    WHERE is_premium IS NULL
       OR receive_notifications IS NULL
       OR is_private_profile IS NULL
       OR has_finished_tutorial IS NULL
       OR has_watch_new_incident_tutorial IS NULL;

    IF v_null_count > 0 THEN
        RAISE WARNING 'Found % rows with NULL boolean values in account table', v_null_count;
    ELSE
        RAISE NOTICE '✅ No NULL values found in boolean columns';
    END IF;
END $$;

-- ============================================================================
-- ÍNDICES (OPCIONAL - DESCOMENTAR SI SE NECESITAN)
-- ============================================================================

-- Crear índice para queries que filtran por is_premium y status
-- CREATE INDEX CONCURRENTLY idx_account_premium_status ON account(is_premium, status) WHERE is_premium = true AND status = 'active';

-- Crear índice para queries que filtran por receive_notifications
-- CREATE INDEX CONCURRENTLY idx_account_notifications ON account(receive_notifications, status) WHERE receive_notifications = true AND status = 'active';

-- ============================================================================
-- COMMIT Y RESUMEN
-- ============================================================================

COMMIT;

-- Mostrar resumen final
SELECT
    '✅ MIGRATION COMPLETED SUCCESSFULLY' AS status,
    NOW() AS completed_at,
    8 AS total_columns_migrated,
    2 AS tables_affected;

-- Query de validación final (ejecutar manualmente después del COMMIT)
-- SELECT
--     column_name,
--     data_type,
--     is_nullable,
--     column_default
-- FROM information_schema.columns
-- WHERE table_name IN ('account', 'account_favorite_locations')
--   AND column_name IN (
--       'is_premium', 'receive_notifications', 'is_private_profile',
--       'has_finished_tutorial', 'has_watch_new_incident_tutorial',
--       'can_update_email', 'can_update_nickname', 'can_update_fullname', 'can_update_birthdate',
--       'status', 'crime', 'traffic_accident'
--   )
-- ORDER BY table_name, column_name;
