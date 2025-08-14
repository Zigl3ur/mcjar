package fabric

import (
	"errors"
	"fmt"

	"github.com/Zigl3ur/mcli/internal/utils"
)

func Handler(version, path string) error {
	url, err := getUrl(version)
	if err != nil {
		return err
	}

	return utils.WriteToFs(url, path)
}

func getUrl(version string) (string, error) {
	loader, err := GetStableLoader()
	if err != nil {
		return "", err
	}

	installer, err := GetStableInstaller()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://meta.fabricmc.net/v2/versions/loader/%s/%s/%s/server/jar", version, loader, installer), nil
}

type FabricVersion struct {
	Versions []struct {
		Version string `json:"version"`
		Stable  bool   `json:"stable"`
	} `json:"game"`
}

func GetVersionsList() (FabricVersion, error) {

	var versions FabricVersion
	if err := utils.GetReqJson("https://meta.fabricmc.net/v2/versions", &versions); err != nil {
		return versions, errors.New("failed to fetch fabric versions")
	}

	return versions, nil
}

func GetStableLoader() (string, error) {

	type LoaderList []struct {
		Version string `json:"version"`
		Stable  bool   `json:"stable"`
	}

	var list LoaderList
	if err := utils.GetReqJson("https://meta.fabricmc.net/v2/versions/loader", &list); err != nil {
		return "", errors.New("failed to fetch fabric loaders")
	}

	for _, l := range list {
		if l.Stable {
			return l.Version, nil
		}
	}

	return "", errors.New("no stable fabric loader found")
}

func GetStableInstaller() (string, error) {

	type InstallerList []struct {
		Version string `json:"version"`
		Stable  bool   `json:"stable"`
	}

	var list InstallerList
	if err := utils.GetReqJson("https://meta.fabricmc.net/v2/versions/installer", &list); err != nil {
		return "", errors.New("failed to fetch fabric installer")
	}

	for _, l := range list {
		if l.Stable {
			return l.Version, nil
		}
	}

	return "", errors.New("no stable fabric installer found")
}
