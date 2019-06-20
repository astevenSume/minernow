----------------------------------------------------
--  `channel_daily`
----------------------------------------------------
ALTER TABLE `channel_daily` DROP `channel_id`;
ALTER TABLE `channel_daily` DROP `ctime`;
ALTER TABLE `channel_daily` DROP `win_lose_money_integer`;
ALTER TABLE `channel_daily` DROP `win_lose_money_decimals`;
ALTER TABLE `channel_daily` DROP `chips_integer`;
ALTER TABLE `channel_daily` DROP `chips_decimals`;
ALTER TABLE `channel_daily` DROP `mtime`;

----------------------------------------------------
--  `game_log`
----------------------------------------------------
ALTER TABLE `game_log` DROP `id`;
ALTER TABLE `game_log` DROP `uid`;
ALTER TABLE `game_log` DROP `channel_id`;
ALTER TABLE `game_log` DROP `account`;
ALTER TABLE `game_log` DROP `log_type`;
ALTER TABLE `game_log` DROP `desc`;
ALTER TABLE `game_log` DROP `ctime`;

----------------------------------------------------
--  `game_order_risk`
----------------------------------------------------
ALTER TABLE `game_order_risk` DROP `id`;
ALTER TABLE `game_order_risk` DROP `alert_id`;
ALTER TABLE `game_order_risk` DROP `uid`;
ALTER TABLE `game_order_risk` DROP `amount`;
ALTER TABLE `game_order_risk` DROP `funds`;
ALTER TABLE `game_order_risk` DROP `pay_type`;
ALTER TABLE `game_order_risk` DROP `pay_account`;
ALTER TABLE `game_order_risk` DROP `order_time`;
ALTER TABLE `game_order_risk` DROP `ctime`;

----------------------------------------------------
--  `game_risk_alert`
----------------------------------------------------
ALTER TABLE `game_risk_alert` DROP `id`;
ALTER TABLE `game_risk_alert` DROP `uid`;
ALTER TABLE `game_risk_alert` DROP `funds`;
ALTER TABLE `game_risk_alert` DROP `eusd_num`;
ALTER TABLE `game_risk_alert` DROP `order_time`;
ALTER TABLE `game_risk_alert` DROP `alert_time`;
ALTER TABLE `game_risk_alert` DROP `do_get`;
ALTER TABLE `game_risk_alert` DROP `warn_grade`;
ALTER TABLE `game_risk_alert` DROP `risk_type`;
ALTER TABLE `game_risk_alert` DROP `order_risk_id`;

----------------------------------------------------
--  `game_transfer`
----------------------------------------------------
ALTER TABLE `game_transfer` DROP `id`;
ALTER TABLE `game_transfer` DROP `uid`;
ALTER TABLE `game_transfer` DROP `channel_id`;
ALTER TABLE `game_transfer` DROP `account`;
ALTER TABLE `game_transfer` DROP `transfer_type`;
ALTER TABLE `game_transfer` DROP `order`;
ALTER TABLE `game_transfer` DROP `game_order`;
ALTER TABLE `game_transfer` DROP `coin_integer`;
ALTER TABLE `game_transfer` DROP `eusd_integer`;
ALTER TABLE `game_transfer` DROP `status`;
ALTER TABLE `game_transfer` DROP `ctime`;
ALTER TABLE `game_transfer` DROP `desc`;
ALTER TABLE `game_transfer` DROP `step`;

----------------------------------------------------
--  `game_user`
----------------------------------------------------
ALTER TABLE `game_user` DROP `uid`;
ALTER TABLE `game_user` DROP `channel_id`;
ALTER TABLE `game_user` DROP `account`;
ALTER TABLE `game_user` DROP `nick_name`;
ALTER TABLE `game_user` DROP `sex`;
ALTER TABLE `game_user` DROP `password`;
ALTER TABLE `game_user` DROP `ctime`;
ALTER TABLE `game_user` DROP `mtime`;
ALTER TABLE `game_user` DROP `status`;

----------------------------------------------------
--  `game_user_daily`
----------------------------------------------------
ALTER TABLE `game_user_daily` DROP `channel_id`;
ALTER TABLE `game_user_daily` DROP `uid`;
ALTER TABLE `game_user_daily` DROP `tax_integer`;
ALTER TABLE `game_user_daily` DROP `tax_decimals`;
ALTER TABLE `game_user_daily` DROP `chips_integer`;
ALTER TABLE `game_user_daily` DROP `chips_decimals`;
ALTER TABLE `game_user_daily` DROP `winlose_integer`;
ALTER TABLE `game_user_daily` DROP `winlose_decimals`;
ALTER TABLE `game_user_daily` DROP `ctime`;
ALTER TABLE `game_user_daily` DROP `mtime`;

