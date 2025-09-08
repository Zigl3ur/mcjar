package fabric

import (
	"errors"
	"fmt"

	"github.com/Zigl3ur/mcjar/internal/utils"
	"github.com/Zigl3ur/mcjar/internal/utils/loader"
)

var baseUrl = "https://meta.fabricmc.net/v2/versions"

func ListHandler(version string, snapshots bool) error {
	rawList, err := getVersionsList()
	if err != nil {
		return err
	}

	vlist := make([]string, 0, len(rawList.Versions))

	for _, v := range rawList.Versions {
		vlist = append(vlist, v.Version)
	}

	versionsMap := utils.SortMcVersions(vlist)
	loader.Stop()

	if !snapshots {
		for _, version := range versionsMap["versions"] {
			fmt.Printf("- %s\n", version)
		}
	} else {
		for _, snapshot := range versionsMap["snapshots"] {
			fmt.Printf("- %s\n", snapshot)
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

func getUrl(version string) (string, error) {
	loader, err := getStableLoader()
	if err != nil {
		return "", err
	}

	installer, err := getStableInstaller()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/loader/%s/%s/%s/server/jar", baseUrl, version, loader, installer), nil
}

type FabricVersion struct {
	Versions []struct {
		Version string `json:"version"`
		Stable  bool   `json:"stable"`
	} `json:"game"`
}

func getVersionsList() (FabricVersion, error) {

	var versions FabricVersion
	if status, err := utils.GetReqJson(baseUrl, &versions); err != nil {
		return versions, fmt.Errorf("failed to fetch Fabric versions from API (HTTP %d): %w", status, err)
	}

	return versions, nil
}

func getStableLoader() (string, error) {

	type LoaderList []struct {
		Version string `json:"version"`
		Stable  bool   `json:"stable"`
	}

	var list LoaderList
	if status, err := utils.GetReqJson(fmt.Sprintf("%s/loader", baseUrl), &list); err != nil {
		return "", fmt.Errorf("failed to fetch Fabric loader versions from API (HTTP %d): %w", status, err)
	}

	for _, loader := range list {
		if loader.Stable {
			return loader.Version, nil
		}
	}

	return "", errors.New("no stable fabric loader found")
}

func getStableInstaller() (string, error) {
	type InstallerList []struct {
		Version string `json:"version"`
		Stable  bool   `json:"stable"`
	}

	var list InstallerList
	if status, err := utils.GetReqJson(fmt.Sprintf("%s/versions/installer", baseUrl), &list); err != nil {
		return "", fmt.Errorf("failed to fetch Fabric installer versions from API (HTTP %d): %w", status, err)
	}

	for _, installer := range list {
		if installer.Stable {
			return installer.Version, nil
		}
	}

	return "", errors.New("no stable Fabric installer version found")
}
