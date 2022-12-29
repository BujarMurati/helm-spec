package helmspec

const OUTPUT_MODE_YAML = "yaml"

type TestSuiteResult struct {
	Succeeded   bool         `json:"succeeded"`
	SpecResults []SpecResult `json:"specResults"`
}

type TestReporter interface {
	Report(outputMode string) (string, error)
}

type HelmTestReporter struct {
	Result TestSuiteResult
}

func (r HelmTestReporter) Report(outputMode string) (output string, err error) {
	return "foo", nil
}

type TestRunner interface {
	Run(specFiles []string) (TestReporter, error)
}

type HelmTestRunner struct{}

func (runner HelmTestRunner) Run(specFiles []string) (rep TestReporter, err error) {
	specs := []*HelmSpec{}
	for _, f := range specFiles {
		spec, err := NewSpec(f)
		if err != nil {
			return nil, err
		}
		specs = append(specs, spec)
	}
	result := TestSuiteResult{
		Succeeded: true,
	}
	for _, spec := range specs {
		r := spec.Execute()
		result.Succeeded = result.Succeeded && r.Succeeded
		result.SpecResults = append(result.SpecResults, r)
	}
	rep = HelmTestReporter{
		Result: result,
	}
	return rep, err
}
