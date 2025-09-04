package utils

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

func Iso8601Format(date string) (string, error) {
	layout := "2006-01-02T15:04:05Z"
	t, err := time.Parse(layout, date)

	if err != nil {
		return "", err
	}

	return t.Format("Jan 2, 2006, 03:04 PM"), nil
}

func SortMcVersions(versions []string) map[string][]string {
	versionsData := make(map[string][]string, 2)
	versionsData["versions"] = make([]string, 0)
	versionsData["snapshots"] = make([]string, 0)

	for _, version := range versions {
		_, unparsed := mcVersionParser(version)
		if unparsed != "" {
			versionsData["snapshots"] = append(versionsData["snapshots"], version)
		} else {
			versionsData["versions"] = append(versionsData["versions"], version)
		}
	}

	slices.SortStableFunc(versionsData["versions"], func(i, j string) int {
		ver0, _ := mcVersionParser(i)
		ver1, _ := mcVersionParser(j)

		for idx := range 3 {
			if ver0[idx] > ver1[idx] {
				return -1
			}
			if ver0[idx] < ver1[idx] {
				return 1
			}
		}

		return 0
	})

	return versionsData
}

func mcVersionParser(version string) ([3]int, string) {
	parts := strings.SplitN(version, ".", 3)

	if len(parts) < 2 {
		return [3]int{}, version
	}

	var mainVersion, subVersion, patch int

	mainVersion, err := strconv.Atoi(parts[0])
	if err != nil {
		return [3]int{}, version
	}
	subVersion, err = strconv.Atoi(parts[1])
	if err != nil {
		return [3]int{}, version
	}

	if len(parts) >= 3 {
		patch, err = strconv.Atoi(parts[2])
		if err != nil {
			return [3]int{}, version
		}
	}

	return [3]int{mainVersion, subVersion, patch}, ""
}

func humanizeByte(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTP"[exp])
}
