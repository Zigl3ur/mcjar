package info

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCommand(parentName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: fmt.Sprintf("See infos about a %s", parentName),
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
