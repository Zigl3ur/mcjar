package utils

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
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
		loader.UpdateMessage(fmt.Sprintf("Downloading: %02d%%, %s", progress, downloadSpeed))
	} else {
		loader.UpdateMessage(fmt.Sprintf("Downloading: %d bytes, %s", wc.Total, downloadSpeed))
	}

	return n, nil
}

func WriteToFs(url, path string) error {
	resp, err := http.Get(url)

	if err != nil {
		return errors.New("failed to download jar")
	}

	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		return errors.New("failed to create file")
	}

	defer out.Close()

	loader.Start("Download starting")

	counter := &WriteCounter{StartTime: time.Now(), ContentLength: resp.ContentLength}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		return errors.New("failed to copy file to fs")
	}

	loader.Stop()
	fmt.Printf("Saved to %s\n", path)

	return nil
}

func GetReqJson(url string, dataJson any) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

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

	defer resp.Body.Close()

	if err := xml.NewDecoder(resp.Body).Decode(&dataXml); err != nil {
		return err
	}

	return nil
}

func humanizeByte(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

// TODO: work for now, but try to improve + tests
func SortMcVersions(versions []string) []string {
	sortedVersion := make([]string, len(versions))
	copy(sortedVersion, versions)

	for i := 0; i < len(versions)-1; i++ {
		for j := range sortedVersion {
			sortedVer, err := mcVersionParser(sortedVersion[j])
			if err != nil {
				continue
			}
			unsortedVer, err := mcVersionParser(sortedVersion[i+1])
			if err != nil {
				continue
			}

			for idx := range 3 {
				if unsortedVer[idx] < sortedVer[idx] {
					sortedVersion[j], sortedVersion[i+1] = sortedVersion[i+1], sortedVersion[j]
					break
				} else if unsortedVer[idx] > sortedVer[idx] {
					break
				}
			}
		}
	}

	return sortedVersion
}

func mcVersionParser(version string) ([3]int, error) {
	parts := strings.SplitN(version, ".", 3)

	var mainVersion, subVersion, patch int

	if len(parts) >= 1 {
		mainVersion, _ = strconv.Atoi(parts[0])
	}
	if len(parts) >= 2 {
		subVersion, _ = strconv.Atoi(parts[1])
	}
	if len(parts) >= 3 {
		patch, _ = strconv.Atoi(parts[2])
	} else {
		return [3]int{}, errors.New("not a correct version format") // if it's a snaphsot or like april fool ver ot whatever that is not following the format "1.12.2"
	}

	return [3]int{mainVersion, subVersion, patch}, nil
}
