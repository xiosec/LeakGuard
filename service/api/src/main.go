package main

import (
	"LeakGuard/config"
	"LeakGuard/databases"
	"LeakGuard/service"
	"flag"
	"fmt"
	"log"
	"time"
)

var (
	ConfigPath string
)

func main() {
	flag.StringVar(&ConfigPath, "config", "config.json", "config file path")
	flag.Parse()

	fmt.Println("\tLeakGuard")
	fmt.Println("https://github.com/xiosec/LeakGuard\n")

	err := config.Load(ConfigPath)
	if err != nil {
		log.Fatal("Error loading configs!")
	}

	err = databases.Init(config.Conf.Elastic)
	if err != nil {
		log.Fatal("Error connecting to elasticsearch!")
	}

	exist, err := databases.IndexExists(config.Conf.Elastic.Index)
	if !exist || err != nil {
		log.Fatal("The desired index was not found!\n" + err.Error())
	}

	fmt.Println("[*] Setting up the web service")
	fmt.Printf("[*] Address :%s:%d\n", config.Conf.Service.Host, config.Conf.Service.Port)
	fmt.Printf("[*] Token   :%s\n", config.Conf.Service.Token)
	fmt.Printf("[*] Time    :%s\n\n", time.Now().Format(time.RFC850))
	service.Run(config.Conf.Service)
}
