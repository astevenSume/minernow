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
`wechat` varchar(100) NOT NULL COMMENT 'Wechat',
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

