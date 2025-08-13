package paper

import (
	"errors"
	"fmt"

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

type PaperVersions struct {
	Versions []struct {
		Version struct {
			Id      string `json:"id"`
			Support struct {
				Status string `json:"status"`
			} `json:"support"`
		} `json:"version"`
		Builds []int `json:"builds"`
	} `json:"versions"`
}

func GetVersionsList() (PaperVersions, error) {

	var versions PaperVersions
	if err := utils.GetReq("https://fill.papermc.io/v3/projects/paper/versions", &versions); err != nil {
		return versions, errors.New("failed to fetch paper versions")
	}

	return versions, nil
}

type PaperBuild struct {
	Id int `json:"id"`
}

func GetBuildList(version string) ([]PaperBuild, error) {

	var builds []PaperBuild
	if err := utils.GetReq(fmt.Sprintf("https://fill.papermc.io/v3/projects/paper/versions/%s/builds?channel=STABLE", version), &builds); err != nil {
		return nil, errors.New("failed to fetch paper build list")
	}

	return builds, nil
}
