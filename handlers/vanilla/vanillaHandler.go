package vanilla

import (
	"errors"
	"fmt"

	"github.com/Zigl3ur/mc-jar-fetcher/utils"
)

func Handler(version, path string) error {
	url, err := getUrl(version)
	if err != nil {
		return err
	}
	return utils.WriteToFs(url, path)
}

type Versions struct {
	Versions []struct {
		Id  string `json:"id"`
		Url string `json:"url"`
	}
}

func getUrl(version string) (string, error) {
	versions, err := GetVersionsList()
	if err != nil {
		return "", err
	}

	var versionUrl string
	for _, v := range versions.Versions {
		if v.Id == version {
			versionUrl = v.Url
			break
		}
	}

	if versionUrl == "" {
		return "", fmt.Errorf("specified version not found (given: %s)", version)
	}

	type DownloadData struct {
		Downloads struct {
			Server struct {
				Url string `json:"url"`
			} `json:"server"`
		} `json:"downloads"`
	}

	var downloadData DownloadData
	if err := utils.GetReq(versionUrl, &downloadData); err != nil {
		return "", errors.New("failed to fetch version details")
	}

	serverUrl := downloadData.Downloads.Server.Url
	if serverUrl == "" {
		return "", fmt.Errorf("no vanilla jar available for specified version (given: %s)", version)
	}

	return serverUrl, nil
}

func GetVersionsList() (Versions, error) {

	var versions Versions
	if err := utils.GetReq("https://launchermeta.mojang.com/mc/game/version_manifest.json", &versions); err != nil {
		return versions, errors.New("failed to fetch version manifest")
	}

	return versions, nil
}
