package paper

import (
	"errors"
	"fmt"

	"github.com/Zigl3ur/mc-jar-fetcher/utils"
)

func Handler(version, build, path string) error {
	url, err := getUrlPaper(version, build)
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

func getUrlPaper(version, build string) (string, error) {
	type PaperUrl struct {
		Downloads struct {
			ServerDefault struct {
				Url string `json:"url"`
			} `json:"server:default"`
		} `json:"downloads"`
	}

	fetchUrl := fmt.Sprintf("https://fill.papermc.io/v3/projects/paper/versions/%s/builds/latest", version)
	if build != "" {
		fetchUrl = fmt.Sprintf("https://fill.papermc.io/v3/projects/paper/versions/%s/builds/%s", version, build)
	}

	var paperUrl PaperUrl
	if err := utils.GetReq(fetchUrl, &paperUrl); err != nil {
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

func GetBuildList(version string) ([]int, error) {
	type PaperBuild struct {
		Id int `json:"id"`
	}

	var builds []PaperBuild
	if err := utils.GetReq(fmt.Sprintf("https://fill.papermc.io/v3/projects/paper/versions/%s/builds?channel=STABLE", version), &builds); err != nil {
		return nil, errors.New("failed to fetch paper build list")
	}

	fmt.Println(builds)

	buildsList := []int{}

	for _, b := range builds {
		buildsList = append(buildsList, b.Id)
	}

	return buildsList, nil
}
