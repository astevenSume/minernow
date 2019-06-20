----------------------------------------------------
--  `game_user_month_report`
----------------------------------------------------
ALTER TABLE `game_user_month_report` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `game_user_month_report` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `game_user_month_report` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'create time';
ALTER TABLE `game_user_month_report` CHANGE `profit` `profit` bigint NOT NULL COMMENT '自己的盈亏金额';
ALTER TABLE `game_user_month_report` CHANGE `agents_profit` `agents_profit` bigint NOT NULL COMMENT '无限下级下级代理的盈亏金额';
ALTER TABLE `game_user_month_report` CHANGE `result_profit` `result_profit` bigint NOT NULL COMMENT '最终的盈亏金额';
ALTER TABLE `game_user_month_report` CHANGE `bet_amount` `bet_amount` bigint NOT NULL COMMENT '投注额';
ALTER TABLE `game_user_month_report` CHANGE `effective_bet_amount` `effective_bet_amount` bigint NOT NULL COMMENT '累计有效投注额(本人加无限下级)';
ALTER TABLE `game_user_month_report` CHANGE `play_game_day` `play_game_day` int NOT NULL COMMENT '玩游戏的天数';
ALTER TABLE `game_user_month_report` CHANGE `is_activity_user` `is_activity_user` bool NOT NULL DEFAULT false COMMENT '是否是活跃用户';
ALTER TABLE `game_user_month_report` CHANGE `agent_level` `agent_level` int unsigned NOT NULL COMMENT '代理等级';
ALTER TABLE `game_user_month_report` CHANGE `up_agent_uid` `up_agent_uid` bigint unsigned NOT NULL COMMENT '上级代理的UID';
ALTER TABLE `game_user_month_report` CHANGE `activity_agent_num` `activity_agent_num` int NOT NULL COMMENT '无限下级代理活跃人数';

----------------------------------------------------
--  `month_dividend_record`
----------------------------------------------------
ALTER TABLE `month_dividend_record` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `month_dividend_record` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `month_dividend_record` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'create time';
ALTER TABLE `month_dividend_record` CHANGE `self_dividend` `self_dividend` bigint NOT NULL COMMENT '自己的分红';
ALTER TABLE `month_dividend_record` CHANGE `agent_dividend` `agent_dividend` bigint NOT NULL COMMENT '要分给代理的分红';
ALTER TABLE `month_dividend_record` CHANGE `result_dividend` `result_dividend` bigint NOT NULL COMMENT '最终获得的分红 self_dividend-agent_dividend';
ALTER TABLE `month_dividend_record` CHANGE `receive_status` `receive_status` int NOT NULL DEFAULT false COMMENT '上级已发放状态是1,2是等待上级发放';
ALTER TABLE `month_dividend_record` CHANGE `received_time` `received_time` bigint NOT NULL COMMENT '领取奖励的时间';
ALTER TABLE `month_dividend_record` CHANGE `pay_status` `pay_status` int NOT NULL DEFAULT false COMMENT '1是已支付状态,2是等待支付状态';
ALTER TABLE `month_dividend_record` CHANGE `level` `level` int unsigned NOT NULL COMMENT 'level';

----------------------------------------------------
--  `profit_report_daily`
----------------------------------------------------
ALTER TABLE `profit_report_daily` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `profit_report_daily` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `profit_report_daily` CHANGE `channel_id` `channel_id` int unsigned NOT NULL COMMENT 'game channel id';
ALTER TABLE `profit_report_daily` CHANGE `bet` `bet` bigint NOT NULL COMMENT '本人的有效投注额';
ALTER TABLE `profit_report_daily` CHANGE `total_valid_bet` `total_valid_bet` bigint NOT NULL COMMENT '累计有效投注额(本人加无限下级)';
ALTER TABLE `profit_report_daily` CHANGE `profit` `profit` bigint NOT NULL COMMENT '盈亏金额';
ALTER TABLE `profit_report_daily` CHANGE `salary` `salary` bigint NOT NULL COMMENT '日工资';
ALTER TABLE `profit_report_daily` CHANGE `self_dividend` `self_dividend` bigint NOT NULL COMMENT '属于自己的月分红';
ALTER TABLE `profit_report_daily` CHANGE `agent_dividend` `agent_dividend` bigint NOT NULL COMMENT '分给下级代理的月分红';
ALTER TABLE `profit_report_daily` CHANGE `result_dividend` `result_dividend` bigint NOT NULL COMMENT '最终获得的月分红,可能为负数';
ALTER TABLE `profit_report_daily` CHANGE `game_withdraw_amount` `game_withdraw_amount` int unsigned NOT NULL COMMENT '游戏提现金额';
ALTER TABLE `profit_report_daily` CHANGE `game_recharge_amount` `game_recharge_amount` int unsigned NOT NULL COMMENT '游戏充值金额';
ALTER TABLE `profit_report_daily` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'create time';

----------------------------------------------------
--  `report_agent_daily`
----------------------------------------------------
ALTER TABLE `report_agent_daily` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `report_agent_daily` CHANGE `sum_withdraw` `sum_withdraw` bigint NOT NULL COMMENT 'sum_withdraw';
ALTER TABLE `report_agent_daily` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'create time';

----------------------------------------------------
--  `report_commission`
----------------------------------------------------
ALTER TABLE `report_commission` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `report_commission` CHANGE `level` `level` int unsigned NOT NULL COMMENT 'level';
ALTER TABLE `report_commission` CHANGE `team_withdraw` `team_withdraw` bigint NOT NULL COMMENT 'team_withdraw';
ALTER TABLE `report_commission` CHANGE `team_can_withdraw` `team_can_withdraw` bigint NOT NULL COMMENT 'team_can_withdraw';
ALTER TABLE `report_commission` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'create time';

----------------------------------------------------
--  `report_eusd_daily`
----------------------------------------------------
ALTER TABLE `report_eusd_daily` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `report_eusd_daily` CHANGE `buy` `buy` bigint NOT NULL COMMENT '购买eusd数量';
ALTER TABLE `report_eusd_daily` CHANGE `sell` `sell` bigint NOT NULL COMMENT '出售eusd数量';
ALTER TABLE `report_eusd_daily` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'create time';

----------------------------------------------------
--  `report_game_record_ag`
----------------------------------------------------
ALTER TABLE `report_game_record_ag` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `report_game_record_ag` CHANGE `account` `account` varchar(50) NOT NULL COMMENT 'account';
ALTER TABLE `report_game_record_ag` CHANGE `game_type` `game_type` varchar(50) NOT NULL COMMENT '游戏ID';
ALTER TABLE `report_game_record_ag` CHANGE `game_name` `game_name` varchar(50) NOT NULL COMMENT '游戏名称';
ALTER TABLE `report_game_record_ag` CHANGE `order_id` `order_id` varchar(50) NOT NULL COMMENT '订单编号';
ALTER TABLE `report_game_record_ag` CHANGE `table_id` `table_id` varchar(50) NOT NULL COMMENT '桌号';
ALTER TABLE `report_game_record_ag` CHANGE `bet` `bet` bigint NOT NULL COMMENT '投注额';
ALTER TABLE `report_game_record_ag` CHANGE `valid_bet` `valid_bet` bigint NOT NULL COMMENT '有效投注额';
ALTER TABLE `report_game_record_ag` CHANGE `profit` `profit` bigint NOT NULL COMMENT '盈亏金额';
ALTER TABLE `report_game_record_ag` CHANGE `bet_time` `bet_time` varchar(32) NOT NULL COMMENT '下注时间';
ALTER TABLE `report_game_record_ag` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'ctime';
ALTER TABLE `report_game_record_ag` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';

----------------------------------------------------
--  `report_game_record_ky`
----------------------------------------------------
ALTER TABLE `report_game_record_ky` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `report_game_record_ky` CHANGE `account` `account` varchar(50) NOT NULL COMMENT 'account';
ALTER TABLE `report_game_record_ky` CHANGE `game_id` `game_id` varchar(50) NOT NULL COMMENT '游戏局号';
ALTER TABLE `report_game_record_ky` CHANGE `game_name` `game_name` varchar(50) NOT NULL COMMENT 'game_name';
ALTER TABLE `report_game_record_ky` CHANGE `server_id` `server_id` int NOT NULL COMMENT '房间ID';
ALTER TABLE `report_game_record_ky` CHANGE `kind_id` `kind_id` varchar(50) NOT NULL COMMENT '游戏ID';
ALTER TABLE `report_game_record_ky` CHANGE `table_id` `table_id` int NOT NULL COMMENT '桌子号';
ALTER TABLE `report_game_record_ky` CHANGE `chair_id` `chair_id` int NOT NULL COMMENT '椅子号';
ALTER TABLE `report_game_record_ky` CHANGE `bet` `bet` bigint NOT NULL COMMENT '投注额';
ALTER TABLE `report_game_record_ky` CHANGE `valid_bet` `valid_bet` bigint NOT NULL COMMENT '有效投注额';
ALTER TABLE `report_game_record_ky` CHANGE `profit` `profit` bigint NOT NULL COMMENT '盈亏金额';
ALTER TABLE `report_game_record_ky` CHANGE `revenue` `revenue` bigint NOT NULL COMMENT '抽水金额';
ALTER TABLE `report_game_record_ky` CHANGE `start_time` `start_time` varchar(32) NOT NULL COMMENT 'start_time';
ALTER TABLE `report_game_record_ky` CHANGE `end_time` `end_time` varchar(32) NOT NULL COMMENT 'end_time';
ALTER TABLE `report_game_record_ky` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'create time';
ALTER TABLE `report_game_record_ky` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';

----------------------------------------------------
--  `report_game_record_rg`
----------------------------------------------------
ALTER TABLE `report_game_record_rg` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `report_game_record_rg` CHANGE `account` `account` varchar(50) NOT NULL COMMENT 'account';
ALTER TABLE `report_game_record_rg` CHANGE `game_name_id` `game_name_id` varchar(50) NOT NULL COMMENT '游戏id';
ALTER TABLE `report_game_record_rg` CHANGE `game_name` `game_name` varchar(50) NOT NULL COMMENT '游戏名称';
ALTER TABLE `report_game_record_rg` CHANGE `game_kind_name` `game_kind_name` varchar(50) NOT NULL COMMENT '玩法名称';
ALTER TABLE `report_game_record_rg` CHANGE `order_id` `order_id` varchar(50) NOT NULL COMMENT '订单编号';
ALTER TABLE `report_game_record_rg` CHANGE `open_date` `open_date` varchar(32) NOT NULL COMMENT '开奖时间';
ALTER TABLE `report_game_record_rg` CHANGE `period_name` `period_name` varchar(50) NOT NULL COMMENT '期号';
ALTER TABLE `report_game_record_rg` CHANGE `open_number` `open_number` varchar(50) NOT NULL COMMENT '开奖号码';
ALTER TABLE `report_game_record_rg` CHANGE `status` `status` tinyint unsigned NOT NULL COMMENT '订单状态';
ALTER TABLE `report_game_record_rg` CHANGE `bet` `bet` bigint NOT NULL COMMENT '投注额';
ALTER TABLE `report_game_record_rg` CHANGE `valid_bet` `valid_bet` bigint NOT NULL COMMENT '有效投注额';
ALTER TABLE `report_game_record_rg` CHANGE `profit` `profit` bigint NOT NULL COMMENT '盈亏金额';
ALTER TABLE `report_game_record_rg` CHANGE `bet_time` `bet_time` varchar(32) NOT NULL COMMENT '下注时间';
ALTER TABLE `report_game_record_rg` CHANGE `bet_content` `bet_content` varchar(50) NOT NULL COMMENT '下注内容';
ALTER TABLE `report_game_record_rg` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'ctime';
ALTER TABLE `report_game_record_rg` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';

----------------------------------------------------
--  `report_game_transfer_daily`
----------------------------------------------------
ALTER TABLE `report_game_transfer_daily` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `report_game_transfer_daily` CHANGE `channel_id` `channel_id` int unsigned NOT NULL COMMENT '渠道ID';
ALTER TABLE `report_game_transfer_daily` CHANGE `recharge` `recharge` bigint NOT NULL COMMENT '游戏充值';
ALTER TABLE `report_game_transfer_daily` CHANGE `withdraw` `withdraw` bigint NOT NULL COMMENT '游戏提现';
ALTER TABLE `report_game_transfer_daily` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'create time';

----------------------------------------------------
--  `report_game_user_daily`
----------------------------------------------------
ALTER TABLE `report_game_user_daily` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `report_game_user_daily` CHANGE `p_uid` `p_uid` bigint unsigned NOT NULL COMMENT 'parent uid';
ALTER TABLE `report_game_user_daily` CHANGE `level` `level` int unsigned NOT NULL COMMENT 'level';
ALTER TABLE `report_game_user_daily` CHANGE `channel_id` `channel_id` int unsigned NOT NULL COMMENT 'game channel id';
ALTER TABLE `report_game_user_daily` CHANGE `bet` `bet` bigint NOT NULL COMMENT '投注金额';
ALTER TABLE `report_game_user_daily` CHANGE `valid_bet` `valid_bet` bigint NOT NULL COMMENT '有效投注额';
ALTER TABLE `report_game_user_daily` CHANGE `total_valid_bet` `total_valid_bet` bigint NOT NULL COMMENT '累计有效投注额(本人加无限下级)';
ALTER TABLE `report_game_user_daily` CHANGE `total_bet_num` `total_bet_num` int NOT NULL COMMENT '累计投注人数(本人加无限下级)';
ALTER TABLE `report_game_user_daily` CHANGE `profit` `profit` bigint NOT NULL COMMENT '盈亏金额';
ALTER TABLE `report_game_user_daily` CHANGE `total_profit` `total_profit` bigint NOT NULL COMMENT '累计盈亏金额(本人加无限下级)';
ALTER TABLE `report_game_user_daily` CHANGE `salary` `salary` bigint NOT NULL COMMENT '日工资';
ALTER TABLE `report_game_user_daily` CHANGE `team_salary` `team_salary` bigint NOT NULL COMMENT '累计日工资(本人加无限下级)';
ALTER TABLE `report_game_user_daily` CHANGE `status` `status` tinyint unsigned NOT NULL COMMENT 'status';
ALTER TABLE `report_game_user_daily` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'create time';

----------------------------------------------------
--  `report_statistic_game_all`
----------------------------------------------------
ALTER TABLE `report_statistic_game_all` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `report_statistic_game_all` CHANGE `channel_id` `channel_id` int unsigned NOT NULL COMMENT 'platform type id';
ALTER TABLE `report_statistic_game_all` CHANGE `newer_nums` `newer_nums` bigint NOT NULL COMMENT 'new player nums';
ALTER TABLE `report_statistic_game_all` CHANGE `bet` `bet` bigint NOT NULL COMMENT 'bet';
ALTER TABLE `report_statistic_game_all` CHANGE `valid_bet` `valid_bet` bigint NOT NULL COMMENT 'valid bet';
ALTER TABLE `report_statistic_game_all` CHANGE `profit` `profit` bigint NOT NULL COMMENT 'profit';
ALTER TABLE `report_statistic_game_all` CHANGE `revenue` `revenue` bigint NOT NULL COMMENT '抽水';
ALTER TABLE `report_statistic_game_all` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'curtime';
ALTER TABLE `report_statistic_game_all` CHANGE `note` `note` varchar(255) NOT NULL COMMENT 'note';
ALTER TABLE `report_statistic_game_all` CHANGE `game_id` `game_id` int unsigned NOT NULL COMMENT 'game id';

----------------------------------------------------
--  `report_statistic_sum`
----------------------------------------------------
ALTER TABLE `report_statistic_sum` CHANGE `id` `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id';
ALTER TABLE `report_statistic_sum` CHANGE `channel_id` `channel_id` int unsigned NOT NULL COMMENT 'platform type id';
ALTER TABLE `report_statistic_sum` CHANGE `channel_positive_nums` `channel_positive_nums` bigint unsigned NOT NULL COMMENT 'today channel positive nums';
ALTER TABLE `report_statistic_sum` CHANGE `channel_salary_daily` `channel_salary_daily` bigint NOT NULL COMMENT 'channel daily salary';
ALTER TABLE `report_statistic_sum` CHANGE `channel_rg_dividend` `channel_rg_dividend` bigint NOT NULL COMMENT 'dividend every month';
ALTER TABLE `report_statistic_sum` CHANGE `channel_withdraw_eusd` `channel_withdraw_eusd` bigint NOT NULL COMMENT 'channel withdraw eusd';
ALTER TABLE `report_statistic_sum` CHANGE `channel_recharge_eusd` `channel_recharge_eusd` bigint NOT NULL COMMENT 'channel recharge eusd';
ALTER TABLE `report_statistic_sum` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'curtime';

----------------------------------------------------
--  `report_team_daily`
----------------------------------------------------
ALTER TABLE `report_team_daily` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `report_team_daily` CHANGE `eusd_buy` `eusd_buy` bigint NOT NULL COMMENT '团队eusd购买金额';
ALTER TABLE `report_team_daily` CHANGE `eusd_sell` `eusd_sell` bigint NOT NULL COMMENT '团队eusd出售金额';
ALTER TABLE `report_team_daily` CHANGE `level` `level` int unsigned NOT NULL COMMENT 'level';
ALTER TABLE `report_team_daily` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'create time';

----------------------------------------------------
--  `report_team_game_transfer_daily`
----------------------------------------------------
ALTER TABLE `report_team_game_transfer_daily` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `report_team_game_transfer_daily` CHANGE `channel_id` `channel_id` int unsigned NOT NULL COMMENT '渠道ID';
ALTER TABLE `report_team_game_transfer_daily` CHANGE `team_recharge` `team_recharge` bigint NOT NULL COMMENT '团队游戏充值';
ALTER TABLE `report_team_game_transfer_daily` CHANGE `team_withdraw` `team_withdraw` bigint NOT NULL COMMENT '团队游戏提现';
ALTER TABLE `report_team_game_transfer_daily` CHANGE `level` `level` int unsigned NOT NULL COMMENT 'level';
ALTER TABLE `report_team_game_transfer_daily` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'create time';

