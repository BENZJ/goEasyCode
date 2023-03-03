package goEasyCode

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type DbConfig struct {
	host     string
	port     string
	user     string
	password string
	dbName   string
}

type YamlDbConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DbName   string `yaml:"dbName"`
}

func (yaml *YamlDbConfig) GetDbConfig() *DbConfig {
	return NewDbConfig().Host(yaml.Host).Port(yaml.Port).User(yaml.User).Password(yaml.Password).DbName(yaml.DbName)
}

func (db *DbConfig) Host(s string) *DbConfig {
	db.host = s
	return db
}

func (db *DbConfig) Port(s string) *DbConfig {
	db.port = s
	return db
}

func (db *DbConfig) User(s string) *DbConfig {
	db.user = s
	return db
}

func (db *DbConfig) Password(s string) *DbConfig {
	db.password = s
	return db
}

func (db *DbConfig) DbName(s string) *DbConfig {
	db.dbName = s
	return db
}

//DbConfig 生成 mysql url
func (db *DbConfig) Url() string {
	return db.user + ":" + db.password + "@tcp(" + db.host + ":" + db.port + ")/" + db.dbName + "?charset=utf8"
}

//新建一个DbConfig
func NewDbConfig() *DbConfig {
	return &DbConfig{}
}

func (db *DbConfig) GetDb(dsn string) (res *sql.DB, err error) {
	if dsn != "" {
		res, err = sql.Open("mysql", dsn)
	} else {
		res, err = sql.Open("mysql", db.Url())
	}

	return
}

func (c *DbConfig) GetYamlConf() *DbConfig {
	//读取resources/application.yaml文件
	yamlFile, err := ioutil.ReadFile("mysql.yaml")
	//若出现错误，打印错误提示
	if err != nil {
		fmt.Println(err.Error())
	}
	conf := &YamlDbConfig{}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		fmt.Println(err.Error())
	}
	return conf.GetDbConfig()
}

/**
*判断文件是否存在
 */
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return true
	}
	return true
}
