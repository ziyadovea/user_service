package main

import (
	"errors"
	"flag"
)

type cliParams struct {
	configPath string
}

func readCLIParams() (cliParams, error) {
	params := cliParams{}

	flag.StringVar(&params.configPath, "config", "", "Path to config file.")
	flag.Parse()

	if params.configPath == "" {
		return cliParams{}, errors.New("empty config path")
	}

	return params, nil
}
