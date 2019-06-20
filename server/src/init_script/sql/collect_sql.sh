#!/bin/bash
SHELL_FOLDER=$(cd "$(dirname "$0")";pwd)

FILE_ADMIN=$SHELL_FOLDER"/otc_admin_tables.sql"
FILE_OTC=$SHELL_FOLDER"/otc_tables.sql"
FILE_EOS=$SHELL_FOLDER"/eos_tables.sql"

#cd到 src
cd $SHELL_FOLDER
cd ../../../
pwd

sh src/tools/dbgenerate/build.sh
sh src/tools/dbgenerate/generate.sh

go fmt src/utils/admin/models/*
go fmt src/utils/common/models/*
go fmt src/utils/eusd/models/*
go fmt src/utils/game/models/*
go fmt src/utils/otc/models/*
go fmt src/utils/usdt/models/*
go fmt src/utils/eos/models/*

echo '' > $FILE_ADMIN
cat src/utils/admin/sql/__all_table_create.sql >> $FILE_ADMIN
cat src/utils/common/sql/__all_table_create.sql >> $FILE_ADMIN

echo '' > $FILE_OTC
cat src/utils/common/sql/__all_table_create.sql >> $FILE_OTC
cat src/utils/eusd/sql/__all_table_create.sql >> $FILE_OTC
cat src/utils/game/sql/__all_table_create.sql >> $FILE_OTC
cat src/utils/otc/sql/__all_table_create.sql >> $FILE_OTC
cat src/utils/usdt/sql/__all_table_create.sql >> $FILE_OTC
cat src/utils/agent/sql/__all_table_create.sql >> $FILE_OTC
cat src/utils/report/sql/__all_table_create.sql >> $FILE_OTC

echo '' > $FILE_EOS
cat src/utils/eos/sql/__all_table_create.sql >> $FILE_EOS

