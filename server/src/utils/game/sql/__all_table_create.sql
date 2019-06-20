-- --------------------------------------------------
--  Table Structure for `models.ChannelDaily`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `channel_daily` (
`channel_id` int unsigned NOT NULL COMMENT 'channel id',
`ctime` bigint NOT NULL COMMENT 'create time',
`win_lose_money_integer` int NOT NULL COMMENT 'win lose money amount integer part',
`win_lose_money_decimals` int NOT NULL COMMENT 'win lose money amount decimals part',
`chips_integer` int NOT NULL COMMENT 'chips amount integer part',
`chips_decimals` int NOT NULL COMMENT 'chips amount decimals part',
`mtime` bigint NOT NULL COMMENT 'modified time',
PRIMARY KEY(`channel_id`,`ctime`)
) ENGINE=InnoDB COMMENT='channel daily table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.GameLog`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `game_log` (
`id` bigint unsigned NOT NULL COMMENT 'id',
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`channel_id` int unsigned NOT NULL COMMENT 'game channel id',
`account` varchar(50) NOT NULL COMMENT 'account',
`log_type` tinyint unsigned NOT NULL COMMENT 'log type',
`desc` varchar(512) NOT NULL COMMENT '',
`ctime` bigint NOT NULL COMMENT 'create time',
PRIMARY KEY(`id`,`ctime`)
) ENGINE=InnoDB COMMENT='game log table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.GameOrderRisk`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `game_order_risk` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`alert_id` bigint NOT NULL COMMENT 'group alert id',
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`amount` bigint NOT NULL COMMENT 'eusd nums',
`funds` bigint NOT NULL COMMENT 'rmb',
`pay_type` tinyint NOT NULL COMMENT 'pay type',
`pay_account` varchar(300) NOT NULL COMMENT 'pay account',
`order_time` bigint NOT NULL COMMENT 'order time',
`ctime` bigint NOT NULL COMMENT 'ctime',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='game order risk table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.GameRiskAlert`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `game_risk_alert` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`funds` bigint NOT NULL COMMENT 'rmb num',
`eusd_num` bigint NOT NULL COMMENT 'withdraw euse num',
`order_time` bigint unsigned NOT NULL COMMENT 'order time',
`alert_time` bigint unsigned NOT NULL COMMENT 'alert risk time',
`do_get` tinyint unsigned NOT NULL COMMENT 'weather get risk',
`warn_grade` tinyint unsigned NOT NULL COMMENT 'warn grade',
`risk_type` tinyint unsigned NOT NULL COMMENT 'risk type',
`order_risk_id` bigint unsigned NOT NULL COMMENT 'order risk id',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='game risk alert table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.GameTransfer`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `game_transfer` (
`id` bigint unsigned NOT NULL COMMENT 'id',
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`channel_id` int unsigned NOT NULL COMMENT 'game channel id',
`account` varchar(50) NOT NULL COMMENT 'account',
`transfer_type` int unsigned NOT NULL COMMENT 'transfer type',
`order` varchar(50) NOT NULL COMMENT 'Order',
`game_order` varchar(50) NOT NULL COMMENT 'Game Order',
`coin_integer` bigint NOT NULL COMMENT 'coin integer part',
`eusd_integer` bigint NOT NULL COMMENT 'eusd integer part',
`status` int unsigned NOT NULL COMMENT 'status',
`ctime` bigint NOT NULL COMMENT '',
`desc` varchar(512) NOT NULL COMMENT '',
`step` varchar(256) NOT NULL COMMENT '',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='game user table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.GameUser`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `game_user` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`channel_id` int unsigned NOT NULL COMMENT 'game channel id',
`account` varchar(50) NOT NULL COMMENT 'account',
`nick_name` varchar(50) NOT NULL COMMENT 'nick name',
`sex` tinyint unsigned NOT NULL COMMENT 'sex',
`password` varchar(100) NOT NULL COMMENT 'password',
`ctime` bigint NOT NULL COMMENT '',
`mtime` bigint NOT NULL COMMENT '',
`status` tinyint unsigned NOT NULL COMMENT '',
PRIMARY KEY(`uid`,`channel_id`)
) ENGINE=InnoDB COMMENT='game user table' DEFAULT CHARSET=utf8;
CREATE INDEX `game_user_channel_id_account` ON `game_user` (`channel_id`, `account`);

-- --------------------------------------------------
--  Table Structure for `models.GameUserDaily`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `game_user_daily` (
`channel_id` int unsigned NOT NULL COMMENT 'game channel id',
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`tax_integer` int NOT NULL COMMENT 'tax amount integer part',
`tax_decimals` int NOT NULL COMMENT 'tax amount decimals part',
`chips_integer` int NOT NULL COMMENT 'chips amount integer part',
`chips_decimals` int NOT NULL COMMENT 'chips amount decimals part',
`winlose_integer` int NOT NULL COMMENT 'winlose amount integer part',
`winlose_decimals` int NOT NULL COMMENT 'winlose amount decimals part',
`ctime` bigint NOT NULL COMMENT 'created time, equals to the begin time of the day',
`mtime` bigint NOT NULL COMMENT 'modified time',
PRIMARY KEY(`channel_id`,`uid`,`ctime`)
) ENGINE=InnoDB COMMENT='game user daily table' DEFAULT CHARSET=utf8;

