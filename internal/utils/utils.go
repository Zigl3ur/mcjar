package utils

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
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
