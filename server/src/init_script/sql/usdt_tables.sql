-- --------------------------------------------------
--  Table Structure for `models.PriKey`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `usdt_prikey` (
`pkid` bigint unsigned NOT NULL COMMENT '',
`pri` varchar(256) NOT NULL COMMENT '',
`address` varchar(100) NOT NULL COMMENT '',
PRIMARY KEY(`pkid`)
) ENGINE=InnoDB COMMENT='' DEFAULT CHARSET=utf8;
CREATE UNIQUE INDEX `usdt_prikey_pri` ON `usdt_prikey` (`pri`);
CREATE UNIQUE INDEX `usdt_prikey_address` ON `usdt_prikey` (`address`);
