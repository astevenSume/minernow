
-- --------------------------------------------------
--  Table Structure for `models.EosAccountKeys`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `eos_account_keys` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`account` varchar(100) NOT NULL COMMENT 'account',
`public_key_owner` varchar(100) NOT NULL COMMENT 'public key owner',
`private_key_owner` varchar(100) NOT NULL COMMENT 'private key owner',
`public_key_active` varchar(100) NOT NULL COMMENT 'public key active',
`private_key_active` varchar(100) NOT NULL COMMENT 'private key active',
`ctime` bigint NOT NULL COMMENT 'created time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='Eos Account keys' DEFAULT CHARSET=utf8;
CREATE UNIQUE INDEX `eos_account_keys_account` ON `eos_account_keys` (`account`);

