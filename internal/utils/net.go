package utils

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Zigl3ur/mcli/internal/utils/loader"
)

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

func WriteToFs(url, outPath string) error {

	dir, filename := filepath.Split(outPath)

	// TODO FIX, if given dir is . its considered as an empty string
	// if dir == "" {
	// 	dir, _ = os.Getwd()
	// }

	stat, err := os.Stat(dir)

	if os.IsNotExist(err) || !stat.IsDir() {
		if AskConfirm(fmt.Sprintf("Directory not found. Create \"%s\" ?", dir)) {
			// todo: mkdirall return nil if dir already exist maybe remove stats
			if err = os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		} else {
			return nil
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer resp.Body.Close()

	file, err := os.Create(dir + filename)
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
