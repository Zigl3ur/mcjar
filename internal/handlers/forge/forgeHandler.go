package forge

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/Zigl3ur/mcli/internal/utils"
	"github.com/Zigl3ur/mcli/internal/utils/loader"
)

func ListHandler(version string, versionChanged, snapshots bool) error {
	rawList, err := getVersionsList()
	if err != nil {
		return err
	}
	loader.Stop()

	vlist := make([]string, 0, len(rawList))

	for key := range rawList {
		vlist = append(vlist, key)
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
			return fmt.Errorf("forge doesnt support this version (given: %s)", version)
		}
	} else if snapshots {
		if len(versionsMap["snapshots"]) > 0 {
			for _, snapshot := range versionsMap["snapshots"] {
				fmt.Printf("- %s\n", snapshot)
			}
		} else {
			return fmt.Errorf("forge doesn't support snapshots")
		}
	} else {
		for _, version := range versionsMap["versions"] {
			fmt.Printf("- %s\n", version)
		}
	}

	return nil
}

func JarHandler(version, build, outPath string, isVerbose bool) error {
	url, err := getUrl(version, build)
	if err != nil {
		return err
	}

	if err = utils.WriteToFs(url, outPath); err != nil {
		return err
	}

	java, err := utils.GetPath("java")
	if err != nil {
		return err
	}

	dir, _ := filepath.Split(outPath)
	cmd := exec.Command(java, "-jar", outPath, "--installServer", dir)

	if isVerbose {
		cmd.Stdout = os.Stdout
	} else {
		loader.Start("Installing forge server")
	}

	if err = cmd.Run(); err != nil {
		loader.Stop()
		return err
	}

	loader.Stop()
	fmt.Printf("Installed forge server at %s\n", dir)

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
	if status, err := utils.GetReqXml("https://maven.minecraftforge.net/net/minecraftforge/forge/maven-metadata.xml", &versions); err != nil {
		return nil, fmt.Errorf("failed to fetch Forge installer versions from API (HTTP %d): %w", status, err)
	}

	versionMap := make(map[string][]string)

	for _, version := range versions.Versioning.Versions {
		parts := strings.Split(version, "-")
		if len(parts) >= 2 {
			v := parts[0]
			b := parts[1]
			versionMap[v] = append(versionMap[v], b)
		}
	}

	for v := range versionMap {
		slices.Reverse(versionMap[v])
	}

	return versionMap, nil
}
