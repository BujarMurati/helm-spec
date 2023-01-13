package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/bujarmurati/helm-spec/internal/helmspec"
	"github.com/bujarmurati/helm-spec/internal/testreport"
	"github.com/urfave/cli/v2"
)

const (
	errFailedToGetAbsolutePathTemplate = "failed to get absolute path of: %v"
	errNotADirectoryTemplate           = "%v is not a directory"
	errNoSpecFilesFoundTemplate        = "no %v files found in %v"
	specFileGlobPattern                = "*_spec.yaml"
	defaultSpecDir                     = "./specs"
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
func validateSpecDirPath(value string) (err error) {
	absPath, err := filepath.Abs(value)
	if err != nil {
		return fmt.Errorf(errFailedToGetAbsolutePathTemplate, value)
	}
	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("spec directory `%v` does not seem to exist: %w", value, err)
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

func validateOutputFormat(o string) (err error) {
	for _, f := range testreport.AllowedOutputFormats {
		if o == f {
			return nil
		}
	}
	return fmt.Errorf("output format must be one of `%v`", testreport.AllowedOutputFormats)
}

func createApp(settings cliSettings) (app *cli.App, err error) {
	app = &cli.App{
		Name:            "helm-spec",
		Usage:           "automated tests for helm charts",
		ArgsUsage:       "<spec directory (default: \"./specs\")>",
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output-format",
				Aliases: []string{"o"},
				Value:   "pretty",
				Usage:   "output format for the report, one of \"pretty\"|\"yaml\"",
			},
		},
		Action: func(cCtx *cli.Context) (err error) {
			var specDir string
			if cCtx.Args().Present() {
				specDir = cCtx.Args().First()
			} else {
				specDir = defaultSpecDir
			}
			if err = validateSpecDirPath(specDir); err != nil {
				return err
			}
			outputFormat := cCtx.String("output-format")
			if outputFormat == "" {
				outputFormat = testreport.OutputFormatPretty
			}
			if err = validateOutputFormat(outputFormat); err != nil {
				return err
			}

			specFiles, err := filepath.Glob(filepath.Join(specDir, specFileGlobPattern))
			if err != nil {
				return err
			}
			result, err := settings.TestRunner.Run(specFiles)
			if err != nil {
				return err
			}
			reportSettings := testreport.TestReportSettings{OutputFormat: outputFormat, UseColor: true}
			report, err := testreport.HelmTestReporter{Result: result}.Report(reportSettings)
			if err != nil {
				return err
			}
			fmt.Fprint(settings.Writer, report)
			if !result.Succeeded {
				return errors.New("test suite failed")
			}
			return nil
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
		log.Fatal(err.Error())
	}
	if err = app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
	}
}
