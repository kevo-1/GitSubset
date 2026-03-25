package internal

import (
	"errors"
	"net/url"
	"strings"
)

type GithubLink struct {
	User string
	Repo string
	Path string
}

func ParseURL(URL string) (GithubLink, error) {
	if strings.HasPrefix(URL, "git@") {
		URL = strings.Replace(URL, "git@github.com:", "https://github.com/", 1)
	}

	parsedURL, err := url.Parse(URL)
	if err != nil {
		return GithubLink{}, err
	}

	if parsedURL.Host != "github.com" {
		return GithubLink{}, errors.New("only github URLs are supported")
	}

	content := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")

	if len(content) != 2 {
		return GithubLink{}, errors.New("invalid github repo url format")
	}

	repo := strings.TrimSuffix(content[1], ".git")

	return GithubLink{
		User: content[0],
		Repo: repo,
		Path: repo,
	}, nil
}