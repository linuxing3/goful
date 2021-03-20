package config

import (
	"fmt"

	"github.com/anmitsu/goful/app"
	"github.com/anmitsu/goful/menu"
	"github.com/spf13/viper"
)

func AddMenuItem(g *app.Goful, key string, item string) {
	accel := viper.GetString(fmt.Sprintf("menu.%s.%s.accel", key, item))
	label := viper.GetString(fmt.Sprintf("menu.%s.%s.label", key, item))
	path := viper.GetString(fmt.Sprintf("menu.%s.%s.path", key, item))
	menu.Add(key, accel, label, func() { g.Shell(path)})
}