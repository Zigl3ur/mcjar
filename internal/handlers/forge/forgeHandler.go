package forge

import (
	"errors"
	"fmt"
	"os/exec"
	"slices"
	"strings"

	"github.com/Zigl3ur/mcli/internal/utils"
	"github.com/Zigl3ur/mcli/internal/utils/loader"
)

func JarHandler(version, build, path string) error {
	url, err := getUrl(version, build)
	if err != nil {
		return err
	}

	if err = utils.WriteToFs(url, path); err != nil {
		return err
	}

	java, err := utils.GetJava()
	if err != nil {
		return err
	}

	destElt := strings.Split(path, "/")
	dest := strings.Join(destElt[:len(destElt)-1], "/")
	cmd := exec.Command(java, "-jar", path, "--installServer", dest)
	loader.Start("Installing forge server")
	// use cmd.Output ? if adding a debug flag and print output ?
	if err = cmd.Run(); err != nil {
		loader.Stop()
		return errors.New("failed to install forge server")
	}

	loader.Stop()
	fmt.Printf("Installed forge server at %s\n", dest)

	return nil
}

func getUrl(version, build string) (string, error) {
	vlist, err := GetVersionsList()
	if err != nil {
		return "", err
	}

	if vlist[version] == nil {
		return "", fmt.Errorf("no forge jar available for provided version (given: %s)", version)
	}

	latestBuild := vlist[version][0]
	url := fmt.Sprintf("https://maven.minecraftforge.net/net/minecraftforge/forge/%s-%s/forge-%s-%s-installer.jar", version, latestBuild, version, latestBuild)

	if build != "latest" {
		if slices.Contains(vlist[version], build) {
			url = fmt.Sprintf("https://maven.minecraftforge.net/net/minecraftforge/forge/%s-%s/forge-%s-%s-installer.jar", version, build, version, build)
		} else {
			return "", fmt.Errorf("no forge jar available for provided version (given: %s, %s)", version, build)
		}
	}

	return url, nil
}

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
