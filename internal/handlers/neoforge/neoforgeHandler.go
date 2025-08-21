package neoforge

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
			log.Fatalf("neoforge doesnt support this version (given: %s)", version)
		}
	} else if snapshots {
		if len(versionsMap["snapshots"]) > 0 {
			for _, s := range versionsMap["snapshots"] {
				fmt.Printf("- %s\n", s)
			}
		} else {
			log.Fatal("neoforge doesn't support snapshots")
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

	java, err := utils.GetJava()
	if err != nil {
		return err
	}

	destElt := strings.Split(path, "/")
	dest := strings.Join(destElt[:len(destElt)-1], "/")
	cmd := exec.Command(java, "-jar", path, "--install-server", dest)
	loader.Start("Installing neoforge server")
	// use cmd.Output ? if adding a debug flag and print output ?
	if err = cmd.Run(); err != nil {
		loader.Stop()
		return errors.New("failed to install neoforge server")
	}

	loader.Stop()
	fmt.Printf("Installed neoforge server at %s\n", dest)

	return nil
}

func getUrl(version, build string) (string, error) {
	vlist, err := getVersionsList()
	if err != nil {
		return "", err
	}

	if vlist[version] == nil {
		return "", fmt.Errorf("no neoforge jar available for provided version (given: %s)", version)
	}

	url := fmt.Sprintf("https://maven.neoforged.net/releases/net/neoforged/neoforge/%s/neoforge-%s-installer.jar", vlist[version][0], vlist[version][0])

	if build != "latest" {
		if slices.Contains(vlist[version], build) {
			url = fmt.Sprintf("https://maven.neoforged.net/releases/net/neoforged/neoforge/%s/neoforge-%s-installer.jar", build, build)
		} else {
			return "", fmt.Errorf("no neoforge jar available for provided version / neoforge version (given: %s, %s)", version, build)
		}
	}

	return url, nil
}

func getVersionsList() (map[string][]string, error) {

	type NeoforgeVersions struct {
		Versions []string `json:"versions"`
	}

	var list NeoforgeVersions
	if err := utils.GetReqJson("https://maven.neoforged.net/api/maven/versions/releases/net/neoforged/neoforge", &list); err != nil {
		return nil, errors.New("failed to fetch neoforge versions")
	}

	versionMap := make(map[string][]string)

	for _, v := range list.Versions {
		// remove april fools versions
		if !strings.HasPrefix(v, "0") {
			parts := strings.SplitN(v, ".", 3)
			version := fmt.Sprintf("1.%s", parts[0])
			build := strings.Join(parts, ".")
			if len(parts) > 1 {
				version = fmt.Sprintf("1.%s.%s", parts[0], parts[1])
			}
			versionMap[version] = append(versionMap[version], build)
		}
	}

	for k := range versionMap {
		slices.Reverse(versionMap[k])
	}

	return versionMap, nil
}
