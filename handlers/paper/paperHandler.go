package paper

import (
	"errors"
	"fmt"

	"github.com/Zigl3ur/mc-jar-fetcher/utils"
)

func Handler(version, path string) error {
	url, err := getUrlPaper(version)
	if err != nil {
		return err
	}
	return utils.WriteToFs(url, path)
}

type PaperVersions struct {
	Projects []struct {
		Project struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"project"`
		Versions map[string][]string `json:"versions"`
	} `json:"projects"`
}

func getUrlPaper(version string) (string, error) {
	type PaperUrl struct {
		Downloads struct {
			ServerDefault struct {
				Url string `json:"url"`
			} `json:"server:default"`
		} `json:"downloads"`
	}

	var paperUrl PaperUrl
	if err := utils.GetReq(fmt.Sprintf("https://fill.papermc.io/v3/projects/paper/versions/%s/builds/latest", version), &paperUrl); err != nil {
		return "", errors.New("paper doesnt support this version")
	}

	return paperUrl.Downloads.ServerDefault.Url, nil
}

func GetVersionsList() ([]string, error) {

	var versions PaperVersions
	if err := utils.GetReq("https://fill.papermc.io/v3/projects", &versions); err != nil {
		fmt.Println(err)
		return nil, errors.New("failed to fetch paper versions")
	}

	var versionList []string

	for _, p := range versions.Projects {
		if p.Project.Id == "paper" {
			for k, v := range p.Versions {
				fmt.Printf("%s - %s\n", k, v)
			}
		}
	}

	return versionList, nil
}
