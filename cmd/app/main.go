package main

import (
	"log"

	"github.com/ziyadovea/task_manager/users/internal/app"
	"github.com/ziyadovea/task_manager/users/internal/config"
)

func main() {
	// read CLI params
	params, err := readCLIParams()
	if err != nil {
		log.Fatal(err)
	}

	// init config
	cfg, err := config.InitConfig(params.configPath)
	if err != nil {
		log.Fatal(err)
	}

	// run, Forrest, run!
	app.Run(cfg)
}
