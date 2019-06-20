----------------------------------------------------
--  `activity_user_conf`
----------------------------------------------------
ALTER TABLE `activity_user_conf` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `activity_user_conf` CHANGE `play_game_day` `play_game_day` int NOT NULL COMMENT '玩了并且有投注的天数 配置';
ALTER TABLE `activity_user_conf` CHANGE `bet_amount` `bet_amount` bigint NOT NULL COMMENT '投注额配置';
ALTER TABLE `activity_user_conf` CHANGE `effective_bet_amount` `effective_bet_amount` bigint NOT NULL COMMENT '累计有效投注额(本人加无限下级)的配置';

----------------------------------------------------
--  `admin_user`
----------------------------------------------------
ALTER TABLE `admin_user` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'user id';
ALTER TABLE `admin_user` CHANGE `name` `name` varchar(100) NOT NULL COMMENT 'user name';
ALTER TABLE `admin_user` CHANGE `email` `email` varchar(100) NOT NULL COMMENT 'user email';
ALTER TABLE `admin_user` CHANGE `status` `status` tinyint NOT NULL COMMENT 'status';
ALTER TABLE `admin_user` CHANGE `ctime` `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time';
ALTER TABLE `admin_user` CHANGE `utime` `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time';
ALTER TABLE `admin_user` CHANGE `dtime` `dtime` bigint NOT NULL DEFAULT 0 COMMENT 'forbit time';
ALTER TABLE `admin_user` CHANGE `login_time` `login_time` bigint NOT NULL DEFAULT 0 COMMENT 'last login time';
ALTER TABLE `admin_user` CHANGE `pwd` `pwd` varchar(100) NOT NULL DEFAULT '' COMMENT 'password';
ALTER TABLE `admin_user` CHANGE `whitelist_ips` `whitelist_ips` varchar(256) NOT NULL DEFAULT '' COMMENT 'whitelist_ips';
ALTER TABLE `admin_user` CHANGE `is_bind` `is_bind` bool NOT NULL DEFAULT false COMMENT 'is_bind';
ALTER TABLE `admin_user` CHANGE `secret_id` `secret_id` varchar(128) NOT NULL DEFAULT '' COMMENT 'secret_id';
ALTER TABLE `admin_user` CHANGE `qr_code` `qr_code` text NOT NULL COMMENT 'qr_code';

----------------------------------------------------
--  `agent_white_list`
----------------------------------------------------
ALTER TABLE `agent_white_list` CHANGE `id` `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `agent_white_list` CHANGE `name` `name` varchar(64) NOT NULL COMMENT 'name';
ALTER TABLE `agent_white_list` CHANGE `commission` `commission` int NOT NULL COMMENT 'commission';
ALTER TABLE `agent_white_list` CHANGE `precision` `precision` int NOT NULL COMMENT 'precision';
ALTER TABLE `agent_white_list` CHANGE `ctime` `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time';
ALTER TABLE `agent_white_list` CHANGE `utime` `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time';

----------------------------------------------------
--  `announcement`
----------------------------------------------------
ALTER TABLE `announcement` CHANGE `id` `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `announcement` CHANGE `type` `type` tinyint NOT NULL COMMENT 'type';
ALTER TABLE `announcement` CHANGE `title` `title` varchar(100) NOT NULL COMMENT 'title';
ALTER TABLE `announcement` CHANGE `content` `content` varchar(500) NOT NULL COMMENT 'content';
ALTER TABLE `announcement` CHANGE `stime` `stime` bigint NOT NULL DEFAULT 0 COMMENT 'start time';
ALTER TABLE `announcement` CHANGE `etime` `etime` bigint NOT NULL DEFAULT 0 COMMENT 'end time';
ALTER TABLE `announcement` CHANGE `ctime` `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time';
ALTER TABLE `announcement` CHANGE `utime` `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time';

----------------------------------------------------
--  `app_channel`
----------------------------------------------------
ALTER TABLE `app_channel` CHANGE `id` `id` int unsigned NOT NULL COMMENT 'id';
ALTER TABLE `app_channel` CHANGE `is_third_hall` `is_third_hall` tinyint NOT NULL COMMENT 'is_third_hall';
ALTER TABLE `app_channel` CHANGE `name` `name` varchar(256) NOT NULL COMMENT 'name';
ALTER TABLE `app_channel` CHANGE `desc` `desc` varchar(256) NOT NULL COMMENT 'desc';
ALTER TABLE `app_channel` CHANGE `exchangeRate` `exchangeRate` int NOT NULL DEFAULT 0 COMMENT 'exchange rate';
ALTER TABLE `app_channel` CHANGE `precision` `precision` int NOT NULL DEFAULT 0 COMMENT 'precision';
ALTER TABLE `app_channel` CHANGE `profit_rate` `profit_rate` int NOT NULL DEFAULT 0 COMMENT 'profit_rate';
ALTER TABLE `app_channel` CHANGE `icon_url` `icon_url` varchar(128) NOT NULL COMMENT 'icon_url';
ALTER TABLE `app_channel` CHANGE `ctime` `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time';
ALTER TABLE `app_channel` CHANGE `utime` `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time';

----------------------------------------------------
--  `app_type`
----------------------------------------------------
ALTER TABLE `app_type` CHANGE `id` `id` int unsigned NOT NULL COMMENT 'id';
ALTER TABLE `app_type` CHANGE `name` `name` varchar(256) NOT NULL COMMENT 'name';
ALTER TABLE `app_type` CHANGE `desc` `desc` varchar(256) NOT NULL COMMENT 'desc';
ALTER TABLE `app_type` CHANGE `ctime` `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time';
ALTER TABLE `app_type` CHANGE `utime` `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time';

----------------------------------------------------
--  `app_version`
----------------------------------------------------
ALTER TABLE `app_version` CHANGE `id` `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `app_version` CHANGE `version` `version` varchar(100) NOT NULL COMMENT 'app version show';
ALTER TABLE `app_version` CHANGE `version_num` `version_num` int NOT NULL COMMENT 'app version num';
ALTER TABLE `app_version` CHANGE `changelog` `changelog` varchar(300) NOT NULL COMMENT 'change log';
ALTER TABLE `app_version` CHANGE `download` `download` varchar(300) NOT NULL COMMENT 'download url';
ALTER TABLE `app_version` CHANGE `system` `system` tinyint NOT NULL COMMENT 'system type';
ALTER TABLE `app_version` CHANGE `status` `status` tinyint NOT NULL COMMENT 'app version status';
ALTER TABLE `app_version` CHANGE `ctime` `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time';
ALTER TABLE `app_version` CHANGE `utime` `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time';
ALTER TABLE `app_version` CHANGE `dtime` `dtime` bigint NOT NULL DEFAULT 0 COMMENT 'update dtime';

----------------------------------------------------
--  `app_whitelist`
----------------------------------------------------
ALTER TABLE `app_whitelist` CHANGE `id` `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `app_whitelist` CHANGE `channel_id` `channel_id` int unsigned NOT NULL COMMENT 'channel_id';
ALTER TABLE `app_whitelist` CHANGE `app_id` `app_id` varchar(16) NOT NULL COMMENT 'app_id';
ALTER TABLE `app_whitelist` CHANGE `ctime` `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time';

----------------------------------------------------
--  `appeal_service`
----------------------------------------------------
ALTER TABLE `appeal_service` CHANGE `id` `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `appeal_service` CHANGE `admin_id` `admin_id` int unsigned NOT NULL COMMENT 'admin_id';
ALTER TABLE `appeal_service` CHANGE `wechat` `wechat` varchar(32) NOT NULL COMMENT 'wechat';
ALTER TABLE `appeal_service` CHANGE `qr_code` `qr_code` varchar(300) NOT NULL COMMENT '';
ALTER TABLE `appeal_service` CHANGE `status` `status` tinyint NOT NULL COMMENT 'status';

----------------------------------------------------
--  `apps`
----------------------------------------------------
ALTER TABLE `apps` CHANGE `id` `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `apps` CHANGE `position` `position` int unsigned NOT NULL COMMENT 'position';
ALTER TABLE `apps` CHANGE `name` `name` varchar(64) NOT NULL COMMENT 'name';
ALTER TABLE `apps` CHANGE `desc` `desc` varchar(256) NOT NULL COMMENT 'desc';
ALTER TABLE `apps` CHANGE `url` `url` varchar(128) NOT NULL COMMENT 'url';
ALTER TABLE `apps` CHANGE `icon_url` `icon_url` varchar(128) NOT NULL COMMENT 'icon_url';
ALTER TABLE `apps` CHANGE `type_id` `type_id` tinyint NOT NULL COMMENT 'type_id';
ALTER TABLE `apps` CHANGE `channel_id` `channel_id` int unsigned NOT NULL COMMENT 'channel_id';
ALTER TABLE `apps` CHANGE `app_id` `app_id` varchar(50) NOT NULL COMMENT 'app_id';
ALTER TABLE `apps` CHANGE `featured` `featured` tinyint NOT NULL COMMENT 'featured';
ALTER TABLE `apps` CHANGE `status` `status` tinyint NOT NULL COMMENT 'status';
ALTER TABLE `apps` CHANGE `orientation` `orientation` tinyint NOT NULL DEFAULT 1 COMMENT 'orientation';
ALTER TABLE `apps` CHANGE `ctime` `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time';
ALTER TABLE `apps` CHANGE `utime` `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time';

----------------------------------------------------
--  `banner`
----------------------------------------------------
ALTER TABLE `banner` CHANGE `id` `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `banner` CHANGE `subject` `subject` varchar(256) NOT NULL COMMENT 'subject';
ALTER TABLE `banner` CHANGE `image` `image` varchar(256) NOT NULL COMMENT 'image';
ALTER TABLE `banner` CHANGE `url` `url` varchar(256) NOT NULL DEFAULT '' COMMENT 'image';
ALTER TABLE `banner` CHANGE `status` `status` tinyint NOT NULL COMMENT 'status';
ALTER TABLE `banner` CHANGE `ctime` `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time';
ALTER TABLE `banner` CHANGE `utime` `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time';
ALTER TABLE `banner` CHANGE `stime` `stime` bigint NOT NULL DEFAULT 0 COMMENT 'start time';
ALTER TABLE `banner` CHANGE `etime` `etime` bigint NOT NULL DEFAULT 0 COMMENT 'end time';

----------------------------------------------------
--  `commissionrates`
----------------------------------------------------
ALTER TABLE `commissionrates` CHANGE `id` `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `commissionrates` CHANGE `min` `min` bigint unsigned NOT NULL COMMENT 'min';
ALTER TABLE `commissionrates` CHANGE `max` `max` bigint unsigned NOT NULL COMMENT 'max';
ALTER TABLE `commissionrates` CHANGE `commission` `commission` int NOT NULL DEFAULT 0 COMMENT 'commission';
ALTER TABLE `commissionrates` CHANGE `precision` `precision` int NOT NULL DEFAULT 0 COMMENT 'precision';
ALTER TABLE `commissionrates` CHANGE `ctime` `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time';
ALTER TABLE `commissionrates` CHANGE `utime` `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time';

----------------------------------------------------
--  `config`
----------------------------------------------------
ALTER TABLE `config` CHANGE `id` `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `config` CHANGE `action` `action` tinyint NOT NULL COMMENT 'action';
ALTER TABLE `config` CHANGE `key` `key` varchar(256) NOT NULL COMMENT 'key';
ALTER TABLE `config` CHANGE `value` `value` text NOT NULL COMMENT 'value';
ALTER TABLE `config` CHANGE `desc` `desc` varchar(256) NOT NULL COMMENT 'descripe';
ALTER TABLE `config` CHANGE `ctime` `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time';

----------------------------------------------------
--  `config_warning`
----------------------------------------------------
ALTER TABLE `config_warning` CHANGE `id` `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `config_warning` CHANGE `type` `type` tinyint NOT NULL COMMENT 'type';
ALTER TABLE `config_warning` CHANGE `national_code` `national_code` varchar(16) NOT NULL COMMENT 'national_code';
ALTER TABLE `config_warning` CHANGE `mobile` `mobile` varchar(32) NOT NULL COMMENT 'mobile';
ALTER TABLE `config_warning` CHANGE `sms_type` `sms_type` tinyint NOT NULL COMMENT 'sms_type';

----------------------------------------------------
--  `endpoint`
----------------------------------------------------
ALTER TABLE `endpoint` CHANGE `id` `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `endpoint` CHANGE `endpoint` `endpoint` varchar(100) NOT NULL COMMENT 'endpoint';
ALTER TABLE `endpoint` CHANGE `ctime` `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time';
ALTER TABLE `endpoint` CHANGE `utime` `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time';

----------------------------------------------------
--  `ip_white_list`
----------------------------------------------------
ALTER TABLE `ip_white_list` CHANGE `id` `id` int unsigned NOT NULL COMMENT 'id';
ALTER TABLE `ip_white_list` CHANGE `ip` `ip` varchar(0) NOT NULL COMMENT 'ip address';

----------------------------------------------------
--  `menu_access`
----------------------------------------------------
ALTER TABLE `menu_access` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `menu_access` CHANGE `role_id` `role_id` bigint unsigned NOT NULL COMMENT '角色id';
ALTER TABLE `menu_access` CHANGE `menu_id` `menu_id` bigint unsigned NOT NULL COMMENT '菜单id';

----------------------------------------------------
--  `menu_conf`
----------------------------------------------------
ALTER TABLE `menu_conf` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `menu_conf` CHANGE `pid` `pid` bigint unsigned NOT NULL COMMENT '父节点的id，一级菜单为0';
ALTER TABLE `menu_conf` CHANGE `level` `level` int NOT NULL COMMENT '层级,1是一级菜单,2是二级菜单';
ALTER TABLE `menu_conf` CHANGE `name` `name` varchar(100) NOT NULL COMMENT '菜单名称';
ALTER TABLE `menu_conf` CHANGE `path` `path` varchar(100) NOT NULL COMMENT '菜单路径';
ALTER TABLE `menu_conf` CHANGE `icon` `icon` varchar(100) NOT NULL COMMENT '一级菜单图标';
ALTER TABLE `menu_conf` CHANGE `hide_in_menu` `hide_in_menu` bool NOT NULL COMMENT '是否隐藏';
ALTER TABLE `menu_conf` CHANGE `component` `component` varchar(100) NOT NULL COMMENT '组件';
ALTER TABLE `menu_conf` CHANGE `order_id` `order_id` int unsigned NOT NULL COMMENT '排序id';
ALTER TABLE `menu_conf` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT '创建时间';
ALTER TABLE `menu_conf` CHANGE `utime` `utime` bigint NOT NULL COMMENT '修改时间';

----------------------------------------------------
--  `month_dividend_position_conf`
----------------------------------------------------
ALTER TABLE `month_dividend_position_conf` CHANGE `id` `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `month_dividend_position_conf` CHANGE `agent_lv` `agent_lv` int NOT NULL COMMENT 'agent_lv';
ALTER TABLE `month_dividend_position_conf` CHANGE `position` `position` int NOT NULL COMMENT 'position';
ALTER TABLE `month_dividend_position_conf` CHANGE `min` `min` bigint NOT NULL COMMENT 'min';
ALTER TABLE `month_dividend_position_conf` CHANGE `max` `max` bigint NOT NULL COMMENT 'max';
ALTER TABLE `month_dividend_position_conf` CHANGE `activity_num` `activity_num` int NOT NULL COMMENT 'activity num';
ALTER TABLE `month_dividend_position_conf` CHANGE `dividend_ratio` `dividend_ratio` int NOT NULL COMMENT 'dividend_ratio num';
ALTER TABLE `month_dividend_position_conf` CHANGE `ctime` `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time';
ALTER TABLE `month_dividend_position_conf` CHANGE `utime` `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time';

----------------------------------------------------
--  `month_dividend_white_list`
----------------------------------------------------
ALTER TABLE `month_dividend_white_list` CHANGE `id` `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `month_dividend_white_list` CHANGE `name` `name` varchar(64) NOT NULL COMMENT 'name';
ALTER TABLE `month_dividend_white_list` CHANGE `dividend_ratio` `dividend_ratio` int NOT NULL COMMENT 'dividend_ratio';
ALTER TABLE `month_dividend_white_list` CHANGE `ctime` `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time';
ALTER TABLE `month_dividend_white_list` CHANGE `utime` `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time';

----------------------------------------------------
--  `operation_log`
----------------------------------------------------
ALTER TABLE `operation_log` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `operation_log` CHANGE `admin_id` `admin_id` bigint unsigned NOT NULL COMMENT 'admin_id';
ALTER TABLE `operation_log` CHANGE `method` `method` varchar(100) NOT NULL DEFAULT '' COMMENT 'req method';
ALTER TABLE `operation_log` CHANGE `route` `route` varchar(100) NOT NULL DEFAULT '' COMMENT 'route ';
ALTER TABLE `operation_log` CHANGE `action` `action` int NOT NULL DEFAULT 0 COMMENT 'action';
ALTER TABLE `operation_log` CHANGE `input` `input` varchar(65535) NOT NULL DEFAULT '' COMMENT 'input';
ALTER TABLE `operation_log` CHANGE `user_agent` `user_agent` varchar(512) NOT NULL DEFAULT '' COMMENT 'user agent';
ALTER TABLE `operation_log` CHANGE `ips` `ips` varchar(100) NOT NULL DEFAULT '' COMMENT 'ip';
ALTER TABLE `operation_log` CHANGE `response_code` `response_code` int NOT NULL DEFAULT 0 COMMENT 'response code';
ALTER TABLE `operation_log` CHANGE `ctime` `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time';

----------------------------------------------------
--  `otc_stat`
----------------------------------------------------
ALTER TABLE `otc_stat` CHANGE `id` `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `otc_stat` CHANGE `date` `date` int unsigned NOT NULL COMMENT 'date';
ALTER TABLE `otc_stat` CHANGE `num_login` `num_login` int unsigned NOT NULL DEFAULT 0 COMMENT '';
ALTER TABLE `otc_stat` CHANGE `num_user_new` `num_user_new` int unsigned NOT NULL DEFAULT 0 COMMENT '';
ALTER TABLE `otc_stat` CHANGE `num_order` `num_order` int unsigned NOT NULL DEFAULT 0 COMMENT 'date';
ALTER TABLE `otc_stat` CHANGE `num_order_deal` `num_order_deal` int unsigned NOT NULL DEFAULT 0 COMMENT '';
ALTER TABLE `otc_stat` CHANGE `num_order_buy` `num_order_buy` int unsigned NOT NULL DEFAULT 0 COMMENT '';
ALTER TABLE `otc_stat` CHANGE `num_order_sell` `num_order_sell` int unsigned NOT NULL DEFAULT 0 COMMENT '';
ALTER TABLE `otc_stat` CHANGE `num_funds` `num_funds` int unsigned NOT NULL DEFAULT 0 COMMENT '';
ALTER TABLE `otc_stat` CHANGE `num_amount` `num_amount` int unsigned NOT NULL DEFAULT 0 COMMENT '';
ALTER TABLE `otc_stat` CHANGE `num_amount_buy` `num_amount_buy` int unsigned NOT NULL DEFAULT 0 COMMENT '';
ALTER TABLE `otc_stat` CHANGE `num_amount_sell` `num_amount_sell` int unsigned NOT NULL DEFAULT 0 COMMENT '';
ALTER TABLE `otc_stat` CHANGE `num_fee_buy` `num_fee_buy` int unsigned NOT NULL DEFAULT 0 COMMENT '';
ALTER TABLE `otc_stat` CHANGE `num_fee_sell` `num_fee_sell` int unsigned NOT NULL DEFAULT 0 COMMENT '';
ALTER TABLE `otc_stat` CHANGE `game_recharge` `game_recharge` int unsigned NOT NULL DEFAULT 0 COMMENT '';
ALTER TABLE `otc_stat` CHANGE `game_withdrawal` `game_withdrawal` int unsigned NOT NULL DEFAULT 0 COMMENT '';
ALTER TABLE `otc_stat` CHANGE `usdt_recharge` `usdt_recharge` int unsigned NOT NULL DEFAULT 0 COMMENT '';
ALTER TABLE `otc_stat` CHANGE `usdt_withdrawal` `usdt_withdrawal` int unsigned NOT NULL DEFAULT 0 COMMENT '';
ALTER TABLE `otc_stat` CHANGE `usdt_fee` `usdt_fee` int unsigned NOT NULL DEFAULT 0 COMMENT '';

----------------------------------------------------
--  `otc_stat_all_people`
----------------------------------------------------
ALTER TABLE `otc_stat_all_people` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `otc_stat_all_people` CHANGE `buy_order` `buy_order` int unsigned NOT NULL DEFAULT 0 COMMENT '';
ALTER TABLE `otc_stat_all_people` CHANGE `sell_order` `sell_order` int unsigned NOT NULL DEFAULT 0 COMMENT '';
ALTER TABLE `otc_stat_all_people` CHANGE `buy_eusd` `buy_eusd` int unsigned NOT NULL DEFAULT 0 COMMENT '';
ALTER TABLE `otc_stat_all_people` CHANGE `sell_eusd` `sell_eusd` int unsigned NOT NULL DEFAULT 0 COMMENT '';
ALTER TABLE `otc_stat_all_people` CHANGE `usdt_recharge` `usdt_recharge` int unsigned NOT NULL DEFAULT 0 COMMENT '';
ALTER TABLE `otc_stat_all_people` CHANGE `usdt_withdrawal` `usdt_withdrawal` int unsigned NOT NULL DEFAULT 0 COMMENT '';

----------------------------------------------------
--  `permission`
----------------------------------------------------
ALTER TABLE `permission` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `permission` CHANGE `slug` `slug` varchar(100) NOT NULL COMMENT 'permission_name';
ALTER TABLE `permission` CHANGE `desc` `desc` varchar(100) NOT NULL COMMENT 'permission_desc';
ALTER TABLE `permission` CHANGE `ctime` `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time';
ALTER TABLE `permission` CHANGE `utime` `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time';
ALTER TABLE `permission` CHANGE `dtime` `dtime` bigint NOT NULL DEFAULT 0 COMMENT 'D time';

----------------------------------------------------
--  `profit_threshold`
----------------------------------------------------
ALTER TABLE `profit_threshold` CHANGE `id` `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `profit_threshold` CHANGE `threshold` `threshold` bigint NOT NULL COMMENT 'threshold';
ALTER TABLE `profit_threshold` CHANGE `admin_id` `admin_id` bigint unsigned NOT NULL COMMENT 'dividend_ratio';
ALTER TABLE `profit_threshold` CHANGE `ctime` `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time';
ALTER TABLE `profit_threshold` CHANGE `utime` `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time';

----------------------------------------------------
--  `role`
----------------------------------------------------
ALTER TABLE `role` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'role id';
ALTER TABLE `role` CHANGE `name` `name` varchar(100) NOT NULL COMMENT 'role name';
ALTER TABLE `role` CHANGE `desc` `desc` varchar(100) NOT NULL COMMENT 'role description';
ALTER TABLE `role` CHANGE `ctime` `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time';
ALTER TABLE `role` CHANGE `utime` `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time';

----------------------------------------------------
--  `role_admin`
----------------------------------------------------
ALTER TABLE `role_admin` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `role_admin` CHANGE `roleid` `roleid` bigint unsigned NOT NULL COMMENT 'role_id';
ALTER TABLE `role_admin` CHANGE `adminid` `adminid` bigint unsigned NOT NULL COMMENT 'admin_id';
ALTER TABLE `role_admin` CHANGE `granted_by` `granted_by` varchar(100) NOT NULL COMMENT 'role granted_by';
ALTER TABLE `role_admin` CHANGE `granted_at` `granted_at` bigint NOT NULL DEFAULT 0 COMMENT 'granted_at time';

----------------------------------------------------
--  `role_permission`
----------------------------------------------------
ALTER TABLE `role_permission` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `role_permission` CHANGE `roleid` `roleid` bigint unsigned NOT NULL COMMENT 'roleid';
ALTER TABLE `role_permission` CHANGE `permissionid` `permissionid` bigint unsigned NOT NULL COMMENT 'pemissionid';

----------------------------------------------------
--  `server_node`
----------------------------------------------------
ALTER TABLE `server_node` CHANGE `app_name` `app_name` varchar(256) NOT NULL COMMENT 'application name';
ALTER TABLE `server_node` CHANGE `region_id` `region_id` bigint NOT NULL COMMENT 'region id';
ALTER TABLE `server_node` CHANGE `server_id` `server_id` bigint NOT NULL COMMENT 'server id';
ALTER TABLE `server_node` CHANGE `last_ping` `last_ping` int unsigned NOT NULL COMMENT 'last ping timestamp';

----------------------------------------------------
--  `smscodes`
----------------------------------------------------
ALTER TABLE `smscodes` CHANGE `id` `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `smscodes` CHANGE `national_code` `national_code` varchar(32) NOT NULL DEFAULT '' COMMENT 'nationalCode';
ALTER TABLE `smscodes` CHANGE `mobile` `mobile` varchar(32) NOT NULL DEFAULT '' COMMENT 'mobile';
ALTER TABLE `smscodes` CHANGE `action` `action` varchar(100) NOT NULL DEFAULT '' COMMENT 'action';
ALTER TABLE `smscodes` CHANGE `code` `code` varchar(16) NOT NULL DEFAULT '' COMMENT 'code';
ALTER TABLE `smscodes` CHANGE `status` `status` tinyint NOT NULL COMMENT 'status';
ALTER TABLE `smscodes` CHANGE `ctime` `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time';
ALTER TABLE `smscodes` CHANGE `etime` `etime` bigint NOT NULL DEFAULT 0 COMMENT 'expired time';

----------------------------------------------------
--  `smstemplates`
----------------------------------------------------
ALTER TABLE `smstemplates` CHANGE `id` `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `smstemplates` CHANGE `name` `name` varchar(100) NOT NULL DEFAULT '' COMMENT 'name';
ALTER TABLE `smstemplates` CHANGE `type` `type` tinyint NOT NULL COMMENT 'type';
ALTER TABLE `smstemplates` CHANGE `template` `template` varchar(256) NOT NULL DEFAULT '' COMMENT 'template';
ALTER TABLE `smstemplates` CHANGE `ctime` `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time';
ALTER TABLE `smstemplates` CHANGE `utime` `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time';

----------------------------------------------------
--  `sys_msg`
----------------------------------------------------
ALTER TABLE `sys_msg` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `sys_msg` CHANGE `key` `key` varchar(200) NOT NULL COMMENT 'System Message Key';
ALTER TABLE `sys_msg` CHANGE `buyer` `buyer` varchar(400) NOT NULL COMMENT 'Buyer Show';
ALTER TABLE `sys_msg` CHANGE `seller` `seller` varchar(400) NOT NULL COMMENT 'Seller Show';
ALTER TABLE `sys_msg` CHANGE `admin` `admin` varchar(400) NOT NULL COMMENT 'Admin Show';
ALTER TABLE `sys_msg` CHANGE `ctime` `ctime` bigint NULL COMMENT '';
ALTER TABLE `sys_msg` CHANGE `utime` `utime` bigint NULL COMMENT '';
ALTER TABLE `sys_msg` CHANGE `dtime` `dtime` bigint NULL COMMENT '';

----------------------------------------------------
--  `sys_notification`
----------------------------------------------------
ALTER TABLE `sys_notification` CHANGE `id` `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `sys_notification` CHANGE `content` `content` varchar(128) NOT NULL COMMENT 'content';
ALTER TABLE `sys_notification` CHANGE `admin_id` `admin_id` int unsigned NOT NULL COMMENT 'admin_id';
ALTER TABLE `sys_notification` CHANGE `status` `status` tinyint NOT NULL COMMENT 'status';
ALTER TABLE `sys_notification` CHANGE `ctime` `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time';
ALTER TABLE `sys_notification` CHANGE `utime` `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time';

----------------------------------------------------
--  `task`
----------------------------------------------------
ALTER TABLE `task` CHANGE `id` `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `task` CHANGE `name` `name` varchar(256) NOT NULL COMMENT 'task name';
ALTER TABLE `task` CHANGE `alia` `alia` varchar(256) NOT NULL COMMENT 'task alia';
ALTER TABLE `task` CHANGE `app_name` `app_name` varchar(256) NOT NULL COMMENT 'application name';
ALTER TABLE `task` CHANGE `func_name` `func_name` varchar(256) NOT NULL COMMENT 'task function name';
ALTER TABLE `task` CHANGE `spec` `spec` varchar(256) NOT NULL COMMENT 'task spec string';
ALTER TABLE `task` CHANGE `status` `status` tinyint unsigned NOT NULL COMMENT 'status';
ALTER TABLE `task` CHANGE `ctime` `ctime` int unsigned NOT NULL COMMENT 'create time';
ALTER TABLE `task` CHANGE `utime` `utime` int unsigned NOT NULL COMMENT 'update time';
ALTER TABLE `task` CHANGE `desc` `desc` varchar(256) NOT NULL COMMENT 'task detail string';

----------------------------------------------------
--  `task_result`
----------------------------------------------------
ALTER TABLE `task_result` CHANGE `id` `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `task_result` CHANGE `app_name` `app_name` varchar(256) NOT NULL COMMENT 'app name';
ALTER TABLE `task_result` CHANGE `region_id` `region_id` bigint NOT NULL COMMENT 'region id';
ALTER TABLE `task_result` CHANGE `server_id` `server_id` bigint NOT NULL COMMENT 'server id';
ALTER TABLE `task_result` CHANGE `name` `name` varchar(256) NOT NULL COMMENT 'name';
ALTER TABLE `task_result` CHANGE `code` `code` int NOT NULL COMMENT 'result code';
ALTER TABLE `task_result` CHANGE `detail` `detail` varchar(256) NOT NULL COMMENT 'task result detail';
ALTER TABLE `task_result` CHANGE `end_time` `end_time` int unsigned NOT NULL COMMENT 'end time';
ALTER TABLE `task_result` CHANGE `begin_time` `begin_time` int unsigned NOT NULL COMMENT 'begin time';
ALTER TABLE `task_result` CHANGE `ctime` `ctime` int unsigned NOT NULL COMMENT 'create time';

----------------------------------------------------
--  `top_agent`
----------------------------------------------------
ALTER TABLE `top_agent` CHANGE `id` `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `top_agent` CHANGE `national_code` `national_code` varchar(16) NOT NULL COMMENT 'national_code';
ALTER TABLE `top_agent` CHANGE `mobile` `mobile` varchar(32) NOT NULL COMMENT 'mobile';
ALTER TABLE `top_agent` CHANGE `status` `status` tinyint NOT NULL COMMENT 'status';
ALTER TABLE `top_agent` CHANGE `ctime` `ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time';
ALTER TABLE `top_agent` CHANGE `utime` `utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time';

