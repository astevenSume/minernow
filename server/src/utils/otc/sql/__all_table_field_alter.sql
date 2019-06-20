----------------------------------------------------
--  `appeal`
----------------------------------------------------
ALTER TABLE `appeal` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `appeal` CHANGE `type` `type` tinyint NOT NULL COMMENT 'type';
ALTER TABLE `appeal` CHANGE `user_id` `user_id` bigint unsigned NOT NULL COMMENT 'user_id';
ALTER TABLE `appeal` CHANGE `admin_id` `admin_id` int unsigned NOT NULL COMMENT 'admin_id';
ALTER TABLE `appeal` CHANGE `order_id` `order_id` bigint unsigned NOT NULL COMMENT 'order_id';
ALTER TABLE `appeal` CHANGE `context` `context` varchar(256) NOT NULL COMMENT 'context';
ALTER TABLE `appeal` CHANGE `status` `status` tinyint NOT NULL COMMENT 'status';
ALTER TABLE `appeal` CHANGE `ctime` `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time';
ALTER TABLE `appeal` CHANGE `utime` `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time';
ALTER TABLE `appeal` CHANGE `wechat` `wechat` varchar(100) NOT NULL COMMENT 'Wechat';

----------------------------------------------------
--  `appeal_deal_log`
----------------------------------------------------
ALTER TABLE `appeal_deal_log` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `appeal_deal_log` CHANGE `appeal_id` `appeal_id` bigint unsigned NOT NULL COMMENT 'appeal_id';
ALTER TABLE `appeal_deal_log` CHANGE `admin_id` `admin_id` int unsigned NOT NULL COMMENT 'admin_id';
ALTER TABLE `appeal_deal_log` CHANGE `order_id` `order_id` bigint unsigned NOT NULL COMMENT 'order_id';
ALTER TABLE `appeal_deal_log` CHANGE `action` `action` tinyint NOT NULL COMMENT 'action';
ALTER TABLE `appeal_deal_log` CHANGE `ctime` `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time';

----------------------------------------------------
--  `commission_calc`
----------------------------------------------------
ALTER TABLE `commission_calc` CHANGE `id` `id` bigint unsigned NOT NULL COMMENT 'id';
ALTER TABLE `commission_calc` CHANGE `start` `start` varchar(16) NOT NULL COMMENT 'commission start';
ALTER TABLE `commission_calc` CHANGE `end` `end` varchar(16) NOT NULL COMMENT 'commission end';
ALTER TABLE `commission_calc` CHANGE `calc_start` `calc_start` varchar(128) NOT NULL COMMENT 'calc start';
ALTER TABLE `commission_calc` CHANGE `calc_end` `calc_end` varchar(128) NOT NULL COMMENT 'calc end';
ALTER TABLE `commission_calc` CHANGE `status` `status` tinyint unsigned NOT NULL COMMENT 'status';
ALTER TABLE `commission_calc` CHANGE `desc` `desc` varchar(256) NOT NULL COMMENT 'desc';

----------------------------------------------------
--  `commission_distribute`
----------------------------------------------------
ALTER TABLE `commission_distribute` CHANGE `id` `id` bigint unsigned NOT NULL COMMENT 'id';
ALTER TABLE `commission_distribute` CHANGE `start` `start` varchar(16) NOT NULL COMMENT 'commission start';
ALTER TABLE `commission_distribute` CHANGE `end` `end` varchar(16) NOT NULL COMMENT 'commission end';
ALTER TABLE `commission_distribute` CHANGE `distribute_start` `distribute_start` varchar(128) NOT NULL COMMENT 'distribute start';
ALTER TABLE `commission_distribute` CHANGE `distribute_end` `distribute_end` varchar(128) NOT NULL COMMENT 'distribute end';
ALTER TABLE `commission_distribute` CHANGE `status` `status` tinyint unsigned NOT NULL COMMENT 'status';
ALTER TABLE `commission_distribute` CHANGE `desc` `desc` varchar(256) NOT NULL COMMENT 'desc';

----------------------------------------------------
--  `commission_stat`
----------------------------------------------------
ALTER TABLE `commission_stat` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'ctime';
ALTER TABLE `commission_stat` CHANGE `tax_integer` `tax_integer` int NOT NULL COMMENT 'tax integer part';
ALTER TABLE `commission_stat` CHANGE `tax_decimals` `tax_decimals` int NOT NULL COMMENT 'tax decimals part';
ALTER TABLE `commission_stat` CHANGE `channel_integer` `channel_integer` int NOT NULL COMMENT 'channel integer part';
ALTER TABLE `commission_stat` CHANGE `channel_decimals` `channel_decimals` int NOT NULL COMMENT 'channel decimals part';
ALTER TABLE `commission_stat` CHANGE `commission_integer` `commission_integer` int NOT NULL COMMENT 'commission integer part';
ALTER TABLE `commission_stat` CHANGE `commission_decimals` `commission_decimals` int NOT NULL COMMENT 'commission decimals part';
ALTER TABLE `commission_stat` CHANGE `profit_integer` `profit_integer` int NOT NULL COMMENT 'profit integer part';
ALTER TABLE `commission_stat` CHANGE `profit_decimals` `profit_decimals` int NOT NULL COMMENT 'profit decimals part';
ALTER TABLE `commission_stat` CHANGE `mtime` `mtime` bigint NOT NULL COMMENT 'modified time';
ALTER TABLE `commission_stat` CHANGE `status` `status` tinyint unsigned NOT NULL COMMENT 'status';

----------------------------------------------------
--  `eos_otc_report`
----------------------------------------------------
ALTER TABLE `eos_otc_report` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'eos otc report id';
ALTER TABLE `eos_otc_report` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'uid';
ALTER TABLE `eos_otc_report` CHANGE `total_order_num` `total_order_num` bigint NOT NULL DEFAULT 0 COMMENT 'total_order_num';
ALTER TABLE `eos_otc_report` CHANGE `success_order_num` `success_order_num` bigint NOT NULL DEFAULT 0 COMMENT 'success_order_num';
ALTER TABLE `eos_otc_report` CHANGE `fail_order_num` `fail_order_num` bigint NOT NULL DEFAULT 0 COMMENT 'fail_order_num';
ALTER TABLE `eos_otc_report` CHANGE `buy_eusd_num` `buy_eusd_num` bigint NOT NULL DEFAULT 0 COMMENT 'buy_eusd_num';
ALTER TABLE `eos_otc_report` CHANGE `sell_eusd_num` `sell_eusd_num` bigint NOT NULL DEFAULT 0 COMMENT 'sell_eusd_num';
ALTER TABLE `eos_otc_report` CHANGE `date` `date` int NOT NULL DEFAULT 0 COMMENT 'date';

----------------------------------------------------
--  `otc_buy`
----------------------------------------------------
ALTER TABLE `otc_buy` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `otc_buy` CHANGE `available` `available` bigint NOT NULL COMMENT 'available token';
ALTER TABLE `otc_buy` CHANGE `frozen` `frozen` bigint unsigned NOT NULL COMMENT 'frozen token';
ALTER TABLE `otc_buy` CHANGE `bought` `bought` bigint unsigned NOT NULL COMMENT 'bought token';
ALTER TABLE `otc_buy` CHANGE `lower_limit_wechat` `lower_limit_wechat` bigint NOT NULL COMMENT 'lower limit';
ALTER TABLE `otc_buy` CHANGE `upper_limit_wechat` `upper_limit_wechat` bigint NOT NULL COMMENT 'upper limit';
ALTER TABLE `otc_buy` CHANGE `lower_limit_bank` `lower_limit_bank` bigint NOT NULL COMMENT 'lower limit';
ALTER TABLE `otc_buy` CHANGE `upper_limit_bank` `upper_limit_bank` bigint NOT NULL COMMENT 'upper limit';
ALTER TABLE `otc_buy` CHANGE `lower_limit_ali` `lower_limit_ali` bigint NOT NULL COMMENT 'lower limit';
ALTER TABLE `otc_buy` CHANGE `upper_limit_ali` `upper_limit_ali` bigint NOT NULL COMMENT 'upper limit';
ALTER TABLE `otc_buy` CHANGE `pay_type` `pay_type` tinyint unsigned NOT NULL COMMENT 'pay type';
ALTER TABLE `otc_buy` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT '';

----------------------------------------------------
--  `otc_exchanger`
----------------------------------------------------
ALTER TABLE `otc_exchanger` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `otc_exchanger` CHANGE `mobile` `mobile` varchar(100) NOT NULL COMMENT 'Mobile';
ALTER TABLE `otc_exchanger` CHANGE `wechat` `wechat` varchar(100) NOT NULL COMMENT 'Wechat';
ALTER TABLE `otc_exchanger` CHANGE `telegram` `telegram` varchar(100) NOT NULL COMMENT 'telegram';
ALTER TABLE `otc_exchanger` CHANGE `from` `from` tinyint NOT NULL COMMENT 'from';
ALTER TABLE `otc_exchanger` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT '';
ALTER TABLE `otc_exchanger` CHANGE `utime` `utime` bigint NOT NULL COMMENT '';

----------------------------------------------------
--  `otc_exchanger_verify`
----------------------------------------------------
ALTER TABLE `otc_exchanger_verify` CHANGE `id` `id` int NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `otc_exchanger_verify` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `otc_exchanger_verify` CHANGE `mobile` `mobile` varchar(100) NOT NULL COMMENT 'Mobile';
ALTER TABLE `otc_exchanger_verify` CHANGE `wechat` `wechat` varchar(100) NOT NULL COMMENT 'Wechat';
ALTER TABLE `otc_exchanger_verify` CHANGE `telegram` `telegram` varchar(100) NOT NULL COMMENT 'telegram';
ALTER TABLE `otc_exchanger_verify` CHANGE `status` `status` tinyint NOT NULL COMMENT 'status';
ALTER TABLE `otc_exchanger_verify` CHANGE `from` `from` tinyint NOT NULL COMMENT 'from';
ALTER TABLE `otc_exchanger_verify` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT '';
ALTER TABLE `otc_exchanger_verify` CHANGE `utime` `utime` bigint NOT NULL COMMENT '';

----------------------------------------------------
--  `otc_msg`
----------------------------------------------------
ALTER TABLE `otc_msg` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `otc_msg` CHANGE `order_id` `order_id` bigint NOT NULL COMMENT 'Order id';
ALTER TABLE `otc_msg` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'Msg Uid';
ALTER TABLE `otc_msg` CHANGE `content` `content` text NOT NULL COMMENT 'Msg Content';
ALTER TABLE `otc_msg` CHANGE `is_read` `is_read` tinyint unsigned NOT NULL COMMENT 'is_read 1--have read 0--unread';
ALTER TABLE `otc_msg` CHANGE `msg_type` `msg_type` varchar(200) NOT NULL COMMENT 'system--system, text--uid';
ALTER TABLE `otc_msg` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT '';

----------------------------------------------------
--  `otc_order`
----------------------------------------------------
ALTER TABLE `otc_order` CHANGE `id` `id` bigint unsigned NOT NULL COMMENT 'id';
ALTER TABLE `otc_order` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT '';
ALTER TABLE `otc_order` CHANGE `uip` `uip` varchar(40) NOT NULL DEFAULT '' COMMENT 'user ip';
ALTER TABLE `otc_order` CHANGE `euid` `euid` bigint unsigned NOT NULL DEFAULT 0 COMMENT '';
ALTER TABLE `otc_order` CHANGE `eip` `eip` varchar(40) NOT NULL COMMENT ' ip';
ALTER TABLE `otc_order` CHANGE `side` `side` tinyint NOT NULL COMMENT '1-buy 2-sell';
ALTER TABLE `otc_order` CHANGE `amount` `amount` bigint NOT NULL COMMENT 'eusd num';
ALTER TABLE `otc_order` CHANGE `price` `price` varchar(100) NOT NULL COMMENT 'eusd => rmb price';
ALTER TABLE `otc_order` CHANGE `funds` `funds` bigint NOT NULL COMMENT 'rmb funds';
ALTER TABLE `otc_order` CHANGE `fee` `fee` bigint NOT NULL COMMENT 'rmb price';
ALTER TABLE `otc_order` CHANGE `pay_id` `pay_id` bigint unsigned NOT NULL COMMENT 'pay id';
ALTER TABLE `otc_order` CHANGE `pay_type` `pay_type` tinyint NOT NULL COMMENT 'pay type';
ALTER TABLE `otc_order` CHANGE `pay_name` `pay_name` varchar(128) NOT NULL COMMENT '';
ALTER TABLE `otc_order` CHANGE `pay_account` `pay_account` varchar(300) NOT NULL COMMENT '';
ALTER TABLE `otc_order` CHANGE `pay_bank` `pay_bank` varchar(128) NOT NULL COMMENT '';
ALTER TABLE `otc_order` CHANGE `pay_bank_branch` `pay_bank_branch` varchar(300) NOT NULL COMMENT '';
ALTER TABLE `otc_order` CHANGE `transfer_id` `transfer_id` bigint unsigned NOT NULL COMMENT '';
ALTER TABLE `otc_order` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT '';
ALTER TABLE `otc_order` CHANGE `pay_time` `pay_time` bigint NOT NULL COMMENT '';
ALTER TABLE `otc_order` CHANGE `finish_time` `finish_time` bigint NOT NULL COMMENT '';
ALTER TABLE `otc_order` CHANGE `utime` `utime` bigint NOT NULL COMMENT '';
ALTER TABLE `otc_order` CHANGE `status` `status` tinyint NOT NULL COMMENT '';
ALTER TABLE `otc_order` CHANGE `epay_id` `epay_id` bigint unsigned NOT NULL COMMENT 'pay id';
ALTER TABLE `otc_order` CHANGE `epay_type` `epay_type` tinyint NOT NULL COMMENT 'pay type';
ALTER TABLE `otc_order` CHANGE `epay_name` `epay_name` varchar(128) NOT NULL COMMENT '';
ALTER TABLE `otc_order` CHANGE `epay_account` `epay_account` varchar(300) NOT NULL COMMENT '';
ALTER TABLE `otc_order` CHANGE `epay_bank` `epay_bank` varchar(128) NOT NULL COMMENT '';
ALTER TABLE `otc_order` CHANGE `epay_bank_branch` `epay_bank_branch` varchar(300) NOT NULL COMMENT '';
ALTER TABLE `otc_order` CHANGE `appeal_status` `appeal_status` tinyint NOT NULL DEFAULT 0 COMMENT '';
ALTER TABLE `otc_order` CHANGE `admin_id` `admin_id` int unsigned NOT NULL DEFAULT 0 COMMENT '';
ALTER TABLE `otc_order` CHANGE `qr_code` `qr_code` varchar(300) NOT NULL COMMENT '';
ALTER TABLE `otc_order` CHANGE `date` `date` int NOT NULL DEFAULT 0 COMMENT '';

----------------------------------------------------
--  `otc_sell`
----------------------------------------------------
ALTER TABLE `otc_sell` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `otc_sell` CHANGE `available` `available` bigint NOT NULL COMMENT 'available token';
ALTER TABLE `otc_sell` CHANGE `frozen` `frozen` bigint NOT NULL COMMENT 'frozen token';
ALTER TABLE `otc_sell` CHANGE `sold` `sold` bigint NOT NULL COMMENT 'Sold token';
ALTER TABLE `otc_sell` CHANGE `lower_limit` `lower_limit` bigint NOT NULL COMMENT 'lower limit';
ALTER TABLE `otc_sell` CHANGE `upper_limit` `upper_limit` bigint NOT NULL COMMENT 'upper limit';
ALTER TABLE `otc_sell` CHANGE `pay_type` `pay_type` tinyint unsigned NOT NULL COMMENT 'pay type';
ALTER TABLE `otc_sell` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT '';

----------------------------------------------------
--  `payment_method`
----------------------------------------------------
ALTER TABLE `payment_method` CHANGE `pmid` `pmid` bigint unsigned NOT NULL COMMENT 'payment method id';
ALTER TABLE `payment_method` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `payment_method` CHANGE `mtype` `mtype` tinyint unsigned NOT NULL COMMENT 'payment method enum';
ALTER TABLE `payment_method` CHANGE `ord` `ord` int unsigned NOT NULL COMMENT 'order serial number';
ALTER TABLE `payment_method` CHANGE `name` `name` varchar(128) NOT NULL COMMENT 'name';
ALTER TABLE `payment_method` CHANGE `account` `account` varchar(128) NOT NULL COMMENT 'account';
ALTER TABLE `payment_method` CHANGE `status` `status` tinyint unsigned NOT NULL COMMENT 'status 0-disable 1-enable ';
ALTER TABLE `payment_method` CHANGE `ctime` `ctime` int unsigned NOT NULL COMMENT 'created time';
ALTER TABLE `payment_method` CHANGE `bank` `bank` varchar(128) NOT NULL COMMENT '';
ALTER TABLE `payment_method` CHANGE `bank_branch` `bank_branch` varchar(128) NOT NULL COMMENT '';
ALTER TABLE `payment_method` CHANGE `qr_code` `qr_code` varchar(256) NOT NULL COMMENT '';
ALTER TABLE `payment_method` CHANGE `qr_code_content` `qr_code_content` text NOT NULL COMMENT 'qr_code_content';
ALTER TABLE `payment_method` CHANGE `low_money_per_tx_limit` `low_money_per_tx_limit` bigint NOT NULL COMMENT '';
ALTER TABLE `payment_method` CHANGE `high_money_per_tx_limit` `high_money_per_tx_limit` bigint NOT NULL COMMENT '';
ALTER TABLE `payment_method` CHANGE `times_per_day_limit` `times_per_day_limit` bigint NOT NULL COMMENT '';
ALTER TABLE `payment_method` CHANGE `money_per_day_limit` `money_per_day_limit` bigint NOT NULL COMMENT '';
ALTER TABLE `payment_method` CHANGE `money_sum_limit` `money_sum_limit` bigint NOT NULL COMMENT '';
ALTER TABLE `payment_method` CHANGE `times_today` `times_today` bigint NOT NULL COMMENT '';
ALTER TABLE `payment_method` CHANGE `money_today` `money_today` bigint NOT NULL COMMENT '';
ALTER TABLE `payment_method` CHANGE `money_sum` `money_sum` bigint unsigned NOT NULL COMMENT '';
ALTER TABLE `payment_method` CHANGE `mtime` `mtime` int unsigned NOT NULL COMMENT 'modified time';
ALTER TABLE `payment_method` CHANGE `use_time` `use_time` int unsigned NOT NULL COMMENT 'use time';

----------------------------------------------------
--  `system_notification`
----------------------------------------------------
ALTER TABLE `system_notification` CHANGE `nid` `nid` bigint unsigned NOT NULL COMMENT 'notification id';
ALTER TABLE `system_notification` CHANGE `notification_type` `notification_type` varchar(100) NOT NULL COMMENT 'Notification Type';
ALTER TABLE `system_notification` CHANGE `content` `content` text NOT NULL COMMENT 'content';
ALTER TABLE `system_notification` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'uid';
ALTER TABLE `system_notification` CHANGE `is_read` `is_read` int NOT NULL COMMENT 'is_read';
ALTER TABLE `system_notification` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT '';

----------------------------------------------------
--  `user`
----------------------------------------------------
ALTER TABLE `user` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `user` CHANGE `national_code` `national_code` varchar(100) NOT NULL COMMENT 'national_code';
ALTER TABLE `user` CHANGE `mobile` `mobile` varchar(100) NOT NULL COMMENT 'Mobile';
ALTER TABLE `user` CHANGE `status` `status` tinyint NOT NULL COMMENT 'status';
ALTER TABLE `user` CHANGE `nick` `nick` varchar(100) NOT NULL COMMENT 'nick';
ALTER TABLE `user` CHANGE `pass` `pass` varchar(100) NOT NULL COMMENT 'password';
ALTER TABLE `user` CHANGE `salt` `salt` varchar(16) NOT NULL COMMENT 'salt';
ALTER TABLE `user` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT '';
ALTER TABLE `user` CHANGE `utime` `utime` bigint NOT NULL COMMENT '';
ALTER TABLE `user` CHANGE `ip` `ip` varchar(100) NOT NULL DEFAULT '' COMMENT 'ip';
ALTER TABLE `user` CHANGE `last_login_time` `last_login_time` bigint NOT NULL DEFAULT 0 COMMENT 'last login time';
ALTER TABLE `user` CHANGE `last_login_ip` `last_login_ip` varchar(100) NOT NULL DEFAULT '' COMMENT 'last login ip';
ALTER TABLE `user` CHANGE `is_exchanger` `is_exchanger` tinyint NOT NULL DEFAULT 0 COMMENT 'last login ip';
ALTER TABLE `user` CHANGE `sign_salt` `sign_salt` varchar(256) NOT NULL COMMENT 'sign salt';

----------------------------------------------------
--  `user_config`
----------------------------------------------------
ALTER TABLE `user_config` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `user_config` CHANGE `wealth_notice` `wealth_notice` bool NOT NULL COMMENT 'wealth_notice';
ALTER TABLE `user_config` CHANGE `order_notice` `order_notice` bool NOT NULL COMMENT 'order_notice';

----------------------------------------------------
--  `user_login_log`
----------------------------------------------------
ALTER TABLE `user_login_log` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `user_login_log` CHANGE `user_id` `user_id` bigint unsigned NOT NULL COMMENT 'user_id';
ALTER TABLE `user_login_log` CHANGE `user_agent` `user_agent` varchar(256) NULL COMMENT '';
ALTER TABLE `user_login_log` CHANGE `ips` `ips` varchar(256) NULL COMMENT 'ips';
ALTER TABLE `user_login_log` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT '';

----------------------------------------------------
--  `user_pay_pass`
----------------------------------------------------
ALTER TABLE `user_pay_pass` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `user_pay_pass` CHANGE `pass` `pass` varchar(256) NOT NULL COMMENT 'password';
ALTER TABLE `user_pay_pass` CHANGE `salt` `salt` varchar(256) NOT NULL COMMENT 'salt';
ALTER TABLE `user_pay_pass` CHANGE `sign_salt` `sign_salt` varchar(256) NOT NULL COMMENT 'sign salt';
ALTER TABLE `user_pay_pass` CHANGE `status` `status` tinyint NOT NULL COMMENT 'status';
ALTER TABLE `user_pay_pass` CHANGE `method` `method` tinyint NOT NULL COMMENT 'setting pwd method';
ALTER TABLE `user_pay_pass` CHANGE `verify_step` `verify_step` tinyint NOT NULL COMMENT '0 pass verify 1 pass&sms verify';
ALTER TABLE `user_pay_pass` CHANGE `timestamp` `timestamp` bigint NOT NULL DEFAULT 0 COMMENT 'timestamp';

