package neoforge

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Zigl3ur/mcli/internal/utils"
)

func Handler(version, build, path string) error {
	url, err := getUrl(version, build)
	if err != nil {
		return err
	}
	return utils.WriteToFs(url, path)
}

func getUrl(version, build string) (string, error) {
	type PaperUrl struct {
		Downloads struct {
			ServerDefault struct {
				Url string `json:"url"`
			} `json:"server:default"`
		} `json:"downloads"`
	}

	fetchUrl := fmt.Sprintf("https://fill.papermc.io/v3/projects/paper/versions/%s/builds/latest", version)
	errorMsg := fmt.Errorf("no paper jar available for provided version (given: %s)", version)

	if build != "latest" {
		fetchUrl = fmt.Sprintf("https://fill.papermc.io/v3/projects/paper/versions/%s/builds/%s", version, build)
		errorMsg = fmt.Errorf("no paper jar available for provided version / build (given: %s, %s)", version, build)
	}

	var paperUrl PaperUrl
	if err := utils.GetReq(fetchUrl, &paperUrl); err != nil {
		return "", errors.New("failed to fetch version details")
	}

	serverUrl := paperUrl.Downloads.ServerDefault.Url
	if serverUrl == "" {
		return "", errorMsg
	}

	return serverUrl, nil
}

func GetVersionsList() ([]string, error) {

	type NeoforgeVersions struct {
		Versions []string `json:"versions"`
	}

	var list NeoforgeVersions
	if err := utils.GetReq("https://maven.neoforged.net/api/maven/versions/releases/net/neoforged/neoforge", &list); err != nil {
		return nil, errors.New("failed to fetch neoforge versions")
	}

	filteredVersions := []string{}
	seenVer := make(map[string]bool)

	for _, v := range list.Versions {
		// remove april fools versions
		if !strings.HasPrefix(v, "0") {
			parts := strings.SplitN(v, ".", 3)

			version := fmt.Sprintf("1.%s", parts[0])
			if len(parts) > 1 {
				version = fmt.Sprintf("1.%s.%s", parts[0], parts[1])
			}

			if !seenVer[version] {
				seenVer[version] = true
				filteredVersions = append(filteredVersions, version)
			}
		}
	}

	return filteredVersions, nil
}
