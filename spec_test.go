package helmspec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadHelmSpec(t *testing.T) {
	spec, err := NewSpec("./testdata/charts/example/specs/example_spec.yaml")
	assert.NoError(t, err)
	assert.Equal(t, "template tests for the `example` helm chart", spec.Title)
}
