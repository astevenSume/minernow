-- --------------------------------------------------
--  Table Structure for `models.EosAccount`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `eos_account` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`uid` bigint unsigned NOT NULL COMMENT 'uid',
`account` varchar(100) NOT NULL COMMENT 'account',
`balance` varchar(100) NOT NULL COMMENT 'balance',
`status` tinyint NOT NULL COMMENT '',
`ctime` bigint NOT NULL COMMENT 'created time',
`utime` bigint NOT NULL COMMENT 'update time',
PRIMARY KEY(`id`,`uid`)
) ENGINE=InnoDB COMMENT='Eos Account' DEFAULT CHARSET=utf8;
CREATE INDEX `eos_account_uid` ON `eos_account` (`uid`);
CREATE INDEX `eos_account_account` ON `eos_account` (`account`);

-- --------------------------------------------------
--  Table Structure for `models.EosOtc`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `eos_otc` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`account` varchar(100) NOT NULL COMMENT 'account',
`status` tinyint NOT NULL COMMENT 'status',
`available` bigint NOT NULL COMMENT 'available',
`trade` bigint NOT NULL COMMENT 'otc trade balance',
`transfer` bigint NOT NULL COMMENT 'transfering balance',
`sell_state` varchar(200) NOT NULL DEFAULT '' COMMENT '',
`sell_pay_type` tinyint unsigned NOT NULL COMMENT '',
`sell_able` bool NOT NULL COMMENT '',
`sell_rmb_day` bigint NOT NULL COMMENT '',
`sell_rmb_today` bigint NOT NULL COMMENT '',
`sell_rmb_lower_limit` bigint NOT NULL COMMENT '',
`sell_utime` bigint NOT NULL COMMENT '',
`buy_able` bool NOT NULL COMMENT '',
`buy_rmb_day` bigint NOT NULL COMMENT '',
`buy_rmb_today` bigint NOT NULL COMMENT '',
`buy_rmb_lower_limit` bigint NOT NULL COMMENT '',
`buy_utime` bigint NOT NULL COMMENT '',
`buy_state` varchar(200) NOT NULL DEFAULT '' COMMENT '',
`ctime` bigint NOT NULL COMMENT 'created time',
`utime` bigint NOT NULL COMMENT 'update time',
PRIMARY KEY(`uid`)
) ENGINE=InnoDB COMMENT='Eos otc' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.EosTransaction`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `eos_transaction` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`type` tinyint unsigned NOT NULL COMMENT 'Transaction type',
`transaction_id` varchar(100) NOT NULL COMMENT 'Transaction id',
`block_num` int unsigned NOT NULL COMMENT 'block num',
`status` tinyint NOT NULL COMMENT 'status',
`payer` varchar(100) NOT NULL COMMENT 'Payer Account',
`receiver` varchar(100) NOT NULL COMMENT 'Receiver Account',
`quantity` varchar(100) NOT NULL COMMENT 'Token quantity',
`memo` varchar(100) NOT NULL COMMENT 'transaction memo',
`ctime` bigint NOT NULL COMMENT 'created time',
`utime` bigint NOT NULL COMMENT 'update time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='Eos Transaction' DEFAULT CHARSET=utf8;
CREATE INDEX `eos_transaction_payer` ON `eos_transaction` (`payer`);
CREATE INDEX `eos_transaction_receiver` ON `eos_transaction` (`receiver`);

-- --------------------------------------------------
--  Table Structure for `models.EosTransactionInfo`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `eos_transaction_info` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`transaction_id` varchar(100) NOT NULL COMMENT 'Transaction id',
`block_num` int unsigned NOT NULL COMMENT 'block num',
`ctime` bigint NOT NULL COMMENT 'created time',
`processed` text NOT NULL COMMENT 'processed info',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='Eos Transaction info' DEFAULT CHARSET=utf8;
CREATE INDEX `eos_transaction_info_transaction_id` ON `eos_transaction_info` (`transaction_id`);

-- --------------------------------------------------
--  Table Structure for `models.EosTxLog`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `eos_tx_log` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`from` varchar(100) NOT NULL COMMENT ' Account',
`from_uid` bigint unsigned NOT NULL COMMENT ' Account',
`to` varchar(100) NOT NULL COMMENT ' Account',
`to_uid` bigint unsigned NOT NULL COMMENT ' Account',
`quantity` bigint NOT NULL COMMENT '',
`status` tinyint NOT NULL COMMENT '',
`log_ids` varchar(100) NOT NULL COMMENT '',
`ctime` bigint NOT NULL COMMENT '',
`txid` bigint unsigned NOT NULL COMMENT 'Transaction id',
`order_id` bigint unsigned NOT NULL COMMENT '',
`utime` bigint NOT NULL COMMENT '',
`sign` varchar(100) NOT NULL DEFAULT '' COMMENT '',
`delay_deal` bool NOT NULL DEFAULT false COMMENT '',
`retry` int NOT NULL DEFAULT 0 COMMENT '',
`memo` varchar(100) NOT NULL DEFAULT '' COMMENT '',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='Eos tx log' DEFAULT CHARSET=utf8;
CREATE INDEX `eos_tx_log_status` ON `eos_tx_log` (`status`);

-- --------------------------------------------------
--  Table Structure for `models.EosUseLog`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `eos_use_log` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`type` tinyint unsigned NOT NULL COMMENT 'Transaction type',
`tid` bigint unsigned NOT NULL COMMENT 'eos_transaction id',
`status` tinyint NOT NULL COMMENT 'status',
`tid_recover` bigint unsigned NOT NULL COMMENT 'eos_transaction id',
`quantity_num` bigint unsigned NOT NULL COMMENT 'Eos Num * 10000',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='Eos use log' DEFAULT CHARSET=utf8;
CREATE INDEX `eos_use_log_status` ON `eos_use_log` (`status`);

-- --------------------------------------------------
--  Table Structure for `models.EosWealth`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `eos_wealth` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`status` tinyint NOT NULL COMMENT '',
`account` varchar(100) NOT NULL COMMENT 'account',
`balance` bigint NOT NULL COMMENT 'balance',
`available` bigint NOT NULL COMMENT 'available balance',
`game` bigint NOT NULL COMMENT 'game balance',
`trade` bigint NOT NULL COMMENT 'trade frozen balance',
`transfer` bigint NOT NULL COMMENT 'transfering balance',
`transfer_game` bigint NOT NULL COMMENT 'transfering to game balance',
`is_exchanger` tinyint NOT NULL COMMENT '',
`ctime` bigint NOT NULL COMMENT 'created time',
`utime` bigint NOT NULL COMMENT 'update time',
PRIMARY KEY(`uid`)
) ENGINE=InnoDB COMMENT='Eos wealth' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.EosWealthLog`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `eos_wealth_log` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`uid` bigint unsigned NOT NULL COMMENT 'uid',
`uid2` bigint unsigned NOT NULL COMMENT 'uid',
`ttype` tinyint unsigned NOT NULL COMMENT 'Transaction type',
`status` tinyint NOT NULL COMMENT 'Transaction status',
`txid` bigint unsigned NOT NULL COMMENT 'Transaction id',
`quantity` bigint NOT NULL COMMENT 'Token quantity',
`ctime` bigint NOT NULL COMMENT 'created time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='Eos wealth log' DEFAULT CHARSET=utf8;
CREATE INDEX `eos_wealth_log_uid_ttype` ON `eos_wealth_log` (`uid`, `ttype`);
CREATE INDEX `eos_wealth_log_status` ON `eos_wealth_log` (`status`);

-- --------------------------------------------------
--  Table Structure for `models.EusdRetire`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `eusd_retire` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`from` varchar(100) NOT NULL COMMENT ' Account',
`from_uid` bigint unsigned NOT NULL COMMENT ' Account',
`quantity` bigint NOT NULL COMMENT '',
`status` tinyint NOT NULL COMMENT '',
`ctime` bigint NOT NULL COMMENT '',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='Eusd Retire' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.PlatformUser`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `platform_user` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`pid` int NOT NULL COMMENT '',
`status` tinyint NOT NULL COMMENT '',
`ctime` int unsigned NOT NULL COMMENT '',
PRIMARY KEY(`uid`)
) ENGINE=InnoDB COMMENT='' DEFAULT CHARSET=utf8;
CREATE INDEX `platform_user_pid_uid` ON `platform_user` (`pid`, `uid`);

-- --------------------------------------------------
--  Table Structure for `models.PlatformUserCate`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `platform_user_cate` (
`id` int NOT NULL AUTO_INCREMENT COMMENT 'id',
`name` varchar(100) NOT NULL COMMENT '',
`dividend` int NOT NULL COMMENT '',
`ctime` int unsigned NOT NULL COMMENT '',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='' DEFAULT CHARSET=utf8;

