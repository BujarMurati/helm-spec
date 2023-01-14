package testreport

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/bujarmurati/helm-spec/internal/helmspec"
	"github.com/pterm/pterm"
)

const pass = "\u2705 passed"
const fail = "\u274c failed"

// returns `passed` or `failed`
func passOrFailNoColor(succeeded bool) string {
	if succeeded {
		return pass
	} else {
		return fail
	}
}

func passOrFail(succeeded bool) string {
	text := passOrFailNoColor(succeeded)
	if succeeded {
		return pterm.FgGreen.Sprint(text)
	} else {
		return pterm.FgRed.Sprint(text)
	}
}

const assertionTmpl = `        {{ passOrFail .Succeeded }} - {{ .Assertion.Description }}
		{{- if (not .Succeeded) }}
		query: 
		    {{ .Assertion.Query }}
		want:
			{{ .Assertion.ExpectedResult }}
		got:
			{{ .ActualResult }}
		{{- end }}`

func prettyAssertionReport(result helmspec.AssertionResult, settings TestReportSettings) (string, error) {
	buf := &strings.Builder{}
	funcMap := make(map[string]any)
	if settings.UseColor {
		funcMap["passOrFail"] = passOrFail
	} else {
		funcMap["passOrFail"] = passOrFailNoColor
	}
	tpl, err := template.New("assertion").Funcs(funcMap).Parse(assertionTmpl)
	if err != nil {
		return "", err
	}
	err = tpl.Execute(buf, result)

	return buf.String(), err
}

const testCaseTmpl = "    {{ passOrFail .Succeeded }} - {{ .Title }}"

func prettyTestCaseReport(result helmspec.TestCaseResult, settings TestReportSettings) (string, error) {
	buf := &strings.Builder{}
	funcMap := make(map[string]any)
	if settings.UseColor {
		funcMap["passOrFail"] = passOrFail
	} else {
		funcMap["passOrFail"] = passOrFailNoColor
	}
	tpl, err := template.New("testCase").Funcs(funcMap).Parse(testCaseTmpl)
	if err != nil {
		return "", err
	}
	err = tpl.Execute(buf, result)
	report := buf.String()
	for _, a := range result.AssertionResults {
		assertionReport, err := prettyAssertionReport(a, settings)
		if err != nil {
			return report, nil
		}
		report += "\n"
		report += assertionReport
	}
	if !result.Succeeded && settings.Verbose {
		manifestLines := strings.Split(result.Manifest, "\n")
		for idx, line := range manifestLines {
			manifestLines[idx] = strings.Repeat(" ", 12) + line
		}
		report += "\n" + strings.Repeat(" ", 8) + "\U0001f4a1 manifest:\n"
		report += strings.Join(manifestLines, "\n")
		if result.Error != nil {
			report += "\n" + strings.Repeat(" ", 8) + "\u26a0\ufe0f error:\n"
			report += strings.Repeat(" ", 12) + result.Error.Error() + "\n"
		}
	}
	return report, err
}

func prettySpecReport(result helmspec.SpecResult, settings TestReportSettings) (string, error) {
	var status string
	if settings.UseColor {
		status = passOrFail(result.Succeeded)
	} else {
		status = passOrFailNoColor(result.Succeeded)
	}

	report := fmt.Sprintf("%v - %v", status, result.Title)
	for _, c := range result.TestCaseResults {
		testCaseReport, err := prettyTestCaseReport(c, settings)
		if err != nil {
			return report, nil
		}
		report += "\n"
		report += testCaseReport
	}
	report += "\n"
	return report, nil
}
