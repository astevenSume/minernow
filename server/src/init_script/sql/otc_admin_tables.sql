
-- --------------------------------------------------
--  Table Structure for `models.ActivityUserConf`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `activity_user_conf` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`play_game_day` int NOT NULL COMMENT '玩了并且有投注的天数 配置',
`bet_amount` bigint NOT NULL COMMENT '投注额配置',
`effective_bet_amount` bigint NOT NULL COMMENT '累计有效投注额(本人加无限下级)的配置',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='月分红活跃用户定义的配置' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.AdminUser`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `admin_user` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'user id',
`name` varchar(100) NOT NULL COMMENT 'user name',
`email` varchar(100) NOT NULL COMMENT 'user email',
`status` tinyint NOT NULL COMMENT 'status',
`ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time',
`utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time',
`dtime` bigint NOT NULL DEFAULT 0 COMMENT 'forbit time',
`login_time` bigint NOT NULL DEFAULT 0 COMMENT 'last login time',
`pwd` varchar(100) NOT NULL DEFAULT '' COMMENT 'password',
`whitelist_ips` varchar(256) NOT NULL DEFAULT '' COMMENT 'whitelist_ips',
`is_bind` bool NOT NULL DEFAULT false COMMENT 'is_bind',
`secret_id` varchar(128) NOT NULL DEFAULT '' COMMENT 'secret_id',
`qr_code` text NOT NULL COMMENT 'qr_code',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='admin user table' DEFAULT CHARSET=utf8;
CREATE UNIQUE INDEX `admin_user_email` ON `admin_user` (`email`);

-- --------------------------------------------------
--  Table Structure for `models.AgentWhiteList`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `agent_white_list` (
`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`name` varchar(64) NOT NULL COMMENT 'name',
`commission` int NOT NULL COMMENT 'commission',
`precision` int NOT NULL COMMENT 'precision',
`ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time',
`utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='agent white list table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.Announcement`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `announcement` (
`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`type` tinyint NOT NULL COMMENT 'type',
`title` varchar(100) NOT NULL COMMENT 'title',
`content` varchar(500) NOT NULL COMMENT 'content',
`stime` bigint NOT NULL DEFAULT 0 COMMENT 'start time',
`etime` bigint NOT NULL DEFAULT 0 COMMENT 'end time',
`ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time',
`utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='system announcement' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.AppChannel`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `app_channel` (
`id` int unsigned NOT NULL COMMENT 'id',
`is_third_hall` tinyint NOT NULL COMMENT 'is_third_hall',
`name` varchar(256) NOT NULL COMMENT 'name',
`desc` varchar(256) NOT NULL COMMENT 'desc',
`exchangeRate` int NOT NULL DEFAULT 0 COMMENT 'exchange rate',
`precision` int NOT NULL DEFAULT 0 COMMENT 'precision',
`profit_rate` int NOT NULL DEFAULT 0 COMMENT 'profit_rate',
`ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time',
`utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='app channel table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.AppType`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `app_type` (
`id` int unsigned NOT NULL COMMENT 'id',
`name` varchar(256) NOT NULL COMMENT 'name',
`desc` varchar(256) NOT NULL COMMENT 'desc',
`ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time',
`utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='app type table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.AppVersion`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `app_version` (
`id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
`version` varchar(100) NOT NULL COMMENT 'app version show',
`version_num` int NOT NULL COMMENT 'app version num',
`changelog` varchar(300) NOT NULL COMMENT 'change log',
`download` varchar(300) NOT NULL COMMENT 'download url',
`system` tinyint NOT NULL COMMENT 'system type',
`status` tinyint NOT NULL COMMENT 'app version status',
`ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time',
`utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time',
`dtime` bigint NOT NULL DEFAULT 0 COMMENT 'update dtime',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='app version table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.AppWhitelist`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `app_whitelist` (
`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`channel_id` int unsigned NOT NULL COMMENT 'channel_id',
`app_id` varchar(16) NOT NULL COMMENT 'app_id',
`ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='app_whitelist table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.AppealService`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `appeal_service` (
`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`admin_id` int unsigned NOT NULL COMMENT 'admin_id',
`wechat` varchar(32) NOT NULL COMMENT 'wechat',
`qr_code` varchar(300) NOT NULL COMMENT '',
`status` tinyint NOT NULL COMMENT 'status',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='appeal_service table' DEFAULT CHARSET=utf8;
CREATE UNIQUE INDEX `appeal_service_admin_id` ON `appeal_service` (`admin_id`);

-- --------------------------------------------------
--  Table Structure for `models.Apps`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `apps` (
`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`position` int unsigned NOT NULL COMMENT 'position',
`name` varchar(64) NOT NULL COMMENT 'name',
`desc` varchar(256) NOT NULL COMMENT 'desc',
`url` varchar(128) NOT NULL COMMENT 'url',
`icon_url` varchar(128) NOT NULL COMMENT 'icon_url',
`type_id` tinyint NOT NULL COMMENT 'type_id',
`channel_id` int unsigned NOT NULL COMMENT 'channel_id',
`app_id` varchar(50) NOT NULL COMMENT 'app_id',
`featured` tinyint NOT NULL COMMENT 'featured',
`status` tinyint NOT NULL COMMENT 'status',
`orientation` tinyint NOT NULL DEFAULT 1 COMMENT 'orientation',
`ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time',
`utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='' DEFAULT CHARSET=utf8;
CREATE INDEX `apps_type_id` ON `apps` (`type_id`);

-- --------------------------------------------------
--  Table Structure for `models.Banner`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `banner` (
`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`subject` varchar(256) NOT NULL COMMENT 'subject',
`image` varchar(256) NOT NULL COMMENT 'image',
`url` varchar(256) NOT NULL DEFAULT '' COMMENT 'image',
`status` tinyint NOT NULL COMMENT 'status',
`ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time',
`utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time',
`stime` bigint NOT NULL DEFAULT 0 COMMENT 'start time',
`etime` bigint NOT NULL DEFAULT 0 COMMENT 'end time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='banner table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.Commissionrates`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `commissionrates` (
`id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
`min` bigint unsigned NOT NULL COMMENT 'min',
`max` bigint unsigned NOT NULL COMMENT 'max',
`commission` int NOT NULL DEFAULT 0 COMMENT 'commission',
`precision` int NOT NULL DEFAULT 0 COMMENT 'precision',
`ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time',
`utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='commissionrates table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.Config`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `config` (
`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`action` tinyint NOT NULL COMMENT 'action',
`key` varchar(256) NOT NULL COMMENT 'key',
`value` text NOT NULL COMMENT 'value',
`desc` varchar(256) NOT NULL COMMENT 'descripe',
`ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='data config, key, value, desc' DEFAULT CHARSET=utf8;
CREATE INDEX `config_action` ON `config` (`action`);
CREATE UNIQUE INDEX `config_key` ON `config` (`key`);

-- --------------------------------------------------
--  Table Structure for `models.ConfigWarning`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `config_warning` (
`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`type` tinyint NOT NULL COMMENT 'type',
`national_code` varchar(16) NOT NULL COMMENT 'national_code',
`mobile` varchar(32) NOT NULL COMMENT 'mobile',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='config_warning table' DEFAULT CHARSET=utf8;
CREATE UNIQUE INDEX `config_warning_type_national_code_mobile` ON `config_warning` (`type`, `national_code`, `mobile`);

-- --------------------------------------------------
--  Table Structure for `models.Endpoint`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `endpoint` (
`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`endpoint` varchar(100) NOT NULL COMMENT 'endpoint',
`ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time',
`utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='endpoint table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.IpWhiteList`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `ip_white_list` (
`id` int unsigned NOT NULL COMMENT 'id',
`ip` varchar(0) NOT NULL COMMENT 'ip address',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='ip white list for otc' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.MonthDividendPositionConf`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `month_dividend_position_conf` (
`id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
`agent_lv` int NOT NULL COMMENT 'agent_lv',
`position` int NOT NULL COMMENT 'position',
`min` bigint NOT NULL COMMENT 'min',
`max` bigint NOT NULL COMMENT 'max',
`activity_num` int NOT NULL COMMENT 'activity num',
`dividend_ratio` int NOT NULL COMMENT 'dividend_ratio num',
`ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time',
`utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='month_dividend_position_conf table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.OperationLog`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `operation_log` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`admin_id` bigint unsigned NOT NULL COMMENT 'admin_id',
`method` varchar(100) NOT NULL DEFAULT '' COMMENT 'req method',
`route` varchar(100) NOT NULL DEFAULT '' COMMENT 'route ',
`action` int NOT NULL DEFAULT 0 COMMENT 'action',
`input` varchar(65535) NOT NULL DEFAULT '' COMMENT 'input',
`user_agent` varchar(512) NOT NULL DEFAULT '' COMMENT 'user agent',
`ips` varchar(100) NOT NULL DEFAULT '' COMMENT 'ip',
`response_code` int NOT NULL DEFAULT 0 COMMENT 'response code',
`ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='OperationLog table' DEFAULT CHARSET=utf8;
CREATE INDEX `operation_log_action_admin_id` ON `operation_log` (`action`, `admin_id`);

-- --------------------------------------------------
--  Table Structure for `models.OtcStat`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `otc_stat` (
`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`date` int unsigned NOT NULL COMMENT 'date',
`num_login` int unsigned NOT NULL DEFAULT 0 COMMENT '',
`num_user_new` int unsigned NOT NULL DEFAULT 0 COMMENT '',
`num_order` int unsigned NOT NULL DEFAULT 0 COMMENT 'date',
`num_order_deal` int unsigned NOT NULL DEFAULT 0 COMMENT '',
`num_order_buy` int unsigned NOT NULL DEFAULT 0 COMMENT '',
`num_order_sell` int unsigned NOT NULL DEFAULT 0 COMMENT '',
`num_funds` int unsigned NOT NULL DEFAULT 0 COMMENT '',
`num_amount` int unsigned NOT NULL DEFAULT 0 COMMENT '',
`num_amount_buy` int unsigned NOT NULL DEFAULT 0 COMMENT '',
`num_amount_sell` int unsigned NOT NULL DEFAULT 0 COMMENT '',
`num_fee_buy` int unsigned NOT NULL DEFAULT 0 COMMENT '',
`num_fee_sell` int unsigned NOT NULL DEFAULT 0 COMMENT '',
`game_recharge` int unsigned NOT NULL DEFAULT 0 COMMENT '',
`game_withdrawal` int unsigned NOT NULL DEFAULT 0 COMMENT '',
`usdt_recharge` int unsigned NOT NULL DEFAULT 0 COMMENT '',
`usdt_withdrawal` int unsigned NOT NULL DEFAULT 0 COMMENT '',
`usdt_fee` int unsigned NOT NULL DEFAULT 0 COMMENT '',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='OtcStat' DEFAULT CHARSET=utf8;
CREATE UNIQUE INDEX `otc_stat_date` ON `otc_stat` (`date`);

-- --------------------------------------------------
--  Table Structure for `models.OtcStatAllPeople`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `otc_stat_all_people` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`buy_order` int unsigned NOT NULL DEFAULT 0 COMMENT '',
`sell_order` int unsigned NOT NULL DEFAULT 0 COMMENT '',
`buy_eusd` int unsigned NOT NULL DEFAULT 0 COMMENT '',
`sell_eusd` int unsigned NOT NULL DEFAULT 0 COMMENT '',
`usdt_recharge` int unsigned NOT NULL DEFAULT 0 COMMENT '',
`usdt_withdrawal` int unsigned NOT NULL DEFAULT 0 COMMENT '',
PRIMARY KEY(`uid`)
) ENGINE=InnoDB COMMENT='OtcStatAllPeople' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.Permission`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `permission` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`slug` varchar(100) NOT NULL COMMENT 'permission_name',
`desc` varchar(100) NOT NULL COMMENT 'permission_desc',
`ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time',
`utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time',
`dtime` bigint NOT NULL DEFAULT 0 COMMENT 'D time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='permission table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.Role`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `role` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'role id',
`name` varchar(100) NOT NULL COMMENT 'role name',
`desc` varchar(100) NOT NULL COMMENT 'role description',
`ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time',
`utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='role table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.RoleAdmin`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `role_admin` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`roleid` bigint unsigned NOT NULL COMMENT 'role_id',
`adminid` bigint unsigned NOT NULL COMMENT 'admin_id',
`granted_by` varchar(100) NOT NULL COMMENT 'role granted_by',
`granted_at` bigint NOT NULL DEFAULT 0 COMMENT 'granted_at time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='role_admin table' DEFAULT CHARSET=utf8;
CREATE INDEX `role_admin_roleid_adminid` ON `role_admin` (`roleid`, `adminid`);

-- --------------------------------------------------
--  Table Structure for `models.RolePermission`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `role_permission` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`roleid` bigint unsigned NOT NULL COMMENT 'roleid',
`permissionid` bigint unsigned NOT NULL COMMENT 'pemissionid',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='role_permission table' DEFAULT CHARSET=utf8;
CREATE INDEX `role_permission_roleid_permissionid` ON `role_permission` (`roleid`, `permissionid`);

-- --------------------------------------------------
--  Table Structure for `models.ServerNode`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `server_node` (
`app_name` varchar(256) NOT NULL COMMENT 'application name',
`region_id` bigint NOT NULL COMMENT 'region id',
`server_id` bigint NOT NULL COMMENT 'server id',
`last_ping` int unsigned NOT NULL COMMENT 'last ping timestamp',
PRIMARY KEY(`app_name`,`region_id`,`server_id`)
) ENGINE=InnoDB COMMENT='server node table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.Smscodes`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `smscodes` (
`id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
`national_code` varchar(32) NOT NULL DEFAULT '' COMMENT 'nationalCode',
`mobile` varchar(32) NOT NULL DEFAULT '' COMMENT 'mobile',
`action` varchar(100) NOT NULL DEFAULT '' COMMENT 'action',
`code` varchar(16) NOT NULL DEFAULT '' COMMENT 'code',
`status` tinyint NOT NULL COMMENT 'status',
`ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time',
`etime` bigint NOT NULL DEFAULT 0 COMMENT 'expired time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='smscodes table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.Smstemplates`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `smstemplates` (
`id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
`name` varchar(100) NOT NULL DEFAULT '' COMMENT 'name',
`action` varchar(100) NOT NULL DEFAULT '' COMMENT 'action',
`template` varchar(256) NOT NULL DEFAULT '' COMMENT 'template',
`ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time',
`utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='smstemplates table' DEFAULT CHARSET=utf8;
CREATE UNIQUE INDEX `smstemplates_action` ON `smstemplates` (`action`);

-- --------------------------------------------------
--  Table Structure for `models.SystemMessage`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `sys_msg` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`key` varchar(200) NOT NULL COMMENT 'System Message Key',
`buyer` varchar(400) NOT NULL COMMENT 'Buyer Show',
`seller` varchar(400) NOT NULL COMMENT 'Seller Show',
`admin` varchar(400) NOT NULL COMMENT 'Admin Show',
`ctime` bigint NULL COMMENT '',
`utime` bigint NULL COMMENT '',
`dtime` bigint NULL COMMENT '',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='System Message' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.SysNotification`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `sys_notification` (
`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`content` varchar(128) NOT NULL COMMENT 'content',
`admin_id` int unsigned NOT NULL COMMENT 'admin_id',
`status` tinyint NOT NULL COMMENT 'status',
`ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time',
`utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='sys_notification table' DEFAULT CHARSET=utf8;
CREATE INDEX `sys_notification_admin_id` ON `sys_notification` (`admin_id`);

-- --------------------------------------------------
--  Table Structure for `models.Task`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `task` (
`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`name` varchar(256) NOT NULL COMMENT 'task name',
`alia` varchar(256) NOT NULL COMMENT 'task alia',
`app_name` varchar(256) NOT NULL COMMENT 'application name',
`func_name` varchar(256) NOT NULL COMMENT 'task function name',
`spec` varchar(256) NOT NULL COMMENT 'task spec string',
`status` tinyint unsigned NOT NULL COMMENT 'status',
`ctime` int unsigned NOT NULL COMMENT 'create time',
`utime` int unsigned NOT NULL COMMENT 'update time',
`desc` varchar(256) NOT NULL COMMENT 'task detail string',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='task table' DEFAULT CHARSET=utf8;
CREATE INDEX `task_name` ON `task` (`name`);
CREATE UNIQUE INDEX `task_app_name_name` ON `task` (`app_name`, `name`);

-- --------------------------------------------------
--  Table Structure for `models.TaskResult`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `task_result` (
`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`app_name` varchar(256) NOT NULL COMMENT 'app name',
`region_id` bigint NOT NULL COMMENT 'region id',
`server_id` bigint NOT NULL COMMENT 'server id',
`name` varchar(256) NOT NULL COMMENT 'name',
`code` int NOT NULL COMMENT 'result code',
`detail` varchar(256) NOT NULL COMMENT 'task result detail',
`begin_time` int unsigned NOT NULL COMMENT 'begin time',
`end_time` int unsigned NOT NULL COMMENT 'end time',
`ctime` int unsigned NOT NULL COMMENT 'create time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='task result table' DEFAULT CHARSET=utf8;
CREATE INDEX `task_result_name_region_id_server_id` ON `task_result` (`name`, `region_id`, `server_id`);
CREATE INDEX `task_result_app_name` ON `task_result` (`app_name`);

-- --------------------------------------------------
--  Table Structure for `models.TopAgent`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `top_agent` (
`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`national_code` varchar(16) NOT NULL COMMENT 'national_code',
`mobile` varchar(32) NOT NULL COMMENT 'mobile',
`status` tinyint NOT NULL COMMENT 'status',
`ctime` bigint NOT NULL DEFAULT 0 COMMENT 'create time',
`utime` bigint NOT NULL DEFAULT 0 COMMENT 'update time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='top_agent table' DEFAULT CHARSET=utf8;
CREATE UNIQUE INDEX `top_agent_mobile` ON `top_agent` (`mobile`);

-- --------------------------------------------------
--  Table Structure for `models.Token`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `token` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`client_type` int unsigned NOT NULL COMMENT 'app type',
`mtime` bigint NOT NULL COMMENT 'access_token modify time',
`access_token` varchar(256) NOT NULL COMMENT 'access_token',
`mac` varchar(100) NOT NULL COMMENT '',
PRIMARY KEY(`uid`,`client_type`,`mac`)
) ENGINE=InnoDB COMMENT='access token table' DEFAULT CHARSET=utf8;

