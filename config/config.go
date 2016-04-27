package config

import (
	"os"

	"github.com/jinzhu/configor"
	"github.com/qor/render"
)

var Config = struct {
	Port uint `default:"7000" env:"PORT"`
	DB   struct {
		Name     string `default:"qor_example"`
		Adapter  string `default:"mysql"`
		User     string
		Password string
	}
}{}

var (
	Root = os.Getenv("GOPATH") + "/src/github.com/qor/qor-example"
	View *render.Render
)

func init() {
	if err := configor.Load(&Config, "config/database.yml"); err != nil {
		panic(err)
	}

	View = render.New()
}
