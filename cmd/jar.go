package cmd

import (
	"fmt"
	"log"

	"github.com/Zigl3ur/mcli/cmd/flags"
	"github.com/Zigl3ur/mcli/internal/handlers/paper"
	"github.com/Zigl3ur/mcli/internal/handlers/purpur"
	"github.com/Zigl3ur/mcli/internal/handlers/vanilla"
	"github.com/spf13/cobra"
)

const invalidServerType string = "Invalid server type, valid ones are [vanilla, paper, spigot, purpur, forge, fabric]"

var jarCmd = &cobra.Command{
	Use:   "jar",
	Short: "Download the server jar file based on specified args",
	Long: `Download a server jar file based on given arguments, for example:
		mcli -v 1.8.9 -t paper -o ~/Downloads/paper_1.8.9.jar`,
	Run: execute,
}

func init() {
	rootCmd.AddCommand(jarCmd)

	jarCmd.Flags().StringP("type", "t", flags.Vanilla.String(), "the server type")
	jarCmd.Flags().StringP("version", "v", "1.21", "the server version")
	jarCmd.Flags().StringP("build", "b", "latest", "the server version build")
	jarCmd.Flags().StringP("output", "o", "server.jar", "the output path for the server jar file")
}

func execute(cmd *cobra.Command, _ []string) {
	serverType := cmd.Flag("type").Value.String()
	version := cmd.Flag("version").Value.String()
	build := cmd.Flag("build").Value.String()
	output := cmd.Flag("output").Value.String()

	fmt.Println("Using Values:")
	fmt.Printf("- type: %s\n", serverType)
	fmt.Printf("- version: %s\n", version)
	if serverType != flags.Vanilla.String() {
		fmt.Printf("- build: %s\n", build)
	}
	fmt.Printf("- output: %s\n", output)

	switch serverType {
	case flags.Vanilla.String():
		if err := vanilla.Handler(version, output); err != nil {
			log.Fatal(err)
		}
	case flags.Paper.String():
		if err := paper.Handler(version, build, output); err != nil {
			log.Fatal(err)
		}
	case flags.Purpur.String():
		if err := purpur.Handler(version, build, output); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal(invalidServerType)
	}
}
