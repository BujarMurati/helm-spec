package helmspec

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteRenderInstructions(t *testing.T) {
	spec, err := NewSpec("./testdata/charts/example/specs/example_spec.yaml")
	assert.NoError(t, err)
	actualManifest, err := spec.TestCases[0].Render.Execute("./testdata/charts/example")
	assert.NoError(t, err)
	expectedManifest, err := os.ReadFile("./testdata/manifest.yaml")
	assert.NoError(t, err)
	assert.Equal(t, string(expectedManifest), actualManifest)
}
