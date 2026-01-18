-- ============================================================================
-- POSTGRESQL SCHEMA FIX - Restaurar Constraints y Defaults desde MySQL
-- ============================================================================
-- Fecha: 2026-01-17
-- Propósito: Corregir la migración de MySQL a PostgreSQL Railway
-- EJECUTAR: psql "DATABASE_URL" -f fix_postgresql_constraints_final.sql
-- ============================================================================

BEGIN;

-- ============================================================================
-- PASO 1: Limpiar datos NULL antes de aplicar NOT NULL constraints
-- ============================================================================

-- incident_reports: vote column (ya limpiado, pero por si acaso)
UPDATE incident_reports SET vote = 1 WHERE vote IS NULL;

-- ============================================================================
-- PASO 2: Aplicar NOT NULL constraints en tablas principales
-- ============================================================================

-- TABLA: account
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

-- TABLA: incident_clusters
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

-- TABLA: incident_reports
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

-- TABLA: incident_subcategories
ALTER TABLE incident_subcategories ALTER COLUMN counter_uses SET NOT NULL;
ALTER TABLE incident_subcategories ALTER COLUMN default_duration_hours SET NOT NULL;

-- TABLA: notifications
ALTER TABLE notifications ALTER COLUMN must_send_as_notification_push SET NOT NULL;
ALTER TABLE notifications ALTER COLUMN must_send_as_notification SET NOT NULL;
ALTER TABLE notifications ALTER COLUMN must_be_processed SET NOT NULL;
ALTER TABLE notifications ALTER COLUMN retry_count SET NOT NULL;

-- TABLA: notification_deliveries
ALTER TABLE notification_deliveries ALTER COLUMN is_read SET NOT NULL;

-- TABLA: account_favorite_locations
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

-- TABLA: influencers (referral system)
ALTER TABLE influencers ALTER COLUMN web_influencer_id SET NOT NULL;
ALTER TABLE influencers ALTER COLUMN referral_code SET NOT NULL;
ALTER TABLE influencers ALTER COLUMN name SET NOT NULL;

-- TABLA: referral_conversions
ALTER TABLE referral_conversions ALTER COLUMN referral_code SET NOT NULL;
ALTER TABLE referral_conversions ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE referral_conversions ALTER COLUMN registered_at SET NOT NULL;

-- TABLA: referral_premium_conversions
ALTER TABLE referral_premium_conversions ALTER COLUMN referral_code SET NOT NULL;
ALTER TABLE referral_premium_conversions ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE referral_premium_conversions ALTER COLUMN amount SET NOT NULL;
ALTER TABLE referral_premium_conversions ALTER COLUMN commission SET NOT NULL;
ALTER TABLE referral_premium_conversions ALTER COLUMN converted_at SET NOT NULL;

COMMIT;

-- ============================================================================
-- VERIFICACIÓN: Ver resumen de cambios aplicados
-- ============================================================================

SELECT
    'VERIFICACIÓN COMPLETA' as mensaje,
    COUNT(*) as total_columnas_verificadas
FROM information_schema.columns
WHERE table_schema = 'public'
  AND table_name IN (
    'account',
    'incident_clusters',
    'incident_reports',
    'incident_subcategories',
    'notifications',
    'notification_deliveries',
    'account_favorite_locations',
    'influencers',
    'referral_conversions',
    'referral_premium_conversions'
  )
  AND is_nullable = 'NO';

-- ============================================================================
-- FIN DEL SCRIPT - Constraints aplicados correctamente
-- ============================================================================
