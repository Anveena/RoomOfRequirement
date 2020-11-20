package ezConfig

import (
	"errors"
	"flag"
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

func ReadConf(configModel interface{}) error {
	confPath := ""
	flag.StringVar(&confPath, "c", "", "Config toml File Path")
	flag.Parse()
	if confPath == "" {
		return errors.New("config path is null")
	}
	tomlData, err := ioutil.ReadFile(confPath)
	if err != nil {
		return err
	}
	return toml.Unmarshal(tomlData, configModel)
}
