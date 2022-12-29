package helmspec

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

func TestHelmTestRunner(t *testing.T) {
	specFiles := []string{"./testdata/charts/example/specs/example_spec.yaml", "./testdata/charts/example/specs/successful_spec.yaml"}
	reporter, err := HelmTestRunner{}.Run(specFiles)
	assert.NoError(t, err)
	assert.NotNil(t, reporter)
}

func TestHelmTestRunnerAbortsIfItFailsToLoadAnySpec(t *testing.T) {
	specFiles := []string{"./testdata/charts/example/specs/example_spec.yaml", "./testdata/charts/example/specs/does_not_exist.yaml"}
	_, err := HelmTestRunner{}.Run(specFiles)
	assert.Error(t, err)
}

func TestTestReporterOutputModeYaml(t *testing.T) {
	spec, err := NewSpec("./testdata/charts/example/specs/successful_spec.yaml")
	assert.NoError(t, err)
	result := spec.Execute()
	testSuiteResult := TestSuiteResult{
		Succeeded:   true,
		SpecResults: []SpecResult{result},
	}
	reporter := HelmTestReporter{Result: testSuiteResult}
	output, err := reporter.Report("yaml")
	assert.NoError(t, err)
	reportedResult := &TestSuiteResult{}
	yaml.Unmarshal([]byte(output), reportedResult)
	assert.Equal(t, testSuiteResult.Succeeded, reportedResult.Succeeded)
	assert.Equal(t, len(testSuiteResult.SpecResults), len(reportedResult.SpecResults))
}
