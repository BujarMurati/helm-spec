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

func TestExecuteTestDoesNotSucceedIfAnyAssertionFails(t *testing.T) {
	spec, err := NewSpec("./testdata/charts/example/specs/example_spec.yaml")
	assert.NoError(t, err)
	result := spec.TestCases[1].Execute(spec.ChartPath)
	assert.NoError(t, result.Error)
	assert.False(t, result.Succeeded)
}

func TestExecuteTestShouldAbortWhenRenderingFailsUnexpectedly(t *testing.T) {
	spec, err := NewSpec("./testdata/charts/example/specs/example_spec.yaml")
	assert.NoError(t, err)
	result := spec.TestCases[0].Execute("./not/an/existing/chart")
	assert.Error(t, result.Error)
	assert.False(t, result.Succeeded)
	// should not run or report assertions if we have an error
	// at the test case level
	assert.Equal(t, 0, len(result.AssertionResults))
}

func TestExecuteTestShouldSucceedOnExpectedFailure(t *testing.T) {
	spec, err := NewSpec("./testdata/charts/example/specs/example_spec.yaml")
	assert.NoError(t, err)
	result := spec.TestCases[2].Execute("./not/an/existing/chart")
	assert.True(t, result.Succeeded)
	// should not run or report assertions if we have an error
	// at the test case level
	assert.Equal(t, 0, len(result.AssertionResults))
}

func TestSpecResultShouldNotSucceedIfAnyTestCaseFails(t *testing.T) {
	spec, err := NewSpec("./testdata/charts/example/specs/example_spec.yaml")
	assert.NoError(t, err)
	result := spec.Execute()
	assert.False(t, result.Succeeded)
}

func TestSpecResultShouldSucceedIfAllTestCasesSucceed(t *testing.T) {
	spec, err := NewSpec("./testdata/charts/example/specs/successful_spec.yaml")
	assert.NoError(t, err)
	result := spec.Execute()
	assert.True(t, result.Succeeded)
}
