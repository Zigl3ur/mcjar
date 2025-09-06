package neoforge

import (
	"errors"
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
			return fmt.Errorf("neoforge doesnt support this version (given: %s)", version)
		}
	} else if snapshots {
		if len(versionsMap["snapshots"]) > 0 {
			for _, s := range versionsMap["snapshots"] {
				fmt.Printf("- %s\n", s)
			}
		} else {
			return errors.New("neoforge doesn't support snapshots")
		}
	} else {
		for _, v := range versionsMap["versions"] {
			fmt.Printf("- %s\n", v)
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
	cmd := exec.Command(java, "-jar", outPath, "--install-server", dir)

	if isVerbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		loader.Start("Installing neoforge server")
	}

	if err = cmd.Run(); err != nil {
		loader.Stop()
		return err
	}

	loader.Stop()
	fmt.Printf("Installed neoforge server at %s\n", dir)

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
	if status, err := utils.GetReqJson("https://maven.neoforged.net/api/maven/versions/releases/net/neoforged/neoforge", &list); err != nil {
		return nil, fmt.Errorf("[%d] Failed to get neoforge versions list", status)
	}

	versionMap := make(map[string][]string)

	for _, version := range list.Versions {
		// remove april fools versions
		if !strings.HasPrefix(version, "0") {
			parts := strings.SplitN(version, ".", 3)
			v := fmt.Sprintf("1.%s", parts[0])
			build := strings.Join(parts, ".")
			if len(parts) > 1 {
				v = fmt.Sprintf("1.%s.%s", parts[0], parts[1])
			}
			versionMap[v] = append(versionMap[v], build)
		}
	}

	for k := range versionMap {
		slices.Reverse(versionMap[k])
	}

	return versionMap, nil
}
