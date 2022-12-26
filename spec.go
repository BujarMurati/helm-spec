package helmspec

import (
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"
)

// a testcase bundles rendering instructions with a list of assertions
// to perform against the rendered output
type TestCase struct {
	// title of the testcase
	Title string `json:"title"`
	// inputs for rendering a helm chart with `helm template`
	Render RenderInstructions `json:"render"`
	// assertions against the rendering output
	Assertions []Assertion `json:"assertions"`
}

type TestCaseResult struct {
	Title            string             `json:"title"`
	Render           RenderInstructions `json:"render"`
	Succeeded        bool               `json:"succeeded"`
	Manifest         string             `json:"manifest"`
	AssertionResults []AssertionResult  `json:"assertionResults"`
	Error            error              `json:"error"`
}

// renders a chart based on render instructions and evaluates assertions
func (t TestCase) Execute(chartPath string) (result TestCaseResult) {
	result.Title = t.Title
	result.Render = t.Render
	result.Manifest, result.Error = t.Render.Execute(chartPath)
	if t.Render.ShouldFailToRender {
		result.Succeeded = result.Error != nil
		return result
	}
	result.Succeeded = true
	for _, assertion := range t.Assertions {
		r := assertion.Evaluate(result.Manifest)
		result.AssertionResults = append(result.AssertionResults, r)
		if !r.Succeeded {
			result.Succeeded = false
		}
	}
	return result
}

// a related group of test cases for the same helm chart
type HelmSpec struct {
	// title
	Title string `json:"title"`
	// path to the helm chart (absolute or relative to the spec file directory)
	ChartPath string `json:"chartPath"`
	// test cases to run for the helm chart
	TestCases []TestCase `json:"testCases"`
}

func NewSpec(filePath string) (spec *HelmSpec, err error) {
	spec = new(HelmSpec)
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return spec, err
	}
	content, err := os.ReadFile(absFilePath)
	if err != nil {
		return spec, err
	}
	err = yaml.Unmarshal(content, spec)
	if err != nil {
		return spec, err
	}
	if !filepath.IsAbs(spec.ChartPath) {
		spec.ChartPath = filepath.Join(filepath.Dir(absFilePath), spec.ChartPath)
	}
	return spec, err
}
