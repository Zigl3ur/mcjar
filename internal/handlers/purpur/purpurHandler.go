package purpur

import (
	"errors"
	"fmt"
	"slices"

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

	vlist, err := GetVersionsList()
	if err != nil {
		return "", err
	}

	if !slices.Contains(vlist, version) {
		return "", fmt.Errorf("no purpur jar available for provided version (given: %s)", version)
	}

	blist, err := GetBuildList(version)
	if err != nil {
		return "", err
	}

	if build == "" {
		build = blist[len(blist)-1]
	} else {
		if !slices.Contains(blist, build) {
			return "", fmt.Errorf("no purpur jar available for provided version / build (given: %s, %s)", version, build)
		}
	}

	return fmt.Sprintf("https://api.purpurmc.org/v2/purpur/%s/%s/download", version, build), nil
}

func GetVersionsList() ([]string, error) {

	type PurpurVersion struct {
		List []string `json:"versions"`
	}

	var versions PurpurVersion
	if err := utils.GetReq("https://api.purpurmc.org/v2/purpur", &versions); err != nil {
		return nil, errors.New("failed to fetch version details")
	}
	return versions.List, nil
}

func GetBuildList(version string) ([]string, error) {

	type PurpurBuilds struct {
		Builds struct {
			List []string `json:"all"`
		} `json:"builds"`
	}

	var builds PurpurBuilds
	if err := utils.GetReq(fmt.Sprintf("https://api.purpurmc.org/v2/purpur/%s", version), &builds); err != nil {
		return nil, errors.New("failed to fetch purpur build list")
	}

	return builds.Builds.List, nil
}
