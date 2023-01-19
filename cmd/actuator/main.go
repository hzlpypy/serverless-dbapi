package main

import (
	"fmt"
	"os"
	"serverless-db/pkg/actuator"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
)

func main() {
	file, err := os.ReadFile("./actuator.yaml")
	if err != nil {
		panic(err)
	}
	config := actuator.Config{}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}
	fmt.Println(config)
	actuator, err := actuator.New(config)
	if err != nil {
		panic(err)
	}
	err = actuator.Run()
	if err != nil {
		panic(err)
	}
}
