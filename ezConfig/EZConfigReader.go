package ezConfig

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
)

func ReadConf(configModel interface{}) error {
	confPath := ""
	flag.StringVar(&confPath, "c", "", "Config JSON File Path")
	flag.Parse()
	if confPath == "" {
		return errors.New("config path is null")
	}
	jsonData, err := ioutil.ReadFile(confPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, configModel)
}
