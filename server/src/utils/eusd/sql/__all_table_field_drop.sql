----------------------------------------------------
--  `eos_account`
----------------------------------------------------
ALTER TABLE `eos_account` DROP `id`;
ALTER TABLE `eos_account` DROP `uid`;
ALTER TABLE `eos_account` DROP `account`;
ALTER TABLE `eos_account` DROP `balance`;
ALTER TABLE `eos_account` DROP `status`;
ALTER TABLE `eos_account` DROP `ctime`;
ALTER TABLE `eos_account` DROP `utime`;

----------------------------------------------------
--  `eos_otc`
----------------------------------------------------
ALTER TABLE `eos_otc` DROP `uid`;
ALTER TABLE `eos_otc` DROP `account`;
ALTER TABLE `eos_otc` DROP `status`;
ALTER TABLE `eos_otc` DROP `available`;
ALTER TABLE `eos_otc` DROP `trade`;
ALTER TABLE `eos_otc` DROP `transfer`;
ALTER TABLE `eos_otc` DROP `sell_state`;
ALTER TABLE `eos_otc` DROP `sell_pay_type`;
ALTER TABLE `eos_otc` DROP `sell_able`;
ALTER TABLE `eos_otc` DROP `sell_rmb_day`;
ALTER TABLE `eos_otc` DROP `sell_rmb_today`;
ALTER TABLE `eos_otc` DROP `sell_rmb_lower_limit`;
ALTER TABLE `eos_otc` DROP `sell_utime`;
ALTER TABLE `eos_otc` DROP `buy_able`;
ALTER TABLE `eos_otc` DROP `buy_rmb_day`;
ALTER TABLE `eos_otc` DROP `buy_rmb_today`;
ALTER TABLE `eos_otc` DROP `buy_rmb_lower_limit`;
ALTER TABLE `eos_otc` DROP `buy_utime`;
ALTER TABLE `eos_otc` DROP `buy_state`;
ALTER TABLE `eos_otc` DROP `ctime`;
ALTER TABLE `eos_otc` DROP `utime`;

----------------------------------------------------
--  `eos_transaction`
----------------------------------------------------
ALTER TABLE `eos_transaction` DROP `id`;
ALTER TABLE `eos_transaction` DROP `type`;
ALTER TABLE `eos_transaction` DROP `transaction_id`;
ALTER TABLE `eos_transaction` DROP `block_num`;
ALTER TABLE `eos_transaction` DROP `status`;
ALTER TABLE `eos_transaction` DROP `payer`;
ALTER TABLE `eos_transaction` DROP `receiver`;
ALTER TABLE `eos_transaction` DROP `quantity`;
ALTER TABLE `eos_transaction` DROP `memo`;
ALTER TABLE `eos_transaction` DROP `ctime`;
ALTER TABLE `eos_transaction` DROP `utime`;

----------------------------------------------------
--  `eos_transaction_info`
----------------------------------------------------
ALTER TABLE `eos_transaction_info` DROP `id`;
ALTER TABLE `eos_transaction_info` DROP `transaction_id`;
ALTER TABLE `eos_transaction_info` DROP `block_num`;
ALTER TABLE `eos_transaction_info` DROP `ctime`;
ALTER TABLE `eos_transaction_info` DROP `processed`;

----------------------------------------------------
--  `eos_tx_log`
----------------------------------------------------
ALTER TABLE `eos_tx_log` DROP `id`;
ALTER TABLE `eos_tx_log` DROP `from`;
ALTER TABLE `eos_tx_log` DROP `from_uid`;
ALTER TABLE `eos_tx_log` DROP `to`;
ALTER TABLE `eos_tx_log` DROP `to_uid`;
ALTER TABLE `eos_tx_log` DROP `quantity`;
ALTER TABLE `eos_tx_log` DROP `status`;
ALTER TABLE `eos_tx_log` DROP `log_ids`;
ALTER TABLE `eos_tx_log` DROP `ctime`;
ALTER TABLE `eos_tx_log` DROP `txid`;
ALTER TABLE `eos_tx_log` DROP `order_id`;
ALTER TABLE `eos_tx_log` DROP `utime`;
ALTER TABLE `eos_tx_log` DROP `sign`;
ALTER TABLE `eos_tx_log` DROP `delay_deal`;
ALTER TABLE `eos_tx_log` DROP `retry`;
ALTER TABLE `eos_tx_log` DROP `memo`;

----------------------------------------------------
--  `eos_use_log`
----------------------------------------------------
ALTER TABLE `eos_use_log` DROP `id`;
ALTER TABLE `eos_use_log` DROP `type`;
ALTER TABLE `eos_use_log` DROP `tid`;
ALTER TABLE `eos_use_log` DROP `status`;
ALTER TABLE `eos_use_log` DROP `tid_recover`;
ALTER TABLE `eos_use_log` DROP `quantity_num`;

----------------------------------------------------
--  `eos_wealth`
----------------------------------------------------
ALTER TABLE `eos_wealth` DROP `uid`;
ALTER TABLE `eos_wealth` DROP `status`;
ALTER TABLE `eos_wealth` DROP `account`;
ALTER TABLE `eos_wealth` DROP `balance`;
ALTER TABLE `eos_wealth` DROP `available`;
ALTER TABLE `eos_wealth` DROP `game`;
ALTER TABLE `eos_wealth` DROP `trade`;
ALTER TABLE `eos_wealth` DROP `transfer`;
ALTER TABLE `eos_wealth` DROP `transfer_game`;
ALTER TABLE `eos_wealth` DROP `is_exchanger`;
ALTER TABLE `eos_wealth` DROP `ctime`;
ALTER TABLE `eos_wealth` DROP `utime`;

----------------------------------------------------
--  `eos_wealth_log`
----------------------------------------------------
ALTER TABLE `eos_wealth_log` DROP `id`;
ALTER TABLE `eos_wealth_log` DROP `uid`;
ALTER TABLE `eos_wealth_log` DROP `uid2`;
ALTER TABLE `eos_wealth_log` DROP `ttype`;
ALTER TABLE `eos_wealth_log` DROP `status`;
ALTER TABLE `eos_wealth_log` DROP `txid`;
ALTER TABLE `eos_wealth_log` DROP `quantity`;
ALTER TABLE `eos_wealth_log` DROP `ctime`;

----------------------------------------------------
--  `eusd_retire`
----------------------------------------------------
ALTER TABLE `eusd_retire` DROP `id`;
ALTER TABLE `eusd_retire` DROP `from`;
ALTER TABLE `eusd_retire` DROP `from_uid`;
ALTER TABLE `eusd_retire` DROP `quantity`;
ALTER TABLE `eusd_retire` DROP `status`;
ALTER TABLE `eusd_retire` DROP `ctime`;

----------------------------------------------------
--  `platform_user`
----------------------------------------------------
ALTER TABLE `platform_user` DROP `uid`;
ALTER TABLE `platform_user` DROP `pid`;
ALTER TABLE `platform_user` DROP `status`;
ALTER TABLE `platform_user` DROP `ctime`;

----------------------------------------------------
--  `platform_user_cate`
----------------------------------------------------
ALTER TABLE `platform_user_cate` DROP `id`;
ALTER TABLE `platform_user_cate` DROP `name`;
ALTER TABLE `platform_user_cate` DROP `dividend`;
ALTER TABLE `platform_user_cate` DROP `ctime`;

