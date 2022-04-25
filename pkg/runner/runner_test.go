package runner

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	testXML := "./example/REST-Project-1-soapui-project.xml"

	t.Run("Successful test", func(t *testing.T) {
		t.Parallel()

		f := MockFetcher{}

		f.FetchFn = func(content *testkube.TestContent) (path string, err error) {
			return testXML, nil
		}

		file, err := createSuccessfulScript()
		assert.NoError(t, err)
		defer file.Close()
		runner := SoapUIRunner{
			Fetcher:        f,
			SoapUIExecPath: file.Name(),
		}
		execution := testkube.Execution{
			Id:              "get_petstore",
			TestName:        "Get Petstore",
			TestNamespace:   "petstore",
			TestType:        "soapui/xml",
			Name:            "Testing GET",
			Envs:            map[string]string{},
			Args:            []string{"-c 'TestCase 1'"},
			Params:          map[string]string{},
			ParamsFile:      "",
			Content:         &testkube.TestContent{},
			StartTime:       time.Time{},
			EndTime:         time.Time{},
			Duration:        "0s",
			ExecutionResult: nil,
			Labels:          map[string]string{},
		}

		res, err := runner.Run(execution)
		assert.NoError(t, err)
		assert.Equal(t, res.Status, testkube.ExecutionStatusPassed)
	})

	t.Run("Failing test", func(t *testing.T) {
		t.Parallel()

		f := MockFetcher{}

		f.FetchFn = func(content *testkube.TestContent) (path string, err error) {
			return testXML, nil
		}

		file, err := createFailingScript()
		assert.NoError(t, err)
		defer file.Close()
		runner := SoapUIRunner{
			Fetcher:        f,
			SoapUIExecPath: file.Name(),
		}
		execution := testkube.Execution{
			Id:            "get_petstore",
			TestName:      "Get Petstore",
			TestNamespace: "petstore",
			TestType:      "soapui/xml",
			Name:          "Testing GET",
			Args:          []string{"-c 'TestCase 1'"},
			Content:       &testkube.TestContent{},
		}

		res, err := runner.Run(execution)
		assert.NoError(t, err)
		assert.Equal(t, res.Status, testkube.ExecutionStatusFailed)
	})
}

func createSuccessfulScript() (*os.File, error) {
	file, err := ioutil.TempFile("", "successful_script")
	if err != nil {
		return nil, err
	}

	_, err = file.WriteString("/bin/sh\nexit 0\n")
	if err != nil {
		return nil, err
	}

	return file, nil
}

func createFailingScript() (*os.File, error) {
	file, err := ioutil.TempFile("", "failing_script")
	if err != nil {
		return nil, err
	}

	_, err = file.WriteString("/bin/sh\nexit 1\n")
	if err != nil {
		return nil, err
	}

	return file, nil
}

type MockFetcher struct {
	FetchFn        func(content *testkube.TestContent) (path string, err error)
	FetchStringFn  func(str string) (path string, err error)
	FetchURIFn     func(uri string) (path string, err error)
	FetchGitDirFn  func(repo *testkube.Repository) (path string, err error)
	FetchGitFileFn func(repo *testkube.Repository) (path string, err error)
}

func (f MockFetcher) Fetch(content *testkube.TestContent) (path string, err error) {
	if f.FetchFn == nil {
		log.Fatal("not implemented")
	}
	return f.FetchFn(content)
}

func (f MockFetcher) FetchString(str string) (path string, err error) {
	if f.FetchStringFn == nil {
		log.Fatal("not implemented")
	}
	return f.FetchStringFn(str)
}

func (f MockFetcher) FetchURI(str string) (path string, err error) {
	if f.FetchURIFn == nil {
		log.Fatal("not implemented")
	}
	return f.FetchURIFn(str)
}

func (f MockFetcher) FetchGitDir(repo *testkube.Repository) (path string, err error) {
	if f.FetchGitDirFn == nil {
		log.Fatal("not implemented")
	}
	return f.FetchGitDir(repo)
}

func (f MockFetcher) FetchGitFile(repo *testkube.Repository) (path string, err error) {
	if f.FetchGitFileFn == nil {
		log.Fatal("not implemented")
	}
	return f.FetchGitFileFn(repo)
}
