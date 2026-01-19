-- =====================================================
-- SCRIPT DE ROLLBACK: Migración de Columnas Booleanas
-- Fecha: 2026-01-18
-- =====================================================
--
-- IMPORTANTE: Este script restaura las tablas desde los backups
-- Solo ejecutar si la migración principal falló o necesita revertirse
--
-- =====================================================

BEGIN;

RAISE NOTICE 'Iniciando rollback de migración de columnas booleanas...';

-- =====================================================
-- VERIFICAR QUE EXISTEN LOS BACKUPS
-- =====================================================

DO $$
DECLARE
    account_backup_exists BOOLEAN;
    locations_backup_exists BOOLEAN;
BEGIN
    SELECT EXISTS (
        SELECT FROM information_schema.tables
        WHERE table_name = 'account_backup_20260118'
    ) INTO account_backup_exists;

    SELECT EXISTS (
        SELECT FROM information_schema.tables
        WHERE table_name = 'account_favorite_locations_backup_20260118'
    ) INTO locations_backup_exists;

    IF NOT account_backup_exists THEN
        RAISE EXCEPTION 'No existe backup de account - ABORTANDO';
    END IF;

    IF NOT locations_backup_exists THEN
        RAISE EXCEPTION 'No existe backup de account_favorite_locations - ABORTANDO';
    END IF;

    RAISE NOTICE '✓ Backups verificados';
END $$;

-- =====================================================
-- FASE 1: RESTAURAR account
-- =====================================================

RAISE NOTICE 'Restaurando tabla account...';

-- Eliminar constraints que referencian account
ALTER TABLE account_achievements DROP CONSTRAINT IF EXISTS fk_account_achievements_account1;
ALTER TABLE account_cluster_saved DROP CONSTRAINT IF EXISTS fk_account_cluster_saved_account1;
ALTER TABLE account_history DROP CONSTRAINT IF EXISTS fk_account_history_account1;
ALTER TABLE account_favorite_locations DROP CONSTRAINT IF EXISTS fk_account_locations_account1;
ALTER TABLE notifications DROP CONSTRAINT IF EXISTS fk_account_notifications_account1;
ALTER TABLE account_premium_payment_history DROP CONSTRAINT IF EXISTS fk_account_premium_payment_history_account1;
ALTER TABLE account_reports DROP CONSTRAINT IF EXISTS fk_account_reports_account1;
ALTER TABLE account_reports DROP CONSTRAINT IF EXISTS fk_account_reports_account2;
ALTER TABLE account_session_history DROP CONSTRAINT IF EXISTS fk_account_session_history_account;
ALTER TABLE device_tokens DROP CONSTRAINT IF EXISTS fk_device_tokens_account1;
ALTER TABLE feedback DROP CONSTRAINT IF EXISTS fk_feedback_account1;
ALTER TABLE incident_clusters DROP CONSTRAINT IF EXISTS fk_incident_clusters_account1;
ALTER TABLE incident_comments DROP CONSTRAINT IF EXISTS fk_incident_comments_account1;
ALTER TABLE incident_flags DROP CONSTRAINT IF EXISTS fk_incident_flags_account1;
ALTER TABLE incident_logs DROP CONSTRAINT IF EXISTS fk_incident_logs_account1;
ALTER TABLE incident_reports DROP CONSTRAINT IF EXISTS fk_incident_reports_account1;
ALTER TABLE incident_votes DROP CONSTRAINT IF EXISTS fk_incident_votes_account1;
ALTER TABLE notification_deliveries DROP CONSTRAINT IF EXISTS fk_notifications_deliveries_account1;
ALTER TABLE referral_conversions DROP CONSTRAINT IF EXISTS referral_conversions_ibfk_1;
ALTER TABLE referral_premium_conversions DROP CONSTRAINT IF EXISTS referral_premium_conversions_ibfk_1;

RAISE NOTICE 'Constraints eliminados';

-- Eliminar tabla account actual
DROP TABLE account;

RAISE NOTICE 'Tabla account eliminada';

-- Restaurar desde backup
CREATE TABLE account AS SELECT * FROM account_backup_20260118;

RAISE NOTICE 'Tabla account restaurada desde backup';

-- Recrear primary key
ALTER TABLE account ADD PRIMARY KEY (account_id);

RAISE NOTICE 'Primary key recreado';

-- Recrear índices
CREATE UNIQUE INDEX email_unique ON account(email);
CREATE INDEX idx_account_credibility ON account(account_id, credibility, status);
CREATE INDEX idx_account_email ON account(email, status);
CREATE INDEX idx_account_nickname ON account(nickname, status);
CREATE INDEX idx_account_premium ON account(is_premium, premium_expired_date, status);
CREATE INDEX idx_account_status ON account(account_id, status);

RAISE NOTICE 'Índices recreados';

-- =====================================================
-- FASE 2: RESTAURAR account_favorite_locations
-- =====================================================

RAISE NOTICE 'Restaurando tabla account_favorite_locations...';

-- Eliminar tabla actual
DROP TABLE account_favorite_locations;

-- Restaurar desde backup
CREATE TABLE account_favorite_locations AS SELECT * FROM account_favorite_locations_backup_20260118;

RAISE NOTICE 'Tabla account_favorite_locations restaurada desde backup';

-- Recrear primary key
ALTER TABLE account_favorite_locations ADD PRIMARY KEY (afl_id);

-- Recrear índice
CREATE INDEX fk_account_locations_account1_idx ON account_favorite_locations(account_id);

RAISE NOTICE 'Índices recreados';

-- =====================================================
-- FASE 3: RECREAR FOREIGN KEYS
-- =====================================================

RAISE NOTICE 'Recreando foreign keys...';

ALTER TABLE account_achievements
    ADD CONSTRAINT fk_account_achievements_account1
    FOREIGN KEY (account_id) REFERENCES account(account_id);

ALTER TABLE account_cluster_saved
    ADD CONSTRAINT fk_account_cluster_saved_account1
    FOREIGN KEY (account_id) REFERENCES account(account_id);

ALTER TABLE account_history
    ADD CONSTRAINT fk_account_history_account1
    FOREIGN KEY (account_id) REFERENCES account(account_id);

ALTER TABLE account_favorite_locations
    ADD CONSTRAINT fk_account_locations_account1
    FOREIGN KEY (account_id) REFERENCES account(account_id);

ALTER TABLE notifications
    ADD CONSTRAINT fk_account_notifications_account1
    FOREIGN KEY (owner_account_id) REFERENCES account(account_id);

ALTER TABLE account_premium_payment_history
    ADD CONSTRAINT fk_account_premium_payment_history_account1
    FOREIGN KEY (account_id) REFERENCES account(account_id);

ALTER TABLE account_reports
    ADD CONSTRAINT fk_account_reports_account1
    FOREIGN KEY (account_id_whos_reporting) REFERENCES account(account_id);

ALTER TABLE account_reports
    ADD CONSTRAINT fk_account_reports_account2
    FOREIGN KEY (account_id) REFERENCES account(account_id);

ALTER TABLE account_session_history
    ADD CONSTRAINT fk_account_session_history_account
    FOREIGN KEY (account_id) REFERENCES account(account_id);

ALTER TABLE device_tokens
    ADD CONSTRAINT fk_device_tokens_account1
    FOREIGN KEY (account_id) REFERENCES account(account_id);

ALTER TABLE feedback
    ADD CONSTRAINT fk_feedback_account1
    FOREIGN KEY (account_id) REFERENCES account(account_id);

ALTER TABLE incident_clusters
    ADD CONSTRAINT fk_incident_clusters_account1
    FOREIGN KEY (account_id) REFERENCES account(account_id);

ALTER TABLE incident_comments
    ADD CONSTRAINT fk_incident_comments_account1
    FOREIGN KEY (account_id) REFERENCES account(account_id);

ALTER TABLE incident_flags
    ADD CONSTRAINT fk_incident_flags_account1
    FOREIGN KEY (account_id) REFERENCES account(account_id);

ALTER TABLE incident_logs
    ADD CONSTRAINT fk_incident_logs_account1
    FOREIGN KEY (account_id) REFERENCES account(account_id);

ALTER TABLE incident_reports
    ADD CONSTRAINT fk_incident_reports_account1
    FOREIGN KEY (account_id) REFERENCES account(account_id);

ALTER TABLE incident_votes
    ADD CONSTRAINT fk_incident_votes_account1
    FOREIGN KEY (account_id) REFERENCES account(account_id);

ALTER TABLE notification_deliveries
    ADD CONSTRAINT fk_notifications_deliveries_account1
    FOREIGN KEY (to_account_id) REFERENCES account(account_id);

ALTER TABLE referral_conversions
    ADD CONSTRAINT referral_conversions_ibfk_1
    FOREIGN KEY (user_id) REFERENCES account(account_id);

ALTER TABLE referral_premium_conversions
    ADD CONSTRAINT referral_premium_conversions_ibfk_1
    FOREIGN KEY (user_id) REFERENCES account(account_id);

RAISE NOTICE 'Foreign keys recreados';

-- =====================================================
-- FASE 4: RECREAR CHECK CONSTRAINTS
-- =====================================================

RAISE NOTICE 'Recreando check constraints...';

ALTER TABLE account
    ADD CONSTRAINT chk_account_role
    CHECK (role::text = ANY (ARRAY['citizen'::character varying, 'admin'::character varying]::text[]));

ALTER TABLE account
    ADD CONSTRAINT chk_account_status
    CHECK (status::text = ANY (ARRAY['pending_activation'::character varying, 'active'::character varying, 'inactive'::character varying, 'blocked'::character varying]::text[]));

ALTER TABLE account_favorite_locations
    ADD CONSTRAINT chk_favorite_locations_latitude
    CHECK (latitude >= '-90'::integer::numeric AND latitude <= 90::numeric);

ALTER TABLE account_favorite_locations
    ADD CONSTRAINT chk_favorite_locations_longitude
    CHECK (longitude >= '-180'::integer::numeric AND longitude <= 180::numeric);

RAISE NOTICE 'Check constraints recreados';

-- =====================================================
-- FASE 5: VERIFICACIÓN
-- =====================================================

DO $$
DECLARE
    account_count INT;
    locations_count INT;
    backup_account_count INT;
    backup_locations_count INT;
BEGIN
    SELECT COUNT(*) INTO account_count FROM account;
    SELECT COUNT(*) INTO locations_count FROM account_favorite_locations;
    SELECT COUNT(*) INTO backup_account_count FROM account_backup_20260118;
    SELECT COUNT(*) INTO backup_locations_count FROM account_favorite_locations_backup_20260118;

    RAISE NOTICE 'Verificación de registros:';
    RAISE NOTICE '  account: % registros (backup: %)', account_count, backup_account_count;
    RAISE NOTICE '  account_favorite_locations: % registros (backup: %)', locations_count, backup_locations_count;

    IF account_count != backup_account_count THEN
        RAISE WARNING 'Cuenta de registros en account no coincide con backup';
    END IF;

    IF locations_count != backup_locations_count THEN
        RAISE WARNING 'Cuenta de registros en account_favorite_locations no coincide con backup';
    END IF;
END $$;

COMMIT;

RAISE NOTICE '==================================================';
RAISE NOTICE 'ROLLBACK COMPLETADO EXITOSAMENTE';
RAISE NOTICE '==================================================';
RAISE NOTICE 'Las tablas han sido restauradas a su estado anterior';
RAISE NOTICE 'Los backups se mantienen disponibles para referencia';
RAISE NOTICE '==================================================';
