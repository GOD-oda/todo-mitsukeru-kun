package github

import (
	"bytes"
	"net/http"
)

const (
	URL = "https://api.github.com/"
)

type Header struct {
	ContentType       string
	Authorization     string
	XGitHubApiVersion string
}

func MakeHeader(token string) Header {
	header := Header{
		ContentType:       "application/vnd.github+json",
		XGitHubApiVersion: "2022-11-28",
	}
	if token != "" {
		header.Authorization = "Bearer " + token
	}

	return header
}

func Get(path string, body string, header Header) (*http.Response, error) {
	req, err := http.NewRequest("GET", URL+path, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", header.ContentType)
	if header.Authorization != "" {
		req.Header.Set("Authorization", header.Authorization)
	}
	req.Header.Set("X-GitHub-Api-Version", header.XGitHubApiVersion)

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, err
}
