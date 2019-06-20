----------------------------------------------------
--  `tmp_game_beters`
----------------------------------------------------
ALTER TABLE `tmp_game_beters` CHANGE `uid` `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `tmp_game_beters` CHANGE `channel_id` `channel_id` int unsigned NOT NULL COMMENT 'channel id';
ALTER TABLE `tmp_game_beters` CHANGE `game_id` `game_id` int unsigned NOT NULL COMMENT 'game id';
ALTER TABLE `tmp_game_beters` CHANGE `ctime` `ctime` bigint NOT NULL COMMENT 'create time';
ALTER TABLE `tmp_game_beters` CHANGE `status` `status` char(20) NOT NULL DEFAULT '' COMMENT 'status';
ALTER TABLE `tmp_game_beters` CHANGE `bet_type` `bet_type` int unsigned NOT NULL COMMENT 'bet type';

