package helmspec

type TestReporter interface {
	Report() (string, error)
}

type TestRunner interface {
	Run(specFiles []string) (TestReporter, error)
}

type HelmTestRunner struct{}

func (runnner *HelmTestRunner) Run(specFiles []string) (rep TestReporter, err error) {
	return rep, err
}
