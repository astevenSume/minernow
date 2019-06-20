-- --------------------------------------------------
--  Table Structure for `models.TmpGamebeters`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `tmp_game_beters` (
`uid` bigint unsigned NOT NULL COMMENT 'user id',
`channel_id` int unsigned NOT NULL COMMENT 'channel id',
`game_id` int unsigned NOT NULL COMMENT 'game id',
`ctime` bigint NOT NULL COMMENT 'create time',
`status` char(20) NOT NULL DEFAULT '' COMMENT 'status',
`bet_type` int unsigned NOT NULL COMMENT 'bet type',
PRIMARY KEY(`uid`,`game_id`,`ctime`)
) ENGINE=InnoDB COMMENT='new player who is not betted' DEFAULT CHARSET=utf8;

