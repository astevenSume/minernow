----------------------------------------------------
--  `eos_account_keys`
----------------------------------------------------
ALTER TABLE `eos_account_keys` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `eos_account_keys` CHANGE `account` `account` varchar(100) NOT NULL COMMENT 'account';
ALTER TABLE `eos_account_keys` CHANGE `public_key_owner` `public_key_owner` varchar(100) NOT NULL COMMENT 'public key owner';
ALTER TABLE `eos_account_keys` CHANGE `private_key_owner` `private_key_owner` varchar(100) NOT NULL COMMENT 'private key owner';
ALTER TABLE `eos_account_keys` CHANGE `public_key_active` `public_key_active` varchar(100) NOT NULL COMMENT 'public key active';
ALTER TABLE `eos_account_keys` CHANGE `private_key_active` `private_key_active` varchar(100) NOT NULL COMMENT 'private key active';
ALTER TABLE `eos_account_keys` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'created time';

