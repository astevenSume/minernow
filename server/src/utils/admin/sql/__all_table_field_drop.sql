----------------------------------------------------
--  `activity_user_conf`
----------------------------------------------------
ALTER TABLE `activity_user_conf` DROP `id`;
ALTER TABLE `activity_user_conf` DROP `play_game_day`;
ALTER TABLE `activity_user_conf` DROP `bet_amount`;
ALTER TABLE `activity_user_conf` DROP `effective_bet_amount`;

----------------------------------------------------
--  `admin_user`
----------------------------------------------------
ALTER TABLE `admin_user` DROP `id`;
ALTER TABLE `admin_user` DROP `name`;
ALTER TABLE `admin_user` DROP `email`;
ALTER TABLE `admin_user` DROP `status`;
ALTER TABLE `admin_user` DROP `ctime`;
ALTER TABLE `admin_user` DROP `utime`;
ALTER TABLE `admin_user` DROP `dtime`;
ALTER TABLE `admin_user` DROP `login_time`;
ALTER TABLE `admin_user` DROP `pwd`;
ALTER TABLE `admin_user` DROP `whitelist_ips`;
ALTER TABLE `admin_user` DROP `is_bind`;
ALTER TABLE `admin_user` DROP `secret_id`;
ALTER TABLE `admin_user` DROP `qr_code`;

----------------------------------------------------
--  `agent_white_list`
----------------------------------------------------
ALTER TABLE `agent_white_list` DROP `id`;
ALTER TABLE `agent_white_list` DROP `name`;
ALTER TABLE `agent_white_list` DROP `commission`;
ALTER TABLE `agent_white_list` DROP `precision`;
ALTER TABLE `agent_white_list` DROP `ctime`;
ALTER TABLE `agent_white_list` DROP `utime`;

----------------------------------------------------
--  `announcement`
----------------------------------------------------
ALTER TABLE `announcement` DROP `id`;
ALTER TABLE `announcement` DROP `type`;
ALTER TABLE `announcement` DROP `title`;
ALTER TABLE `announcement` DROP `content`;
ALTER TABLE `announcement` DROP `stime`;
ALTER TABLE `announcement` DROP `etime`;
ALTER TABLE `announcement` DROP `ctime`;
ALTER TABLE `announcement` DROP `utime`;

----------------------------------------------------
--  `app_channel`
----------------------------------------------------
ALTER TABLE `app_channel` DROP `id`;
ALTER TABLE `app_channel` DROP `is_third_hall`;
ALTER TABLE `app_channel` DROP `name`;
ALTER TABLE `app_channel` DROP `desc`;
ALTER TABLE `app_channel` DROP `exchangeRate`;
ALTER TABLE `app_channel` DROP `precision`;
ALTER TABLE `app_channel` DROP `profit_rate`;
ALTER TABLE `app_channel` DROP `icon_url`;
ALTER TABLE `app_channel` DROP `ctime`;
ALTER TABLE `app_channel` DROP `utime`;

----------------------------------------------------
--  `app_type`
----------------------------------------------------
ALTER TABLE `app_type` DROP `id`;
ALTER TABLE `app_type` DROP `name`;
ALTER TABLE `app_type` DROP `desc`;
ALTER TABLE `app_type` DROP `ctime`;
ALTER TABLE `app_type` DROP `utime`;

----------------------------------------------------
--  `app_version`
----------------------------------------------------
ALTER TABLE `app_version` DROP `id`;
ALTER TABLE `app_version` DROP `version`;
ALTER TABLE `app_version` DROP `version_num`;
ALTER TABLE `app_version` DROP `changelog`;
ALTER TABLE `app_version` DROP `download`;
ALTER TABLE `app_version` DROP `system`;
ALTER TABLE `app_version` DROP `status`;
ALTER TABLE `app_version` DROP `ctime`;
ALTER TABLE `app_version` DROP `utime`;
ALTER TABLE `app_version` DROP `dtime`;

----------------------------------------------------
--  `app_whitelist`
----------------------------------------------------
ALTER TABLE `app_whitelist` DROP `id`;
ALTER TABLE `app_whitelist` DROP `channel_id`;
ALTER TABLE `app_whitelist` DROP `app_id`;
ALTER TABLE `app_whitelist` DROP `ctime`;

----------------------------------------------------
--  `appeal_service`
----------------------------------------------------
ALTER TABLE `appeal_service` DROP `id`;
ALTER TABLE `appeal_service` DROP `admin_id`;
ALTER TABLE `appeal_service` DROP `wechat`;
ALTER TABLE `appeal_service` DROP `qr_code`;
ALTER TABLE `appeal_service` DROP `status`;

----------------------------------------------------
--  `apps`
----------------------------------------------------
ALTER TABLE `apps` DROP `id`;
ALTER TABLE `apps` DROP `position`;
ALTER TABLE `apps` DROP `name`;
ALTER TABLE `apps` DROP `desc`;
ALTER TABLE `apps` DROP `url`;
ALTER TABLE `apps` DROP `icon_url`;
ALTER TABLE `apps` DROP `type_id`;
ALTER TABLE `apps` DROP `channel_id`;
ALTER TABLE `apps` DROP `app_id`;
ALTER TABLE `apps` DROP `featured`;
ALTER TABLE `apps` DROP `status`;
ALTER TABLE `apps` DROP `orientation`;
ALTER TABLE `apps` DROP `ctime`;
ALTER TABLE `apps` DROP `utime`;

----------------------------------------------------
--  `banner`
----------------------------------------------------
ALTER TABLE `banner` DROP `id`;
ALTER TABLE `banner` DROP `subject`;
ALTER TABLE `banner` DROP `image`;
ALTER TABLE `banner` DROP `url`;
ALTER TABLE `banner` DROP `status`;
ALTER TABLE `banner` DROP `ctime`;
ALTER TABLE `banner` DROP `utime`;
ALTER TABLE `banner` DROP `stime`;
ALTER TABLE `banner` DROP `etime`;

----------------------------------------------------
--  `commissionrates`
----------------------------------------------------
ALTER TABLE `commissionrates` DROP `id`;
ALTER TABLE `commissionrates` DROP `min`;
ALTER TABLE `commissionrates` DROP `max`;
ALTER TABLE `commissionrates` DROP `commission`;
ALTER TABLE `commissionrates` DROP `precision`;
ALTER TABLE `commissionrates` DROP `ctime`;
ALTER TABLE `commissionrates` DROP `utime`;

----------------------------------------------------
--  `config`
----------------------------------------------------
ALTER TABLE `config` DROP `id`;
ALTER TABLE `config` DROP `action`;
ALTER TABLE `config` DROP `key`;
ALTER TABLE `config` DROP `value`;
ALTER TABLE `config` DROP `desc`;
ALTER TABLE `config` DROP `ctime`;

----------------------------------------------------
--  `config_warning`
----------------------------------------------------
ALTER TABLE `config_warning` DROP `id`;
ALTER TABLE `config_warning` DROP `type`;
ALTER TABLE `config_warning` DROP `national_code`;
ALTER TABLE `config_warning` DROP `mobile`;
ALTER TABLE `config_warning` DROP `sms_type`;

----------------------------------------------------
--  `endpoint`
----------------------------------------------------
ALTER TABLE `endpoint` DROP `id`;
ALTER TABLE `endpoint` DROP `endpoint`;
ALTER TABLE `endpoint` DROP `ctime`;
ALTER TABLE `endpoint` DROP `utime`;

----------------------------------------------------
--  `ip_white_list`
----------------------------------------------------
ALTER TABLE `ip_white_list` DROP `id`;
ALTER TABLE `ip_white_list` DROP `ip`;

----------------------------------------------------
--  `menu_access`
----------------------------------------------------
ALTER TABLE `menu_access` DROP `id`;
ALTER TABLE `menu_access` DROP `role_id`;
ALTER TABLE `menu_access` DROP `menu_id`;

----------------------------------------------------
--  `menu_conf`
----------------------------------------------------
ALTER TABLE `menu_conf` DROP `id`;
ALTER TABLE `menu_conf` DROP `pid`;
ALTER TABLE `menu_conf` DROP `level`;
ALTER TABLE `menu_conf` DROP `name`;
ALTER TABLE `menu_conf` DROP `path`;
ALTER TABLE `menu_conf` DROP `icon`;
ALTER TABLE `menu_conf` DROP `hide_in_menu`;
ALTER TABLE `menu_conf` DROP `component`;
ALTER TABLE `menu_conf` DROP `order_id`;
ALTER TABLE `menu_conf` DROP `ctime`;
ALTER TABLE `menu_conf` DROP `utime`;

----------------------------------------------------
--  `month_dividend_position_conf`
----------------------------------------------------
ALTER TABLE `month_dividend_position_conf` DROP `id`;
ALTER TABLE `month_dividend_position_conf` DROP `agent_lv`;
ALTER TABLE `month_dividend_position_conf` DROP `position`;
ALTER TABLE `month_dividend_position_conf` DROP `min`;
ALTER TABLE `month_dividend_position_conf` DROP `max`;
ALTER TABLE `month_dividend_position_conf` DROP `activity_num`;
ALTER TABLE `month_dividend_position_conf` DROP `dividend_ratio`;
ALTER TABLE `month_dividend_position_conf` DROP `ctime`;
ALTER TABLE `month_dividend_position_conf` DROP `utime`;

----------------------------------------------------
--  `month_dividend_white_list`
----------------------------------------------------
ALTER TABLE `month_dividend_white_list` DROP `id`;
ALTER TABLE `month_dividend_white_list` DROP `name`;
ALTER TABLE `month_dividend_white_list` DROP `dividend_ratio`;
ALTER TABLE `month_dividend_white_list` DROP `ctime`;
ALTER TABLE `month_dividend_white_list` DROP `utime`;

----------------------------------------------------
--  `operation_log`
----------------------------------------------------
ALTER TABLE `operation_log` DROP `id`;
ALTER TABLE `operation_log` DROP `admin_id`;
ALTER TABLE `operation_log` DROP `method`;
ALTER TABLE `operation_log` DROP `route`;
ALTER TABLE `operation_log` DROP `action`;
ALTER TABLE `operation_log` DROP `input`;
ALTER TABLE `operation_log` DROP `user_agent`;
ALTER TABLE `operation_log` DROP `ips`;
ALTER TABLE `operation_log` DROP `response_code`;
ALTER TABLE `operation_log` DROP `ctime`;

----------------------------------------------------
--  `otc_stat`
----------------------------------------------------
ALTER TABLE `otc_stat` DROP `id`;
ALTER TABLE `otc_stat` DROP `date`;
ALTER TABLE `otc_stat` DROP `num_login`;
ALTER TABLE `otc_stat` DROP `num_user_new`;
ALTER TABLE `otc_stat` DROP `num_order`;
ALTER TABLE `otc_stat` DROP `num_order_deal`;
ALTER TABLE `otc_stat` DROP `num_order_buy`;
ALTER TABLE `otc_stat` DROP `num_order_sell`;
ALTER TABLE `otc_stat` DROP `num_funds`;
ALTER TABLE `otc_stat` DROP `num_amount`;
ALTER TABLE `otc_stat` DROP `num_amount_buy`;
ALTER TABLE `otc_stat` DROP `num_amount_sell`;
ALTER TABLE `otc_stat` DROP `num_fee_buy`;
ALTER TABLE `otc_stat` DROP `num_fee_sell`;
ALTER TABLE `otc_stat` DROP `game_recharge`;
ALTER TABLE `otc_stat` DROP `game_withdrawal`;
ALTER TABLE `otc_stat` DROP `usdt_recharge`;
ALTER TABLE `otc_stat` DROP `usdt_withdrawal`;
ALTER TABLE `otc_stat` DROP `usdt_fee`;

----------------------------------------------------
--  `otc_stat_all_people`
----------------------------------------------------
ALTER TABLE `otc_stat_all_people` DROP `uid`;
ALTER TABLE `otc_stat_all_people` DROP `buy_order`;
ALTER TABLE `otc_stat_all_people` DROP `sell_order`;
ALTER TABLE `otc_stat_all_people` DROP `buy_eusd`;
ALTER TABLE `otc_stat_all_people` DROP `sell_eusd`;
ALTER TABLE `otc_stat_all_people` DROP `usdt_recharge`;
ALTER TABLE `otc_stat_all_people` DROP `usdt_withdrawal`;

----------------------------------------------------
--  `permission`
----------------------------------------------------
ALTER TABLE `permission` DROP `id`;
ALTER TABLE `permission` DROP `slug`;
ALTER TABLE `permission` DROP `desc`;
ALTER TABLE `permission` DROP `ctime`;
ALTER TABLE `permission` DROP `utime`;
ALTER TABLE `permission` DROP `dtime`;

----------------------------------------------------
--  `profit_threshold`
----------------------------------------------------
ALTER TABLE `profit_threshold` DROP `id`;
ALTER TABLE `profit_threshold` DROP `threshold`;
ALTER TABLE `profit_threshold` DROP `admin_id`;
ALTER TABLE `profit_threshold` DROP `ctime`;
ALTER TABLE `profit_threshold` DROP `utime`;

----------------------------------------------------
--  `role`
----------------------------------------------------
ALTER TABLE `role` DROP `id`;
ALTER TABLE `role` DROP `name`;
ALTER TABLE `role` DROP `desc`;
ALTER TABLE `role` DROP `ctime`;
ALTER TABLE `role` DROP `utime`;

----------------------------------------------------
--  `role_admin`
----------------------------------------------------
ALTER TABLE `role_admin` DROP `id`;
ALTER TABLE `role_admin` DROP `roleid`;
ALTER TABLE `role_admin` DROP `adminid`;
ALTER TABLE `role_admin` DROP `granted_by`;
ALTER TABLE `role_admin` DROP `granted_at`;

----------------------------------------------------
--  `role_permission`
----------------------------------------------------
ALTER TABLE `role_permission` DROP `id`;
ALTER TABLE `role_permission` DROP `roleid`;
ALTER TABLE `role_permission` DROP `permissionid`;

----------------------------------------------------
--  `server_node`
----------------------------------------------------
ALTER TABLE `server_node` DROP `app_name`;
ALTER TABLE `server_node` DROP `region_id`;
ALTER TABLE `server_node` DROP `server_id`;
ALTER TABLE `server_node` DROP `last_ping`;

----------------------------------------------------
--  `smscodes`
----------------------------------------------------
ALTER TABLE `smscodes` DROP `id`;
ALTER TABLE `smscodes` DROP `national_code`;
ALTER TABLE `smscodes` DROP `mobile`;
ALTER TABLE `smscodes` DROP `action`;
ALTER TABLE `smscodes` DROP `code`;
ALTER TABLE `smscodes` DROP `status`;
ALTER TABLE `smscodes` DROP `ctime`;
ALTER TABLE `smscodes` DROP `etime`;

----------------------------------------------------
--  `smstemplates`
----------------------------------------------------
ALTER TABLE `smstemplates` DROP `id`;
ALTER TABLE `smstemplates` DROP `name`;
ALTER TABLE `smstemplates` DROP `type`;
ALTER TABLE `smstemplates` DROP `template`;
ALTER TABLE `smstemplates` DROP `ctime`;
ALTER TABLE `smstemplates` DROP `utime`;

----------------------------------------------------
--  `sys_msg`
----------------------------------------------------
ALTER TABLE `sys_msg` DROP `id`;
ALTER TABLE `sys_msg` DROP `key`;
ALTER TABLE `sys_msg` DROP `buyer`;
ALTER TABLE `sys_msg` DROP `seller`;
ALTER TABLE `sys_msg` DROP `admin`;
ALTER TABLE `sys_msg` DROP `ctime`;
ALTER TABLE `sys_msg` DROP `utime`;
ALTER TABLE `sys_msg` DROP `dtime`;

----------------------------------------------------
--  `sys_notification`
----------------------------------------------------
ALTER TABLE `sys_notification` DROP `id`;
ALTER TABLE `sys_notification` DROP `content`;
ALTER TABLE `sys_notification` DROP `admin_id`;
ALTER TABLE `sys_notification` DROP `status`;
ALTER TABLE `sys_notification` DROP `ctime`;
ALTER TABLE `sys_notification` DROP `utime`;

----------------------------------------------------
--  `task`
----------------------------------------------------
ALTER TABLE `task` DROP `id`;
ALTER TABLE `task` DROP `name`;
ALTER TABLE `task` DROP `alia`;
ALTER TABLE `task` DROP `app_name`;
ALTER TABLE `task` DROP `func_name`;
ALTER TABLE `task` DROP `spec`;
ALTER TABLE `task` DROP `status`;
ALTER TABLE `task` DROP `ctime`;
ALTER TABLE `task` DROP `utime`;
ALTER TABLE `task` DROP `desc`;

----------------------------------------------------
--  `task_result`
----------------------------------------------------
ALTER TABLE `task_result` DROP `id`;
ALTER TABLE `task_result` DROP `app_name`;
ALTER TABLE `task_result` DROP `region_id`;
ALTER TABLE `task_result` DROP `server_id`;
ALTER TABLE `task_result` DROP `name`;
ALTER TABLE `task_result` DROP `code`;
ALTER TABLE `task_result` DROP `detail`;
ALTER TABLE `task_result` DROP `end_time`;
ALTER TABLE `task_result` DROP `begin_time`;
ALTER TABLE `task_result` DROP `ctime`;

----------------------------------------------------
--  `top_agent`
----------------------------------------------------
ALTER TABLE `top_agent` DROP `id`;
ALTER TABLE `top_agent` DROP `national_code`;
ALTER TABLE `top_agent` DROP `mobile`;
ALTER TABLE `top_agent` DROP `status`;
ALTER TABLE `top_agent` DROP `ctime`;
ALTER TABLE `top_agent` DROP `utime`;

