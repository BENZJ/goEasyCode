package goEasyCode

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

var typeForMysqlToJava = map[string]string{
	"int":                "Integer",
	"integer":            "Integer",
	"tinyint":            "Integer",
	"smallint":           "Integer",
	"mediumint":          "Integer",
	"bigint":             "Integer",
	"int unsigned":       "Integer",
	"integer unsigned":   "Integer",
	"tinyint unsigned":   "Integer",
	"smallint unsigned":  "Integer",
	"mediumint unsigned": "Integer",
	"bigint unsigned":    "Integer",
	"bit":                "Integer",
	"bool":               "Boolean",
	"enum":               "String",
	"set":                "String",
	"varchar":            "String",
	"char":               "String",
	"tinytext":           "String",
	"mediumtext":         "String",
	"text":               "String",
	"longtext":           "String",
	"blob":               "String",
	"tinyblob":           "String",
	"mediumblob":         "String",
	"longblob":           "String",
	"date":               "Date", // time.Time or string
	"datetime":           "Date", // time.Time or string
	"timestamp":          "Date", // time.Time or string
	"time":               "Date", // time.Time or string
	"float":              "Double",
	"double":             "Double",
	"decimal":            "Double",
	"binary":             "String",
	"varbinary":          "String",
}

type Column struct {
	ColumnName     string //java变量名称
	RealColumnName string //数据库中列名称
	Type           string //java类型
	RealType       string //数据库类型
	Nullable       string
	TableName      string
	ColumnComment  string
	Tag            string
}

type table struct {
	TableName string
	ClassName string
	Columns   []Column
}

//map for converting mysql type to golang types

type Table2Struct struct {
	dsn            string
	db             *sql.DB
	table          string
	prefix         string
	err            error
	realNameMethod string
	enableJsonTag  bool   // 是否添加json的tag, 默认不添加
	packageName    string // 生成struct的包名(默认为空的话, 则取名为: package model)
	tagKey         string // tag字段的key值,默认是orm
	dateToTime     bool   // 是否将 date相关字段转换为 time.Time,默认否
	columns        []Column
}

func NewTable2Struct() *Table2Struct {
	return &Table2Struct{}
}

func (t *Table2Struct) Dsn(d string) *Table2Struct {
	t.dsn = d
	return t
}

func (t *Table2Struct) TagKey(r string) *Table2Struct {
	t.tagKey = r
	return t
}

func (t *Table2Struct) PackageName(r string) *Table2Struct {
	t.packageName = r
	return t
}

func (t *Table2Struct) RealNameMethod(r string) *Table2Struct {
	t.realNameMethod = r
	return t
}

func (t *Table2Struct) DB(d *sql.DB) *Table2Struct {
	t.db = d
	return t
}

func (t *Table2Struct) Table(tab string) *Table2Struct {
	t.table = tab
	return t
}

func (t *Table2Struct) GetTable() string {
	return t.table
}

func (t *Table2Struct) Prefix(p string) *Table2Struct {
	t.prefix = p
	return t
}

func (t *Table2Struct) EnableJsonTag(p bool) *Table2Struct {
	t.enableJsonTag = p
	return t
}

func (t *Table2Struct) DateToTime(d bool) *Table2Struct {
	t.dateToTime = d
	return t
}

func (t *Table2Struct) GetClaseName() string {
	return UpCaseFirst(CamelCase(t.table))
}

func (t *Table2Struct) InitColumns() { //初始化表结构
	var sqlStr = `SELECT COLUMN_NAME,DATA_TYPE,IS_NULLABLE,TABLE_NAME,COLUMN_COMMENT
		FROM information_schema.COLUMNS 
		WHERE table_schema = DATABASE()`
	if t.table != "" {
		sqlStr += fmt.Sprintf(" AND TABLE_NAME = '%s'", t.prefix+t.table)
	}
	// sql排序
	sqlStr += " order by TABLE_NAME asc, ORDINAL_POSITION asc"

	rows, err := t.db.Query(sqlStr)
	if err != nil {
		log.Println("Error reading table information: ", err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		col := Column{}
		err = rows.Scan(&col.RealColumnName, &col.RealType, &col.Nullable, &col.TableName, &col.ColumnComment)

		if err != nil {
			log.Println(err.Error())
			return
		}
		col.Type = typeForMysqlToJava[col.RealType]
		col.ColumnName = CamelCase(col.RealColumnName)
		t.columns = append(t.columns, col)
	}
	return
}

func (t *Table2Struct) GetColumns() []Column {
	return t.columns
}

func CamelCase(str string) string {
	// 是否有表前缀, 设置了就先去除表前缀
	var text string
	//for _, p := range strings.Split(name, "_") {
	for i, p := range strings.Split(str, "_") {
		if i == 0 {
			text += strings.ToLower(p[0:])
		} else {
			// 字段首字母大写的同时, 是否要把其他字母转换为小写
			switch len(p) {
			case 0:
			case 1:
				text += strings.ToUpper(p[0:1])
			default:
				text += strings.ToUpper(p[0:1]) + strings.ToLower(p[1:])
			}
		}
	}
	return text
}

/**
*首字母大写
 */
func UpCaseFirst(str string) (text string) {
	text = strings.ToUpper(str[0:1]) + str[1:]
	return
}

func Generate(dbconf *DbConfig, t2t *Table2Struct, tempType string) (res string, fileName string, err error) {
	var db *sql.DB
	if len(t2t.dsn) > 0 {
		db, err = dbconf.GetDb(t2t.dsn)
	} else {
		db, err = dbconf.GetDb("")
	}
	if err != nil {
		return
	}
	t2t.DB(db).Table(t2t.table).InitColumns()
	p := table{
		TableName: t2t.GetTable(),
		ClassName: t2t.GetClaseName(),
		Columns:   t2t.GetColumns(),
	}
	defer db.Close()

	var tmpl *template.Template
	if tempType == "Do" {
		tmpl, err = template.New("tmp").Parse(DoTempl)
		fileName = p.ClassName + "DO.java"
	} else if tempType == "Mapper" {
		tmpl, err = template.New("tmp").Parse(MapperTempl)
		fileName = p.ClassName + "Mapper.xml"
	} else if tempType == "Dao" {
		tmpl, err = template.New("tmp").Parse(DaoTempl)
		fileName = p.ClassName + "Mapper.java"
	}
	if err != nil {
		return
	}
	var b1 bytes.Buffer
	tmpl.Execute(&b1, p)
	res = b1.String()
	return
}
