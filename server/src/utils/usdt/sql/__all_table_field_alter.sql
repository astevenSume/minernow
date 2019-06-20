----------------------------------------------------
--  `market_prices`
----------------------------------------------------
ALTER TABLE `market_prices` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'market otc price record id';
ALTER TABLE `market_prices` CHANGE `market` `market` tinyint NOT NULL COMMENT 'market';
ALTER TABLE `market_prices` CHANGE `currency` `currency` varchar(100) NOT NULL COMMENT 'currency';
ALTER TABLE `market_prices` CHANGE `trade_method` `trade_method` tinyint NOT NULL COMMENT 'trade method';
ALTER TABLE `market_prices` CHANGE `pow_price` `pow_price` bigint unsigned NOT NULL COMMENT 'pow price';
ALTER TABLE `market_prices` CHANGE `pow` `pow` int NOT NULL DEFAULT 4 COMMENT 'pow_price = (real price) * pow';
ALTER TABLE `market_prices` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'created time';

----------------------------------------------------
--  `prices`
----------------------------------------------------
ALTER TABLE `prices` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'otc price record id';
ALTER TABLE `prices` CHANGE `currency` `currency` varchar(100) NOT NULL COMMENT 'currency';
ALTER TABLE `prices` CHANGE `pow_price` `pow_price` bigint unsigned NOT NULL COMMENT 'pow price';
ALTER TABLE `prices` CHANGE `pow` `pow` int NOT NULL DEFAULT 4 COMMENT 'pow_price = (real price) * pow';
ALTER TABLE `prices` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'created time';

----------------------------------------------------
--  `usdt_account`
----------------------------------------------------
ALTER TABLE `usdt_account` CHANGE `uaid` `uaid` bigint unsigned NOT NULL COMMENT '';
ALTER TABLE `usdt_account` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT '';
ALTER TABLE `usdt_account` CHANGE `status` `status` tinyint unsigned NOT NULL COMMENT 'status 0-locked 1-working';
ALTER TABLE `usdt_account` CHANGE `available_integer` `available_integer` bigint NOT NULL COMMENT 'available amount integer part';
ALTER TABLE `usdt_account` CHANGE `frozen_integer` `frozen_integer` bigint NOT NULL COMMENT 'frozen amount integer part';
ALTER TABLE `usdt_account` CHANGE `transfer_frozen_integer` `transfer_frozen_integer` bigint NOT NULL COMMENT 'frozen amount integer part';
ALTER TABLE `usdt_account` CHANGE `mortgaged_integer` `mortgaged_integer` bigint NOT NULL COMMENT 'mortgaged amount integer part';
ALTER TABLE `usdt_account` CHANGE `btc_available_integer` `btc_available_integer` bigint NOT NULL COMMENT 'available btc amount integer part';
ALTER TABLE `usdt_account` CHANGE `btc_frozen_integer` `btc_frozen_integer` bigint NOT NULL COMMENT 'frozen btc amount integer part';
ALTER TABLE `usdt_account` CHANGE `btc_mortgaged_integer` `btc_mortgaged_integer` bigint NOT NULL COMMENT 'mortgaged btc amount integer part';
ALTER TABLE `usdt_account` CHANGE `waiting_cash_sweep_integer` `waiting_cash_sweep_integer` bigint NOT NULL COMMENT 'waiting cash sweep integer part';
ALTER TABLE `usdt_account` CHANGE `cash_sweep_integer` `cash_sweep_integer` bigint NOT NULL COMMENT 'cash sweep integer part';
ALTER TABLE `usdt_account` CHANGE `owned_by_platform_integer` `owned_by_platform_integer` bigint NOT NULL COMMENT 'owned by platform integer part';
ALTER TABLE `usdt_account` CHANGE `sweep_status` `sweep_status` tinyint unsigned NOT NULL COMMENT 'sweep status';
ALTER TABLE `usdt_account` CHANGE `pkid` `pkid` bigint unsigned NOT NULL COMMENT 'usdt private key id';
ALTER TABLE `usdt_account` CHANGE `address` `address` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `usdt_account` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'created time';
ALTER TABLE `usdt_account` CHANGE `mtime` `mtime` bigint NOT NULL COMMENT 'update time';
ALTER TABLE `usdt_account` CHANGE `sign` `sign` varchar(256) NOT NULL COMMENT 'sign';

----------------------------------------------------
--  `usdt_onchain_balance`
----------------------------------------------------
ALTER TABLE `usdt_onchain_balance` CHANGE `address` `address` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_balance` CHANGE `property_id` `property_id` int unsigned NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_balance` CHANGE `pending_pos` `pending_pos` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_balance` CHANGE `reserved` `reserved` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_balance` CHANGE `divisible` `divisible` bool NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_balance` CHANGE `amount_integer` `amount_integer` bigint NOT NULL COMMENT 'amount integer part';
ALTER TABLE `usdt_onchain_balance` CHANGE `frozen` `frozen` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_balance` CHANGE `pending_neg` `pending_neg` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_balance` CHANGE `mtime` `mtime` bigint NOT NULL COMMENT '';

----------------------------------------------------
--  `usdt_onchain_data`
----------------------------------------------------
ALTER TABLE `usdt_onchain_data` CHANGE `address` `address` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_data` CHANGE `attr_type` `attr_type` int unsigned NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_data` CHANGE `data_int64` `data_int64` bigint NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_data` CHANGE `data_str` `data_str` varchar(256) NOT NULL COMMENT '';

----------------------------------------------------
--  `usdt_onchain_log`
----------------------------------------------------
ALTER TABLE `usdt_onchain_log` CHANGE `oclid` `oclid` bigint unsigned NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_log` CHANGE `from` `from` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_log` CHANGE `to` `to` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_log` CHANGE `tx` `tx` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_log` CHANGE `status` `status` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_log` CHANGE `pushed` `pushed` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_log` CHANGE `signedTx` `signedTx` text NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_log` CHANGE `amount_integer` `amount_integer` bigint NOT NULL COMMENT 'amount integer part';
ALTER TABLE `usdt_onchain_log` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'created time';

----------------------------------------------------
--  `usdt_onchain_sync_pos`
----------------------------------------------------
ALTER TABLE `usdt_onchain_sync_pos` CHANGE `address` `address` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_sync_pos` CHANGE `page` `page` int unsigned NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_sync_pos` CHANGE `tx_id` `tx_id` varchar(128) NOT NULL COMMENT '';

----------------------------------------------------
--  `usdt_onchain_transaction`
----------------------------------------------------
ALTER TABLE `usdt_onchain_transaction` CHANGE `tx_id` `tx_id` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_transaction` CHANGE `uaid` `uaid` bigint unsigned NOT NULL COMMENT 'usdt account id';
ALTER TABLE `usdt_onchain_transaction` CHANGE `type` `type` tinyint unsigned NOT NULL COMMENT 'Transaction type';
ALTER TABLE `usdt_onchain_transaction` CHANGE `property_id` `property_id` int unsigned NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_transaction` CHANGE `property_name` `property_name` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_transaction` CHANGE `tx_type` `tx_type` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_transaction` CHANGE `tx_type_int` `tx_type_int` int NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_transaction` CHANGE `amount_integer` `amount_integer` bigint NOT NULL COMMENT 'amount integer part';
ALTER TABLE `usdt_onchain_transaction` CHANGE `block` `block` int unsigned NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_transaction` CHANGE `block_hash` `block_hash` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_transaction` CHANGE `block_time` `block_time` bigint NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_transaction` CHANGE `confirmations` `confirmations` int unsigned NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_transaction` CHANGE `divisible` `divisible` bool NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_transaction` CHANGE `fee_amount_integer` `fee_amount_integer` bigint NOT NULL COMMENT 'fee amount integer part';
ALTER TABLE `usdt_onchain_transaction` CHANGE `is_mine` `is_mine` bool NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_transaction` CHANGE `position_in_block` `position_in_block` int unsigned NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_transaction` CHANGE `referenceaddress` `referenceaddress` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_transaction` CHANGE `sending_address` `sending_address` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_transaction` CHANGE `version` `version` int NOT NULL COMMENT '';
ALTER TABLE `usdt_onchain_transaction` CHANGE `mtime` `mtime` bigint NOT NULL COMMENT '';

----------------------------------------------------
--  `usdt_prikey`
----------------------------------------------------
ALTER TABLE `usdt_prikey` CHANGE `pkid` `pkid` bigint unsigned NOT NULL COMMENT '';
ALTER TABLE `usdt_prikey` CHANGE `pri` `pri` varchar(256) NOT NULL COMMENT '';
ALTER TABLE `usdt_prikey` CHANGE `address` `address` varchar(100) NOT NULL COMMENT '';

----------------------------------------------------
--  `usdt_sweep_log`
----------------------------------------------------
ALTER TABLE `usdt_sweep_log` CHANGE `id` `id` bigint unsigned NOT NULL COMMENT 'id';
ALTER TABLE `usdt_sweep_log` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'uid';
ALTER TABLE `usdt_sweep_log` CHANGE `ttype` `ttype` int unsigned NOT NULL COMMENT 'Transaction type';
ALTER TABLE `usdt_sweep_log` CHANGE `status` `status` int unsigned NOT NULL COMMENT 'Transaction status';
ALTER TABLE `usdt_sweep_log` CHANGE `from` `from` varchar(256) NOT NULL COMMENT 'sender address';
ALTER TABLE `usdt_sweep_log` CHANGE `to` `to` varchar(256) NOT NULL COMMENT 'receiver address';
ALTER TABLE `usdt_sweep_log` CHANGE `txid` `txid` varchar(256) NOT NULL COMMENT 'Transaction id';
ALTER TABLE `usdt_sweep_log` CHANGE `amount_integer` `amount_integer` bigint NOT NULL COMMENT 'amount integer';
ALTER TABLE `usdt_sweep_log` CHANGE `fee_integer` `fee_integer` bigint NOT NULL COMMENT 'fee integer';
ALTER TABLE `usdt_sweep_log` CHANGE `fee_onchain_integer` `fee_onchain_integer` bigint NOT NULL COMMENT 'fee onchain';
ALTER TABLE `usdt_sweep_log` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'created time';
ALTER TABLE `usdt_sweep_log` CHANGE `utime` `utime` bigint NOT NULL COMMENT 'update time';
ALTER TABLE `usdt_sweep_log` CHANGE `step` `step` varchar(64) NOT NULL COMMENT 'desc';
ALTER TABLE `usdt_sweep_log` CHANGE `desc` `desc` varchar(256) NOT NULL COMMENT 'desc';

----------------------------------------------------
--  `usdt_transaction`
----------------------------------------------------
ALTER TABLE `usdt_transaction` CHANGE `tx_id` `tx_id` varchar(100) NOT NULL COMMENT '';
ALTER TABLE `usdt_transaction` CHANGE `uaid` `uaid` bigint unsigned NOT NULL COMMENT 'usdt account id';
ALTER TABLE `usdt_transaction` CHANGE `type` `type` tinyint unsigned NOT NULL COMMENT 'Transaction type';
ALTER TABLE `usdt_transaction` CHANGE `block_num` `block_num` int unsigned NOT NULL COMMENT 'block num';
ALTER TABLE `usdt_transaction` CHANGE `status` `status` int NOT NULL COMMENT 'Transaction status';
ALTER TABLE `usdt_transaction` CHANGE `payer` `payer` varchar(100) NOT NULL COMMENT 'Payer Account';
ALTER TABLE `usdt_transaction` CHANGE `receiver` `receiver` varchar(100) NOT NULL COMMENT 'Receiver Account';
ALTER TABLE `usdt_transaction` CHANGE `amount_integer` `amount_integer` bigint NOT NULL COMMENT 'amount integer part';
ALTER TABLE `usdt_transaction` CHANGE `fee` `fee` varchar(100) NOT NULL COMMENT 'transaction Fee';
ALTER TABLE `usdt_transaction` CHANGE `memo` `memo` varchar(100) NOT NULL COMMENT 'transaction memo';
ALTER TABLE `usdt_transaction` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'created time';
ALTER TABLE `usdt_transaction` CHANGE `utime` `utime` bigint NOT NULL COMMENT 'update time';

----------------------------------------------------
--  `usdt_wealth_log`
----------------------------------------------------
ALTER TABLE `usdt_wealth_log` CHANGE `id` `id` bigint unsigned NOT NULL COMMENT 'id';
ALTER TABLE `usdt_wealth_log` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'uid';
ALTER TABLE `usdt_wealth_log` CHANGE `ttype` `ttype` int unsigned NOT NULL COMMENT 'Transaction type';
ALTER TABLE `usdt_wealth_log` CHANGE `status` `status` int unsigned NOT NULL COMMENT 'Transaction status';
ALTER TABLE `usdt_wealth_log` CHANGE `from` `from` varchar(256) NOT NULL COMMENT 'sender address';
ALTER TABLE `usdt_wealth_log` CHANGE `to` `to` varchar(256) NOT NULL COMMENT 'receiver address';
ALTER TABLE `usdt_wealth_log` CHANGE `txid` `txid` varchar(256) NOT NULL COMMENT 'Transaction id';
ALTER TABLE `usdt_wealth_log` CHANGE `amount_integer` `amount_integer` bigint NOT NULL COMMENT 'amount integer';
ALTER TABLE `usdt_wealth_log` CHANGE `fee_integer` `fee_integer` bigint NOT NULL COMMENT 'fee integer';
ALTER TABLE `usdt_wealth_log` CHANGE `fee_usdt_integer` `fee_usdt_integer` bigint NOT NULL COMMENT 'fee usdt integer';
ALTER TABLE `usdt_wealth_log` CHANGE `fee_onchain_integer` `fee_onchain_integer` bigint NOT NULL COMMENT 'fee onchain';
ALTER TABLE `usdt_wealth_log` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'created time';
ALTER TABLE `usdt_wealth_log` CHANGE `utime` `utime` bigint NOT NULL COMMENT 'update time';
ALTER TABLE `usdt_wealth_log` CHANGE `step` `step` varchar(64) NOT NULL COMMENT 'desc';
ALTER TABLE `usdt_wealth_log` CHANGE `desc` `desc` varchar(256) NOT NULL COMMENT 'desc';
ALTER TABLE `usdt_wealth_log` CHANGE `sign` `sign` varchar(256) NOT NULL COMMENT 'sign';
ALTER TABLE `usdt_wealth_log` CHANGE `memo` `memo` varchar(256) NOT NULL COMMENT 'memo';

