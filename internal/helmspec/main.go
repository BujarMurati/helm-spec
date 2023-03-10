package helmspec

type TestSuiteResult struct {
	Succeeded   bool         `json:"succeeded"`
	SpecResults []SpecResult `json:"specResults"`
}

type TestRunner interface {
	Run(specFiles []string) (TestSuiteResult, error)
}

type HelmTestRunner struct{}

func (runner HelmTestRunner) Run(specFiles []string) (result TestSuiteResult, err error) {
	specs := []*HelmSpec{}
	for _, f := range specFiles {
		spec, err := NewSpec(f)
		if err != nil {
			return result, err
		}
		specs = append(specs, spec)
	}
	result = TestSuiteResult{
		Succeeded: true,
	}
	for _, spec := range specs {
		r := spec.Execute()
		result.Succeeded = result.Succeeded && r.Succeeded
		result.SpecResults = append(result.SpecResults, r)
	}
	return result, err
}
