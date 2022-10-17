package runner

import (
	"errors"
	"os"
	"testing"

	"github.com/kubeshop/testkube-executor-soapui/pkg/mock"
	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	testXML := "./example/REST-Project-1-soapui-project.xml"
	f := mock.Fetcher{}
	f.FetchFn = func(content *testkube.TestContent) (path string, err error) {
		return testXML, nil
	}

	e := testkube.Execution{
		Id:            "get_petstore",
		TestName:      "Get Petstore",
		TestNamespace: "petstore",
		TestType:      "soapui/xml",
		Name:          "Testing GET",
		Args:          []string{"-c 'TestCase 1'"},
		Content:       &testkube.TestContent{},
	}

	tests := []struct {
		name            string
		scraper         func(id string, directories []string) error
		testFileCreator func() (*os.File, error)
		execution       testkube.Execution
		expectedError   string
		expectedStatus  testkube.ExecutionStatus
	}{
		{
			name:            "Successful test, successful scraper",
			scraper:         func(id string, directories []string) error { return nil },
			testFileCreator: createSuccessfulScript,
			execution:       e,
			expectedError:   "",
			expectedStatus:  *testkube.ExecutionStatusPassed,
		},
		{
			name:            "Failing test, successful scraper",
			scraper:         func(id string, directories []string) error { return nil },
			testFileCreator: createFailingScript,
			execution:       e,
			expectedError:   "",
			expectedStatus:  *testkube.ExecutionStatusFailed,
		},
		{
			name:            "Successful test, failing scraper",
			scraper:         func(id string, directories []string) error { return errors.New("Scraping failed") },
			testFileCreator: createSuccessfulScript,
			execution:       e,
			expectedError:   "failed getting artifacts: Scraping failed",
			expectedStatus:  *testkube.ExecutionStatusPassed,
		},
		{
			name:            "Failing test, failing scraper",
			scraper:         func(id string, directories []string) error { return errors.New("Scraping failed") },
			testFileCreator: createFailingScript,
			execution:       e,
			expectedError:   "failed getting artifacts: Scraping failed",
			expectedStatus:  *testkube.ExecutionStatusFailed,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			s := mock.Scraper{}
			s.ScrapeFn = test.scraper

			file, err := test.testFileCreator()
			assert.NoError(t, err)
			defer file.Close()

			runner := SoapUIRunner{
				Fetcher:        f,
				SoapUIExecPath: file.Name(),
				Scraper:        s,
			}

			res, err := runner.Run(test.execution)
			assert.EqualError(t, err, test.expectedError)
			assert.Equal(t, test.expectedStatus, *res.Status)
		})
	}
}

func createSuccessfulScript() (*os.File, error) {
	file, err := os.CreateTemp("", "successful_script")
	if err != nil {
		return nil, err
	}

	_, err = file.WriteString("#!/bin/sh\nexit 0\n")
	if err != nil {
		return nil, err
	}

	return file, nil
}

func createFailingScript() (*os.File, error) {
	file, err := os.CreateTemp("", "failing_script")
	if err != nil {
		return nil, err
	}

	_, err = file.WriteString("#!/bin/sh\nexit 1\n")
	if err != nil {
		return nil, err
	}

	return file, nil
}
