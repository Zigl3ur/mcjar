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

	if dir == "" {
		dir, _ = os.Getwd()
	} else {
		if err := CheckDir(dir); err != nil {
			return err
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer resp.Body.Close()

	file, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer file.Close()

	loader.Start("Download starting")

	counter := &WriteCounter{StartTime: time.Now(), ContentLength: resp.ContentLength, filename: filename}
	if _, err = io.Copy(file, io.TeeReader(resp.Body, counter)); err != nil {
		return err
	}

	loader.Stop()
	fmt.Printf("Saved %s in %s\n", filename, dir)

	return nil
}

func GetReqJson(url string, dataJson any) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, errors.New(resp.Status)
	}

	//nolint:errcheck
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&dataJson); err != nil {
		return resp.StatusCode, err
	}

	return resp.StatusCode, nil
}

func GetReqXml(url string, dataXml any) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, errors.New(resp.Status)
	}

	//nolint:errcheck
	defer resp.Body.Close()

	if err := xml.NewDecoder(resp.Body).Decode(&dataXml); err != nil {
		return resp.StatusCode, err
	}

	return resp.StatusCode, nil
}
