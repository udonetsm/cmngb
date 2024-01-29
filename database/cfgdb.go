package database

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/udonetsm/cmngb/flags"
	"github.com/udonetsm/cmngb/models"
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

func CreateTableForUser(db *gorm.DB, e *models.Entries, y *YAMLObject) {
	s := fmt.Sprintf("create user %s with password '%s'; create table %s_entries"+
		"(id varchar(20) primary key, contact jsonb); "+
		"alter table %s_entries owner to %s", e.Owner, y.Pass, e.Owner, e.Owner, e.Owner)
	e.Error = db.Exec(s).Error
}

// duck typing for load data base connection config
type CfgDBGetter interface {
	YAMLCfg(string)
	GetDB(*models.Entries) *gorm.DB
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
}

func authFailed(err error) bool {
	return strings.Contains(err.Error(), "28P01")
}

// build database connection string
// using object built on YAMLCfg function
// and get database usin built config
func (y *YAMLObject) GetDB(e *models.Entries) (db *gorm.DB) {
	var err error
	var dsn string
	if e.Ok {
		y.User = flags.DBUSER
		y.Pass = flags.DBPASS
		dsn = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", y.User, y.Pass, y.Host, y.Port, y.DBNM, y.SSLM)
	} else {
		dsn = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", e.Owner, y.Pass, y.Host, y.Port, y.DBNM, y.SSLM)
	}
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
func LoadCfgAndGetDB(yg CfgDBGetter, path string, e *models.Entries) (db *gorm.DB) {
	yg.YAMLCfg(path)
	db = yg.GetDB(e)
	return
}
