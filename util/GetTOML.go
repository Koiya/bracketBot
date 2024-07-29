package util

import (
	"fmt"
	"github.com/pelletier/go-toml"
)

func GetTOML(x string) string {
	config, err := toml.LoadFile("config.toml")
	if err != nil {
		fmt.Println("Error ", err.Error())
		return ""
	}
	result := config.Get(x).(string)
	return result
}
