package cmd

import (
	"os"

	"github.com/Zigl3ur/mcjar/internal/cli/commands/addons"
	"github.com/Zigl3ur/mcjar/internal/cli/commands/jar"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "mcjar",
	Short:        "A Simple cli tool to easily download jar for minecraft as server software and addons (plugins/mods/modpacks/data packs)",
	Long:         "A Simple cli tool to easily download jar for minecraft as server software and addons (plugins/mods/modpacks/data packs)",
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
	rootCmd.AddCommand(jar.NewCommand(), addons.NewCommand())
}
