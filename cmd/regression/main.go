package main

import (
	"os"

	"github.com/src-d/regression-borges"

	"github.com/jessevdk/go-flags"
	"gopkg.in/src-d/go-log.v1"
	"gopkg.in/src-d/regression-core.v0"
)

var description = `Borges regression tester.

This tool executes borges pack with several versions and compares times and resource usage. There should be at least two versions specified as arguments in the following way:

* v0.12.1 - release name from github (https://github.com/src-d/borges/releases). The binary will be downloaded.
* latest - latest release from github. The binary will be downloaded.
* remote:master - any tag or branch from borges repository. The binary will be built automatically.
* local:fix/some-bug - tag or branch from the repository in the current directory. The binary will be built.
* pull:266 - code from pull request #266 from borges repo. Binary is built.
* /path/to/borges - a borges binary built locally.

The repositories and downloaded/built borges binaries are cached by default in "repos" and "binaries" repositories from the current directory.
`

type Options struct {
	regression.Config

	CSV bool `long:"csv" description:"save csv files with last result"`
}

func main() {
	options := Options{
		Config: regression.NewConfig(),
	}

	parser := flags.NewParser(&options, flags.Default)
	parser.LongDescription = description

	args, err := parser.Parse()
	if err != nil {
		if err, ok := err.(*flags.Error); ok {
			if err.Type == flags.ErrHelp {
				os.Exit(0)
			}
		}

		log.Errorf(err, "Could not parse arguments")
		os.Exit(1)
	}

	config := options.Config

	if config.ShowRepos {
		repos, err := regression.NewRepositories(config)
		if err != nil {
			log.Errorf(err, "Can not initialize repositories")
			panic(err)
		}

		repos.ShowRepos()
		os.Exit(0)
	}

	if len(args) < 1 {
		log.Errorf(nil, "There should be at least one version")
		os.Exit(1)
	}

	config.Versions = args

	test, err := borges.NewTest(config)
	if err != nil {
		panic(err)
	}

	log.Infof("Preparing run")
	err = test.Prepare()
	if err != nil {
		log.Errorf(err, "Could not prepare environment")
		os.Exit(1)
	}

	err = test.Run()
	if err != nil {
		panic(err)
	}

	res := test.GetResults()
	if res && options.CSV {
		test.SaveLatestCSV()
	}

	err = test.Stop()
	if err != nil {
		panic(err)
	}

	if !res {
		os.Exit(1)
	}
}
