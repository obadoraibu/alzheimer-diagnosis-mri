package main

import (
	"flag"
	"github.com/obadoraibu/go-auth/internal/app"
	"log"
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

	if err := app.Run(mainConfigPath, dbConfigPath); err != nil {
		log.Fatal("cannot run the app")
	}
}
