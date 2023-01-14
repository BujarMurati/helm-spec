package testreport

import (
	"fmt"
	"strings"

	"github.com/bujarmurati/helm-spec/internal/helmspec"
	"sigs.k8s.io/yaml"
)

const OutputFormatYAML = "yaml"
const OutputFormatPretty = "pretty"

var AllowedOutputFormats = [...]string{OutputFormatYAML, OutputFormatPretty}

type TestReportSettings struct {
	OutputFormat string
	UseColor     bool
	Verbose      bool
}

type TestReporter interface {
	Report(result helmspec.TestSuiteResult, settings TestReportSettings) (string, error)
}

type HelmTestReporter struct{}

func (r HelmTestReporter) Report(result helmspec.TestSuiteResult, settings TestReportSettings) (output string, err error) {
	switch settings.OutputFormat {
	case OutputFormatYAML:
		content, err := yaml.Marshal(result)
		return string(content), err
	case OutputFormatPretty:
		var status string
		if settings.UseColor {
			status = passOrFail(result.Succeeded)
		} else {
			status = passOrFailNoColor(result.Succeeded)
		}
		report := fmt.Sprintf("testsuite %v\n", status)
		report += strings.Repeat("=", 32) + " details " + strings.Repeat("=", 32)
		report += "\n\n"
		for _, s := range result.SpecResults {
			r, err := prettySpecReport(s, settings)
			if err != nil {
				return report, err
			}
			report += r
		}
		return report, nil
	default:
		return "", fmt.Errorf("unsupported output format `%v`", settings.OutputFormat)
	}
}
