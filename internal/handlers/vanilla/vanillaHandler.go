package vanilla

import (
	"fmt"

	"github.com/Zigl3ur/mcli/internal/utils"
	"github.com/Zigl3ur/mcli/internal/utils/loader"
)

func ListHandler(snapshots bool) error {
	rawList, err := getVersionsList()
	if err != nil {
		return err
	}

	vlist := make([]string, 0, len(rawList.Versions))

	for _, v := range rawList.Versions {
		vlist = append(vlist, v.Id)
	}

	versionsMap := utils.SortMcVersions(vlist)

	loader.Stop()

	if !snapshots {
		for _, v := range versionsMap["versions"] {
			fmt.Printf("- %s\n", v)
		}
	} else {
		for _, v := range versionsMap["snapshots"] {
			fmt.Printf("- %s\n", v)
		}
	}

	return nil
}

func JarHandler(version, outPath string) error {
	url, err := getUrl(version)
	if err != nil {
		return err
	}
	return utils.WriteToFs(url, outPath)
}

type Versions struct {
	Versions []struct {
		Id  string `json:"id"`
		Url string `json:"url"`
	}
}

func getUrl(version string) (string, error) {
	versions, err := getVersionsList()
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
			} `jsons:"server"`
		} `json:"downloads"`
	}

	var downloadData DownloadData
	if status, err := utils.GetReqJson(versionUrl, &downloadData); err != nil {
		return "", fmt.Errorf("failed to fetch Vanilla download url from API (HTTP %d): %w", status, err)
	}

	serverUrl := downloadData.Downloads.Server.Url
	if serverUrl == "" {
		return "", err
	}

	return serverUrl, nil
}

func getVersionsList() (Versions, error) {

	var versions Versions
	if status, err := utils.GetReqJson("https://launchermeta.mojang.com/mc/game/version_manifest.json", &versions); err != nil {
		return versions, fmt.Errorf("failed to fetch Vanilla versions from API (HTTP %d): %w", status, err)
	}

	return versions, nil
}
