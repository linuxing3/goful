package config

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/anmitsu/goful/app"
	"github.com/spf13/viper"
)

func CustomizeCommand(g *app.Goful) {
	CustomizeConfig(g, "external-command")
}

func CustomizeBookmark(g *app.Goful) {
	CustomizeConfig(g, "bookmark")
}

func CustomizeConfig(g *app.Goful, section string) {
	path := fmt.Sprintf("menu.%s.list", section)
	list := viper.GetStringSlice(path)
	for _, item := range list {
		if runtime.GOOS == "windows" {
			AddMenuItem(g, section, item)
		} else {
			if !strings.Contains(item, "win") {
				AddMenuItem(g, section, item)
			}
		}
	}
}
