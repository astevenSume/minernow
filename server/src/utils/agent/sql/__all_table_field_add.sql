----------------------------------------------------
--  `agent`
----------------------------------------------------
ALTER TABLE `agent` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `agent` ADD `sum_salary` bigint NOT NULL COMMENT 'sum_salary' AFTER `uid`;
ALTER TABLE `agent` ADD `sum_withdraw` bigint NOT NULL COMMENT 'sum_withdraw' AFTER `sum_salary`;
ALTER TABLE `agent` ADD `sum_can_withdraw` bigint NOT NULL COMMENT 'sum_can_withdraw' AFTER `sum_withdraw`;
ALTER TABLE `agent` ADD `ctime` bigint NOT NULL COMMENT '' AFTER `sum_can_withdraw`;
ALTER TABLE `agent` ADD `mtime` bigint NOT NULL COMMENT '' AFTER `ctime`;
ALTER TABLE `agent` ADD `pwd` varchar(64) NOT NULL COMMENT 'password' AFTER `mtime`;

----------------------------------------------------
--  `agent_channel_commission`
----------------------------------------------------
ALTER TABLE `agent_channel_commission` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `agent_channel_commission` ADD `channel_id` int unsigned NOT NULL COMMENT 'channel id' AFTER `uid`;
ALTER TABLE `agent_channel_commission` ADD `ctime` bigint NOT NULL COMMENT '' AFTER `channel_id`;
ALTER TABLE `agent_channel_commission` ADD `mtime` bigint NOT NULL COMMENT '' AFTER `ctime`;
ALTER TABLE `agent_channel_commission` ADD `integer` int NOT NULL COMMENT 'commission integer part' AFTER `mtime`;
ALTER TABLE `agent_channel_commission` ADD `decimals` int NOT NULL COMMENT 'commission decimals part' AFTER `integer`;
ALTER TABLE `agent_channel_commission` ADD `status` tinyint unsigned NOT NULL COMMENT 'commission status' AFTER `decimals`;

----------------------------------------------------
--  `agent_path`
----------------------------------------------------
ALTER TABLE `agent_path` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `agent_path` ADD `level` int unsigned NOT NULL COMMENT 'agent level' AFTER `uid`;
ALTER TABLE `agent_path` ADD `sn` int unsigned NOT NULL COMMENT 'agent serial number on specific level' AFTER `level`;
ALTER TABLE `agent_path` ADD `path` text NOT NULL COMMENT 'user agent path' AFTER `sn`;
ALTER TABLE `agent_path` ADD `ctime` bigint NOT NULL COMMENT '' AFTER `path`;
ALTER TABLE `agent_path` ADD `mtime` bigint NOT NULL COMMENT '' AFTER `ctime`;
ALTER TABLE `agent_path` ADD `invite_code` varchar(100) NOT NULL COMMENT 'invite code' AFTER `mtime`;
ALTER TABLE `agent_path` ADD `whitelist_id` int unsigned NOT NULL COMMENT 'agent commission whitelist id' AFTER `invite_code`;
ALTER TABLE `agent_path` ADD `invite_num` int unsigned NOT NULL COMMENT 'invite number' AFTER `whitelist_id`;
ALTER TABLE `agent_path` ADD `parent_uid` bigint unsigned NOT NULL COMMENT 'parent uid' AFTER `invite_num`;
ALTER TABLE `agent_path` ADD `dividend_position` int unsigned NOT NULL DEFAULT 0 COMMENT 'month dividend position' AFTER `parent_uid`;

----------------------------------------------------
--  `agent_withdraw`
----------------------------------------------------
ALTER TABLE `agent_withdraw` ADD `id` bigint unsigned NOT NULL COMMENT 'agent withdraw id';
ALTER TABLE `agent_withdraw` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id' AFTER `id`;
ALTER TABLE `agent_withdraw` ADD `amount` bigint NOT NULL COMMENT 'amount' AFTER `uid`;
ALTER TABLE `agent_withdraw` ADD `status` tinyint unsigned NOT NULL COMMENT 'status' AFTER `amount`;
ALTER TABLE `agent_withdraw` ADD `ctime` bigint NOT NULL COMMENT '' AFTER `status`;
ALTER TABLE `agent_withdraw` ADD `mtime` bigint NOT NULL COMMENT '' AFTER `ctime`;
ALTER TABLE `agent_withdraw` ADD `desc` varchar(256) NULL COMMENT '' AFTER `mtime`;

----------------------------------------------------
--  `invite_code`
----------------------------------------------------
ALTER TABLE `invite_code` ADD `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `invite_code` ADD `code` varchar(16) NOT NULL COMMENT 'code' AFTER `id`;
ALTER TABLE `invite_code` ADD `status` tinyint unsigned NOT NULL COMMENT 'status' AFTER `code`;

