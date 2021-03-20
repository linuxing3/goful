package config

import (
	"fmt"

	"github.com/anmitsu/goful/app"
	"github.com/anmitsu/goful/menu"

	"github.com/spf13/viper"
)

type MenuItem struct {
	accel string
	label string
	path  string
}

func SetConfigItem(g *app.Goful, section string, item string, value interface{}) {
	viper.Set(fmt.Sprintf("menu.%s.%s", section, item), value)
	viper.WriteConfig()
}

func GetConfigItem(g *app.Goful, section string, item string) map[string]interface{} {
	value := viper.GetStringMap(fmt.Sprintf("menu.%s.%s", section, item))
	return value
}

func AddMenuItem(g *app.Goful, section string, item string) {
	// TODO: 将结构转变为结构体
	// c := MenuItem{}
	// data := GetConfigItem(g, section, item)
	// FillStruct(data, c)
	// mapstruct.Map2Struct(data, &c)
	accel := viper.GetString(fmt.Sprintf("menu.%s.%s.accel", section, item))
	label := viper.GetString(fmt.Sprintf("menu.%s.%s.label", section, item))
	path := viper.GetString(fmt.Sprintf("menu.%s.%s.path", section, item))
	c := MenuItem{
		accel: accel,
		label: label,
		path: path,
	}
	if section == "bookmark" {
		menu.Add(section, c.accel, c.label, func() { g.Dir().Chdir(c.path) })
	} else if section == "command" {
		menu.Add(section, c.accel, c.label, func() { g.Shell(c.path) })
	}
}
