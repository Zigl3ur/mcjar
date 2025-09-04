package utils

import (
	"fmt"
	"strings"
)

// facets builder create facets string with given data,
// facets is a query param for modrinth api
func FacetsBuilder(versions []string, loader, projectType string) string {
	elt := make([]string, 0, 3)

	if len(versions) > 0 {
		velt := make([]string, 0, len(versions))

		for _, v := range versions {
			velt = append(velt, fmt.Sprintf("\"versions:%s\"", v))
		}

		elt = append(elt, fmt.Sprintf("[%s]", strings.Join(velt, ",")))
	}

	if loader != "" {
		elt = append(elt, fmt.Sprintf("[\"categories:%s\"]", loader))
	}

	if projectType != "" {
		elt = append(elt, fmt.Sprintf("[\"project_type:%s\"]", projectType))
	}

	if len(elt) == 0 {
		return ""
	}

	return fmt.Sprintf("[%s]", strings.Join(elt, ","))
}
