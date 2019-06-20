----------------------------------------------------
--  `channel_daily`
----------------------------------------------------
ALTER TABLE `channel_daily` ADD `channel_id` int unsigned NOT NULL COMMENT 'channel id';
ALTER TABLE `channel_daily` ADD `ctime` bigint NOT NULL COMMENT 'create time' AFTER `channel_id`;
ALTER TABLE `channel_daily` ADD `win_lose_money_integer` int NOT NULL COMMENT 'win lose money amount integer part' AFTER `ctime`;
ALTER TABLE `channel_daily` ADD `win_lose_money_decimals` int NOT NULL COMMENT 'win lose money amount decimals part' AFTER `win_lose_money_integer`;
ALTER TABLE `channel_daily` ADD `chips_integer` int NOT NULL COMMENT 'chips amount integer part' AFTER `win_lose_money_decimals`;
ALTER TABLE `channel_daily` ADD `chips_decimals` int NOT NULL COMMENT 'chips amount decimals part' AFTER `chips_integer`;
ALTER TABLE `channel_daily` ADD `mtime` bigint NOT NULL COMMENT 'modified time' AFTER `chips_decimals`;

----------------------------------------------------
--  `game_log`
----------------------------------------------------
ALTER TABLE `game_log` ADD `id` bigint unsigned NOT NULL COMMENT 'id';
ALTER TABLE `game_log` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id' AFTER `id`;
ALTER TABLE `game_log` ADD `channel_id` int unsigned NOT NULL COMMENT 'game channel id' AFTER `uid`;
ALTER TABLE `game_log` ADD `account` varchar(50) NOT NULL COMMENT 'account' AFTER `channel_id`;
ALTER TABLE `game_log` ADD `log_type` tinyint unsigned NOT NULL COMMENT 'log type' AFTER `account`;
ALTER TABLE `game_log` ADD `desc` varchar(512) NOT NULL COMMENT '' AFTER `log_type`;
ALTER TABLE `game_log` ADD `ctime` bigint NOT NULL COMMENT 'create time' AFTER `desc`;

----------------------------------------------------
--  `game_order_risk`
----------------------------------------------------
ALTER TABLE `game_order_risk` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `game_order_risk` ADD `alert_id` bigint NOT NULL COMMENT 'group alert id' AFTER `id`;
ALTER TABLE `game_order_risk` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id' AFTER `alert_id`;
ALTER TABLE `game_order_risk` ADD `amount` bigint NOT NULL COMMENT 'eusd nums' AFTER `uid`;
ALTER TABLE `game_order_risk` ADD `funds` bigint NOT NULL COMMENT 'rmb' AFTER `amount`;
ALTER TABLE `game_order_risk` ADD `pay_type` tinyint NOT NULL COMMENT 'pay type' AFTER `funds`;
ALTER TABLE `game_order_risk` ADD `pay_account` varchar(300) NOT NULL COMMENT 'pay account' AFTER `pay_type`;
ALTER TABLE `game_order_risk` ADD `order_time` bigint NOT NULL COMMENT 'order time' AFTER `pay_account`;
ALTER TABLE `game_order_risk` ADD `ctime` bigint NOT NULL COMMENT 'ctime' AFTER `order_time`;

----------------------------------------------------
--  `game_risk_alert`
----------------------------------------------------
ALTER TABLE `game_risk_alert` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `game_risk_alert` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id' AFTER `id`;
ALTER TABLE `game_risk_alert` ADD `funds` bigint NOT NULL COMMENT 'rmb num' AFTER `uid`;
ALTER TABLE `game_risk_alert` ADD `eusd_num` bigint NOT NULL COMMENT 'withdraw euse num' AFTER `funds`;
ALTER TABLE `game_risk_alert` ADD `order_time` bigint unsigned NOT NULL COMMENT 'order time' AFTER `eusd_num`;
ALTER TABLE `game_risk_alert` ADD `alert_time` bigint unsigned NOT NULL COMMENT 'alert risk time' AFTER `order_time`;
ALTER TABLE `game_risk_alert` ADD `do_get` tinyint unsigned NOT NULL COMMENT 'weather get risk' AFTER `alert_time`;
ALTER TABLE `game_risk_alert` ADD `warn_grade` tinyint unsigned NOT NULL COMMENT 'warn grade' AFTER `do_get`;
ALTER TABLE `game_risk_alert` ADD `risk_type` tinyint unsigned NOT NULL COMMENT 'risk type' AFTER `warn_grade`;
ALTER TABLE `game_risk_alert` ADD `order_risk_id` bigint unsigned NOT NULL COMMENT 'order risk id' AFTER `risk_type`;

----------------------------------------------------
--  `game_transfer`
----------------------------------------------------
ALTER TABLE `game_transfer` ADD `id` bigint unsigned NOT NULL COMMENT 'id';
ALTER TABLE `game_transfer` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id' AFTER `id`;
ALTER TABLE `game_transfer` ADD `channel_id` int unsigned NOT NULL COMMENT 'game channel id' AFTER `uid`;
ALTER TABLE `game_transfer` ADD `account` varchar(50) NOT NULL COMMENT 'account' AFTER `channel_id`;
ALTER TABLE `game_transfer` ADD `transfer_type` int unsigned NOT NULL COMMENT 'transfer type' AFTER `account`;
ALTER TABLE `game_transfer` ADD `order` varchar(50) NOT NULL COMMENT 'Order' AFTER `transfer_type`;
ALTER TABLE `game_transfer` ADD `game_order` varchar(50) NOT NULL COMMENT 'Game Order' AFTER `order`;
ALTER TABLE `game_transfer` ADD `coin_integer` bigint NOT NULL COMMENT 'coin integer part' AFTER `game_order`;
ALTER TABLE `game_transfer` ADD `eusd_integer` bigint NOT NULL COMMENT 'eusd integer part' AFTER `coin_integer`;
ALTER TABLE `game_transfer` ADD `status` int unsigned NOT NULL COMMENT 'status' AFTER `eusd_integer`;
ALTER TABLE `game_transfer` ADD `ctime` bigint NOT NULL COMMENT '' AFTER `status`;
ALTER TABLE `game_transfer` ADD `desc` varchar(512) NOT NULL COMMENT '' AFTER `ctime`;
ALTER TABLE `game_transfer` ADD `step` varchar(256) NOT NULL COMMENT '' AFTER `desc`;

----------------------------------------------------
--  `game_user`
----------------------------------------------------
ALTER TABLE `game_user` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `game_user` ADD `channel_id` int unsigned NOT NULL COMMENT 'game channel id' AFTER `uid`;
ALTER TABLE `game_user` ADD `account` varchar(50) NOT NULL COMMENT 'account' AFTER `channel_id`;
ALTER TABLE `game_user` ADD `nick_name` varchar(50) NOT NULL COMMENT 'nick name' AFTER `account`;
ALTER TABLE `game_user` ADD `sex` tinyint unsigned NOT NULL COMMENT 'sex' AFTER `nick_name`;
ALTER TABLE `game_user` ADD `password` varchar(100) NOT NULL COMMENT 'password' AFTER `sex`;
ALTER TABLE `game_user` ADD `ctime` bigint NOT NULL COMMENT '' AFTER `password`;
ALTER TABLE `game_user` ADD `mtime` bigint NOT NULL COMMENT '' AFTER `ctime`;
ALTER TABLE `game_user` ADD `status` tinyint unsigned NOT NULL COMMENT '' AFTER `mtime`;

----------------------------------------------------
--  `game_user_daily`
----------------------------------------------------
ALTER TABLE `game_user_daily` ADD `channel_id` int unsigned NOT NULL COMMENT 'game channel id';
ALTER TABLE `game_user_daily` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id' AFTER `channel_id`;
ALTER TABLE `game_user_daily` ADD `tax_integer` int NOT NULL COMMENT 'tax amount integer part' AFTER `uid`;
ALTER TABLE `game_user_daily` ADD `tax_decimals` int NOT NULL COMMENT 'tax amount decimals part' AFTER `tax_integer`;
ALTER TABLE `game_user_daily` ADD `chips_integer` int NOT NULL COMMENT 'chips amount integer part' AFTER `tax_decimals`;
ALTER TABLE `game_user_daily` ADD `chips_decimals` int NOT NULL COMMENT 'chips amount decimals part' AFTER `chips_integer`;
ALTER TABLE `game_user_daily` ADD `winlose_integer` int NOT NULL COMMENT 'winlose amount integer part' AFTER `chips_decimals`;
ALTER TABLE `game_user_daily` ADD `winlose_decimals` int NOT NULL COMMENT 'winlose amount decimals part' AFTER `winlose_integer`;
ALTER TABLE `game_user_daily` ADD `ctime` bigint NOT NULL COMMENT 'created time, equals to the begin time of the day' AFTER `winlose_decimals`;
ALTER TABLE `game_user_daily` ADD `mtime` bigint NOT NULL COMMENT 'modified time' AFTER `ctime`;

