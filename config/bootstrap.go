package config

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"time"
)

//go:embed *.yaml
var configs embed.FS

func init() {
	env := os.Getenv("ENV")
	vp := viper.New()

	configFileStream, err := configs.ReadFile("app." + env + ".yaml")
	fmt.Println(configFileStream)
	if err != nil {
		panic(err)
	}
	vp.SetConfigType("yaml")
	err = vp.ReadConfig(bytes.NewReader(configFileStream))
	if err != nil {
		panic(err)
	}
	err = vp.UnmarshalKey("app", &App)
	if err != nil {
		panic(err)
	}

	err = vp.UnmarshalKey("database", &Database)
	if err != nil {
		panic(err)
	}
	Database.MaxLifeTime *= time.Second
}
