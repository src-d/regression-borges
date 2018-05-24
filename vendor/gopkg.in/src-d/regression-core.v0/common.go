package regression

import "io/ioutil"

func createTempDir() (string, error) {
	dir, err := ioutil.TempDir("", "regression-")
	if err != nil {
		return "", err
	}

	return dir, nil
}
