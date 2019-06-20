----------------------------------------------------
--  `tmp_game_beters`
----------------------------------------------------
ALTER TABLE `tmp_game_beters` ADD `uid` bigint unsigned NOT NULL COMMENT 'user id';
ALTER TABLE `tmp_game_beters` ADD `channel_id` int unsigned NOT NULL COMMENT 'channel id' AFTER `uid`;
ALTER TABLE `tmp_game_beters` ADD `game_id` int unsigned NOT NULL COMMENT 'game id' AFTER `channel_id`;
ALTER TABLE `tmp_game_beters` ADD `ctime` bigint NOT NULL COMMENT 'create time' AFTER `game_id`;
ALTER TABLE `tmp_game_beters` ADD `status` char(20) NOT NULL DEFAULT '' COMMENT 'status' AFTER `ctime`;
ALTER TABLE `tmp_game_beters` ADD `bet_type` int unsigned NOT NULL COMMENT 'bet type' AFTER `status`;

