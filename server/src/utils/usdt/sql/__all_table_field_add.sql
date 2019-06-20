----------------------------------------------------
--  `market_prices`
----------------------------------------------------
ALTER TABLE `market_prices` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'market otc price record id';
ALTER TABLE `market_prices` ADD `market` tinyint NOT NULL COMMENT 'market' AFTER `id`;
ALTER TABLE `market_prices` ADD `currency` varchar(100) NOT NULL COMMENT 'currency' AFTER `market`;
ALTER TABLE `market_prices` ADD `trade_method` tinyint NOT NULL COMMENT 'trade method' AFTER `currency`;
ALTER TABLE `market_prices` ADD `pow_price` bigint unsigned NOT NULL COMMENT 'pow price' AFTER `trade_method`;
ALTER TABLE `market_prices` ADD `pow` int NOT NULL DEFAULT 4 COMMENT 'pow_price = (real price) * pow' AFTER `pow_price`;
ALTER TABLE `market_prices` ADD `ctime` bigint NOT NULL COMMENT 'created time' AFTER `pow`;

----------------------------------------------------
--  `prices`
----------------------------------------------------
ALTER TABLE `prices` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'otc price record id';
ALTER TABLE `prices` ADD `currency` varchar(100) NOT NULL COMMENT 'currency' AFTER `id`;
ALTER TABLE `prices` ADD `pow_price` bigint unsigned NOT NULL COMMENT 'pow price' AFTER `currency`;
ALTER TABLE `prices` ADD `pow` int NOT NULL DEFAULT 4 COMMENT 'pow_price = (real price) * pow' AFTER `pow_price`;
ALTER TABLE `prices` ADD `ctime` bigint NOT NULL COMMENT 'created time' AFTER `pow`;

----------------------------------------------------
--  `usdt_account`
----------------------------------------------------
ALTER TABLE `usdt_account` ADD `uaid` bigint unsigned NOT NULL COMMENT '';
ALTER TABLE `usdt_account` ADD `uid` bigint unsigned NOT NULL COMMENT '' AFTER `uaid`;
ALTER TABLE `usdt_account` ADD `status` tinyint unsigned NOT NULL COMMENT 'status 0-locked 1-working' AFTER `uid`;
ALTER TABLE `usdt_account` ADD `available_integer` bigint NOT NULL COMMENT 'available amount integer part' AFTER `status`;
ALTER TABLE `usdt_account` ADD `frozen_integer` bigint NOT NULL COMMENT 'frozen amount integer part' AFTER `available_integer`;
ALTER TABLE `usdt_account` ADD `transfer_frozen_integer` bigint NOT NULL COMMENT 'frozen amount integer part' AFTER `frozen_integer`;
ALTER TABLE `usdt_account` ADD `mortgaged_integer` bigint NOT NULL COMMENT 'mortgaged amount integer part' AFTER `transfer_frozen_integer`;
ALTER TABLE `usdt_account` ADD `btc_available_integer` bigint NOT NULL COMMENT 'available btc amount integer part' AFTER `mortgaged_integer`;
ALTER TABLE `usdt_account` ADD `btc_frozen_integer` bigint NOT NULL COMMENT 'frozen btc amount integer part' AFTER `btc_available_integer`;
ALTER TABLE `usdt_account` ADD `btc_mortgaged_integer` bigint NOT NULL COMMENT 'mortgaged btc amount integer part' AFTER `btc_frozen_integer`;
ALTER TABLE `usdt_account` ADD `waiting_cash_sweep_integer` bigint NOT NULL COMMENT 'waiting cash sweep integer part' AFTER `btc_mortgaged_integer`;
ALTER TABLE `usdt_account` ADD `cash_sweep_integer` bigint NOT NULL COMMENT 'cash sweep integer part' AFTER `waiting_cash_sweep_integer`;
ALTER TABLE `usdt_account` ADD `owned_by_platform_integer` bigint NOT NULL COMMENT 'owned by platform integer part' AFTER `cash_sweep_integer`;
ALTER TABLE `usdt_account` ADD `sweep_status` tinyint unsigned NOT NULL COMMENT 'sweep status' AFTER `owned_by_platform_integer`;
ALTER TABLE `usdt_account` ADD `pkid` bigint unsigned NOT NULL COMMENT 'usdt private key id' AFTER `sweep_status`;
ALTER TABLE `usdt_account` ADD `address` varchar(100) NOT NULL COMMENT '' AFTER `pkid`;
ALTER TABLE `usdt_account` ADD `ctime` bigint NOT NULL COMMENT 'created time' AFTER `address`;
ALTER TABLE `usdt_account` ADD `mtime` bigint NOT NULL COMMENT 'update time' AFTER `ctime`;
ALTER TABLE `usdt_account` ADD `sign` varchar(256) NOT NULL COMMENT 'sign' AFTER `mtime`;

----------------------------------------------------
--  `usdt_onchain_balance`
----------------------------------------------------
ALTER TABLE `usdt_onchain_balance` ADD `address` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_balance` ADD `property_id` int unsigned NOT NULL COMMENT '' AFTER `address`;
ALTER TABLE `usdt_onchain_balance` ADD `pending_pos` varchar(100) NOT NULL COMMENT '' AFTER `property_id`;
ALTER TABLE `usdt_onchain_balance` ADD `reserved` varchar(100) NOT NULL COMMENT '' AFTER `pending_pos`;
ALTER TABLE `usdt_onchain_balance` ADD `divisible` bool NOT NULL COMMENT '' AFTER `reserved`;
ALTER TABLE `usdt_onchain_balance` ADD `amount_integer` bigint NOT NULL COMMENT 'amount integer part' AFTER `divisible`;
ALTER TABLE `usdt_onchain_balance` ADD `frozen` varchar(100) NOT NULL COMMENT '' AFTER `amount_integer`;
ALTER TABLE `usdt_onchain_balance` ADD `pending_neg` varchar(100) NOT NULL COMMENT '' AFTER `frozen`;
ALTER TABLE `usdt_onchain_balance` ADD `mtime` bigint NOT NULL COMMENT '' AFTER `pending_neg`;

----------------------------------------------------
--  `usdt_onchain_data`
----------------------------------------------------
ALTER TABLE `usdt_onchain_data` ADD `address` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_data` ADD `attr_type` int unsigned NOT NULL COMMENT '' AFTER `address`;
ALTER TABLE `usdt_onchain_data` ADD `data_int64` bigint NOT NULL COMMENT '' AFTER `attr_type`;
ALTER TABLE `usdt_onchain_data` ADD `data_str` varchar(256) NOT NULL COMMENT '' AFTER `data_int64`;

----------------------------------------------------
--  `usdt_onchain_log`
----------------------------------------------------
ALTER TABLE `usdt_onchain_log` ADD `oclid` bigint unsigned NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_log` ADD `from` varchar(100) NOT NULL COMMENT '' AFTER `oclid`;
ALTER TABLE `usdt_onchain_log` ADD `to` varchar(100) NOT NULL COMMENT '' AFTER `from`;
ALTER TABLE `usdt_onchain_log` ADD `tx` varchar(100) NOT NULL COMMENT '' AFTER `to`;
ALTER TABLE `usdt_onchain_log` ADD `status` varchar(100) NOT NULL COMMENT '' AFTER `tx`;
ALTER TABLE `usdt_onchain_log` ADD `pushed` varchar(100) NOT NULL COMMENT '' AFTER `status`;
ALTER TABLE `usdt_onchain_log` ADD `signedTx` text NOT NULL COMMENT '' AFTER `pushed`;
ALTER TABLE `usdt_onchain_log` ADD `amount_integer` bigint NOT NULL COMMENT 'amount integer part' AFTER `signedTx`;
ALTER TABLE `usdt_onchain_log` ADD `ctime` bigint NOT NULL COMMENT 'created time' AFTER `amount_integer`;

----------------------------------------------------
--  `usdt_onchain_sync_pos`
----------------------------------------------------
ALTER TABLE `usdt_onchain_sync_pos` ADD `address` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_sync_pos` ADD `page` int unsigned NOT NULL COMMENT '' AFTER `address`;
ALTER TABLE `usdt_onchain_sync_pos` ADD `tx_id` varchar(128) NOT NULL COMMENT '' AFTER `page`;

----------------------------------------------------
--  `usdt_onchain_transaction`
----------------------------------------------------
ALTER TABLE `usdt_onchain_transaction` ADD `tx_id` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_transaction` ADD `uaid` bigint unsigned NOT NULL COMMENT 'usdt account id' AFTER `tx_id`;
ALTER TABLE `usdt_onchain_transaction` ADD `type` tinyint unsigned NOT NULL COMMENT 'Transaction type' AFTER `uaid`;
ALTER TABLE `usdt_onchain_transaction` ADD `property_id` int unsigned NOT NULL COMMENT '' AFTER `type`;
ALTER TABLE `usdt_onchain_transaction` ADD `property_name` varchar(100) NOT NULL COMMENT '' AFTER `property_id`;
ALTER TABLE `usdt_onchain_transaction` ADD `tx_type` varchar(100) NOT NULL COMMENT '' AFTER `property_name`;
ALTER TABLE `usdt_onchain_transaction` ADD `tx_type_int` int NOT NULL COMMENT '' AFTER `tx_type`;
ALTER TABLE `usdt_onchain_transaction` ADD `amount_integer` bigint NOT NULL COMMENT 'amount integer part' AFTER `tx_type_int`;
ALTER TABLE `usdt_onchain_transaction` ADD `block` int unsigned NOT NULL COMMENT '' AFTER `amount_integer`;
ALTER TABLE `usdt_onchain_transaction` ADD `block_hash` varchar(100) NOT NULL COMMENT '' AFTER `block`;
ALTER TABLE `usdt_onchain_transaction` ADD `block_time` bigint NOT NULL COMMENT '' AFTER `block_hash`;
ALTER TABLE `usdt_onchain_transaction` ADD `confirmations` int unsigned NOT NULL COMMENT '' AFTER `block_time`;
ALTER TABLE `usdt_onchain_transaction` ADD `divisible` bool NOT NULL COMMENT '' AFTER `confirmations`;
ALTER TABLE `usdt_onchain_transaction` ADD `fee_amount_integer` bigint NOT NULL COMMENT 'fee amount integer part' AFTER `divisible`;
ALTER TABLE `usdt_onchain_transaction` ADD `is_mine` bool NOT NULL COMMENT '' AFTER `fee_amount_integer`;
ALTER TABLE `usdt_onchain_transaction` ADD `position_in_block` int unsigned NOT NULL COMMENT '' AFTER `is_mine`;
ALTER TABLE `usdt_onchain_transaction` ADD `referenceaddress` varchar(100) NOT NULL COMMENT '' AFTER `position_in_block`;
ALTER TABLE `usdt_onchain_transaction` ADD `sending_address` varchar(100) NOT NULL COMMENT '' AFTER `referenceaddress`;
ALTER TABLE `usdt_onchain_transaction` ADD `version` int NOT NULL COMMENT '' AFTER `sending_address`;
ALTER TABLE `usdt_onchain_transaction` ADD `mtime` bigint NOT NULL COMMENT '' AFTER `version`;

----------------------------------------------------
--  `usdt_prikey`
----------------------------------------------------
ALTER TABLE `usdt_prikey` ADD `pkid` bigint unsigned NOT NULL COMMENT '';
ALTER TABLE `usdt_prikey` ADD `pri` varchar(256) NOT NULL COMMENT '' AFTER `pkid`;
ALTER TABLE `usdt_prikey` ADD `address` varchar(100) NOT NULL COMMENT '' AFTER `pri`;

----------------------------------------------------
--  `usdt_sweep_log`
----------------------------------------------------
ALTER TABLE `usdt_sweep_log` ADD `id` bigint unsigned NOT NULL COMMENT 'id';
ALTER TABLE `usdt_sweep_log` ADD `uid` bigint unsigned NOT NULL COMMENT 'uid' AFTER `id`;
ALTER TABLE `usdt_sweep_log` ADD `ttype` int unsigned NOT NULL COMMENT 'Transaction type' AFTER `uid`;
ALTER TABLE `usdt_sweep_log` ADD `status` int unsigned NOT NULL COMMENT 'Transaction status' AFTER `ttype`;
ALTER TABLE `usdt_sweep_log` ADD `from` varchar(256) NOT NULL COMMENT 'sender address' AFTER `status`;
ALTER TABLE `usdt_sweep_log` ADD `to` varchar(256) NOT NULL COMMENT 'receiver address' AFTER `from`;
ALTER TABLE `usdt_sweep_log` ADD `txid` varchar(256) NOT NULL COMMENT 'Transaction id' AFTER `to`;
ALTER TABLE `usdt_sweep_log` ADD `amount_integer` bigint NOT NULL COMMENT 'amount integer' AFTER `txid`;
ALTER TABLE `usdt_sweep_log` ADD `fee_integer` bigint NOT NULL COMMENT 'fee integer' AFTER `amount_integer`;
ALTER TABLE `usdt_sweep_log` ADD `fee_onchain_integer` bigint NOT NULL COMMENT 'fee onchain' AFTER `fee_integer`;
ALTER TABLE `usdt_sweep_log` ADD `ctime` bigint NOT NULL COMMENT 'created time' AFTER `fee_onchain_integer`;
ALTER TABLE `usdt_sweep_log` ADD `utime` bigint NOT NULL COMMENT 'update time' AFTER `ctime`;
ALTER TABLE `usdt_sweep_log` ADD `step` varchar(64) NOT NULL COMMENT 'desc' AFTER `utime`;
ALTER TABLE `usdt_sweep_log` ADD `desc` varchar(256) NOT NULL COMMENT 'desc' AFTER `step`;

----------------------------------------------------
--  `usdt_transaction`
----------------------------------------------------
ALTER TABLE `usdt_transaction` ADD `tx_id` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `usdt_transaction` ADD `uaid` bigint unsigned NOT NULL COMMENT 'usdt account id' AFTER `tx_id`;
ALTER TABLE `usdt_transaction` ADD `type` tinyint unsigned NOT NULL COMMENT 'Transaction type' AFTER `uaid`;
ALTER TABLE `usdt_transaction` ADD `block_num` int unsigned NOT NULL COMMENT 'block num' AFTER `type`;
ALTER TABLE `usdt_transaction` ADD `status` int NOT NULL COMMENT 'Transaction status' AFTER `block_num`;
ALTER TABLE `usdt_transaction` ADD `payer` varchar(100) NOT NULL COMMENT 'Payer Account' AFTER `status`;
ALTER TABLE `usdt_transaction` ADD `receiver` varchar(100) NOT NULL COMMENT 'Receiver Account' AFTER `payer`;
ALTER TABLE `usdt_transaction` ADD `amount_integer` bigint NOT NULL COMMENT 'amount integer part' AFTER `receiver`;
ALTER TABLE `usdt_transaction` ADD `fee` varchar(100) NOT NULL COMMENT 'transaction Fee' AFTER `amount_integer`;
ALTER TABLE `usdt_transaction` ADD `memo` varchar(100) NOT NULL COMMENT 'transaction memo' AFTER `fee`;
ALTER TABLE `usdt_transaction` ADD `ctime` bigint NOT NULL COMMENT 'created time' AFTER `memo`;
ALTER TABLE `usdt_transaction` ADD `utime` bigint NOT NULL COMMENT 'update time' AFTER `ctime`;

----------------------------------------------------
--  `usdt_wealth_log`
----------------------------------------------------
ALTER TABLE `usdt_wealth_log` ADD `id` bigint unsigned NOT NULL COMMENT 'id';
ALTER TABLE `usdt_wealth_log` ADD `uid` bigint unsigned NOT NULL COMMENT 'uid' AFTER `id`;
ALTER TABLE `usdt_wealth_log` ADD `ttype` int unsigned NOT NULL COMMENT 'Transaction type' AFTER `uid`;
ALTER TABLE `usdt_wealth_log` ADD `status` int unsigned NOT NULL COMMENT 'Transaction status' AFTER `ttype`;
ALTER TABLE `usdt_wealth_log` ADD `from` varchar(256) NOT NULL COMMENT 'sender address' AFTER `status`;
ALTER TABLE `usdt_wealth_log` ADD `to` varchar(256) NOT NULL COMMENT 'receiver address' AFTER `from`;
ALTER TABLE `usdt_wealth_log` ADD `txid` varchar(256) NOT NULL COMMENT 'Transaction id' AFTER `to`;
ALTER TABLE `usdt_wealth_log` ADD `amount_integer` bigint NOT NULL COMMENT 'amount integer' AFTER `txid`;
ALTER TABLE `usdt_wealth_log` ADD `fee_integer` bigint NOT NULL COMMENT 'fee integer' AFTER `amount_integer`;
ALTER TABLE `usdt_wealth_log` ADD `fee_usdt_integer` bigint NOT NULL COMMENT 'fee usdt integer' AFTER `fee_integer`;
ALTER TABLE `usdt_wealth_log` ADD `fee_onchain_integer` bigint NOT NULL COMMENT 'fee onchain' AFTER `fee_usdt_integer`;
ALTER TABLE `usdt_wealth_log` ADD `ctime` bigint NOT NULL COMMENT 'created time' AFTER `fee_onchain_integer`;
ALTER TABLE `usdt_wealth_log` ADD `utime` bigint NOT NULL COMMENT 'update time' AFTER `ctime`;
ALTER TABLE `usdt_wealth_log` ADD `step` varchar(64) NOT NULL COMMENT 'desc' AFTER `utime`;
ALTER TABLE `usdt_wealth_log` ADD `desc` varchar(256) NOT NULL COMMENT 'desc' AFTER `step`;
ALTER TABLE `usdt_wealth_log` ADD `sign` varchar(256) NOT NULL COMMENT 'sign' AFTER `desc`;
ALTER TABLE `usdt_wealth_log` ADD `memo` varchar(256) NOT NULL COMMENT 'memo' AFTER `sign`;

