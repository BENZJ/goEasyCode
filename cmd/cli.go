package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/benzj/goEasyCode"
)

func main() {
	parser()
}

func parser() {
	dsn := flag.String("dsn", "", "数据库dsn配置")
	table := flag.String("table", "", "表名称")
	template := flag.String("template", "Do", "生成的类型默认是Do类,可选Do,Mapper,Dao")
	version := flag.Bool("version", false, "版本号")
	v := flag.Bool("v", false, "版本号")
	enableJsonTag := flag.Bool("enableJsonTag", false, "是否添加json的tag,默认false")
	h := flag.Bool("h", false, "帮助")
	help := flag.Bool("help", false, "帮助")

	// 开始
	flag.Parse()

	if *h || *help {
		flag.Usage()
		return
	}

	// 版本号
	if *version || *v {
		fmt.Println(fmt.Sprintf("\n version: %s\n %s\n using -h param for more help \n",
			goEasyCode.VERSION, goEasyCode.VERSION_TEXT))
		return
	}

	dbconf := false
	var db *goEasyCode.DbConfig
	if goEasyCode.Exists("mysql.yaml") {
		dbconf = true
		db = goEasyCode.NewDbConfig().GetYamlConf()
	}

	// 初始化
	t2t := goEasyCode.NewTable2Struct()
	t2t.Table(*table).
		// 是否添加json tag
		EnableJsonTag(*enableJsonTag).
		// 数据库dsn
		Dsn(*dsn)

	var res string
	var fileName string
	var err error
	if dbconf {
		res, fileName, err = goEasyCode.Generate(db, t2t, *template)
	} else {
		res, fileName, err = goEasyCode.Generate(nil, t2t, *template)
	}

	if err != nil {
		log.Println(err.Error())
	}
	fmt.Printf("文件名称:  %s", fileName)
	print(res)
}
