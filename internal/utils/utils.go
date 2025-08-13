package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type WriteCounter struct {
	Total         uint64
	ContentLength int64
}

var loadingGlyphs = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)

	spinnerIndex := (wc.Total / 1024000) % uint64(len(loadingGlyphs))

	if wc.ContentLength > 0 {
		progress := int64(wc.Total*100) / wc.ContentLength
		fmt.Printf("\r%s Downloading: %02d%%", Loading(spinnerIndex), progress)
	} else {
		fmt.Printf("\r%s Downloading: %d bytes", Loading(spinnerIndex), wc.Total)
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

	counter := &WriteCounter{ContentLength: resp.ContentLength}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		return errors.New("failed to copy file to fs")
	}

	ClearLine()
	fmt.Printf("Saved to %s\n", path)

	return nil
}

func GetReq(url string, dataJson any) error {
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

func Loading(index uint64) string {
	return loadingGlyphs[index%uint64(len(loadingGlyphs))]
}

func ClearLine() {
	fmt.Print("\r\033[K")
}
