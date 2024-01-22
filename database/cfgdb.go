package database

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/udonetsm/cmngb/flags"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type YAMLObject struct {
	Host  string `yaml:"host"`
	Port  string `yaml:"port"`
	User  string `yaml:"user"`
	Pass  string `yaml:"password"`
	SSLM  string `yaml:"sslmode"`
	DBNM  string `yaml:"dbname"`
	Error error  `yaml:"-"`
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
		y.Error = err
		return
	}
	err = yaml.Unmarshal(data, y)
	if err != nil {
		y.Error = err
		return
	}
	y.User = flags.NAME
	y.Pass = flags.PASS
}

func authFailed(err error) bool {
	return strings.Contains(err.Error(), "28P01")
}

// build database connection string
// using object built on YAMLCfg function
// and get database usin built config
func (y *YAMLObject) GetDB() (db *gorm.DB) {
	var err error
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", y.User, y.Pass, y.Host, y.Port, y.DBNM, y.SSLM)
	dialector := postgres.Open(dsn)
	db, err = gorm.Open(dialector)
	if err != nil {
		if authFailed(err) {
			y.Error = errors.New("AUTH FAILED")
		} else {
			y.Error = err
		}
	}
	return
}

// using duck typing for load database connection
func LoadCfgAndGetDB(yg CfgDBGetter, path string) (db *gorm.DB) {
	yg.YAMLCfg(path)
	db = yg.GetDB()
	return
}
