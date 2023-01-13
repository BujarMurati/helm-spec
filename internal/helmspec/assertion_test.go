package helmspec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssertionEvaluation(t *testing.T) {
	type testCase struct {
		title          string
		query          string
		expectedResult string
		actualResult   string
		shouldError    bool
	}

	testCases := []testCase{
		{
			title:          "happy path",
			query:          ".foo",
			expectedResult: "bar",
			actualResult:   "bar",
		},
		{
			title:          "regular failure",
			query:          ".bar",
			expectedResult: "foo",
			actualResult:   "null",
		},
		{
			title:          "invalid query",
			query:          "not a valid yq query",
			expectedResult: "irrelevant",
			shouldError:    true,
		},
	}

	manifest := `
foo: bar
`

	for _, c := range testCases {
		t.Run(c.title, func(t *testing.T) {
			a := Assertion{
				Description:    "description",
				Query:          c.query,
				ExpectedResult: c.expectedResult,
			}
			result := a.Evaluate(manifest)
			if c.shouldError {
				assert.Error(t, result.Error)
			} else {
				assert.Equal(t, c.actualResult, result.ActualResult)
				assert.Equal(t, c.actualResult == c.expectedResult, result.Succeeded)
				assert.NoError(t, result.Error)
			}
		})
	}
}
