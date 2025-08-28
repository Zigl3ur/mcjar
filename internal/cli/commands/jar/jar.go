package jar

import (
	"fmt"

	"github.com/Zigl3ur/mcli/internal/cli/commands/jar/list"
	"github.com/Zigl3ur/mcli/internal/cli/flags"
	"github.com/Zigl3ur/mcli/internal/handlers/fabric"
	"github.com/Zigl3ur/mcli/internal/handlers/forge"
	"github.com/Zigl3ur/mcli/internal/handlers/neoforge"
	"github.com/Zigl3ur/mcli/internal/handlers/paper"
	"github.com/Zigl3ur/mcli/internal/handlers/purpur"
	"github.com/Zigl3ur/mcli/internal/handlers/vanilla"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jar",
		Short: "Download a server jar file based on given args",
		Long:  "Download a server jar file based on given args",
		RunE:  execute,
	}

	cmd.Flags().StringP("type", "t", "", "the server type")
	cmd.Flags().StringP("version", "v", "1.21", "the server version")
	cmd.Flags().StringP("build", "b", "latest", "the server version build")
	cmd.Flags().StringP("output", "o", "server.jar", "the output path for the server jar file")

	cmd.Flags().SortFlags = false
	cmd.AddCommand(list.NewCommand())

	return cmd
}

func execute(cmd *cobra.Command, args []string) error {

	serverType := cmd.Flag("type").Value.String()
	version := cmd.Flag("version").Value.String()
	build := cmd.Flag("build").Value.String()
	output := cmd.Flag("output").Value.String()
	isVerbose, _ := cmd.Flags().GetBool("verbose")

	if !cmd.Flag("output").Changed {
		if serverType != flags.Vanilla.String() {
			output = fmt.Sprintf("%s-%s-%s.jar", serverType, version, build)
		} else {
			output = fmt.Sprintf("%s-%s.jar", serverType, version)
		}
	}

	if serverType != "" {
		fmt.Println("Using Values:")
		fmt.Printf("- type: %s\n", serverType)
		fmt.Printf("- version: %s\n", version)
		if serverType != flags.Vanilla.String() {
			fmt.Printf("- build: %s\n", build)
		}
		fmt.Printf("- output: %s\n", output)
	}

	switch serverType {
	case flags.Vanilla.String():
		return vanilla.JarHandler(version, output)
	case flags.Paper.String():
		return paper.JarHandler(flags.Paper.String(), version, build, output)
	case flags.Folia.String():
		return paper.JarHandler(flags.Folia.String(), version, build, output)
	case flags.Velocity.String():
		return paper.JarHandler(flags.Velocity.String(), version, build, output)
	case flags.Purpur.String():
		return purpur.JarHandler(version, build, output)
	case flags.Fabric.String():
		return fabric.JarHandler(version, output)
	case flags.Neoforge.String():
		return neoforge.JarHandler(version, build, output, isVerbose)
	case flags.Forge.String():
		return forge.JarHandler(version, build, output, isVerbose)
	default:
		//nolint:errcheck
		cmd.Usage()
		return nil
	}
}
