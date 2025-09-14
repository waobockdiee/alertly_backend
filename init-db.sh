#!/bin/bash

# Script para inicializar la base de datos con el esquema correcto

DB_HOST="alertly-main-db.cluster-c3qmq4y86s84.us-west-2.rds.amazonaws.com"
DB_USER="adminalertly"
DB_PASS="Po1Ng2O3;"
DB_NAME="alertly"

echo "🔧 Inicializando base de datos en $DB_HOST..."

# Crear un archivo SQL temporal con el esquema completo
cat > /tmp/init-alertly.sql << 'EOF'
-- Usar la base de datos
USE alertly;

-- Eliminar tablas si existen (para reinicio limpio)
SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS incident_categories;
DROP TABLE IF EXISTS incident_subcategories;
SET FOREIGN_KEY_CHECKS = 1;

-- Crear tabla incident_categories con la estructura correcta del código
CREATE TABLE IF NOT EXISTS `incident_categories` (
  `inca_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NOT NULL,
  `description` TEXT NULL,
  `icon` VARCHAR(255) NULL,  -- El código espera 'icon', no 'icon_uri'
  `code` VARCHAR(45) NULL,
  `border_color` VARCHAR(7) NULL,
  `default_circle_range` INT NULL,
  `max_circle_range` INT NULL,
  PRIMARY KEY (`inca_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Insertar categorías de ejemplo
INSERT INTO `incident_categories` (`inca_id`, `name`, `description`, `icon`, `code`, `border_color`) VALUES
(1, 'Crime', 'Incidents involving illegal activities', 'crime', 'crime', '#FF0000'),
(2, 'Traffic Accident', 'Incidents on roadways', 'traffic_accident', 'traffic_accident', '#FFA500'),
(3, 'Medical Emergency', 'Urgent health situations', 'medical_emergency', 'medical_emergency', '#00FF00'),
(4, 'Fire', 'Fire incidents', 'fire_incident', 'fire_incident', '#FF4500'),
(5, 'Vandalism', 'Property damage', 'vandalism', 'vandalism', '#800080'),
(6, 'Suspicious Activity', 'Unusual behavior', 'suspicious_activity', 'suspicious_activity', '#FFFF00'),
(7, 'Infrastructure Issue', 'Public infrastructure problems', 'infrastructure_issues', 'infrastructure_issues', '#808080'),
(8, 'Extreme Weather', 'Severe weather conditions', 'extreme_weather', 'extreme_weather', '#0000FF'),
(9, 'Community Event', 'Local events', 'community_events', 'community_events', '#00FFFF'),
(10, 'Dangerous Wildlife', 'Wildlife encounters', 'dangerous_wildlife_sighting', 'dangerous_wildlife_sighting', '#8B4513'),
(11, 'Positive Actions', 'Good deeds', 'positive_actions', 'positive_actions', '#32CD32'),
(12, 'Lost Pet', 'Missing pets', 'lost_pet', 'lost_pet', '#FF69B4');

-- Verificar que se crearon correctamente
SELECT COUNT(*) as total_categories FROM incident_categories;
EOF

echo "📝 Archivo SQL creado en /tmp/init-alertly.sql"
echo "✅ Base de datos lista para inicialización"

# Nota: Para ejecutar manualmente:
# mysql -h $DB_HOST -u $DB_USER -p"$DB_PASS" < /tmp/init-alertly.sql

