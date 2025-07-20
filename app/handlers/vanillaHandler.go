package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

func GetUrlVanilla(version string) (string, error) {
	resp, err := http.Get("https://launchermeta.mojang.com/mc/game/version_manifest.json")

	if err != nil {
		return "", errors.New("Failed to get version manifest")
	}

	defer resp.Body.Close()

	type Versions struct {
		Versions []struct {
			Id  string `json:"id"`
			Url string `json:"url"`
		}
	}

	var versions Versions
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return "", errors.New("Failed to decode versions list")
	}

	var versionUrl string
	for _, v := range versions.Versions {
		if v.Id == version {
			versionUrl = v.Url
			break
		}
	}

	if versionUrl == "" {
		return "", errors.New("Specified version not found")
	}

	resp, err = http.Get(versionUrl)

	if err != nil {
		return "", errors.New("Failed to get download url")
	}

	defer resp.Body.Close()

	type DownloadData struct {
		Downloads struct {
			Server struct {
				Url string `json:"url"`
			} `json:"server"`
		} `json:"downloads"`
	}

	var downloadData DownloadData
	if err := json.NewDecoder(resp.Body).Decode(&downloadData); err != nil {
		return "", errors.New("Failed to decode download url")
	}

	return downloadData.Downloads.Server.Url, nil
}
