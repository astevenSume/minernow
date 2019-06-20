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

