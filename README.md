# regression-borges

**regression-borges** is a tool than runs different versions of `borges` and compares its resource consumption.

```
Usage:
  regression [OPTIONS]

Borges regression tester.

This tool executes borges pack with several versions and compares times and
resource usage. There should be at least two versions specified as arguments in
the following way:

* v0.12.1 - release name from github
(https://github.com/src-d/borges/releases). The binary will be downloaded.
* remote:master - any tag or branch from borges repository. The binary will be
built automatically.
* local:fix/some-bug - tag or branch from the repository in the current
directory. The binary will be built.
* pull:266 - code from pull request #266 from borges repo. Binary is built.
* /path/to/borges - a borges binary built locally.

The repositories and downloaded/built borges binaries are cached by default in
"repos" and "binaries" repositories from the current directory.


Application Options:
      --binaries=   Directory to store binaries (default: binaries)
                    [$REG_BINARIES]
      --repos=      Directory to store repositories (default: repos)
                    [$REG_REPOS]
      --url=        URL to the tool repo [$REG_GITURL]
      --gitport=    Port for local git server (default: 9418) [$REG_GITPORT]
  -c, --complexity= Complexity of the repositories to test (default: 1)
                    [$REG_COMPLEXITY]
  -n, --repeat=     Number of times a test is run (default: 3) [$REG_REPEAT]
      --show-repos  List available repositories to test

Help Options:
  -h, --help        Show this help message
```

## License

Licensed under the terms of the Apache License Version 2.0. See the `LICENSE`
file for the full license text.

