----------------------------------------------------
--  `activity_user_conf`
----------------------------------------------------
ALTER TABLE `activity_user_conf` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `activity_user_conf` ADD `play_game_day` int NOT NULL COMMENT '玩了并且有投注的天数 配置' AFTER `id`;
ALTER TABLE `activity_user_conf` ADD `bet_amount` bigint NOT NULL COMMENT '投注额配置' AFTER `play_game_day`;
ALTER TABLE `activity_user_conf` ADD `effective_bet_amount` bigint NOT NULL COMMENT '累计有效投注额(本人加无限下级)的配置' AFTER `bet_amount`;

----------------------------------------------------
--  `admin_user`
----------------------------------------------------
ALTER TABLE `admin_user` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'user id';
ALTER TABLE `admin_user` ADD `name` varchar(100) NOT NULL COMMENT 'user name' AFTER `id`;
ALTER TABLE `admin_user` ADD `email` varchar(100) NOT NULL COMMENT 'user email' AFTER `name`;
ALTER TABLE `admin_user` ADD `status` tinyint NOT NULL COMMENT 'status' AFTER `email`;
ALTER TABLE `admin_user` ADD `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time' AFTER `status`;
ALTER TABLE `admin_user` ADD `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time' AFTER `ctime`;
ALTER TABLE `admin_user` ADD `dtime` bigint NOT NULL DEFAULT 0 COMMENT 'forbit time' AFTER `utime`;
ALTER TABLE `admin_user` ADD `login_time` bigint NOT NULL DEFAULT 0 COMMENT 'last login time' AFTER `dtime`;
ALTER TABLE `admin_user` ADD `pwd` varchar(100) NOT NULL DEFAULT '' COMMENT 'password' AFTER `login_time`;
ALTER TABLE `admin_user` ADD `whitelist_ips` varchar(256) NOT NULL DEFAULT '' COMMENT 'whitelist_ips' AFTER `pwd`;
ALTER TABLE `admin_user` ADD `is_bind` bool NOT NULL DEFAULT false COMMENT 'is_bind' AFTER `whitelist_ips`;
ALTER TABLE `admin_user` ADD `secret_id` varchar(128) NOT NULL DEFAULT '' COMMENT 'secret_id' AFTER `is_bind`;
ALTER TABLE `admin_user` ADD `qr_code` text NOT NULL COMMENT 'qr_code' AFTER `secret_id`;

----------------------------------------------------
--  `agent_white_list`
----------------------------------------------------
ALTER TABLE `agent_white_list` ADD `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `agent_white_list` ADD `name` varchar(64) NOT NULL COMMENT 'name' AFTER `id`;
ALTER TABLE `agent_white_list` ADD `commission` int NOT NULL COMMENT 'commission' AFTER `name`;
ALTER TABLE `agent_white_list` ADD `precision` int NOT NULL COMMENT 'precision' AFTER `commission`;
ALTER TABLE `agent_white_list` ADD `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time' AFTER `precision`;
ALTER TABLE `agent_white_list` ADD `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time' AFTER `ctime`;

----------------------------------------------------
--  `announcement`
----------------------------------------------------
ALTER TABLE `announcement` ADD `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `announcement` ADD `type` tinyint NOT NULL COMMENT 'type' AFTER `id`;
ALTER TABLE `announcement` ADD `title` varchar(100) NOT NULL COMMENT 'title' AFTER `type`;
ALTER TABLE `announcement` ADD `content` varchar(500) NOT NULL COMMENT 'content' AFTER `title`;
ALTER TABLE `announcement` ADD `stime` bigint NOT NULL DEFAULT 0 COMMENT 'start time' AFTER `content`;
ALTER TABLE `announcement` ADD `etime` bigint NOT NULL DEFAULT 0 COMMENT 'end time' AFTER `stime`;
ALTER TABLE `announcement` ADD `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time' AFTER `etime`;
ALTER TABLE `announcement` ADD `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time' AFTER `ctime`;

----------------------------------------------------
--  `app_channel`
----------------------------------------------------
ALTER TABLE `app_channel` ADD `id` int unsigned NOT NULL COMMENT 'id';
ALTER TABLE `app_channel` ADD `is_third_hall` tinyint NOT NULL COMMENT 'is_third_hall' AFTER `id`;
ALTER TABLE `app_channel` ADD `name` varchar(256) NOT NULL COMMENT 'name' AFTER `is_third_hall`;
ALTER TABLE `app_channel` ADD `desc` varchar(256) NOT NULL COMMENT 'desc' AFTER `name`;
ALTER TABLE `app_channel` ADD `exchangeRate` int NOT NULL DEFAULT 0 COMMENT 'exchange rate' AFTER `desc`;
ALTER TABLE `app_channel` ADD `precision` int NOT NULL DEFAULT 0 COMMENT 'precision' AFTER `exchangeRate`;
ALTER TABLE `app_channel` ADD `profit_rate` int NOT NULL DEFAULT 0 COMMENT 'profit_rate' AFTER `precision`;
ALTER TABLE `app_channel` ADD `icon_url` varchar(128) NOT NULL COMMENT 'icon_url' AFTER `profit_rate`;
ALTER TABLE `app_channel` ADD `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time' AFTER `icon_url`;
ALTER TABLE `app_channel` ADD `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time' AFTER `ctime`;

----------------------------------------------------
--  `app_type`
----------------------------------------------------
ALTER TABLE `app_type` ADD `id` int unsigned NOT NULL COMMENT 'id';
ALTER TABLE `app_type` ADD `name` varchar(256) NOT NULL COMMENT 'name' AFTER `id`;
ALTER TABLE `app_type` ADD `desc` varchar(256) NOT NULL COMMENT 'desc' AFTER `name`;
ALTER TABLE `app_type` ADD `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time' AFTER `desc`;
ALTER TABLE `app_type` ADD `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time' AFTER `ctime`;

----------------------------------------------------
--  `app_version`
----------------------------------------------------
ALTER TABLE `app_version` ADD `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `app_version` ADD `version` varchar(100) NOT NULL COMMENT 'app version show' AFTER `id`;
ALTER TABLE `app_version` ADD `version_num` int NOT NULL COMMENT 'app version num' AFTER `version`;
ALTER TABLE `app_version` ADD `changelog` varchar(300) NOT NULL COMMENT 'change log' AFTER `version_num`;
ALTER TABLE `app_version` ADD `download` varchar(300) NOT NULL COMMENT 'download url' AFTER `changelog`;
ALTER TABLE `app_version` ADD `system` tinyint NOT NULL COMMENT 'system type' AFTER `download`;
ALTER TABLE `app_version` ADD `status` tinyint NOT NULL COMMENT 'app version status' AFTER `system`;
ALTER TABLE `app_version` ADD `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time' AFTER `status`;
ALTER TABLE `app_version` ADD `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time' AFTER `ctime`;
ALTER TABLE `app_version` ADD `dtime` bigint NOT NULL DEFAULT 0 COMMENT 'update dtime' AFTER `utime`;

----------------------------------------------------
--  `app_whitelist`
----------------------------------------------------
ALTER TABLE `app_whitelist` ADD `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `app_whitelist` ADD `channel_id` int unsigned NOT NULL COMMENT 'channel_id' AFTER `id`;
ALTER TABLE `app_whitelist` ADD `app_id` varchar(16) NOT NULL COMMENT 'app_id' AFTER `channel_id`;
ALTER TABLE `app_whitelist` ADD `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time' AFTER `app_id`;

----------------------------------------------------
--  `appeal_service`
----------------------------------------------------
ALTER TABLE `appeal_service` ADD `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `appeal_service` ADD `admin_id` int unsigned NOT NULL COMMENT 'admin_id' AFTER `id`;
ALTER TABLE `appeal_service` ADD `wechat` varchar(32) NOT NULL COMMENT 'wechat' AFTER `admin_id`;
ALTER TABLE `appeal_service` ADD `qr_code` varchar(300) NOT NULL COMMENT '' AFTER `wechat`;
ALTER TABLE `appeal_service` ADD `status` tinyint NOT NULL COMMENT 'status' AFTER `qr_code`;

----------------------------------------------------
--  `apps`
----------------------------------------------------
ALTER TABLE `apps` ADD `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `apps` ADD `position` int unsigned NOT NULL COMMENT 'position' AFTER `id`;
ALTER TABLE `apps` ADD `name` varchar(64) NOT NULL COMMENT 'name' AFTER `position`;
ALTER TABLE `apps` ADD `desc` varchar(256) NOT NULL COMMENT 'desc' AFTER `name`;
ALTER TABLE `apps` ADD `url` varchar(128) NOT NULL COMMENT 'url' AFTER `desc`;
ALTER TABLE `apps` ADD `icon_url` varchar(128) NOT NULL COMMENT 'icon_url' AFTER `url`;
ALTER TABLE `apps` ADD `type_id` tinyint NOT NULL COMMENT 'type_id' AFTER `icon_url`;
ALTER TABLE `apps` ADD `channel_id` int unsigned NOT NULL COMMENT 'channel_id' AFTER `type_id`;
ALTER TABLE `apps` ADD `app_id` varchar(50) NOT NULL COMMENT 'app_id' AFTER `channel_id`;
ALTER TABLE `apps` ADD `featured` tinyint NOT NULL COMMENT 'featured' AFTER `app_id`;
ALTER TABLE `apps` ADD `status` tinyint NOT NULL COMMENT 'status' AFTER `featured`;
ALTER TABLE `apps` ADD `orientation` tinyint NOT NULL DEFAULT 1 COMMENT 'orientation' AFTER `status`;
ALTER TABLE `apps` ADD `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time' AFTER `orientation`;
ALTER TABLE `apps` ADD `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time' AFTER `ctime`;

----------------------------------------------------
--  `banner`
----------------------------------------------------
ALTER TABLE `banner` ADD `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `banner` ADD `subject` varchar(256) NOT NULL COMMENT 'subject' AFTER `id`;
ALTER TABLE `banner` ADD `image` varchar(256) NOT NULL COMMENT 'image' AFTER `subject`;
ALTER TABLE `banner` ADD `url` varchar(256) NOT NULL DEFAULT '' COMMENT 'image' AFTER `image`;
ALTER TABLE `banner` ADD `status` tinyint NOT NULL COMMENT 'status' AFTER `url`;
ALTER TABLE `banner` ADD `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time' AFTER `status`;
ALTER TABLE `banner` ADD `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time' AFTER `ctime`;
ALTER TABLE `banner` ADD `stime` bigint NOT NULL DEFAULT 0 COMMENT 'start time' AFTER `utime`;
ALTER TABLE `banner` ADD `etime` bigint NOT NULL DEFAULT 0 COMMENT 'end time' AFTER `stime`;

----------------------------------------------------
--  `commissionrates`
----------------------------------------------------
ALTER TABLE `commissionrates` ADD `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `commissionrates` ADD `min` bigint unsigned NOT NULL COMMENT 'min' AFTER `id`;
ALTER TABLE `commissionrates` ADD `max` bigint unsigned NOT NULL COMMENT 'max' AFTER `min`;
ALTER TABLE `commissionrates` ADD `commission` int NOT NULL DEFAULT 0 COMMENT 'commission' AFTER `max`;
ALTER TABLE `commissionrates` ADD `precision` int NOT NULL DEFAULT 0 COMMENT 'precision' AFTER `commission`;
ALTER TABLE `commissionrates` ADD `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time' AFTER `precision`;
ALTER TABLE `commissionrates` ADD `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time' AFTER `ctime`;

----------------------------------------------------
--  `config`
----------------------------------------------------
ALTER TABLE `config` ADD `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `config` ADD `action` tinyint NOT NULL COMMENT 'action' AFTER `id`;
ALTER TABLE `config` ADD `key` varchar(256) NOT NULL COMMENT 'key' AFTER `action`;
ALTER TABLE `config` ADD `value` text NOT NULL COMMENT 'value' AFTER `key`;
ALTER TABLE `config` ADD `desc` varchar(256) NOT NULL COMMENT 'descripe' AFTER `value`;
ALTER TABLE `config` ADD `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time' AFTER `desc`;

----------------------------------------------------
--  `config_warning`
----------------------------------------------------
ALTER TABLE `config_warning` ADD `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `config_warning` ADD `type` tinyint NOT NULL COMMENT 'type' AFTER `id`;
ALTER TABLE `config_warning` ADD `national_code` varchar(16) NOT NULL COMMENT 'national_code' AFTER `type`;
ALTER TABLE `config_warning` ADD `mobile` varchar(32) NOT NULL COMMENT 'mobile' AFTER `national_code`;
ALTER TABLE `config_warning` ADD `sms_type` tinyint NOT NULL COMMENT 'sms_type' AFTER `mobile`;

----------------------------------------------------
--  `endpoint`
----------------------------------------------------
ALTER TABLE `endpoint` ADD `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `endpoint` ADD `endpoint` varchar(100) NOT NULL COMMENT 'endpoint' AFTER `id`;
ALTER TABLE `endpoint` ADD `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time' AFTER `endpoint`;
ALTER TABLE `endpoint` ADD `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time' AFTER `ctime`;

----------------------------------------------------
--  `ip_white_list`
----------------------------------------------------
ALTER TABLE `ip_white_list` ADD `id` int unsigned NOT NULL COMMENT 'id';
ALTER TABLE `ip_white_list` ADD `ip` varchar(0) NOT NULL COMMENT 'ip address' AFTER `id`;

----------------------------------------------------
--  `menu_access`
----------------------------------------------------
ALTER TABLE `menu_access` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `menu_access` ADD `role_id` bigint unsigned NOT NULL COMMENT '角色id' AFTER `id`;
ALTER TABLE `menu_access` ADD `menu_id` bigint unsigned NOT NULL COMMENT '菜单id' AFTER `role_id`;

----------------------------------------------------
--  `menu_conf`
----------------------------------------------------
ALTER TABLE `menu_conf` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `menu_conf` ADD `pid` bigint unsigned NOT NULL COMMENT '父节点的id，一级菜单为0' AFTER `id`;
ALTER TABLE `menu_conf` ADD `level` int NOT NULL COMMENT '层级,1是一级菜单,2是二级菜单' AFTER `pid`;
ALTER TABLE `menu_conf` ADD `name` varchar(100) NOT NULL COMMENT '菜单名称' AFTER `level`;
ALTER TABLE `menu_conf` ADD `path` varchar(100) NOT NULL COMMENT '菜单路径' AFTER `name`;
ALTER TABLE `menu_conf` ADD `icon` varchar(100) NOT NULL COMMENT '一级菜单图标' AFTER `path`;
ALTER TABLE `menu_conf` ADD `hide_in_menu` bool NOT NULL COMMENT '是否隐藏' AFTER `icon`;
ALTER TABLE `menu_conf` ADD `component` varchar(100) NOT NULL COMMENT '组件' AFTER `hide_in_menu`;
ALTER TABLE `menu_conf` ADD `order_id` int unsigned NOT NULL COMMENT '排序id' AFTER `component`;
ALTER TABLE `menu_conf` ADD `ctime` bigint NOT NULL COMMENT '创建时间' AFTER `order_id`;
ALTER TABLE `menu_conf` ADD `utime` bigint NOT NULL COMMENT '修改时间' AFTER `ctime`;

----------------------------------------------------
--  `month_dividend_position_conf`
----------------------------------------------------
ALTER TABLE `month_dividend_position_conf` ADD `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `month_dividend_position_conf` ADD `agent_lv` int NOT NULL COMMENT 'agent_lv' AFTER `id`;
ALTER TABLE `month_dividend_position_conf` ADD `position` int NOT NULL COMMENT 'position' AFTER `agent_lv`;
ALTER TABLE `month_dividend_position_conf` ADD `min` bigint NOT NULL COMMENT 'min' AFTER `position`;
ALTER TABLE `month_dividend_position_conf` ADD `max` bigint NOT NULL COMMENT 'max' AFTER `min`;
ALTER TABLE `month_dividend_position_conf` ADD `activity_num` int NOT NULL COMMENT 'activity num' AFTER `max`;
ALTER TABLE `month_dividend_position_conf` ADD `dividend_ratio` int NOT NULL COMMENT 'dividend_ratio num' AFTER `activity_num`;
ALTER TABLE `month_dividend_position_conf` ADD `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time' AFTER `dividend_ratio`;
ALTER TABLE `month_dividend_position_conf` ADD `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time' AFTER `ctime`;

----------------------------------------------------
--  `month_dividend_white_list`
----------------------------------------------------
ALTER TABLE `month_dividend_white_list` ADD `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `month_dividend_white_list` ADD `name` varchar(64) NOT NULL COMMENT 'name' AFTER `id`;
ALTER TABLE `month_dividend_white_list` ADD `dividend_ratio` int NOT NULL COMMENT 'dividend_ratio' AFTER `name`;
ALTER TABLE `month_dividend_white_list` ADD `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time' AFTER `dividend_ratio`;
ALTER TABLE `month_dividend_white_list` ADD `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time' AFTER `ctime`;

----------------------------------------------------
--  `operation_log`
----------------------------------------------------
ALTER TABLE `operation_log` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `operation_log` ADD `admin_id` bigint unsigned NOT NULL COMMENT 'admin_id' AFTER `id`;
ALTER TABLE `operation_log` ADD `method` varchar(100) NOT NULL DEFAULT '' COMMENT 'req method' AFTER `admin_id`;
ALTER TABLE `operation_log` ADD `route` varchar(100) NOT NULL DEFAULT '' COMMENT 'route ' AFTER `method`;
ALTER TABLE `operation_log` ADD `action` int NOT NULL DEFAULT 0 COMMENT 'action' AFTER `route`;
ALTER TABLE `operation_log` ADD `input` varchar(65535) NOT NULL DEFAULT '' COMMENT 'input' AFTER `action`;
ALTER TABLE `operation_log` ADD `user_agent` varchar(512) NOT NULL DEFAULT '' COMMENT 'user agent' AFTER `input`;
ALTER TABLE `operation_log` ADD `ips` varchar(100) NOT NULL DEFAULT '' COMMENT 'ip' AFTER `user_agent`;
ALTER TABLE `operation_log` ADD `response_code` int NOT NULL DEFAULT 0 COMMENT 'response code' AFTER `ips`;
ALTER TABLE `operation_log` ADD `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time' AFTER `response_code`;

----------------------------------------------------
--  `otc_stat`
----------------------------------------------------
ALTER TABLE `otc_stat` ADD `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `otc_stat` ADD `date` int unsigned NOT NULL COMMENT 'date' AFTER `id`;
ALTER TABLE `otc_stat` ADD `num_login` int unsigned NOT NULL DEFAULT 0 COMMENT '' AFTER `date`;
ALTER TABLE `otc_stat` ADD `num_user_new` int unsigned NOT NULL DEFAULT 0 COMMENT '' AFTER `num_login`;
ALTER TABLE `otc_stat` ADD `num_order` int unsigned NOT NULL DEFAULT 0 COMMENT 'date' AFTER `num_user_new`;
ALTER TABLE `otc_stat` ADD `num_order_deal` int unsigned NOT NULL DEFAULT 0 COMMENT '' AFTER `num_order`;
ALTER TABLE `otc_stat` ADD `num_order_buy` int unsigned NOT NULL DEFAULT 0 COMMENT '' AFTER `num_order_deal`;
ALTER TABLE `otc_stat` ADD `num_order_sell` int unsigned NOT NULL DEFAULT 0 COMMENT '' AFTER `num_order_buy`;
ALTER TABLE `otc_stat` ADD `num_funds` int unsigned NOT NULL DEFAULT 0 COMMENT '' AFTER `num_order_sell`;
ALTER TABLE `otc_stat` ADD `num_amount` int unsigned NOT NULL DEFAULT 0 COMMENT '' AFTER `num_funds`;
ALTER TABLE `otc_stat` ADD `num_amount_buy` int unsigned NOT NULL DEFAULT 0 COMMENT '' AFTER `num_amount`;
ALTER TABLE `otc_stat` ADD `num_amount_sell` int unsigned NOT NULL DEFAULT 0 COMMENT '' AFTER `num_amount_buy`;
ALTER TABLE `otc_stat` ADD `num_fee_buy` int unsigned NOT NULL DEFAULT 0 COMMENT '' AFTER `num_amount_sell`;
ALTER TABLE `otc_stat` ADD `num_fee_sell` int unsigned NOT NULL DEFAULT 0 COMMENT '' AFTER `num_fee_buy`;
ALTER TABLE `otc_stat` ADD `game_recharge` int unsigned NOT NULL DEFAULT 0 COMMENT '' AFTER `num_fee_sell`;
ALTER TABLE `otc_stat` ADD `game_withdrawal` int unsigned NOT NULL DEFAULT 0 COMMENT '' AFTER `game_recharge`;
ALTER TABLE `otc_stat` ADD `usdt_recharge` int unsigned NOT NULL DEFAULT 0 COMMENT '' AFTER `game_withdrawal`;
ALTER TABLE `otc_stat` ADD `usdt_withdrawal` int unsigned NOT NULL DEFAULT 0 COMMENT '' AFTER `usdt_recharge`;
ALTER TABLE `otc_stat` ADD `usdt_fee` int unsigned NOT NULL DEFAULT 0 COMMENT '' AFTER `usdt_withdrawal`;

----------------------------------------------------
--  `otc_stat_all_people`
----------------------------------------------------
ALTER TABLE `otc_stat_all_people` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `otc_stat_all_people` ADD `buy_order` int unsigned NOT NULL DEFAULT 0 COMMENT '' AFTER `uid`;
ALTER TABLE `otc_stat_all_people` ADD `sell_order` int unsigned NOT NULL DEFAULT 0 COMMENT '' AFTER `buy_order`;
ALTER TABLE `otc_stat_all_people` ADD `buy_eusd` int unsigned NOT NULL DEFAULT 0 COMMENT '' AFTER `sell_order`;
ALTER TABLE `otc_stat_all_people` ADD `sell_eusd` int unsigned NOT NULL DEFAULT 0 COMMENT '' AFTER `buy_eusd`;
ALTER TABLE `otc_stat_all_people` ADD `usdt_recharge` int unsigned NOT NULL DEFAULT 0 COMMENT '' AFTER `sell_eusd`;
ALTER TABLE `otc_stat_all_people` ADD `usdt_withdrawal` int unsigned NOT NULL DEFAULT 0 COMMENT '' AFTER `usdt_recharge`;

----------------------------------------------------
--  `permission`
----------------------------------------------------
ALTER TABLE `permission` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `permission` ADD `slug` varchar(100) NOT NULL COMMENT 'permission_name' AFTER `id`;
ALTER TABLE `permission` ADD `desc` varchar(100) NOT NULL COMMENT 'permission_desc' AFTER `slug`;
ALTER TABLE `permission` ADD `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time' AFTER `desc`;
ALTER TABLE `permission` ADD `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time' AFTER `ctime`;
ALTER TABLE `permission` ADD `dtime` bigint NOT NULL DEFAULT 0 COMMENT 'D time' AFTER `utime`;

----------------------------------------------------
--  `profit_threshold`
----------------------------------------------------
ALTER TABLE `profit_threshold` ADD `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `profit_threshold` ADD `threshold` bigint NOT NULL COMMENT 'threshold' AFTER `id`;
ALTER TABLE `profit_threshold` ADD `admin_id` bigint unsigned NOT NULL COMMENT 'dividend_ratio' AFTER `threshold`;
ALTER TABLE `profit_threshold` ADD `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time' AFTER `admin_id`;
ALTER TABLE `profit_threshold` ADD `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time' AFTER `ctime`;

----------------------------------------------------
--  `role`
----------------------------------------------------
ALTER TABLE `role` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'role id';
ALTER TABLE `role` ADD `name` varchar(100) NOT NULL COMMENT 'role name' AFTER `id`;
ALTER TABLE `role` ADD `desc` varchar(100) NOT NULL COMMENT 'role description' AFTER `name`;
ALTER TABLE `role` ADD `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time' AFTER `desc`;
ALTER TABLE `role` ADD `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time' AFTER `ctime`;

----------------------------------------------------
--  `role_admin`
----------------------------------------------------
ALTER TABLE `role_admin` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `role_admin` ADD `roleid` bigint unsigned NOT NULL COMMENT 'role_id' AFTER `id`;
ALTER TABLE `role_admin` ADD `adminid` bigint unsigned NOT NULL COMMENT 'admin_id' AFTER `roleid`;
ALTER TABLE `role_admin` ADD `granted_by` varchar(100) NOT NULL COMMENT 'role granted_by' AFTER `adminid`;
ALTER TABLE `role_admin` ADD `granted_at` bigint NOT NULL DEFAULT 0 COMMENT 'granted_at time' AFTER `granted_by`;

----------------------------------------------------
--  `role_permission`
----------------------------------------------------
ALTER TABLE `role_permission` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `role_permission` ADD `roleid` bigint unsigned NOT NULL COMMENT 'roleid' AFTER `id`;
ALTER TABLE `role_permission` ADD `permissionid` bigint unsigned NOT NULL COMMENT 'pemissionid' AFTER `roleid`;

----------------------------------------------------
--  `server_node`
----------------------------------------------------
ALTER TABLE `server_node` ADD `app_name` varchar(256) NOT NULL COMMENT 'application name';
ALTER TABLE `server_node` ADD `region_id` bigint NOT NULL COMMENT 'region id' AFTER `app_name`;
ALTER TABLE `server_node` ADD `server_id` bigint NOT NULL COMMENT 'server id' AFTER `region_id`;
ALTER TABLE `server_node` ADD `last_ping` int unsigned NOT NULL COMMENT 'last ping timestamp' AFTER `server_id`;

----------------------------------------------------
--  `smscodes`
----------------------------------------------------
ALTER TABLE `smscodes` ADD `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `smscodes` ADD `national_code` varchar(32) NOT NULL DEFAULT '' COMMENT 'nationalCode' AFTER `id`;
ALTER TABLE `smscodes` ADD `mobile` varchar(32) NOT NULL DEFAULT '' COMMENT 'mobile' AFTER `national_code`;
ALTER TABLE `smscodes` ADD `action` varchar(100) NOT NULL DEFAULT '' COMMENT 'action' AFTER `mobile`;
ALTER TABLE `smscodes` ADD `code` varchar(16) NOT NULL DEFAULT '' COMMENT 'code' AFTER `action`;
ALTER TABLE `smscodes` ADD `status` tinyint NOT NULL COMMENT 'status' AFTER `code`;
ALTER TABLE `smscodes` ADD `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time' AFTER `status`;
ALTER TABLE `smscodes` ADD `etime` bigint NOT NULL DEFAULT 0 COMMENT 'expired time' AFTER `ctime`;

----------------------------------------------------
--  `smstemplates`
----------------------------------------------------
ALTER TABLE `smstemplates` ADD `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `smstemplates` ADD `name` varchar(100) NOT NULL DEFAULT '' COMMENT 'name' AFTER `id`;
ALTER TABLE `smstemplates` ADD `type` tinyint NOT NULL COMMENT 'type' AFTER `name`;
ALTER TABLE `smstemplates` ADD `template` varchar(256) NOT NULL DEFAULT '' COMMENT 'template' AFTER `type`;
ALTER TABLE `smstemplates` ADD `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time' AFTER `template`;
ALTER TABLE `smstemplates` ADD `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time' AFTER `ctime`;

----------------------------------------------------
--  `sys_msg`
----------------------------------------------------
ALTER TABLE `sys_msg` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `sys_msg` ADD `key` varchar(200) NOT NULL COMMENT 'System Message Key' AFTER `id`;
ALTER TABLE `sys_msg` ADD `buyer` varchar(400) NOT NULL COMMENT 'Buyer Show' AFTER `key`;
ALTER TABLE `sys_msg` ADD `seller` varchar(400) NOT NULL COMMENT 'Seller Show' AFTER `buyer`;
ALTER TABLE `sys_msg` ADD `admin` varchar(400) NOT NULL COMMENT 'Admin Show' AFTER `seller`;
ALTER TABLE `sys_msg` ADD `ctime` bigint NULL COMMENT '' AFTER `admin`;
ALTER TABLE `sys_msg` ADD `utime` bigint NULL COMMENT '' AFTER `ctime`;
ALTER TABLE `sys_msg` ADD `dtime` bigint NULL COMMENT '' AFTER `utime`;

----------------------------------------------------
--  `sys_notification`
----------------------------------------------------
ALTER TABLE `sys_notification` ADD `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `sys_notification` ADD `content` varchar(128) NOT NULL COMMENT 'content' AFTER `id`;
ALTER TABLE `sys_notification` ADD `admin_id` int unsigned NOT NULL COMMENT 'admin_id' AFTER `content`;
ALTER TABLE `sys_notification` ADD `status` tinyint NOT NULL COMMENT 'status' AFTER `admin_id`;
ALTER TABLE `sys_notification` ADD `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time' AFTER `status`;
ALTER TABLE `sys_notification` ADD `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time' AFTER `ctime`;

----------------------------------------------------
--  `task`
----------------------------------------------------
ALTER TABLE `task` ADD `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `task` ADD `name` varchar(256) NOT NULL COMMENT 'task name' AFTER `id`;
ALTER TABLE `task` ADD `alia` varchar(256) NOT NULL COMMENT 'task alia' AFTER `name`;
ALTER TABLE `task` ADD `app_name` varchar(256) NOT NULL COMMENT 'application name' AFTER `alia`;
ALTER TABLE `task` ADD `func_name` varchar(256) NOT NULL COMMENT 'task function name' AFTER `app_name`;
ALTER TABLE `task` ADD `spec` varchar(256) NOT NULL COMMENT 'task spec string' AFTER `func_name`;
ALTER TABLE `task` ADD `status` tinyint unsigned NOT NULL COMMENT 'status' AFTER `spec`;
ALTER TABLE `task` ADD `ctime` int unsigned NOT NULL COMMENT 'create time' AFTER `status`;
ALTER TABLE `task` ADD `utime` int unsigned NOT NULL COMMENT 'update time' AFTER `ctime`;
ALTER TABLE `task` ADD `desc` varchar(256) NOT NULL COMMENT 'task detail string' AFTER `utime`;

----------------------------------------------------
--  `task_result`
----------------------------------------------------
ALTER TABLE `task_result` ADD `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `task_result` ADD `app_name` varchar(256) NOT NULL COMMENT 'app name' AFTER `id`;
ALTER TABLE `task_result` ADD `region_id` bigint NOT NULL COMMENT 'region id' AFTER `app_name`;
ALTER TABLE `task_result` ADD `server_id` bigint NOT NULL COMMENT 'server id' AFTER `region_id`;
ALTER TABLE `task_result` ADD `name` varchar(256) NOT NULL COMMENT 'name' AFTER `server_id`;
ALTER TABLE `task_result` ADD `code` int NOT NULL COMMENT 'result code' AFTER `name`;
ALTER TABLE `task_result` ADD `detail` varchar(256) NOT NULL COMMENT 'task result detail' AFTER `code`;
ALTER TABLE `task_result` ADD `end_time` int unsigned NOT NULL COMMENT 'end time' AFTER `detail`;
ALTER TABLE `task_result` ADD `begin_time` int unsigned NOT NULL COMMENT 'begin time' AFTER `end_time`;
ALTER TABLE `task_result` ADD `ctime` int unsigned NOT NULL COMMENT 'create time' AFTER `begin_time`;

----------------------------------------------------
--  `top_agent`
----------------------------------------------------
ALTER TABLE `top_agent` ADD `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `top_agent` ADD `national_code` varchar(16) NOT NULL COMMENT 'national_code' AFTER `id`;
ALTER TABLE `top_agent` ADD `mobile` varchar(32) NOT NULL COMMENT 'mobile' AFTER `national_code`;
ALTER TABLE `top_agent` ADD `status` tinyint NOT NULL COMMENT 'status' AFTER `mobile`;
ALTER TABLE `top_agent` ADD `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time' AFTER `status`;
ALTER TABLE `top_agent` ADD `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time' AFTER `ctime`;

