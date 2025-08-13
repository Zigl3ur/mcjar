package fabric

import (
	"errors"
	"fmt"

	"github.com/Zigl3ur/mcli/utils"
	"github.com/spf13/pflag"
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

	if pflag.Lookup("build").Changed {
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

type FabricVersion struct {
	Versions []struct {
		Version string `json:"version"`
		Stable  bool   `json:"stable"`
	} `json:"game"`
}

func GetVersionsList() (FabricVersion, error) {

	var versions FabricVersion
	if err := utils.GetReq("https://meta.fabricmc.net/v2/versions", &versions); err != nil {
		return versions, errors.New("failed to fetch fabric versions")
	}

	return versions, nil
}

func GetLoader() (string, error) {

	type LoaderList []struct {
		Version string `json:"version"`
		Stable  bool   `json:"stable"`
	}

	var list LoaderList
	if err := utils.GetReq("https://meta.fabricmc.net/v2/versions/loader", &list); err != nil {
		return list, errors.New("failed to fetch fabric loaders")
	}

	return versions, nil
}
