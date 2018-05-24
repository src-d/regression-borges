package borges

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/src-d/regression-core.v0"
)

const defaultBorgesURL = "https://github.com/src-d/borges"

func NewBorges(
	config regression.Config,
	version string,
	releases *regression.Releases,
) (*regression.Binary, error) {
	url := config.GitURL
	if url == "" {
		url = defaultBorgesURL
	}

	path, err := gitURLToProjectPath(url)
	if err != nil {
		return nil, err
	}

	tool := regression.Tool{
		Name:        "borges",
		GitURL:      url,
		ProjectPath: path,
	}

	return regression.NewBinary(config, tool, version, releases), nil
}

var sshRegex = regexp.MustCompile(`([a-zA-Z0-9_-]+)@([^:]+):(.+)`)

func gitURLToProjectPath(url string) (string, error) {
	url = strings.ToLower(url)
	if strings.HasPrefix(url, "git://") {
		url = url[6:]
	} else if strings.HasPrefix(url, "http://") {
		url = url[7:]
	} else if strings.HasPrefix(url, "https://") {
		url = url[8:]
	} else if idx := strings.Index(url, "//"); idx >= 0 && idx+2 < len(url) {
		url = url[idx+2:]
	} else if sshRegex.MatchString(url) {
		matches := sshRegex.FindStringSubmatch(url)
		url = filepath.Join(matches[2], matches[3])
	} else {
		return "", fmt.Errorf("unable to convert URL %q to path", url)
	}

	if strings.HasSuffix(url, ".git") {
		url = url[:len(url)-4]
	}

	return filepath.FromSlash(url), nil
}
