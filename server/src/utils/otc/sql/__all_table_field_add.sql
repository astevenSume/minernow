----------------------------------------------------
--  `appeal`
----------------------------------------------------
ALTER TABLE `appeal` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `appeal` ADD `type` tinyint NOT NULL COMMENT 'type' AFTER `id`;
ALTER TABLE `appeal` ADD `user_id` bigint unsigned NOT NULL COMMENT 'user_id' AFTER `type`;
ALTER TABLE `appeal` ADD `admin_id` int unsigned NOT NULL COMMENT 'admin_id' AFTER `user_id`;
ALTER TABLE `appeal` ADD `order_id` bigint unsigned NOT NULL COMMENT 'order_id' AFTER `admin_id`;
ALTER TABLE `appeal` ADD `context` varchar(256) NOT NULL COMMENT 'context' AFTER `order_id`;
ALTER TABLE `appeal` ADD `status` tinyint NOT NULL COMMENT 'status' AFTER `context`;
ALTER TABLE `appeal` ADD `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time' AFTER `status`;
ALTER TABLE `appeal` ADD `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time' AFTER `ctime`;
ALTER TABLE `appeal` ADD `wechat` varchar(100) NOT NULL COMMENT 'Wechat' AFTER `utime`;

----------------------------------------------------
--  `appeal_deal_log`
----------------------------------------------------
ALTER TABLE `appeal_deal_log` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `appeal_deal_log` ADD `appeal_id` bigint unsigned NOT NULL COMMENT 'appeal_id' AFTER `id`;
ALTER TABLE `appeal_deal_log` ADD `admin_id` int unsigned NOT NULL COMMENT 'admin_id' AFTER `appeal_id`;
ALTER TABLE `appeal_deal_log` ADD `order_id` bigint unsigned NOT NULL COMMENT 'order_id' AFTER `admin_id`;
ALTER TABLE `appeal_deal_log` ADD `action` tinyint NOT NULL COMMENT 'action' AFTER `order_id`;
ALTER TABLE `appeal_deal_log` ADD `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time' AFTER `action`;

----------------------------------------------------
--  `commission_calc`
----------------------------------------------------
ALTER TABLE `commission_calc` ADD `id` bigint unsigned NOT NULL COMMENT 'id';
ALTER TABLE `commission_calc` ADD `start` varchar(16) NOT NULL COMMENT 'commission start' AFTER `id`;
ALTER TABLE `commission_calc` ADD `end` varchar(16) NOT NULL COMMENT 'commission end' AFTER `start`;
ALTER TABLE `commission_calc` ADD `calc_start` varchar(128) NOT NULL COMMENT 'calc start' AFTER `end`;
ALTER TABLE `commission_calc` ADD `calc_end` varchar(128) NOT NULL COMMENT 'calc end' AFTER `calc_start`;
ALTER TABLE `commission_calc` ADD `status` tinyint unsigned NOT NULL COMMENT 'status' AFTER `calc_end`;
ALTER TABLE `commission_calc` ADD `desc` varchar(256) NOT NULL COMMENT 'desc' AFTER `status`;

----------------------------------------------------
--  `commission_distribute`
----------------------------------------------------
ALTER TABLE `commission_distribute` ADD `id` bigint unsigned NOT NULL COMMENT 'id';
ALTER TABLE `commission_distribute` ADD `start` varchar(16) NOT NULL COMMENT 'commission start' AFTER `id`;
ALTER TABLE `commission_distribute` ADD `end` varchar(16) NOT NULL COMMENT 'commission end' AFTER `start`;
ALTER TABLE `commission_distribute` ADD `distribute_start` varchar(128) NOT NULL COMMENT 'distribute start' AFTER `end`;
ALTER TABLE `commission_distribute` ADD `distribute_end` varchar(128) NOT NULL COMMENT 'distribute end' AFTER `distribute_start`;
ALTER TABLE `commission_distribute` ADD `status` tinyint unsigned NOT NULL COMMENT 'status' AFTER `distribute_end`;
ALTER TABLE `commission_distribute` ADD `desc` varchar(256) NOT NULL COMMENT 'desc' AFTER `status`;

----------------------------------------------------
--  `commission_stat`
----------------------------------------------------
ALTER TABLE `commission_stat` ADD `ctime` bigint NOT NULL COMMENT 'ctime';
ALTER TABLE `commission_stat` ADD `tax_integer` int NOT NULL COMMENT 'tax integer part' AFTER `ctime`;
ALTER TABLE `commission_stat` ADD `tax_decimals` int NOT NULL COMMENT 'tax decimals part' AFTER `tax_integer`;
ALTER TABLE `commission_stat` ADD `channel_integer` int NOT NULL COMMENT 'channel integer part' AFTER `tax_decimals`;
ALTER TABLE `commission_stat` ADD `channel_decimals` int NOT NULL COMMENT 'channel decimals part' AFTER `channel_integer`;
ALTER TABLE `commission_stat` ADD `commission_integer` int NOT NULL COMMENT 'commission integer part' AFTER `channel_decimals`;
ALTER TABLE `commission_stat` ADD `commission_decimals` int NOT NULL COMMENT 'commission decimals part' AFTER `commission_integer`;
ALTER TABLE `commission_stat` ADD `profit_integer` int NOT NULL COMMENT 'profit integer part' AFTER `commission_decimals`;
ALTER TABLE `commission_stat` ADD `profit_decimals` int NOT NULL COMMENT 'profit decimals part' AFTER `profit_integer`;
ALTER TABLE `commission_stat` ADD `mtime` bigint NOT NULL COMMENT 'modified time' AFTER `profit_decimals`;
ALTER TABLE `commission_stat` ADD `status` tinyint unsigned NOT NULL COMMENT 'status' AFTER `mtime`;

----------------------------------------------------
--  `eos_otc_report`
----------------------------------------------------
ALTER TABLE `eos_otc_report` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'eos otc report id';
ALTER TABLE `eos_otc_report` ADD `uid` bigint unsigned NOT NULL COMMENT 'uid' AFTER `id`;
ALTER TABLE `eos_otc_report` ADD `total_order_num` bigint NOT NULL DEFAULT 0 COMMENT 'total_order_num' AFTER `uid`;
ALTER TABLE `eos_otc_report` ADD `success_order_num` bigint NOT NULL DEFAULT 0 COMMENT 'success_order_num' AFTER `total_order_num`;
ALTER TABLE `eos_otc_report` ADD `fail_order_num` bigint NOT NULL DEFAULT 0 COMMENT 'fail_order_num' AFTER `success_order_num`;
ALTER TABLE `eos_otc_report` ADD `buy_eusd_num` bigint NOT NULL DEFAULT 0 COMMENT 'buy_eusd_num' AFTER `fail_order_num`;
ALTER TABLE `eos_otc_report` ADD `sell_eusd_num` bigint NOT NULL DEFAULT 0 COMMENT 'sell_eusd_num' AFTER `buy_eusd_num`;
ALTER TABLE `eos_otc_report` ADD `date` int NOT NULL DEFAULT 0 COMMENT 'date' AFTER `sell_eusd_num`;

----------------------------------------------------
--  `otc_buy`
----------------------------------------------------
ALTER TABLE `otc_buy` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `otc_buy` ADD `available` bigint NOT NULL COMMENT 'available token' AFTER `uid`;
ALTER TABLE `otc_buy` ADD `frozen` bigint unsigned NOT NULL COMMENT 'frozen token' AFTER `available`;
ALTER TABLE `otc_buy` ADD `bought` bigint unsigned NOT NULL COMMENT 'bought token' AFTER `frozen`;
ALTER TABLE `otc_buy` ADD `lower_limit_wechat` bigint NOT NULL COMMENT 'lower limit' AFTER `bought`;
ALTER TABLE `otc_buy` ADD `upper_limit_wechat` bigint NOT NULL COMMENT 'upper limit' AFTER `lower_limit_wechat`;
ALTER TABLE `otc_buy` ADD `lower_limit_bank` bigint NOT NULL COMMENT 'lower limit' AFTER `upper_limit_wechat`;
ALTER TABLE `otc_buy` ADD `upper_limit_bank` bigint NOT NULL COMMENT 'upper limit' AFTER `lower_limit_bank`;
ALTER TABLE `otc_buy` ADD `lower_limit_ali` bigint NOT NULL COMMENT 'lower limit' AFTER `upper_limit_bank`;
ALTER TABLE `otc_buy` ADD `upper_limit_ali` bigint NOT NULL COMMENT 'upper limit' AFTER `lower_limit_ali`;
ALTER TABLE `otc_buy` ADD `pay_type` tinyint unsigned NOT NULL COMMENT 'pay type' AFTER `upper_limit_ali`;
ALTER TABLE `otc_buy` ADD `ctime` bigint NOT NULL COMMENT '' AFTER `pay_type`;

----------------------------------------------------
--  `otc_exchanger`
----------------------------------------------------
ALTER TABLE `otc_exchanger` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `otc_exchanger` ADD `mobile` varchar(100) NOT NULL COMMENT 'Mobile' AFTER `uid`;
ALTER TABLE `otc_exchanger` ADD `wechat` varchar(100) NOT NULL COMMENT 'Wechat' AFTER `mobile`;
ALTER TABLE `otc_exchanger` ADD `telegram` varchar(100) NOT NULL COMMENT 'telegram' AFTER `wechat`;
ALTER TABLE `otc_exchanger` ADD `from` tinyint NOT NULL COMMENT 'from' AFTER `telegram`;
ALTER TABLE `otc_exchanger` ADD `ctime` bigint NOT NULL COMMENT '' AFTER `from`;
ALTER TABLE `otc_exchanger` ADD `utime` bigint NOT NULL COMMENT '' AFTER `ctime`;

----------------------------------------------------
--  `otc_exchanger_verify`
----------------------------------------------------
ALTER TABLE `otc_exchanger_verify` ADD `id` int NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `otc_exchanger_verify` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id' AFTER `id`;
ALTER TABLE `otc_exchanger_verify` ADD `mobile` varchar(100) NOT NULL COMMENT 'Mobile' AFTER `uid`;
ALTER TABLE `otc_exchanger_verify` ADD `wechat` varchar(100) NOT NULL COMMENT 'Wechat' AFTER `mobile`;
ALTER TABLE `otc_exchanger_verify` ADD `telegram` varchar(100) NOT NULL COMMENT 'telegram' AFTER `wechat`;
ALTER TABLE `otc_exchanger_verify` ADD `status` tinyint NOT NULL COMMENT 'status' AFTER `telegram`;
ALTER TABLE `otc_exchanger_verify` ADD `from` tinyint NOT NULL COMMENT 'from' AFTER `status`;
ALTER TABLE `otc_exchanger_verify` ADD `ctime` bigint NOT NULL COMMENT '' AFTER `from`;
ALTER TABLE `otc_exchanger_verify` ADD `utime` bigint NOT NULL COMMENT '' AFTER `ctime`;

----------------------------------------------------
--  `otc_msg`
----------------------------------------------------
ALTER TABLE `otc_msg` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `otc_msg` ADD `order_id` bigint NOT NULL COMMENT 'Order id' AFTER `id`;
ALTER TABLE `otc_msg` ADD `uid` bigint unsigned NOT NULL COMMENT 'Msg Uid' AFTER `order_id`;
ALTER TABLE `otc_msg` ADD `content` text NOT NULL COMMENT 'Msg Content' AFTER `uid`;
ALTER TABLE `otc_msg` ADD `is_read` tinyint unsigned NOT NULL COMMENT 'is_read 1--have read 0--unread' AFTER `content`;
ALTER TABLE `otc_msg` ADD `msg_type` varchar(200) NOT NULL COMMENT 'system--system, text--uid' AFTER `is_read`;
ALTER TABLE `otc_msg` ADD `ctime` bigint NOT NULL COMMENT '' AFTER `msg_type`;

----------------------------------------------------
--  `otc_order`
----------------------------------------------------
ALTER TABLE `otc_order` ADD `id` bigint unsigned NOT NULL COMMENT 'id';
ALTER TABLE `otc_order` ADD `uid` bigint unsigned NOT NULL COMMENT '' AFTER `id`;
ALTER TABLE `otc_order` ADD `uip` varchar(40) NOT NULL DEFAULT '' COMMENT 'user ip' AFTER `uid`;
ALTER TABLE `otc_order` ADD `euid` bigint unsigned NOT NULL DEFAULT 0 COMMENT '' AFTER `uip`;
ALTER TABLE `otc_order` ADD `eip` varchar(40) NOT NULL COMMENT ' ip' AFTER `euid`;
ALTER TABLE `otc_order` ADD `side` tinyint NOT NULL COMMENT '1-buy 2-sell' AFTER `eip`;
ALTER TABLE `otc_order` ADD `amount` bigint NOT NULL COMMENT 'eusd num' AFTER `side`;
ALTER TABLE `otc_order` ADD `price` varchar(100) NOT NULL COMMENT 'eusd => rmb price' AFTER `amount`;
ALTER TABLE `otc_order` ADD `funds` bigint NOT NULL COMMENT 'rmb funds' AFTER `price`;
ALTER TABLE `otc_order` ADD `fee` bigint NOT NULL COMMENT 'rmb price' AFTER `funds`;
ALTER TABLE `otc_order` ADD `pay_id` bigint unsigned NOT NULL COMMENT 'pay id' AFTER `fee`;
ALTER TABLE `otc_order` ADD `pay_type` tinyint NOT NULL COMMENT 'pay type' AFTER `pay_id`;
ALTER TABLE `otc_order` ADD `pay_name` varchar(128) NOT NULL COMMENT '' AFTER `pay_type`;
ALTER TABLE `otc_order` ADD `pay_account` varchar(300) NOT NULL COMMENT '' AFTER `pay_name`;
ALTER TABLE `otc_order` ADD `pay_bank` varchar(128) NOT NULL COMMENT '' AFTER `pay_account`;
ALTER TABLE `otc_order` ADD `pay_bank_branch` varchar(300) NOT NULL COMMENT '' AFTER `pay_bank`;
ALTER TABLE `otc_order` ADD `transfer_id` bigint unsigned NOT NULL COMMENT '' AFTER `pay_bank_branch`;
ALTER TABLE `otc_order` ADD `ctime` bigint NOT NULL COMMENT '' AFTER `transfer_id`;
ALTER TABLE `otc_order` ADD `pay_time` bigint NOT NULL COMMENT '' AFTER `ctime`;
ALTER TABLE `otc_order` ADD `finish_time` bigint NOT NULL COMMENT '' AFTER `pay_time`;
ALTER TABLE `otc_order` ADD `utime` bigint NOT NULL COMMENT '' AFTER `finish_time`;
ALTER TABLE `otc_order` ADD `status` tinyint NOT NULL COMMENT '' AFTER `utime`;
ALTER TABLE `otc_order` ADD `epay_id` bigint unsigned NOT NULL COMMENT 'pay id' AFTER `status`;
ALTER TABLE `otc_order` ADD `epay_type` tinyint NOT NULL COMMENT 'pay type' AFTER `epay_id`;
ALTER TABLE `otc_order` ADD `epay_name` varchar(128) NOT NULL COMMENT '' AFTER `epay_type`;
ALTER TABLE `otc_order` ADD `epay_account` varchar(300) NOT NULL COMMENT '' AFTER `epay_name`;
ALTER TABLE `otc_order` ADD `epay_bank` varchar(128) NOT NULL COMMENT '' AFTER `epay_account`;
ALTER TABLE `otc_order` ADD `epay_bank_branch` varchar(300) NOT NULL COMMENT '' AFTER `epay_bank`;
ALTER TABLE `otc_order` ADD `appeal_status` tinyint NOT NULL DEFAULT 0 COMMENT '' AFTER `epay_bank_branch`;
ALTER TABLE `otc_order` ADD `admin_id` int unsigned NOT NULL DEFAULT 0 COMMENT '' AFTER `appeal_status`;
ALTER TABLE `otc_order` ADD `qr_code` varchar(300) NOT NULL COMMENT '' AFTER `admin_id`;
ALTER TABLE `otc_order` ADD `date` int NOT NULL DEFAULT 0 COMMENT '' AFTER `qr_code`;

----------------------------------------------------
--  `otc_sell`
----------------------------------------------------
ALTER TABLE `otc_sell` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `otc_sell` ADD `available` bigint NOT NULL COMMENT 'available token' AFTER `uid`;
ALTER TABLE `otc_sell` ADD `frozen` bigint NOT NULL COMMENT 'frozen token' AFTER `available`;
ALTER TABLE `otc_sell` ADD `sold` bigint NOT NULL COMMENT 'Sold token' AFTER `frozen`;
ALTER TABLE `otc_sell` ADD `lower_limit` bigint NOT NULL COMMENT 'lower limit' AFTER `sold`;
ALTER TABLE `otc_sell` ADD `upper_limit` bigint NOT NULL COMMENT 'upper limit' AFTER `lower_limit`;
ALTER TABLE `otc_sell` ADD `pay_type` tinyint unsigned NOT NULL COMMENT 'pay type' AFTER `upper_limit`;
ALTER TABLE `otc_sell` ADD `ctime` bigint NOT NULL COMMENT '' AFTER `pay_type`;

----------------------------------------------------
--  `payment_method`
----------------------------------------------------
ALTER TABLE `payment_method` ADD `pmid` bigint unsigned NOT NULL COMMENT 'payment method id';
ALTER TABLE `payment_method` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id' AFTER `pmid`;
ALTER TABLE `payment_method` ADD `mtype` tinyint unsigned NOT NULL COMMENT 'payment method enum' AFTER `uid`;
ALTER TABLE `payment_method` ADD `ord` int unsigned NOT NULL COMMENT 'order serial number' AFTER `mtype`;
ALTER TABLE `payment_method` ADD `name` varchar(128) NOT NULL COMMENT 'name' AFTER `ord`;
ALTER TABLE `payment_method` ADD `account` varchar(128) NOT NULL COMMENT 'account' AFTER `name`;
ALTER TABLE `payment_method` ADD `status` tinyint unsigned NOT NULL COMMENT 'status 0-disable 1-enable ' AFTER `account`;
ALTER TABLE `payment_method` ADD `ctime` int unsigned NOT NULL COMMENT 'created time' AFTER `status`;
ALTER TABLE `payment_method` ADD `bank` varchar(128) NOT NULL COMMENT '' AFTER `ctime`;
ALTER TABLE `payment_method` ADD `bank_branch` varchar(128) NOT NULL COMMENT '' AFTER `bank`;
ALTER TABLE `payment_method` ADD `qr_code` varchar(256) NOT NULL COMMENT '' AFTER `bank_branch`;
ALTER TABLE `payment_method` ADD `qr_code_content` text NOT NULL COMMENT 'qr_code_content' AFTER `qr_code`;
ALTER TABLE `payment_method` ADD `low_money_per_tx_limit` bigint NOT NULL COMMENT '' AFTER `qr_code_content`;
ALTER TABLE `payment_method` ADD `high_money_per_tx_limit` bigint NOT NULL COMMENT '' AFTER `low_money_per_tx_limit`;
ALTER TABLE `payment_method` ADD `times_per_day_limit` bigint NOT NULL COMMENT '' AFTER `high_money_per_tx_limit`;
ALTER TABLE `payment_method` ADD `money_per_day_limit` bigint NOT NULL COMMENT '' AFTER `times_per_day_limit`;
ALTER TABLE `payment_method` ADD `money_sum_limit` bigint NOT NULL COMMENT '' AFTER `money_per_day_limit`;
ALTER TABLE `payment_method` ADD `times_today` bigint NOT NULL COMMENT '' AFTER `money_sum_limit`;
ALTER TABLE `payment_method` ADD `money_today` bigint NOT NULL COMMENT '' AFTER `times_today`;
ALTER TABLE `payment_method` ADD `money_sum` bigint unsigned NOT NULL COMMENT '' AFTER `money_today`;
ALTER TABLE `payment_method` ADD `mtime` int unsigned NOT NULL COMMENT 'modified time' AFTER `money_sum`;
ALTER TABLE `payment_method` ADD `use_time` int unsigned NOT NULL COMMENT 'use time' AFTER `mtime`;

----------------------------------------------------
--  `system_notification`
----------------------------------------------------
ALTER TABLE `system_notification` ADD `nid` bigint unsigned NOT NULL COMMENT 'notification id';
ALTER TABLE `system_notification` ADD `notification_type` varchar(100) NOT NULL COMMENT 'Notification Type' AFTER `nid`;
ALTER TABLE `system_notification` ADD `content` text NOT NULL COMMENT 'content' AFTER `notification_type`;
ALTER TABLE `system_notification` ADD `uid` bigint unsigned NOT NULL COMMENT 'uid' AFTER `content`;
ALTER TABLE `system_notification` ADD `is_read` int NOT NULL COMMENT 'is_read' AFTER `uid`;
ALTER TABLE `system_notification` ADD `ctime` bigint NOT NULL COMMENT '' AFTER `is_read`;

----------------------------------------------------
--  `user`
----------------------------------------------------
ALTER TABLE `user` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `user` ADD `national_code` varchar(100) NOT NULL COMMENT 'national_code' AFTER `uid`;
ALTER TABLE `user` ADD `mobile` varchar(100) NOT NULL COMMENT 'Mobile' AFTER `national_code`;
ALTER TABLE `user` ADD `status` tinyint NOT NULL COMMENT 'status' AFTER `mobile`;
ALTER TABLE `user` ADD `nick` varchar(100) NOT NULL COMMENT 'nick' AFTER `status`;
ALTER TABLE `user` ADD `pass` varchar(100) NOT NULL COMMENT 'password' AFTER `nick`;
ALTER TABLE `user` ADD `salt` varchar(16) NOT NULL COMMENT 'salt' AFTER `pass`;
ALTER TABLE `user` ADD `ctime` bigint NOT NULL COMMENT '' AFTER `salt`;
ALTER TABLE `user` ADD `utime` bigint NOT NULL COMMENT '' AFTER `ctime`;
ALTER TABLE `user` ADD `ip` varchar(100) NOT NULL DEFAULT '' COMMENT 'ip' AFTER `utime`;
ALTER TABLE `user` ADD `last_login_time` bigint NOT NULL DEFAULT 0 COMMENT 'last login time' AFTER `ip`;
ALTER TABLE `user` ADD `last_login_ip` varchar(100) NOT NULL DEFAULT '' COMMENT 'last login ip' AFTER `last_login_time`;
ALTER TABLE `user` ADD `is_exchanger` tinyint NOT NULL DEFAULT 0 COMMENT 'last login ip' AFTER `last_login_ip`;
ALTER TABLE `user` ADD `sign_salt` varchar(256) NOT NULL COMMENT 'sign salt' AFTER `is_exchanger`;

----------------------------------------------------
--  `user_config`
----------------------------------------------------
ALTER TABLE `user_config` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `user_config` ADD `wealth_notice` bool NOT NULL COMMENT 'wealth_notice' AFTER `uid`;
ALTER TABLE `user_config` ADD `order_notice` bool NOT NULL COMMENT 'order_notice' AFTER `wealth_notice`;

----------------------------------------------------
--  `user_login_log`
----------------------------------------------------
ALTER TABLE `user_login_log` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `user_login_log` ADD `user_id` bigint unsigned NOT NULL COMMENT 'user_id' AFTER `id`;
ALTER TABLE `user_login_log` ADD `user_agent` varchar(256) NULL COMMENT '' AFTER `user_id`;
ALTER TABLE `user_login_log` ADD `ips` varchar(256) NULL COMMENT 'ips' AFTER `user_agent`;
ALTER TABLE `user_login_log` ADD `ctime` bigint NOT NULL COMMENT '' AFTER `ips`;

----------------------------------------------------
--  `user_pay_pass`
----------------------------------------------------
ALTER TABLE `user_pay_pass` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `user_pay_pass` ADD `pass` varchar(256) NOT NULL COMMENT 'password' AFTER `uid`;
ALTER TABLE `user_pay_pass` ADD `salt` varchar(256) NOT NULL COMMENT 'salt' AFTER `pass`;
ALTER TABLE `user_pay_pass` ADD `sign_salt` varchar(256) NOT NULL COMMENT 'sign salt' AFTER `salt`;
ALTER TABLE `user_pay_pass` ADD `status` tinyint NOT NULL COMMENT 'status' AFTER `sign_salt`;
ALTER TABLE `user_pay_pass` ADD `method` tinyint NOT NULL COMMENT 'setting pwd method' AFTER `status`;
ALTER TABLE `user_pay_pass` ADD `verify_step` tinyint NOT NULL COMMENT '0 pass verify 1 pass&sms verify' AFTER `method`;
ALTER TABLE `user_pay_pass` ADD `timestamp` bigint NOT NULL DEFAULT 0 COMMENT 'timestamp' AFTER `verify_step`;

