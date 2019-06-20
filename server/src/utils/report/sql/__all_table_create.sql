-- --------------------------------------------------
--  Table Structure for `models.GameUserMonthReport`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `game_user_month_report` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`ctime` bigint NOT NULL COMMENT 'create time',
`profit` bigint NOT NULL COMMENT '自己的盈亏金额',
`agents_profit` bigint NOT NULL COMMENT '无限下级下级代理的盈亏金额',
`result_profit` bigint NOT NULL COMMENT '最终的盈亏金额',
`bet_amount` bigint NOT NULL COMMENT '投注额',
`effective_bet_amount` bigint NOT NULL COMMENT '累计有效投注额(本人加无限下级)',
`play_game_day` int NOT NULL COMMENT '玩游戏的天数',
`is_activity_user` bool NOT NULL DEFAULT false COMMENT '是否是活跃用户',
`agent_level` int unsigned NOT NULL COMMENT '代理等级',
`up_agent_uid` bigint unsigned NOT NULL COMMENT '上级代理的UID',
`activity_agent_num` int NOT NULL COMMENT '无限下级代理活跃人数',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='game_user_month_report table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.MonthDividendRecord`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `month_dividend_record` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`ctime` bigint NOT NULL COMMENT 'create time',
`self_dividend` bigint NOT NULL COMMENT '自己的分红',
`agent_dividend` bigint NOT NULL COMMENT '要分给代理的分红',
`result_dividend` bigint NOT NULL COMMENT '最终获得的分红 self_dividend-agent_dividend',
`receive_status` int NOT NULL DEFAULT false COMMENT '上级已发放状态是1,2是等待上级发放',
`received_time` bigint NOT NULL COMMENT '领取奖励的时间',
`pay_status` int NOT NULL DEFAULT false COMMENT '1是已支付状态,2是等待支付状态',
`level` int unsigned NOT NULL COMMENT 'level',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='month_dividend_record table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.ProfitReportDaily`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `profit_report_daily` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`channel_id` int unsigned NOT NULL COMMENT 'game channel id',
`bet` bigint NOT NULL COMMENT '本人的有效投注额',
`total_valid_bet` bigint NOT NULL COMMENT '累计有效投注额(本人加无限下级)',
`profit` bigint NOT NULL COMMENT '盈亏金额',
`salary` bigint NOT NULL COMMENT '日工资',
`self_dividend` bigint NOT NULL COMMENT '属于自己的月分红',
`agent_dividend` bigint NOT NULL COMMENT '分给下级代理的月分红',
`result_dividend` bigint NOT NULL COMMENT '最终获得的月分红,可能为负数',
`game_withdraw_amount` int unsigned NOT NULL COMMENT '游戏提现金额',
`game_recharge_amount` int unsigned NOT NULL COMMENT '游戏充值金额',
`ctime` bigint NOT NULL COMMENT 'create time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='profit_report_daily table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.ReportAgentDaily`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `report_agent_daily` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`sum_withdraw` bigint NOT NULL COMMENT 'sum_withdraw',
`ctime` bigint NOT NULL COMMENT 'create time',
PRIMARY KEY(`uid`,`ctime`)
) ENGINE=InnoDB COMMENT='report_agent_daily table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.ReportCommission`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `report_commission` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`level` int unsigned NOT NULL COMMENT 'level',
`team_withdraw` bigint NOT NULL COMMENT 'team_withdraw',
`team_can_withdraw` bigint NOT NULL COMMENT 'team_can_withdraw',
`ctime` bigint NOT NULL COMMENT 'create time',
PRIMARY KEY(`uid`)
) ENGINE=InnoDB COMMENT='report_commission table' DEFAULT CHARSET=utf8;
CREATE INDEX `report_commission_level` ON `report_commission` (`level`);

-- --------------------------------------------------
--  Table Structure for `models.ReportEusdDaily`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `report_eusd_daily` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`buy` bigint NOT NULL COMMENT '购买eusd数量',
`sell` bigint NOT NULL COMMENT '出售eusd数量',
`ctime` bigint NOT NULL COMMENT 'create time',
PRIMARY KEY(`uid`,`ctime`)
) ENGINE=InnoDB COMMENT='report_eusd_daily table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.ReportGameRecordAg`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `report_game_record_ag` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`account` varchar(50) NOT NULL COMMENT 'account',
`game_type` varchar(50) NOT NULL COMMENT '游戏ID',
`game_name` varchar(50) NOT NULL COMMENT '游戏名称',
`order_id` varchar(50) NOT NULL COMMENT '订单编号',
`table_id` varchar(50) NOT NULL COMMENT '桌号',
`bet` bigint NOT NULL COMMENT '投注额',
`valid_bet` bigint NOT NULL COMMENT '有效投注额',
`profit` bigint NOT NULL COMMENT '盈亏金额',
`bet_time` varchar(32) NOT NULL COMMENT '下注时间',
`ctime` bigint NOT NULL COMMENT 'ctime',
`uid` bigint unsigned NOT NULL COMMENT 'user id',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='report_game_record_Ag table' DEFAULT CHARSET=utf8;
CREATE INDEX `report_game_record_ag_uid_bet_valid_bet_profit_ctime` ON `report_game_record_ag` (`uid`, `bet`, `valid_bet`, `profit`, `ctime`);

-- --------------------------------------------------
--  Table Structure for `models.ReportGameRecordKy`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `report_game_record_ky` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`account` varchar(50) NOT NULL COMMENT 'account',
`game_id` varchar(50) NOT NULL COMMENT '游戏局号',
`game_name` varchar(50) NOT NULL COMMENT 'game_name',
`server_id` int NOT NULL COMMENT '房间ID',
`kind_id` varchar(50) NOT NULL COMMENT '游戏ID',
`table_id` int NOT NULL COMMENT '桌子号',
`chair_id` int NOT NULL COMMENT '椅子号',
`bet` bigint NOT NULL COMMENT '投注额',
`valid_bet` bigint NOT NULL COMMENT '有效投注额',
`profit` bigint NOT NULL COMMENT '盈亏金额',
`revenue` bigint NOT NULL COMMENT '抽水金额',
`start_time` varchar(32) NOT NULL COMMENT 'start_time',
`end_time` varchar(32) NOT NULL COMMENT 'end_time',
`ctime` bigint NOT NULL COMMENT 'create time',
`uid` bigint unsigned NOT NULL COMMENT 'user id',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='report_game_record_ky table' DEFAULT CHARSET=utf8;
CREATE INDEX `report_game_record_ky_uid_bet_valid_bet_profit_ctime` ON `report_game_record_ky` (`uid`, `bet`, `valid_bet`, `profit`, `ctime`);

-- --------------------------------------------------
--  Table Structure for `models.ReportGameRecordRg`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `report_game_record_rg` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`account` varchar(50) NOT NULL COMMENT 'account',
`game_name_id` varchar(50) NOT NULL COMMENT '游戏id',
`game_name` varchar(50) NOT NULL COMMENT '游戏名称',
`game_kind_name` varchar(50) NOT NULL COMMENT '玩法名称',
`order_id` varchar(50) NOT NULL COMMENT '订单编号',
`open_date` varchar(32) NOT NULL COMMENT '开奖时间',
`period_name` varchar(50) NOT NULL COMMENT '期号',
`open_number` varchar(50) NOT NULL COMMENT '开奖号码',
`status` tinyint unsigned NOT NULL COMMENT '订单状态',
`bet` bigint NOT NULL COMMENT '投注额',
`valid_bet` bigint NOT NULL COMMENT '有效投注额',
`profit` bigint NOT NULL COMMENT '盈亏金额',
`bet_time` varchar(32) NOT NULL COMMENT '下注时间',
`bet_content` varchar(50) NOT NULL COMMENT '下注内容',
`ctime` bigint NOT NULL COMMENT 'ctime',
`uid` bigint unsigned NOT NULL COMMENT 'user id',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='report_game_record_Rg table' DEFAULT CHARSET=utf8;
CREATE INDEX `report_game_record_rg_uid_bet_valid_bet_profit_ctime` ON `report_game_record_rg` (`uid`, `bet`, `valid_bet`, `profit`, `ctime`);

-- --------------------------------------------------
--  Table Structure for `models.ReportGameTransferDaily`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `report_game_transfer_daily` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`channel_id` int unsigned NOT NULL COMMENT '渠道ID',
`recharge` bigint NOT NULL COMMENT '游戏充值',
`withdraw` bigint NOT NULL COMMENT '游戏提现',
`ctime` bigint NOT NULL COMMENT 'create time',
PRIMARY KEY(`uid`,`channel_id`,`ctime`)
) ENGINE=InnoDB COMMENT='report_game_transfer_daily table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.ReportGameUserDaily`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `report_game_user_daily` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`p_uid` bigint unsigned NOT NULL COMMENT 'parent uid',
`level` int unsigned NOT NULL COMMENT 'level',
`channel_id` int unsigned NOT NULL COMMENT 'game channel id',
`bet` bigint NOT NULL COMMENT '投注金额',
`valid_bet` bigint NOT NULL COMMENT '有效投注额',
`total_valid_bet` bigint NOT NULL COMMENT '累计有效投注额(本人加无限下级)',
`total_bet_num` int NOT NULL COMMENT '累计投注人数(本人加无限下级)',
`profit` bigint NOT NULL COMMENT '盈亏金额',
`total_profit` bigint NOT NULL COMMENT '累计盈亏金额(本人加无限下级)',
`salary` bigint NOT NULL COMMENT '日工资',
`team_salary` bigint NOT NULL COMMENT '累计日工资(本人加无限下级)',
`status` tinyint unsigned NOT NULL COMMENT 'status',
`ctime` bigint NOT NULL COMMENT 'create time',
PRIMARY KEY(`uid`,`channel_id`,`ctime`)
) ENGINE=InnoDB COMMENT='report_game_user_daily table' DEFAULT CHARSET=utf8;
CREATE INDEX `report_game_user_daily_level` ON `report_game_user_daily` (`level`);

-- --------------------------------------------------
--  Table Structure for `models.ReportStatisticGameAll`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `report_statistic_game_all` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`channel_id` int unsigned NOT NULL COMMENT 'platform type id',
`newer_nums` bigint NOT NULL COMMENT 'new player nums',
`bet` bigint NOT NULL COMMENT 'bet',
`valid_bet` bigint NOT NULL COMMENT 'valid bet',
`profit` bigint NOT NULL COMMENT 'profit',
`revenue` bigint NOT NULL COMMENT '抽水',
`ctime` bigint NOT NULL COMMENT 'curtime',
`note` varchar(255) NOT NULL COMMENT 'note',
`game_id` int unsigned NOT NULL COMMENT 'game id',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='report_statistic_game_all table' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.ReportStatisticSum`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `report_statistic_sum` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
`channel_id` int unsigned NOT NULL COMMENT 'platform type id',
`channel_positive_nums` bigint unsigned NOT NULL COMMENT 'today channel positive nums',
`channel_salary_daily` bigint NOT NULL COMMENT 'channel daily salary',
`channel_rg_dividend` bigint NOT NULL COMMENT 'dividend every month',
`channel_withdraw_eusd` bigint NOT NULL COMMENT 'channel withdraw eusd',
`channel_recharge_eusd` bigint NOT NULL COMMENT 'channel recharge eusd',
`ctime` bigint NOT NULL COMMENT 'curtime',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='report_statistic_sum table every day' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.ReportTeamDaily`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `report_team_daily` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`eusd_buy` bigint NOT NULL COMMENT '团队eusd购买金额',
`eusd_sell` bigint NOT NULL COMMENT '团队eusd出售金额',
`level` int unsigned NOT NULL COMMENT 'level',
`ctime` bigint NOT NULL COMMENT 'create time',
PRIMARY KEY(`uid`,`ctime`)
) ENGINE=InnoDB COMMENT='report_team_daily table' DEFAULT CHARSET=utf8;
CREATE INDEX `report_team_daily_level` ON `report_team_daily` (`level`);

-- --------------------------------------------------
--  Table Structure for `models.ReportTeamGameTransferDaily`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `report_team_game_transfer_daily` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`channel_id` int unsigned NOT NULL COMMENT '渠道ID',
`team_recharge` bigint NOT NULL COMMENT '团队游戏充值',
`team_withdraw` bigint NOT NULL COMMENT '团队游戏提现',
`level` int unsigned NOT NULL COMMENT 'level',
`ctime` bigint NOT NULL COMMENT 'create time',
PRIMARY KEY(`uid`,`channel_id`,`ctime`)
) ENGINE=InnoDB COMMENT='report_team_game_transfer_daily table' DEFAULT CHARSET=utf8;
CREATE INDEX `report_team_game_transfer_daily_level` ON `report_team_game_transfer_daily` (`level`);

