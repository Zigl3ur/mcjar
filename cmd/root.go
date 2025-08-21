/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Zigl3ur/mcli/internal/cli/commands/jar"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mcli",
	Short: "Simple cli tool to easily manage minecraft server",
	Long:  "A simple cli tool to make minecraft server management simple by providing jar download, world saving, and rcon",
	Run:   execute,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("verbose", "v", false, "show debug logs")
	rootCmd.Flags().Bool("version", false, "display mcli version")

	rootCmd.Flags().SortFlags = false
	rootCmd.AddCommand(jar.NewCommand())
}

func execute(cmd *cobra.Command, args []string) {
	versionToggle, _ := cmd.Flags().GetBool("version")

	if versionToggle {
		fmt.Println("mcli version 0.0.1")
	} else {
		if err := cmd.Usage(); err != nil {
			// shouldn't happen
			log.Fatal("failed to print help message")
		}
	}
}
