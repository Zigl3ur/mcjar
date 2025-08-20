package fabric

import (
	"errors"
	"fmt"
	"log"

	"github.com/Zigl3ur/mcli/internal/utils"
	"github.com/Zigl3ur/mcli/internal/utils/loader"
)

func VersionsListHandler(version string, snapshots bool) {
	rawList, err := getVersionsList()
	if err != nil {
		log.Fatal(err)
	}

	vlist := make([]string, 0, len(rawList.Versions))

	for _, v := range rawList.Versions {
		vlist = append(vlist, v.Version)
	}

	versionsMap := utils.SortMcVersions(vlist)
	loader.Stop()

	if !snapshots {
		for _, v := range versionsMap["versions"] {
			fmt.Printf("- %s\n", v)
		}
	} else {
		for _, s := range versionsMap["snapshots"] {
			fmt.Printf("- %s\n", s)
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

func getUrl(version string) (string, error) {
	loader, err := getStableLoader()
	if err != nil {
		return "", err
	}

	installer, err := getStableInstaller()
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

func getVersionsList() (FabricVersion, error) {

	var versions FabricVersion
	if err := utils.GetReqJson("https://meta.fabricmc.net/v2/versions", &versions); err != nil {
		return versions, errors.New("failed to fetch fabric versions")
	}

	return versions, nil
}

func getStableLoader() (string, error) {

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

func getStableInstaller() (string, error) {

	type InstallerList []struct {
		Version string `json:"version"`
		Stable  bool   `json:"stable"`
	}

	var list InstallerList
	if err := utils.GetReqJson("https://meta.fabricmc.net/v2/versions/installer", &list); err != nil {
		return "", errors.New("failed to fetch fabric installers")
	}

	for _, l := range list {
		if l.Stable {
			return l.Version, nil
		}
	}

	return "", errors.New("no stable fabric installer found")
}
