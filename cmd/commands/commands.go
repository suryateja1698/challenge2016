package commands

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"challenge/pkg/services"

	"challenge/pkg/utils"

	"github.com/spf13/cobra"
)

type Config struct {
	StateFile  string
	CitiesFile string
	LogLevel   string
	LogFormat  string
}

type App struct {
	service  *services.DistributionService
	jsonRepo *utils.JSONUtil
	csvRepo  *utils.CSVUtil
	config   *Config
	logger   *slog.Logger
}

func NewApp(logger *slog.Logger, config *Config) *App {
	return &App{
		service:  services.NewDistributionService(logger),
		jsonRepo: utils.NewJSONUtil(logger),
		csvRepo:  utils.NewCSVUtil(logger),
		config:   config,
		logger:   logger,
	}
}

func InitLogger(level string) *slog.Logger {
	var logLevel slog.Level
	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)

	return slog.New(handler)
}

// LoadData loads cities and state
func (a *App) LoadData() error {
	a.logger.Info("starting distribution system",
		slog.String("version", "1.0.0"),
		slog.String("log_level", a.config.LogLevel),
		slog.String("log_format", a.config.LogFormat))

	// Load cities if file exists
	if a.config.CitiesFile != "" {
		if _, err := os.Stat(a.config.CitiesFile); err == nil {
			locations, err := a.csvRepo.LoadCities(a.config.CitiesFile)
			if err != nil {
				a.logger.Warn("could not load cities file",
					slog.String("error", err.Error()))
			} else {
				a.service.SetLocations(locations)
			}
		} else {
			a.logger.Debug("cities file not found, skipping",
				slog.String("filename", a.config.CitiesFile))
		}
	}

	// Load state if file exists
	if a.config.StateFile != "" {
		if a.jsonRepo.Exists(a.config.StateFile) {
			distributors, err := a.jsonRepo.Load(a.config.StateFile)
			if err != nil {
				a.logger.Warn("could not load state file",
					slog.String("error", err.Error()))
			} else {
				a.service.SetDistributors(distributors)
			}
		} else {
			a.logger.Debug("state file not found, starting fresh",
				slog.String("filename", a.config.StateFile))
		}
	}

	return nil
}

func (a *App) SaveData() error {
	if a.config.StateFile != "" {
		return a.jsonRepo.Save(a.config.StateFile, a.service.GetDistributors())
	}
	return nil
}

func (a *App) BuildRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "distrib",
		Short: "Distribution Permission Management System",
		Long:  "A CLI tool to manage movie distribution permissions across geographical territories",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			a.LoadData()
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if err := a.SaveData(); err != nil {
				a.logger.Error("failed to save state",
					slog.String("error", err.Error()))
			}
		},
	}

	rootCmd.AddCommand(
		a.buildAddCommand(),
		a.buildIncludeCommand(),
		a.buildExcludeCommand(),
		a.buildCheckCommand(),
		a.buildShowCommand(),
		a.buildListCommand(),
	)

	return rootCmd
}

func (a *App) buildAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [distributor-name]",
		Short: "Add a new distributor",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			parent, _ := cmd.Flags().GetString("parent")
			err := a.service.AddDistributor(args[0], parent)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Added distributor: %s\n", args[0])
			if parent != "" {
				fmt.Printf("Parent: %s\n", parent)
			}
		},
	}
	cmd.Flags().StringP("parent", "p", "", "Parent distributor")
	return cmd
}

func (a *App) buildIncludeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "include [distributor] [location]",
		Short: "Add an include permission for a distributor",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			err := a.service.AddPermission(args[0], true, args[1])
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Added INCLUDE permission for %s: %s\n", args[0], args[1])
		},
	}
}

func (a *App) buildExcludeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "exclude [distributor] [location]",
		Short: "Add an exclude permission for a distributor",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			err := a.service.AddPermission(args[0], false, args[1])
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Added EXCLUDE permission for %s: %s\n", args[0], args[1])
		},
	}
}

func (a *App) buildCheckCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "check [distributor] [location]",
		Short: "Check if a distributor can distribute in a location",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			can, err := a.service.CanDistribute(args[0], args[1])
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			if can {
				fmt.Println("Yes can distribute")
			} else {
				fmt.Println("Can't distribute")
			}
		},
	}
}

func (a *App) buildShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show [distributor]",
		Short: "Show distributor information",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dist, err := a.service.GetDistributor(args[0])
			if err != nil {
				fmt.Printf("Distributor %s not found\n", args[0])
				os.Exit(1)
			}

			fmt.Printf("\n┌─ Distributor: %s\n", dist.Name)
			if dist.Parent != "" {
				fmt.Printf("├─ Parent: %s\n", dist.Parent)
			}
			fmt.Printf("├─ Permissions:\n")

			if len(dist.Permissions) == 0 {
				fmt.Printf("│  (none)\n")
			} else {
				for i, perm := range dist.Permissions {
					prefix := "├──"
					if i == len(dist.Permissions)-1 {
						prefix = "└──"
					}

					permType := "EXCLUDE"
					if perm.IsInclude {
						permType = "INCLUDE"
					}
					fmt.Printf("│  %s %s: %s\n", prefix, permType, perm.Location.String())
				}
			}
			fmt.Println()
		},
	}
}

func (a *App) buildListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all distributors",
		Run: func(cmd *cobra.Command, args []string) {
			distributors := a.service.GetDistributors()
			if len(distributors) == 0 {
				fmt.Println("No distributors found")
				return
			}

			fmt.Println("\nDistributors:")
			for name, dist := range distributors {
				parent := "root"
				if dist.Parent != "" {
					parent = dist.Parent
				}
				fmt.Printf("  • %s (parent: %s, permissions: %d)\n", name, parent, len(dist.Permissions))
			}
			fmt.Println()
		},
	}
}
