package borges

import (
	"gopkg.in/src-d/regression-core.v0"
)

func NewToolBorges() regression.Tool {
	steps := []regression.BuildStep{
		{
			Dir:     "",
			Command: "make",
			Args:    []string{"packages"},
		},
	}

	return regression.Tool{
		Name:        "borges",
		GitURL:      "https://github.com/src-d/borges",
		ProjectPath: "github.com/src-d/borges",
		BuildSteps:  steps,
	}
}

func NewBorges(
	config regression.Config,
	version string,
	releases *regression.Releases,
) *regression.Binary {
	return regression.NewBinary(config, NewToolBorges(), version, releases)
}
