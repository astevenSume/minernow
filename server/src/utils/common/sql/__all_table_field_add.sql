----------------------------------------------------
--  `token`
----------------------------------------------------
ALTER TABLE `token` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `token` ADD `client_type` int unsigned NOT NULL COMMENT 'app type' AFTER `uid`;
ALTER TABLE `token` ADD `mtime` bigint NOT NULL COMMENT 'access_token modify time' AFTER `client_type`;
ALTER TABLE `token` ADD `access_token` varchar(256) NOT NULL COMMENT 'access_token' AFTER `mtime`;
ALTER TABLE `token` ADD `mac` varchar(100) NOT NULL COMMENT '' AFTER `access_token`;

