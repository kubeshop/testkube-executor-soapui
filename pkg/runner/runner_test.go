package runner

import (
	"log"
	"testing"
	"time"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {

	t.Run("runner should run test based on execution data", func(t *testing.T) {
		// given
		runner := NewRunner()
		execution := testkube.NewQueuedExecution()
		execution.Content = testkube.NewStringTestContent("hello I'm test content")

		// when
		result, err := runner.Run(*execution)

		// then
		assert.NoError(t, err)
		assert.Equal(t, result.Status, testkube.ExecutionStatusPassed)
	})

	t.Run("Run GET to Swagger petstore", func(t *testing.T) {
		f := MockFetcher{}

		f.FetchFn = func(content *testkube.TestContent) (path string, err error) {
			return "./test_projects/REST-Project-2-soapui-project.xml", nil
		}

		runner := SoapUIRunner{
			Fetcher: f,
		}
		execution := testkube.Execution{
			Id:            "get_petstore",
			TestName:      "Get Petstore",
			TestNamespace: "swagger",
			TestType:      "soapui/rest/xml", // TODO update
			Name:          "Testing GET",
			Envs:          map[string]string{},
			Args:          []string{},
			Params:        map[string]string{},
			ParamsFile:    "",
			Content: &testkube.TestContent{
				Type_:      "type",
				Repository: nil, // *Repository `json:"repository,omitempty"`
				Data:       "",
				Uri:        "",
			},
			StartTime:       time.Time{},
			EndTime:         time.Time{},
			Duration:        "1min",
			ExecutionResult: nil,
			Labels:          map[string]string{},
		}
		// execution.Content = testkube.NewStringTestContent("hello I'm test content")

		runner.Run(execution)
	})

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
