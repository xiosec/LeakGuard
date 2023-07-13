package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Elastic struct {
	Addresses []string `json:"addresses"`
	Index     string   `json:"index"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
}

type Service struct {
	Host  string `json:"host"`
	Port  uint   `json:"port"`
	Token string `json:"token"`
}

type Config struct {
	Elastic Elastic `json:"elastic"`
	Service Service `json:"service"`
}

var Conf Config

func Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &Conf)

	return err
}
