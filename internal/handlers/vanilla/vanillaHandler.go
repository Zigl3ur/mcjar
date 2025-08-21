package vanilla

import (
	"errors"
	"fmt"
	"log"

	"github.com/Zigl3ur/mcli/internal/utils"
	"github.com/Zigl3ur/mcli/internal/utils/loader"
)

func ListHandler(snapshots bool) {
	rawList, err := getVersionsList()
	if err != nil {
		log.Fatal(err)
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
}

func JarHandler(version, path string) error {
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
	if err := utils.GetReqJson(versionUrl, &downloadData); err != nil {
		return "", errors.New("failed to fetch version details")
	}

	serverUrl := downloadData.Downloads.Server.Url
	if serverUrl == "" {
		return "", fmt.Errorf("no vanilla jar available for specified version (given: %s)", version)
	}

	return serverUrl, nil
}

func getVersionsList() (Versions, error) {

	var versions Versions
	if err := utils.GetReqJson("https://launchermeta.mojang.com/mc/game/version_manifest.json", &versions); err != nil {
		return versions, errors.New("failed to fetch version manifest")
	}

	return versions, nil
}
