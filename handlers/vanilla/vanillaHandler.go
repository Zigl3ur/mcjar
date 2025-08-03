package vanilla

import (
	"errors"

	"github.com/Zigl3ur/mc-jar-fetcher/utils"
)

func VanillaHandler(version, path string) error {
	url, err := getUrlVanilla(version)
	if err != nil {
		return err
	}
	return utils.WriteToFs(url, path)
}

func getUrlVanilla(version string) (string, error) {
	type Versions struct {
		Versions []struct {
			Id  string `json:"id"`
			Url string `json:"url"`
		}
	}

	var versions Versions
	if err := utils.GetReq("https://launchermeta.mojang.com/mc/game/version_manifest.json", &versions); err != nil {
		return "", errors.New("failed to fetch version manifest")
	}

	var versionUrl string
	for _, v := range versions.Versions {
		if v.Id == version {
			versionUrl = v.Url
			break
		}
	}

	if versionUrl == "" {
		return "", errors.New("specified version not found")
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
		return "", errors.New("failed to get download url")
	}

	return downloadData.Downloads.Server.Url, nil
}
