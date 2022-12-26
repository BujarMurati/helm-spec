package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

type errCapture struct {
	Ctx *cli.Context
	Err error
}

func testSettings(e *errCapture) cliSettings {

	return cliSettings{
		Reader:    strings.NewReader(""),
		Writer:    &strings.Builder{},
		ErrWriter: &strings.Builder{},
		ExitErrHandler: func(cCtx *cli.Context, err error) {
			e.Ctx = cCtx
			e.Err = err
		},
	}
}

// runs the CLI with args. Outputs and errors  are captured
func testRun(t *testing.T, args []string) (settings cliSettings, e *errCapture, err error) {
	t.Helper()
	e = &errCapture{}
	settings = testSettings(e)
	app, err := createApp(settings)
	assert.NoError(t, err)
	return settings, e, app.Run(args)
}

func TestHasTestCommand(t *testing.T) {
	args := []string{"helm-spec", "test"}
	_, _, err := testRun(t, args)
	assert.NoError(t, err)
}

func TestAcceptsTestSuitePathFlag(t *testing.T) {
	args := []string{"helm-spec", "test", "--testsuite=./testdata/specs"}
	_, _, err := testRun(t, args)
	assert.NoError(t, err)
}

func TestValidatesTestSuitePath(t *testing.T) {
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
	}

	for _, c := range testCases {
		t.Run(c.title, func(t *testing.T) {
			args := []string{"helm-spec", "test", fmt.Sprintf("--t=%v", c.path)}
			_, _, err := testRun(t, args)
			if c.valid {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, c.errMsgFragment)
			}
		})
	}
}
