package use

import (
	"log"
	"os"

	"github.com/udonetsm/client/models"
	"gopkg.in/yaml.v2"
)

func ReadYamlFile(path string) *models.YAMLObject {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	y := &models.YAMLObject{}
	err = yaml.Unmarshal(data, y)
	return y
}
