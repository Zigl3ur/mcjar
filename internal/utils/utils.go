package utils

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Zigl3ur/mcli/internal/utils/loader"
)

const InvalidServerType string = "Invalid server type, valid ones are [vanilla, paper, spigot, purpur, forge, neoforge, fabric]"

type WriteCounter struct {
	Total         uint64
	ContentLength int64
	StartTime     time.Time
	filename      string
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)

	var downloadSpeed string
	since := time.Since(wc.StartTime).Seconds()
	if since > 0 {
		downloadSpeed = fmt.Sprintf("%s/s", humanizeByte(int64((float64(wc.Total) / since))))
	}

	if wc.ContentLength > 0 {
		progress := int64(wc.Total*100) / wc.ContentLength
		loader.UpdateMessage(fmt.Sprintf("Downloading %s: %02d%%, %s", wc.filename, progress, downloadSpeed))
	} else {
		loader.UpdateMessage(fmt.Sprintf("Downloading %s: %d bytes, %s", wc.filename, wc.Total, downloadSpeed))
	}

	return n, nil
}

func WriteToFs(url, dir, filename string) error {

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if AskConfirm("Directory not found. Create it ?") {
			if err = os.Mkdir(dir, 0755); err != nil {
				return err
			}
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer resp.Body.Close()

	// dir will always end with "/" cause its checked in commands prerun func
	filepath := dir + filename
	file, err := os.Create(filepath)
	if err != nil {
		return errors.New("failed to create output file")
	}

	//nolint:errcheck
	defer file.Close()

	loader.Start("Download starting")

	counter := &WriteCounter{StartTime: time.Now(), ContentLength: resp.ContentLength, filename: filename}
	if _, err = io.Copy(file, io.TeeReader(resp.Body, counter)); err != nil {
		return errors.New("failed to create output file")
	}

	loader.Stop()
	fmt.Printf("Saved %s in %s\n", filename, dir)

	return nil
}

func GetReqJson(url string, dataJson any) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&dataJson); err != nil {
		return err
	}

	return nil
}

func GetReqXml(url string, dataXml any) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer resp.Body.Close()

	if err := xml.NewDecoder(resp.Body).Decode(&dataXml); err != nil {
		return err
	}

	return nil
}

func humanizeByte(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTP"[exp])
}

func SortMcVersions(versions []string) map[string][]string {
	versionsData := make(map[string][]string, 2)
	versionsData["versions"] = make([]string, 0)
	versionsData["snapshots"] = make([]string, 0)

	for _, version := range versions {
		_, unparsed := mcVersionParser(version)
		if unparsed != "" {
			versionsData["snapshots"] = append(versionsData["snapshots"], version)
		} else {
			versionsData["versions"] = append(versionsData["versions"], version)
		}
	}

	slices.SortStableFunc(versionsData["versions"], func(i, j string) int {
		ver0, _ := mcVersionParser(i)
		ver1, _ := mcVersionParser(j)

		for idx := range 3 {
			if ver0[idx] > ver1[idx] {
				return -1
			}
			if ver0[idx] < ver1[idx] {
				return 1
			}
		}

		return 0
	})

	return versionsData
}

func mcVersionParser(version string) ([3]int, string) {
	parts := strings.SplitN(version, ".", 3)

	if len(parts) < 2 {
		return [3]int{}, version
	}

	var mainVersion, subVersion, patch int

	mainVersion, err := strconv.Atoi(parts[0])
	if err != nil {
		return [3]int{}, version
	}
	subVersion, err = strconv.Atoi(parts[1])
	if err != nil {
		return [3]int{}, version
	}

	if len(parts) >= 3 {
		patch, err = strconv.Atoi(parts[2])
		if err != nil {
			return [3]int{}, version
		}
	}

	return [3]int{mainVersion, subVersion, patch}, ""
}

func GetPath(file string) (string, error) {
	path, err := exec.LookPath(file)
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return "", fmt.Errorf("%s not found in PATH, please install it and retry", file)
		} else {
			return "", err
		}
	}

	return path, nil
}

// facets builder create facets string with given data,
// facets is a query param for modrinth api
func FacetsBuilder(versions []string, loader, projectType string) string {
	elt := make([]string, 0, 3)

	if len(versions) > 0 {
		velt := make([]string, 0, len(versions))

		for _, v := range versions {
			velt = append(velt, fmt.Sprintf("\"versions:%s\"", v))
		}

		elt = append(elt, fmt.Sprintf("[%s]", strings.Join(velt, ",")))
	}

	if loader != "" {
		elt = append(elt, fmt.Sprintf("[\"categories:%s\"]", loader))
	}

	if projectType != "" {
		elt = append(elt, fmt.Sprintf("[\"project_type:%s\"]", projectType))
	}

	if len(elt) == 0 {
		return ""
	}

	return fmt.Sprintf("[%s]", strings.Join(elt, ","))
}

func Iso8601Format(date string) (string, error) {
	layout := "2006-01-02T15:04:05Z"
	t, err := time.Parse(layout, date)

	if err != nil {
		return "", err
	}

	return t.Format("Jan 2, 2006, 03:04 PM"), nil
}

func AskConfirm(message string) bool {
	var answer string

	fmt.Printf("%s [Y/n]: ", message)
	//nolint:errcheck
	fmt.Scanf("%s", &answer)

	if strings.ToUpper(answer) == "Y" {
		return true
	}

	fmt.Println("Aborted.")
	return false
}
