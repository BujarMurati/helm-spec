package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bujarmurati/helm-spec/internal/helmspec"
	"github.com/bujarmurati/helm-spec/internal/testreport"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

type errCapture struct {
	Ctx *cli.Context
	Err error
}

type mockTestRunner struct {
	Result    helmspec.TestSuiteResult
	SpecFiles []string
	HasRun    bool
}

func (m *mockTestRunner) Run(specFiles []string) (r helmspec.TestSuiteResult, err error) {
	m.SpecFiles = specFiles
	m.HasRun = true
	return m.Result, nil
}

type mockTestReporter struct {
	Settings testreport.TestReportSettings
}

func (m *mockTestReporter) Report(_ helmspec.TestSuiteResult, settings testreport.TestReportSettings) (string, error) {
	m.Settings = settings
	return "output", nil
}

type testCLISettings struct {
	cliSettings
	*errCapture
}

func newTestCLISettings() testCLISettings {
	e := &errCapture{}
	runner := &mockTestRunner{}
	runner.Result.Succeeded = true
	return testCLISettings{
		cliSettings{
			Reader:    strings.NewReader(""),
			Writer:    &strings.Builder{},
			ErrWriter: &strings.Builder{},
			ExitErrHandler: func(cCtx *cli.Context, err error) {
				e.Ctx = cCtx
				e.Err = err
			},
			TestRunner:   runner,
			TestReporter: &mockTestReporter{},
		},
		e,
	}
}

// runs the CLI with args. Outputs and errors  are captured
func testRun(t *testing.T, args []string) (settings testCLISettings, err error) {
	t.Helper()
	settings = newTestCLISettings()
	app, err := createApp(settings.cliSettings)
	assert.NoError(t, err)
	return settings, app.Run(args)
}

// just a sanity check to fail fast if something is utterly broken
func TestSanity(t *testing.T) {
	args := []string{"helm-spec", "--help"}
	_, err := testRun(t, args)
	assert.NoError(t, err)
}

func TestAcceptsSpecDirArg(t *testing.T) {
	args := []string{"helm-spec", "./testdata/specs"}
	_, err := testRun(t, args)
	assert.NoError(t, err)
}

func TestValidatesSpecDirArg(t *testing.T) {
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
			_, err := testRun(t, args)
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
	settings, err := testRun(t, args)
	assert.NoError(t, err)
	assert.True(t, settings.TestRunner.(*mockTestRunner).HasRun)
	assert.ElementsMatch(t, []string{filepath.Join(specDir, "example_spec.yaml")}, settings.TestRunner.(*mockTestRunner).SpecFiles)
}

func TestExitCodeOnFailure(t *testing.T) {
	specDir, err := filepath.Abs("./testdata/specs")
	assert.NoError(t, err)
	args := []string{"helm-spec", specDir}
	settings := newTestCLISettings()
	settings.TestRunner.(*mockTestRunner).Result.Succeeded = false
	app, err := createApp(settings.cliSettings)
	assert.NoError(t, err)
	err = app.Run(args)
	assert.Error(t, err)
}

func TestReportPrettyWithColorByDefault(t *testing.T) {
	specDir, err := filepath.Abs("./testdata/specs")
	assert.NoError(t, err)
	args := []string{"helm-spec", specDir}
	settings, err := testRun(t, args)
	assert.NoError(t, err)
	reportSettings := settings.TestReporter.(*mockTestReporter).Settings
	assert.True(t, reportSettings.UseColor)
	assert.Equal(t, testreport.OutputFormatPretty, reportSettings.OutputFormat)
}

func TestDisableColorOutputViaFlag(t *testing.T) {
	specDir, err := filepath.Abs("./testdata/specs")
	assert.NoError(t, err)
	args := []string{"helm-spec", "--no-color", specDir}
	settings, err := testRun(t, args)
	assert.NoError(t, err)
	reportSettings := settings.TestReporter.(*mockTestReporter).Settings
	assert.False(t, reportSettings.UseColor)
}

func TestDisableColorOutputViaEnv(t *testing.T) {
	specDir, err := filepath.Abs("./testdata/specs")
	assert.NoError(t, err)
	args := []string{"helm-spec", specDir}
	for _, env := range []string{"NO_COLOR", "HELM_SPEC_NO_COLOR"} {
		t.Run(env, func(t *testing.T) {
			os.Setenv(env, "")
			defer os.Unsetenv(env)
			settings, err := testRun(t, args)
			assert.NoError(t, err)
			reportSettings := settings.TestReporter.(*mockTestReporter).Settings
			assert.False(t, reportSettings.UseColor)
		})
	}
}

func TestVerboseMode(t *testing.T) {
	specDir, err := filepath.Abs("./testdata/specs")
	assert.NoError(t, err)
	args := []string{"helm-spec", "--verbose", specDir}
	settings, err := testRun(t, args)
	assert.NoError(t, err)
	reportSettings := settings.TestReporter.(*mockTestReporter).Settings
	assert.True(t, reportSettings.Verbose)
}

func TestVersion(t *testing.T) {
	args := []string{"helm-spec", "--version"}
	version = "0.1.0"
	settings, err := testRun(t, args)
	assert.NoError(t, err)
	output := settings.cliSettings.Writer.(*strings.Builder).String()
	assert.Contains(t, output, version)
}
