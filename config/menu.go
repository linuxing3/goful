package config

import (
	"fmt"
	"reflect"

	"github.com/anmitsu/goful/app"
	"github.com/anmitsu/goful/menu"
	"github.com/linuxing3/goful/util"

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

func GetConfigItem(g *app.Goful, section string, item string) map[string]string {
	return viper.GetStringMapString(fmt.Sprintf("menu.%s.%s", section, item))
}

func MapConfigItemToMenuItem(section, item string) MenuItem{
	m := MenuItem{}
	m.accel = viper.GetString(fmt.Sprintf("menu.%s.%s.accel", section, item))
	m.label = viper.GetString(fmt.Sprintf("menu.%s.%s.label", section, item))
	m.path = viper.GetString(fmt.Sprintf("menu.%s.%s.path", section, item))
	return m
}

func AddMenuItem(g *app.Goful, section string, item string) {

	c := MapConfigItemToMenuItem(section, item)
	// debug stuct
	// DebugMenuItem(c)
	if section == "bookmark" {
		menu.Add(section, c.accel, c.label, func() { g.Dir().Chdir(c.path) })
	} else if section == "external-command" || section == "command" {
		menu.Add(section, c.accel, c.label, func() { g.Shell(c.path) })
	} else if section == "editor" {
		menu.Add(section, c.accel, c.label, func() { g.Spawn(c.path) })
	} else if section == "git" {
		menu.Add(section, c.accel, c.label, func() { g.Spawn(c.path) })
	}
}

func DebugMenuItem(c interface{}) {
	// TODO: 将结构转变为结构体
	t := reflect.ValueOf(c).Elem()
	util.ValueToString(t)
	// for k, v := range data {
	// 	t.FieldByName(k).Set(reflect.ValueOf(v))
	// }
}
