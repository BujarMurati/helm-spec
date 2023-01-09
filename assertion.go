package helmspec

import (
	"strings"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	logging "gopkg.in/op/go-logging.v1"
)

// checks the output of a `yq` query against rendered manifests
type Assertion struct {
	// human-readable description of what the assertion tests
	Description string
	// a [yq] query to perform against the rendering output
	// The output will contain all rendered manifests with document separators
	// [yq]: https://mikefarah.gitbook.io/yq/
	Query string
	// a string that the output of the `yq` query must equal in order for the test to pass
	ExpectedResult string
}

func EvalYQ(expression string, input string) (result string, err error) {
	evaluator := yqlib.NewStringEvaluator()
	encoder := yqlib.NewYamlEncoder(2, false, yqlib.NewDefaultYamlPreferences())
	decoder := yqlib.NewYamlDecoder(yqlib.NewDefaultYamlPreferences())
	return evaluator.Evaluate(expression, input, encoder, decoder)
}

type AssertionResult struct {
	Assertion    Assertion `json:"assertion"`
	Succeeded    bool      `json:"succeeded"`
	ActualResult string    `json:"actualResult"`
	Error        error     `json:"error"`
}

func (a Assertion) Evaluate(manifest string) (result AssertionResult) {
	actualResult, err := EvalYQ(a.Query, manifest)
	result.ActualResult = strings.TrimSpace(actualResult)
	result.Error = err
	result.Assertion = a
	result.Succeeded = result.ActualResult == a.ExpectedResult
	return result
}

func init() {
	// quiet down the noisy yq logger
	logger := yqlib.GetLogger()
	logging.SetLevel(logging.ERROR, logger.Module)
}
