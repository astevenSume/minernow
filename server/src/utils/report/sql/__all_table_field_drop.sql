----------------------------------------------------
--  `game_user_month_report`
----------------------------------------------------
ALTER TABLE `game_user_month_report` DROP `id`;
ALTER TABLE `game_user_month_report` DROP `uid`;
ALTER TABLE `game_user_month_report` DROP `ctime`;
ALTER TABLE `game_user_month_report` DROP `profit`;
ALTER TABLE `game_user_month_report` DROP `agents_profit`;
ALTER TABLE `game_user_month_report` DROP `result_profit`;
ALTER TABLE `game_user_month_report` DROP `bet_amount`;
ALTER TABLE `game_user_month_report` DROP `effective_bet_amount`;
ALTER TABLE `game_user_month_report` DROP `play_game_day`;
ALTER TABLE `game_user_month_report` DROP `is_activity_user`;
ALTER TABLE `game_user_month_report` DROP `agent_level`;
ALTER TABLE `game_user_month_report` DROP `up_agent_uid`;
ALTER TABLE `game_user_month_report` DROP `activity_agent_num`;

----------------------------------------------------
--  `month_dividend_record`
----------------------------------------------------
ALTER TABLE `month_dividend_record` DROP `id`;
ALTER TABLE `month_dividend_record` DROP `uid`;
ALTER TABLE `month_dividend_record` DROP `ctime`;
ALTER TABLE `month_dividend_record` DROP `self_dividend`;
ALTER TABLE `month_dividend_record` DROP `agent_dividend`;
ALTER TABLE `month_dividend_record` DROP `result_dividend`;
ALTER TABLE `month_dividend_record` DROP `receive_status`;
ALTER TABLE `month_dividend_record` DROP `received_time`;
ALTER TABLE `month_dividend_record` DROP `pay_status`;
ALTER TABLE `month_dividend_record` DROP `level`;

----------------------------------------------------
--  `profit_report_daily`
----------------------------------------------------
ALTER TABLE `profit_report_daily` DROP `id`;
ALTER TABLE `profit_report_daily` DROP `uid`;
ALTER TABLE `profit_report_daily` DROP `channel_id`;
ALTER TABLE `profit_report_daily` DROP `bet`;
ALTER TABLE `profit_report_daily` DROP `total_valid_bet`;
ALTER TABLE `profit_report_daily` DROP `profit`;
ALTER TABLE `profit_report_daily` DROP `salary`;
ALTER TABLE `profit_report_daily` DROP `self_dividend`;
ALTER TABLE `profit_report_daily` DROP `agent_dividend`;
ALTER TABLE `profit_report_daily` DROP `result_dividend`;
ALTER TABLE `profit_report_daily` DROP `game_withdraw_amount`;
ALTER TABLE `profit_report_daily` DROP `game_recharge_amount`;
ALTER TABLE `profit_report_daily` DROP `ctime`;

----------------------------------------------------
--  `report_agent_daily`
----------------------------------------------------
ALTER TABLE `report_agent_daily` DROP `uid`;
ALTER TABLE `report_agent_daily` DROP `sum_withdraw`;
ALTER TABLE `report_agent_daily` DROP `ctime`;

----------------------------------------------------
--  `report_commission`
----------------------------------------------------
ALTER TABLE `report_commission` DROP `uid`;
ALTER TABLE `report_commission` DROP `level`;
ALTER TABLE `report_commission` DROP `team_withdraw`;
ALTER TABLE `report_commission` DROP `team_can_withdraw`;
ALTER TABLE `report_commission` DROP `ctime`;

----------------------------------------------------
--  `report_eusd_daily`
----------------------------------------------------
ALTER TABLE `report_eusd_daily` DROP `uid`;
ALTER TABLE `report_eusd_daily` DROP `buy`;
ALTER TABLE `report_eusd_daily` DROP `sell`;
ALTER TABLE `report_eusd_daily` DROP `ctime`;

----------------------------------------------------
--  `report_game_record_ag`
----------------------------------------------------
ALTER TABLE `report_game_record_ag` DROP `id`;
ALTER TABLE `report_game_record_ag` DROP `account`;
ALTER TABLE `report_game_record_ag` DROP `game_type`;
ALTER TABLE `report_game_record_ag` DROP `game_name`;
ALTER TABLE `report_game_record_ag` DROP `order_id`;
ALTER TABLE `report_game_record_ag` DROP `table_id`;
ALTER TABLE `report_game_record_ag` DROP `bet`;
ALTER TABLE `report_game_record_ag` DROP `valid_bet`;
ALTER TABLE `report_game_record_ag` DROP `profit`;
ALTER TABLE `report_game_record_ag` DROP `bet_time`;
ALTER TABLE `report_game_record_ag` DROP `ctime`;
ALTER TABLE `report_game_record_ag` DROP `uid`;

----------------------------------------------------
--  `report_game_record_ky`
----------------------------------------------------
ALTER TABLE `report_game_record_ky` DROP `id`;
ALTER TABLE `report_game_record_ky` DROP `account`;
ALTER TABLE `report_game_record_ky` DROP `game_id`;
ALTER TABLE `report_game_record_ky` DROP `game_name`;
ALTER TABLE `report_game_record_ky` DROP `server_id`;
ALTER TABLE `report_game_record_ky` DROP `kind_id`;
ALTER TABLE `report_game_record_ky` DROP `table_id`;
ALTER TABLE `report_game_record_ky` DROP `chair_id`;
ALTER TABLE `report_game_record_ky` DROP `bet`;
ALTER TABLE `report_game_record_ky` DROP `valid_bet`;
ALTER TABLE `report_game_record_ky` DROP `profit`;
ALTER TABLE `report_game_record_ky` DROP `revenue`;
ALTER TABLE `report_game_record_ky` DROP `start_time`;
ALTER TABLE `report_game_record_ky` DROP `end_time`;
ALTER TABLE `report_game_record_ky` DROP `ctime`;
ALTER TABLE `report_game_record_ky` DROP `uid`;

----------------------------------------------------
--  `report_game_record_rg`
----------------------------------------------------
ALTER TABLE `report_game_record_rg` DROP `id`;
ALTER TABLE `report_game_record_rg` DROP `account`;
ALTER TABLE `report_game_record_rg` DROP `game_name_id`;
ALTER TABLE `report_game_record_rg` DROP `game_name`;
ALTER TABLE `report_game_record_rg` DROP `game_kind_name`;
ALTER TABLE `report_game_record_rg` DROP `order_id`;
ALTER TABLE `report_game_record_rg` DROP `open_date`;
ALTER TABLE `report_game_record_rg` DROP `period_name`;
ALTER TABLE `report_game_record_rg` DROP `open_number`;
ALTER TABLE `report_game_record_rg` DROP `status`;
ALTER TABLE `report_game_record_rg` DROP `bet`;
ALTER TABLE `report_game_record_rg` DROP `valid_bet`;
ALTER TABLE `report_game_record_rg` DROP `profit`;
ALTER TABLE `report_game_record_rg` DROP `bet_time`;
ALTER TABLE `report_game_record_rg` DROP `bet_content`;
ALTER TABLE `report_game_record_rg` DROP `ctime`;
ALTER TABLE `report_game_record_rg` DROP `uid`;

----------------------------------------------------
--  `report_game_transfer_daily`
----------------------------------------------------
ALTER TABLE `report_game_transfer_daily` DROP `uid`;
ALTER TABLE `report_game_transfer_daily` DROP `channel_id`;
ALTER TABLE `report_game_transfer_daily` DROP `recharge`;
ALTER TABLE `report_game_transfer_daily` DROP `withdraw`;
ALTER TABLE `report_game_transfer_daily` DROP `ctime`;

----------------------------------------------------
--  `report_game_user_daily`
----------------------------------------------------
ALTER TABLE `report_game_user_daily` DROP `uid`;
ALTER TABLE `report_game_user_daily` DROP `p_uid`;
ALTER TABLE `report_game_user_daily` DROP `level`;
ALTER TABLE `report_game_user_daily` DROP `channel_id`;
ALTER TABLE `report_game_user_daily` DROP `bet`;
ALTER TABLE `report_game_user_daily` DROP `valid_bet`;
ALTER TABLE `report_game_user_daily` DROP `total_valid_bet`;
ALTER TABLE `report_game_user_daily` DROP `total_bet_num`;
ALTER TABLE `report_game_user_daily` DROP `profit`;
ALTER TABLE `report_game_user_daily` DROP `total_profit`;
ALTER TABLE `report_game_user_daily` DROP `salary`;
ALTER TABLE `report_game_user_daily` DROP `team_salary`;
ALTER TABLE `report_game_user_daily` DROP `status`;
ALTER TABLE `report_game_user_daily` DROP `ctime`;

----------------------------------------------------
--  `report_statistic_game_all`
----------------------------------------------------
ALTER TABLE `report_statistic_game_all` DROP `id`;
ALTER TABLE `report_statistic_game_all` DROP `channel_id`;
ALTER TABLE `report_statistic_game_all` DROP `newer_nums`;
ALTER TABLE `report_statistic_game_all` DROP `bet`;
ALTER TABLE `report_statistic_game_all` DROP `valid_bet`;
ALTER TABLE `report_statistic_game_all` DROP `profit`;
ALTER TABLE `report_statistic_game_all` DROP `revenue`;
ALTER TABLE `report_statistic_game_all` DROP `ctime`;
ALTER TABLE `report_statistic_game_all` DROP `note`;
ALTER TABLE `report_statistic_game_all` DROP `game_id`;

----------------------------------------------------
--  `report_statistic_sum`
----------------------------------------------------
ALTER TABLE `report_statistic_sum` DROP `id`;
ALTER TABLE `report_statistic_sum` DROP `channel_id`;
ALTER TABLE `report_statistic_sum` DROP `channel_positive_nums`;
ALTER TABLE `report_statistic_sum` DROP `channel_salary_daily`;
ALTER TABLE `report_statistic_sum` DROP `channel_rg_dividend`;
ALTER TABLE `report_statistic_sum` DROP `channel_withdraw_eusd`;
ALTER TABLE `report_statistic_sum` DROP `channel_recharge_eusd`;
ALTER TABLE `report_statistic_sum` DROP `ctime`;

----------------------------------------------------
--  `report_team_daily`
----------------------------------------------------
ALTER TABLE `report_team_daily` DROP `uid`;
ALTER TABLE `report_team_daily` DROP `eusd_buy`;
ALTER TABLE `report_team_daily` DROP `eusd_sell`;
ALTER TABLE `report_team_daily` DROP `level`;
ALTER TABLE `report_team_daily` DROP `ctime`;

----------------------------------------------------
--  `report_team_game_transfer_daily`
----------------------------------------------------
ALTER TABLE `report_team_game_transfer_daily` DROP `uid`;
ALTER TABLE `report_team_game_transfer_daily` DROP `channel_id`;
ALTER TABLE `report_team_game_transfer_daily` DROP `team_recharge`;
ALTER TABLE `report_team_game_transfer_daily` DROP `team_withdraw`;
ALTER TABLE `report_team_game_transfer_daily` DROP `level`;
ALTER TABLE `report_team_game_transfer_daily` DROP `ctime`;

