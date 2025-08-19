package utils

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Zigl3ur/mcli/internal/utils/loader"
)

const InvalidServerType string = "Invalid server type, valid ones are [vanilla, paper, spigot, purpur, forge, neoforge, fabric]"

// if its not following conventionnal release name like "1.12.2" (probably a snapshot or whatever april fools versions
var mcVersionParseError error = errors.New("failed to parse mc version")

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

func SortMcVersions(versions []string) []string {
	slices.SortStableFunc(versions, func(i, j string) int {
		ver0, err0 := mcVersionParser(i)
		ver1, err1 := mcVersionParser(j)

		if err0 != nil || err1 != nil {
			return 0
		}

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
	return versions
}

func mcVersionParser(version string) ([3]int, error) {
	parts := strings.SplitN(version, ".", 3)

	if len(parts) < 2 {
		return [3]int{}, mcVersionParseError
	}

	var mainVersion, subVersion, patch int

	mainVersion, err := strconv.Atoi(parts[0])
	if err != nil {
		return [3]int{}, mcVersionParseError
	}
	subVersion, err = strconv.Atoi(parts[1])
	if err != nil {
		return [3]int{}, mcVersionParseError
	}

	if len(parts) >= 3 {
		patch, err = strconv.Atoi(parts[2])
		if err != nil {
			return [3]int{}, mcVersionParseError
		}
	}

	return [3]int{mainVersion, subVersion, patch}, nil
}
