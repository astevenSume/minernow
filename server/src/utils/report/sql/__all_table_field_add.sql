----------------------------------------------------
--  `game_user_month_report`
----------------------------------------------------
ALTER TABLE `game_user_month_report` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `game_user_month_report` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id' AFTER `id`;
ALTER TABLE `game_user_month_report` ADD `ctime` bigint NOT NULL COMMENT 'create time' AFTER `uid`;
ALTER TABLE `game_user_month_report` ADD `profit` bigint NOT NULL COMMENT '自己的盈亏金额' AFTER `ctime`;
ALTER TABLE `game_user_month_report` ADD `agents_profit` bigint NOT NULL COMMENT '无限下级下级代理的盈亏金额' AFTER `profit`;
ALTER TABLE `game_user_month_report` ADD `result_profit` bigint NOT NULL COMMENT '最终的盈亏金额' AFTER `agents_profit`;
ALTER TABLE `game_user_month_report` ADD `bet_amount` bigint NOT NULL COMMENT '投注额' AFTER `result_profit`;
ALTER TABLE `game_user_month_report` ADD `effective_bet_amount` bigint NOT NULL COMMENT '累计有效投注额(本人加无限下级)' AFTER `bet_amount`;
ALTER TABLE `game_user_month_report` ADD `play_game_day` int NOT NULL COMMENT '玩游戏的天数' AFTER `effective_bet_amount`;
ALTER TABLE `game_user_month_report` ADD `is_activity_user` bool NOT NULL DEFAULT false COMMENT '是否是活跃用户' AFTER `play_game_day`;
ALTER TABLE `game_user_month_report` ADD `agent_level` int unsigned NOT NULL COMMENT '代理等级' AFTER `is_activity_user`;
ALTER TABLE `game_user_month_report` ADD `up_agent_uid` bigint unsigned NOT NULL COMMENT '上级代理的UID' AFTER `agent_level`;
ALTER TABLE `game_user_month_report` ADD `activity_agent_num` int NOT NULL COMMENT '无限下级代理活跃人数' AFTER `up_agent_uid`;

----------------------------------------------------
--  `month_dividend_record`
----------------------------------------------------
ALTER TABLE `month_dividend_record` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `month_dividend_record` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id' AFTER `id`;
ALTER TABLE `month_dividend_record` ADD `ctime` bigint NOT NULL COMMENT 'create time' AFTER `uid`;
ALTER TABLE `month_dividend_record` ADD `self_dividend` bigint NOT NULL COMMENT '自己的分红' AFTER `ctime`;
ALTER TABLE `month_dividend_record` ADD `agent_dividend` bigint NOT NULL COMMENT '要分给代理的分红' AFTER `self_dividend`;
ALTER TABLE `month_dividend_record` ADD `result_dividend` bigint NOT NULL COMMENT '最终获得的分红 self_dividend-agent_dividend' AFTER `agent_dividend`;
ALTER TABLE `month_dividend_record` ADD `receive_status` int NOT NULL DEFAULT false COMMENT '上级已发放状态是1,2是等待上级发放' AFTER `result_dividend`;
ALTER TABLE `month_dividend_record` ADD `received_time` bigint NOT NULL COMMENT '领取奖励的时间' AFTER `receive_status`;
ALTER TABLE `month_dividend_record` ADD `pay_status` int NOT NULL DEFAULT false COMMENT '1是已支付状态,2是等待支付状态' AFTER `received_time`;
ALTER TABLE `month_dividend_record` ADD `level` int unsigned NOT NULL COMMENT 'level' AFTER `pay_status`;

----------------------------------------------------
--  `profit_report_daily`
----------------------------------------------------
ALTER TABLE `profit_report_daily` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `profit_report_daily` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id' AFTER `id`;
ALTER TABLE `profit_report_daily` ADD `channel_id` int unsigned NOT NULL COMMENT 'game channel id' AFTER `uid`;
ALTER TABLE `profit_report_daily` ADD `bet` bigint NOT NULL COMMENT '本人的有效投注额' AFTER `channel_id`;
ALTER TABLE `profit_report_daily` ADD `total_valid_bet` bigint NOT NULL COMMENT '累计有效投注额(本人加无限下级)' AFTER `bet`;
ALTER TABLE `profit_report_daily` ADD `profit` bigint NOT NULL COMMENT '盈亏金额' AFTER `total_valid_bet`;
ALTER TABLE `profit_report_daily` ADD `salary` bigint NOT NULL COMMENT '日工资' AFTER `profit`;
ALTER TABLE `profit_report_daily` ADD `self_dividend` bigint NOT NULL COMMENT '属于自己的月分红' AFTER `salary`;
ALTER TABLE `profit_report_daily` ADD `agent_dividend` bigint NOT NULL COMMENT '分给下级代理的月分红' AFTER `self_dividend`;
ALTER TABLE `profit_report_daily` ADD `result_dividend` bigint NOT NULL COMMENT '最终获得的月分红,可能为负数' AFTER `agent_dividend`;
ALTER TABLE `profit_report_daily` ADD `game_withdraw_amount` int unsigned NOT NULL COMMENT '游戏提现金额' AFTER `result_dividend`;
ALTER TABLE `profit_report_daily` ADD `game_recharge_amount` int unsigned NOT NULL COMMENT '游戏充值金额' AFTER `game_withdraw_amount`;
ALTER TABLE `profit_report_daily` ADD `ctime` bigint NOT NULL COMMENT 'create time' AFTER `game_recharge_amount`;

----------------------------------------------------
--  `report_agent_daily`
----------------------------------------------------
ALTER TABLE `report_agent_daily` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `report_agent_daily` ADD `sum_withdraw` bigint NOT NULL COMMENT 'sum_withdraw' AFTER `uid`;
ALTER TABLE `report_agent_daily` ADD `ctime` bigint NOT NULL COMMENT 'create time' AFTER `sum_withdraw`;

----------------------------------------------------
--  `report_commission`
----------------------------------------------------
ALTER TABLE `report_commission` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `report_commission` ADD `level` int unsigned NOT NULL COMMENT 'level' AFTER `uid`;
ALTER TABLE `report_commission` ADD `team_withdraw` bigint NOT NULL COMMENT 'team_withdraw' AFTER `level`;
ALTER TABLE `report_commission` ADD `team_can_withdraw` bigint NOT NULL COMMENT 'team_can_withdraw' AFTER `team_withdraw`;
ALTER TABLE `report_commission` ADD `ctime` bigint NOT NULL COMMENT 'create time' AFTER `team_can_withdraw`;

----------------------------------------------------
--  `report_eusd_daily`
----------------------------------------------------
ALTER TABLE `report_eusd_daily` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `report_eusd_daily` ADD `buy` bigint NOT NULL COMMENT '购买eusd数量' AFTER `uid`;
ALTER TABLE `report_eusd_daily` ADD `sell` bigint NOT NULL COMMENT '出售eusd数量' AFTER `buy`;
ALTER TABLE `report_eusd_daily` ADD `ctime` bigint NOT NULL COMMENT 'create time' AFTER `sell`;

----------------------------------------------------
--  `report_game_record_ag`
----------------------------------------------------
ALTER TABLE `report_game_record_ag` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `report_game_record_ag` ADD `account` varchar(50) NOT NULL COMMENT 'account' AFTER `id`;
ALTER TABLE `report_game_record_ag` ADD `game_type` varchar(50) NOT NULL COMMENT '游戏ID' AFTER `account`;
ALTER TABLE `report_game_record_ag` ADD `game_name` varchar(50) NOT NULL COMMENT '游戏名称' AFTER `game_type`;
ALTER TABLE `report_game_record_ag` ADD `order_id` varchar(50) NOT NULL COMMENT '订单编号' AFTER `game_name`;
ALTER TABLE `report_game_record_ag` ADD `table_id` varchar(50) NOT NULL COMMENT '桌号' AFTER `order_id`;
ALTER TABLE `report_game_record_ag` ADD `bet` bigint NOT NULL COMMENT '投注额' AFTER `table_id`;
ALTER TABLE `report_game_record_ag` ADD `valid_bet` bigint NOT NULL COMMENT '有效投注额' AFTER `bet`;
ALTER TABLE `report_game_record_ag` ADD `profit` bigint NOT NULL COMMENT '盈亏金额' AFTER `valid_bet`;
ALTER TABLE `report_game_record_ag` ADD `bet_time` varchar(32) NOT NULL COMMENT '下注时间' AFTER `profit`;
ALTER TABLE `report_game_record_ag` ADD `ctime` bigint NOT NULL COMMENT 'ctime' AFTER `bet_time`;
ALTER TABLE `report_game_record_ag` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id' AFTER `ctime`;

----------------------------------------------------
--  `report_game_record_ky`
----------------------------------------------------
ALTER TABLE `report_game_record_ky` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `report_game_record_ky` ADD `account` varchar(50) NOT NULL COMMENT 'account' AFTER `id`;
ALTER TABLE `report_game_record_ky` ADD `game_id` varchar(50) NOT NULL COMMENT '游戏局号' AFTER `account`;
ALTER TABLE `report_game_record_ky` ADD `game_name` varchar(50) NOT NULL COMMENT 'game_name' AFTER `game_id`;
ALTER TABLE `report_game_record_ky` ADD `server_id` int NOT NULL COMMENT '房间ID' AFTER `game_name`;
ALTER TABLE `report_game_record_ky` ADD `kind_id` varchar(50) NOT NULL COMMENT '游戏ID' AFTER `server_id`;
ALTER TABLE `report_game_record_ky` ADD `table_id` int NOT NULL COMMENT '桌子号' AFTER `kind_id`;
ALTER TABLE `report_game_record_ky` ADD `chair_id` int NOT NULL COMMENT '椅子号' AFTER `table_id`;
ALTER TABLE `report_game_record_ky` ADD `bet` bigint NOT NULL COMMENT '投注额' AFTER `chair_id`;
ALTER TABLE `report_game_record_ky` ADD `valid_bet` bigint NOT NULL COMMENT '有效投注额' AFTER `bet`;
ALTER TABLE `report_game_record_ky` ADD `profit` bigint NOT NULL COMMENT '盈亏金额' AFTER `valid_bet`;
ALTER TABLE `report_game_record_ky` ADD `revenue` bigint NOT NULL COMMENT '抽水金额' AFTER `profit`;
ALTER TABLE `report_game_record_ky` ADD `start_time` varchar(32) NOT NULL COMMENT 'start_time' AFTER `revenue`;
ALTER TABLE `report_game_record_ky` ADD `end_time` varchar(32) NOT NULL COMMENT 'end_time' AFTER `start_time`;
ALTER TABLE `report_game_record_ky` ADD `ctime` bigint NOT NULL COMMENT 'create time' AFTER `end_time`;
ALTER TABLE `report_game_record_ky` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id' AFTER `ctime`;

----------------------------------------------------
--  `report_game_record_rg`
----------------------------------------------------
ALTER TABLE `report_game_record_rg` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `report_game_record_rg` ADD `account` varchar(50) NOT NULL COMMENT 'account' AFTER `id`;
ALTER TABLE `report_game_record_rg` ADD `game_name_id` varchar(50) NOT NULL COMMENT '游戏id' AFTER `account`;
ALTER TABLE `report_game_record_rg` ADD `game_name` varchar(50) NOT NULL COMMENT '游戏名称' AFTER `game_name_id`;
ALTER TABLE `report_game_record_rg` ADD `game_kind_name` varchar(50) NOT NULL COMMENT '玩法名称' AFTER `game_name`;
ALTER TABLE `report_game_record_rg` ADD `order_id` varchar(50) NOT NULL COMMENT '订单编号' AFTER `game_kind_name`;
ALTER TABLE `report_game_record_rg` ADD `open_date` varchar(32) NOT NULL COMMENT '开奖时间' AFTER `order_id`;
ALTER TABLE `report_game_record_rg` ADD `period_name` varchar(50) NOT NULL COMMENT '期号' AFTER `open_date`;
ALTER TABLE `report_game_record_rg` ADD `open_number` varchar(50) NOT NULL COMMENT '开奖号码' AFTER `period_name`;
ALTER TABLE `report_game_record_rg` ADD `status` tinyint unsigned NOT NULL COMMENT '订单状态' AFTER `open_number`;
ALTER TABLE `report_game_record_rg` ADD `bet` bigint NOT NULL COMMENT '投注额' AFTER `status`;
ALTER TABLE `report_game_record_rg` ADD `valid_bet` bigint NOT NULL COMMENT '有效投注额' AFTER `bet`;
ALTER TABLE `report_game_record_rg` ADD `profit` bigint NOT NULL COMMENT '盈亏金额' AFTER `valid_bet`;
ALTER TABLE `report_game_record_rg` ADD `bet_time` varchar(32) NOT NULL COMMENT '下注时间' AFTER `profit`;
ALTER TABLE `report_game_record_rg` ADD `bet_content` varchar(50) NOT NULL COMMENT '下注内容' AFTER `bet_time`;
ALTER TABLE `report_game_record_rg` ADD `ctime` bigint NOT NULL COMMENT 'ctime' AFTER `bet_content`;
ALTER TABLE `report_game_record_rg` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id' AFTER `ctime`;

----------------------------------------------------
--  `report_game_transfer_daily`
----------------------------------------------------
ALTER TABLE `report_game_transfer_daily` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `report_game_transfer_daily` ADD `channel_id` int unsigned NOT NULL COMMENT '渠道ID' AFTER `uid`;
ALTER TABLE `report_game_transfer_daily` ADD `recharge` bigint NOT NULL COMMENT '游戏充值' AFTER `channel_id`;
ALTER TABLE `report_game_transfer_daily` ADD `withdraw` bigint NOT NULL COMMENT '游戏提现' AFTER `recharge`;
ALTER TABLE `report_game_transfer_daily` ADD `ctime` bigint NOT NULL COMMENT 'create time' AFTER `withdraw`;

----------------------------------------------------
--  `report_game_user_daily`
----------------------------------------------------
ALTER TABLE `report_game_user_daily` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `report_game_user_daily` ADD `p_uid` bigint unsigned NOT NULL COMMENT 'parent uid' AFTER `uid`;
ALTER TABLE `report_game_user_daily` ADD `level` int unsigned NOT NULL COMMENT 'level' AFTER `p_uid`;
ALTER TABLE `report_game_user_daily` ADD `channel_id` int unsigned NOT NULL COMMENT 'game channel id' AFTER `level`;
ALTER TABLE `report_game_user_daily` ADD `bet` bigint NOT NULL COMMENT '投注金额' AFTER `channel_id`;
ALTER TABLE `report_game_user_daily` ADD `valid_bet` bigint NOT NULL COMMENT '有效投注额' AFTER `bet`;
ALTER TABLE `report_game_user_daily` ADD `total_valid_bet` bigint NOT NULL COMMENT '累计有效投注额(本人加无限下级)' AFTER `valid_bet`;
ALTER TABLE `report_game_user_daily` ADD `total_bet_num` int NOT NULL COMMENT '累计投注人数(本人加无限下级)' AFTER `total_valid_bet`;
ALTER TABLE `report_game_user_daily` ADD `profit` bigint NOT NULL COMMENT '盈亏金额' AFTER `total_bet_num`;
ALTER TABLE `report_game_user_daily` ADD `total_profit` bigint NOT NULL COMMENT '累计盈亏金额(本人加无限下级)' AFTER `profit`;
ALTER TABLE `report_game_user_daily` ADD `salary` bigint NOT NULL COMMENT '日工资' AFTER `total_profit`;
ALTER TABLE `report_game_user_daily` ADD `team_salary` bigint NOT NULL COMMENT '累计日工资(本人加无限下级)' AFTER `salary`;
ALTER TABLE `report_game_user_daily` ADD `status` tinyint unsigned NOT NULL COMMENT 'status' AFTER `team_salary`;
ALTER TABLE `report_game_user_daily` ADD `ctime` bigint NOT NULL COMMENT 'create time' AFTER `status`;

----------------------------------------------------
--  `report_statistic_game_all`
----------------------------------------------------
ALTER TABLE `report_statistic_game_all` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `report_statistic_game_all` ADD `channel_id` int unsigned NOT NULL COMMENT 'platform type id' AFTER `id`;
ALTER TABLE `report_statistic_game_all` ADD `newer_nums` bigint NOT NULL COMMENT 'new player nums' AFTER `channel_id`;
ALTER TABLE `report_statistic_game_all` ADD `bet` bigint NOT NULL COMMENT 'bet' AFTER `newer_nums`;
ALTER TABLE `report_statistic_game_all` ADD `valid_bet` bigint NOT NULL COMMENT 'valid bet' AFTER `bet`;
ALTER TABLE `report_statistic_game_all` ADD `profit` bigint NOT NULL COMMENT 'profit' AFTER `valid_bet`;
ALTER TABLE `report_statistic_game_all` ADD `revenue` bigint NOT NULL COMMENT '抽水' AFTER `profit`;
ALTER TABLE `report_statistic_game_all` ADD `ctime` bigint NOT NULL COMMENT 'curtime' AFTER `revenue`;
ALTER TABLE `report_statistic_game_all` ADD `note` varchar(255) NOT NULL COMMENT 'note' AFTER `ctime`;
ALTER TABLE `report_statistic_game_all` ADD `game_id` int unsigned NOT NULL COMMENT 'game id' AFTER `note`;

----------------------------------------------------
--  `report_statistic_sum`
----------------------------------------------------
ALTER TABLE `report_statistic_sum` ADD `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `report_statistic_sum` ADD `channel_id` int unsigned NOT NULL COMMENT 'platform type id' AFTER `id`;
ALTER TABLE `report_statistic_sum` ADD `channel_positive_nums` bigint unsigned NOT NULL COMMENT 'today channel positive nums' AFTER `channel_id`;
ALTER TABLE `report_statistic_sum` ADD `channel_salary_daily` bigint NOT NULL COMMENT 'channel daily salary' AFTER `channel_positive_nums`;
ALTER TABLE `report_statistic_sum` ADD `channel_rg_dividend` bigint NOT NULL COMMENT 'dividend every month' AFTER `channel_salary_daily`;
ALTER TABLE `report_statistic_sum` ADD `channel_withdraw_eusd` bigint NOT NULL COMMENT 'channel withdraw eusd' AFTER `channel_rg_dividend`;
ALTER TABLE `report_statistic_sum` ADD `channel_recharge_eusd` bigint NOT NULL COMMENT 'channel recharge eusd' AFTER `channel_withdraw_eusd`;
ALTER TABLE `report_statistic_sum` ADD `ctime` bigint NOT NULL COMMENT 'curtime' AFTER `channel_recharge_eusd`;

----------------------------------------------------
--  `report_team_daily`
----------------------------------------------------
ALTER TABLE `report_team_daily` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `report_team_daily` ADD `eusd_buy` bigint NOT NULL COMMENT '团队eusd购买金额' AFTER `uid`;
ALTER TABLE `report_team_daily` ADD `eusd_sell` bigint NOT NULL COMMENT '团队eusd出售金额' AFTER `eusd_buy`;
ALTER TABLE `report_team_daily` ADD `level` int unsigned NOT NULL COMMENT 'level' AFTER `eusd_sell`;
ALTER TABLE `report_team_daily` ADD `ctime` bigint NOT NULL COMMENT 'create time' AFTER `level`;

----------------------------------------------------
--  `report_team_game_transfer_daily`
----------------------------------------------------
ALTER TABLE `report_team_game_transfer_daily` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `report_team_game_transfer_daily` ADD `channel_id` int unsigned NOT NULL COMMENT '渠道ID' AFTER `uid`;
ALTER TABLE `report_team_game_transfer_daily` ADD `team_recharge` bigint NOT NULL COMMENT '团队游戏充值' AFTER `channel_id`;
ALTER TABLE `report_team_game_transfer_daily` ADD `team_withdraw` bigint NOT NULL COMMENT '团队游戏提现' AFTER `team_recharge`;
ALTER TABLE `report_team_game_transfer_daily` ADD `level` int unsigned NOT NULL COMMENT 'level' AFTER `team_withdraw`;
ALTER TABLE `report_team_game_transfer_daily` ADD `ctime` bigint NOT NULL COMMENT 'create time' AFTER `level`;

