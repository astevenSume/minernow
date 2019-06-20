----------------------------------------------------
--  `market_prices`
----------------------------------------------------
ALTER TABLE `market_prices` DROP `id`;
ALTER TABLE `market_prices` DROP `market`;
ALTER TABLE `market_prices` DROP `currency`;
ALTER TABLE `market_prices` DROP `trade_method`;
ALTER TABLE `market_prices` DROP `pow_price`;
ALTER TABLE `market_prices` DROP `pow`;
ALTER TABLE `market_prices` DROP `ctime`;

----------------------------------------------------
--  `prices`
----------------------------------------------------
ALTER TABLE `prices` DROP `id`;
ALTER TABLE `prices` DROP `currency`;
ALTER TABLE `prices` DROP `pow_price`;
ALTER TABLE `prices` DROP `pow`;
ALTER TABLE `prices` DROP `ctime`;

----------------------------------------------------
--  `usdt_account`
----------------------------------------------------
ALTER TABLE `usdt_account` DROP `uaid`;
ALTER TABLE `usdt_account` DROP `uid`;
ALTER TABLE `usdt_account` DROP `status`;
ALTER TABLE `usdt_account` DROP `available_integer`;
ALTER TABLE `usdt_account` DROP `frozen_integer`;
ALTER TABLE `usdt_account` DROP `transfer_frozen_integer`;
ALTER TABLE `usdt_account` DROP `mortgaged_integer`;
ALTER TABLE `usdt_account` DROP `btc_available_integer`;
ALTER TABLE `usdt_account` DROP `btc_frozen_integer`;
ALTER TABLE `usdt_account` DROP `btc_mortgaged_integer`;
ALTER TABLE `usdt_account` DROP `waiting_cash_sweep_integer`;
ALTER TABLE `usdt_account` DROP `cash_sweep_integer`;
ALTER TABLE `usdt_account` DROP `owned_by_platform_integer`;
ALTER TABLE `usdt_account` DROP `sweep_status`;
ALTER TABLE `usdt_account` DROP `pkid`;
ALTER TABLE `usdt_account` DROP `address`;
ALTER TABLE `usdt_account` DROP `ctime`;
ALTER TABLE `usdt_account` DROP `mtime`;
ALTER TABLE `usdt_account` DROP `sign`;

----------------------------------------------------
--  `usdt_onchain_balance`
----------------------------------------------------
ALTER TABLE `usdt_onchain_balance` DROP `address`;
ALTER TABLE `usdt_onchain_balance` DROP `property_id`;
ALTER TABLE `usdt_onchain_balance` DROP `pending_pos`;
ALTER TABLE `usdt_onchain_balance` DROP `reserved`;
ALTER TABLE `usdt_onchain_balance` DROP `divisible`;
ALTER TABLE `usdt_onchain_balance` DROP `amount_integer`;
ALTER TABLE `usdt_onchain_balance` DROP `frozen`;
ALTER TABLE `usdt_onchain_balance` DROP `pending_neg`;
ALTER TABLE `usdt_onchain_balance` DROP `mtime`;

----------------------------------------------------
--  `usdt_onchain_data`
----------------------------------------------------
ALTER TABLE `usdt_onchain_data` DROP `address`;
ALTER TABLE `usdt_onchain_data` DROP `attr_type`;
ALTER TABLE `usdt_onchain_data` DROP `data_int64`;
ALTER TABLE `usdt_onchain_data` DROP `data_str`;

----------------------------------------------------
--  `usdt_onchain_log`
----------------------------------------------------
ALTER TABLE `usdt_onchain_log` DROP `oclid`;
ALTER TABLE `usdt_onchain_log` DROP `from`;
ALTER TABLE `usdt_onchain_log` DROP `to`;
ALTER TABLE `usdt_onchain_log` DROP `tx`;
ALTER TABLE `usdt_onchain_log` DROP `status`;
ALTER TABLE `usdt_onchain_log` DROP `pushed`;
ALTER TABLE `usdt_onchain_log` DROP `signedTx`;
ALTER TABLE `usdt_onchain_log` DROP `amount_integer`;
ALTER TABLE `usdt_onchain_log` DROP `ctime`;

----------------------------------------------------
--  `usdt_onchain_sync_pos`
----------------------------------------------------
ALTER TABLE `usdt_onchain_sync_pos` DROP `address`;
ALTER TABLE `usdt_onchain_sync_pos` DROP `page`;
ALTER TABLE `usdt_onchain_sync_pos` DROP `tx_id`;

----------------------------------------------------
--  `usdt_onchain_transaction`
----------------------------------------------------
ALTER TABLE `usdt_onchain_transaction` DROP `tx_id`;
ALTER TABLE `usdt_onchain_transaction` DROP `uaid`;
ALTER TABLE `usdt_onchain_transaction` DROP `type`;
ALTER TABLE `usdt_onchain_transaction` DROP `property_id`;
ALTER TABLE `usdt_onchain_transaction` DROP `property_name`;
ALTER TABLE `usdt_onchain_transaction` DROP `tx_type`;
ALTER TABLE `usdt_onchain_transaction` DROP `tx_type_int`;
ALTER TABLE `usdt_onchain_transaction` DROP `amount_integer`;
ALTER TABLE `usdt_onchain_transaction` DROP `block`;
ALTER TABLE `usdt_onchain_transaction` DROP `block_hash`;
ALTER TABLE `usdt_onchain_transaction` DROP `block_time`;
ALTER TABLE `usdt_onchain_transaction` DROP `confirmations`;
ALTER TABLE `usdt_onchain_transaction` DROP `divisible`;
ALTER TABLE `usdt_onchain_transaction` DROP `fee_amount_integer`;
ALTER TABLE `usdt_onchain_transaction` DROP `is_mine`;
ALTER TABLE `usdt_onchain_transaction` DROP `position_in_block`;
ALTER TABLE `usdt_onchain_transaction` DROP `referenceaddress`;
ALTER TABLE `usdt_onchain_transaction` DROP `sending_address`;
ALTER TABLE `usdt_onchain_transaction` DROP `version`;
ALTER TABLE `usdt_onchain_transaction` DROP `mtime`;

----------------------------------------------------
--  `usdt_prikey`
----------------------------------------------------
ALTER TABLE `usdt_prikey` DROP `pkid`;
ALTER TABLE `usdt_prikey` DROP `pri`;
ALTER TABLE `usdt_prikey` DROP `address`;

----------------------------------------------------
--  `usdt_sweep_log`
----------------------------------------------------
ALTER TABLE `usdt_sweep_log` DROP `id`;
ALTER TABLE `usdt_sweep_log` DROP `uid`;
ALTER TABLE `usdt_sweep_log` DROP `ttype`;
ALTER TABLE `usdt_sweep_log` DROP `status`;
ALTER TABLE `usdt_sweep_log` DROP `from`;
ALTER TABLE `usdt_sweep_log` DROP `to`;
ALTER TABLE `usdt_sweep_log` DROP `txid`;
ALTER TABLE `usdt_sweep_log` DROP `amount_integer`;
ALTER TABLE `usdt_sweep_log` DROP `fee_integer`;
ALTER TABLE `usdt_sweep_log` DROP `fee_onchain_integer`;
ALTER TABLE `usdt_sweep_log` DROP `ctime`;
ALTER TABLE `usdt_sweep_log` DROP `utime`;
ALTER TABLE `usdt_sweep_log` DROP `step`;
ALTER TABLE `usdt_sweep_log` DROP `desc`;

----------------------------------------------------
--  `usdt_transaction`
----------------------------------------------------
ALTER TABLE `usdt_transaction` DROP `tx_id`;
ALTER TABLE `usdt_transaction` DROP `uaid`;
ALTER TABLE `usdt_transaction` DROP `type`;
ALTER TABLE `usdt_transaction` DROP `block_num`;
ALTER TABLE `usdt_transaction` DROP `status`;
ALTER TABLE `usdt_transaction` DROP `payer`;
ALTER TABLE `usdt_transaction` DROP `receiver`;
ALTER TABLE `usdt_transaction` DROP `amount_integer`;
ALTER TABLE `usdt_transaction` DROP `fee`;
ALTER TABLE `usdt_transaction` DROP `memo`;
ALTER TABLE `usdt_transaction` DROP `ctime`;
ALTER TABLE `usdt_transaction` DROP `utime`;

----------------------------------------------------
--  `usdt_wealth_log`
----------------------------------------------------
ALTER TABLE `usdt_wealth_log` DROP `id`;
ALTER TABLE `usdt_wealth_log` DROP `uid`;
ALTER TABLE `usdt_wealth_log` DROP `ttype`;
ALTER TABLE `usdt_wealth_log` DROP `status`;
ALTER TABLE `usdt_wealth_log` DROP `from`;
ALTER TABLE `usdt_wealth_log` DROP `to`;
ALTER TABLE `usdt_wealth_log` DROP `txid`;
ALTER TABLE `usdt_wealth_log` DROP `amount_integer`;
ALTER TABLE `usdt_wealth_log` DROP `fee_integer`;
ALTER TABLE `usdt_wealth_log` DROP `fee_usdt_integer`;
ALTER TABLE `usdt_wealth_log` DROP `fee_onchain_integer`;
ALTER TABLE `usdt_wealth_log` DROP `ctime`;
ALTER TABLE `usdt_wealth_log` DROP `utime`;
ALTER TABLE `usdt_wealth_log` DROP `step`;
ALTER TABLE `usdt_wealth_log` DROP `desc`;
ALTER TABLE `usdt_wealth_log` DROP `sign`;
ALTER TABLE `usdt_wealth_log` DROP `memo`;

