package database

import (
	"fmt"

	"github.com/udonetsm/client/models"
	"gorm.io/gorm"
)

func LoadDb() *gorm.DB {
	y := &models.YAMLObject{}
	db := models.LoadCfgAndGetDB(y, "database/cfg.yaml")
	return db
}

func UpdateNumber(rj *models.RequestJSON) {
	db := LoadDb()
	db.Model(&models.JSONObject{}).Where("number=?", rj.Target).Update("entries_id", rj.Upgrade)
}

func Create(c *models.Contact) {
	o := &models.JSONObject{}
	o.Pack(c)
	// db := LoadDb()
	fmt.Println(o)
}
