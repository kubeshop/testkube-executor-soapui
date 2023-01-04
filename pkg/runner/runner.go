package runner

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/executor"
	"github.com/kubeshop/testkube/pkg/executor/content"
	"github.com/kubeshop/testkube/pkg/executor/output"
	"github.com/kubeshop/testkube/pkg/executor/runner"
	"github.com/kubeshop/testkube/pkg/executor/scraper"
	"github.com/kubeshop/testkube/pkg/executor/secret"
)

type Params struct {
	Endpoint        string // RUNNER_ENDPOINT
	AccessKeyID     string // RUNNER_ACCESSKEYID
	SecretAccessKey string // RUNNER_SECRETACCESSKEY
	Location        string // RUNNER_LOCATION
	Token           string // RUNNER_TOKEN
	Ssl             bool   // RUNNER_SSL
	DataDir         string // RUNNER_DATADIR
}

const FailureMessage string = "finished with status [FAILED]"

// NewRunner creates a new SoapUIRunner
func NewRunner() (*SoapUIRunner, error) {
	var params Params
	err := envconfig.Process("runner", &params)
	if err != nil {
		return nil, err
	}

	return &SoapUIRunner{
		SoapUIExecPath: "/usr/local/SmartBear/EntryPoint.sh",
		SoapUILogsPath: "/home/soapui/.soapuios/logs",
		Fetcher:        content.NewFetcher(""),
		Scraper: scraper.NewMinioScraper(
			params.Endpoint,
			params.AccessKeyID,
			params.SecretAccessKey,
			params.Location,
			params.Token,
			params.Ssl,
		),
		DataDir: params.DataDir,
	}, nil
}

// SoapUIRunner runs SoapUI tests
type SoapUIRunner struct {
	SoapUIExecPath string
	SoapUILogsPath string
	Fetcher        content.ContentFetcher
	Scraper        scraper.Scraper
	DataDir        string
}

// Run executes the test and returns the test results
func (r *SoapUIRunner) Run(execution testkube.Execution) (result testkube.ExecutionResult, err error) {
	testFile, err := r.Fetcher.Fetch(execution.Content)
	if err != nil {
		return result, err
	}

	output.PrintEvent("created content path", testFile)
	setUpEnvironment(execution.Args, testFile)

	if execution.Content.IsDir() {
		return testkube.ExecutionResult{}, errors.New("SoapUI executor only tests one project per execution, a directory of projects was given")
	}

	output.PrintEvent("running SoapUI tests")
	result = r.runSoapUI(&execution)

	output.PrintEvent("scraping for log files")
	if err = r.Scraper.Scrape(execution.Id, []string{r.SoapUILogsPath}); err != nil {
		return result, fmt.Errorf("failed getting artifacts: %w", err)
	}

	return result, nil
}

// setUpEnvironment sets up the COMMAND_LINE environment variable to
// contain the incoming arguments and to point to the test file path
func setUpEnvironment(args []string, testFilePath string) {
	args = append(args, testFilePath)
	os.Setenv("COMMAND_LINE", strings.Join(args, " "))
}

// runSoapUI runs the SoapUI executable and returns the output
func (r *SoapUIRunner) runSoapUI(execution *testkube.Execution) testkube.ExecutionResult {

	envManager := secret.NewEnvManagerWithVars(execution.Variables)
	envManager.GetVars(envManager.Variables)
	for _, env := range envManager.Variables {
		os.Setenv(env.Name, env.Value)
	}

	runPath := ""
	if execution.Content.Repository != nil && execution.Content.Repository.WorkingDir != "" {
		runPath = filepath.Join(r.DataDir, "repo", execution.Content.Repository.WorkingDir)
	}

	output, err := executor.Run(runPath, "/bin/sh", envManager, r.SoapUIExecPath)
	output = envManager.Obfuscate(output)
	if err != nil {
		return testkube.ExecutionResult{
			Status:       testkube.ExecutionStatusFailed,
			ErrorMessage: err.Error(),
		}
	}
	if strings.Contains(string(output), FailureMessage) {
		return testkube.ExecutionResult{
			Status:       testkube.ExecutionStatusFailed,
			ErrorMessage: FailureMessage,
			Output:       string(output),
		}
	}

	return testkube.ExecutionResult{
		Status: testkube.ExecutionStatusPassed,
		Output: string(output),
	}
}

// GetType returns runner type
func (r *SoapUIRunner) GetType() runner.Type {
	return runner.TypeMain
}
