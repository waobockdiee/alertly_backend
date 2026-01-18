-- ============================================================================
-- POSTGRESQL SCHEMA FIX - Restaurar Constraints y Defaults desde MySQL
-- ============================================================================
-- Fecha: 2026-01-17
-- Propósito: Corregir la migración de MySQL a PostgreSQL Railway
-- Problema: La migración automática no preservó correctamente:
--   1. Columnas NOT NULL
--   2. DEFAULT values
--   3. Tipos de datos (ENUM → VARCHAR sin constraints)
--
-- IMPORTANTE: Ejecutar en Railway PostgreSQL ANTES de usar la BD en producción
-- ============================================================================

-- ============================================================================
-- 1. TABLA: account
-- ============================================================================

-- Agregar constraints NOT NULL que faltan
ALTER TABLE account ALTER COLUMN email SET NOT NULL;
ALTER TABLE account ALTER COLUMN password SET NOT NULL;
ALTER TABLE account ALTER COLUMN nickname SET NOT NULL;
ALTER TABLE account ALTER COLUMN role SET NOT NULL;
ALTER TABLE account ALTER COLUMN status SET NOT NULL;
ALTER TABLE account ALTER COLUMN credibility SET NOT NULL;
ALTER TABLE account ALTER COLUMN is_private_profile SET NOT NULL;
ALTER TABLE account ALTER COLUMN score SET NOT NULL;
ALTER TABLE account ALTER COLUMN is_premium SET NOT NULL;
ALTER TABLE account ALTER COLUMN counter_total_incidents_created SET NOT NULL;
ALTER TABLE account ALTER COLUMN counter_total_votes_made SET NOT NULL;
ALTER TABLE account ALTER COLUMN counter_total_comments_made SET NOT NULL;
ALTER TABLE account ALTER COLUMN counter_total_locations SET NOT NULL;
ALTER TABLE account ALTER COLUMN counter_total_flags SET NOT NULL;
ALTER TABLE account ALTER COLUMN counter_total_medals SET NOT NULL;
ALTER TABLE account ALTER COLUMN has_finished_tutorial SET NOT NULL;
ALTER TABLE account ALTER COLUMN has_watch_new_incident_tutorial SET NOT NULL;
ALTER TABLE account ALTER COLUMN thumbnail_url SET NOT NULL;
ALTER TABLE account ALTER COLUMN counter_new_notifications SET NOT NULL;
ALTER TABLE account ALTER COLUMN crime SET NOT NULL;
ALTER TABLE account ALTER COLUMN traffic_accident SET NOT NULL;
ALTER TABLE account ALTER COLUMN medical_emergency SET NOT NULL;
ALTER TABLE account ALTER COLUMN fire_incident SET NOT NULL;
ALTER TABLE account ALTER COLUMN vandalism SET NOT NULL;
ALTER TABLE account ALTER COLUMN suspicious_activity SET NOT NULL;
ALTER TABLE account ALTER COLUMN infrastructure_issues SET NOT NULL;
ALTER TABLE account ALTER COLUMN extreme_weather SET NOT NULL;
ALTER TABLE account ALTER COLUMN community_events SET NOT NULL;
ALTER TABLE account ALTER COLUMN dangerous_wildlife_sighting SET NOT NULL;
ALTER TABLE account ALTER COLUMN positive_actions SET NOT NULL;
ALTER TABLE account ALTER COLUMN lost_pet SET NOT NULL;
ALTER TABLE account ALTER COLUMN incident_as_update SET NOT NULL;
ALTER TABLE account ALTER COLUMN can_update_email SET NOT NULL;
ALTER TABLE account ALTER COLUMN can_update_nickname SET NOT NULL;
ALTER TABLE account ALTER COLUMN can_update_fullname SET NOT NULL;
ALTER TABLE account ALTER COLUMN can_update_birthdate SET NOT NULL;
ALTER TABLE account ALTER COLUMN receive_notifications SET NOT NULL;

-- Corregir defaults que están como NULL::type en lugar de verdaderos defaults
ALTER TABLE account ALTER COLUMN email DROP DEFAULT;
ALTER TABLE account ALTER COLUMN email SET DEFAULT NULL;
ALTER TABLE account ALTER COLUMN password DROP DEFAULT;
ALTER TABLE account ALTER COLUMN password SET DEFAULT NULL;
ALTER TABLE account ALTER COLUMN nickname DROP DEFAULT;
ALTER TABLE account ALTER COLUMN nickname SET DEFAULT NULL;

-- ============================================================================
-- 2. TABLA: incident_clusters
-- ============================================================================

-- Agregar NOT NULL constraints
ALTER TABLE incident_clusters ALTER COLUMN insu_id SET NOT NULL;
ALTER TABLE incident_clusters ALTER COLUMN incident_count SET NOT NULL;
ALTER TABLE incident_clusters ALTER COLUMN is_active SET NOT NULL;
ALTER TABLE incident_clusters ALTER COLUMN counter_total_comments SET NOT NULL;
ALTER TABLE incident_clusters ALTER COLUMN counter_total_votes SET NOT NULL;
ALTER TABLE incident_clusters ALTER COLUMN counter_total_views SET NOT NULL;
ALTER TABLE incident_clusters ALTER COLUMN counter_total_flags SET NOT NULL;
ALTER TABLE incident_clusters ALTER COLUMN counter_total_votes_true SET NOT NULL;
ALTER TABLE incident_clusters ALTER COLUMN counter_total_votes_false SET NOT NULL;
ALTER TABLE incident_clusters ALTER COLUMN credibility SET NOT NULL;
ALTER TABLE incident_clusters ALTER COLUMN score_true SET NOT NULL;
ALTER TABLE incident_clusters ALTER COLUMN score_false SET NOT NULL;
ALTER TABLE incident_clusters ALTER COLUMN account_id SET NOT NULL;

-- ============================================================================
-- 3. TABLA: incident_reports
-- ============================================================================

-- Agregar NOT NULL constraints
ALTER TABLE incident_reports ALTER COLUMN account_id SET NOT NULL;
ALTER TABLE incident_reports ALTER COLUMN incl_id SET NOT NULL;
ALTER TABLE incident_reports ALTER COLUMN insu_id SET NOT NULL;
ALTER TABLE incident_reports ALTER COLUMN counter_total_comments SET NOT NULL;
ALTER TABLE incident_reports ALTER COLUMN counter_total_votes SET NOT NULL;
ALTER TABLE incident_reports ALTER COLUMN counter_total_views SET NOT NULL;
ALTER TABLE incident_reports ALTER COLUMN counter_total_flags SET NOT NULL;
ALTER TABLE incident_reports ALTER COLUMN is_anonymous SET NOT NULL;
ALTER TABLE incident_reports ALTER COLUMN counter_total_votes_true SET NOT NULL;
ALTER TABLE incident_reports ALTER COLUMN counter_total_votes_fake SET NOT NULL;
ALTER TABLE incident_reports ALTER COLUMN is_active SET NOT NULL;
ALTER TABLE incident_reports ALTER COLUMN vote SET NOT NULL;
ALTER TABLE incident_reports ALTER COLUMN credibility SET NOT NULL;

-- ============================================================================
-- 4. TABLA: incident_subcategories
-- ============================================================================

-- Agregar NOT NULL constraints
ALTER TABLE incident_subcategories ALTER COLUMN inca_id SET NOT NULL;
ALTER TABLE incident_subcategories ALTER COLUMN counter_uses SET NOT NULL;
ALTER TABLE incident_subcategories ALTER COLUMN default_duration_hours SET NOT NULL;

-- ============================================================================
-- 5. TABLA: notifications
-- ============================================================================

-- Agregar NOT NULL constraints
ALTER TABLE notifications ALTER COLUMN must_send_as_notification_push SET NOT NULL;
ALTER TABLE notifications ALTER COLUMN must_send_as_notification SET NOT NULL;
ALTER TABLE notifications ALTER COLUMN must_be_processed SET NOT NULL;
ALTER TABLE notifications ALTER COLUMN retry_count SET NOT NULL;

-- ============================================================================
-- 6. TABLA: notification_deliveries
-- ============================================================================

-- Agregar NOT NULL constraints
ALTER TABLE notification_deliveries ALTER COLUMN is_read SET NOT NULL;
ALTER TABLE notification_deliveries ALTER COLUMN to_account_id SET NOT NULL;
ALTER TABLE notification_deliveries ALTER COLUMN noti_id SET NOT NULL;

-- ============================================================================
-- 7. TABLA: account_favorite_locations
-- ============================================================================

-- Agregar NOT NULL constraints
ALTER TABLE account_favorite_locations ALTER COLUMN account_id SET NOT NULL;
ALTER TABLE account_favorite_locations ALTER COLUMN status SET NOT NULL;
ALTER TABLE account_favorite_locations ALTER COLUMN crime SET NOT NULL;
ALTER TABLE account_favorite_locations ALTER COLUMN traffic_accident SET NOT NULL;
ALTER TABLE account_favorite_locations ALTER COLUMN medical_emergency SET NOT NULL;
ALTER TABLE account_favorite_locations ALTER COLUMN fire_incident SET NOT NULL;
ALTER TABLE account_favorite_locations ALTER COLUMN vandalism SET NOT NULL;
ALTER TABLE account_favorite_locations ALTER COLUMN suspicious_activity SET NOT NULL;
ALTER TABLE account_favorite_locations ALTER COLUMN infrastructure_issues SET NOT NULL;
ALTER TABLE account_favorite_locations ALTER COLUMN extreme_weather SET NOT NULL;
ALTER TABLE account_favorite_locations ALTER COLUMN community_events SET NOT NULL;
ALTER TABLE account_favorite_locations ALTER COLUMN dangerous_wildlife_sighting SET NOT NULL;
ALTER TABLE account_favorite_locations ALTER COLUMN positive_actions SET NOT NULL;
ALTER TABLE account_favorite_locations ALTER COLUMN lost_pet SET NOT NULL;
ALTER TABLE account_favorite_locations ALTER COLUMN radius SET NOT NULL;

-- ============================================================================
-- 8. TABLA: incident_votes
-- ============================================================================

ALTER TABLE incident_votes ALTER COLUMN account_id SET NOT NULL;
ALTER TABLE incident_votes ALTER COLUMN inre_id SET NOT NULL;

-- ============================================================================
-- 9. TABLA: account_session_history
-- ============================================================================

ALTER TABLE account_session_history ALTER COLUMN account_id SET NOT NULL;

-- ============================================================================
-- 10. TABLA: account_cluster_saved
-- ============================================================================

ALTER TABLE account_cluster_saved ALTER COLUMN incl_id SET NOT NULL;
ALTER TABLE account_cluster_saved ALTER COLUMN account_id SET NOT NULL;

-- ============================================================================
-- 11. TABLA: account_reports
-- ============================================================================

ALTER TABLE account_reports ALTER COLUMN account_id_whos_reporting SET NOT NULL;
ALTER TABLE account_reports ALTER COLUMN account_id SET NOT NULL;

-- ============================================================================
-- 12. TABLA: account_history
-- ============================================================================

ALTER TABLE account_history ALTER COLUMN account_id SET NOT NULL;
ALTER TABLE account_history ALTER COLUMN incl_id SET NOT NULL;

-- ============================================================================
-- 13. TABLA: device_tokens
-- ============================================================================

ALTER TABLE device_tokens ALTER COLUMN account_id SET NOT NULL;

-- ============================================================================
-- 14. TABLA: feedback
-- ============================================================================

ALTER TABLE feedback ALTER COLUMN account_id SET NOT NULL;

-- ============================================================================
-- 15. TABLA: incident_flags
-- ============================================================================

ALTER TABLE incident_flags ALTER COLUMN inre_id SET NOT NULL;
ALTER TABLE incident_flags ALTER COLUMN account_id SET NOT NULL;

-- ============================================================================
-- 16. TABLA: incident_logs
-- ============================================================================

ALTER TABLE incident_logs ALTER COLUMN inre_id SET NOT NULL;
ALTER TABLE incident_logs ALTER COLUMN account_id SET NOT NULL;

-- ============================================================================
-- 17. TABLA: account_premium_payment_history
-- ============================================================================

ALTER TABLE account_premium_payment_history ALTER COLUMN account_id SET NOT NULL;

-- ============================================================================
-- 18. TABLA: account_achievements
-- ============================================================================

ALTER TABLE account_achievements ALTER COLUMN account_id SET NOT NULL;

-- ============================================================================
-- 19. TABLA: influencers (referral system)
-- ============================================================================

ALTER TABLE influencers ALTER COLUMN web_influencer_id SET NOT NULL;
ALTER TABLE influencers ALTER COLUMN referral_code SET NOT NULL;
ALTER TABLE influencers ALTER COLUMN name SET NOT NULL;

-- ============================================================================
-- 20. TABLA: referral_conversions
-- ============================================================================

ALTER TABLE referral_conversions ALTER COLUMN referral_code SET NOT NULL;
ALTER TABLE referral_conversions ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE referral_conversions ALTER COLUMN registered_at SET NOT NULL;

-- ============================================================================
-- 21. TABLA: referral_premium_conversions
-- ============================================================================

ALTER TABLE referral_premium_conversions ALTER COLUMN referral_code SET NOT NULL;
ALTER TABLE referral_premium_conversions ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE referral_premium_conversions ALTER COLUMN amount SET NOT NULL;
ALTER TABLE referral_premium_conversions ALTER COLUMN commission SET NOT NULL;
ALTER TABLE referral_premium_conversions ALTER COLUMN converted_at SET NOT NULL;

-- ============================================================================
-- VERIFICACIONES FINALES
-- ============================================================================

-- Verificar que los cambios se aplicaron correctamente
SELECT
    table_name,
    column_name,
    is_nullable,
    column_default,
    data_type
FROM information_schema.columns
WHERE table_schema = 'public'
  AND table_name IN (
    'account',
    'incident_clusters',
    'incident_reports',
    'incident_subcategories',
    'notifications',
    'notification_deliveries',
    'account_favorite_locations'
  )
  AND column_name IN (
    'account_id', 'email', 'password', 'nickname', 'role', 'status',
    'insu_id', 'incident_count', 'is_active', 'credibility',
    'counter_total_votes', 'counter_total_comments'
  )
ORDER BY table_name, column_name;

-- ============================================================================
-- FIN DEL SCRIPT
-- ============================================================================
