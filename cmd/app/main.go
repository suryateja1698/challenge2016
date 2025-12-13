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

	rootCmd.PersistentFlags().StringVarP(&config.StateFile, "state", "s", "current_state.json", "State file to persist data")
	rootCmd.PersistentFlags().StringVarP(&config.CitiesFile, "cities", "c", "cities.csv", "Cities CSV file")
	rootCmd.PersistentFlags().StringVarP(&config.LogLevel, "log-level", "l", "info", "Log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().StringVarP(&config.LogFormat, "log-format", "f", "json", "Log format JSON")

	if err := rootCmd.Execute(); err != nil {
		logger.Error("command execution failed", "error", err.Error())
		fmt.Println(err)
		os.Exit(1)
	}
}
