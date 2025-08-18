package jar

import (
	"fmt"
	"log"

	"github.com/Zigl3ur/mcli/internal/cli/commands/jar/list"
	"github.com/Zigl3ur/mcli/internal/cli/flags"
	"github.com/Zigl3ur/mcli/internal/handlers/fabric"
	"github.com/Zigl3ur/mcli/internal/handlers/neoforge"
	"github.com/Zigl3ur/mcli/internal/handlers/paper"
	"github.com/Zigl3ur/mcli/internal/handlers/purpur"
	"github.com/Zigl3ur/mcli/internal/handlers/vanilla"
	"github.com/Zigl3ur/mcli/internal/utils"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jar",
		Short: "Download the server jar file based on specified args",
		Long: `Download a server jar file based on given arguments, for example:
		mcli -v 1.8.9 -t paper -o ~/Downloads/paper_1.8.9.jar`,
		Run: execute,
	}

	cmd.Flags().StringP("type", "t", flags.Vanilla.String(), "the server type")
	cmd.Flags().StringP("version", "v", "1.21", "the server version")
	cmd.Flags().StringP("build", "b", "latest", "the server version build")
	cmd.Flags().StringP("output", "o", "server.jar", "the output path for the server jar file")

	cmd.AddCommand(list.NewCommand())

	return cmd
}

func execute(cmd *cobra.Command, args []string) {
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
	case flags.Fabric.String():
		if err := fabric.Handler(version, output); err != nil {
			log.Fatal(err)
		}
	case flags.Neoforge.String():
		if err := neoforge.Handler(version, build, output); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal(utils.InvalidServerType)
	}
}
