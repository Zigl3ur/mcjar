package paper

import (
	"errors"
	"fmt"
	"slices"

	"github.com/Zigl3ur/mcli/internal/utils"
	"github.com/Zigl3ur/mcli/internal/utils/loader"
)

func ListHandler(project, version string, versionChanged, snapshots bool) error {
	rawList, err := getVersionsList(project)
	if err != nil {
		return err
	}

	vlist := make([]string, 0, len(rawList.Versions))

	for _, v := range rawList.Versions {
		vlist = append(vlist, v.Version.Id)
	}

	versionsMap := utils.SortMcVersions(vlist)
	loader.Stop()

	if versionChanged {
		if slices.Contains(versionsMap["versions"], version) || slices.Contains(versionsMap["snapshots"], version) {
			blist, err := getBuildList(project, version)
			if err != nil {
				return err
			}
			fmt.Printf("- %s\n", version)
			for _, b := range blist {
				fmt.Printf("  - %d\n", b.Id)
			}
		} else {
			return fmt.Errorf("paper doesnt support this version (given: %s)", version)
		}
	} else if snapshots {
		if len(versionsMap["snapshots"]) > 0 {
			for _, s := range versionsMap["snapshots"] {
				fmt.Printf("- %s\n", s)
			}
		} else {
			return fmt.Errorf("%s doesn't support snapshots", project)
		}
	} else {
		for _, v := range versionsMap["versions"] {
			fmt.Printf("- %s\n", v)
		}
	}

	return nil
}

func JarHandler(project, version, build, outPath string) error {
	url, err := getUrl(project, version, build)
	if err != nil {
		return err
	}
	return utils.WriteToFs(url, outPath)
}

func getUrl(project, version, build string) (string, error) {
	type PaperUrl struct {
		Downloads struct {
			ServerDefault struct {
				Url string `json:"url"`
			} `json:"server:default"`
		} `json:"downloads"`
	}

	fetchUrl := fmt.Sprintf("https://fill.papermc.io/v3/projects/%s/versions/%s/builds/latest", project, version)
	errorMsg := fmt.Errorf("no %s jar available for provided version (given: %s)", project, version)

	if build != "latest" {
		fetchUrl = fmt.Sprintf("https://fill.papermc.io/v3/projects/%s/versions/%s/builds/%s", project, version, build)
		errorMsg = fmt.Errorf("no %s jar available for provided version / build (given: %s, %s)", project, version, build)
	}

	var paperUrl PaperUrl
	if err := utils.GetReqJson(fetchUrl, &paperUrl); err != nil {
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

func getVersionsList(project string) (PaperVersions, error) {

	var versions PaperVersions
	if err := utils.GetReqJson(fmt.Sprintf("https://fill.papermc.io/v3/projects/%s/versions", project), &versions); err != nil {
		return versions, fmt.Errorf("failed to fetch %s version list", project)
	}

	return versions, nil
}

type PaperBuild struct {
	Id int `json:"id"`
}

func getBuildList(project, version string) ([]PaperBuild, error) {

	var builds []PaperBuild
	if err := utils.GetReqJson(fmt.Sprintf("https://fill.papermc.io/v3/projects/%s/versions/%s/builds?channel=STABLE", project, version), &builds); err != nil {
		return nil, fmt.Errorf("failed to fetch %s build list", project)
	}

	return builds, nil
}
