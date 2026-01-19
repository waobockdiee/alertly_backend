-- ROLLBACK Script: Revertir migración de columnas BOOLEAN → tipos originales
-- Fecha: 2026-01-18
-- Descripción: Revierte los cambios de 001_standardize_boolean_columns.sql
-- IMPORTANTE: Solo ejecutar si la migración causa problemas críticos

BEGIN;

-- ============================================================================
-- TABLA: account - REVERTIR A TIPOS ORIGINALES
-- ============================================================================

-- 1. is_premium: BOOLEAN → SMALLINT
ALTER TABLE account
ALTER COLUMN is_premium TYPE SMALLINT
USING CASE WHEN is_premium THEN 1 ELSE 0 END;

-- 2. can_update_email: BOOLEAN → SMALLINT
ALTER TABLE account
ALTER COLUMN can_update_email TYPE SMALLINT
USING CASE WHEN can_update_email THEN 1 ELSE 0 END;

-- 3. can_update_nickname: BOOLEAN → SMALLINT
ALTER TABLE account
ALTER COLUMN can_update_nickname TYPE SMALLINT
USING CASE WHEN can_update_nickname THEN 1 ELSE 0 END;

-- 4. can_update_fullname: BOOLEAN → SMALLINT
ALTER TABLE account
ALTER COLUMN can_update_fullname TYPE SMALLINT
USING CASE WHEN can_update_fullname THEN 1 ELSE 0 END;

-- 5. can_update_birthdate: BOOLEAN → SMALLINT
ALTER TABLE account
ALTER COLUMN can_update_birthdate TYPE SMALLINT
USING CASE WHEN can_update_birthdate THEN 1 ELSE 0 END;

-- 6. has_finished_tutorial: BOOLEAN → CHAR(2)
ALTER TABLE account
ALTER COLUMN has_finished_tutorial TYPE CHAR(2)
USING CASE WHEN has_finished_tutorial THEN '1' ELSE '0' END;

-- 7. has_watch_new_incident_tutorial: BOOLEAN → CHAR(2)
ALTER TABLE account
ALTER COLUMN has_watch_new_incident_tutorial TYPE CHAR(2)
USING CASE WHEN has_watch_new_incident_tutorial THEN '1' ELSE '0' END;

-- ============================================================================
-- TABLA: account_favorite_locations - REVERTIR A TIPOS ORIGINALES
-- ============================================================================

-- 8. status: BOOLEAN → SMALLINT
ALTER TABLE account_favorite_locations
ALTER COLUMN status TYPE SMALLINT
USING CASE WHEN status THEN 1 ELSE 0 END;

-- ============================================================================
-- VERIFICACIÓN DE ROLLBACK
-- ============================================================================

DO $$
BEGIN
    RAISE NOTICE '⚠️  ROLLBACK COMPLETED - Database reverted to original types';
    RAISE NOTICE '⚠️  Recuerde actualizar el código Go para usar dbtypes.BoolToInt()';
END $$;

COMMIT;

-- Mostrar resumen final
SELECT
    '⚠️  ROLLBACK COMPLETED' AS status,
    NOW() AS completed_at,
    'All boolean columns reverted to original types (SMALLINT/CHAR)' AS message;
