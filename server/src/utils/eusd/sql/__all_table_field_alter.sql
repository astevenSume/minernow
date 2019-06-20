----------------------------------------------------
--  `eos_account`
----------------------------------------------------
ALTER TABLE `eos_account` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `eos_account` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'uid';
ALTER TABLE `eos_account` CHANGE `account` `account` varchar(100) NOT NULL COMMENT 'account';
ALTER TABLE `eos_account` CHANGE `balance` `balance` varchar(100) NOT NULL COMMENT 'balance';
ALTER TABLE `eos_account` CHANGE `status` `status` tinyint NOT NULL COMMENT '';
ALTER TABLE `eos_account` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'created time';
ALTER TABLE `eos_account` CHANGE `utime` `utime` bigint NOT NULL COMMENT 'update time';

----------------------------------------------------
--  `eos_otc`
----------------------------------------------------
ALTER TABLE `eos_otc` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `eos_otc` CHANGE `account` `account` varchar(100) NOT NULL COMMENT 'account';
ALTER TABLE `eos_otc` CHANGE `status` `status` tinyint NOT NULL COMMENT 'status';
ALTER TABLE `eos_otc` CHANGE `available` `available` bigint NOT NULL COMMENT 'available';
ALTER TABLE `eos_otc` CHANGE `trade` `trade` bigint NOT NULL COMMENT 'otc trade balance';
ALTER TABLE `eos_otc` CHANGE `transfer` `transfer` bigint NOT NULL COMMENT 'transfering balance';
ALTER TABLE `eos_otc` CHANGE `sell_state` `sell_state` varchar(200) NOT NULL DEFAULT '' COMMENT '';
ALTER TABLE `eos_otc` CHANGE `sell_pay_type` `sell_pay_type` tinyint unsigned NOT NULL COMMENT '';
ALTER TABLE `eos_otc` CHANGE `sell_able` `sell_able` bool NOT NULL COMMENT '';
ALTER TABLE `eos_otc` CHANGE `sell_rmb_day` `sell_rmb_day` bigint NOT NULL COMMENT '';
ALTER TABLE `eos_otc` CHANGE `sell_rmb_today` `sell_rmb_today` bigint NOT NULL COMMENT '';
ALTER TABLE `eos_otc` CHANGE `sell_rmb_lower_limit` `sell_rmb_lower_limit` bigint NOT NULL COMMENT '';
ALTER TABLE `eos_otc` CHANGE `sell_utime` `sell_utime` bigint NOT NULL COMMENT '';
ALTER TABLE `eos_otc` CHANGE `buy_able` `buy_able` bool NOT NULL COMMENT '';
ALTER TABLE `eos_otc` CHANGE `buy_rmb_day` `buy_rmb_day` bigint NOT NULL COMMENT '';
ALTER TABLE `eos_otc` CHANGE `buy_rmb_today` `buy_rmb_today` bigint NOT NULL COMMENT '';
ALTER TABLE `eos_otc` CHANGE `buy_rmb_lower_limit` `buy_rmb_lower_limit` bigint NOT NULL COMMENT '';
ALTER TABLE `eos_otc` CHANGE `buy_utime` `buy_utime` bigint NOT NULL COMMENT '';
ALTER TABLE `eos_otc` CHANGE `buy_state` `buy_state` varchar(200) NOT NULL DEFAULT '' COMMENT '';
ALTER TABLE `eos_otc` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'created time';
ALTER TABLE `eos_otc` CHANGE `utime` `utime` bigint NOT NULL COMMENT 'update time';

----------------------------------------------------
--  `eos_transaction`
----------------------------------------------------
ALTER TABLE `eos_transaction` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `eos_transaction` CHANGE `type` `type` tinyint unsigned NOT NULL COMMENT 'Transaction type';
ALTER TABLE `eos_transaction` CHANGE `transaction_id` `transaction_id` varchar(100) NOT NULL COMMENT 'Transaction id';
ALTER TABLE `eos_transaction` CHANGE `block_num` `block_num` int unsigned NOT NULL COMMENT 'block num';
ALTER TABLE `eos_transaction` CHANGE `status` `status` tinyint NOT NULL COMMENT 'status';
ALTER TABLE `eos_transaction` CHANGE `payer` `payer` varchar(100) NOT NULL COMMENT 'Payer Account';
ALTER TABLE `eos_transaction` CHANGE `receiver` `receiver` varchar(100) NOT NULL COMMENT 'Receiver Account';
ALTER TABLE `eos_transaction` CHANGE `quantity` `quantity` varchar(100) NOT NULL COMMENT 'Token quantity';
ALTER TABLE `eos_transaction` CHANGE `memo` `memo` varchar(100) NOT NULL COMMENT 'transaction memo';
ALTER TABLE `eos_transaction` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'created time';
ALTER TABLE `eos_transaction` CHANGE `utime` `utime` bigint NOT NULL COMMENT 'update time';

----------------------------------------------------
--  `eos_transaction_info`
----------------------------------------------------
ALTER TABLE `eos_transaction_info` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `eos_transaction_info` CHANGE `transaction_id` `transaction_id` varchar(100) NOT NULL COMMENT 'Transaction id';
ALTER TABLE `eos_transaction_info` CHANGE `block_num` `block_num` int unsigned NOT NULL COMMENT 'block num';
ALTER TABLE `eos_transaction_info` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'created time';
ALTER TABLE `eos_transaction_info` CHANGE `processed` `processed` text NOT NULL COMMENT 'processed info';

----------------------------------------------------
--  `eos_tx_log`
----------------------------------------------------
ALTER TABLE `eos_tx_log` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `eos_tx_log` CHANGE `from` `from` varchar(100) NOT NULL COMMENT ' Account';
ALTER TABLE `eos_tx_log` CHANGE `from_uid` `from_uid` bigint unsigned NOT NULL COMMENT ' Account';
ALTER TABLE `eos_tx_log` CHANGE `to` `to` varchar(100) NOT NULL COMMENT ' Account';
ALTER TABLE `eos_tx_log` CHANGE `to_uid` `to_uid` bigint unsigned NOT NULL COMMENT ' Account';
ALTER TABLE `eos_tx_log` CHANGE `quantity` `quantity` bigint NOT NULL COMMENT '';
ALTER TABLE `eos_tx_log` CHANGE `status` `status` tinyint NOT NULL COMMENT '';
ALTER TABLE `eos_tx_log` CHANGE `log_ids` `log_ids` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `eos_tx_log` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT '';
ALTER TABLE `eos_tx_log` CHANGE `txid` `txid` bigint unsigned NOT NULL COMMENT 'Transaction id';
ALTER TABLE `eos_tx_log` CHANGE `order_id` `order_id` bigint unsigned NOT NULL COMMENT '';
ALTER TABLE `eos_tx_log` CHANGE `utime` `utime` bigint NOT NULL COMMENT '';
ALTER TABLE `eos_tx_log` CHANGE `sign` `sign` varchar(100) NOT NULL DEFAULT '' COMMENT '';
ALTER TABLE `eos_tx_log` CHANGE `delay_deal` `delay_deal` bool NOT NULL DEFAULT false COMMENT '';
ALTER TABLE `eos_tx_log` CHANGE `retry` `retry` int NOT NULL DEFAULT 0 COMMENT '';
ALTER TABLE `eos_tx_log` CHANGE `memo` `memo` varchar(100) NOT NULL DEFAULT '' COMMENT '';

----------------------------------------------------
--  `eos_use_log`
----------------------------------------------------
ALTER TABLE `eos_use_log` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `eos_use_log` CHANGE `type` `type` tinyint unsigned NOT NULL COMMENT 'Transaction type';
ALTER TABLE `eos_use_log` CHANGE `tid` `tid` bigint unsigned NOT NULL COMMENT 'eos_transaction id';
ALTER TABLE `eos_use_log` CHANGE `status` `status` tinyint NOT NULL COMMENT 'status';
ALTER TABLE `eos_use_log` CHANGE `tid_recover` `tid_recover` bigint unsigned NOT NULL COMMENT 'eos_transaction id';
ALTER TABLE `eos_use_log` CHANGE `quantity_num` `quantity_num` bigint unsigned NOT NULL COMMENT 'Eos Num * 10000';

----------------------------------------------------
--  `eos_wealth`
----------------------------------------------------
ALTER TABLE `eos_wealth` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `eos_wealth` CHANGE `status` `status` tinyint NOT NULL COMMENT '';
ALTER TABLE `eos_wealth` CHANGE `account` `account` varchar(100) NOT NULL COMMENT 'account';
ALTER TABLE `eos_wealth` CHANGE `balance` `balance` bigint NOT NULL COMMENT 'balance';
ALTER TABLE `eos_wealth` CHANGE `available` `available` bigint NOT NULL COMMENT 'available balance';
ALTER TABLE `eos_wealth` CHANGE `game` `game` bigint NOT NULL COMMENT 'game balance';
ALTER TABLE `eos_wealth` CHANGE `trade` `trade` bigint NOT NULL COMMENT 'trade frozen balance';
ALTER TABLE `eos_wealth` CHANGE `transfer` `transfer` bigint NOT NULL COMMENT 'transfering balance';
ALTER TABLE `eos_wealth` CHANGE `transfer_game` `transfer_game` bigint NOT NULL COMMENT 'transfering to game balance';
ALTER TABLE `eos_wealth` CHANGE `is_exchanger` `is_exchanger` tinyint NOT NULL COMMENT '';
ALTER TABLE `eos_wealth` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'created time';
ALTER TABLE `eos_wealth` CHANGE `utime` `utime` bigint NOT NULL COMMENT 'update time';

----------------------------------------------------
--  `eos_wealth_log`
----------------------------------------------------
ALTER TABLE `eos_wealth_log` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `eos_wealth_log` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'uid';
ALTER TABLE `eos_wealth_log` CHANGE `uid2` `uid2` bigint unsigned NOT NULL COMMENT 'uid';
ALTER TABLE `eos_wealth_log` CHANGE `ttype` `ttype` tinyint unsigned NOT NULL COMMENT 'Transaction type';
ALTER TABLE `eos_wealth_log` CHANGE `status` `status` tinyint NOT NULL COMMENT 'Transaction status';
ALTER TABLE `eos_wealth_log` CHANGE `txid` `txid` bigint unsigned NOT NULL COMMENT 'Transaction id';
ALTER TABLE `eos_wealth_log` CHANGE `quantity` `quantity` bigint NOT NULL COMMENT 'Token quantity';
ALTER TABLE `eos_wealth_log` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'created time';

----------------------------------------------------
--  `eusd_retire`
----------------------------------------------------
ALTER TABLE `eusd_retire` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `eusd_retire` CHANGE `from` `from` varchar(100) NOT NULL COMMENT ' Account';
ALTER TABLE `eusd_retire` CHANGE `from_uid` `from_uid` bigint unsigned NOT NULL COMMENT ' Account';
ALTER TABLE `eusd_retire` CHANGE `quantity` `quantity` bigint NOT NULL COMMENT '';
ALTER TABLE `eusd_retire` CHANGE `status` `status` tinyint NOT NULL COMMENT '';
ALTER TABLE `eusd_retire` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT '';

----------------------------------------------------
--  `platform_user`
----------------------------------------------------
ALTER TABLE `platform_user` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `platform_user` CHANGE `pid` `pid` int NOT NULL COMMENT '';
ALTER TABLE `platform_user` CHANGE `status` `status` tinyint NOT NULL COMMENT '';
ALTER TABLE `platform_user` CHANGE `ctime` `ctime` int unsigned NOT NULL COMMENT '';

----------------------------------------------------
--  `platform_user_cate`
----------------------------------------------------
ALTER TABLE `platform_user_cate` CHANGE `id` `id` int NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `platform_user_cate` CHANGE `name` `name` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `platform_user_cate` CHANGE `dividend` `dividend` int NOT NULL COMMENT '';
ALTER TABLE `platform_user_cate` CHANGE `ctime` `ctime` int unsigned NOT NULL COMMENT '';

