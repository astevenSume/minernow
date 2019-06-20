
-- --------------------------------------------------
--  Table Structure for `models.Token`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `token` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`client_type` int unsigned NOT NULL COMMENT 'app type',
`mtime` bigint NOT NULL COMMENT 'access_token modify time',
`access_token` varchar(256) NOT NULL COMMENT 'access_token',
`mac` varchar(100) NOT NULL COMMENT '',
PRIMARY KEY(`uid`,`client_type`,`mac`)
) ENGINE=InnoDB COMMENT='access token table' DEFAULT CHARSET=utf8;

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
--  Table Structure for `models.GameRiskWithdraw`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `game_risk_withdraw` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`eusd_num` bigint NOT NULL COMMENT 'withdraw euse num',
`alert_time` bigint unsigned NOT NULL COMMENT 'alert risk time',
`do_get` tinyint unsigned NOT NULL COMMENT 'weather get risk',
`warn_grade` tinyint unsigned NOT NULL COMMENT 'warn grade',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='game withdraw risk alert table' DEFAULT CHARSET=utf8;

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

-- --------------------------------------------------
--  Table Structure for `models.Appeal`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `appeal` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`type` tinyint NOT NULL COMMENT 'type',
`user_id` bigint unsigned NOT NULL COMMENT 'user_id',
`admin_id` int unsigned NOT NULL COMMENT 'admin_id',
`order_id` bigint unsigned NOT NULL COMMENT 'order_id',
`context` varchar(256) NOT NULL COMMENT 'context',
`status` tinyint NOT NULL COMMENT 'status',
`ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time',
`utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='appeal table' DEFAULT CHARSET=utf8;
CREATE INDEX `appeal_order_id_type_status` ON `appeal` (`order_id`, `type`, `status`);

-- --------------------------------------------------
--  Table Structure for `models.AppealDealLog`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `appeal_deal_log` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`appeal_id` bigint unsigned NOT NULL COMMENT 'appeal_id',
`admin_id` int unsigned NOT NULL COMMENT 'admin_id',
`order_id` bigint unsigned NOT NULL COMMENT 'order_id',
`action` tinyint NOT NULL COMMENT 'action',
`ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='appeal_deal_log table' DEFAULT CHARSET=utf8;
CREATE INDEX `appeal_deal_log_appeal_id_order_id_admin_id` ON `appeal_deal_log` (`appeal_id`, `order_id`, `admin_id`);

-- --------------------------------------------------
--  Table Structure for `models.CommissionCalc`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `commission_calc` (
`id` bigint unsigned NOT NULL COMMENT 'id',
`start` varchar(16) NOT NULL COMMENT 'commission start',
`end` varchar(16) NOT NULL COMMENT 'commission end',
`calc_start` varchar(128) NOT NULL COMMENT 'calc start',
`calc_end` varchar(128) NOT NULL COMMENT 'calc end',
`status` tinyint unsigned NOT NULL COMMENT 'status',
`desc` varchar(256) NOT NULL COMMENT 'desc',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='commission log table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.CommissionDistribute`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `commission_distribute` (
`id` bigint unsigned NOT NULL COMMENT 'id',
`start` varchar(16) NOT NULL COMMENT 'commission start',
`end` varchar(16) NOT NULL COMMENT 'commission end',
`distribute_start` varchar(128) NOT NULL COMMENT 'distribute start',
`distribute_end` varchar(128) NOT NULL COMMENT 'distribute end',
`status` tinyint unsigned NOT NULL COMMENT 'status',
`desc` varchar(256) NOT NULL COMMENT 'desc',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='commission distribute table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.CommissionStat`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `commission_stat` (
`ctime` bigint NOT NULL COMMENT 'ctime',
`tax_integer` int NOT NULL COMMENT 'tax integer part',
`tax_decimals` int NOT NULL COMMENT 'tax decimals part',
`channel_integer` int NOT NULL COMMENT 'channel integer part',
`channel_decimals` int NOT NULL COMMENT 'channel decimals part',
`commission_integer` int NOT NULL COMMENT 'commission integer part',
`commission_decimals` int NOT NULL COMMENT 'commission decimals part',
`profit_integer` int NOT NULL COMMENT 'profit integer part',
`profit_decimals` int NOT NULL COMMENT 'profit decimals part',
`mtime` bigint NOT NULL COMMENT 'modified time',
`status` tinyint unsigned NOT NULL COMMENT 'status',
PRIMARY KEY(`ctime`)
) ENGINE=InnoDB COMMENT='commission statistic daily table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.EosOtcReport`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `eos_otc_report` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'eos otc report id',
`uid` bigint unsigned NOT NULL COMMENT 'uid',
`total_order_num` bigint NOT NULL DEFAULT 0 COMMENT 'total_order_num',
`success_order_num` bigint NOT NULL DEFAULT 0 COMMENT 'success_order_num',
`fail_order_num` bigint NOT NULL DEFAULT 0 COMMENT 'fail_order_num',
`buy_eusd_num` bigint NOT NULL DEFAULT 0 COMMENT 'buy_eusd_num',
`sell_eusd_num` bigint NOT NULL DEFAULT 0 COMMENT 'sell_eusd_num',
`date` int NOT NULL DEFAULT 0 COMMENT 'date',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='eos otc report' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.MarketPrices`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `market_prices` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'market otc price record id',
`market` tinyint NOT NULL COMMENT 'market',
`currency` varchar(100) NOT NULL COMMENT 'currency',
`trade_method` tinyint NOT NULL COMMENT 'trade method',
`pow_price` bigint unsigned NOT NULL COMMENT 'pow price',
`pow` int NOT NULL DEFAULT 4 COMMENT 'pow_price = (real price) * pow',
`ctime` bigint NOT NULL COMMENT 'created time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='market otc ex price' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.OtcBuy`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `otc_buy` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`available` bigint NOT NULL COMMENT 'available token',
`frozen` bigint unsigned NOT NULL COMMENT 'frozen token',
`bought` bigint unsigned NOT NULL COMMENT 'bought token',
`lower_limit_wechat` bigint NOT NULL COMMENT 'lower limit',
`upper_limit_wechat` bigint NOT NULL COMMENT 'upper limit',
`lower_limit_bank` bigint NOT NULL COMMENT 'lower limit',
`upper_limit_bank` bigint NOT NULL COMMENT 'upper limit',
`lower_limit_ali` bigint NOT NULL COMMENT 'lower limit',
`upper_limit_ali` bigint NOT NULL COMMENT 'upper limit',
`pay_type` tinyint unsigned NOT NULL COMMENT 'pay type',
`ctime` bigint NOT NULL COMMENT '',
PRIMARY KEY(`uid`)
) ENGINE=InnoDB COMMENT='otc_buy' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.OtcExchanger`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `otc_exchanger` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`mobile` varchar(100) NOT NULL COMMENT 'Mobile',
`wechat` varchar(100) NOT NULL COMMENT 'Wechat',
`telegram` varchar(100) NOT NULL COMMENT 'telegram',
`from` tinyint NOT NULL COMMENT 'from',
`ctime` bigint NOT NULL COMMENT '',
`utime` bigint NOT NULL COMMENT '',
PRIMARY KEY(`uid`)
) ENGINE=InnoDB COMMENT='otc exchanger' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.OtcExchangerVerify`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `otc_exchanger_verify` (
`id` int NOT NULL AUTO_INCREMENT COMMENT 'id',
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`mobile` varchar(100) NOT NULL COMMENT 'Mobile',
`wechat` varchar(100) NOT NULL COMMENT 'Wechat',
`telegram` varchar(100) NOT NULL COMMENT 'telegram',
`status` tinyint NOT NULL COMMENT 'status',
`from` tinyint NOT NULL COMMENT 'from',
`ctime` bigint NOT NULL COMMENT '',
`utime` bigint NOT NULL COMMENT '',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='otc exchanger verify' DEFAULT CHARSET=utf8;
CREATE INDEX `otc_exchanger_verify_uid_mobile` ON `otc_exchanger_verify` (`uid`, `mobile`);

-- --------------------------------------------------
--  Table Structure for `models.OtcMsg`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `otc_msg` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`order_id` bigint NOT NULL COMMENT 'Order id',
`uid` bigint unsigned NOT NULL COMMENT 'Msg Uid',
`content` text NOT NULL COMMENT 'Msg Content',
`is_read` tinyint unsigned NOT NULL COMMENT 'is_read 1--have read 0--unread',
`msg_type` varchar(200) NOT NULL COMMENT 'system--system, text--uid',
`ctime` bigint NOT NULL COMMENT '',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='Otc Message' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.OtcOrder`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `otc_order` (
`id` bigint unsigned NOT NULL COMMENT 'id',
`uid` bigint unsigned NOT NULL COMMENT '',
`uip` varchar(40) NOT NULL DEFAULT '' COMMENT 'user ip',
`euid` bigint unsigned NOT NULL DEFAULT 0 COMMENT '',
`eip` varchar(40) NOT NULL COMMENT ' ip',
`side` tinyint NOT NULL COMMENT '1-buy 2-sell',
`amount` bigint NOT NULL COMMENT 'eusd num',
`price` varchar(100) NOT NULL COMMENT 'eusd => rmb price',
`funds` bigint NOT NULL COMMENT 'rmb funds',
`fee` bigint NOT NULL COMMENT 'rmb price',
`pay_id` bigint unsigned NOT NULL COMMENT 'pay id',
`pay_type` tinyint NOT NULL COMMENT 'pay type',
`pay_name` varchar(128) NOT NULL COMMENT '',
`pay_account` varchar(300) NOT NULL COMMENT '',
`pay_bank` varchar(128) NOT NULL COMMENT '',
`pay_bank_branch` varchar(300) NOT NULL COMMENT '',
`transfer_id` bigint unsigned NOT NULL COMMENT '',
`ctime` bigint NOT NULL COMMENT '',
`pay_time` bigint NOT NULL COMMENT '',
`finish_time` bigint NOT NULL COMMENT '',
`utime` bigint NOT NULL COMMENT '',
`status` tinyint NOT NULL COMMENT '',
`epay_id` bigint unsigned NOT NULL COMMENT 'pay id',
`epay_type` tinyint NOT NULL COMMENT 'pay type',
`epay_name` varchar(128) NOT NULL COMMENT '',
`epay_account` varchar(300) NOT NULL COMMENT '',
`epay_bank` varchar(128) NOT NULL COMMENT '',
`epay_bank_branch` varchar(300) NOT NULL COMMENT '',
`appeal_status` tinyint NOT NULL DEFAULT 0 COMMENT '',
`admin_id` int unsigned NOT NULL DEFAULT 0 COMMENT '',
`qr_code` varchar(300) NOT NULL COMMENT '',
`date` int NOT NULL DEFAULT 0 COMMENT '',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='otc_sell' DEFAULT CHARSET=utf8;
CREATE INDEX `otc_order_uid_side_status` ON `otc_order` (`uid`, `side`, `status`);
CREATE INDEX `otc_order_euid_side_pay_type` ON `otc_order` (`euid`, `side`, `pay_type`);
CREATE INDEX `otc_order_euid_side_date` ON `otc_order` (`euid`, `side`, `date`);

-- --------------------------------------------------
--  Table Structure for `models.OtcSell`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `otc_sell` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`available` bigint NOT NULL COMMENT 'available token',
`frozen` bigint NOT NULL COMMENT 'frozen token',
`sold` bigint NOT NULL COMMENT 'Sold token',
`lower_limit` bigint NOT NULL COMMENT 'lower limit',
`upper_limit` bigint NOT NULL COMMENT 'upper limit',
`pay_type` tinyint unsigned NOT NULL COMMENT 'pay type',
`ctime` bigint NOT NULL COMMENT '',
PRIMARY KEY(`uid`)
) ENGINE=InnoDB COMMENT='otc_sell' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.PaymentMethod`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `payment_method` (
`pmid` bigint unsigned NOT NULL COMMENT 'payment method id',
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`mtype` tinyint unsigned NOT NULL COMMENT 'payment method enum',
`ord` int unsigned NOT NULL COMMENT 'order serial number',
`name` varchar(128) NOT NULL COMMENT 'name',
`account` varchar(128) NOT NULL COMMENT 'account',
`status` tinyint unsigned NOT NULL COMMENT 'status 0-disable 1-enable ',
`ctime` int unsigned NOT NULL COMMENT 'created time',
`bank` varchar(128) NOT NULL COMMENT '',
`bank_branch` varchar(128) NOT NULL COMMENT '',
`qr_code` varchar(256) NOT NULL COMMENT '',
`qr_code_content` text NOT NULL COMMENT 'qr_code_content',
`low_money_per_tx_limit` bigint NOT NULL COMMENT '',
`high_money_per_tx_limit` bigint NOT NULL COMMENT '',
`times_per_day_limit` bigint NOT NULL COMMENT '',
`money_per_day_limit` bigint NOT NULL COMMENT '',
`money_sum_limit` bigint NOT NULL COMMENT '',
`times_today` bigint NOT NULL COMMENT '',
`money_today` bigint NOT NULL COMMENT '',
`money_sum` bigint unsigned NOT NULL COMMENT '',
`mtime` int unsigned NOT NULL COMMENT 'modified time',
`use_time` int unsigned NOT NULL COMMENT 'use time',
PRIMARY KEY(`pmid`)
) ENGINE=InnoDB COMMENT='user payment methods table' DEFAULT CHARSET=utf8;
CREATE INDEX `payment_method_uid_mtype` ON `payment_method` (`uid`, `mtype`);

-- --------------------------------------------------
--  Table Structure for `models.Prices`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `prices` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'otc price record id',
`currency` varchar(100) NOT NULL COMMENT 'currency',
`pow_price` bigint unsigned NOT NULL COMMENT 'pow price',
`pow` int NOT NULL DEFAULT 4 COMMENT 'pow_price = (real price) * pow',
`ctime` bigint NOT NULL COMMENT 'created time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='system otc ex price' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.SystemNotification`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `system_notification` (
`nid` bigint unsigned NOT NULL COMMENT 'notification id',
`notification_type` varchar(100) NOT NULL COMMENT 'Notification Type',
`content` text NOT NULL COMMENT 'content',
`uid` bigint unsigned NOT NULL COMMENT 'uid',
`is_read` int NOT NULL COMMENT 'is_read',
`ctime` bigint NOT NULL COMMENT '',
PRIMARY KEY(`nid`)
) ENGINE=InnoDB COMMENT='System Notification' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.User`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `user` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`national_code` varchar(100) NOT NULL COMMENT 'national_code',
`mobile` varchar(100) NOT NULL COMMENT 'Mobile',
`status` tinyint NOT NULL COMMENT 'status',
`nick` varchar(100) NOT NULL COMMENT 'nick',
`pass` varchar(100) NOT NULL COMMENT 'password',
`salt` varchar(16) NOT NULL COMMENT 'salt',
`ctime` bigint NOT NULL COMMENT '',
`utime` bigint NOT NULL COMMENT '',
`ip` varchar(100) NOT NULL DEFAULT '' COMMENT 'ip',
`last_login_time` bigint NOT NULL DEFAULT 0 COMMENT 'last login time',
`last_login_ip` varchar(100) NOT NULL DEFAULT '' COMMENT 'last login ip',
`is_exchanger` tinyint NOT NULL DEFAULT 0 COMMENT 'last login ip',
`sign_salt` varchar(256) NOT NULL COMMENT 'sign salt',
PRIMARY KEY(`uid`)
) ENGINE=InnoDB COMMENT='user table' DEFAULT CHARSET=utf8;
CREATE UNIQUE INDEX `user_national_code_mobile` ON `user` (`national_code`, `mobile`);

-- --------------------------------------------------
--  Table Structure for `models.UserConfig`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `user_config` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`wealth_notice` bool NOT NULL COMMENT 'wealth_notice',
`order_notice` bool NOT NULL COMMENT 'order_notice',
PRIMARY KEY(`uid`)
) ENGINE=InnoDB COMMENT='User Config' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.UserLoginLog`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `user_login_log` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`user_id` bigint unsigned NOT NULL COMMENT 'user_id',
`user_agent` varchar(256) NULL COMMENT '',
`ips` varchar(256) NULL COMMENT 'ips',
`ctime` bigint NOT NULL COMMENT '',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='user login log table' DEFAULT CHARSET=utf8;
CREATE INDEX `user_login_log_ctime` ON `user_login_log` (`ctime`);

-- --------------------------------------------------
--  Table Structure for `models.UserPayPassword`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `user_pay_pass` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`pass` varchar(256) NOT NULL COMMENT 'password',
`salt` varchar(256) NOT NULL COMMENT 'salt',
`sign_salt` varchar(256) NOT NULL COMMENT 'sign salt',
`status` tinyint NOT NULL COMMENT 'status',
`method` tinyint NOT NULL COMMENT 'setting pwd method',
`verify_step` tinyint NOT NULL COMMENT '0 pass verify 1 pass&sms verify',
`timestamp` bigint NOT NULL DEFAULT 0 COMMENT 'timestamp',
PRIMARY KEY(`uid`)
) ENGINE=InnoDB COMMENT='User Pay Password' DEFAULT CHARSET=utf8;
CREATE UNIQUE INDEX `user_pay_pass_uid` ON `user_pay_pass` (`uid`);

-- --------------------------------------------------
--  Table Structure for `models.UsdtAccount`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `usdt_account` (
`uaid` bigint unsigned NOT NULL COMMENT '',
`uid` bigint unsigned NOT NULL COMMENT '',
`status` tinyint unsigned NOT NULL COMMENT 'status 0-locked 1-working',
`available_integer` bigint NOT NULL COMMENT 'available amount integer part',
`frozen_integer` bigint NOT NULL COMMENT 'frozen amount integer part',
`transfer_frozen_integer` bigint NOT NULL COMMENT 'frozen amount integer part',
`mortgaged_integer` bigint NOT NULL COMMENT 'mortgaged amount integer part',
`btc_available_integer` bigint NOT NULL COMMENT 'available btc amount integer part',
`btc_frozen_integer` bigint NOT NULL COMMENT 'frozen btc amount integer part',
`btc_mortgaged_integer` bigint NOT NULL COMMENT 'mortgaged btc amount integer part',
`waiting_cash_sweep_integer` bigint NOT NULL COMMENT 'waiting cash sweep integer part',
`cash_sweep_integer` bigint NOT NULL COMMENT 'cash sweep integer part',
`owned_by_platform_integer` bigint NOT NULL COMMENT 'owned by platform integer part',
`sweep_status` tinyint unsigned NOT NULL COMMENT 'sweep status',
`pkid` bigint unsigned NOT NULL COMMENT 'usdt private key id',
`address` varchar(100) NOT NULL COMMENT '',
`ctime` bigint NOT NULL COMMENT 'created time',
`mtime` bigint NOT NULL COMMENT 'update time',
PRIMARY KEY(`uaid`)
) ENGINE=InnoDB COMMENT='usdt account table' DEFAULT CHARSET=utf8;
CREATE INDEX `usdt_account_uid` ON `usdt_account` (`uid`);
CREATE INDEX `usdt_account_pkid` ON `usdt_account` (`pkid`);

-- --------------------------------------------------
--  Table Structure for `models.UsdtOnchainBalance`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `usdt_onchain_balance` (
`address` varchar(100) NOT NULL COMMENT '',
`property_id` int unsigned NOT NULL COMMENT '',
`pending_pos` varchar(100) NOT NULL COMMENT '',
`reserved` varchar(100) NOT NULL COMMENT '',
`divisible` bool NOT NULL COMMENT '',
`amount_integer` bigint NOT NULL COMMENT 'amount integer part',
`frozen` varchar(100) NOT NULL COMMENT '',
`pending_neg` varchar(100) NOT NULL COMMENT '',
`mtime` bigint NOT NULL COMMENT '',
PRIMARY KEY(`address`,`property_id`)
) ENGINE=InnoDB COMMENT='the balance on chain' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.UsdtOnChainData`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `usdt_onchain_data` (
`address` varchar(100) NOT NULL COMMENT '',
`attr_type` int unsigned NOT NULL COMMENT '',
`data_int64` bigint NOT NULL COMMENT '',
`data_str` varchar(256) NOT NULL COMMENT '',
PRIMARY KEY(`address`,`attr_type`)
) ENGINE=InnoDB COMMENT='configurations of usdt on chain.' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.UsdtOnchainLog`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `usdt_onchain_log` (
`oclid` bigint unsigned NOT NULL COMMENT '',
`from` varchar(100) NOT NULL COMMENT '',
`to` varchar(100) NOT NULL COMMENT '',
`tx` varchar(100) NOT NULL COMMENT '',
`status` varchar(100) NOT NULL COMMENT '',
`pushed` varchar(100) NOT NULL COMMENT '',
`signedTx` text NOT NULL COMMENT '',
`amount_integer` bigint NOT NULL COMMENT 'amount integer part',
`ctime` bigint NOT NULL COMMENT 'created time',
PRIMARY KEY(`oclid`)
) ENGINE=InnoDB COMMENT='the log of tx on chain' DEFAULT CHARSET=utf8;
CREATE INDEX `usdt_onchain_log_from` ON `usdt_onchain_log` (`from`);
CREATE INDEX `usdt_onchain_log_to` ON `usdt_onchain_log` (`to`);
CREATE INDEX `usdt_onchain_log_ctime` ON `usdt_onchain_log` (`ctime`);

-- --------------------------------------------------
--  Table Structure for `models.UsdtOnChainSyncPos`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `usdt_onchain_sync_pos` (
`address` varchar(100) NOT NULL COMMENT '',
`page` int unsigned NOT NULL COMMENT '',
`tx_id` varchar(128) NOT NULL COMMENT '',
PRIMARY KEY(`address`)
) ENGINE=InnoDB COMMENT='configurations of usdt on chain.' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.UsdtOnchainTransaction`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `usdt_onchain_transaction` (
`tx_id` varchar(100) NOT NULL COMMENT '',
`uaid` bigint unsigned NOT NULL COMMENT 'usdt account id',
`type` tinyint unsigned NOT NULL COMMENT 'Transaction type',
`property_id` int unsigned NOT NULL COMMENT '',
`property_name` varchar(100) NOT NULL COMMENT '',
`tx_type` varchar(100) NOT NULL COMMENT '',
`tx_type_int` int NOT NULL COMMENT '',
`amount_integer` bigint NOT NULL COMMENT 'amount integer part',
`block` int unsigned NOT NULL COMMENT '',
`block_hash` varchar(100) NOT NULL COMMENT '',
`block_time` bigint NOT NULL COMMENT '',
`confirmations` int unsigned NOT NULL COMMENT '',
`divisible` bool NOT NULL COMMENT '',
`fee_amount_integer` bigint NOT NULL COMMENT 'fee amount integer part',
`is_mine` bool NOT NULL COMMENT '',
`position_in_block` int unsigned NOT NULL COMMENT '',
`referenceaddress` varchar(100) NOT NULL COMMENT '',
`sending_address` varchar(100) NOT NULL COMMENT '',
`version` int NOT NULL COMMENT '',
`mtime` bigint NOT NULL COMMENT '',
PRIMARY KEY(`tx_id`,`property_id`)
) ENGINE=InnoDB COMMENT='the tx on chain' DEFAULT CHARSET=utf8;
CREATE INDEX `usdt_onchain_transaction_uaid_type` ON `usdt_onchain_transaction` (`uaid`, `type`);

-- --------------------------------------------------
--  Table Structure for `models.PriKey`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `usdt_prikey` (
`pkid` bigint unsigned NOT NULL COMMENT '',
`pri` varchar(256) NOT NULL COMMENT '',
`address` varchar(100) NOT NULL COMMENT '',
PRIMARY KEY(`pkid`)
) ENGINE=InnoDB COMMENT='' DEFAULT CHARSET=utf8;
CREATE UNIQUE INDEX `usdt_prikey_pri` ON `usdt_prikey` (`pri`);
CREATE UNIQUE INDEX `usdt_prikey_address` ON `usdt_prikey` (`address`);

-- --------------------------------------------------
--  Table Structure for `models.UsdtSweepLog`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `usdt_sweep_log` (
`id` bigint unsigned NOT NULL COMMENT 'id',
`uid` bigint unsigned NOT NULL COMMENT 'uid',
`ttype` int unsigned NOT NULL COMMENT 'Transaction type',
`status` int unsigned NOT NULL COMMENT 'Transaction status',
`from` varchar(256) NOT NULL COMMENT 'sender address',
`to` varchar(256) NOT NULL COMMENT 'receiver address',
`txid` varchar(256) NOT NULL COMMENT 'Transaction id',
`amount_integer` bigint NOT NULL COMMENT 'amount integer',
`fee_integer` bigint NOT NULL COMMENT 'fee integer',
`fee_onchain_integer` bigint NOT NULL COMMENT 'fee onchain',
`ctime` bigint NOT NULL COMMENT 'created time',
`utime` bigint NOT NULL COMMENT 'update time',
`step` varchar(64) NOT NULL COMMENT 'desc',
`desc` varchar(256) NOT NULL COMMENT 'desc',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='Usdt sweep log' DEFAULT CHARSET=utf8;
CREATE INDEX `usdt_sweep_log_uid_ttype` ON `usdt_sweep_log` (`uid`, `ttype`);
CREATE INDEX `usdt_sweep_log_txid` ON `usdt_sweep_log` (`txid`);

-- --------------------------------------------------
--  Table Structure for `models.UsdtTransaction`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `usdt_transaction` (
`tx_id` varchar(100) NOT NULL COMMENT '',
`uaid` bigint unsigned NOT NULL COMMENT 'usdt account id',
`type` tinyint unsigned NOT NULL COMMENT 'Transaction type',
`block_num` int unsigned NOT NULL COMMENT 'block num',
`status` int NOT NULL COMMENT 'Transaction status',
`payer` varchar(100) NOT NULL COMMENT 'Payer Account',
`receiver` varchar(100) NOT NULL COMMENT 'Receiver Account',
`amount_integer` bigint NOT NULL COMMENT 'amount integer part',
`fee` varchar(100) NOT NULL COMMENT 'transaction Fee',
`memo` varchar(100) NOT NULL COMMENT 'transaction memo',
`ctime` bigint NOT NULL COMMENT 'created time',
`utime` bigint NOT NULL COMMENT 'update time',
PRIMARY KEY(`tx_id`,`uaid`)
) ENGINE=InnoDB COMMENT='Usdt Transaction' DEFAULT CHARSET=utf8;
CREATE INDEX `usdt_transaction_uaid_type` ON `usdt_transaction` (`uaid`, `type`);

-- --------------------------------------------------
--  Table Structure for `models.UsdtWealthLog`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `usdt_wealth_log` (
`id` bigint unsigned NOT NULL COMMENT 'id',
`uid` bigint unsigned NOT NULL COMMENT 'uid',
`ttype` int unsigned NOT NULL COMMENT 'Transaction type',
`status` int unsigned NOT NULL COMMENT 'Transaction status',
`from` varchar(256) NOT NULL COMMENT 'sender address',
`to` varchar(256) NOT NULL COMMENT 'receiver address',
`txid` varchar(256) NOT NULL COMMENT 'Transaction id',
`amount_integer` bigint NOT NULL COMMENT 'amount integer',
`fee_integer` bigint NOT NULL COMMENT 'fee integer',
`fee_usdt_integer` bigint NOT NULL COMMENT 'fee usdt integer',
`fee_onchain_integer` bigint NOT NULL COMMENT 'fee onchain',
`ctime` bigint NOT NULL COMMENT 'created time',
`utime` bigint NOT NULL COMMENT 'update time',
`step` varchar(64) NOT NULL COMMENT 'desc',
`desc` varchar(256) NOT NULL COMMENT 'desc',
`sign` varchar(256) NOT NULL COMMENT 'sign',
`memo` varchar(256) NOT NULL COMMENT 'memo',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='Usdt wealth log' DEFAULT CHARSET=utf8;
CREATE INDEX `usdt_wealth_log_uid_ttype` ON `usdt_wealth_log` (`uid`, `ttype`);
CREATE INDEX `usdt_wealth_log_txid` ON `usdt_wealth_log` (`txid`);
CREATE INDEX `usdt_wealth_log_ctime` ON `usdt_wealth_log` (`ctime`);

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
`dividend_position` int NOT NULL DEFAULT 0 COMMENT 'month dividend position',
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

-- --------------------------------------------------
--  Table Structure for `models.GameUserMonthReport`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `game_user_month_report` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`ctime` bigint NOT NULL COMMENT 'create time',
`profit` bigint NOT NULL COMMENT '自己的盈亏金额',
`agents_profit` bigint NOT NULL COMMENT '无限下级下级代理的盈亏金额',
`result_profit` bigint NOT NULL COMMENT '最终的盈亏金额',
`bet_amount` bigint NOT NULL COMMENT '投注额',
`effective_bet_amount` bigint NOT NULL COMMENT '累计有效投注额(本人加无限下级)',
`play_game_day` int NOT NULL COMMENT '玩游戏的天数',
`is_activity_user` bool NOT NULL DEFAULT false COMMENT '是否是活跃用户',
`agent_level` int unsigned NOT NULL COMMENT '代理等级',
`up_agent_uid` bigint unsigned NOT NULL COMMENT '上级代理的UID',
`activity_agent_num` int NOT NULL COMMENT '无限下级代理活跃人数',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='game_user_month_report table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.MonthDividendRecord`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `month_dividend_record` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`ctime` bigint NOT NULL COMMENT 'create time',
`self_dividend` bigint NOT NULL COMMENT '自己的分红',
`agent_dividend` bigint NOT NULL COMMENT '要分给代理的分红',
`result_dividend` bigint NOT NULL COMMENT '最终获得的分红 self_dividend-agent_dividend',
`receive_status` int NOT NULL DEFAULT false COMMENT '领取状态1是已领取2是未领取3是等待上级发放',
`received_time` bigint NOT NULL COMMENT '领取奖励的时间',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='month_dividend_record table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.ProfitReportDaily`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `profit_report_daily` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`channel_id` int unsigned NOT NULL COMMENT 'game channel id',
`bet` bigint NOT NULL COMMENT '本人的有效投注额',
`total_valid_bet` bigint NOT NULL COMMENT '累计有效投注额(本人加无限下级)',
`profit` bigint NOT NULL COMMENT '盈亏金额',
`salary` bigint NOT NULL COMMENT '日工资',
`self_dividend` bigint NOT NULL COMMENT '属于自己的月分红',
`agent_dividend` bigint NOT NULL COMMENT '分给下级代理的月分红',
`result_dividend` bigint NOT NULL COMMENT '最终获得的月分红,可能为负数',
`withdraw_amount` int unsigned NOT NULL COMMENT '佣金提现值',
`game_recharge_amount` int unsigned NOT NULL COMMENT '游戏充值金额',
`ctime` bigint NOT NULL COMMENT 'create time',
PRIMARY KEY(`uid`,`channel_id`,`ctime`)
) ENGINE=InnoDB COMMENT='profit_report_daily table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.ReportGameRecordAg`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `report_game_record_ag` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`account` varchar(50) NOT NULL COMMENT 'account',
`game_type` varchar(50) NOT NULL COMMENT '游戏ID',
`game_name` varchar(50) NOT NULL COMMENT '游戏名称',
`order_id` varchar(50) NOT NULL COMMENT '订单编号',
`table_id` varchar(50) NOT NULL COMMENT '桌号',
`bet` bigint NOT NULL COMMENT '投注额',
`valid_bet` bigint NOT NULL COMMENT '有效投注额',
`profit` bigint NOT NULL COMMENT '盈亏金额',
`bet_time` varchar(32) NOT NULL COMMENT '下注时间',
`ctime` bigint NOT NULL COMMENT 'ctime',
`uid` bigint unsigned NOT NULL COMMENT 'user id',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='report_game_record_Ag table' DEFAULT CHARSET=utf8;
CREATE INDEX `report_game_record_ag_uid_bet_valid_bet_profit_ctime` ON `report_game_record_ag` (`uid`, `bet`, `valid_bet`, `profit`, `ctime`);

-- --------------------------------------------------
--  Table Structure for `models.ReportGameRecordKy`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `report_game_record_ky` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`account` varchar(50) NOT NULL COMMENT 'account',
`game_id` varchar(50) NOT NULL COMMENT '游戏局号',
`game_name` varchar(50) NOT NULL COMMENT 'game_name',
`server_id` int NOT NULL COMMENT '房间ID',
`kind_id` varchar(50) NOT NULL COMMENT '游戏ID',
`table_id` int NOT NULL COMMENT '桌子号',
`chair_id` int NOT NULL COMMENT '椅子号',
`bet` bigint NOT NULL COMMENT '投注额',
`valid_bet` bigint NOT NULL COMMENT '有效投注额',
`profit` bigint NOT NULL COMMENT '盈亏金额',
`revenue` bigint NOT NULL COMMENT '抽水金额',
`start_time` varchar(32) NOT NULL COMMENT 'start_time',
`end_time` varchar(32) NOT NULL COMMENT 'end_time',
`ctime` bigint NOT NULL COMMENT 'create time',
`uid` bigint unsigned NOT NULL COMMENT 'user id',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='report_game_record_ky table' DEFAULT CHARSET=utf8;
CREATE INDEX `report_game_record_ky_uid_bet_valid_bet_profit_ctime` ON `report_game_record_ky` (`uid`, `bet`, `valid_bet`, `profit`, `ctime`);

-- --------------------------------------------------
--  Table Structure for `models.ReportGameRecordRg`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `report_game_record_rg` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`account` varchar(50) NOT NULL COMMENT 'account',
`game_name_id` varchar(50) NOT NULL COMMENT '游戏id',
`game_name` varchar(50) NOT NULL COMMENT '游戏名称',
`game_kind_name` varchar(50) NOT NULL COMMENT '玩法名称',
`order_id` varchar(50) NOT NULL COMMENT '订单编号',
`open_date` varchar(32) NOT NULL COMMENT '开奖时间',
`period_name` varchar(50) NOT NULL COMMENT '期号',
`open_number` varchar(50) NOT NULL COMMENT '开奖号码',
`status` tinyint unsigned NOT NULL COMMENT '订单状态',
`bet` bigint NOT NULL COMMENT '投注额',
`valid_bet` bigint NOT NULL COMMENT '有效投注额',
`profit` bigint NOT NULL COMMENT '盈亏金额',
`bet_time` varchar(32) NOT NULL COMMENT '下注时间',
`bet_content` varchar(50) NOT NULL COMMENT '下注内容',
`ctime` bigint NOT NULL COMMENT 'ctime',
`uid` bigint unsigned NOT NULL COMMENT 'user id',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='report_game_record_Rg table' DEFAULT CHARSET=utf8;
CREATE INDEX `report_game_record_rg_uid_bet_valid_bet_profit_ctime` ON `report_game_record_rg` (`uid`, `bet`, `valid_bet`, `profit`, `ctime`);

-- --------------------------------------------------
--  Table Structure for `models.ReportGameUserDaily`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `report_game_user_daily` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`channel_id` int unsigned NOT NULL COMMENT 'game channel id',
`bet` bigint NOT NULL COMMENT '投注金额',
`valid_bet` bigint NOT NULL COMMENT '有效投注额',
`total_valid_bet` bigint NOT NULL COMMENT '累计有效投注额(本人加无限下级)',
`profit` bigint NOT NULL COMMENT '盈亏金额',
`total_profit` bigint NOT NULL COMMENT '累计盈亏金额(本人加无限下级)',
`salary` bigint NOT NULL COMMENT '日工资',
`status` tinyint unsigned NOT NULL COMMENT 'status',
`ctime` bigint NOT NULL COMMENT 'create time',
PRIMARY KEY(`uid`,`channel_id`,`ctime`)
) ENGINE=InnoDB COMMENT='report_game_user_daily table' DEFAULT CHARSET=utf8;

