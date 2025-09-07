package cmd

import (
	"os"

	"github.com/Zigl3ur/mcli/internal/cli/commands/addons"
	"github.com/Zigl3ur/mcli/internal/cli/commands/backup"
	"github.com/Zigl3ur/mcli/internal/cli/commands/jar"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "mcli",
	Short:        "Simple cli tool to easily manage minecraft server",
	Long:         "A simple cli tool to make minecraft server management simple by providing jar download, backup and many other tools",
	Version:      "0.0.1",
	SilenceUsage: true, // do not show usage on errors
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().Bool("verbose", false, "show debug logs")

	rootCmd.Flags().SortFlags = false
	rootCmd.AddCommand(jar.NewCommand(), addons.NewCommand(), backup.NewCommand())
}
