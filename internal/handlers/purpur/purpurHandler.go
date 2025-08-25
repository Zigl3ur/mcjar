package purpur

import (
	"errors"
	"fmt"
	"log"
	"slices"

	"github.com/Zigl3ur/mcli/internal/utils"
	"github.com/Zigl3ur/mcli/internal/utils/loader"
)

func ListHandler(version string, versionChanged, snapshots bool) error {
	rawList, err := getVersionsList()
	if err != nil {
		return err
	}

	vlist := make([]string, 0, len(rawList))

	vlist = append(vlist, rawList...)

	versionsMap := utils.SortMcVersions(vlist)
	loader.Stop()

	if versionChanged {
		if slices.Contains(versionsMap["versions"], version) || slices.Contains(versionsMap["snapshots"], version) {
			blist, err := getBuildList(version)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("- %s\n", version)
			for _, b := range blist {
				fmt.Printf("  - %s\n", b)
			}
		} else {
			return fmt.Errorf("purpur doesn't support this version (given: %s)", version)
		}
	} else if snapshots {
		if len(versionsMap["snapshots"]) > 0 {
			for _, s := range versionsMap["snapshots"] {
				fmt.Printf("- %s\n", s)
			}
		} else {
			return errors.New("purpur doesn't support snapshots")
		}
	} else {
		for _, v := range versionsMap["versions"] {
			fmt.Printf("- %s\n", v)
		}
	}

	return nil
}

func JarHandler(version, build, path string) error {
	url, err := getUrl(version, build)
	if err != nil {
		return err
	}
	return utils.WriteToFs(url, path)
}

func getUrl(version, build string) (string, error) {

	vlist, err := getVersionsList()
	if err != nil {
		return "", err
	}

	if !slices.Contains(vlist, version) {
		return "", fmt.Errorf("no purpur jar available for provided version (given: %s)", version)
	}

	blist, err := getBuildList(version)
	if err != nil {
		return "", err
	}

	if build == "latest" {
		// if no build provided get the latest one
		build = blist[len(blist)-1]
	} else {
		if !slices.Contains(blist, build) {
			return "", fmt.Errorf("no purpur jar available for provided version / build (given: %s, %s)", version, build)
		}
	}

	return fmt.Sprintf("https://api.purpurmc.org/v2/purpur/%s/%s/download", version, build), nil
}

func getVersionsList() ([]string, error) {

	type PurpurVersion struct {
		List []string `json:"versions"`
	}

	var versions PurpurVersion
	if err := utils.GetReqJson("https://api.purpurmc.org/v2/purpur", &versions); err != nil {
		return nil, errors.New("failed to fetch version details")
	}

	slices.Reverse(versions.List)
	return versions.List, nil
}

func getBuildList(version string) ([]string, error) {

	type PurpurBuilds struct {
		Builds struct {
			List []string `json:"all"`
		} `json:"builds"`
	}

	var builds PurpurBuilds
	if err := utils.GetReqJson(fmt.Sprintf("https://api.purpurmc.org/v2/purpur/%s", version), &builds); err != nil {
		return nil, errors.New("failed to fetch purpur build list")
	}

	slices.Reverse(builds.Builds.List)

	return builds.Builds.List, nil
}
