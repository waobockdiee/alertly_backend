-- Performance Indexes for Alertly Database
-- Execute this script to add optimized indexes for better performance

USE alertly;

-- ✅ Índice compuesto para geolocalización de clusters
-- Optimiza queries de getclustersbylocation y getclusterbyradius
CREATE INDEX IF NOT EXISTS idx_clusters_location_time ON incident_clusters 
(center_latitude, center_longitude, start_time, end_time, is_active);

-- ✅ Índice para clustering de incidentes
-- Optimiza la búsqueda de clusters existentes por subcategoría y tiempo
CREATE INDEX IF NOT EXISTS idx_clusters_insu_created ON incident_clusters 
(insu_id, created_at, is_active);

-- ✅ Índice para votos de incidentes
-- Optimiza la verificación de si un usuario ya votó
CREATE INDEX IF NOT EXISTS idx_reports_incl_account ON incident_reports 
(incl_id, account_id, is_active);

-- ✅ Índice para búsqueda de incidentes por usuario
-- Optimiza queries de historial de usuario
CREATE INDEX IF NOT EXISTS idx_reports_account_created ON incident_reports 
(account_id, created_at DESC);

-- ✅ Índice para lugares favoritos
-- Optimiza búsqueda de lugares guardados por usuario
CREATE INDEX IF NOT EXISTS idx_favorite_locations_account ON account_favorite_locations 
(account_id, status);

-- ✅ Índice para notificaciones
-- Optimiza queries de notificaciones pendientes
CREATE INDEX IF NOT EXISTS idx_notifications_owner_processed ON notifications 
(owner_account_id, must_be_processed, created_at);

-- ✅ Índice para device tokens
-- Optimiza búsqueda de tokens por usuario
CREATE INDEX IF NOT EXISTS idx_device_tokens_account ON device_tokens 
(account_id, device_token);

-- ✅ Índice para comentarios
-- Optimiza queries de comentarios por cluster
CREATE INDEX IF NOT EXISTS idx_comments_cluster_created ON cluster_comments 
(incl_id, created_at DESC);

-- ✅ Índice para reportes de usuarios
-- Optimiza queries de moderación
CREATE INDEX IF NOT EXISTS idx_account_reports_target ON account_reports 
(account_id, created_at DESC);

-- ✅ Índice para historial de cuenta
-- Optimiza queries de historial de usuario
CREATE INDEX IF NOT EXISTS idx_account_history_user ON account_history 
(account_id, created_at DESC);

-- ✅ Índice para clusters guardados
-- Optimiza queries de incidentes guardados
CREATE INDEX IF NOT EXISTS idx_cluster_saved_account ON account_cluster_saved 
(account_id, created_at DESC);

-- ✅ Índice para categorías (aunque son pocas, mejora performance)
CREATE INDEX IF NOT EXISTS idx_categories_code ON incident_categories 
(code);

-- ✅ Índice para subcategorías
CREATE INDEX IF NOT EXISTS idx_subcategories_category ON incident_subcategories 
(inca_id, code);

-- ✅ Índice para credibilidad de usuarios
-- Optimiza cálculos de credibilidad
CREATE INDEX IF NOT EXISTS idx_account_credibility ON account 
(credibility DESC, status);

-- ✅ Índice para usuarios premium
-- Optimiza queries de usuarios premium
CREATE INDEX IF NOT EXISTS idx_account_premium ON account 
(is_premium, status);

-- ✅ Índice para búsqueda por email
-- Optimiza login y validaciones
CREATE INDEX IF NOT EXISTS idx_account_email ON account 
(email, status);

-- ✅ Índice para activación de cuenta
-- Optimiza queries de activación
CREATE INDEX IF NOT EXISTS idx_account_activation ON account 
(activation_code, status);

-- ✅ Índice para notificaciones push
-- Optimiza envío de notificaciones
CREATE INDEX IF NOT EXISTS idx_notification_deliveries_account ON notification_deliveries 
(to_account_id, is_read, created_at DESC);

-- Verificar que los índices se crearon correctamente
SHOW INDEX FROM incident_clusters;
SHOW INDEX FROM incident_reports;
SHOW INDEX FROM account_favorite_locations;
SHOW INDEX FROM notifications;
SHOW INDEX FROM device_tokens;

-- Estadísticas de performance (ejecutar después de crear los índices)
ANALYZE TABLE incident_clusters;
ANALYZE TABLE incident_reports;
ANALYZE TABLE account_favorite_locations;
ANALYZE TABLE notifications;
ANALYZE TABLE device_tokens;
