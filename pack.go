package borges

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/src-d/regression-core.v0"
)

type Pack struct {
	*regression.Executor
	test   bool
	binary string
	repo   string
	files  []os.FileInfo
}

func NewPack(binary, repo string) (*Pack, error) {
	return &Pack{
		Executor: new(regression.Executor),
		binary:   binary,
		repo:     repo,
	}, nil
}

func (p *Pack) Run() error {
	list, err := createList(p.repo)
	if err != nil {
		return err
	}
	defer os.Remove(list)

	dir, err := createTempDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	lArg := fmt.Sprintf("%s", list)
	dArg := fmt.Sprintf("--root-repositories-dir=%s", dir)
	tArg := fmt.Sprintf("--timeout=4h")

	executor, err := regression.NewExecutor(p.binary, "pack", dArg, tArg, lArg)
	if err != nil {
		return err
	}

	p.Executor = executor

	err = p.Executor.Run()
	if err != nil {
		return err
	}

	files, err := fileInfo(dir)
	if err != nil {
		return err
	}

	p.files = files
	p.test = true

	return nil
}

func (p *Pack) Files() ([]os.FileInfo, error) {
	if !p.Executed {
		return nil, regression.ErrNotRun
	}

	return p.files, nil
}

func (p *Pack) Result() (*PackResult, error) {
	var size int64

	for _, file := range p.files {
		size += file.Size()
	}

	rusage, err := p.Rusage()
	if err != nil {
		return nil, err
	}

	wall, err := p.Wall()
	if err != nil {
		return nil, err
	}

	packResult := &PackResult{
		Result: &regression.Result{
			Memory: rusage.Maxrss * 1024,
			Wtime:  wall,
			Stime:  time.Duration(rusage.Stime.Nano()),
			Utime:  time.Duration(rusage.Utime.Nano()),
		},
		Files:    p.files,
		FileSize: size,
	}

	return packResult, nil
}

type PackResult struct {
	*regression.Result
	Files    []os.FileInfo
	FileSize int64 // bytes
}

func NewPackResult() *PackResult {
	return &PackResult{Result: new(regression.Result)}
}

type PackComparison struct {
	FileSize float64
}

const FileSize = "file_size"

func (p *PackResult) Compare(q *PackResult) PackComparison {
	return PackComparison{
		FileSize: percent(p.FileSize, q.FileSize),
	}
}

func (p *PackResult) ComparePrint(q *PackResult, allowance float64) bool {
	ok := p.Result.ComparePrint(q.Result, allowance)
	c := p.Compare(q)

	if c.FileSize > allowance {
		ok = false
	}
	fmt.Printf(regression.CompareFormat,
		"FileSize",
		regression.ToMiB(p.FileSize),
		regression.ToMiB(q.FileSize),
		c.FileSize,
		allowance > c.FileSize,
	)

	return ok
}

func createList(repo string) (string, error) {
	tmpFile, err := ioutil.TempFile("", "packer-list")
	if err != nil {
		return "", err
	}

	_, err = tmpFile.WriteString(repo)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", err
	}

	err = tmpFile.Close()
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	return tmpFile.Name(), nil
}

func createTempDir() (string, error) {
	dir, err := ioutil.TempDir("", "packer-dir")
	if err != nil {
		return "", err
	}

	return dir, nil
}

func fileInfo(dir string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dir)
}

func percent(a, b int64) float64 {
	diff := b - a
	return (float64(diff) / float64(a)) * 100
}
