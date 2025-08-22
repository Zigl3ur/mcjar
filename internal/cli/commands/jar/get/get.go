package get

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCommand(parentName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: fmt.Sprintf("Download a %s", parentName),
		Run: func(cmd *cobra.Command, args []string) {
			execute(cmd, args, parentName)
		},
	}

	cmd.Flags().StringP("destination", "d", ".", "the folder where to put the downloaded jar")

	cmd.Flags().SortFlags = false

	return cmd
}

func execute(cmd *cobra.Command, args []string, parentName string) {

}
