----------------------------------------------------
--  `eos_account`
----------------------------------------------------
ALTER TABLE `eos_account` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `eos_account` ADD `uid` bigint unsigned NOT NULL COMMENT 'uid' AFTER `id`;
ALTER TABLE `eos_account` ADD `account` varchar(100) NOT NULL COMMENT 'account' AFTER `uid`;
ALTER TABLE `eos_account` ADD `balance` varchar(100) NOT NULL COMMENT 'balance' AFTER `account`;
ALTER TABLE `eos_account` ADD `status` tinyint NOT NULL COMMENT '' AFTER `balance`;
ALTER TABLE `eos_account` ADD `ctime` bigint NOT NULL COMMENT 'created time' AFTER `status`;
ALTER TABLE `eos_account` ADD `utime` bigint NOT NULL COMMENT 'update time' AFTER `ctime`;

----------------------------------------------------
--  `eos_otc`
----------------------------------------------------
ALTER TABLE `eos_otc` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `eos_otc` ADD `account` varchar(100) NOT NULL COMMENT 'account' AFTER `uid`;
ALTER TABLE `eos_otc` ADD `status` tinyint NOT NULL COMMENT 'status' AFTER `account`;
ALTER TABLE `eos_otc` ADD `available` bigint NOT NULL COMMENT 'available' AFTER `status`;
ALTER TABLE `eos_otc` ADD `trade` bigint NOT NULL COMMENT 'otc trade balance' AFTER `available`;
ALTER TABLE `eos_otc` ADD `transfer` bigint NOT NULL COMMENT 'transfering balance' AFTER `trade`;
ALTER TABLE `eos_otc` ADD `sell_state` varchar(200) NOT NULL DEFAULT '' COMMENT '' AFTER `transfer`;
ALTER TABLE `eos_otc` ADD `sell_pay_type` tinyint unsigned NOT NULL COMMENT '' AFTER `sell_state`;
ALTER TABLE `eos_otc` ADD `sell_able` bool NOT NULL COMMENT '' AFTER `sell_pay_type`;
ALTER TABLE `eos_otc` ADD `sell_rmb_day` bigint NOT NULL COMMENT '' AFTER `sell_able`;
ALTER TABLE `eos_otc` ADD `sell_rmb_today` bigint NOT NULL COMMENT '' AFTER `sell_rmb_day`;
ALTER TABLE `eos_otc` ADD `sell_rmb_lower_limit` bigint NOT NULL COMMENT '' AFTER `sell_rmb_today`;
ALTER TABLE `eos_otc` ADD `sell_utime` bigint NOT NULL COMMENT '' AFTER `sell_rmb_lower_limit`;
ALTER TABLE `eos_otc` ADD `buy_able` bool NOT NULL COMMENT '' AFTER `sell_utime`;
ALTER TABLE `eos_otc` ADD `buy_rmb_day` bigint NOT NULL COMMENT '' AFTER `buy_able`;
ALTER TABLE `eos_otc` ADD `buy_rmb_today` bigint NOT NULL COMMENT '' AFTER `buy_rmb_day`;
ALTER TABLE `eos_otc` ADD `buy_rmb_lower_limit` bigint NOT NULL COMMENT '' AFTER `buy_rmb_today`;
ALTER TABLE `eos_otc` ADD `buy_utime` bigint NOT NULL COMMENT '' AFTER `buy_rmb_lower_limit`;
ALTER TABLE `eos_otc` ADD `buy_state` varchar(200) NOT NULL DEFAULT '' COMMENT '' AFTER `buy_utime`;
ALTER TABLE `eos_otc` ADD `ctime` bigint NOT NULL COMMENT 'created time' AFTER `buy_state`;
ALTER TABLE `eos_otc` ADD `utime` bigint NOT NULL COMMENT 'update time' AFTER `ctime`;

----------------------------------------------------
--  `eos_transaction`
----------------------------------------------------
ALTER TABLE `eos_transaction` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `eos_transaction` ADD `type` tinyint unsigned NOT NULL COMMENT 'Transaction type' AFTER `id`;
ALTER TABLE `eos_transaction` ADD `transaction_id` varchar(100) NOT NULL COMMENT 'Transaction id' AFTER `type`;
ALTER TABLE `eos_transaction` ADD `block_num` int unsigned NOT NULL COMMENT 'block num' AFTER `transaction_id`;
ALTER TABLE `eos_transaction` ADD `status` tinyint NOT NULL COMMENT 'status' AFTER `block_num`;
ALTER TABLE `eos_transaction` ADD `payer` varchar(100) NOT NULL COMMENT 'Payer Account' AFTER `status`;
ALTER TABLE `eos_transaction` ADD `receiver` varchar(100) NOT NULL COMMENT 'Receiver Account' AFTER `payer`;
ALTER TABLE `eos_transaction` ADD `quantity` varchar(100) NOT NULL COMMENT 'Token quantity' AFTER `receiver`;
ALTER TABLE `eos_transaction` ADD `memo` varchar(100) NOT NULL COMMENT 'transaction memo' AFTER `quantity`;
ALTER TABLE `eos_transaction` ADD `ctime` bigint NOT NULL COMMENT 'created time' AFTER `memo`;
ALTER TABLE `eos_transaction` ADD `utime` bigint NOT NULL COMMENT 'update time' AFTER `ctime`;

----------------------------------------------------
--  `eos_transaction_info`
----------------------------------------------------
ALTER TABLE `eos_transaction_info` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `eos_transaction_info` ADD `transaction_id` varchar(100) NOT NULL COMMENT 'Transaction id' AFTER `id`;
ALTER TABLE `eos_transaction_info` ADD `block_num` int unsigned NOT NULL COMMENT 'block num' AFTER `transaction_id`;
ALTER TABLE `eos_transaction_info` ADD `ctime` bigint NOT NULL COMMENT 'created time' AFTER `block_num`;
ALTER TABLE `eos_transaction_info` ADD `processed` text NOT NULL COMMENT 'processed info' AFTER `ctime`;

----------------------------------------------------
--  `eos_tx_log`
----------------------------------------------------
ALTER TABLE `eos_tx_log` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `eos_tx_log` ADD `from` varchar(100) NOT NULL COMMENT ' Account' AFTER `id`;
ALTER TABLE `eos_tx_log` ADD `from_uid` bigint unsigned NOT NULL COMMENT ' Account' AFTER `from`;
ALTER TABLE `eos_tx_log` ADD `to` varchar(100) NOT NULL COMMENT ' Account' AFTER `from_uid`;
ALTER TABLE `eos_tx_log` ADD `to_uid` bigint unsigned NOT NULL COMMENT ' Account' AFTER `to`;
ALTER TABLE `eos_tx_log` ADD `quantity` bigint NOT NULL COMMENT '' AFTER `to_uid`;
ALTER TABLE `eos_tx_log` ADD `status` tinyint NOT NULL COMMENT '' AFTER `quantity`;
ALTER TABLE `eos_tx_log` ADD `log_ids` varchar(100) NOT NULL COMMENT '' AFTER `status`;
ALTER TABLE `eos_tx_log` ADD `ctime` bigint NOT NULL COMMENT '' AFTER `log_ids`;
ALTER TABLE `eos_tx_log` ADD `txid` bigint unsigned NOT NULL COMMENT 'Transaction id' AFTER `ctime`;
ALTER TABLE `eos_tx_log` ADD `order_id` bigint unsigned NOT NULL COMMENT '' AFTER `txid`;
ALTER TABLE `eos_tx_log` ADD `utime` bigint NOT NULL COMMENT '' AFTER `order_id`;
ALTER TABLE `eos_tx_log` ADD `sign` varchar(100) NOT NULL DEFAULT '' COMMENT '' AFTER `utime`;
ALTER TABLE `eos_tx_log` ADD `delay_deal` bool NOT NULL DEFAULT false COMMENT '' AFTER `sign`;
ALTER TABLE `eos_tx_log` ADD `retry` int NOT NULL DEFAULT 0 COMMENT '' AFTER `delay_deal`;
ALTER TABLE `eos_tx_log` ADD `memo` varchar(100) NOT NULL DEFAULT '' COMMENT '' AFTER `retry`;

----------------------------------------------------
--  `eos_use_log`
----------------------------------------------------
ALTER TABLE `eos_use_log` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `eos_use_log` ADD `type` tinyint unsigned NOT NULL COMMENT 'Transaction type' AFTER `id`;
ALTER TABLE `eos_use_log` ADD `tid` bigint unsigned NOT NULL COMMENT 'eos_transaction id' AFTER `type`;
ALTER TABLE `eos_use_log` ADD `status` tinyint NOT NULL COMMENT 'status' AFTER `tid`;
ALTER TABLE `eos_use_log` ADD `tid_recover` bigint unsigned NOT NULL COMMENT 'eos_transaction id' AFTER `status`;
ALTER TABLE `eos_use_log` ADD `quantity_num` bigint unsigned NOT NULL COMMENT 'Eos Num * 10000' AFTER `tid_recover`;

----------------------------------------------------
--  `eos_wealth`
----------------------------------------------------
ALTER TABLE `eos_wealth` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `eos_wealth` ADD `status` tinyint NOT NULL COMMENT '' AFTER `uid`;
ALTER TABLE `eos_wealth` ADD `account` varchar(100) NOT NULL COMMENT 'account' AFTER `status`;
ALTER TABLE `eos_wealth` ADD `balance` bigint NOT NULL COMMENT 'balance' AFTER `account`;
ALTER TABLE `eos_wealth` ADD `available` bigint NOT NULL COMMENT 'available balance' AFTER `balance`;
ALTER TABLE `eos_wealth` ADD `game` bigint NOT NULL COMMENT 'game balance' AFTER `available`;
ALTER TABLE `eos_wealth` ADD `trade` bigint NOT NULL COMMENT 'trade frozen balance' AFTER `game`;
ALTER TABLE `eos_wealth` ADD `transfer` bigint NOT NULL COMMENT 'transfering balance' AFTER `trade`;
ALTER TABLE `eos_wealth` ADD `transfer_game` bigint NOT NULL COMMENT 'transfering to game balance' AFTER `transfer`;
ALTER TABLE `eos_wealth` ADD `is_exchanger` tinyint NOT NULL COMMENT '' AFTER `transfer_game`;
ALTER TABLE `eos_wealth` ADD `ctime` bigint NOT NULL COMMENT 'created time' AFTER `is_exchanger`;
ALTER TABLE `eos_wealth` ADD `utime` bigint NOT NULL COMMENT 'update time' AFTER `ctime`;

----------------------------------------------------
--  `eos_wealth_log`
----------------------------------------------------
ALTER TABLE `eos_wealth_log` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `eos_wealth_log` ADD `uid` bigint unsigned NOT NULL COMMENT 'uid' AFTER `id`;
ALTER TABLE `eos_wealth_log` ADD `uid2` bigint unsigned NOT NULL COMMENT 'uid' AFTER `uid`;
ALTER TABLE `eos_wealth_log` ADD `ttype` tinyint unsigned NOT NULL COMMENT 'Transaction type' AFTER `uid2`;
ALTER TABLE `eos_wealth_log` ADD `status` tinyint NOT NULL COMMENT 'Transaction status' AFTER `ttype`;
ALTER TABLE `eos_wealth_log` ADD `txid` bigint unsigned NOT NULL COMMENT 'Transaction id' AFTER `status`;
ALTER TABLE `eos_wealth_log` ADD `quantity` bigint NOT NULL COMMENT 'Token quantity' AFTER `txid`;
ALTER TABLE `eos_wealth_log` ADD `ctime` bigint NOT NULL COMMENT 'created time' AFTER `quantity`;

----------------------------------------------------
--  `eusd_retire`
----------------------------------------------------
ALTER TABLE `eusd_retire` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `eusd_retire` ADD `from` varchar(100) NOT NULL COMMENT ' Account' AFTER `id`;
ALTER TABLE `eusd_retire` ADD `from_uid` bigint unsigned NOT NULL COMMENT ' Account' AFTER `from`;
ALTER TABLE `eusd_retire` ADD `quantity` bigint NOT NULL COMMENT '' AFTER `from_uid`;
ALTER TABLE `eusd_retire` ADD `status` tinyint NOT NULL COMMENT '' AFTER `quantity`;
ALTER TABLE `eusd_retire` ADD `ctime` bigint NOT NULL COMMENT '' AFTER `status`;

----------------------------------------------------
--  `platform_user`
----------------------------------------------------
ALTER TABLE `platform_user` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `platform_user` ADD `pid` int NOT NULL COMMENT '' AFTER `uid`;
ALTER TABLE `platform_user` ADD `status` tinyint NOT NULL COMMENT '' AFTER `pid`;
ALTER TABLE `platform_user` ADD `ctime` int unsigned NOT NULL COMMENT '' AFTER `status`;

----------------------------------------------------
--  `platform_user_cate`
----------------------------------------------------
ALTER TABLE `platform_user_cate` ADD `id` int NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `platform_user_cate` ADD `name` varchar(100) NOT NULL COMMENT '' AFTER `id`;
ALTER TABLE `platform_user_cate` ADD `dividend` int NOT NULL COMMENT '' AFTER `name`;
ALTER TABLE `platform_user_cate` ADD `ctime` int unsigned NOT NULL COMMENT '' AFTER `dividend`;

