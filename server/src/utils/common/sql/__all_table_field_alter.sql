----------------------------------------------------
--  `token`
----------------------------------------------------
ALTER TABLE `token` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `token` CHANGE `client_type` `client_type` int unsigned NOT NULL COMMENT 'app type';
ALTER TABLE `token` CHANGE `mtime` `mtime` bigint NOT NULL COMMENT 'access_token modify time';
ALTER TABLE `token` CHANGE `access_token` `access_token` varchar(256) NOT NULL COMMENT 'access_token';
ALTER TABLE `token` CHANGE `mac` `mac` varchar(100) NOT NULL COMMENT '';

