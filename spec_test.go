package helmspec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadHelmSpec(t *testing.T) {
	spec, err := NewSpec("./testdata/charts/example/specs/example_spec.yaml")
	assert.NoError(t, err)
	assert.Equal(t, "template tests for the `example` helm chart", spec.Title)
}

func TestChartPathIsRelativeToSpecFile(t *testing.T) {
	expectedChartPath, err := filepath.Abs("./testdata/charts/example")
	assert.NoError(t, err)
	spec, err := NewSpec("./testdata/charts/example/specs/example_spec.yaml")
	assert.NoError(t, err)
	assert.Equal(t, expectedChartPath, spec.ChartPath)
}

func TestExecuteTestCaseHappyPath(t *testing.T) {
	spec, err := NewSpec("./testdata/charts/example/specs/example_spec.yaml")
	assert.NoError(t, err)
	result := spec.TestCases[0].Execute(spec.ChartPath)
	assert.NoError(t, result.Error)
	assert.Equal(t, spec.TestCases[0].Title, result.Title)
	assert.Equal(t, spec.TestCases[0].Render, result.Render)
	assert.True(t, result.Succeeded)
	expectedManifest, err := os.ReadFile("./testdata/example_spec_0_manifest.yaml")
	assert.NoError(t, err)
	assert.Equal(t, string(expectedManifest), result.Manifest)
	assert.Equal(t, 1, len(result.AssertionResults))
	assert.True(t, result.AssertionResults[0].Succeeded)
}
