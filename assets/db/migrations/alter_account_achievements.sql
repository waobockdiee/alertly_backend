ALTER TABLE `alertly`.`account_achievements`
ADD COLUMN `badge_threshold` INT UNSIGNED NULL AFTER `icon_url`,
MODIFY COLUMN `description` TEXT NULL;

ALTER TABLE `alertly`.`account_achievements`
ADD UNIQUE INDEX `uq_account_badge` (`account_id` ASC, `type` ASC, `name` ASC) VISIBLE;