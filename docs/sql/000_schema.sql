-- Create Table User
DROP TABLE IF EXISTS `user`;
CREATE TABLE IF NOT EXISTS `user` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `fk_role_id` INT COMMENT 'Foreign Key To role Id',
  `email` VARCHAR(255) NOT NULL DEFAULT '',
  `username` VARCHAR(255) NOT NULL DEFAULT '',
  `password` VARCHAR(255) NOT NULL DEFAULT '',
  `display_name` VARCHAR(255) NOT NULL DEFAULT '',

  -- Utility columns
  `status` SMALLINT NOT NULL DEFAULT '1',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` VARCHAR(255),
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_by` VARCHAR(255),
  `deleted_at`TIMESTAMP,
  `deleted_by` VARCHAR(255),
  PRIMARY KEY (`id`),
  UNIQUE (`username`)
) ENGINE = INNODB COMMENT='User table';

-- Create Table Url
DROP TABLE IF EXISTS `url`;
CREATE TABLE IF NOT EXISTS `url` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `original_url` VARCHAR(255) NOT NULL DEFAULT '',
  `shorten_url` VARCHAR(255) NOT NULL DEFAULT '',
  `visit` INT NOT NULL DEFAULT 0,
  `user_id` INT NOT NULL,

  -- Utility columns
  `status` SMALLINT NOT NULL DEFAULT '1',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` VARCHAR(255),
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_by` VARCHAR(255),
  `deleted_at`TIMESTAMP,
  `deleted_by` VARCHAR(255),
  PRIMARY KEY (`id`),
  UNIQUE (`shorten_url`)
) ENGINE = INNODB COMMENT='url table';

-- Create Table Role
DROP TABLE IF EXISTS `role`;
CREATE TABLE IF NOT EXISTS `role` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL DEFAULT '',
  `type` VARCHAR(255) NOT NULL DEFAULT '',
  `rank` INT NOT NULL DEFAULT 0,

  -- Utility columns
  `status` SMALLINT NOT NULL DEFAULT '1',
  `flag` INT NOT NULL DEFAULT '0',
  `meta` VARCHAR(255),
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` VARCHAR(255),
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_by` VARCHAR(255),
  `deleted_at`TIMESTAMP,
  `deleted_by` VARCHAR(255),
  PRIMARY KEY (`id`)
);

-- Create Table Action
DROP TABLE IF EXISTS `action`;
CREATE TABLE IF NOT EXISTS `action` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `code` VARCHAR(255) NOT NULL DEFAULT '',
    `description` VARCHAR(255) NOT NULL DEFAULT '',

    -- Utility columns
    `status` SMALLINT NOT NULL DEFAULT '1',
    `flag` INT NOT NULL DEFAULT '0',
    `meta` VARCHAR(255),
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` VARCHAR(255),
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `updated_by` VARCHAR(255),
    `deleted_at`TIMESTAMP,
    `deleted_by` VARCHAR(255),
    PRIMARY KEY (`id`),
    UNIQUE (`code`)
) ENGINE = INNODB COMMENT='all actions that can be performed';

-- Create Table Role-Action
DROP TABLE IF EXISTS `role_action`;
CREATE TABLE IF NOT EXISTS `role_action` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `fk_role_id` INT NOT NULL,
    `fk_action_id` INT NOT NULL,

    -- Utility columns
    `status` SMALLINT NOT NULL DEFAULT '1',
    `flag` INT NOT NULL DEFAULT '0',
    `meta` VARCHAR(255),
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` VARCHAR(255),
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `updated_by` VARCHAR(255),
    `deleted_at`TIMESTAMP,
    `deleted_by` VARCHAR(255),
    PRIMARY KEY (`id`),
    UNIQUE `role_action_unique_index` (`fk_role_id`, `fk_action_id`)
) ENGINE = INNODB COMMENT='relation between role and action';
