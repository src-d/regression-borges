package regression

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func createTempDir() (string, error) {
	dir, err := ioutil.TempDir("", "regression-")
	if err != nil {
		return "", err
	}

	return dir, nil
}

func recursiveCopy(src, dst string) error {
	stat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		err = os.MkdirAll(dst, 0700)
		if err != nil {
			return err
		}

		files, err := ioutil.ReadDir(src)
		if err != nil {
			return err
		}

		for _, file := range files {
			srcPath := filepath.Join(src, file.Name())
			dstPath := filepath.Join(dst, file.Name())

			err = recursiveCopy(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	} else {
		err = copyFile(src, dst, stat.Mode())
		if err != nil {
			return err
		}
	}

	return nil
}

func copyFile(source, destination string, mode os.FileMode) error {
	exist, err := fileExist(source)
	if err != nil {
		return err
	}
	if !exist {
		return ErrBinaryNotFound.New()
	}

	orig, err := os.Open(source)
	if err != nil {
		return err
	}

	dir := filepath.Dir(destination)
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		return err
	}

	dst, err := os.Create(destination)
	if err != nil {
		return err
	}
	dst.Chmod(mode)
	defer dst.Close()

	_, err = io.Copy(dst, orig)
	if err != nil {
		dst.Close()
		os.Remove(dst.Name())
		return err
	}

	return nil
}
