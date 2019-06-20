package main

import (
	"common"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

const (
	COMMENT      = "comment"
	PK           = "pk"
	NULL         = "null"
	TYPE         = "type"
	DEFAULT      = "default"
	MODEL_NAME   = "model_name"
	INDEX        = "index"
	INDEX_LEN    = "index_len"
	UNIQUE_INDEX = "unique index"
	SN           = "sn"
	SIZE         = "size"
	PROTO        = "proto"
	DB           = "db"
	AUTO         = "auto"
)

var (
	MODEL_DIR, SQL_DIR, CONFIG_PATH string
)

//数据库数据类型与model数据类型的映射
var DATA2MODEL_TYPE_MAP = map[string]string{
	"tinyint unsigned":  "uint8",
	"tinyint":           "int8",
	"smallint unsigned": "uint16",
	"smallint":          "int16",
	"int unsigned":      "uint32",
	"int":               "int32",
	"bigint unsigned":   "uint64",
	"bigint":            "int64",
	"char":              "string",
	"varchar":           "string",
	"blob":              "string",
	"text":              "string",
	"bool":              "bool",
}

//表定义
type table struct {
	name         string         //表名
	comment      string         //表注释
	columns      []*column      //列列表
	indexs       []*index       //索引定义
	uniqueIndexs []*uniqueIndex //索引定义
	modelName    string
}

func NewTable(name string) *table {
	return &table{
		name: name,
	}
}

//列
type column struct {
	name                                string
	comment                             string
	sType                               string
	iSize                               int
	iDefault                            interface{}
	isPk, isNull, isProto, isDb, isAuto bool
	modelName                           string
	sn                                  int
	indexLen                            int
}

func NewColumn() *column {
	return &column{
		isDb: true,
	}
}

func (c *column) parse(columnMap map[string]interface{}) (err error) {
	for k, v := range columnMap {
		switch k {
		case TYPE:
			{
				if sType, ok := v.(string); ok {
					if _, ok := DATA2MODEL_TYPE_MAP[sType]; ok {
						c.sType = sType
					} else {
						common.LogFuncCritical("unsupported TYPE %s", sType)
						return ERR_PARSE_FAILED
					}
				} else {
					common.LogFuncCritical("parse TYPE %v failed", v)
					return ERR_PARSE_FAILED
				}
			}
		case PK:
			{
				if isPk, ok := v.(bool); ok {
					c.isPk = isPk
				} else {
					common.LogFuncCritical("parse PK %v failed", v)
					return ERR_PARSE_FAILED
				}
			}
		case AUTO:
			{
				if isAuto, ok := v.(bool); ok {
					c.isAuto = isAuto
				} else {
					common.LogFuncCritical("parse AUTO %v failed", v)
					return ERR_PARSE_FAILED
				}
			}
		case PROTO:
			{
				if isProto, ok := v.(bool); ok {
					c.isProto = isProto
				} else {
					common.LogFuncCritical("parse PROTO %v failed", v)
					return ERR_PARSE_FAILED
				}
			}
		case DB:
			{
				if isDb, ok := v.(bool); ok {
					c.isDb = isDb
				} else {
					common.LogFuncCritical("parse DB %v failed", v)
					return ERR_PARSE_FAILED
				}
			}
		case NULL:
			{
				if isNull, ok := v.(bool); ok {
					c.isNull = isNull
				} else {
					common.LogFuncCritical("parse NULL %v failed", v)
					return ERR_PARSE_FAILED
				}
			}
		case DEFAULT:
			{
				c.iDefault = v
			}
		case MODEL_NAME:
			{
				if modelName, ok := v.(string); ok {
					c.modelName = modelName
				} else {
					common.LogFuncCritical("parse MODEL_NAME %v failed", v)
					return ERR_PARSE_FAILED
				}
			}
		case COMMENT:
			{
				if comment, ok := v.(string); ok {
					c.comment = comment
				} else {
					common.LogFuncCritical("parse COMMENT %v failed", v)
					return ERR_PARSE_FAILED
				}
			}
		case SN:
			{
				//c.sn = v.(int)
				if sn, ok := v.(float64); ok {
					c.sn = int(sn)
				} else {
					common.LogFuncCritical("parse SN %v failed", v)
					return ERR_PARSE_FAILED
				}
			}
		case SIZE:
			{
				if size, ok := v.(float64); ok {
					c.iSize = int(size)
				} else {
					common.LogFuncCritical("parse SIZE %v failed", v)
					return ERR_PARSE_FAILED
				}
			}
		case INDEX_LEN:
			{
				if indexLen, ok := v.(float64); ok {
					c.indexLen = int(indexLen)
				} else {
					common.LogFuncCritical("parse INDEX_LEN %v failed", v)
					return ERR_PARSE_FAILED
				}
			}
		}
	}

	return
}

//索引
type index struct {
	columns []string
}

func NewIndex() *index {
	return &index{}
}

//unique index
type uniqueIndex struct {
	columns []string
}

func NewUniqueIndex() *uniqueIndex {
	return &uniqueIndex{}
}

type TableSlice []*table

func (t *table) parseIndex(indexSlice []interface{}) (err error) {
	for _, indexs := range indexSlice {
		idx := NewIndex()
		if indexStrs, ok := indexs.([]interface{}); ok {
			for _, indexStr := range indexStrs {
				if s, ok := indexStr.(string); ok {
					idx.columns = append(idx.columns, s)
				} else {
					common.LogFuncCritical("parse INDEX %v failed", indexStr)
					return ERR_PARSE_FAILED
				}
			}
		} else {
			common.LogFuncCritical("parse INDEX %v failed.", indexStrs)
			return ERR_PARSE_FAILED
		}

		t.indexs = append(t.indexs, idx)
	}

	return
}

func (t *table) parseUniqueIndex(indexSlice []interface{}) (err error) {
	for _, indexs := range indexSlice {
		idx := NewUniqueIndex()
		if indexStrs, ok := indexs.([]interface{}); ok {
			for _, indexStr := range indexStrs {
				if s, ok := indexStr.(string); ok {
					idx.columns = append(idx.columns, s)
				} else {
					common.LogFuncCritical("parse INDEX %v failed", indexStr)
					return ERR_PARSE_FAILED
				}
			}
		} else {
			common.LogFuncCritical("parse INDEX %v failed.", indexStrs)
			return ERR_PARSE_FAILED
		}

		t.uniqueIndexs = append(t.uniqueIndexs, idx)
	}

	return
}

//生成model文件
func (t *table) outputModelCode() (err error) {
	source := "package models\n\n"

	modelCode := fmt.Sprintf("//auto_models_start\n type %s struct{\n", t.modelName)

	//生成数据库结构体
	isExistPk := false
	pks := []*column{}
	for _, c := range t.columns {
		if !c.isDb {
			continue
		}

		modelCode += fmt.Sprintf("\t%s %s `orm:\"column(%s)", c.modelName, DATA2MODEL_TYPE_MAP[c.sType], c.name)

		if c.isNull {
			modelCode += ";null"
		}

		if c.isPk {
			if !isExistPk { //orm只支持一个pk，原因是什么？
				modelCode += ";pk"
			}
			pks = append(pks, c)
			isExistPk = true
		}

		if c.iSize > 0 {
			modelCode += fmt.Sprintf(";size(%d)", c.iSize)
		}

		//if c.isProto {
		modelCode += fmt.Sprintf("\" json:\"%s", c.name)
		//}

		modelCode += ",omitempty\"`\n"
	}

	modelCode += "}\n\n"

	//if len(pks) > 1 {
	//	source += "import \"common\"\n\n"
	//}

	source += modelCode

	//TableName function
	tableNameCode := fmt.Sprintf("func (this *%s) TableName() string {\n    return \"%s\"\n}\n\n", t.modelName, t.name)
	source += tableNameCode

	//	//如果多字段联合主键，就生成update函数
	//	updateCode := ""
	//	if len(pks) > 1 {
	//		updateCode += fmt.Sprintf("func (this *%s) Update() (num int64, err error) {\n", t.modelName)
	//		updateCode += fmt.Sprintf("    res, err := DbOrm.Raw(\"UPDATE %s SET ", t.name)
	//		first := true
	//		for _, c := range t.columns {
	//			if isPk(c, pks) {
	//				continue
	//			}
	//
	//			if !first {
	//				updateCode += ","
	//			}
	//
	//			updateCode += fmt.Sprintf("%s=?", c.name)
	//
	//			if first {
	//				first = !first
	//			}
	//		}
	//
	//		updateCode += " WHERE "
	//		first = true
	//		for _, pk := range pks {
	//
	//			if !first {
	//				updateCode += " AND "
	//			}
	//
	//			updateCode += fmt.Sprintf("%s=?", pk.name)
	//
	//			if first {
	//				first = !first
	//			}
	//		}
	//
	//		updateCode += "\""
	//
	//		for _, c := range t.columns {
	//			if c.isPk {
	//				continue
	//			}
	//
	//			updateCode += fmt.Sprintf(",this.%s", c.modelName)
	//		}
	//
	//		for _, pk := range pks {
	//			updateCode += fmt.Sprintf(",this.%s", pk.modelName)
	//		}
	//
	//		updateCode += ").Exec()\n"
	//		updateCode += `    if err == nil {
	//        num, _ = res.RowsAffected()
	//    }
	//
	//	return
	//}`
	//	}

	//source += updateCode

	//常量代码
	constCode := "//table " + t.name + " name and attributes defination.\n"
	constCode += "const TABLE_" + t.modelName + " = \"" + t.name + "\"\n"
	for _, c := range t.columns {
		constCode += "const COLUMN_" + t.modelName + "_" + c.modelName + " = \"" + c.name + "\"\n"
	}
	for _, c := range t.columns {
		constCode += "const ATTRIBUTE_" + t.modelName + "_" + c.modelName + " = \"" + c.modelName + "\"\n"
	}
	constCode += "\n"
	source += constCode

	//生成消息结构体
	msgCode := "\n//消息定义\ntype WSMsg" + t.modelName + " struct{\n"
	isNeedOutputMsg := false
	for _, c := range t.columns {
		if c.isProto {
			msgCode += fmt.Sprintf("\t%s %s `json:\"%s,omitempty\"`\n", c.modelName, DATA2MODEL_TYPE_MAP[c.sType], c.name)
			isNeedOutputMsg = true
		}
	}
	msgCode += "}\n"

	//生成数据库结构转换为消息结构的函数
	msgCode += fmt.Sprintf("\n//convert db struct to message struct\nfunc %s2WSMsg%s(%s *%s) (msg *WSMsg%s) {\n", t.modelName, t.modelName, t.name, t.modelName,
		t.modelName)
	msgCode += fmt.Sprintf("\tmsg = &WSMsg%s{\n", t.modelName)

	for _, c := range t.columns {
		if c.isDb && c.isProto {
			msgCode += fmt.Sprintf("\t\t%s : %s.%s,\n", c.modelName, t.name, c.modelName)
			isNeedOutputMsg = true
		}

	}

	msgCode += "\t}\n\treturn\n}\n"

	source += "//auto_models_end\n"

	if isNeedOutputMsg {
		//source += msgCode
	}

	err = ioutil.WriteFile(MODEL_DIR+t.name+".go", []byte(source), os.FileMode(0666))
	if err != nil {
		common.LogFuncCritical("write to %s failed : %v", MODEL_DIR+t.name+".go", err)
		return
	}

	return
}

func isPk(c *column, pks []*column) bool {
	for _, pk := range pks {
		if pk.name == c.name {
			return true
		}
	}
	return false
}

func outputSql(createSql, alterSql, dropSql, addSql, deleteSql, dropTableSql string) (err error) {
	file := []string{
		"__all_table_create.sql",
		"__all_table_field_alter.sql",
		"__all_table_field_drop.sql",
		"__all_table_field_add.sql",
		"__all_table_field_delete.sql",
		"__all_table_drop.sql",
	}
	sql := []string{
		createSql, alterSql, dropSql, addSql, deleteSql, dropTableSql,
	}
	for i := 0; i < len(file); i++ {
		err = ioutil.WriteFile(SQL_DIR+file[i], []byte(sql[i]), os.FileMode(0666))
		if err != nil {
			common.LogFuncCritical("write to %s failed : %v", SQL_DIR+file[i], err)
			return
		}
	}

	return
}

func (t *table) generateSQL() (err error, createSql, alterSql, dropSql string) {
	//create sql
	err, createSql = t.generateCreateSQL()
	if err != nil {
		return
	}

	//alter sql
	err, alterSql = t.generateAlterSQL()
	if err != nil {
		return
	}

	//drop sql
	err, dropSql = t.generateDropSQL()
	if err != nil {
		return
	}

	return
}

func (t *table) generateCreateSQL() (err error, sql string) {
	sql += "-- --------------------------------------------------\n"
	sql += fmt.Sprintf("--  Table Structure for `models.%s`\n", t.modelName)
	sql += "-- --------------------------------------------------\n"
	sql += fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (\n", t.name)
	var primaryKeyStr string
	for i, c := range t.columns {
		if !c.isDb {
			continue
		}

		if c.sType == "varchar" || c.sType == "char" {
			sql += fmt.Sprintf("`%s` %s(%d)", c.name, c.sType, c.iSize)
		} else {
			sql += fmt.Sprintf("`%s` %s", c.name, c.sType)
		}

		if c.isNull {
			sql += " NULL"
		} else {
			sql += " NOT NULL"
		}

		if c.isPk {
			if len(primaryKeyStr) == 0 {
				primaryKeyStr = fmt.Sprintf("`%s`", c.name)
			} else {
				primaryKeyStr += fmt.Sprintf(",`%s`", c.name)
			}
		}

		if c.isAuto {
			sql += " AUTO_INCREMENT"
		}

		if c.iDefault != nil {
			if c.sType == "varchar" || c.sType == "char" {
				sql += " DEFAULT ''"
			} else {
				sql += fmt.Sprintf(" DEFAULT %v", c.iDefault)
			}

		}

		if i == len(t.columns)-1 && len(primaryKeyStr) == 0 {
			sql += fmt.Sprintf(" COMMENT '%s'\n", c.comment)
		} else {
			sql += fmt.Sprintf(" COMMENT '%s',\n", c.comment)
		}

	}

	if len(primaryKeyStr) > 0 {
		sql += fmt.Sprintf("PRIMARY KEY(%s)\n", primaryKeyStr)
	}

	sql += fmt.Sprintf(") ENGINE=InnoDB COMMENT='%s' DEFAULT CHARSET=utf8;\n", t.comment)
	if len(t.indexs) > 0 {
		for _, index := range t.indexs {
			if len(index.columns) > 0 {
				indexName := t.name
				var indexContent string
				var indexLen int
				for i, c := range index.columns {
					var found bool
					for _, tc := range t.columns {
						if tc.name == c {
							found = true
							if tc.indexLen > 0 {
								indexLen = tc.indexLen
							}
						}
					}

					if !found { //index name no found
						common.LogFuncCritical("read file %s failed : %v", CONFIG_PATH, err)
						os.Exit(-1)
					}

					indexName += "_" + c
					if i == 0 {
						indexContent += fmt.Sprintf("`%s`", c)
					} else {
						indexContent += fmt.Sprintf(", `%s`", c)
					}
				}

				sql += fmt.Sprintf("CREATE INDEX `%s` ON `%s` (%s", indexName, t.name, indexContent)

				if indexLen > 0 {
					sql += fmt.Sprintf("(%d));\n", indexLen)
				} else {
					sql += ");\n"
				}
			}
		}
	}

	if len(t.uniqueIndexs) > 0 {
		for _, index := range t.uniqueIndexs {
			if len(index.columns) > 0 {
				indexName := t.name
				var indexContent string
				for i, c := range index.columns {
					indexName += "_" + c
					if i == 0 {
						indexContent += fmt.Sprintf("`%s`", c)
					} else {
						indexContent += fmt.Sprintf(", `%s`", c)
					}
				}

				sql += fmt.Sprintf("CREATE UNIQUE INDEX `%s` ON `%s` (%s);\n", indexName, t.name, indexContent)
			}
		}
	}

	sql += "\n"

	return
}

func (t *table) generateAlterSQL() (err error, sql string) {
	sql += "----------------------------------------------------\n"
	sql += fmt.Sprintf("--  `%s`\n", t.name)
	sql += "----------------------------------------------------\n"

	for _, c := range t.columns {
		if !c.isDb {
			continue
		}

		sql += fmt.Sprintf("ALTER TABLE `%s` CHANGE `%s` ", t.name, c.name)
		if c.sType == "varchar" || c.sType == "char" {
			sql += fmt.Sprintf("`%s` %s(%d)", c.name, c.sType, c.iSize)
		} else {
			sql += fmt.Sprintf("`%s` %s", c.name, c.sType)
		}

		if c.isNull {
			sql += " NULL"
		} else {
			sql += " NOT NULL"
		}

		if c.iDefault != nil {
			if c.sType == "varchar" || c.sType == "char" {
				sql += " DEFAULT ''"
			} else {
				sql += fmt.Sprintf(" DEFAULT %v", c.iDefault)
			}
		}

		if c.isAuto {
			sql += " AUTO_INCREMENT"
		}

		sql += fmt.Sprintf(" COMMENT '%s';\n", c.comment)
	}
	sql += "\n"

	return
}

func (t *table) generateDropSQL() (err error, sql string) {
	sql += "----------------------------------------------------\n"
	sql += fmt.Sprintf("--  `%s`\n", t.name)
	sql += "----------------------------------------------------\n"

	for _, c := range t.columns {
		if !c.isDb {
			continue
		}

		sql += fmt.Sprintf("ALTER TABLE `%s` DROP `%s`;\n", t.name, c.name)
	}

	sql += "\n"

	return
}

func (t *table) generateAddSQL() (err error, sql string) {

	var preCName string
	sql += "----------------------------------------------------\n"
	sql += fmt.Sprintf("--  `%s`\n", t.name)
	sql += "----------------------------------------------------\n"

	for _, c := range t.columns {
		if !c.isDb {
			continue
		}

		sql += fmt.Sprintf("ALTER TABLE `%s` ADD ", t.name)
		if c.sType == "varchar" || c.sType == "char" {
			sql += fmt.Sprintf("`%s` %s(%d)", c.name, c.sType, c.iSize)
		} else {
			sql += fmt.Sprintf("`%s` %s", c.name, c.sType)
		}

		if c.isNull {
			sql += " NULL"
		} else {
			sql += " NOT NULL"
		}

		if c.iDefault != nil {
			if c.sType == "varchar" || c.sType == "char" {
				sql += " DEFAULT ''"
			} else {
				sql += fmt.Sprintf(" DEFAULT %v", c.iDefault)
			}

		}

		if c.isAuto {
			sql += " AUTO_INCREMENT"
		}

		if len(preCName) > 0 {
			sql += fmt.Sprintf(" COMMENT '%s' AFTER `%s`;\n", c.comment, preCName)
		} else {
			sql += fmt.Sprintf(" COMMENT '%s';\n", c.comment)
		}

		preCName = c.name
	}
	sql += "\n"

	return
}

func (t *table) generateDeleteSQL() (err error, sql string) {
	sql += fmt.Sprintf("DELETE FROM `%s`;\n", t.name)
	return
}

func (t *table) generateDropTableSQL() (err error, sql string) {
	sql += fmt.Sprintf("DROP TABLE `%s`;\n", t.name)
	return
}

var ERR_PARSE_FAILED = errors.New("parse failed!")

func main() {

	flag.StringVar(&MODEL_DIR, "models-dir", "../../models", "directory of models")
	flag.StringVar(&SQL_DIR, "sql-dir", "../../sql", "directory of sql")
	flag.StringVar(&CONFIG_PATH, "config", "./", "directory of config file")
	flag.Parse()
	//common.LogFuncFlags()
	//common.LogFuncPrintf("MODEL_DIR %s, SQL_DIR %s\n", MODEL_DIR, SQL_DIR)

	var content map[string]interface{}
	if data, err := ioutil.ReadFile(CONFIG_PATH); err != nil {
		common.LogFuncCritical("read file %s failed : %v", CONFIG_PATH, err)
		return
	} else {
		if err = json.Unmarshal(data, &content); err != nil {
			common.LogFuncCritical("json.Unmarshal failed : %v", err)
			return
		} else {
			var tableSlice TableSlice
			for k, v := range content {
				//common.LogFuncPrintf("%v : %v\n", k, v)
				maps, ok := v.(map[string]interface{})
				if !ok {
					common.LogFuncCritical("convert failed\n")
					return
				}

				var (
					t = NewTable(k)
				)

				for k1, v1 := range maps {
					switch k1 {
					case MODEL_NAME:
						{
							if t.modelName, ok = v1.(string); !ok {
								common.LogFuncCritical("parse COMMENT %v 失败", v)
								return
							}
						}
					case COMMENT:
						{
							if t.comment, ok = v1.(string); !ok {
								common.LogFuncCritical("parse COMMENT %v 失败", v)
								return
							}
						}

					case INDEX:
						{
							err = t.parseIndex(v1.([]interface{}))
							if err != nil {
								return
							}
						}

					case UNIQUE_INDEX:
						{
							err = t.parseUniqueIndex(v1.([]interface{}))
							if err != nil {
								return
							}
						}

					default:
						m, ok := v1.(map[string]interface{})
						if !ok {
							common.LogFuncCritical("convert %v failed", v1)
							return
						}
						c := NewColumn()
						if err = c.parse(m); err != nil {
							common.LogFuncCritical("convert failed : %v\n", err)
							return
						} else {
							c.name = k1
							t.columns = append(t.columns, c)
						}
					}
				}

				if t != nil {
					sort.Slice(t.columns, func(i, j int) bool {
						return t.columns[i].sn < t.columns[j].sn
					})
					tableSlice = append(tableSlice, t)
				}
			}

			//sort tableSlice
			sort.Slice(tableSlice, func(i, j int) bool {
				return tableSlice[i].name < tableSlice[j].name

			})

			//导出common代码
			err = outputModelCommonCode(tableSlice)
			if err != nil {
				return
			}

			//导出model代码源文件
			for _, t := range tableSlice {
				t.outputModelCode()
			}

			//导出sql
			var (
				createSql, alterSql, dropSql, addSql, deleteSql, dropTableSql string
			)

			for _, t := range tableSlice {
				var sql string
				err, sql = t.generateCreateSQL()
				if err != nil {
					return
				}
				createSql += sql

				sql = ""
				err, sql = t.generateAlterSQL()
				if err != nil {
					return
				}
				alterSql += sql

				sql = ""
				err, sql = t.generateDropSQL()
				if err != nil {
					return
				}
				dropSql += sql

				sql = ""
				err, sql = t.generateAddSQL()
				if err != nil {
					return
				}
				addSql += sql

				sql = ""
				err, sql = t.generateDeleteSQL()
				if err != nil {
					return
				}
				deleteSql += sql

				sql = ""
				err, sql = t.generateDropTableSQL()
				if err != nil {
					return
				}
				dropTableSql += sql
			}

			outputSql(createSql, alterSql, dropSql, addSql, deleteSql, dropTableSql)
		}
	}
}

//生成model文件
func outputModelCommonCode(tables TableSlice) (err error) {
	code := "package models\n\n" +
		"import \"github.com/astaxie/beego/orm\"\n\n"

	code += "func ModelsInit() (err error) {\n" +
		"	orm.RegisterModel(\n"
	for _, t := range tables {
		code += "\t\tnew(" + t.modelName + "),\n"
	}
	code += "\t)\n" +
		"\n\treturn\n" +
		"}\n"

	err = ioutil.WriteFile(MODEL_DIR+"common.go", []byte(code), os.FileMode(0666))
	if err != nil {
		common.LogFuncCritical("write to %s failed : %v", MODEL_DIR+"common.go", err)
		return
	}

	return
}
