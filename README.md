# regression-borges

**regression-borges** is a tool than runs different versions of `borges` and compares its resource consumption.

```
Usage:
  regression [OPTIONS]

Borges regression tester.

This tool executes borges pack with several versions and compares times and
resource usage. There should be at least two versions specified as arguments in
the following way:

* v0.12.1 - release name from github (https://github.com/src-d/borges/releases).
The binary will be downloaded.
* latest - latest release from github. The binary will be downloaded.
* remote:master - any tag or branch from borges repository. The binary will be
built automatically.
* local:fix/some-bug - tag or branch from the repository in the current directory.
The binary will be built.
* local:HEAD - current state of the repository. Binary is built.
* pull:266 - code from pull request #266 from borges repo. Binary is built.
* /path/to/borges - a borges binary built locally.

The repositories and downloaded/built borges binaries are cached by default in
"repos" and "binaries" repositories from the current directory.


Application Options:
      --binaries=   Directory to store binaries (default: binaries) [$REG_BINARIES]
      --repos=      Directory to store repositories (default: repos) [$REG_REPOS]
      --url=        URL to the tool repo [$REG_GITURL]
      --gitport=    Port for local git server (default: 9418) [$REG_GITPORT]
      --repos-file= YAML file with the list of repos [$REG_REPOS_FILE]
  -c, --complexity= Complexity of the repositories to test (default: 1)
                    [$REG_COMPLEXITY]
  -n, --repeat=     Number of times a test is run (default: 3) [$REG_REPEAT]
      --show-repos  List available repositories to test
  -t, --token=      Token used to connect to the API [$REG_TOKEN]
      --csv         save csv files with last result

Help Options:
  -h, --help        Show this help message
```

## Repos file example

```yaml
---
- name:        cangallo
  url:         git://github.com/jfontan/cangallo.git
  description: Small repository that should be fast to clone
  complexity:  0

- name:        octoprint-tft
  url:         https://github.com/mcuadros/OctoPrint-TFT
  description: Small repository that should be fast to clone
  complexity:  0

- name:        upsilon
  url:         git://github.com/upsilonproject/upsilon-common.git
  description: Average repository
  complexity:  1

- name:        numpy
  url:         git://github.com/numpy/numpy.git
  description: Average repository
  complexity:  2

- name:        tensorflow
  url:         git://github.com/tensorflow/tensorflow.git
  description: Average repository
  complexity:  3

- name:        bismuth
  url:         git://github.com/hclivess/Bismuth.git
  description: Big files repo (100Mb)
  complexity:  4
```

## License

Licensed under the terms of the Apache License Version 2.0. See the `LICENSE`
file for the full license text.

