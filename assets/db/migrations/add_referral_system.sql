-- =====================================================
-- ALERTLY REFERRAL SYSTEM - DATABASE MIGRATION
-- =====================================================
-- Fecha: 2025-11-02
-- Propósito: Sistema de referrals para influencers
-- Tablas: influencers, referral_conversions,
--         referral_premium_conversions, referral_metrics_cache
-- =====================================================

USE `alertly`;

-- =====================================================
-- Tabla 1: influencers
-- Almacena información de influencers/marketers
-- =====================================================
CREATE TABLE IF NOT EXISTS `influencers` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `web_influencer_id` INT UNSIGNED NOT NULL COMMENT 'ID de marketing_influencers del backend web',
    `referral_code` VARCHAR(20) NOT NULL COMMENT 'Código único del influencer (ej: INF-IG0001)',
    `name` VARCHAR(255) NOT NULL COMMENT 'Nombre del influencer',
    `platform` ENUM('Instagram', 'TikTok', 'Reddit', 'Other') NOT NULL COMMENT 'Plataforma principal del influencer',
    `is_active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '1=activo, 0=inactivo',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`),
    UNIQUE KEY `unique_referral_code` (`referral_code`),
    INDEX `idx_referral_code` (`referral_code`),
    INDEX `idx_is_active` (`is_active`),
    INDEX `idx_web_influencer_id` (`web_influencer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Tabla de influencers que promocionan la app';

-- =====================================================
-- Tabla 2: referral_conversions
-- Registra cada usuario que se registró con código
-- =====================================================
CREATE TABLE IF NOT EXISTS `referral_conversions` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `referral_code` VARCHAR(20) NOT NULL COMMENT 'Código del influencer usado',
    `user_id` INT UNSIGNED NOT NULL COMMENT 'ID del usuario registrado (account.account_id)',
    `registered_at` DATETIME NOT NULL COMMENT 'Fecha y hora del registro',
    `platform` ENUM('iOS', 'Android') NOT NULL COMMENT 'Plataforma del dispositivo',
    `earnings` DECIMAL(10,2) NOT NULL DEFAULT 0.10 COMMENT 'Comisión por registro ($0.10 CAD)',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`),
    UNIQUE KEY `unique_user_referral` (`user_id`) COMMENT 'Un usuario solo puede usar UN código',
    INDEX `idx_referral_code` (`referral_code`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_registered_at` (`registered_at`),

    FOREIGN KEY (`user_id`) REFERENCES `account` (`account_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Conversiones de registro con código de referral';

-- =====================================================
-- Tabla 3: referral_premium_conversions
-- Registra suscripciones premium de usuarios referidos
-- =====================================================
CREATE TABLE IF NOT EXISTS `referral_premium_conversions` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `referral_code` VARCHAR(20) NOT NULL COMMENT 'Código del influencer',
    `user_id` INT UNSIGNED NOT NULL COMMENT 'ID del usuario que se suscribió',
    `conversion_id` BIGINT UNSIGNED NULL COMMENT 'FK a referral_conversions',
    `subscription_type` ENUM('monthly', 'yearly') NOT NULL COMMENT 'Tipo de suscripción',
    `amount` DECIMAL(10,2) NOT NULL COMMENT 'Monto de la suscripción en CAD',
    `commission` DECIMAL(10,2) NOT NULL COMMENT 'Comisión calculada (15% del amount)',
    `commission_percentage` DECIMAL(5,2) NOT NULL DEFAULT 15.00 COMMENT 'Porcentaje de comisión',
    `converted_at` DATETIME NOT NULL COMMENT 'Fecha y hora de la conversión',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`),
    INDEX `idx_referral_code` (`referral_code`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_converted_at` (`converted_at`),
    INDEX `idx_conversion_id` (`conversion_id`),

    FOREIGN KEY (`user_id`) REFERENCES `account` (`account_id`) ON DELETE CASCADE,
    FOREIGN KEY (`conversion_id`) REFERENCES `referral_conversions` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Conversiones premium de usuarios referidos';

-- =====================================================
-- Tabla 4: referral_metrics_cache (OPCIONAL)
-- Cache para optimizar queries de métricas
-- =====================================================
CREATE TABLE IF NOT EXISTS `referral_metrics_cache` (
    `referral_code` VARCHAR(20) NOT NULL,
    `total_registrations` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Total de registros',
    `total_premium_conversions` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Total conversiones premium',
    `total_earnings` DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT 'Total de ganancias en CAD',
    `last_updated` DATETIME NOT NULL COMMENT 'Última actualización del cache',

    PRIMARY KEY (`referral_code`),
    INDEX `idx_last_updated` (`last_updated`),
    INDEX `idx_total_earnings` (`total_earnings`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Cache de métricas para optimizar consultas';

-- =====================================================
-- DATOS DE PRUEBA
-- =====================================================

-- Influencers de prueba
INSERT INTO `influencers` (`web_influencer_id`, `referral_code`, `name`, `platform`, `is_active`) VALUES
(1, 'INF-IG0001', 'John Doe', 'Instagram', 1),
(2, 'INF-TT0002', 'Jane Smith', 'TikTok', 1),
(3, 'INF-RD0003', 'Bob Johnson', 'Reddit', 1),
(4, 'INF-IG0004', 'Alice Williams', 'Instagram', 1),
(5, 'INF-TT0005', 'Charlie Brown', 'TikTok', 0)  -- Inactivo para testing
ON DUPLICATE KEY UPDATE updated_at = CURRENT_TIMESTAMP;

-- =====================================================
-- VERIFICACIÓN
-- =====================================================

-- Mostrar tablas creadas
SELECT 'Tablas de referral system creadas:' as Message;
SHOW TABLES LIKE 'influencers';
SHOW TABLES LIKE 'referral_conversions';
SHOW TABLES LIKE 'referral_premium_conversions';
SHOW TABLES LIKE 'referral_metrics_cache';

-- Contar influencers insertados
SELECT COUNT(*) as total_influencers FROM influencers;
SELECT * FROM influencers;

-- =====================================================
-- NOTAS IMPORTANTES
-- =====================================================
-- 1. API Key generado: 0560f9d085d2fe20fd7ddcd51024abb508d68f07935ad875b8769773b62d5071
-- 2. Agregar a backend/.env: REFERRAL_API_KEY=0560f9d085d2fe20fd7ddcd51024abb508d68f07935ad875b8769773b62d5071
-- 3. Compartir API Key con backend web para autenticación
-- 4. La tabla referral_metrics_cache es opcional pero recomendada
-- 5. Los índices están optimizados para las queries más frecuentes
-- =====================================================
