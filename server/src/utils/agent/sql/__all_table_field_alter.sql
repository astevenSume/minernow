----------------------------------------------------
--  `agent`
----------------------------------------------------
ALTER TABLE `agent` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `agent` CHANGE `sum_salary` `sum_salary` bigint NOT NULL COMMENT 'sum_salary';
ALTER TABLE `agent` CHANGE `sum_withdraw` `sum_withdraw` bigint NOT NULL COMMENT 'sum_withdraw';
ALTER TABLE `agent` CHANGE `sum_can_withdraw` `sum_can_withdraw` bigint NOT NULL COMMENT 'sum_can_withdraw';
ALTER TABLE `agent` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT '';
ALTER TABLE `agent` CHANGE `mtime` `mtime` bigint NOT NULL COMMENT '';
ALTER TABLE `agent` CHANGE `pwd` `pwd` varchar(64) NOT NULL COMMENT 'password';

----------------------------------------------------
--  `agent_channel_commission`
----------------------------------------------------
ALTER TABLE `agent_channel_commission` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `agent_channel_commission` CHANGE `channel_id` `channel_id` int unsigned NOT NULL COMMENT 'channel id';
ALTER TABLE `agent_channel_commission` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT '';
ALTER TABLE `agent_channel_commission` CHANGE `mtime` `mtime` bigint NOT NULL COMMENT '';
ALTER TABLE `agent_channel_commission` CHANGE `integer` `integer` int NOT NULL COMMENT 'commission integer part';
ALTER TABLE `agent_channel_commission` CHANGE `decimals` `decimals` int NOT NULL COMMENT 'commission decimals part';
ALTER TABLE `agent_channel_commission` CHANGE `status` `status` tinyint unsigned NOT NULL COMMENT 'commission status';

----------------------------------------------------
--  `agent_path`
----------------------------------------------------
ALTER TABLE `agent_path` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `agent_path` CHANGE `level` `level` int unsigned NOT NULL COMMENT 'agent level';
ALTER TABLE `agent_path` CHANGE `sn` `sn` int unsigned NOT NULL COMMENT 'agent serial number on specific level';
ALTER TABLE `agent_path` CHANGE `path` `path` text NOT NULL COMMENT 'user agent path';
ALTER TABLE `agent_path` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT '';
ALTER TABLE `agent_path` CHANGE `mtime` `mtime` bigint NOT NULL COMMENT '';
ALTER TABLE `agent_path` CHANGE `invite_code` `invite_code` varchar(100) NOT NULL COMMENT 'invite code';
ALTER TABLE `agent_path` CHANGE `whitelist_id` `whitelist_id` int unsigned NOT NULL COMMENT 'agent commission whitelist id';
ALTER TABLE `agent_path` CHANGE `invite_num` `invite_num` int unsigned NOT NULL COMMENT 'invite number';
ALTER TABLE `agent_path` CHANGE `parent_uid` `parent_uid` bigint unsigned NOT NULL COMMENT 'parent uid';
ALTER TABLE `agent_path` CHANGE `dividend_position` `dividend_position` int unsigned NOT NULL DEFAULT 0 COMMENT 'month dividend position';

----------------------------------------------------
--  `agent_withdraw`
----------------------------------------------------
ALTER TABLE `agent_withdraw` CHANGE `id` `id` bigint unsigned NOT NULL COMMENT 'agent withdraw id';
ALTER TABLE `agent_withdraw` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `agent_withdraw` CHANGE `amount` `amount` bigint NOT NULL COMMENT 'amount';
ALTER TABLE `agent_withdraw` CHANGE `status` `status` tinyint unsigned NOT NULL COMMENT 'status';
ALTER TABLE `agent_withdraw` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT '';
ALTER TABLE `agent_withdraw` CHANGE `mtime` `mtime` bigint NOT NULL COMMENT '';
ALTER TABLE `agent_withdraw` CHANGE `desc` `desc` varchar(256) NULL COMMENT '';

----------------------------------------------------
--  `invite_code`
----------------------------------------------------
ALTER TABLE `invite_code` CHANGE `id` `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `invite_code` CHANGE `code` `code` varchar(16) NOT NULL COMMENT 'code';
ALTER TABLE `invite_code` CHANGE `status` `status` tinyint unsigned NOT NULL COMMENT 'status';

