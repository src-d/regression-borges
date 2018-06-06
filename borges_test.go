package borges

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBorges(t *testing.T) {
	require := require.New(t)

	borges := NewToolBorges()
	require.NotNil(borges)

	require.Equal("borges", borges.Name)
	require.NotEmpty(borges.GitURL)
	require.NotEmpty(borges.ProjectPath)
	require.NotNil(borges.BuildSteps)
}
