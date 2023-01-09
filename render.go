package helmspec

import (
	"os/exec"
	"strings"
)

// inputs for rendering a the chart with `helm template`
type RenderInstructions struct {
	// the release name to pass to `helm template`
	ReleaseName string `json:"releaseName"`
	// the release namespace to pass to `helm template`
	Namespace string `json:"namespace"`
	// all user-supplied values in one inline yaml document
	Values string `json´:"values"`
	// extra arguments passed through to the helm CLI, i.e. ["--set-file", "foo=foo.txt"]
	ExtraArgs []string `json:"extraArgs"`
	// require rendering to fail for the test to pass
	ShouldFailToRender bool `json:"shouldFailToRender"`
}

// runs helm dependency build and helm template, returning the
// rendered manifest or error
func (r RenderInstructions) Execute(chartPath string) (manifest string, err error) {
	helmDepBuildArgs := []string{"dependency", "build", chartPath}
	helmDepBuild := exec.Command("helm", helmDepBuildArgs...)
	err = helmDepBuild.Run()
	if err != nil {
		return "", err
	}

	helmTemplateArgs := []string{"template"}
	if r.ReleaseName != "" {
		helmTemplateArgs = append(helmTemplateArgs, r.ReleaseName)
	}
	helmTemplateArgs = append(helmTemplateArgs, chartPath)
	if r.Namespace != "" {
		helmTemplateArgs = append(helmTemplateArgs, "-n", r.Namespace)
	}
	helmTemplateArgs = append(helmTemplateArgs, r.ExtraArgs...)
	helmTemplateArgs = append(helmTemplateArgs, "-f", "-")
	helmTemplate := exec.Command("helm", helmTemplateArgs...)
	helmTemplate.Stdin = strings.NewReader(r.Values)
	out := &strings.Builder{}
	helmTemplate.Stdout = out
	err = helmTemplate.Run()
	return out.String(), err
}
