CREATE TABLE IF NOT EXISTS users(
    `id` BINARY(16) NOT NULL,
    `identity` VARCHAR(127) NOT NULL,
    `password` BINARY(64) NOT NULL,
    `salt` BINARY(16) NOT NULL,
    `status` ENUM('enabled', 'disabled', 'suspended') NOT NULL,
    `ctime` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `mtime` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(`id`),
    UNIQUE(`identity`)
) ENGINE InnoDB COLLATE 'utf8mb4_unicode_ci' CHARACTER SET 'utf8mb4';
