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
	Dest          string
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	fmt.Fprintf(os.Stdout, "\rDownloading: %02d%%", int64(wc.Total*100)/wc.ContentLength)
	if wc.Total == uint64(wc.ContentLength) {
		fmt.Fprintf(os.Stdout, "\nSaved to %s\n", wc.Dest)
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

	counter := &WriteCounter{ContentLength: resp.ContentLength, Dest: path}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		return errors.New("failed to copy file to fs")
	}

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
