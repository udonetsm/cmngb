package cfgsrv

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// this object for load database connection config
type YAMLObject struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"password"`
	SSLM string `yaml:"sslmode"`
	DBNM string `yaml:"dbname"`
}
type Ntrs struct {
	Number, Object string
}

// duck typing for load data base connection config
type CfgDBGetter interface {
	YAMLCfg(string)
	GetDB() *gorm.DB
}

// this method read config from target .yaml file and unpack it in object
func (y *YAMLObject) YAMLCfg(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(data, y)
}

// build database connection string
// using object built on YAMLCfg function
// and get database usin built config
func (y *YAMLObject) GetDB() (db *gorm.DB) {
	var err error
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", y.User, y.Pass, y.Host, y.Port, y.DBNM, y.SSLM)
	db, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Println(err)
		return
	}
	return
}

// using duck typing for load database connection
func LoadCfgAndGetDB(yg CfgDBGetter, path string) (db *gorm.DB) {
	yg.YAMLCfg(path)
	db = yg.GetDB()
	return
}
