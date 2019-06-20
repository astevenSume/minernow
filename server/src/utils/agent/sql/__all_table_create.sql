-- --------------------------------------------------
--  Table Structure for `models.Agent`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `agent` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`sum_salary` bigint NOT NULL COMMENT 'sum_salary',
`sum_withdraw` bigint NOT NULL COMMENT 'sum_withdraw',
`sum_can_withdraw` bigint NOT NULL COMMENT 'sum_can_withdraw',
`ctime` bigint NOT NULL COMMENT '',
`mtime` bigint NOT NULL COMMENT '',
`pwd` varchar(64) NOT NULL COMMENT 'password',
PRIMARY KEY(`uid`)
) ENGINE=InnoDB COMMENT='agent table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.AgentChannelCommission`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `agent_channel_commission` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`channel_id` int unsigned NOT NULL COMMENT 'channel id',
`ctime` bigint NOT NULL COMMENT '',
`mtime` bigint NOT NULL COMMENT '',
`integer` int NOT NULL COMMENT 'commission integer part',
`decimals` int NOT NULL COMMENT 'commission decimals part',
`status` tinyint unsigned NOT NULL COMMENT 'commission status',
PRIMARY KEY(`uid`,`channel_id`,`ctime`,`mtime`)
) ENGINE=InnoDB COMMENT='agent channel commission table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.AgentPath`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `agent_path` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`level` int unsigned NOT NULL COMMENT 'agent level',
`sn` int unsigned NOT NULL COMMENT 'agent serial number on specific level',
`path` text NOT NULL COMMENT 'user agent path',
`ctime` bigint NOT NULL COMMENT '',
`mtime` bigint NOT NULL COMMENT '',
`invite_code` varchar(100) NOT NULL COMMENT 'invite code',
`whitelist_id` int unsigned NOT NULL COMMENT 'agent commission whitelist id',
`invite_num` int unsigned NOT NULL COMMENT 'invite number',
`parent_uid` bigint unsigned NOT NULL COMMENT 'parent uid',
`dividend_position` int unsigned NOT NULL DEFAULT 0 COMMENT 'month dividend position',
PRIMARY KEY(`uid`)
) ENGINE=InnoDB COMMENT='agent path table' DEFAULT CHARSET=utf8;
CREATE INDEX `agent_path_level` ON `agent_path` (`level`);
CREATE INDEX `agent_path_sn` ON `agent_path` (`sn`);
CREATE INDEX `agent_path_path` ON `agent_path` (`path`(1024));
CREATE INDEX `agent_path_whitelist_id` ON `agent_path` (`whitelist_id`);
CREATE UNIQUE INDEX `agent_path_invite_code` ON `agent_path` (`invite_code`);

-- --------------------------------------------------
--  Table Structure for `models.AgentWithdraw`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `agent_withdraw` (
`id` bigint unsigned NOT NULL COMMENT 'agent withdraw id',
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`amount` bigint NOT NULL COMMENT 'amount',
`status` tinyint unsigned NOT NULL COMMENT 'status',
`ctime` bigint NOT NULL COMMENT '',
`mtime` bigint NOT NULL COMMENT '',
`desc` varchar(256) NULL COMMENT '',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='agent withdraw table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.InviteCode`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `invite_code` (
`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`code` varchar(16) NOT NULL COMMENT 'code',
`status` tinyint unsigned NOT NULL COMMENT 'status',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='invite code table' DEFAULT CHARSET=utf8;
CREATE UNIQUE INDEX `invite_code_code` ON `invite_code` (`code`);

