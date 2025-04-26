package main

import (
	"flag"
	"log"

	"github.com/obadoraibu/go-auth/internal/app"
	"github.com/sirupsen/logrus"
)

var (
	mainConfigPath string
	dbConfigPath   string
)

func init() {
	flag.StringVar(&mainConfigPath, "mainCfgPath", "configs/main.yml", "path to the main config file")
	flag.StringVar(&dbConfigPath, "dbCfgPath", "configs/db.yml", "path to the database config file")
}

func main() {
	flag.Parse()
	logrus.Println("main.go run")
	if err := app.Run(mainConfigPath, dbConfigPath); err != nil {
		log.Fatal("cannot run the app")
	}
}
