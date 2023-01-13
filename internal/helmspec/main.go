package helmspec

type TestSuiteResult struct {
	Succeeded   bool         `json:"succeeded"`
	SpecResults []SpecResult `json:"specResults"`
}

type TestReporter interface {
	Report(outputMode string) (string, error)
}

type TestRunner interface {
	Run(specFiles []string) (TestSuiteResult, error)
}

type HelmTestRunner struct{}

func (runner HelmTestRunner) Run(specFiles []string) (res TestSuiteResult, err error) {
	specs := []*HelmSpec{}
	for _, f := range specFiles {
		spec, err := NewSpec(f)
		if err != nil {
			return res, err
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
	return res, err
}
