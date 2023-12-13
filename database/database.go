package database

import (
	"fmt"
	"log"

	"github.com/udonetsm/server/use"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Unumber() {
	y := use.ReadYamlFile("database/cfg.yaml")
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", y.User, y.Pass, y.Host, y.Port, y.DBNM, y.SSLM)
	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
}
