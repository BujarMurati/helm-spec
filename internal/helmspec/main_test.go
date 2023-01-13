package helmspec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelmTestRunner(t *testing.T) {
	specFiles := []string{"./testdata/charts/example/specs/example_spec.yaml", "./testdata/charts/example/specs/successful_spec.yaml"}
	result, err := HelmTestRunner{}.Run(specFiles)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestHelmTestRunnerAbortsIfItFailsToLoadAnySpec(t *testing.T) {
	specFiles := []string{"./testdata/charts/example/specs/example_spec.yaml", "./testdata/charts/example/specs/does_not_exist.yaml"}
	_, err := HelmTestRunner{}.Run(specFiles)
	assert.Error(t, err)
}
