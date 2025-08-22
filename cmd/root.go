package cmd

import (
	"os"

	"github.com/Zigl3ur/mcli/internal/cli/commands/datapack"
	"github.com/Zigl3ur/mcli/internal/cli/commands/jar"
	"github.com/Zigl3ur/mcli/internal/cli/commands/mod"
	"github.com/Zigl3ur/mcli/internal/cli/commands/modpack"
	"github.com/Zigl3ur/mcli/internal/cli/commands/plugin"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "mcli",
	Short:   "Simple cli tool to easily manage minecraft server",
	Long:    "A simple cli tool to make minecraft server management simple by providing jar download, backup and many other tools",
	Version: "0.0.1",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("verbose", "v", false, "show debug logs")

	rootCmd.Flags().SortFlags = false
	rootCmd.AddCommand(jar.NewCommand())
	rootCmd.AddCommand(plugin.NewCommand())
	rootCmd.AddCommand(mod.NewCommand())
	rootCmd.AddCommand(modpack.NewCommand())
	rootCmd.AddCommand(datapack.NewCommand())

}
