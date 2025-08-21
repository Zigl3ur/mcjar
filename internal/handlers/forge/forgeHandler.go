package forge

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"slices"
	"strings"

	"github.com/Zigl3ur/mcli/internal/utils"
	"github.com/Zigl3ur/mcli/internal/utils/loader"
)

func ListHandler(version string, versionChanged, snapshots bool) {
	rawList, err := getVersionsList()
	if err != nil {
		log.Fatal(err)
	}
	loader.Stop()

	vlist := make([]string, 0, len(rawList))

	for k := range rawList {
		vlist = append(vlist, k)
	}

	versionsMap := utils.SortMcVersions(vlist)
	loader.Stop()

	if versionChanged {
		if slices.Contains(versionsMap["versions"], version) || slices.Contains(versionsMap["snapshots"], version) {
			fmt.Printf("- %s\n", version)
			for _, b := range rawList[version] {
				fmt.Printf("  - %s\n", b)
			}
		} else {
			log.Fatalf("forge doesnt support this version (given: %s)", version)
		}
	} else if snapshots {
		if len(versionsMap["snapshots"]) > 0 {
			for _, s := range versionsMap["snapshots"] {
				fmt.Printf("- %s\n", s)
			}
		} else {
			log.Fatal("forge doesn't support snapshots")
		}
	} else {
		for _, v := range versionsMap["versions"] {
			fmt.Printf("- %s\n", v)
		}
	}
}

func JarHandler(version, build, path string) error {
	url, err := getUrl(version, build)
	if err != nil {
		return err
	}

	if err = utils.WriteToFs(url, path); err != nil {
		return err
	}

	java, err := utils.GetPath("java")
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
	vlist, err := getVersionsList()
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

func getVersionsList() (map[string][]string, error) {
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

	for v := range versionMap {
		slices.Reverse(versionMap[v])
	}

	return versionMap, nil
}
