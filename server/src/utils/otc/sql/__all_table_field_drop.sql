----------------------------------------------------
--  `appeal`
----------------------------------------------------
ALTER TABLE `appeal` DROP `id`;
ALTER TABLE `appeal` DROP `type`;
ALTER TABLE `appeal` DROP `user_id`;
ALTER TABLE `appeal` DROP `admin_id`;
ALTER TABLE `appeal` DROP `order_id`;
ALTER TABLE `appeal` DROP `context`;
ALTER TABLE `appeal` DROP `status`;
ALTER TABLE `appeal` DROP `ctime`;
ALTER TABLE `appeal` DROP `utime`;
ALTER TABLE `appeal` DROP `wechat`;

----------------------------------------------------
--  `appeal_deal_log`
----------------------------------------------------
ALTER TABLE `appeal_deal_log` DROP `id`;
ALTER TABLE `appeal_deal_log` DROP `appeal_id`;
ALTER TABLE `appeal_deal_log` DROP `admin_id`;
ALTER TABLE `appeal_deal_log` DROP `order_id`;
ALTER TABLE `appeal_deal_log` DROP `action`;
ALTER TABLE `appeal_deal_log` DROP `ctime`;

----------------------------------------------------
--  `commission_calc`
----------------------------------------------------
ALTER TABLE `commission_calc` DROP `id`;
ALTER TABLE `commission_calc` DROP `start`;
ALTER TABLE `commission_calc` DROP `end`;
ALTER TABLE `commission_calc` DROP `calc_start`;
ALTER TABLE `commission_calc` DROP `calc_end`;
ALTER TABLE `commission_calc` DROP `status`;
ALTER TABLE `commission_calc` DROP `desc`;

----------------------------------------------------
--  `commission_distribute`
----------------------------------------------------
ALTER TABLE `commission_distribute` DROP `id`;
ALTER TABLE `commission_distribute` DROP `start`;
ALTER TABLE `commission_distribute` DROP `end`;
ALTER TABLE `commission_distribute` DROP `distribute_start`;
ALTER TABLE `commission_distribute` DROP `distribute_end`;
ALTER TABLE `commission_distribute` DROP `status`;
ALTER TABLE `commission_distribute` DROP `desc`;

----------------------------------------------------
--  `commission_stat`
----------------------------------------------------
ALTER TABLE `commission_stat` DROP `ctime`;
ALTER TABLE `commission_stat` DROP `tax_integer`;
ALTER TABLE `commission_stat` DROP `tax_decimals`;
ALTER TABLE `commission_stat` DROP `channel_integer`;
ALTER TABLE `commission_stat` DROP `channel_decimals`;
ALTER TABLE `commission_stat` DROP `commission_integer`;
ALTER TABLE `commission_stat` DROP `commission_decimals`;
ALTER TABLE `commission_stat` DROP `profit_integer`;
ALTER TABLE `commission_stat` DROP `profit_decimals`;
ALTER TABLE `commission_stat` DROP `mtime`;
ALTER TABLE `commission_stat` DROP `status`;

----------------------------------------------------
--  `eos_otc_report`
----------------------------------------------------
ALTER TABLE `eos_otc_report` DROP `id`;
ALTER TABLE `eos_otc_report` DROP `uid`;
ALTER TABLE `eos_otc_report` DROP `total_order_num`;
ALTER TABLE `eos_otc_report` DROP `success_order_num`;
ALTER TABLE `eos_otc_report` DROP `fail_order_num`;
ALTER TABLE `eos_otc_report` DROP `buy_eusd_num`;
ALTER TABLE `eos_otc_report` DROP `sell_eusd_num`;
ALTER TABLE `eos_otc_report` DROP `date`;

----------------------------------------------------
--  `otc_buy`
----------------------------------------------------
ALTER TABLE `otc_buy` DROP `uid`;
ALTER TABLE `otc_buy` DROP `available`;
ALTER TABLE `otc_buy` DROP `frozen`;
ALTER TABLE `otc_buy` DROP `bought`;
ALTER TABLE `otc_buy` DROP `lower_limit_wechat`;
ALTER TABLE `otc_buy` DROP `upper_limit_wechat`;
ALTER TABLE `otc_buy` DROP `lower_limit_bank`;
ALTER TABLE `otc_buy` DROP `upper_limit_bank`;
ALTER TABLE `otc_buy` DROP `lower_limit_ali`;
ALTER TABLE `otc_buy` DROP `upper_limit_ali`;
ALTER TABLE `otc_buy` DROP `pay_type`;
ALTER TABLE `otc_buy` DROP `ctime`;

----------------------------------------------------
--  `otc_exchanger`
----------------------------------------------------
ALTER TABLE `otc_exchanger` DROP `uid`;
ALTER TABLE `otc_exchanger` DROP `mobile`;
ALTER TABLE `otc_exchanger` DROP `wechat`;
ALTER TABLE `otc_exchanger` DROP `telegram`;
ALTER TABLE `otc_exchanger` DROP `from`;
ALTER TABLE `otc_exchanger` DROP `ctime`;
ALTER TABLE `otc_exchanger` DROP `utime`;

----------------------------------------------------
--  `otc_exchanger_verify`
----------------------------------------------------
ALTER TABLE `otc_exchanger_verify` DROP `id`;
ALTER TABLE `otc_exchanger_verify` DROP `uid`;
ALTER TABLE `otc_exchanger_verify` DROP `mobile`;
ALTER TABLE `otc_exchanger_verify` DROP `wechat`;
ALTER TABLE `otc_exchanger_verify` DROP `telegram`;
ALTER TABLE `otc_exchanger_verify` DROP `status`;
ALTER TABLE `otc_exchanger_verify` DROP `from`;
ALTER TABLE `otc_exchanger_verify` DROP `ctime`;
ALTER TABLE `otc_exchanger_verify` DROP `utime`;

----------------------------------------------------
--  `otc_msg`
----------------------------------------------------
ALTER TABLE `otc_msg` DROP `id`;
ALTER TABLE `otc_msg` DROP `order_id`;
ALTER TABLE `otc_msg` DROP `uid`;
ALTER TABLE `otc_msg` DROP `content`;
ALTER TABLE `otc_msg` DROP `is_read`;
ALTER TABLE `otc_msg` DROP `msg_type`;
ALTER TABLE `otc_msg` DROP `ctime`;

----------------------------------------------------
--  `otc_order`
----------------------------------------------------
ALTER TABLE `otc_order` DROP `id`;
ALTER TABLE `otc_order` DROP `uid`;
ALTER TABLE `otc_order` DROP `uip`;
ALTER TABLE `otc_order` DROP `euid`;
ALTER TABLE `otc_order` DROP `eip`;
ALTER TABLE `otc_order` DROP `side`;
ALTER TABLE `otc_order` DROP `amount`;
ALTER TABLE `otc_order` DROP `price`;
ALTER TABLE `otc_order` DROP `funds`;
ALTER TABLE `otc_order` DROP `fee`;
ALTER TABLE `otc_order` DROP `pay_id`;
ALTER TABLE `otc_order` DROP `pay_type`;
ALTER TABLE `otc_order` DROP `pay_name`;
ALTER TABLE `otc_order` DROP `pay_account`;
ALTER TABLE `otc_order` DROP `pay_bank`;
ALTER TABLE `otc_order` DROP `pay_bank_branch`;
ALTER TABLE `otc_order` DROP `transfer_id`;
ALTER TABLE `otc_order` DROP `ctime`;
ALTER TABLE `otc_order` DROP `pay_time`;
ALTER TABLE `otc_order` DROP `finish_time`;
ALTER TABLE `otc_order` DROP `utime`;
ALTER TABLE `otc_order` DROP `status`;
ALTER TABLE `otc_order` DROP `epay_id`;
ALTER TABLE `otc_order` DROP `epay_type`;
ALTER TABLE `otc_order` DROP `epay_name`;
ALTER TABLE `otc_order` DROP `epay_account`;
ALTER TABLE `otc_order` DROP `epay_bank`;
ALTER TABLE `otc_order` DROP `epay_bank_branch`;
ALTER TABLE `otc_order` DROP `appeal_status`;
ALTER TABLE `otc_order` DROP `admin_id`;
ALTER TABLE `otc_order` DROP `qr_code`;
ALTER TABLE `otc_order` DROP `date`;

----------------------------------------------------
--  `otc_sell`
----------------------------------------------------
ALTER TABLE `otc_sell` DROP `uid`;
ALTER TABLE `otc_sell` DROP `available`;
ALTER TABLE `otc_sell` DROP `frozen`;
ALTER TABLE `otc_sell` DROP `sold`;
ALTER TABLE `otc_sell` DROP `lower_limit`;
ALTER TABLE `otc_sell` DROP `upper_limit`;
ALTER TABLE `otc_sell` DROP `pay_type`;
ALTER TABLE `otc_sell` DROP `ctime`;

----------------------------------------------------
--  `payment_method`
----------------------------------------------------
ALTER TABLE `payment_method` DROP `pmid`;
ALTER TABLE `payment_method` DROP `uid`;
ALTER TABLE `payment_method` DROP `mtype`;
ALTER TABLE `payment_method` DROP `ord`;
ALTER TABLE `payment_method` DROP `name`;
ALTER TABLE `payment_method` DROP `account`;
ALTER TABLE `payment_method` DROP `status`;
ALTER TABLE `payment_method` DROP `ctime`;
ALTER TABLE `payment_method` DROP `bank`;
ALTER TABLE `payment_method` DROP `bank_branch`;
ALTER TABLE `payment_method` DROP `qr_code`;
ALTER TABLE `payment_method` DROP `qr_code_content`;
ALTER TABLE `payment_method` DROP `low_money_per_tx_limit`;
ALTER TABLE `payment_method` DROP `high_money_per_tx_limit`;
ALTER TABLE `payment_method` DROP `times_per_day_limit`;
ALTER TABLE `payment_method` DROP `money_per_day_limit`;
ALTER TABLE `payment_method` DROP `money_sum_limit`;
ALTER TABLE `payment_method` DROP `times_today`;
ALTER TABLE `payment_method` DROP `money_today`;
ALTER TABLE `payment_method` DROP `money_sum`;
ALTER TABLE `payment_method` DROP `mtime`;
ALTER TABLE `payment_method` DROP `use_time`;

----------------------------------------------------
--  `system_notification`
----------------------------------------------------
ALTER TABLE `system_notification` DROP `nid`;
ALTER TABLE `system_notification` DROP `notification_type`;
ALTER TABLE `system_notification` DROP `content`;
ALTER TABLE `system_notification` DROP `uid`;
ALTER TABLE `system_notification` DROP `is_read`;
ALTER TABLE `system_notification` DROP `ctime`;

----------------------------------------------------
--  `user`
----------------------------------------------------
ALTER TABLE `user` DROP `uid`;
ALTER TABLE `user` DROP `national_code`;
ALTER TABLE `user` DROP `mobile`;
ALTER TABLE `user` DROP `status`;
ALTER TABLE `user` DROP `nick`;
ALTER TABLE `user` DROP `pass`;
ALTER TABLE `user` DROP `salt`;
ALTER TABLE `user` DROP `ctime`;
ALTER TABLE `user` DROP `utime`;
ALTER TABLE `user` DROP `ip`;
ALTER TABLE `user` DROP `last_login_time`;
ALTER TABLE `user` DROP `last_login_ip`;
ALTER TABLE `user` DROP `is_exchanger`;
ALTER TABLE `user` DROP `sign_salt`;

----------------------------------------------------
--  `user_config`
----------------------------------------------------
ALTER TABLE `user_config` DROP `uid`;
ALTER TABLE `user_config` DROP `wealth_notice`;
ALTER TABLE `user_config` DROP `order_notice`;

----------------------------------------------------
--  `user_login_log`
----------------------------------------------------
ALTER TABLE `user_login_log` DROP `id`;
ALTER TABLE `user_login_log` DROP `user_id`;
ALTER TABLE `user_login_log` DROP `user_agent`;
ALTER TABLE `user_login_log` DROP `ips`;
ALTER TABLE `user_login_log` DROP `ctime`;

----------------------------------------------------
--  `user_pay_pass`
----------------------------------------------------
ALTER TABLE `user_pay_pass` DROP `uid`;
ALTER TABLE `user_pay_pass` DROP `pass`;
ALTER TABLE `user_pay_pass` DROP `salt`;
ALTER TABLE `user_pay_pass` DROP `sign_salt`;
ALTER TABLE `user_pay_pass` DROP `status`;
ALTER TABLE `user_pay_pass` DROP `method`;
ALTER TABLE `user_pay_pass` DROP `verify_step`;
ALTER TABLE `user_pay_pass` DROP `timestamp`;

