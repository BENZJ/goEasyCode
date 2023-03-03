package test

import (
	"log"
	"testing"

	"github.com/benzj/goEasyCode"
)

func TestGenerate(t *testing.T) {
	var dbconf *goEasyCode.DbConfig
	if goEasyCode.Exists("mysql.yaml") {
		dbconf = goEasyCode.NewDbConfig().GetYamlConf()
	}
	t2t := goEasyCode.NewTable2Struct()
	t2t.Table("cookie_info")
	res, _, err := goEasyCode.Generate(dbconf, t2t, "Mapper")

	if err != nil {
		log.Println(err.Error())
		return
	}
	print(res)
}
