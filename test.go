package borges

import (
	"fmt"
	"math"
	"time"

	"gopkg.in/src-d/go-log.v1"
	"gopkg.in/src-d/regression-core.v0"
)

type packResults map[string][]*PackResult
type versionResults map[string]packResults

type Test struct {
	config  regression.Config
	repos   *regression.Repositories
	server  *regression.GitServer
	borges  map[string]*regression.Binary
	results versionResults
}

func NewTest(config regression.Config) (*Test, error) {
	repos, err := regression.NewRepositories(config)
	if err != nil {
		return nil, err
	}

	server := regression.NewGitServer(config)

	return &Test{
		config: config,
		repos:  repos,
		server: server,
	}, nil
}

func (t *Test) Prepare() error {
	err := t.prepareServer()
	if err != nil {
		return err
	}

	err = t.prepareBorges()
	return err
}

func (t *Test) Stop() error {
	return t.server.Stop()
}

func (t *Test) Run() error {
	results := make(versionResults)

	for _, version := range t.config.Versions {
		_, ok := results[version]
		if !ok {
			results[version] = make(packResults)
		}

		borges, ok := t.borges[version]
		if !ok {
			panic("borges not initialized. Was Prepare called?")
		}

		log.With(log.Fields{"version": version}).Infof("Running version tests")

		times := t.config.Repeat
		if times < 1 {
			times = 1
		}

		for _, repo := range t.repos.Names() {
			results[version][repo] = make([]*PackResult, times)
			for i := 0; i < times; i++ {
				// TODO: do not stop on errors

				result, err := t.runTest(borges, repo)
				results[version][repo][i] = result

				if err != nil {
					return err
				}
			}
		}
	}

	t.results = results

	return nil
}

func (t *Test) GetResults() bool {
	if len(t.config.Versions) < 1 {
		panic("there should be at least one version")
	}

	ok := true
	if len(t.config.Versions) > 1 {
		ok = ok && t.CompareVersions()
	}

	return ok
}

func (t *Test) SaveLatestCSV() {
	version := t.config.Versions[len(t.config.Versions)-1]
	for _, repo := range t.repos.Names() {
		res := getAverageResult(t.results[version][repo])
		if err := res.SaveAllCSV(fmt.Sprintf("plot_%s_", repo)); err != nil {
			panic(err)
		}
	}
}

func (t *Test) CompareVersions() bool {
	versions := t.config.Versions
	ok := true
	for i, version := range versions[0 : len(versions)-1] {
		fmt.Printf("#### Comparing %s - %s ####\n", version, versions[i+1])
		a := t.results[versions[i]]
		b := t.results[versions[i+1]]

		for _, repo := range t.repos.Names() {
			fmt.Printf("## Repo %s ##\n", repo)

			repoA := getAverageResult(a[repo])
			repoB := getAverageResult(b[repo])

			c := repoA.ComparePrint(repoB, 10.0)
			ok = ok && c
		}
	}

	return ok
}

func (t *Test) runTest(
	borges *regression.Binary,
	repo string,
) (*PackResult, error) {
	url := t.server.Url(repo)
	log.Infof("Executing pack test for %s", repo)

	pack, err := NewPack(borges.Path, url)
	if err != nil {
		log.Errorf(err, "Could not execute pack")
		return nil, err
	}

	err = pack.Run()
	out, _ := pack.Out()
	if err != nil {
		log.With(log.Fields{
			"repo":   repo,
			"borges": borges.Path,
			"url":    url,
			"output": out}).Errorf(err, "Could not execute pack")
		return nil, err
	}

	var fileSize int64
	for _, f := range pack.files {
		fileSize += f.Size()
	}

	rusage, err := pack.Rusage()
	if err != nil {
		return nil, err
	}

	wall, err := pack.Wall()
	if err != nil {
		return nil, err
	}

	log.With(log.Fields{
		"wall":     wall,
		"memory":   rusage.Maxrss,
		"fileSize": fileSize,
	}).Infof("Finished pack")

	return pack.Result()
}

func (t *Test) prepareServer() error {
	log.Infof("Downloading repositories")
	err := t.repos.Download()
	if err != nil {
		return err
	}

	log.Infof("Starting git server")
	err = t.server.Start()
	return err
}

func (t *Test) prepareBorges() error {
	log.Infof("Preparing borges binaries")
	releases := regression.NewReleases("src-d", "borges", t.config.GitHubToken)

	t.borges = make(map[string]*regression.Binary, len(t.config.Versions))
	for _, version := range t.config.Versions {
		b := NewBorges(t.config, version, releases)
		err := b.Download()
		if err != nil {
			return err
		}

		t.borges[version] = b
	}

	return nil
}

func getAverageResult(rs []*PackResult) *PackResult {
	agg := &PackResult{}

	// Discard first for warmup
	rs = rs[1:]

	for _, r := range rs {
		agg.Memory += r.Memory
		agg.Wtime += r.Wtime
		agg.Stime += r.Stime
		agg.Utime += r.Utime
		agg.FileSize = r.FileSize
	}

	agg.Memory = int64(math.Round(float64(agg.Memory) / float64(len(rs))))
	agg.Wtime = time.Duration(math.Round(float64(agg.Wtime) / float64(len(rs))))
	agg.Stime = time.Duration(math.Round(float64(agg.Stime) / float64(len(rs))))
	agg.Utime = time.Duration(math.Round(float64(agg.Utime) / float64(len(rs))))
	agg.FileSize = int64(math.Round(float64(agg.FileSize) / float64(len(rs))))

	return agg
}
