package main

import (
	"path/filepath"
	"strings"
	"testing"

	helmspec "github.com/bujarmurati/helm-spec"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

type errCapture struct {
	Ctx *cli.Context
	Err error
}

type mockTestReporter struct{}

func (m mockTestReporter) Report(outputMode string) (string, error) {
	return "success!", nil
}

type mockTestRunner struct {
	SpecFiles []string
	HasRun    bool
}

func (m *mockTestRunner) Run(specFiles []string) (r helmspec.TestReporter, err error) {
	m.SpecFiles = specFiles
	m.HasRun = true
	return mockTestReporter{}, nil
}

// runs the CLI with args. Outputs and errors  are captured
func testRun(t *testing.T, args []string) (runner *mockTestRunner, out *strings.Builder, err error) {
	t.Helper()
	e := &errCapture{}
	runner = &mockTestRunner{}
	out = &strings.Builder{}
	settings := cliSettings{
		Reader:    strings.NewReader(""),
		Writer:    out,
		ErrWriter: out,
		ExitErrHandler: func(cCtx *cli.Context, err error) {
			e.Ctx = cCtx
			e.Err = err
		},
		TestRunner: runner,
	}
	app, err := createApp(settings)
	assert.NoError(t, err)
	return runner, out, app.Run(args)
}

// just a sanity check to fail fast if something is utterly broken
func TestSanity(t *testing.T) {
	args := []string{"helm-spec", "--help"}
	_, _, err := testRun(t, args)
	assert.NoError(t, err)
}

func TestAcceptsTestSuitePathArg(t *testing.T) {
	args := []string{"helm-spec", "./testdata/specs"}
	_, _, err := testRun(t, args)
	assert.NoError(t, err)
}

func TestValidatesTestArg(t *testing.T) {
	type testCase struct {
		title          string
		path           string
		valid          bool
		errMsgFragment string
	}
	testCases := []testCase{
		{
			title:          "non-existant path",
			path:           "./foo",
			errMsgFragment: "no such file or directory",
		},
		{
			title:          "not a directory",
			path:           "./main_test.go",
			errMsgFragment: "not a directory",
		},
		{
			title:          "directory without *_spec.yaml files",
			path:           "./testdata",
			errMsgFragment: "no *_spec.yaml files",
		},
	}

	for _, c := range testCases {
		t.Run(c.title, func(t *testing.T) {
			args := []string{"helm-spec", c.path}
			_, _, err := testRun(t, args)
			if c.valid {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, c.errMsgFragment)
			}
		})
	}
}

func TestTestCommandExecutesTestRunner(t *testing.T) {
	specDir, err := filepath.Abs("./testdata/specs")
	assert.NoError(t, err)
	args := []string{"helm-spec", specDir}
	runner, _, err := testRun(t, args)
	assert.NoError(t, err)
	assert.True(t, runner.HasRun)
	assert.ElementsMatch(t, []string{filepath.Join(specDir, "example_spec.yaml")}, runner.SpecFiles)
}

func TestTestCommandsWritesTestReport(t *testing.T) {
	specDir, err := filepath.Abs("./testdata/specs")
	assert.NoError(t, err)
	args := []string{"helm-spec", specDir}
	_, out, err := testRun(t, args)
	assert.NoError(t, err)
	assert.Equal(t, "success!", out.String())
}
