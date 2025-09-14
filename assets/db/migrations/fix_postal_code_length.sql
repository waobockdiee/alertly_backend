-- Fix postal_code column length for better compatibility
-- This migration increases the postal_code field length to accommodate longer postal codes

-- Update incident_reports table
ALTER TABLE `alertly`.`incident_reports` 
MODIFY COLUMN `postal_code` VARCHAR(12) NULL;

-- Update incident_clusters table (already has VARCHAR(12), but ensuring consistency)
ALTER TABLE `alertly`.`incident_clusters` 
MODIFY COLUMN `postal_code` VARCHAR(12) NULL;

-- Update account_favorite_locations table
ALTER TABLE `alertly`.`account_favorite_locations` 
MODIFY COLUMN `postal_code` VARCHAR(12) NULL;

-- Add comment for documentation
ALTER TABLE `alertly`.`incident_reports` 
MODIFY COLUMN `postal_code` VARCHAR(12) NULL COMMENT 'Postal/ZIP code - increased to 12 chars for international codes';

ALTER TABLE `alertly`.`incident_clusters` 
MODIFY COLUMN `postal_code` VARCHAR(12) NULL COMMENT 'Postal/ZIP code - increased to 12 chars for international codes';

ALTER TABLE `alertly`.`account_favorite_locations` 
MODIFY COLUMN `postal_code` VARCHAR(12) NULL COMMENT 'Postal/ZIP code - increased to 12 chars for international codes';
