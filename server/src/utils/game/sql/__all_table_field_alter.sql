----------------------------------------------------
--  `channel_daily`
----------------------------------------------------
ALTER TABLE `channel_daily` CHANGE `channel_id` `channel_id` int unsigned NOT NULL COMMENT 'channel id';
ALTER TABLE `channel_daily` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'create time';
ALTER TABLE `channel_daily` CHANGE `win_lose_money_integer` `win_lose_money_integer` int NOT NULL COMMENT 'win lose money amount integer part';
ALTER TABLE `channel_daily` CHANGE `win_lose_money_decimals` `win_lose_money_decimals` int NOT NULL COMMENT 'win lose money amount decimals part';
ALTER TABLE `channel_daily` CHANGE `chips_integer` `chips_integer` int NOT NULL COMMENT 'chips amount integer part';
ALTER TABLE `channel_daily` CHANGE `chips_decimals` `chips_decimals` int NOT NULL COMMENT 'chips amount decimals part';
ALTER TABLE `channel_daily` CHANGE `mtime` `mtime` bigint NOT NULL COMMENT 'modified time';

----------------------------------------------------
--  `game_log`
----------------------------------------------------
ALTER TABLE `game_log` CHANGE `id` `id` bigint unsigned NOT NULL COMMENT 'id';
ALTER TABLE `game_log` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `game_log` CHANGE `channel_id` `channel_id` int unsigned NOT NULL COMMENT 'game channel id';
ALTER TABLE `game_log` CHANGE `account` `account` varchar(50) NOT NULL COMMENT 'account';
ALTER TABLE `game_log` CHANGE `log_type` `log_type` tinyint unsigned NOT NULL COMMENT 'log type';
ALTER TABLE `game_log` CHANGE `desc` `desc` varchar(512) NOT NULL COMMENT '';
ALTER TABLE `game_log` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'create time';

----------------------------------------------------
--  `game_order_risk`
----------------------------------------------------
ALTER TABLE `game_order_risk` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `game_order_risk` CHANGE `alert_id` `alert_id` bigint NOT NULL COMMENT 'group alert id';
ALTER TABLE `game_order_risk` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `game_order_risk` CHANGE `amount` `amount` bigint NOT NULL COMMENT 'eusd nums';
ALTER TABLE `game_order_risk` CHANGE `funds` `funds` bigint NOT NULL COMMENT 'rmb';
ALTER TABLE `game_order_risk` CHANGE `pay_type` `pay_type` tinyint NOT NULL COMMENT 'pay type';
ALTER TABLE `game_order_risk` CHANGE `pay_account` `pay_account` varchar(300) NOT NULL COMMENT 'pay account';
ALTER TABLE `game_order_risk` CHANGE `order_time` `order_time` bigint NOT NULL COMMENT 'order time';
ALTER TABLE `game_order_risk` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'ctime';

----------------------------------------------------
--  `game_risk_alert`
----------------------------------------------------
ALTER TABLE `game_risk_alert` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `game_risk_alert` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `game_risk_alert` CHANGE `funds` `funds` bigint NOT NULL COMMENT 'rmb num';
ALTER TABLE `game_risk_alert` CHANGE `eusd_num` `eusd_num` bigint NOT NULL COMMENT 'withdraw euse num';
ALTER TABLE `game_risk_alert` CHANGE `order_time` `order_time` bigint unsigned NOT NULL COMMENT 'order time';
ALTER TABLE `game_risk_alert` CHANGE `alert_time` `alert_time` bigint unsigned NOT NULL COMMENT 'alert risk time';
ALTER TABLE `game_risk_alert` CHANGE `do_get` `do_get` tinyint unsigned NOT NULL COMMENT 'weather get risk';
ALTER TABLE `game_risk_alert` CHANGE `warn_grade` `warn_grade` tinyint unsigned NOT NULL COMMENT 'warn grade';
ALTER TABLE `game_risk_alert` CHANGE `risk_type` `risk_type` tinyint unsigned NOT NULL COMMENT 'risk type';
ALTER TABLE `game_risk_alert` CHANGE `order_risk_id` `order_risk_id` bigint unsigned NOT NULL COMMENT 'order risk id';

----------------------------------------------------
--  `game_transfer`
----------------------------------------------------
ALTER TABLE `game_transfer` CHANGE `id` `id` bigint unsigned NOT NULL COMMENT 'id';
ALTER TABLE `game_transfer` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `game_transfer` CHANGE `channel_id` `channel_id` int unsigned NOT NULL COMMENT 'game channel id';
ALTER TABLE `game_transfer` CHANGE `account` `account` varchar(50) NOT NULL COMMENT 'account';
ALTER TABLE `game_transfer` CHANGE `transfer_type` `transfer_type` int unsigned NOT NULL COMMENT 'transfer type';
ALTER TABLE `game_transfer` CHANGE `order` `order` varchar(50) NOT NULL COMMENT 'Order';
ALTER TABLE `game_transfer` CHANGE `game_order` `game_order` varchar(50) NOT NULL COMMENT 'Game Order';
ALTER TABLE `game_transfer` CHANGE `coin_integer` `coin_integer` bigint NOT NULL COMMENT 'coin integer part';
ALTER TABLE `game_transfer` CHANGE `eusd_integer` `eusd_integer` bigint NOT NULL COMMENT 'eusd integer part';
ALTER TABLE `game_transfer` CHANGE `status` `status` int unsigned NOT NULL COMMENT 'status';
ALTER TABLE `game_transfer` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT '';
ALTER TABLE `game_transfer` CHANGE `desc` `desc` varchar(512) NOT NULL COMMENT '';
ALTER TABLE `game_transfer` CHANGE `step` `step` varchar(256) NOT NULL COMMENT '';

----------------------------------------------------
--  `game_user`
----------------------------------------------------
ALTER TABLE `game_user` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `game_user` CHANGE `channel_id` `channel_id` int unsigned NOT NULL COMMENT 'game channel id';
ALTER TABLE `game_user` CHANGE `account` `account` varchar(50) NOT NULL COMMENT 'account';
ALTER TABLE `game_user` CHANGE `nick_name` `nick_name` varchar(50) NOT NULL COMMENT 'nick name';
ALTER TABLE `game_user` CHANGE `sex` `sex` tinyint unsigned NOT NULL COMMENT 'sex';
ALTER TABLE `game_user` CHANGE `password` `password` varchar(100) NOT NULL COMMENT 'password';
ALTER TABLE `game_user` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT '';
ALTER TABLE `game_user` CHANGE `mtime` `mtime` bigint NOT NULL COMMENT '';
ALTER TABLE `game_user` CHANGE `status` `status` tinyint unsigned NOT NULL COMMENT '';

----------------------------------------------------
--  `game_user_daily`
----------------------------------------------------
ALTER TABLE `game_user_daily` CHANGE `channel_id` `channel_id` int unsigned NOT NULL COMMENT 'game channel id';
ALTER TABLE `game_user_daily` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `game_user_daily` CHANGE `tax_integer` `tax_integer` int NOT NULL COMMENT 'tax amount integer part';
ALTER TABLE `game_user_daily` CHANGE `tax_decimals` `tax_decimals` int NOT NULL COMMENT 'tax amount decimals part';
ALTER TABLE `game_user_daily` CHANGE `chips_integer` `chips_integer` int NOT NULL COMMENT 'chips amount integer part';
ALTER TABLE `game_user_daily` CHANGE `chips_decimals` `chips_decimals` int NOT NULL COMMENT 'chips amount decimals part';
ALTER TABLE `game_user_daily` CHANGE `winlose_integer` `winlose_integer` int NOT NULL COMMENT 'winlose amount integer part';
ALTER TABLE `game_user_daily` CHANGE `winlose_decimals` `winlose_decimals` int NOT NULL COMMENT 'winlose amount decimals part';
ALTER TABLE `game_user_daily` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'created time, equals to the begin time of the day';
ALTER TABLE `game_user_daily` CHANGE `mtime` `mtime` bigint NOT NULL COMMENT 'modified time';

