----------------------------------------------------
--  `agent`
----------------------------------------------------
ALTER TABLE `agent` DROP `uid`;
ALTER TABLE `agent` DROP `sum_salary`;
ALTER TABLE `agent` DROP `sum_withdraw`;
ALTER TABLE `agent` DROP `sum_can_withdraw`;
ALTER TABLE `agent` DROP `ctime`;
ALTER TABLE `agent` DROP `mtime`;
ALTER TABLE `agent` DROP `pwd`;

----------------------------------------------------
--  `agent_channel_commission`
----------------------------------------------------
ALTER TABLE `agent_channel_commission` DROP `uid`;
ALTER TABLE `agent_channel_commission` DROP `channel_id`;
ALTER TABLE `agent_channel_commission` DROP `ctime`;
ALTER TABLE `agent_channel_commission` DROP `mtime`;
ALTER TABLE `agent_channel_commission` DROP `integer`;
ALTER TABLE `agent_channel_commission` DROP `decimals`;
ALTER TABLE `agent_channel_commission` DROP `status`;

----------------------------------------------------
--  `agent_path`
----------------------------------------------------
ALTER TABLE `agent_path` DROP `uid`;
ALTER TABLE `agent_path` DROP `level`;
ALTER TABLE `agent_path` DROP `sn`;
ALTER TABLE `agent_path` DROP `path`;
ALTER TABLE `agent_path` DROP `ctime`;
ALTER TABLE `agent_path` DROP `mtime`;
ALTER TABLE `agent_path` DROP `invite_code`;
ALTER TABLE `agent_path` DROP `whitelist_id`;
ALTER TABLE `agent_path` DROP `invite_num`;
ALTER TABLE `agent_path` DROP `parent_uid`;
ALTER TABLE `agent_path` DROP `dividend_position`;

----------------------------------------------------
--  `agent_withdraw`
----------------------------------------------------
ALTER TABLE `agent_withdraw` DROP `id`;
ALTER TABLE `agent_withdraw` DROP `uid`;
ALTER TABLE `agent_withdraw` DROP `amount`;
ALTER TABLE `agent_withdraw` DROP `status`;
ALTER TABLE `agent_withdraw` DROP `ctime`;
ALTER TABLE `agent_withdraw` DROP `mtime`;
ALTER TABLE `agent_withdraw` DROP `desc`;

----------------------------------------------------
--  `invite_code`
----------------------------------------------------
ALTER TABLE `invite_code` DROP `id`;
ALTER TABLE `invite_code` DROP `code`;
ALTER TABLE `invite_code` DROP `status`;

