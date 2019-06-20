----------------------------------------------------
--  `eos_account_keys`
----------------------------------------------------
ALTER TABLE `eos_account_keys` DROP `id`;
ALTER TABLE `eos_account_keys` DROP `account`;
ALTER TABLE `eos_account_keys` DROP `public_key_owner`;
ALTER TABLE `eos_account_keys` DROP `private_key_owner`;
ALTER TABLE `eos_account_keys` DROP `public_key_active`;
ALTER TABLE `eos_account_keys` DROP `private_key_active`;
ALTER TABLE `eos_account_keys` DROP `ctime`;

