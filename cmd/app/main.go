package main

import (
	"fmt"
	"os"

	"challenge/cmd/commands"
)

func main() {
	config := &commands.Config{
		StateFile:  "current_state.json",
		CitiesFile: "cities.csv",
		LogLevel:   "info",
	}

	logger := commands.InitLogger(config.LogLevel)
	app := commands.NewApp(logger, config)
	rootCmd := app.BuildRootCommand()
	if err := rootCmd.Execute(); err != nil {
		logger.Error("command execution failed", "error", err.Error())
		fmt.Println(err)
		os.Exit(1)
	}
}
