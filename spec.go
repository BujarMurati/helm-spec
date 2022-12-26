package helmspec

import (
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"
)

// inputs for rendering a the chart with `helm template`
type RenderInstructions struct {
	// the release name to pass to `helm template`
	ReleaseName string `json:"releaseName"`
	// the release namespace to pass to `helm template`
	Namespace string `json:"namespace"`
	// all user-supplied values in one inline yaml document
	Values string `json´:"values"`
	// extra arguments passed through to the helm CLI, i.e. ["--set-file", "foo=foo.txt"]
	ExtraArgs []string `json:"extraArgs"`
}

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
