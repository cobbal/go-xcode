package xcodeproj

import (
	"encoding/json"
	"fmt"

	"github.com/bitrise-io/go-utils/command/rubyscript"
)

// ScanBlueprintIdentifier ...
func ScanBlueprintIdentifier(projectPth, scheme string) (string, error) {
	runner := rubyscript.New(bluePrintScriptContent)

	bundleInstallCmd, err := runner.BundleInstallCommand(targetInfoGemfileContent, "")
	if err != nil {
		return "", fmt.Errorf("failed to create bundle install command, error: %s", err)
	}

	if out, err := bundleInstallCmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
		return "", fmt.Errorf("bundle install failed, output: %s, error: %s", out, err)
	}

	runCmd, err := runner.RunScriptCommand()
	if err != nil {
		return "", fmt.Errorf("failed to create script runner command, error: %s", err)
	}

	envsToAppend := []string{
		"PROEJECTPATH=" + projectPth,
		"SCHEME_NAME=" + scheme,
	}
	envs := append(runCmd.GetCmd().Env, envsToAppend...)

	runCmd.SetEnvs(envs...)

	out, err := runCmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run scheme configuration analyzer script, output: %s, error: %s", out, err)
	}

	type OutputModel struct {
		Data  string `json:"data"`
		Error string `json:"error"`
	}
	var output OutputModel
	if err := json.Unmarshal([]byte(out), &output); err != nil {
		out = clearRubyScriptOutput(out)
		if err := json.Unmarshal([]byte(out), &output); err != nil {
			return "", fmt.Errorf("failed to unmarshal output: %s", out)
		}
	}

	if output.Error != "" {
		return "", fmt.Errorf("failed to get provisioning profile - bundle id mapping, error: %s", output.Error)
	}

	return output.Data, nil
}
