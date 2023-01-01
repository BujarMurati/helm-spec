package helmspec

import (
	"fmt"
	"strings"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"sigs.k8s.io/yaml"
)

const OutputFormatYAML = "yaml"
const OutputFormatPretty = "pretty"
const treeIndentationLevel = 4

var treeVerticalString = pterm.DefaultTree.TreeStyle.Sprint(pterm.DefaultTree.VerticalString)

var AllowedOutputFormats = [...]string{OutputFormatYAML, OutputFormatPretty}

type HelmTestReporter struct {
	Result TestSuiteResult
}

// returns color-coded `passed` or `failed`
func passOrFail(succeeded bool) string {
	if succeeded {
		return pterm.FgGreen.Sprint("passed")
	} else {
		return pterm.FgRed.Sprint("failed")
	}
}

func treeIndent(content string, treeLevel int, verticalContinuationsAtLevels []int, verticalString string) string {
	lines := strings.Split(content, "\n")
	spaceCount := (treeLevel + 1) * treeIndentationLevel
	prefix := strings.Repeat(" ", spaceCount)
	for _, l := range verticalContinuationsAtLevels {
		indexToReplace := l * treeIndentationLevel
		prefix = prefix[:indexToReplace] + verticalString + prefix[indexToReplace+1:]
	}
	indentedLines := []string{}
	for _, line := range lines {
		indentedLines = append(indentedLines, prefix+line)
	}
	return strings.Join(indentedLines, "\n")
}

func (r HelmTestReporter) Report(outputFormat string) (output string, err error) {
	switch {
	case outputFormat == OutputFormatYAML:
		content, err := yaml.Marshal(r.Result)
		return string(content), err
	case outputFormat == OutputFormatPretty:
		resultTree := pterm.LeveledList{}
		for specIndex, specResult := range r.Result.SpecResults {
			resultTree = append(resultTree, pterm.LeveledListItem{
				Level: 0,
				Text:  fmt.Sprintf("spec `%v`: %v", specResult.Title, passOrFail(specResult.Succeeded)),
			})
			for testCaseIndex, testCaseResult := range specResult.TestCaseResults {
				resultTree = append(resultTree, pterm.LeveledListItem{
					Level: 1,
					Text:  fmt.Sprintf("testCase `%v`: %v", testCaseResult.Title, passOrFail(testCaseResult.Succeeded)),
				})
				for _, assertionResult := range testCaseResult.AssertionResults {
					resultTree = append(resultTree, pterm.LeveledListItem{
						Level: 2,
						Text:  fmt.Sprintf("assertion `%v`: %v", assertionResult.Assertion.Description, passOrFail(assertionResult.Succeeded)),
					})
					if !assertionResult.Succeeded {
						resultTree = append(resultTree, pterm.LeveledListItem{
							Level: 3,
							Text:  fmt.Sprintf("query: `%v`", assertionResult.Assertion.Query),
						})
						resultTree = append(resultTree, pterm.LeveledListItem{
							Level: 3,
							Text:  fmt.Sprintf("expected: `%v`", assertionResult.Assertion.ExpectedResult),
						})
						resultTree = append(resultTree, pterm.LeveledListItem{
							Level: 3,
							Text:  fmt.Sprintf("actual: `%v`", assertionResult.ActualResult),
						})
					}
				}
				if testCaseResult.Error != nil {
					resultTree = append(resultTree, pterm.LeveledListItem{
						Level: 2,
						Text:  fmt.Sprintf("rendering error: `%v`", testCaseResult.Error),
					})
				}
				if !testCaseResult.Succeeded {
					verticalContinuationsAtLevels := []int{}
					isLastSpec := specIndex+1 == len(r.Result.SpecResults)
					if !isLastSpec {
						verticalContinuationsAtLevels = append(verticalContinuationsAtLevels, 0)
					}
					isLastTestCaseInSpec := testCaseIndex+1 == len(specResult.TestCaseResults)
					if !isLastTestCaseInSpec {
						verticalContinuationsAtLevels = append(verticalContinuationsAtLevels, 1)
					}
					manifestBox := pterm.DefaultBox.WithTitle(fmt.Sprintf("manifest rendered for spec `%v` test case `%v`", specResult.Title, testCaseResult.Title)).Sprint(testCaseResult.Manifest)
					resultTree = append(resultTree, pterm.LeveledListItem{
						Level: 2,
						// the pterm styling seems to "swallow"(?!) whitespace, thus`treeVerticalString+"   "`
						Text: "\n" + treeIndent(manifestBox, 2, verticalContinuationsAtLevels, treeVerticalString+"   "),
					})
				}
			}
		}

		root := putils.TreeFromLeveledList(resultTree)
		return pterm.DefaultTree.WithIndent(treeIndentationLevel).WithRoot(root).Srender()
	default:
		return "", fmt.Errorf("unsupported output format `%v`", outputFormat)
	}
}
