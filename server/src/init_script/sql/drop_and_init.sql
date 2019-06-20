
drop database otc;
drop database otc_admin;
drop database `default`;
drop database eos;
drop database usdt;

source create_db.sql;

use eos;
source eos_tables.sql;
source eos_init.sql;

use otc;
source otc_tables.sql;
source otc_init.sql;

use otc_admin;
source otc_admin_tables.sql;
source otc_admin_init.sql;

use usdt;
source usdt_tables.sql;
source usdt_init.sql;
