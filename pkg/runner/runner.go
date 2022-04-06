package runner

import (
	"errors"
	"os"
	"os/exec"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/executor/content"
	"github.com/kubeshop/testkube/pkg/executor/output"
)

func NewRunner() *SoapUIRunner {
	return &SoapUIRunner{
		Fetcher: content.NewFetcher(""),
	}
}

// SoapUIRunner runs SoapUI tests
type SoapUIRunner struct {
	Fetcher content.ContentFetcher
}

// Run executes the test and returns the test results
func (r *SoapUIRunner) Run(execution testkube.Execution) (result testkube.ExecutionResult, err error) {
	testFile, err := r.Fetcher.Fetch(execution.Content)
	if err != nil {
		return result, err
	}

	output.PrintEvent("created content path", testFile)
	setUpEnvironment(testFile)

	if execution.Content.IsDir() {
		return testkube.ExecutionResult{}, errors.New("SoapUI executor only tests one project per execution, a directory of projects was given")
	}

	return runSoapUI(), nil
}

// setTestFile sets up the COMMAND_LINE environment variable to
// point to the test file path
func setUpEnvironment(testFilePath string) {
	os.Setenv("COMMAND_LINE", testFilePath)
}

// runSoapUI runs SoapUI tests and returns the output
func runSoapUI() testkube.ExecutionResult {
	output, err := exec.Command("/bin/sh", "/usr/local/SmartBear/EntryPoint.sh").Output()
	if err != nil {
		return testkube.ExecutionResult{
			Status:       testkube.ExecutionStatusFailed,
			ErrorMessage: err.Error(),
		}
	}

	return testkube.ExecutionResult{
		Status: testkube.ExecutionStatusPassed,
		Output: string(output),
	}
}
