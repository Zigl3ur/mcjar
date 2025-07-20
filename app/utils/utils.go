package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type WriteCounter struct {
	Total         uint64
	ContentLength int64
	Dest          string
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	fmt.Fprintf(os.Stdout, "\rDownloaded: %d / %d bytes", wc.Total, wc.ContentLength)
	if wc.Total == uint64(wc.ContentLength) {
		fmt.Fprintf(os.Stdout, "\nSaved to %s\n", wc.Dest)
	}
	return n, nil
}

func WriteToFs(url, path string) error {

	resp, err := http.Get(url)

	if err != nil {
		return errors.New("Failed to download jar")
	}

	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		return errors.New("Failed to create file")
	}

	defer out.Close()

	counter := &WriteCounter{ContentLength: resp.ContentLength, Dest: path}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return errors.New("Failed to copy file to fs")
	}

	return nil
}
