package main

import (
	"github.com/lin-snow/ech0/cmd"
	"github.com/lin-snow/ech0/internal/config"
	logUtil "github.com/lin-snow/ech0/internal/util/log"
)

func init() {
	// Logger
	logUtil.InitLogger()

	// Config
	config.LoadAppConfig()
}

func main() {
	cmd.Execute()
}
