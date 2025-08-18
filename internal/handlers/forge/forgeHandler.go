package forge

import (
	"errors"
	"strings"

	"github.com/Zigl3ur/mcli/internal/utils"
)

// func Handler(version, path string) error {
// url, err := getUrl(version)
// if err != nil {
// return err
// }
//
// return utils.WriteToFs(url, path)
// }

// func getUrl(version string) (string, error) {

// 	return fmt.Sprintf("https://meta.fabricmc.net/v2/versions/loader/%s/%s/%s/server/jar", version, loader, installer), nil
// }

func GetVersionsList() (map[string][]string, error) {
	type ForgeVersions struct {
		Versioning struct {
			Latest   string   `xml:"latest"`
			Release  string   `xml:"release"`
			Versions []string `xml:"versions>version"`
		} `xml:"versioning"`
	}

	var versions ForgeVersions
	if err := utils.GetReqXml("https://maven.minecraftforge.net/net/minecraftforge/forge/maven-metadata.xml", &versions); err != nil {
		return nil, errors.New("failed to fetch forge versions")
	}

	versionMap := make(map[string][]string)

	for _, v := range versions.Versioning.Versions {
		parts := strings.Split(v, "-")
		if len(parts) >= 2 {
			version := parts[0]
			build := parts[1]
			versionMap[version] = append(versionMap[version], build)
		}
	}

	return versionMap, nil
}
