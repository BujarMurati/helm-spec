package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	helmspec "github.com/bujarmurati/helm-spec"
	"github.com/urfave/cli/v2"
)

const (
	errFailedToGetAbsolutePathTemplate = "failed to get absolute path of: %v"
	errNotADirectoryTemplate           = "%v is not a directory"
	errNoSpecFilesFoundTemplate        = "no %v files found in %v"
	specFileGlobPattern                = "*_spec.yaml"
)

type cliSettings struct {
	Reader         io.Reader
	Writer         io.Writer
	ErrWriter      io.Writer
	ExitErrHandler cli.ExitErrHandlerFunc
	TestRunner     helmspec.TestRunner
}

var defaultSettings = cliSettings{
	Reader:         os.Stdin,
	Writer:         os.Stdout,
	ErrWriter:      os.Stderr,
	ExitErrHandler: nil,
	TestRunner:     &helmspec.HelmTestRunner{},
}

// validates that the path is an existing directory containing spec files
func validateTestSuitePath(cCtx *cli.Context, value string) (err error) {
	absPath, err := filepath.Abs(value)
	if err != nil {
		return fmt.Errorf(errFailedToGetAbsolutePathTemplate, value)
	}
	info, err := os.Stat(absPath)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf(errNotADirectoryTemplate, absPath)
	}
	specFiles, err := filepath.Glob(filepath.Join(absPath, specFileGlobPattern))
	if err != nil {
		return err
	}
	if len(specFiles) == 0 {
		return fmt.Errorf(errNoSpecFilesFoundTemplate, specFileGlobPattern, absPath)
	}
	return nil
}

func createApp(settings cliSettings) (app *cli.App, err error) {
	app = &cli.App{
		Name:           "helm-spec",
		DefaultCommand: "test",
		Commands: []*cli.Command{
			{
				Name: "test",
				Flags: []cli.Flag{
					&cli.PathFlag{
						Usage:       "path to a directory containing *_spec.yaml files",
						Name:        "testsuite",
						Aliases:     []string{"t"},
						Value:       "./specs",
						DefaultText: "./specs",
						Required:    false,
						TakesFile:   false,
						Action:      validateTestSuitePath,
					},
				},
				Action: func(cCtx *cli.Context) (err error) {
					specDir := cCtx.Path("testsuite")
					specFiles, err := filepath.Glob(filepath.Join(specDir, specFileGlobPattern))
					if err != nil {
						return err
					}
					reporter, err := settings.TestRunner.Run(specFiles)
					if err != nil {
						return err
					}
					report, err := reporter.Report("yaml")
					if err != nil {
						return err
					}
					fmt.Fprint(settings.Writer, report)
					return nil
				},
			},
		},
		Reader:         settings.Reader,
		Writer:         settings.Writer,
		ErrWriter:      settings.ErrWriter,
		ExitErrHandler: settings.ExitErrHandler,
	}
	return app, nil
}

func main() {
	app, err := createApp(defaultSettings)
	if err != nil {
		panic("failed to initialize app")
	}
	app.Run(os.Args)
}
