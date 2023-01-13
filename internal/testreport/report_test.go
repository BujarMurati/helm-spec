package testreport

import (
	"testing"

	"github.com/bujarmurati/helm-spec/internal/helmspec"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

func TestTestReporterOutputModeYaml(t *testing.T) {
	spec, err := helmspec.NewSpec("../helmspec/testdata/charts/example/specs/successful_spec.yaml")
	assert.NoError(t, err)
	result := spec.Execute()
	testSuiteResult := helmspec.TestSuiteResult{
		Succeeded:   true,
		SpecResults: []helmspec.SpecResult{result},
	}
	reporter := HelmTestReporter{Result: testSuiteResult}
	output, err := reporter.Report("yaml")
	assert.NoError(t, err)
	reportedResult := &helmspec.TestSuiteResult{}
	yaml.Unmarshal([]byte(output), reportedResult)
	assert.Equal(t, testSuiteResult.Succeeded, reportedResult.Succeeded)
	assert.Equal(t, len(testSuiteResult.SpecResults), len(reportedResult.SpecResults))
}

func TestReporterOutputModePretty(t *testing.T) {
	spec, err := helmspec.NewSpec("../helmspec/testdata/charts/example/specs/successful_spec.yaml")
	assert.NoError(t, err)
	result := spec.Execute()
	testSuiteResult := helmspec.TestSuiteResult{
		Succeeded:   true,
		SpecResults: []helmspec.SpecResult{result},
	}
	reporter := HelmTestReporter{Result: testSuiteResult}
	_, err = reporter.Report("pretty")
	assert.NoError(t, err)
}

func TestTreeIndent(t *testing.T) {
	content := "foo"
	output := treeIndent(content, 2, []int{0, 1}, "|")
	assert.Equal(t, "|   |       foo", output)
}
