package testreport

import (
	"strings"
	"testing"

	"github.com/bujarmurati/helm-spec/internal/helmspec"
	"github.com/stretchr/testify/assert"
)

func TestPrettySuccessfulAssertionReport(t *testing.T) {
	settings := TestReportSettings{
		UseColor:     true,
		OutputFormat: "pretty",
	}
	res := helmspec.AssertionResult{
		Succeeded:    true,
		ActualResult: "bar",
		Assertion: helmspec.Assertion{
			Description:    "deployment name should be `bar`",
			ExpectedResult: "bar",
			Query:          "select(.kind==\"Deployment\") | .metadata.name",
		},
	}
	output, err := prettyAssertionReport(res, settings)
	t.Cleanup(func() {
		t.Logf("\n%v", output)
	})
	assert.NoError(t, err)
	lines := strings.Split(output, "\n")
	assert.Equal(t, 1, len(lines))
	heading := strings.TrimSpace(lines[0])
	assert.Contains(t, heading, pass)
}

func TestPrettyFailedAssertionReport(t *testing.T) {
	settings := TestReportSettings{
		UseColor:     true,
		OutputFormat: "pretty",
	}
	res := helmspec.AssertionResult{
		Succeeded:    false,
		ActualResult: "foo",
		Assertion: helmspec.Assertion{
			Description:    "deployment name should be `bar`",
			ExpectedResult: "bar",
			Query:          "select(.kind==\"Deployment\") | .metadata.name",
		},
	}
	output, err := prettyAssertionReport(res, settings)
	t.Cleanup(func() {
		t.Logf("\n%v", output)
	})
	assert.NoError(t, err)
	lines := strings.Split(output, "\n")
	assert.Equal(t, 7, len(lines))
	heading := strings.TrimSpace(lines[0])
	assert.Contains(t, heading, fail)
	query := strings.TrimSpace(lines[2])
	assert.Contains(t, query, res.Assertion.Query)
	want := strings.TrimSpace(lines[4])
	assert.Contains(t, want, res.Assertion.ExpectedResult)
	got := strings.TrimSpace(lines[6])
	assert.Contains(t, got, res.ActualResult)
}

func TestPrettySuccessfulTestCaseReport(t *testing.T) {
	settings := TestReportSettings{
		UseColor:     true,
		OutputFormat: "pretty",
	}
	res := helmspec.TestCaseResult{
		Title:     "deployment",
		Succeeded: true,
		AssertionResults: []helmspec.AssertionResult{{
			Succeeded:    true,
			ActualResult: "bar",
			Assertion: helmspec.Assertion{
				Description:    "deployment name should be `bar`",
				ExpectedResult: "bar",
				Query:          "select(.kind==\"Deployment\") | .metadata.name",
			}},
		}}
	output, err := prettyTestCaseReport(res, settings)
	t.Cleanup(func() {
		t.Logf("\n%v", output)
	})
	assert.NoError(t, err)
	lines := strings.Split(output, "\n")
	assert.Equal(t, 2, len(lines))
	heading := strings.TrimSpace(lines[0])
	assert.Contains(t, heading, pass)
	assertion := strings.TrimSpace(lines[1])
	assert.Contains(t, assertion, pass)
}

func TestPrettyFailedTestCaseReport(t *testing.T) {
	settings := TestReportSettings{
		UseColor:     true,
		OutputFormat: "pretty",
	}
	res := helmspec.TestCaseResult{
		Title:     "deployment",
		Succeeded: true,
		AssertionResults: []helmspec.AssertionResult{
			{
				Succeeded:    true,
				ActualResult: "bar",
				Assertion: helmspec.Assertion{
					Description:    "deployment name should be `bar`",
					ExpectedResult: "bar",
					Query:          "select(.kind==\"Deployment\") | .metadata.name",
				},
			},
			{
				Succeeded:    false,
				ActualResult: "not bar",
				Assertion: helmspec.Assertion{
					Description:    "deployment should have label `foo: bar`",
					ExpectedResult: "bar",
					Query:          "select(.kind==\"Deployment\") | .labels.foo",
				},
			},
		}}
	output, err := prettyTestCaseReport(res, settings)
	t.Cleanup(func() {
		t.Logf("\n%v", output)
	})
	assert.NoError(t, err)
	lines := strings.Split(output, "\n")
	assert.Equal(t, 9, len(lines))
	heading := strings.TrimSpace(lines[0])
	assert.Contains(t, heading, pass)
	successfulAssertion := strings.TrimSpace(lines[1])
	assert.Contains(t, successfulAssertion, pass)
	failedAssertions := strings.TrimSpace(lines[2])
	assert.Contains(t, failedAssertions, fail)
}
