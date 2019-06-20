-- --------------------------------------------------
--  Table Structure for `models.MarketPrices`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `market_prices` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'market otc price record id',
`market` tinyint NOT NULL COMMENT 'market',
`currency` varchar(100) NOT NULL COMMENT 'currency',
`trade_method` tinyint NOT NULL COMMENT 'trade method',
`pow_price` bigint unsigned NOT NULL COMMENT 'pow price',
`pow` int NOT NULL DEFAULT 4 COMMENT 'pow_price = (real price) * pow',
`ctime` bigint NOT NULL COMMENT 'created time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='market otc ex price' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.Prices`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `prices` (
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'otc price record id',
`currency` varchar(100) NOT NULL COMMENT 'currency',
`pow_price` bigint unsigned NOT NULL COMMENT 'pow price',
`pow` int NOT NULL DEFAULT 4 COMMENT 'pow_price = (real price) * pow',
`ctime` bigint NOT NULL COMMENT 'created time',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='system otc ex price' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.UsdtAccount`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `usdt_account` (
`uaid` bigint unsigned NOT NULL COMMENT '',
`uid` bigint unsigned NOT NULL COMMENT '',
`status` tinyint unsigned NOT NULL COMMENT 'status 0-locked 1-working',
`available_integer` bigint NOT NULL COMMENT 'available amount integer part',
`frozen_integer` bigint NOT NULL COMMENT 'frozen amount integer part',
`transfer_frozen_integer` bigint NOT NULL COMMENT 'frozen amount integer part',
`mortgaged_integer` bigint NOT NULL COMMENT 'mortgaged amount integer part',
`btc_available_integer` bigint NOT NULL COMMENT 'available btc amount integer part',
`btc_frozen_integer` bigint NOT NULL COMMENT 'frozen btc amount integer part',
`btc_mortgaged_integer` bigint NOT NULL COMMENT 'mortgaged btc amount integer part',
`waiting_cash_sweep_integer` bigint NOT NULL COMMENT 'waiting cash sweep integer part',
`cash_sweep_integer` bigint NOT NULL COMMENT 'cash sweep integer part',
`owned_by_platform_integer` bigint NOT NULL COMMENT 'owned by platform integer part',
`sweep_status` tinyint unsigned NOT NULL COMMENT 'sweep status',
`pkid` bigint unsigned NOT NULL COMMENT 'usdt private key id',
`address` varchar(100) NOT NULL COMMENT '',
`ctime` bigint NOT NULL COMMENT 'created time',
`mtime` bigint NOT NULL COMMENT 'update time',
`sign` varchar(256) NOT NULL COMMENT 'sign',
PRIMARY KEY(`uaid`)
) ENGINE=InnoDB COMMENT='usdt account table' DEFAULT CHARSET=utf8;
CREATE INDEX `usdt_account_uid` ON `usdt_account` (`uid`);
CREATE INDEX `usdt_account_pkid` ON `usdt_account` (`pkid`);

-- --------------------------------------------------
--  Table Structure for `models.UsdtOnchainBalance`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `usdt_onchain_balance` (
`address` varchar(100) NOT NULL COMMENT '',
`property_id` int unsigned NOT NULL COMMENT '',
`pending_pos` varchar(100) NOT NULL COMMENT '',
`reserved` varchar(100) NOT NULL COMMENT '',
`divisible` bool NOT NULL COMMENT '',
`amount_integer` bigint NOT NULL COMMENT 'amount integer part',
`frozen` varchar(100) NOT NULL COMMENT '',
`pending_neg` varchar(100) NOT NULL COMMENT '',
`mtime` bigint NOT NULL COMMENT '',
PRIMARY KEY(`address`,`property_id`)
) ENGINE=InnoDB COMMENT='the balance on chain' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.UsdtOnChainData`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `usdt_onchain_data` (
`address` varchar(100) NOT NULL COMMENT '',
`attr_type` int unsigned NOT NULL COMMENT '',
`data_int64` bigint NOT NULL COMMENT '',
`data_str` varchar(256) NOT NULL COMMENT '',
PRIMARY KEY(`address`,`attr_type`)
) ENGINE=InnoDB COMMENT='configurations of usdt on chain.' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.UsdtOnchainLog`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `usdt_onchain_log` (
`oclid` bigint unsigned NOT NULL COMMENT '',
`from` varchar(100) NOT NULL COMMENT '',
`to` varchar(100) NOT NULL COMMENT '',
`tx` varchar(100) NOT NULL COMMENT '',
`status` varchar(100) NOT NULL COMMENT '',
`pushed` varchar(100) NOT NULL COMMENT '',
`signedTx` text NOT NULL COMMENT '',
`amount_integer` bigint NOT NULL COMMENT 'amount integer part',
`ctime` bigint NOT NULL COMMENT 'created time',
PRIMARY KEY(`oclid`)
) ENGINE=InnoDB COMMENT='the log of tx on chain' DEFAULT CHARSET=utf8;
CREATE INDEX `usdt_onchain_log_from` ON `usdt_onchain_log` (`from`);
CREATE INDEX `usdt_onchain_log_to` ON `usdt_onchain_log` (`to`);
CREATE INDEX `usdt_onchain_log_ctime` ON `usdt_onchain_log` (`ctime`);

-- --------------------------------------------------
--  Table Structure for `models.UsdtOnChainSyncPos`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `usdt_onchain_sync_pos` (
`address` varchar(100) NOT NULL COMMENT '',
`page` int unsigned NOT NULL COMMENT '',
`tx_id` varchar(128) NOT NULL COMMENT '',
PRIMARY KEY(`address`)
) ENGINE=InnoDB COMMENT='configurations of usdt on chain.' DEFAULT CHARSET=utf8;

-- --------------------------------------------------
--  Table Structure for `models.UsdtOnchainTransaction`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `usdt_onchain_transaction` (
`tx_id` varchar(100) NOT NULL COMMENT '',
`uaid` bigint unsigned NOT NULL COMMENT 'usdt account id',
`type` tinyint unsigned NOT NULL COMMENT 'Transaction type',
`property_id` int unsigned NOT NULL COMMENT '',
`property_name` varchar(100) NOT NULL COMMENT '',
`tx_type` varchar(100) NOT NULL COMMENT '',
`tx_type_int` int NOT NULL COMMENT '',
`amount_integer` bigint NOT NULL COMMENT 'amount integer part',
`block` int unsigned NOT NULL COMMENT '',
`block_hash` varchar(100) NOT NULL COMMENT '',
`block_time` bigint NOT NULL COMMENT '',
`confirmations` int unsigned NOT NULL COMMENT '',
`divisible` bool NOT NULL COMMENT '',
`fee_amount_integer` bigint NOT NULL COMMENT 'fee amount integer part',
`is_mine` bool NOT NULL COMMENT '',
`position_in_block` int unsigned NOT NULL COMMENT '',
`referenceaddress` varchar(100) NOT NULL COMMENT '',
`sending_address` varchar(100) NOT NULL COMMENT '',
`version` int NOT NULL COMMENT '',
`mtime` bigint NOT NULL COMMENT '',
PRIMARY KEY(`tx_id`,`property_id`)
) ENGINE=InnoDB COMMENT='the tx on chain' DEFAULT CHARSET=utf8;
CREATE INDEX `usdt_onchain_transaction_uaid_type` ON `usdt_onchain_transaction` (`uaid`, `type`);

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

-- --------------------------------------------------
--  Table Structure for `models.UsdtSweepLog`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `usdt_sweep_log` (
`id` bigint unsigned NOT NULL COMMENT 'id',
`uid` bigint unsigned NOT NULL COMMENT 'uid',
`ttype` int unsigned NOT NULL COMMENT 'Transaction type',
`status` int unsigned NOT NULL COMMENT 'Transaction status',
`from` varchar(256) NOT NULL COMMENT 'sender address',
`to` varchar(256) NOT NULL COMMENT 'receiver address',
`txid` varchar(256) NOT NULL COMMENT 'Transaction id',
`amount_integer` bigint NOT NULL COMMENT 'amount integer',
`fee_integer` bigint NOT NULL COMMENT 'fee integer',
`fee_onchain_integer` bigint NOT NULL COMMENT 'fee onchain',
`ctime` bigint NOT NULL COMMENT 'created time',
`utime` bigint NOT NULL COMMENT 'update time',
`step` varchar(64) NOT NULL COMMENT 'desc',
`desc` varchar(256) NOT NULL COMMENT 'desc',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='Usdt sweep log' DEFAULT CHARSET=utf8;
CREATE INDEX `usdt_sweep_log_uid_ttype` ON `usdt_sweep_log` (`uid`, `ttype`);
CREATE INDEX `usdt_sweep_log_txid` ON `usdt_sweep_log` (`txid`);

-- --------------------------------------------------
--  Table Structure for `models.UsdtTransaction`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `usdt_transaction` (
`tx_id` varchar(100) NOT NULL COMMENT '',
`uaid` bigint unsigned NOT NULL COMMENT 'usdt account id',
`type` tinyint unsigned NOT NULL COMMENT 'Transaction type',
`block_num` int unsigned NOT NULL COMMENT 'block num',
`status` int NOT NULL COMMENT 'Transaction status',
`payer` varchar(100) NOT NULL COMMENT 'Payer Account',
`receiver` varchar(100) NOT NULL COMMENT 'Receiver Account',
`amount_integer` bigint NOT NULL COMMENT 'amount integer part',
`fee` varchar(100) NOT NULL COMMENT 'transaction Fee',
`memo` varchar(100) NOT NULL COMMENT 'transaction memo',
`ctime` bigint NOT NULL COMMENT 'created time',
`utime` bigint NOT NULL COMMENT 'update time',
PRIMARY KEY(`tx_id`,`uaid`)
) ENGINE=InnoDB COMMENT='Usdt Transaction' DEFAULT CHARSET=utf8;
CREATE INDEX `usdt_transaction_uaid_type` ON `usdt_transaction` (`uaid`, `type`);

-- --------------------------------------------------
--  Table Structure for `models.UsdtWealthLog`
-- --------------------------------------------------
CREATE TABLE IF NOT EXISTS `usdt_wealth_log` (
`id` bigint unsigned NOT NULL COMMENT 'id',
`uid` bigint unsigned NOT NULL COMMENT 'uid',
`ttype` int unsigned NOT NULL COMMENT 'Transaction type',
`status` int unsigned NOT NULL COMMENT 'Transaction status',
`from` varchar(256) NOT NULL COMMENT 'sender address',
`to` varchar(256) NOT NULL COMMENT 'receiver address',
`txid` varchar(256) NOT NULL COMMENT 'Transaction id',
`amount_integer` bigint NOT NULL COMMENT 'amount integer',
`fee_integer` bigint NOT NULL COMMENT 'fee integer',
`fee_usdt_integer` bigint NOT NULL COMMENT 'fee usdt integer',
`fee_onchain_integer` bigint NOT NULL COMMENT 'fee onchain',
`ctime` bigint NOT NULL COMMENT 'created time',
`utime` bigint NOT NULL COMMENT 'update time',
`step` varchar(64) NOT NULL COMMENT 'desc',
`desc` varchar(256) NOT NULL COMMENT 'desc',
`sign` varchar(256) NOT NULL COMMENT 'sign',
`memo` varchar(256) NOT NULL COMMENT 'memo',
PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='Usdt wealth log' DEFAULT CHARSET=utf8;
CREATE INDEX `usdt_wealth_log_uid_ttype` ON `usdt_wealth_log` (`uid`, `ttype`);
CREATE INDEX `usdt_wealth_log_txid` ON `usdt_wealth_log` (`txid`);
CREATE INDEX `usdt_wealth_log_ctime` ON `usdt_wealth_log` (`ctime`);

