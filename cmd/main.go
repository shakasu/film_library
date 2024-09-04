package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type config struct {
	Port string `yaml:"port"`
	Db   struct {
		Username string `yaml:"username"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Dbname   string `yaml:"dbname"`
	}
}

func main() {

	var c config
	c.getConfig()

	fmt.Println(c)
	//srv := new(Server)
	//go func() {
	//	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
	//		fmt.Printf("error occured while running http server: %s", err.Error())
	//	}
	//}()
}

func (cfg *config) getConfig() *config {

	yamlFile, err := os.ReadFile("config/config.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return cfg
}
