ALTER TABLE `alertly`.`account` 
ADD COLUMN `device_token` VARCHAR(255) NULL AFTER `counter_new_notifications`,
ADD COLUMN `last_active` TIMESTAMP NULL AFTER `device_token`;