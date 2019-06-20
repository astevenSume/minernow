----------------------------------------------------
--  `eos_account_keys`
----------------------------------------------------
ALTER TABLE `eos_account_keys` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `eos_account_keys` ADD `account` varchar(100) NOT NULL COMMENT 'account' AFTER `id`;
ALTER TABLE `eos_account_keys` ADD `public_key_owner` varchar(100) NOT NULL COMMENT 'public key owner' AFTER `account`;
ALTER TABLE `eos_account_keys` ADD `private_key_owner` varchar(100) NOT NULL COMMENT 'private key owner' AFTER `public_key_owner`;
ALTER TABLE `eos_account_keys` ADD `public_key_active` varchar(100) NOT NULL COMMENT 'public key active' AFTER `private_key_owner`;
ALTER TABLE `eos_account_keys` ADD `private_key_active` varchar(100) NOT NULL COMMENT 'private key active' AFTER `public_key_active`;
ALTER TABLE `eos_account_keys` ADD `ctime` bigint NOT NULL COMMENT 'created time' AFTER `private_key_active`;

