DELETE  FROM user where national_code !="1";

DELETE  FROM eos_account where public_key_owner = '';
DELETE  FROM eos_wealth_log;

truncate table otc_sell;
truncate table otc_buy;
truncate table otc_order;
truncate table eos_otc;
truncate table eos_wealth;

