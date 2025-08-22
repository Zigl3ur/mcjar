package search

import (
	"fmt"
	"log"

	"github.com/Zigl3ur/mcli/internal/handlers/modrinth"
	"github.com/Zigl3ur/mcli/internal/utils"
	"github.com/Zigl3ur/mcli/internal/utils/loader"
	"github.com/spf13/cobra"
)

func NewCommand(parentName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: fmt.Sprintf("Search for %ss", parentName),
		Run: func(cmd *cobra.Command, args []string) {
			execute(cmd, args, parentName)
		},
	}

	cmd.Flags().Int("limit", 10, "number of results to display")
	cmd.Flags().StringP("index", "i", "relevance", "sort search result")
	cmd.Flags().StringArrayP("versions", "v", []string{}, "versions that results items must support")
	cmd.Flags().StringP("loader", "l", "", "the minecraft loader")

	cmd.Flags().SortFlags = false

	return cmd
}

func execute(cmd *cobra.Command, args []string, parentName string) {
	if len(args) < 1 {
		cmd.Usage()
		return
	}

	query := args[0]

	limit, _ := cmd.Flags().GetInt("limit")
	index := cmd.Flag("index").Value.String()
	versions, _ := cmd.Flags().GetStringArray("versions")
	mcLoader, _ := cmd.Flags().GetString("loader")

	loader.Start(fmt.Sprintf("Searching for \"%s\"", query))

	facets := utils.FacetsBuilder(versions, mcLoader, parentName)

	results, err := modrinth.Search(query, index, facets, limit)
	loader.Stop()
	if err != nil {
		log.Fatal(err)
	}

	if results.TotalHits == 0 {
		fmt.Printf("No results found for \"%s\"\n", query)
	} else {
		fmt.Printf("Found %d results for \"%s\" (limit: %d)\n", results.TotalHits, query, limit)
		for _, r := range results.Results {
			fmt.Printf("  - %s, %s\n", r.Title, r.Description)
		}
	}

}
