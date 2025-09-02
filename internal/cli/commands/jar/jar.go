package jar

import (
	"fmt"
	"path/filepath"
	"slices"

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
		Use:     "jar",
		Short:   "Download a server jar file based on given args",
		Long:    "Download a server jar file based on given args",
		PreRunE: validate,
		RunE:    execute,
	}

	cmd.Flags().StringP("type", "t", "", "the server type (required)")
	cmd.Flags().StringP("version", "v", "1.21", "the server version")
	cmd.Flags().StringP("build", "b", "latest", "the server version build")
	cmd.Flags().StringP("destination", "d", "", "the folder destination for the server jar file")

	//nolint:errcheck
	cmd.MarkFlagRequired("type")
	//nolint:errcheck
	cmd.MarkFlagDirname("destination")

	cmd.Flags().SortFlags = false
	cmd.AddCommand(list.NewCommand())

	return cmd
}

func validate(cmd *cobra.Command, args []string) error {
	serverType, _ := cmd.Flags().GetString("type")

	if cmd.Flag("type").Changed && !slices.Contains(flags.ValidServerType, serverType) {
		return fmt.Errorf("invalid type, valid ones are %s", flags.ValidServerType)
	}

	return nil
}

func execute(cmd *cobra.Command, args []string) error {
	serverType, _ := cmd.Flags().GetString("type")
	version, _ := cmd.Flags().GetString("version")
	build, _ := cmd.Flags().GetString("build")
	dir, _ := cmd.Flags().GetString("destination")
	isVerbose, _ := cmd.Flags().GetBool("verbose")

	var filename string
	if serverType != flags.Vanilla.String() {
		filename = fmt.Sprintf("%s-%s-%s.jar", serverType, version, build)
	} else {
		filename = fmt.Sprintf("%s-%s.jar", serverType, version)
	}

	outPath := filepath.Join(dir, filename)

	fmt.Println("Using Values:")
	fmt.Printf("- type: %s\n", serverType)
	fmt.Printf("- version: %s\n", version)
	if serverType != flags.Vanilla.String() {
		fmt.Printf("- build: %s\n", build)
	}
	fmt.Printf("- output: %s\n", outPath)

	switch serverType {
	case flags.Vanilla.String():
		return vanilla.JarHandler(version, outPath)
	case flags.Paper.String():
		return paper.JarHandler(flags.Paper.String(), version, build, outPath)
	case flags.Folia.String():
		return paper.JarHandler(flags.Folia.String(), version, build, outPath)
	case flags.Velocity.String():
		return paper.JarHandler(flags.Velocity.String(), version, build, outPath)
	case flags.Purpur.String():
		return purpur.JarHandler(version, build, outPath)
	case flags.Fabric.String():
		return fabric.JarHandler(version, outPath)
	case flags.Neoforge.String():
		return neoforge.JarHandler(version, build, outPath, isVerbose)
	case flags.Forge.String():
		return forge.JarHandler(version, build, outPath, isVerbose)
	default:
		//nolint:errcheck
		cmd.Usage()
		return nil
	}
}
